package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"net/http"
	"strings"
	"time"
	configuration "todo-app/appconfig"
	"todo-app/internal/todo-app/domain"
)

type ITodoRepository interface {
	Upsert(ctx context.Context, entity *domain.Todo) (int, error)
	FindAll(ctx context.Context) ([]*domain.Todo, int, error)
	FindById(ctx context.Context, todoId string) (*domain.Todo, int, error)
	FindByIds(ctx context.Context, todoIds []string) ([]*domain.Todo, int, error)
}

type todoRepository struct {
	todoBucket  *gocb.Bucket
	todoCluster *gocb.Cluster
}

func NewTodoRepository(cluster *gocb.Cluster) ITodoRepository {
	return &todoRepository{
		todoBucket:  cluster.Bucket(configuration.CouchbaseBucket),
		todoCluster: cluster,
	}
}

func (r *todoRepository) Upsert(ctx context.Context, entity *domain.Todo) (int, error) {
	_, err := r.todoBucket.DefaultCollection().Upsert(entity.Title, entity, nil)

	if err != nil {
		fmt.Printf("context: %+v - todoRepository.Upsert Error: %+v", ctx, err)
		return http.StatusInternalServerError, errors.New("error occurred on upsert operation")
	}
	return http.StatusOK, nil
}

func (r *todoRepository) FindAll(ctx context.Context) ([]*domain.Todo, int, error) {
	queryStr := fmt.Sprintf("SELECT t.* FROM `%s` t", r.todoBucket.Name())

	queryResult, err := r.todoCluster.Query(queryStr,
		&gocb.QueryOptions{Timeout: time.Second * 10, Adhoc: true, Readonly: true})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			fmt.Printf("todoRepository.FindAll INFO: %+v\n", err)
			return nil, http.StatusOK, nil
		}

		fmt.Printf("todoRepository.FindAll context: %+v ERROR : %#v\n", ctx, err)
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("internal server error occurred"))
	}

	var todos []*domain.Todo

	for queryResult.Next() {
		var todo domain.Todo
		err = queryResult.Row(&todo)

		if err == nil {
			todos = append(todos, &todo)
		}
	}

	return todos, http.StatusOK, nil
}

func (r *todoRepository) FindById(ctx context.Context, todoId string) (*domain.Todo, int, error) {
	var entity domain.Todo

	queryResult, err := r.todoBucket.DefaultCollection().Get(todoId,
		&gocb.GetOptions{Timeout: time.Second * 1})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			fmt.Printf("todoRepository.FindById INFO: %#v\n", err)
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found error occurred for giving id: %s", todoId))
		}

		fmt.Printf("todoRepository.FindById context: %+v ERROR : %#v\n", ctx, err)
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("internal server error occurred for giving id: %s", todoId))
	}

	err = queryResult.Content(&entity)

	if err != nil {
		fmt.Printf("todoRepository.FindById context: %+v Content ERROR : %#v\n", ctx, err)
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("internal server error occurred for giving id: %s", todoId))
	}

	return &entity, http.StatusOK, nil
}

func (r *todoRepository) FindByIds(ctx context.Context, todoIds []string) ([]*domain.Todo, int, error) {

	queryStr := fmt.Sprintf("SELECT t.* FROM `%s` t USE KEYS [\"%s\"]", r.todoBucket.Name(), strings.Join(todoIds, "\",\""))

	queryResult, err := r.todoCluster.Query(queryStr,
		&gocb.QueryOptions{Timeout: time.Second * 10, Adhoc: true, Readonly: true})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			fmt.Printf("todoRepository.FindByIds INFO: %+v\n", err)
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("not found error occurred"))
		}

		fmt.Printf("todoRepository.FindByIds context: %+v ERROR : %#v\n", ctx, err)
		return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("internal server error occurred"))
	}

	var todos []*domain.Todo

	for queryResult.Next() {
		var todo domain.Todo
		err = queryResult.Row(&todo)

		if err == nil {
			todos = append(todos, &todo)
		}
	}

	if len(todos) == 0 {
		fmt.Printf("todoRepository.FindByIds context: %+v not found delivery by given ids:[\"%s\"]\n", strings.Join(todoIds, "\",\""))

		return nil, http.StatusOK, nil
	}

	return todos, http.StatusOK, nil
}
