package user

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string    `json:"id,omitempty"`
	EmailId   string    `json:"emailId"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) createUser(db *sql.DB) (*User, error) {
	log.Println("Adding user to DB")
	var insertQuery = `INSERT INTO user_details (id, email, password, created_at)
		VALUES(NULL,$1, $2, CURRENT_TIMESTAMP)`
	enc_pass, err := user.Encrypt_password()
	if err != nil {
		return nil, err
	}
	row, err := db.Query(insertQuery, user.EmailId, enc_pass)
	if err != nil {
		log.Println("Exception while adding user: ", err.Error())
		return nil, err
	}
	for row.Next() {
		err = row.Scan(&user.Id, &user.EmailId)
		if err != nil {
			log.Println("Unable to extract rows: ", err.Error())
			return nil, err
		}
	}
	return user, nil
}

func (user *User) CheckIfUserExist(db *sql.DB) (*User, error) {
	userQuery := `SELECT email, password FROM user_details WHERE email = $1`
	rows, err := db.Query(userQuery, user.EmailId)
	if err != nil {
		log.Println("Something went wrong while geting user: ", err.Error())
		return nil, err
	}
	log.Printf("{%+v}", rows)
	dbUser := User{}
	for rows.Next() {
		rows.Scan(&dbUser.EmailId, &dbUser.Password)
	}
	log.Println("User Data user: ", dbUser.EmailId)
	return &dbUser, nil
}

func (user *User) ListAllUser(db *sql.DB) ([]User, error) {
	sqlQuery := `SELECT email, created_at FROM user_details`

	rows, err := db.Query(sqlQuery)

	defer rows.Close()

	var userList []User
	if err != nil {
		log.Println("Something went wrong while executing query: ", err.Error())
	}
	for rows.Next() {
		var dbUser User
		if err := rows.Scan(&dbUser.EmailId, &dbUser.CreatedAt); err != nil {
			log.Println("Error scanning row: ", err.Error())
			return nil, err
		}
		userList = append(userList, dbUser)
	}
	println("Size of user List ", len(userList))
	return userList, nil
}

// Security Related Utility Method
var MY_SECRET = []byte("D2953AFCC7938B14DD1B969BB4535")

func (user *User) ParseJwtToken(tokenString string) (*string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(MY_SECRET), nil
	})
	if err != nil {
		log.Println("Error in parsing JWT")
		return nil, err
	}

	if token.Valid {
		userName, err := token.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		return &userName, nil
	} else {
		return nil, errors.New("Unauthorized")
	}
}

func (user *User) Create_auth_token() (*string, error) {

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)), //For 15 days
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.EmailId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_str, err := token.SignedString(MY_SECRET)
	if err != nil {
		log.Println("Something went wrong while generating token ", err.Error())
		return nil, err
	}
	return &token_str, nil
}

func (user *User) Compare_password(enc_pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(enc_pass), []byte(user.Password))
	if err != nil {
		log.Println("Exception ", err.Error())
		return false
	}
	return true
}

func (user *User) Encrypt_password() (*string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while encrypting password", err.Error())
		return nil, err
	}
	retStr := string(enc)
	return &retStr, nil
}
