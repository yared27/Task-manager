package data
import (
	"errors"
	"strconv"
	"sync"
	"task_manager/models"
	"time"
)

var ErrorNotFound = errors.New("Task not found")

type service struct {
	mu sync.RWMutex
	store map[int]models.Task
	next int
}

var svc = &service{
	store:make(map[int]models.Task),
	next:1,
}
func ListTasks() []models.Task {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	tasks := make([]models.Task,0, len(svc.store))

	for _, t := range svc.store{
		tasks = append(tasks,t)
	}

	return  tasks
}

func GetTask(id int)(models.Task, error){
	svc.mu.Lock()
	defer svc.mu.Unlock()
	t, ok := svc.store[id]
	if !ok{
		return models.Task{}, ErrorNotFound
	}
	return  t,nil
}


func CreateTask(title, description string, due time.Time, status string) models.Task{
	svc.mu.Lock()
	defer svc.mu.Unlock()
	t := models.Task{
		ID:  strconv.Itoa(svc.next),
		Title: title,
		Description: description,
		DueDate: due,
		Status: status,
	}
	svc.store[svc.next] = t
	svc.next++
	return t
}

func UpdateTask(id int, title, description string, due time.Time, status string) (models.Task, error){
	svc.mu.Lock()
	defer svc.mu.Unlock()
	_,ok := svc.store[id]

	if !ok {
		return models.Task{}, ErrorNotFound
	}

	t := models.Task{
		ID: strconv.Itoa(id),
		Title: title,
		Description: description,
		DueDate: due,
		Status: status,
	}
	svc.store[id] = t
	return t,nil

}

func DeleteTask(id int) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	_,ok := svc.store[id]

	if !ok{
		return ErrorNotFound
	}

	delete(svc.store, id)
	return nil
}