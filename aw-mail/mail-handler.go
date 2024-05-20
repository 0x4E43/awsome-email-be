package awmail

import (
	"0x4E43/email-app-be/global"
	"0x4E43/email-app-be/logger"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EmailDBConfig struct {
	ConDB *sql.DB
}

// Custom logger
var log = logger.Log

func (emailDBConfig *EmailDBConfig) EmailSenderHandler(c echo.Context) error {
	var emailBody Email
	//get the body
	// var body = c.Request().Body
	if err := json.NewDecoder(c.Request().Body).Decode(&emailBody); err != nil {
		log.Println(err.Error())
		res := global.PrepareResponse("Invalid data", http.StatusBadRequest, nil)
		return c.JSON(res.Status, res)
	}
	log.Printf("%+v", emailBody)
	if emailBody.MailType == 1 && (&emailBody.To == nil || emailBody.To == "") {
		log.Println("TO is required for normal email")
		res := global.PrepareResponse("Invalid data", http.StatusBadRequest, nil)
		return c.JSON(res.Status, res)
	}

	if emailBody.MailType == 0 {
		emailBody.SendTestEmail(emailDBConfig.ConDB)
	}
	println("{} {}", emailBody.To, emailBody.MailType)

	return c.String(http.StatusOK, "Hello Email")
}

