package api

import (
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/Stiffjobs/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return ErrInvalidID()
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return ErrNotFound()
	}
	user, err := getAuthedUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	if err := h.store.Booking.Update(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Message: "updated"})
}

func (h *BookingHandler) HandleGetListBooking(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetList(c.Context(), db.Map{})
	if err != nil {
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}
