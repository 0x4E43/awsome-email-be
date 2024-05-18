package db

import (
	"database/sql"
)

type DBCon struct {
	DB *sql.DB
}

func (d *DBCon) CreateRequiredTables() error {
	var sqls []string
	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS user_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(100),
		password VARCHAR(50),
		created_at TIMESTAMP
	);`)

	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS email_configs (
		id INT PRIMARY KEY,
		smtp_host VARCHAR(50),
		smtp_pass VARCHAR(100),
		smtp_from VARCHAR(100),
		is_default boolean,
		created_at TIMESTAMP
	);`)

	sqls = append(sqls, ` INSERT INTO user_details(id, email, password, created_at)
		SELECT 1, 'talk2nimai@gmail.com','$2a$10$BH5Ya0R/NefDs58YEkg7Vu/CO6tnbgKKKhOjuG4nbtmaH7QKtztOG', CURRENT_TIMESTAMP
		WHERE NOT EXISTS (SELECT 1 FROM user_details WHERE id=1);`)
	for _, query := range sqls {
		_, err := d.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
