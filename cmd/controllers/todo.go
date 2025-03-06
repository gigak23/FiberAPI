package controllers

import (
	"os"
	"time"

	"github.com/gigak23/FiberAPI.git/cmd/config"
	"github.com/gigak23/FiberAPI.git/cmd/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTodos(c *fiber.Ctx) error {
	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	//Create empty query to filter data we want
	query := bson.D{}

	//Stores all the data we queried
	cursor, err := todoCollection.Find(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot store documents",
			"error":   err.Error(),
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

	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

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
	data.ID = primitive.NewObjectID()
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

	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	//Gets the users id provided in url
	paramID := c.Params("id")

	//Convert to objectID
	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Not a number",
			"error":   err,
		})
	}

	//Create emtpy todo struct
	todo := &models.Todo{}

	//Query the todo data based on id number
	query := bson.D{{Key: "_id", Value: id}}

	//Find the todo in database and store it in todo struct
	//Otherwise store it in error and handle it
	err = todoCollection.FindOne(c.Context(), query).Decode(todo)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "cannot decode data",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})

}

func UpdateTodo(c *fiber.Ctx) error {

	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	paramID := c.Params("id")

	//Convert to objectID
	id, err := primitive.ObjectIDFromHex(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Not a number",
			"error":   err,
		})
	}

	//Create instance of models.Todo struct
	data := new(models.Todo)
	err = c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "could not parse",
			"error":   err,
		})
	}

	//Query todo by id (_id is default key in MongoDB)
	query := bson.D{{Key: "_id", Value: id}}

	//Update todo variable
	var dataToUpdate bson.D

	//Change values based if they are holding memory address or nil (Whether client requested to update this value)
	//Append bson.E key-value pairs to bson.D slice
	if data.Task != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "task", Value: data.Task})
	}
	if data.Completed != nil {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "completed", Value: data.Completed})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	//Then we assign the key-value pairs to bson.D
	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	//Find and update todo in MongoDB
	err = todoCollection.FindOneAndUpdate(c.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update todo",
			"error":   err,
		})
	}

	//Create an instance of Todo struct
	todo := &models.Todo{}

	//Query the todo by id and store the data in todo variable
	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})

}

func DeleteTodo(c *fiber.Ctx) error {

	todoCollection := config.MI.DB.Collection(os.Getenv("TODO_COLLECTION"))

	paramID := c.Params("id")

	//Convert to objectID
	id, err := primitive.ObjectIDFromHex(paramID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID not a number",
			"error":   err,
		})
	}

	//Query todo by id
	query := bson.D{{Key: "_id", Value: id}}

	//Find in MongoDB and delete by query filter
	err = todoCollection.FindOneAndDelete(c.Context(), query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete todo",
			"error":   err,
		})
	}

	//Returns a successful response but with no body or contents
	return c.SendStatus(fiber.StatusNoContent)

}
