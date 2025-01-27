package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var (
	db          *sql.DB
	redisClient *redis.Client
)

var jwtSecret = []byte("your-secret-key-here") // Replace with a secure secret

func init() {
	// Initialize PostgreSQL
	var err error
	connStr := "user=kai password=Keeru5463 dbname=wira_rd sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	// Add CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			next.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/dashboard", AuthMiddleware(DashboardHandler)).Methods("GET")
	return r
}

func generateJWTToken(accID int) (string, error) {
	// Create claims with account ID and expiry
	claims := jwt.MapClaims{
		"acc_id": accID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginHandler handles username/password + 2FA
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	totpCode := r.FormValue("twoFA_code")

	// Get account from DB
	var accID int
	var storedPassword, secret2FA string
	err := db.QueryRow(`
        SELECT acc_id, encrypted_password, secretkey_2fa 
        FROM Account WHERE username = $1`, username).
		Scan(&accID, &storedPassword, &secret2FA)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify TOTP
	if !totp.Validate(totpCode, secret2FA) {
		http.Error(w, "Invalid 2FA code", http.StatusUnauthorized)
		return
	}

	// Create session token
	token, err := generateJWTToken(accID) // Now properly defined
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}
	expiry := time.Now().Add(24 * time.Hour)

	// Store session in DB
	_, err = db.Exec(`
        INSERT INTO Session (session_id, acc_id, expiry_datetime)
        VALUES ($1, $2, $3)`, token, accID, expiry)
	if err != nil {
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

// AuthMiddleware validates JWT and checks session in the database
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Validate token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		token := tokenParts[1]

		// Parse and validate JWT
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !parsedToken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get account ID from JWT
		accIDFloat, ok := claims["acc_id"].(float64)
		if !ok {
			http.Error(w, "Invalid account ID in token", http.StatusUnauthorized)
			return
		}
		accID := int(accIDFloat)

		// Verify session exists in database
		var dbExpiry time.Time
		err = db.QueryRow(`
            SELECT expiry_datetime 
            FROM Session 
            WHERE session_id = $1 AND acc_id = $2`,
			token, accID,
		).Scan(&dbExpiry)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Session not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Add account ID to context
		ctx := context.WithValue(r.Context(), "accID", accID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	classID := r.URL.Query().Get("class_id")
	if classID == "" {
		http.Error(w, "Missing 'class_id' parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		http.Error(w, "Invalid 'page' parameter", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
		return
	}
	searchQuery := r.URL.Query().Get("search") // Get search parameter

	// Check Redis cache first
	cacheKey := fmt.Sprintf("rankings:%s:%d:%d:%s", classID, page, limit, searchQuery)
	cachedData, err := redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// Log cache hit
		fmt.Println("Cache hit for key:", cacheKey)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	} else {
		// Log cache miss
		fmt.Println("Cache miss for key:", cacheKey)
	}

	// Query database
	query := `
        SELECT a.username, c.class_id, s.reward_score 
        FROM Scores s
        JOIN Character c ON s.char_id = c.char_id
        JOIN Account a ON c.acc_id = a.acc_id
        WHERE c.class_id = $1
        AND (a.username ILIKE '%' || $2 || '%')
        ORDER BY s.reward_score DESC
        LIMIT $3 OFFSET $4`

	rows, err := db.Query(query, classID, searchQuery, limit, (page-1)*limit)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process rows into JSON
	type Ranking struct {
		Username    string `json:"username"`
		ClassID     int    `json:"class_id"`
		RewardScore int    `json:"reward_score"`
	}
	var rankings []Ranking

	for rows.Next() {
		var r Ranking
		if err := rows.Scan(&r.Username, &r.ClassID, &r.RewardScore); err != nil {
			http.Error(w, "Failed to read row", http.StatusInternalServerError)
			return
		}
		rankings = append(rankings, r)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	// Convert to JSON
	jsonData, err := json.Marshal(rankings)
	if err != nil {
		http.Error(w, "JSON serialization failed", http.StatusInternalServerError)
		return
	}

	// Cache indefinitely
	err = redisClient.Set(context.Background(), cacheKey, jsonData, 0).Err()
	if err != nil {
		fmt.Println("Error caching data:", err)
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
