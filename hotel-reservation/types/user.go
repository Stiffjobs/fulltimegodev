package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcyptCost          = 12
	minFirstNameLength = 2
	minLastNameLength  = 7
	minPasswordLength  = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.FirstName) > 0 {
		m["firstName"] = params.FirstName
	}
	if len(params.LastName) > 0 {
		m["lastName"] = params.LastName
	}
	return m
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstName should be at least %d characters", minFirstNameLength)
	}

	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastName should be at least %d characters", minLastNameLength)
	}
	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password should be at least %d characters", minPasswordLength)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is not valid", params.Email)
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	return emailRegex.MatchString(email)

}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcyptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
