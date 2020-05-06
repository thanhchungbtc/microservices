package database

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticketing/auth/model"
	"ticketing/auth/services"
)

type UserRepository struct {
	db *Database
}

func (repo *UserRepository) Exists(email string) (result bool, err error) {
	collection := repo.db.Collection("users")
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

func (repo *UserRepository) Create(user model.User) (*model.User, error) {
	collection := repo.db.Collection("users")
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

func (repo *UserRepository) List() (users []model.User, err error) {
	collection := repo.db.Collection("users")

	cur, err := collection.Find(context.Background(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user model.User
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

func (repo *UserRepository) Get(email string) (*model.User, error) {
	collection := repo.db.Collection("users")
	cur := collection.FindOne(context.Background(), bson.M{
		"email": email,
	})
	if err := cur.Err(); err != nil {
		return nil, err
	}
	var user model.User
	if err := cur.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetJWT(user *model.User) (string, error) {
	type claims struct {
		*model.User
		jwt.StandardClaims
	}
	cl := claims{
		user,
		jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}
