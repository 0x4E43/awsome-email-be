package awmail

import (
	"0x4E43/email-app-be/global"
	"0x4E43/email-app-be/logger"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type EmailDBConfig struct {
	ConDB *sql.DB
}

// Custom logger
var log = logger.Log

func (emailDBConfig *EmailDBConfig) EmailSenderHandler(c echo.Context) error {
	var res *global.Response
	var emailBody Email
	//get the body
	// var body = c.Request().Body
	if err := json.NewDecoder(c.Request().Body).Decode(&emailBody); err != nil {
		log.Println(err.Error())
		res := global.PrepareResponse("Invalid data", http.StatusBadRequest, nil)
		return c.JSON(res.Status, res)
	}
	log.Printf("%+v", emailBody)
	if emailBody.MailType == 0 && emailBody.To == "" {
		log.Println("TO is required for normal email")

		res := global.PrepareResponse("Invalid data", http.StatusBadRequest, nil)
		return c.JSON(res.Status, res)
	}
	if emailBody.MailType == 0 { //Normal Email
		var emailRec []string
		for _, email := range strings.Split(emailBody.To, ",") {
			if email != "" {
				emailRec = append(emailRec, strings.TrimSpace(email))
			}
		}
		log.Println("Email recipients", emailRec)
		err := emailBody.sendEmail(emailRec, emailDBConfig.ConDB)
		if err != nil {
			res = global.PrepareResponse("Something went wrong", http.StatusInternalServerError, nil)
			return c.JSON(res.Status, res)
		}
		log.Println("Email request added for normal email")
		res = global.PrepareResponse("Email request set successfully", http.StatusOK, emailBody)
	}

	if emailBody.MailType == 1 { //test emaikks
		err := emailBody.SendTestEmail(emailDBConfig.ConDB)
		if err != nil {
			res = global.PrepareResponse("Something went wrong", http.StatusInternalServerError, nil)
			return c.JSON(res.Status, res)
		}
	}
	println("{} {}", emailBody.To, emailBody.MailType)

	return c.JSON(res.Status, res)
}
