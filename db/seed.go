package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	connStr := "user=kai password=Keeru5463 dbname=wira_rd sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create test account with known credentials
	createTestAccount(db)

	// Generate 12,499 additional accounts
	for i := 1; i < 12500; i++ {
		username := gofakeit.Username() + strconv.Itoa(i)
		email := fmt.Sprintf("%s%d@example.com", gofakeit.Username(), i)

		// Generate real password hash
		password := fmt.Sprintf("%s123", username) // Password = username + "123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		secret2FA := gofakeit.UUID()

		var accID int
		err = db.QueryRow(`
            INSERT INTO Account (username, email, encrypted_password, secretkey_2fa)
            VALUES ($1, $2, $3, $4) RETURNING acc_id`,
			username, email, string(hashedPassword), secret2FA,
		).Scan(&accID)
		if err != nil {
			log.Fatal(err)
		}

		// Insert characters and scores (same as before)
		for classID := 1; classID <= 8; classID++ {
			var charID int
			err := db.QueryRow(`
                INSERT INTO Character (acc_id, class_id)
                VALUES ($1, $2) RETURNING char_id`,
				accID, classID,
			).Scan(&charID)
			if err != nil {
				log.Fatal(err)
			}

			score := rand.Intn(10000)
			_, err = db.Exec(`
                INSERT INTO Scores (char_id, reward_score)
                VALUES ($1, $2)`,
				charID, score,
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("Fake data inserted successfully!")
}

func createTestAccount(db *sql.DB) {
	// Known test credentials
	testAccount := struct {
		username  string
		email     string
		password  string
		secret2FA string
	}{
		username:  "testuser",
		email:     "test@example.com",
		password:  "testpassword",     // Raw password
		secret2FA: "JBSWY3DPEHPK3PXP", // TOTP secret for code 123456 (test value)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testAccount.password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash test password:", err)
	}

	// Insert test account
	var accID int
	err = db.QueryRow(`
        INSERT INTO Account (username, email, encrypted_password, secretkey_2fa)
        VALUES ($1, $2, $3, $4) RETURNING acc_id`,
		testAccount.username, testAccount.email, string(hashedPassword), testAccount.secret2FA,
	).Scan(&accID)
	if err != nil {
		log.Fatal(err)
	}

	// Add test account characters
	for classID := 1; classID <= 8; classID++ {
		var charID int
		err := db.QueryRow(`
            INSERT INTO Character (acc_id, class_id)
            VALUES ($1, $2) RETURNING char_id`,
			accID, classID,
		).Scan(&charID)
		if err != nil {
			log.Fatal(err)
		}

		// Test scores (high values for visibility)
		score := 10000 - (classID * 100) // Makes test account rank high
		_, err = db.Exec(`
            INSERT INTO Scores (char_id, reward_score)
            VALUES ($1, $2)`,
			charID, score,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
