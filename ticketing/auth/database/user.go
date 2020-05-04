package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticketing/auth/services"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *Database) CreateUser(user User) (*User, error) {
	collection := db.Collection("users")
	passwordService := services.NewPassword()
	user.Password = passwordService.ToHash(user.Password)
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v", result.InsertedID)
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &user, err
}

func (db *Database) IsUserExists(email string) (result bool, err error) {
	collection := db.Collection("users")
	cur := collection.FindOne(context.Background(), bson.M{
		"email": email,
	})
	if err := cur.Err(); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (db *Database) GetUsers() (users []User, err error) {
	collection := db.Collection("users")

	cur, err := collection.Find(context.Background(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
