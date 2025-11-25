package config
import (
	"context"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
func ConnectDB() (*mongo.Client, error){
	godotenv.Load()
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in .env")
	}



	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx,options.Client().ApplyURI(mongoURI))

	if err != nil {
		return nil,err
		 }
	
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	DB = client
	return client,nil
}

func GetCollection(collectionName string) *mongo.Collection{
	dbName := os.Getenv("DB_NAME")
	return DB.Database(dbName).Collection(collectionName)
}

