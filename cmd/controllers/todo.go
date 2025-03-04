package controllers

import (
	"github.com/gofiber/fiber/v2"
	//"fmt"
	"slices"
	"strconv"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

// Creates a slice of pointer to Todo struct
var todos = []*Todo{
	{
		ID:        1,
		Task:      "Grind",
		Completed: true,
	},
	{
		ID:        2,
		Task:      "Learn Golang",
		Completed: false,
	},
	{
		ID:        3,
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

func CreateTodo(c *fiber.Ctx) error {

	//Holds the task value
	type Request struct {
		Task string `json:"task"`
	}

	var body Request

	//Tries to convert assign JSON value into Request struct body
	//error if this cannot be done
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "could not parse",
		})
	}

	//Creates a new todo in the todos struct
	//with the task provided form user
	todo := &Todo{
		ID:        len(todos) + 1,
		Task:      body.Task,
		Completed: false,
	}

	//Adds new todo into todos slice
	todos = append(todos, todo)

	//Displays the list of todos on server
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    todo,
	})

}

func GetTodo(c *fiber.Ctx) error {

	//Gets the users id provided in url
	paramID := c.Params("id")

	//Converts id to int and checks if conversion is valid
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Not a number",
		})
	}

	//Checks if id matches any todo from the todos struct
	for _, todo := range todos {
		if id == todo.ID {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"data": fiber.Map{
					"todo": todo,
				},
			})
		}
	}

	//If id does not match, then state id cannot be found
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"sucdess": false,
		"messgae": "Could not find ID",
	})
}

func UpdateTodo(c *fiber.Ctx) error {

	paramID := c.Params("id")

	id, err := strconv.Atoi(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Not a number",
		})
	}

	//Create pointer field so we know which one to update
	type Request struct {
		Task      *string `json:"task"`
		Completed *bool   `json:"completed"`
	}

	var body Request

	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "could not parse",
		})
	}

	//Set new todo struct = original todos struct values
	var todo *Todo

	for _, t := range todos {
		if id == t.ID {
			todo = t
			break
		}
	}

	//Change values based if they are holding memory address or nil (Whether client requested to update this value)
	if body.Task != nil {
		todo.Task = *body.Task
	}
	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})

}

func DeleteTodo(c *fiber.Ctx) error {

	paramID := c.Params("id")

	id, err := strconv.Atoi(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID not a number",
		})
	}

	//Removes a task struct from todos slice
	for i, t := range todos {
		if t.ID == id {

			//todos = append(todos[:i], todos[i+1:]...)
			todos = slices.Delete(todos, i, i+1)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"data": fiber.Map{
					"todo removed": todos,
				},
			})
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "ID not found",
	})

}
