package api

import (
	"errors"
	"fmt"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err

	}
	user, err := h.userStore.GetByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentails")
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password)); err != nil {
		return fmt.Errorf("invalid credentails")
	}
	fmt.Println("user authenticated: ", user)
	return nil
}
