package user

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// user role enum
const (
	ADMIN = 0
	TEST  = 1
)

type User struct {
	Id        int       `json:"id,omitempty"`
	EmailId   string    `json:"emailId"`
	Password  string    `json:"password"`
	UserType  int       `json:"userType"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) createUser(db *sql.DB) (*User, error) {
	log.Println("Adding user to DB")
	var insertQuery = `INSERT INTO user_details (id, email, password, user_type, created_at)
		VALUES(NULL,$1, $2, $3, CURRENT_TIMESTAMP)`
	enc_pass, err := user.Encrypt_password()
	if err != nil {
		return nil, err
	}
	row, err := db.Query(insertQuery, user.EmailId, enc_pass, user.UserType)
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
	userQuery := `SELECT email, password, user_type FROM user_details WHERE email = $1`
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
	sqlQuery := `SELECT id, email, created_at, user_type FROM user_details`

	rows, err := db.Query(sqlQuery)

	var userList []User
	if err != nil {
		log.Println("Something went wrong while executing query: ", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var dbUser User
		if err := rows.Scan(&dbUser.Id, &dbUser.EmailId, &dbUser.CreatedAt, &dbUser.UserType); err != nil {
			log.Println("Error scanning row: ", err.Error())
			return nil, err
		}
		userList = append(userList, dbUser)
	}
	println("Size of user List ", len(userList))
	return userList, nil
}

func (user *User) DeleteUser(db *sql.DB, userId int) error {
	sqlQuery := "DELETE FROM user_details where id" + strconv.Itoa(userId)

	_, err := db.Exec(sqlQuery)
	if err != nil {
		log.Println("Error while deleting user ", err.Error())
		return err
	}
	return nil
}

// Security Related Utility Method

func (user *User) ParseJwtToken(tokenString string) (*string, error) {
	tokenString = strings.TrimSpace(tokenString)
	log.Println("Token :", tokenString)
	var MY_SECRET = "D2953AFCC7938B14DD1B969BB4535"
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(MY_SECRET), nil
	})
	if err != nil {
		log.Println("Error in parsing JWT : ", err.Error())
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
	var MY_SECRET = "D2953AFCC7938B14DD1B969BB4535"
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)), //For 15 days
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.EmailId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_str, err := token.SignedString([]byte(MY_SECRET))
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
