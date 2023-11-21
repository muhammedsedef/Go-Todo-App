package request

import "todo-app/internal/todo-app/application/handler"

type SaveTodoRequest struct {
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	AddressInfo Address `json:`
}

type Address struct {
	Name   string `json:"name"`
	CityId int64  `json:"cityId"`
}

func (request *SaveTodoRequest) ToCommand() handler.SaveTodoCommand {
	return handler.SaveTodoCommand{
		Title:   request.Title,
		Content: request.Content,
	}
}
