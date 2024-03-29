package main

import (
	"context"
	"flag"
	"log"

	"github.com/Stiffjobs/hotel-reservation/api"
	"github.com/Stiffjobs/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	//initialization
	var (
		app          = fiber.New(config)
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", api.AdminAuth)
		apiv1Public    = app.Group("/api")
		hotelHandler   = api.NewHotelHandler(store)
	)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello World",
		})
	})

	//auth handlers
	apiv1Public.Post("/auth", authHandler.HandleAuthenticate)

	//Versioned
	//user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetListHotel)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandleGetListRoom)

	//book handler
	apiv1.Get("/room", roomHandler.HandleGetListRoom)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// TODO: cancel a booking

	//bookings handler
	admin.Get("/booking", bookingHandler.HandleGetListBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	apiv1.Put("/booking/:id", bookingHandler.HandleCancelBooking)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	log.Fatal(app.Listen(*listenAddr))
}
