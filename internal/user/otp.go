package user

import (
	"time"
	"errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/ayo-ajayi/teamsync/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	URL       string             `json:"url" bson:"url"`
	Secret    string             `json:"secret" bson:"secret"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}


type OTPManager struct {
	db 						   db.IDatabase
	Issuer                          string
	SignUpOtpValidityInSecs         uint
	ForgotPasswordOtpValidityInSecs uint
}


type IOTPManager interface {
	GenerateSignUpOTP(email string) (string, error)
	GenerateForgotPasswordOTP(email string) (string, error)
}

func NewOTPManager(db db.IDatabase, issuer string, signUpOtpValidityInSecs uint, forgotPasswordOtpValidityInSecs uint) *OTPManager {
	return &OTPManager{db: db, Issuer: issuer, SignUpOtpValidityInSecs: signUpOtpValidityInSecs, ForgotPasswordOtpValidityInSecs: forgotPasswordOtpValidityInSecs}
}

func (om *OTPManager) createOTP(email string, otpValidityInSecs uint) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      om.Issuer,
		AccountName: email,
		Period:      otpValidityInSecs,
	})
	if err != nil {
		return "", err
	}
	err = om.saveOTP(&OTP{
		Email:     email,
		URL:       key.URL(),
		Secret:    key.Secret(),
		ExpiresAt: time.Now().Add(time.Duration(otpValidityInSecs) * time.Second),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}
	otp, err := getOTPFromSecret(key.Secret(), otpValidityInSecs, time.Now())
	if err != nil {
		return "", err
	}
	if otp == "" {
		return "", errors.New("otp is empty")
	}
	return otp, nil
}


func (om *OTPManager) saveOTP(otp *OTP) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := om.db.InsertOne(ctx, otp)
	if err != nil {
		return err
	}
	return nil
}

func getOTPFromSecret(secret string, otpValidityInSecs uint, t time.Time) (string, error) {
	return totp.GenerateCodeCustom(secret, t, totp.ValidateOpts{
		Period:    otpValidityInSecs,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
}

func (om *OTPManager)GenerateSignUpOTP(email string) (string, error) {
	return om.createOTP(email, om.SignUpOtpValidityInSecs)
}

func (om *OTPManager)GenerateForgotPasswordOTP(email string) (string, error) {
	return om.createOTP(email, om.ForgotPasswordOtpValidityInSecs)
}