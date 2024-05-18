package main

import (
	awmail "0x4E43/email-app-be/aw-mail"
	cache "0x4E43/email-app-be/cache"
	db_utils "0x4E43/email-app-be/db"
	user "0x4E43/email-app-be/user"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {

	var server = db_utils.DBCon{}
	db, err := sql.Open("sqlite3", "./db/email.db")
	if err != nil {
		log.Panic("Something went wrong while connecting to db", err.Error())
	}
	if err != nil {
		println("Something went wrong while opeing DB", err.Error())
	}
	server.DB = db
	if err := server.CreateRequiredTables(); err != nil {
		log.Panic("Failed to create tables : ", err.Error())
	}
	defer db.Close()

	cache.LoadUserCache(db)
	// defer close
	log.Println("DB Connected: ", db)
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_unix} ${remote_ip} ${method} ${uri} ${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	//disable CSRF
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c echo.Context) bool {
			// Return true to skip CSRF protection for certain routes or requests
			// For example:
			// return c.Request().Method == http.MethodGet // Skip CSRF for GET requests
			return true // Skip CSRF for all requests (use with caution)
		},
	}))
	var userAPI = new(user.UserAPI)
	userAPI.ConDB = db
	//USER RELATED ENDPOINTS
	e.POST("/user/login", userAPI.UserLoginHandler) //public endpoint

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/send-email", awmail.EmailSenderHandler)
	// restricted endpoints
	r := e.Group("/app")
	r.Use(AuthenticateAPIV2)
	r.POST("/user/create", userAPI.UserCreateHandler)
	r.GET("/user/list-all", userAPI.ListAllUserHandler)
	r.DELETE("/user/delete/:userId", userAPI.UserDeleteHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

// middleware
func AuthenticateAPI(c echo.Context) error {
	log.Println("Auth Midddleware")
	return nil
}

func AuthenticateAPIV2(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Get the header and extract
		header := c.Request().Header.Get("Authorization")
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			tokenSlice := strings.Split(header, "Bearer ")
			if len(tokenSlice) != 2 {
				// Invalid Authorization header format
				// Handle the error accordingly
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			// Extract the JWT token
			token := strings.TrimSpace(tokenSlice[1])

			// parse Jwt Token
			var user user.User
			userName, err := user.ParseJwtToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Authentication Required"})

			}
			if cache.IsUserInCache(*userName) {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Please login"})
			}
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Authorization required"})
		}
	}
}
