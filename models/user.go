package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

type ResponseUser struct {
	Error   error
	Message string
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, user *User) (string, *ResponseUser) {
	//get user from email
	loginUser := new(User)
	result := db.Where("email = ?", user.Email).First(loginUser)
	if result.Error != nil {
		return "", &ResponseUser{
			Error:   result.Error,
			Message: "Email is not found",
		}
	}

	//compare password with is already hashed with the password in the database
	err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password))
	if err != nil {
		return "", &ResponseUser{
			Error:   err,
			Message: "Password is incorrect",
		}
	}

	//if passed will return jwt token to fiber then fiber will use it to authenticate the user via cookies(context)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = loginUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", &ResponseUser{
			Error:   err,
			Message: "Something went wrong",
		}
	}
	return t, nil

}
