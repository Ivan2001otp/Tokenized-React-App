package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect() error{
	if client!=nil{
		return nil;
	}

	mongoUrl := os.Getenv("MONGO_URL")

	if mongoUrl==""{
		log.Panic("Mongo Url not found!");
		return nil;
	}

	ctx := context.Background();
	client,err 	:= mongo.NewClient(options.Client().ApplyURI(mongoUrl))

	if err!=nil{
		return err;
	}

	err = client.Connect(ctx)
	if err!=nil{
		return err;
	}

	return nil;
}

func Close() error{
	if client==nil{
		return nil;
	}

	ctx := context.Background();

	err := client.Disconnect(ctx)

	if err!=nil{
		return err;
	}
	client = nil;

	return nil;
}

func GetCollectionByName(collectionName string) *mongo.Collection{
	if client==nil{
		log.Println("GetCollectionByName->Mongo client not connected")
		log.Panic("Mongo client not connected");
		return nil;
	}

	var dbName string = os.Getenv("DATABASE_NAME")

	if dbName==""{
		log.Panic("Database Name not found !");
		return nil;
	}
	return client.Database(dbName).Collection(collectionName)
}