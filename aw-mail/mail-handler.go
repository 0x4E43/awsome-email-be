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

func (emailDbConfig *EmailDBConfig) AddEmailConfigHandler(c echo.Context) error {
	//check if all fields are there orn not
	var emailConfig EmailConfig

	if err := json.NewDecoder(c.Request().Body).Decode(&emailConfig); err != nil {
		log.Println("Error while decoding response body")
		res := global.PrepareResponse("Something went wrong", http.StatusInternalServerError, nil)
		return c.JSON(res.Status, res)
	}

	if emailConfig.SmtpFrom == "" || emailConfig.SmtpHOST == "" || emailConfig.SmtpPass == "" || emailConfig.SmtpPort == 0 {
		log.Println("Not all the fields are provided in the request")
		res := global.PrepareResponse("Please provide all the fields", http.StatusBadRequest, nil)
		return c.JSON(res.Status, res)
	}
	newEmailConfig, err := emailConfig.checkIfConfigHostExists(emailDbConfig.ConDB)
	if newEmailConfig != nil {
		res := global.PrepareResponse("Email config with same host already exists", http.StatusOK, nil)
		return c.JSON(res.Status, res)
	}
	newEmailConfig, err = emailConfig.AddNewEmailConfig(emailDbConfig.ConDB)
	if err != nil {
		res := global.PrepareResponse("Something went wrong", http.StatusInternalServerError, nil)
		return c.JSON(res.Status, res)
	}
	res := global.PrepareResponse("Email config saved successfully", http.StatusOK, *newEmailConfig)
	return c.JSON(res.Status, res)
}
