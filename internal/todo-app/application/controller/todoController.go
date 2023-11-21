package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-app/internal/todo-app/application/controller/request"
	"todo-app/internal/todo-app/application/handler"
)

type ITodo interface {
	GetAll(ctx *gin.Context)
	Save(ctx *gin.Context)
}

type todoController struct {
	todoCommandHandler handler.ITodoCommandHandler
}

func NewTodoController(todoCommandHandler handler.ITodoCommandHandler) ITodo {
	return &todoController{
		todoCommandHandler: todoCommandHandler,
	}
}

func (controller *todoController) GetAll(ctx *gin.Context) {
	fmt.Println("todoController.saveTodo INFO START")

	entityList, statusCode, err := controller.todoCommandHandler.GetAll(ctx)

	if err != nil {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entityList)
}

func (controller *todoController) Save(ctx *gin.Context) {
	var saveTodoRequest request.SaveTodoRequest
	err := ctx.BindJSON(&saveTodoRequest)

	if err != nil {
		fmt.Println("todoController.SaveTodo ERROR - Request bind json error occurred")
		ctx.JSON(http.StatusInternalServerError, errors.New("request bind json error occurred"))
		return
	}

	fmt.Printf("todoController.saveTodo Request: %#v\n", saveTodoRequest)

	statusCode, handlerError := controller.todoCommandHandler.Save(ctx, saveTodoRequest.ToCommand())
	if handlerError != nil {
		fmt.Printf("todoController.saveTodo error: %+v\n", handlerError)
		ctx.JSON(statusCode, gin.H{"error": handlerError.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "todo successfully created"})
}
