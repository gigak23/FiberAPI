package main

import (
	"github.com/gigak23/FiberAPI.git/cmd/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	// setup routes
	setupRoutes(app)

	for _, route := range app.GetRoutes() {
		println("Registered route:", route.Method, route.Path)
	}

	//Listen on server 3000 and catch error
	err := app.Listen(":3000")

	//handle error
	if err != nil {
		panic(err)
	}
}

func setupRoutes(app *fiber.App) {

	// give response when at /
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "YOU ARE NOT AT API ENDPOINT",
		})
	})

	// api group at /api
	api := app.Group("/api")

	// give response when at /api
	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint 😉",
		})
	})

	//Created a nested group - /api/todos and passes it
	//into the TodoRoute function from routes package
	//connects todo routes
	routes.TodoRoute(api.Group("/todos"))

	//c.BodyParser(&struct{})
	// This will recieve requested data from the client and populate the struct with data
	// Then we can give a response back to client using
	// return c.JSON(struct)

}
