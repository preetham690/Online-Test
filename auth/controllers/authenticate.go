package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"onlinetest/auth/models"
	user_mod "onlinetest/auth/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client
var userCollection *mongo.Collection

// here we are connecting the mongoDB
func init() {
	opts := options.Client().ApplyURI("mongodb+srv://onlinetestvsnco:columbus@cluster0.upwds9z.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(context.Background())

	userCollection = client.Database("admin").Collection("users")
	if userCollection == nil {
		log.Fatal("Failed to initialize user collections")
	}

	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("unimplemented method"))
		return
	}

	u := new(models.User)

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		//glog.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there seems to be some thing went wrong.Please try again or contact admin"))
		return
	}

	err = u.Validate()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	//here we are just hashing the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := user_mod.User{
		User_ID:  primitive.NewObjectID(), // Assign a new ObjectID (if needed)
		Name:     u.Name,
		Email:    u.Email,
		Mobile:   u.Mobile,
		Password: string(hashedPass), //hashed password storing as a string
		IsAdmin:  u.IsAdmin,
	}

	//here we are inserting the data into the database
	//result, err := .DB.Insert(, u)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := userCollection.InsertOne(ctx, user)
	// if err != nil {
	// 	http.Error(w, "ERROR while inserting : "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	fmt.Fprintln(w, "User Registered Successfully")

	defer client.Disconnect(context.Background())

}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user user_mod.User
	//here we are finding the user in the database using his/her name
	err := userCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid Username Or Password", http.StatusUnauthorized)
		return
	}

	//comparing the hashed password with the password entered
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid Username Or Password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Authentication Successful")
}
