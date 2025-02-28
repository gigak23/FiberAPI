package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        string `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{
		ID:        "1",
		Task:      "Learn Golang",
		Completed: false,
	},
	{
		ID:        "2",
		Task:      "Do homework",
		Completed: false,
	},
}

func GetTodos(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todos": todos,
		},
	})
}
