package position

import (
	"time"

	"github.com/ayo-ajayi/teamsync/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PositionRepo struct {
	db db.IDatabase
}

func NewPositionRepo(db db.IDatabase) *PositionRepo {
	return &PositionRepo{db: db}
}
type Position struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	SalaryPerHour float64 `json:"salary_per_hour" bson:"salary_per_hour"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"` 
}

func (pr *PositionRepo)CreatePosition(position *Position)error{
	ctx, cancel:=db.DBReqContext(5)
	defer cancel()
	_, err := pr.db.InsertOne(ctx, position)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PositionRepo)GetPosition(filter interface{})(*Position, error){
	ctx, cancel:=db.DBReqContext(5)
	defer cancel()
	var position Position
	err := pr.db.FindOne(ctx, filter).Decode(&position)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (pr *PositionRepo)UpdatePosition(filter interface{}, update interface{})(error){
	ctx, cancel:=db.DBReqContext(5)
	defer cancel()
	_, err:=pr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PositionRepo)GetPositions(filter interface{})([]*Position, error){
	ctx, cancel:=db.DBReqContext(5)
	defer cancel()
	var positions []*Position
	cursor, err:=ur.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &positions)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (pr *PositionRepo)DeletePosition(filter interface{})(error){
	ctx, cancel:=db.DBReqContext(5)
	defer cancel()
	_, err:=pr.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}


type IPositionRepo interface {
	CreatePosition(position *Position)error
	GetPosition(filter interface{})(*Position, error)
	UpdatePosition(filter interface{}, update interface{})(error)
	GetPositions(filter interface{})([]*Position, error)
	DeletePosition(filter interface{})(error)
}