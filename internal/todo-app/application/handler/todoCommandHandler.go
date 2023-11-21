package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"todo-app/internal/todo-app/application/repository"
	"todo-app/internal/todo-app/domain"
)

type ITodoCommandHandler interface {
	Save(ctx context.Context, command SaveTodoCommand) (int, error)
	GetAll(ctx context.Context) ([]*domain.Todo, int, error)
}

type todoCommandHandler struct {
	todoRepository repository.ITodoRepository
}

func NewTodoCommandHandler(todoRepository repository.ITodoRepository) ITodoCommandHandler {
	return &todoCommandHandler{todoRepository: todoRepository}
}

func (handler *todoCommandHandler) Save(ctx context.Context, command SaveTodoCommand) (int, error) {
	fmt.Printf("todoCommandHandler.Save INFO START - command: %+v\n", command)

	entity, statusCode, err := handler.todoRepository.FindById(ctx, command.Title)
	if entity != nil {
		fmt.Printf("todoCommandHandler.Save WARNING! Title %s exist\n", command.Title)
		return http.StatusBadRequest, errors.New(fmt.Sprintf("given title: %s is already exist", command.Title))
	}

	if err != nil {
		if statusCode == http.StatusNotFound {
			entity = &domain.Todo{}
		} else {
			fmt.Printf("todoCommandHandler.Save ERROR: %+v\n", err)
			return http.StatusInternalServerError, errors.New("internal server error")
		}
	}

	handler.buildEntityFromCommand(entity, command)

	return handler.todoRepository.Upsert(ctx, entity)
}

func (handler *todoCommandHandler) buildEntityFromCommand(entity *domain.Todo, command SaveTodoCommand) {
	entity.Title = command.Title
	entity.Content = command.Content
}

func (handler *todoCommandHandler) GetAll(ctx context.Context) ([]*domain.Todo, int, error) {
	entityList, status, err := handler.todoRepository.FindAll(ctx)
	if err != nil {
		return nil, status, err
	}

	if len(entityList) == 0 {
		return nil, http.StatusNotFound, errors.New("not found error occurred")
	}

	return entityList, status, nil
}
