package awmail

import (
	"database/sql"
	"net/smtp"
	"regexp"
	"strconv"

	"github.com/jordan-wright/email"
)

type EmailConfig struct {
	Id        string `json:"id"`
	SmtpHOST  string `json:"smtpHost"`
	SmtpPass  string `json:"smtpPass"`
	SmtpPort  int    `json:"smtpPort"`
	SmtpFrom  string `json:"smtpFrom"`
	IsDefault bool   `json:"isDefault"`
}

type Email struct {
	To       string `json:"toAddress"`
	Subject  string `json:"subject"`
	Body     string `json:"mailBody"`
	MailType int    `json:"mailType"`
}

func (email *Email) SendTestEmail(db *sql.DB) error {
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

	log.Printf("Hello Test EMail {%v}", userList)
	return email.sendEmail(userList, db)
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
	contacts = *getValidEmails(contacts)
	log.Println("Contacts: ", getValidEmails(contacts), " Size", len(contacts))
	em.From = emailConfig.SmtpFrom
	em.To = contacts
	em.Subject = cMail.Subject
	em.HTML = []byte(cMail.Body)

	err = em.Send(emailConfig.SmtpHOST+":"+strconv.Itoa(emailConfig.SmtpPort), smtp.PlainAuth("", emailConfig.SmtpFrom, emailConfig.SmtpPass, emailConfig.SmtpHOST))
	log.Printf("{%+v}", emailConfig)
	if err != nil {
		log.Println("Something went wrong while sending email ", err.Error())
		// log.Panicf("{%+v}", err)
		return err
	}
	log.Println("Email sent successfully")
	return nil
}

func getValidEmails(contacts []string) *[]string {
	// Regular expression for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex pattern
	regex := regexp.MustCompile(emailRegex)

	// Check if the email matches the regex pattern
	var emailList []string
	for _, email := range contacts {
		if regex.MatchString(email) {
			emailList = append(emailList, email)
		}
	}
	return &emailList
}

func (emailConfig *EmailConfig) AddNewEmailConfig(db *sql.DB) (*EmailConfig, error) {
	if emailConfig.IsDefault {
		query := "UPDATE email_configs SET is_default = 0 WHERE is_default =1"
		_, err := db.Exec(query)
		if err != nil {
			log.Println("Error while setting defaults for emailconfig ", err.Error())
			return nil, err
		}
	}

	query := `INSERT INTO email_configs (smtp_host, smtp_pass, smtp_port, smtp_from, is_default, created_at) 
			 VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)  RETURNING id;`
	result, err := db.Exec(query, emailConfig.SmtpHOST, emailConfig.SmtpPass, emailConfig.SmtpPort, emailConfig.SmtpFrom, emailConfig.IsDefault)

	if err != nil {
		log.Println("Error while writing data into DB ", err.Error())
		return nil, err
	}

	id, _ := result.LastInsertId()
	emailConfig.Id = strconv.FormatInt(id, 10)
	return emailConfig, nil
}

func (emailConfig *EmailConfig) checkIfConfigHostExists(db *sql.DB) (*EmailConfig, error) {
	query := "SELECT smtp_host FROM email_configs WHERE smtp_host = $1 LIMIT 1"
	row, err := db.Query(query, emailConfig.SmtpHOST)
	defer row.Close()

	if err != nil {
		log.Println("Something went wrong ", err.Error())
		return nil, err
	}
	tempConfig := EmailConfig{}
	for row.Next() {
		if err := row.Scan(&tempConfig.SmtpHOST); err != nil {
			log.Println("Something went wrong while scan ", err.Error())
			return nil, err
		}
		return &tempConfig, nil
	}
	return nil, nil
}
