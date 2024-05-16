package awmail

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

var EmailBody struct{
	To string`json:"to"`
	From string `json:"from"`
}

func EmailSenderHandler(c echo.Context) error{
	//get the body
	// var body = c.Request().Body
	if err := json.NewDecoder(c.Request().Body).Decode(&EmailBody); err != nil {
        return err // Handle error if any
    }
    println("{} {}", EmailBody.To, EmailBody.From)
    return c.String(http.StatusOK, "Hello Email")
}