package team

import (
	"github.com/ayo-ajayi/teamsync/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Team struct {
	Id          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	ManagerId   primitive.ObjectID   `json:"manager_id" bson:"manager_id"`
	StaffIds    []primitive.ObjectID `json:"staff_ids" bson:"staff_ids"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

type ITeamRepo interface {
	CreateTeam(team *Team) error
	UpdateTeam(filter interface{}, update interface{}) error
	GetTeam(filter interface{}) (*Team, error)
	GetTeams(filter interface{}) ([]*Team, error)
}

type TeamRepo struct {
	db db.IDatabase
}

func NewTeamRepo(db db.IDatabase) *TeamRepo {
	return &TeamRepo{db: db}
}

func (tr *TeamRepo) CreateTeam(team *Team) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := tr.db.InsertOne(ctx, team)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TeamRepo) UpdateTeam(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := tr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TeamRepo) GetTeam(filter interface{}) (*Team, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var team *Team
	err := tr.db.FindOne(ctx, filter).Decode(team)
	if err != nil {
		return nil, err
	}
	return team, nil
}


func (tr *TeamRepo) GetTeams(filter interface{}) ([]*Team, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var teams []*Team
	cursor, err := tr.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, err
	}
	return teams, nil
}