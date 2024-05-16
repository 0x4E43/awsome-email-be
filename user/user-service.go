package user

import (
	"database/sql"
	"log"
)


type User struct{
	Id string `json:"id,omitempty"`
	EmailId string `json:"emailId"`
	Password string `json:"password"`
}

func (user User) createUser(db *sql.DB) (*User, error){
	log.Println("Adding user to DB")
	var insertQuery = `INSERT INTO user_details (username, email, created_at)
		VALUES($1, $2, CURRENT_TIMESTAMP)`;
	row, err := db.Query(insertQuery, user.EmailId, user.Password)
	if err != nil {
		log.Println("Exception while adding user: ", err.Error())
		return nil, err
	}
	for row.Next(){
		err = row.Scan(&user.Id, &user.EmailId);
		if err != nil {
			log.Println("Unable to extract rows: ", err.Error())
			return nil, err
		}
	}
	return &user, nil
}