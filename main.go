package main

import (
	"fmt"
	"log"
	"net/http"
	authenticate "onlinetest/auth/controllers"
)

//var userCollection *mongo.Collection

func main() {
	//for connecting mongoDB
	// opts := options.Client().ApplyURI("mongodb+srv://onlinetestvsnco:columbus@cluster0.upwds9z.mongodb.net/?retryWrites=true&w=majority")
	// client, err := mongo.Connect(context.Background(), opts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer client.Disconnect(context.Background())

	// userCollection = client.Database("admin").Collection("users")
	// if userCollection == nil {
	// 	log.Fatal("Failed to initialize user collections")
	// }

	//handlers
	http.HandleFunc("/register", authenticate.Register)
	http.HandleFunc("/login", authenticate.Login)

	port := 3000
	fmt.Printf("Server is running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		fmt.Sprintln(err, "error with the port")
		log.Fatal(err)
	}
}
