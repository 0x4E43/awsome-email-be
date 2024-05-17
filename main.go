package main

import (
	awmail "0x4E43/email-app-be/aw-mail"
	db_utils "0x4E43/email-app-be/db"
	user "0x4E43/email-app-be/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {

	var server = db_utils.DBCon{}
	db, err := sql.Open("sqlite3", "./db/email.db")
	if err != nil {
		log.Panic("Something went wrong while connecting to db", err.Error())
	}
	if err!=nil{
		println("Something went wrong while opeing DB", err.Error())
	}
	server.DB =  db
	if err:= server.CreateRequiredTables(); err!=nil{
		log.Panic("Failed to create tables : ", err.Error());
	}
	defer db.Close();
	// defer close
	println("DB Connected: ", db)
	e := echo.New()

	var userAPI =  new(user.UserAPI)
	userAPI.ConDB = db
	//USER RELATED ENDPOINTS
	e.POST("/user/create", userAPI.UserCreateHandler)
	e.POST("/user/login", userAPI.UserLoginHandler)
	e.GET("/user/list-all", userAPI.ListAllUserHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/send-email", awmail.EmailSenderHandler )
	e.Logger.Fatal(e.Start(":1323"))
}