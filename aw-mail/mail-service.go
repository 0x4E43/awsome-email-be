package awmail

import (
	"database/sql"
	"net/smtp"
	"strconv"

	"github.com/jordan-wright/email"
)

type EmailConfig struct {
	SmtpHOST  string `json:"smtpHost"`
	SmtpPass  string `json:"smtpPass"`
	SmtpPort  int    `json:"smtpPort"`
	SmtpFrom  string `json:"smtpFrom"`
	IsDefault bool   `json:"isDefault"`
}

type Email struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	MailType int    `json:"mailType"`
}

func (email *Email) SendTestEmail(db *sql.DB) {
	// find user from user details where usertype = 1

	sqlQuery := "SELECT email FROM user_details WHERE user_type = 0"

	var userList []string

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Println("Error while query ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Println("Error while geting data ", err.Error())
		}
		userList = append(userList, email)
	}

	email.sendEmail(userList, db)
	log.Printf("Hello Test EMail {%v}", userList)
}

func (cMail *Email) sendEmail(contacts []string, db *sql.DB) error {
	emailConfigQuery := "SELECT smtp_host, smtp_pass, smtp_from, smtp_port FROM email_configs WHERE is_default = 1"
	var emailConfig EmailConfig
	rows, err := db.Query(emailConfigQuery)
	if err != nil {
		log.Println("Error while query ", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&emailConfig.SmtpHOST, &emailConfig.SmtpPass, &emailConfig.SmtpFrom, &emailConfig.SmtpPort)
		if err != nil {
			log.Println("Error while geting data ", err.Error())
			return err
		}
	}

	//send email
	em := email.NewEmail()

	em.From = emailConfig.SmtpFrom
	em.To = contacts
	em.Subject = cMail.Subject
	em.HTML = []byte(cMail.Body)

	err = em.Send(emailConfig.SmtpHOST+":"+strconv.Itoa(emailConfig.SmtpPort), smtp.PlainAuth("", emailConfig.SmtpFrom, emailConfig.SmtpPass, emailConfig.SmtpHOST))
	log.Printf("{%+v}", emailConfig)
	if err != nil {
		log.Println("Something went wrong while sending email ", err.Error())
		return err
	}
	log.Println("Email sent successfully")
	return nil
}
