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
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})

	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	log.Fatal(app.Listen(*listenAddr))
}
