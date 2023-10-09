package user

import (
	"github.com/ayo-ajayi/teamsync/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepo struct {
	db db.IDatabase
}

func NewUserRepo(db db.IDatabase) *UserRepo {
	return &UserRepo{db: db}
}

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	Firstname  string             `json:"firstname" bson:"firstname"`
	Lastname   string             `json:"lastname" bson:"lastname"`
	Role       Role               `json:"role" bson:"role"`
	IsVerified bool               `json:"is_verified" bson:"is_verified"`
	PositionID primitive.ObjectID `json:"position_id" bson:"position_id"` 
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Role string

const (
	Admin  Role = "admin"
	Manager Role = "manager"
	Staff  Role = "staff"
)

func (ur *UserRepo) UserIsExists(filter interface{}) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	err := ur.db.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ur *UserRepo) CreateUser(user *User) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := ur.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUser(filter interface{}) (*User, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var user User
	err := ur.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) UpdateUser(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := ur.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUsers(filter interface{}) ([]*User, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var users []*User
	cursor, err := ur.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

type IUserRepo interface {
	CreateUser(user *User) error
	GetUser(filter interface{}) (*User, error)
	UpdateUser(filter interface{}, update interface{}) error
	GetUsers(filter interface{}) ([]*User, error)
	UserIsExists(filter interface{}) (bool, error)
}

