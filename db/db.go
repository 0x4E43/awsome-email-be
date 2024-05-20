package db

import (
	"0x4E43/email-app-be/logger"
	"database/sql"
	"os"
)

// Custom logger
var log = logger.Log

type DBCon struct {
	DB *sql.DB
}

type OsEnv struct {
	AdminUser string
	AdminPass string
	SmtpFrom  string
	SmtpHost  string
	SmtpPort  string
	SmtpPass  string
}

func getEnvs() *OsEnv {
	osEnv := OsEnv{
		AdminUser: os.Getenv("DEFAULT_ADMIN_EMAIL"),
		AdminPass: os.Getenv("DEFAULT_ADMIN_PASS"),
		SmtpHost:  os.Getenv("SMTP_DEFAULT_HOST"),
		SmtpPass:  os.Getenv("SMTP_DEFAULT_PASS"),
		SmtpFrom:  os.Getenv("SMTP_DEFAULT_FROM"),
		SmtpPort:  os.Getenv("SMTP_DEFAULT_PORT"),
	}
	return &osEnv
}

func (d *DBCon) CreateRequiredTables() error {
	var sqls []string
	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS user_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(100),
		password VARCHAR(50),
		user_type INTEGER,
		created_at TIMESTAMP
	);`)

	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS email_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		smtp_host VARCHAR(50),
		smtp_pass VARCHAR(100),
		smtp_from VARCHAR(100),
		smtp_port INTEGER,
		is_default BOOLEAN,
		created_at TIMESTAMP
	);`)

	//get ebv variables
	osEnvs := getEnvs()
	sqls = append(sqls, ` INSERT INTO user_details(id, email, password, created_at, user_type)
		SELECT 1, '`+osEnvs.AdminUser+`','`+osEnvs.AdminPass+`', CURRENT_TIMESTAMP, 1
		WHERE NOT EXISTS (SELECT 1 FROM user_details WHERE id=1);`)

	sqls = append(sqls, ` INSERT INTO email_configs(id, smtp_host, smtp_pass, smtp_from, smtp_port, is_default, created_at)
		SELECT 1, '`+osEnvs.SmtpHost+`','`+osEnvs.SmtpPass+`', '`+osEnvs.SmtpFrom+`','`+osEnvs.SmtpPort+`', 1, CURRENT_TIMESTAMP
		WHERE NOT EXISTS (SELECT 1 FROM email_configs WHERE id=1);`)
	for _, query := range sqls {
		// log.Println(query)
		_, err := d.DB.Exec(query)
		if err != nil {
			log.Printf("%+v", err)
			log.Fatal(err.Error())
			return err
		}
		// log.Println("Query Executed: ", query)
	}
	return nil
}
