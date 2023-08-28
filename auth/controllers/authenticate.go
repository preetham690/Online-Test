package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	user_mod "onlinetest/auth/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	defer client.Disconnect(context.Background())

	userCollection = client.Database("admin").Collection("users")
	if userCollection == nil {
		log.Fatal("Failed to initialize user collections")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	//username := r.FormValue("username")
	password := r.FormValue("password")
	//isAdmin := false

	//here we are just hashing the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	// var user = new User({
	// 	Username: username,
	// 	Password: string(hashedPass),
	// 	isAdmin:  isAdmin,
	// })

	w.Header().Set("Content-Type", "application/json")

	var user user_mod.User

	json.NewDecoder(r.Body).Decode(&user)
	user.Password = string(hashedPass)

	//here we are inserting the data into the database
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := userCollection.InsertOne(ctx, user)

	// if err != nil {
	// 	http.Error(w, "ERROR while inserting error", http.StatusInternalServerError)
	// 	return
	// }
	json.NewEncoder(w).Encode(result)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User Registered Successfully")

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
