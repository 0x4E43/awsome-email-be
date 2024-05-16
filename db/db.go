package db

import "database/sql"

type DBCon struct{
	DB *sql.DB
}


func (d *DBCon) CreateRequiredTables() error{
	var sqls []string 
	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS user_details (
		id INT PRIMARY KEY,
		username VARCHAR(50),
		email VARCHAR(100),
		created_at TIMESTAMP
	);`);

	sqls = append(sqls, `CREATE TABLE IF NOT EXISTS email_configs (
		id INT PRIMARY KEY,
		smtp_host VARCHAR(50),
		smtp_pass VARCHAR(100),
		smtp_from VARCHAR(100),
		created_at TIMESTAMP
	);`);

	sqls = append(sqls, ` INSERT INTO user_details(id, username, email, created_at)
		SELECT 1, 'nimai', 'talk2nimai@gmail.com', CURRENT_TIMESTAMP
		WHERE NOT EXISTS (SELECT 1 FROM user_details WHERE id=1);`)
	for _ , query :=range sqls{
		_, err := d.DB.Exec(query)
		if err != nil{
			return err
		}
	}
	return nil
}