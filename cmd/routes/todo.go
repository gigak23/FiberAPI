package routes

import (
	"github.com/gigak23/FiberAPI.git/cmd/controllers"
	"github.com/gofiber/fiber/v2"
)

// Takes the grouped routes such as an agrument and displays the
// to-do tasks from the controllers GetTodos function at the route passed in
// which is api/todos
func TodoRoute(route fiber.Router) {
	route.Get("", controllers.GetTodos)
}
