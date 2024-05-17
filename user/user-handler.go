package user

import (
	"0x4E43/email-app-be/cache"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAPI struct {
	ConDB *sql.DB
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

func (userApi *UserAPI) UserCreateHandler(c echo.Context) error {
	var user = User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.String(http.StatusForbidden, "Something went wrong") // Handle error if any
	}
	log.Println("UserName: ", user.EmailId, " Pass: ", user.Password)
	added_user, err := user.createUser(userApi.ConDB)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, added_user)
}

func (userApi *UserAPI) UserLoginHandler(c echo.Context) error {
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
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	isAuthenticated := user.Compare_password(dbUser.Password)
	if !isAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid Creds"})
	}
	token, err := dbUser.Create_auth_token()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
	}
	//set up cache for authorization
	userEmail := dbUser.EmailId
	cache.AddUserToCache(userEmail)
	return c.JSON(http.StatusOK, map[string]string{"token": *token})
}

// func (userApi *UserAPI) ListAllUserHandler(c echo.Context) error {
// 	var user = User{}
// 	userList, err := user.ListAllUser(userApi.ConDB)
// 	if err != nil {
// 		log.Println(" Exception while listing user : ", err.Error())
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
// 	}
// 	println("Size of user List 2 ", len(userList))
// 	var res = Response{
// 		Message: "User list successfull",
// 		Status:  http.StatusOK,
// 		Data:    userList,
// 	}
// 	return c.JSON(http.StatusOK, res)
// }

func (userApi *UserAPI) ListAllUserHandler(c echo.Context) error {
	var user = User{}
	userList, err := user.ListAllUser(userApi.ConDB)
	if err != nil {
		log.Println("Exception while listing user:", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
	}

	log.Println("Size of user List:", len(userList))

	var res = Response{
		Message: "User list successful",
		Status:  http.StatusOK,
		Data:    userList,
	}

	// Send the response
	return c.JSON(http.StatusOK, res)
}
