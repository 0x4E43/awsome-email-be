package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAPI struct{
	ConDB *sql.DB
}

func (userApi *UserAPI)UserCreateHandler(c echo.Context) error{
	var user = User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
        return c.String(http.StatusForbidden, "Something went wrong") // Handle error if any
    } 
	log.Println("UserName: ", user.EmailId, " Pass: ", user.Password)
	added_user, err := user.createUser(userApi.ConDB)
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, added_user)
}

func (userApi *UserAPI)UserLoginHandler(c echo.Context) error{
	var user = User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
        return c.String(http.StatusForbidden, "Something went wrong") // Handle error if any
    } 
	log.Println("UserName: ", user.EmailId, " Pass: ", user.Password)
	//check user
	dbUser, err := user.CheckIfUserExist(userApi.ConDB)
	if dbUser.EmailId == "" {
		return c.JSON(http.StatusOK, "No User Found")
	} 
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}
	isAuthenticated := user.Compare_password(dbUser.Password)
	if !isAuthenticated{
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid Creds"})
	}
	token, err := dbUser.Create_auth_token()
	if err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": *token})
}
