package shift

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShiftRepo interface{}

type Shift struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ShiftReport string             `json:"shift_report" bson:"shift_report"`
	HoursWorked int                `json:"hours_worked" bson:"hours_worked"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
