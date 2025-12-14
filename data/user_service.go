package data

import (
	"context"
	"errors"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrorUserNotFound = errors.New("user not found")

type UserService struct {
	Collection *mongo.Collection
}

// Constructor
func NewUserService(col *mongo.Collection) *UserService {
	return &UserService{Collection: col}
}

// ---------------- REGISTER ----------------
func (s *UserService) Register(username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user exists
	count, err := s.Collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// First user becomes admin
	totalUsers, err := s.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	role := "user"
	if totalUsers == 0 {
		role = "admin"
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Password:  string(hashed),
		Role:      role,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = s.Collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = "" // never return hash
	return &user, nil
}

// ---------------- LOGIN ----------------
func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	user.Password = ""
	return &user, nil
}

// ---------------- PROMOTE ----------------
func (s *UserService) PromoteUser(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.Collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrorUserNotFound
	}

	return nil
}
