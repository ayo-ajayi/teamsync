package position

import (
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PositionService struct {
	repo IPositionRepo
}

func NewPositionService(repo IPositionRepo) *PositionService {
	return &PositionService{repo: repo}
}

func (ps *PositionService) CreatePosition(position *Position) error {

	position.CreatedAt = time.Now()
	position.UpdatedAt = time.Now()
	err:=ps.repo.CreatePosition(position)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PositionService) GetPosition(id primitive.ObjectID) (*Position, error) {
	position, err:=ps.repo.GetPosition(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return position, nil
}

func (ps *PositionService) UpdatePosition(position *Position) error {
	oldPosition, err:=ps.repo.GetPosition(bson.M{"_id": position.Id})
	if err != nil {
		return err
	}
	position.CreatedAt = oldPosition.CreatedAt
	position.UpdatedAt = time.Now()
	err = ps.repo.UpdatePosition(bson.M{"_id": position.Id}, bson.M{"$set": position})
	if err != nil {
		return err
	}
	return nil
}

func (ps *PositionService) GetPositions() ([]*Position, error) {
	positions, err:=ps.repo.GetPositions(bson.M{})
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (ps *PositionService) DeletePosition(id primitive.ObjectID) error {
	err:=ps.repo.DeletePosition(bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
