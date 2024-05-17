package user

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)


type User struct{
	Id string `json:"id,omitempty"`
	EmailId string `json:"emailId"`
	Password string `json:"password"`
}

func (user User) createUser(db *sql.DB) (*User, error){
	log.Println("Adding user to DB")
	var insertQuery = `INSERT INTO user_details (id, email, password, created_at)
		VALUES(NULL,$1, $2, CURRENT_TIMESTAMP)`;
	enc_pass, err := Encrypt_password(user.Password)
	if err != nil{
		return nil, err
	}
	row, err := db.Query(insertQuery, enc_pass, user.Password)
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

//METHOD TO ENCRYPT STRING
func Encrypt_password(pass string) (*string, error){
	enc, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while encrypting password", err.Error())
		return nil, err
	}
	retStr := string(enc) 
	return &retStr, nil
}