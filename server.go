package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
	server := echo.New()
	validate := validator.New()

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	server.POST(InsertTodo_Path, func(c echo.Context) error {
		args, err := MapAndValidate[InsertTodo_Request](c, validate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, Todo{Title: args.Title})
	})

	server.Logger.Fatal(server.Start(":5080"))
}

func MapAndValidate[T any](c echo.Context, validate *validator.Validate) (*T, error) {
	args := new(T)
	if err := c.Bind(args); err != nil {
		return args, err
	}

	if err := validate.Struct(args); err != nil {
		return args, err
	}

	return args, nil
}
