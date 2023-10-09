package user

import (
	"errors"
	"time"

	"github.com/ayo-ajayi/teamsync/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessDetails struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	AccessUuid string             `json:"access_uuid" bson:"access_uuid"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id"`
	ExpireAt   time.Time          `json:"expire_at" bson:"expire_at"`
}

type TokenManager struct {
	accessTokenSecret          string
	accessTokenValidityInHours int64
	db                         db.IDatabase
}

type TokenDetails struct {
	AccessToken string `json:"access_token"`
	AcessUuid   string `json:"-"`
	AtExpires   int64  `json:"at_expires"`
}

func NewTokenManager(accessTokenSecret string, accessTokenValidityInHours int64, db db.IDatabase) *TokenManager {
	return &TokenManager{accessTokenSecret: accessTokenSecret, accessTokenValidityInHours: accessTokenValidityInHours, db: db}
}

func createAccessToken(userId primitive.ObjectID, uuid string, expires int64, secret string)(string, error){
	claims:=jwt.MapClaims{}
	claims["user_id"]=userId
	claims["access_uuid"]=uuid
	claims["exp"]=expires
	claims["authorized"] = true
	at:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(secret))
}

func (tm *TokenManager)GenerateAccessToken(userId primitive.ObjectID)(*TokenDetails, error){
	td:=&TokenDetails{}
	td.AtExpires= time.Now().Add(time.Hour * time.Duration(tm.accessTokenValidityInHours)).Unix()
	td.AcessUuid=uuid.New().String()
	
	
	accessToken, err:=createAccessToken(userId, td.AcessUuid, td.AtExpires, tm.accessTokenSecret)
	if err != nil {
		return nil, err
	}
	if accessToken == ""{
		return nil, errors.New("access token is empty")
	}
	td.AccessToken=accessToken
	return td, nil
}

func (tm *TokenManager)SaveAccessToken(userId primitive.ObjectID, td *TokenDetails)error{
	ctx, cancel:= db.DBReqContext(5)
	defer cancel()
	_, err:=tm.db.InsertOne(ctx, &AccessDetails{
		AccessUuid: td.AcessUuid,
		UserId: userId,
		ExpireAt: time.Unix(td.AtExpires, 0),
	})
	return err
}

func (tm *TokenManager)DeleteAccessToken(userId primitive.ObjectID)error{
	ctx, cancel:= db.DBReqContext(5)
	defer cancel()
	_, err:=tm.db.DeleteOne(ctx, bson.M{"user_id": userId})
	return err
}

func (tm *TokenManager)FindAccessToken(uuid string)(*AccessDetails, error){
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var accessDetails AccessDetails
	err:=tm.db.FindOne(ctx, bson.M{"access_uuid": uuid}).Decode(&accessDetails)
	if err != nil {
		return nil, err
	}
	return &accessDetails, nil
}


func (tm *TokenManager)ValidateToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}

func (tm *TokenManager) ExtractTokenMetadata(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized")
	}
	accessUuid, ok := claims["access_uuid"].(string)
	if !ok || accessUuid == "" {
		return nil, errors.New("unauthorized")
	}
	userId, ok := claims["user_id"].(primitive.ObjectID)
	if !ok || userId == primitive.NilObjectID {
		return nil, errors.New("unauthorized")
	}
	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}


type IMiddlewareTokenManager interface {
	ExtractTokenMetadata(token *jwt.Token) (*AccessDetails, error)
	ValidateToken(token string, secret string) (*jwt.Token, error)
	FindAccessToken(uuid string)(*AccessDetails, error)
}