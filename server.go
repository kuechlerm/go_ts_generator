package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const InsertTodo_Path = "/todos"

type InsertTodo_Request struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type InsertTodo_Response struct {
	Id string `json:"id"`
}

type Todo struct {
	Title string `json:"title"`
}

func Init_server() {
	server := gin.Default()
	validate := validator.New()

	server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	server.POST(InsertTodo_Path, func(c *gin.Context) {
		args, err := MapAndValidate[InsertTodo_Request](c, validate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Todo{Title: args.Title})
	})

	// server.Logger.Fatal(server.Start(":5080"))
	server.Run(":5080")
}

func MapAndValidate[T any](c *gin.Context, validate *validator.Validate) (*T, error) {
	args := new(T)
	if err := c.Bind(args); err != nil {
		return args, err
	}

	if err := validate.Struct(args); err != nil {
		return args, err
	}

	return args, nil
}
