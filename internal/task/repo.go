package task

import "go.mongodb.org/mongo-driver/bson/primitive"
import "time"
import 	"github.com/ayo-ajayi/teamsync/internal/db"

type Task struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status TaskStatus `json:"status" bson:"status"`
	TaskType  TaskType `json:"task_type" bson:"task_type"`
	Description string `json:"description" bson:"description"`
	AssignedBy   primitive.ObjectID `json:"assigned_by" bson:"assigned_by"`
	AssignedTo   []primitive.ObjectID `json:"assigned_to" bson:"assigned_to"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type TaskType string

const (
	TeamTask       TaskType = "team"
	IndividualTask TaskType = "individual"
)

//create personal task by setting assigned by and assigned to oneself

type TaskStatus string

const (
	Completed  TaskStatus = "completed"
	InProgress TaskStatus = "inprogress"
	NotStarted TaskStatus = "notstarted"
	Accepted   TaskStatus = "accepted"
	Rejected   TaskStatus = "rejected"
)

// manager can create task

type TaskRepo struct{
	db db.IDatabase
}

func NewTaskRepo(db db.IDatabase) *TaskRepo {
	return &TaskRepo{db: db}
}

func (tr *TaskRepo) CreateTask(task *Task) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := tr.db.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskRepo) UpdateTask(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_,err := tr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskRepo) GetTask(filter interface{}) (*Task, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var task *Task
	err := tr.db.FindOne(ctx, filter).Decode(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *TaskRepo) GetTasks(filter interface{}) ([]*Task, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var tasks []*Task
	cursor, err := tr.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

type ITaskRepo interface {
	CreateTask(task *Task) error 
	UpdateTask(filter interface{}, update interface{}) error 
	GetTask(filter interface{}) (*Task, error) 
	GetTasks(filter interface{}) ([]*Task, error)
}