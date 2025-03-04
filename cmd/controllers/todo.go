package controllers

import (
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/gigak23/FiberAPI.git/cmd/config"
	"github.com/gigak23/FiberAPI.git/cmd/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTodos(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	//Create empty query to filter data we want
	var query bson.D

	//Stores all the data we queried
	cursor, err := todoCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot store documents",
			"error":   err,
		})
	}

	var todos []models.Todo = make([]models.Todo, 0)

	//Assings all data in cursor into a Todo struct in todos slice
	err = cursor.All(c.Context(), &todos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot assign ",
			"error":   err.Error(),
		})
	}

	//returns queried todos
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todos": todos,
		},
	})
}

func CreateTodo(c *fiber.Ctx) error {

	todoCollection := config.MI.DB.Collection("TODO_COLLECTION")

	data := new(models.Todo)

	//Tries to convert assign JSON value into Todo struct data
	//error if this cannot be done
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "could not parse",
			"error":   err.Error(),
		})
	}

	//Assign default values to Todo struct data
	data.ID = nil
	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	//Inserts the data into the collection in mongoDB
	result, err := todoCollection.InsertOne(c.Context(), data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert todo",
			"error":   err.Error(),
		})
	}

	//Get the inserted data
	todo := &models.Todo{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	//Displays the created todo
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
