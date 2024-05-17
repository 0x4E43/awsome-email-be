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

func (user *User) createUser(db *sql.DB) (*User, error){
	log.Println("Adding user to DB")
	var insertQuery = `INSERT INTO user_details (id, email, password, created_at)
		VALUES(NULL,$1, $2, CURRENT_TIMESTAMP)`;
	enc_pass, err := user.Encrypt_password()
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
	return user, nil
}

func (user *User) CheckIfUserExist(db *sql.DB) (*User, error){
	userQuery := `SELECT email, password FROM user_details WHERE email = $1`;
	rows, err := db.Query(userQuery, user.EmailId)
	if err!= nil {
		log.Println("Something went wrong while geting user: ", err.Error())
		return nil, err
	}
	log.Printf("{%+v}", rows)
	dbUser := User{}
	for rows.Next(){
		rows.Scan(&dbUser.EmailId, &dbUser.Password)
	}
	log.Println("User Data user: ", dbUser.EmailId)
	return &dbUser, nil
}


// Security Related Utility Method


func (user *User) Create_auth_token() string{
	return ""
}


func (user *User)Compare_password(enc_pass string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(enc_pass), []byte(user.Password))
	if err != nil {
		log.Println("Exception ", err.Error())
		return false;
	}
	return true
}



func (user *User)Encrypt_password() (*string, error){
	enc, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while encrypting password", err.Error())
		return nil, err
	}
	retStr := string(enc) 
	return &retStr, nil
}

