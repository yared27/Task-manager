package data

import (
	"context"
	"errors"
	"task_manager/models"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrorNotFound = errors.New("Task not found")

type TaskService struct {
	Collection *mongo.Collection
}

func  NewTaskService(col *mongo.Collection) *TaskService{
	return &TaskService{Collection: col}

}
func (s *TaskService)ListTasks() ([]models.Task,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor,err := s.Collection.Find(ctx, bson.M{})
	if err != nil{
		return nil,err
	}
	defer cursor.Close(ctx)
	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetTask(id string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// parse the id as an integer and lookup in the in-memory store
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return models.Task{}, err
	}

	var task models.Task
	err = s.Collection.FindOne(ctx, bson.M{"_id":objectID}).Decode(&task)

	return task, err
}

func (s *TaskService) CreateTask(task models.Task) (primitive.ObjectID, error){
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	
	defer cancel()
	now := primitive.NewDateTimeFromTime(time.Now())
	task.CreatedAt = now
	task.UpdatedAt = now

	task.ID = primitive.NilObjectID
	result, err := s.Collection.InsertOne(ctx, task)

	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("failed to get inserted ID")
	}
	return oid, nil
}

func (s *TaskService)UpdateTask(id string, updatedTask *models.Task)error{
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	updatedTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updated := bson.M{
			"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"dueDate":    updatedTask.DueDate,
			"status":      updatedTask.Status,
			"updatedAt":  updatedTask.UpdatedAt,
		},
	}
	result, err := s.Collection.UpdateOne(ctx, bson.M{"_id": objectID}, updated)

	if err != nil {	
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}


func (s *TaskService)DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := s.Collection.DeleteOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}