package database

import (
	model "Backend/Model"
	"context"
	"log"
	"os"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
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

func FetchRefreshTokenForAuthtokenUpdate(collectionName,user_name,email,last_name string,)(string,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var fetchedRefreshToken string;

	var collection *mongo.Collection = GetCollectionByName(collectionName);

	filter := bson.M{
		"user_name":user_name,
		"last_name":last_name,
		"email":email,
	}

	var result bson.M;
	err := collection.FindOne(ctx,filter).Decode(&result)
	defer cancel();

	if err!=nil{
		log.Println("FetchRefreshTokenForAuthtokenUpdate->err1")
		log.Fatal(err);
		return "",err;
	}

	fetchedRefreshToken = result["refresh_token"].(string)

	return fetchedRefreshToken,nil;
}

func SetRefreshTokenAsEmpty(collectionName,last_name,user_name,email string,) (error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var collection *mongo.Collection = GetCollectionByName(collectionName);

	filter := bson.M{"user_name":user_name,"last_name":last_name,"email":email}
	update := bson.M{"$set":bson.M{"refresh_token":""}}

	_,err := collection.UpdateOne(ctx,filter,update)

	defer cancel();

	return err;
	
}

func FetchUserCountByCredential(collectionName,email,user_name string,)(int64,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var collection *mongo.Collection = GetCollectionByName(collectionName);

	filter := bson.M{
		"$or":[]interface{}{
			bson.M{"user_name":user_name},
			bson.M{"email":email},
		},
	}

	count,err := collection.CountDocuments(ctx,filter,)

	defer cancel();

	if err!=nil{
		log.Panic(err)
		return -1,err;
	}	

	log.Println("Count of documents is ",count)
	return count,nil;
}

func DeleteUserByCredentials(collectionName,user_name,last_name,email string)(error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);

	var collection *mongo.Collection = GetCollectionByName(collectionName)
	filter := bson.M{"user_name":user_name,"last_name":last_name,"email":email};

	deleteResult,err := collection.DeleteOne(ctx,filter,);
	defer cancel();

	if err!=nil{
		log.Println("DeleteUserByCredentials -> could not delete user!");
		log.Panic(err);
		return err;
	}

	log.Println("The number of deleted documents is ",deleteResult.DeletedCount)
	return nil;
}

func SaveUserCredential(collectionName string,newUser model.User,)(interface{},error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

	var collection *mongo.Collection = GetCollectionByName(collectionName)

	result,err := collection.InsertOne(ctx,newUser)

	defer cancel();

	if err!=nil{
		log.Println("Could not save new user ,in signup!");
		log.Fatal(err)
		return nil,err
	}
	
	log.Println("Inserted ID : ",result.InsertedID);

	return result.InsertedID,nil;
}

func FetchUserBySpecificCredential(collectionName,credential string)(*model.User,error){
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

	var user model.User

	var collection *mongo.Collection = GetCollectionByName(collectionName)
	err := collection.FindOne(ctx,bson.M{"email":credential}).Decode(&user)

	defer cancel();

	if err!=nil{
		log.Panic(err)
		return nil,err;
	}

	return &user,nil;
}
/*
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
	var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second);


	 var collection *mongo.Collection = GetCollectionByName(collectionName);

	 result,err := collection.DeleteOne(ctx,bson.M{"refresh_token":refreshToken},);
	defer cancel();

	 if err!=nil{
		log.Println("DeleteRefreshToken-> something went wrong 1");
		log.Fatal(err);
		return err;
	 }

	 if result.DeletedCount==0{
		log.Println("Refresh token Document not found to delete");
	 }else{
		log.Println("Refresh token deleted");
	 }

	 return nil;
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
*/