# WIRA's Ranking Dashboard

## Instructions
### How to run the project
1. Gitclone the project
2. Change the database configuration in the api/handlers.go file ( right now still use my local database configuration ) :
   - run the Posgres DB
3. Go to root project :
   - run [ go mod tidy ]
   - run [ go run main.go ]
4. Go to wira frontend :
   - run [ npm install ( For windows ) / sudo npm install ( For MacOS ) ]
   - run [ npm run dev ( For windows ) / sudo npm run dev ( For MacOS ) ]
5. - Start redis server ( run redis-server )
6. Can use Postman / Swagger to test the backend API
7. Go to localhost:5173 to test the frontend side
8. API endpoint :
   - "http://localhost:8080/login" for login and 2FA
   - "http://localhost:8080/dashboard?class_id=1&page=1&limit=20" for player's scores ( change the query params )

