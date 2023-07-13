package middleware 

import (
	"golang-react-todo/models"
	"net/http"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init(){
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv(){
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("error loading the .env file")
	}
}

func createDBInstance(){
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)

	client,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(),nil)
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	collection = client.Database(dbName).Collection(collName)
	fmt.Println("collection instance created")

}

func GetAllTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*") // Cors
	w.Header().Set("Access-Control-Allow-Methods","GET")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	
	var task models.ToDoList
	json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}

func TaskComplete(w http.ResponseWriter, r *http.Request){
	
	w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	taskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func UndoTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := mux.Vars(r)
	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request){
	// w.Header().Set("Context-Type","application/x-www-from-urlencoded")
	// w.Header().Set("Access-Control-Allow-Origin","*")
	// w.Header().Set("Access-Control-Allow-Methods","DELETE")
	// w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	
	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)

}

// --------------------------------------------------------------------

func getAllTasks() []primitive.M {
	cur,err := collection.Find(context.Background(),bson.D{{}})
	if err!=nil{
		log.Fatal(err)
	}
	var results []primitive.M
	for cur.Next(context.Background()){
		var result bson.M
		err := cur.Decode(&result)
		if err!=nil {
			log.Fatal(err)
		}
		results = append(results,result)
	}
	if err:=cur.Err(); err!=nil{
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return results	
}

func taskComplete(task string){
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":true}}
	result,err := collection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("modified count:",result.ModifiedCount)
}

func insertOneTask(task models.ToDoList){
	if task.Task ==""{
		return
	}
	fmt.Println(task.Task)
	insertedResult,err := collection.InsertOne(context.Background(),task)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("insrted a single record with id", insertedResult.InsertedID)
}

func undoTask(task string){
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":false}}
	result,err := collection.UpdateOne(context.Background(),filter,update)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("modified count:",result.ModifiedCount)
}

func deleteOneTask(task string){
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	d,err := collection.DeleteOne(context.Background(),filter)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Deleted Document",d)
}

func deleteAllTasks() int64 {
	d,err := collection.DeleteMany(context.Background(),bson.D{{}},nil)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("deleted document",d)
	return d.DeletedCount
}

