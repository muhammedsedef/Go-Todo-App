package main

import (
	"github.com/gin-gonic/gin"
	configuration "todo-app/appconfig"
	"todo-app/internal/todo-app/application/controller"
	"todo-app/internal/todo-app/application/handler"
	"todo-app/internal/todo-app/application/repository"
	"todo-app/pkg/couchbase"
	"todo-app/pkg/server"
)

func main() {
	engine := gin.New()

	//couchbase
	couchbaseCluster := couchbase.ConnectCluster(
		configuration.CouchbaseHost,
		configuration.CouchbaseUsername,
		configuration.CouchbasePassword,
	)

	todoRepository := repository.NewTodoRepository(couchbaseCluster)
	todoCommandHandler := handler.NewTodoCommandHandler(todoRepository)
	todoController := controller.NewTodoController(todoCommandHandler)

	todoGroup := engine.Group("api/v1/todo")
	todoGroup.GET("", todoController.GetAll)
	todoGroup.POST("", todoController.Save)

	server.NewServer(engine).StartHttpServer()
}
