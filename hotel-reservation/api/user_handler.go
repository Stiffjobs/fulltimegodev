package api

import (
	"errors"
	"net/http"

	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound("user")
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetList(c.Context())
	if err != nil {
		return ErrResourceNotFound("users")
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}
	insertedUser, err := h.userStore.Create(c.Context(), user)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		update types.UpdateUserParams
		userID = c.Params("id")
	)

	if err := c.BodyParser(&update); err != nil {
		return ErrBadRequest()
	}
	if err := h.userStore.Update(c.Context(), userID, update); err != nil {
		return ErrResourceNotFound("user")
	}

	return c.JSON(map[string]string{"updated": userID})
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.Delete(c.Context(), userID); err != nil {
		return ErrResourceNotFound("user")
	}
	return c.JSON(map[string]string{"deleted": userID})
}
