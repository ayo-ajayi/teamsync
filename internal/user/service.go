package user

import (
	"errors"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	otpManager   IOTPManager
	userRepo     IUserRepo
	emailManager IEmailManager
}

func NewUserService(userRepo IUserRepo, otpManager IOTPManager) *UserService {
	return &UserService{userRepo: userRepo, otpManager: otpManager}
}

func (us *UserService) SignUpAndSendVerifyEmail(user *User) error {
	var wg sync.WaitGroup
	var userExists bool
	var passwordHash string
	var otp string

	errCh := make(chan error, 3)
	wg.Add(1)
	go func() {
		defer wg.Done()
		exists, err := us.userRepo.UserIsExists(bson.M{"email": user.Email})
		if err != nil {
			errCh <- err
			return
		}
		userExists = exists
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		hash, err := HashPassword(user.Password)
		if err != nil {
			errCh <- err
			return
		}
		if hash == "" {
			errCh <- errors.New("password hash is empty")
			return
		}
		passwordHash = hash
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		genOtp, err := us.otpManager.GenerateSignUpOTP(user.Email)
		if err != nil {
			errCh <- err
			return
		}
		if genOtp == "" {
			errCh <- errors.New("otp is empty")
			return
		}
		otp = genOtp

	}()
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	if userExists {
		return errors.New("user already exists")
	}
	user.Password = passwordHash
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Email != os.Getenv("ADMIN_EMAIL") {
		user.Role = Staff
	}
	if err := us.userRepo.CreateUser(user); err != nil {
		return err
	}
	if err := us.emailManager.SendSignUpOTP(user.Email, user.Firstname, otp); err != nil {
		return err
	}
	return nil

}

func (us *UserService) Profile(id primitive.ObjectID) (*User, error) {
	user, err := us.userRepo.GetUser(bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUsers(currentUser primitive.ObjectID) ([]*string, error) {
	users, err := us.userRepo.GetUsers(bson.M{})
	if err != nil {
		return nil, err
	}
	var emails []*string
	for _, user := range users {
		if user.Id == currentUser {
			continue
		}
		emails = append(emails, &user.Email)
	}
	return emails, nil
}

func (us *UserService) UpdateUserPosition(id primitive.ObjectID, positionID primitive.ObjectID) error {
	user, err := us.userRepo.GetUser(bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}
	//check if position exists
	user.PositionID = positionID
	user.UpdatedAt = time.Now()
	if err := us.userRepo.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": user}); err != nil {
		return err
	}
	return nil
}


func (us *UserService)UpdateUserRole(id primitive.ObjectID, role Role) error {
	user, err := us.userRepo.GetUser(bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}
	user.Role = role
	user.UpdatedAt = time.Now()
	if err := us.userRepo.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": user}); err != nil {
		return err
	}
	return nil
}