package database

import (
	model "Backend/Model"
	models "Backend/Model"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	
	client,err 	:= mongo.NewClient(options.Client().ApplyURI(mongoUrl))


	if err!=nil{
		return err;
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second);

	defer cancel();

	err = client.Connect(ctx)
	if err!=nil{
		return err;
	}

	log.Println("Mongodb connected successfully!");
	return nil;
}

func Close() error{
	if client==nil{
		return nil;
	}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second);

	err := client.Disconnect(ctx)

	defer cancel();

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

func FetchUserById(uuid string,collectionName string) (*models.User,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var foundUser models.User

	var collection *mongo.Collection = GetCollectionByName(collectionName);
	err := collection.FindOne(ctx,bson.M{"user_id":uuid}).Decode(&foundUser);

	if err!=nil{
		defer cancel();
		log.Println("Could not fetch the user by id!");
		return nil,err;
	}

	defer cancel();

	return &foundUser,nil
}

func CheckRefreshTokenAlreadyexists(collectionName string,refreshToken string) (bool,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var collection *mongo.Collection = GetCollectionByName(collectionName);


	count,err := collection.CountDocuments(ctx,bson.M{"refresh_token":refreshToken})
	defer cancel();

	if err!=nil{
		log.Println("CheckRefreshTokenAlreadyexists->something went wrong 1");
		return false,err;
	}

	if count>0{
		return true,nil;
	}

	return false,nil;
}

func SaveRefreshToken(collectionName string,refreshToken string) error{
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);
	var collection *mongo.Collection = GetCollectionByName(collectionName);
	var tokenObj *model.TokenSign;

	err := collection.FindOne(ctx,bson.M{"refresh_token":refreshToken}).Decode(&tokenObj)

	if err!=nil{
		defer cancel();
		log.Println("SaveRefreshToken->The token already exists!");
		log.Println(refreshToken);
		return err;
	}

	result,err := collection.InsertOne(ctx,tokenObj);

	if err!=nil{
		defer cancel();
		log.Println("SaveRefreshToken->Could not save refresh token!")
		return err;
	}

	defer cancel();
	log.Println("The refresh token saved successfully ,",result);

	return nil;
}

func DeleteRefreshToken(collectionName string ,refreshToken string)error{
	
}


func FetchUserByName(userName string,collectionName string)(*models.User,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var foundUser models.User;
	var collection *mongo.Collection = GetCollectionByName(collectionName);

	err := collection.FindOne(ctx,bson.M{"user_name":userName}).Decode(&foundUser);
	
	defer cancel();

	if err!=nil{
		log.Println("FetchUserByName->the username does not exists");
		return nil,err;
	}

	return &foundUser,nil;
}