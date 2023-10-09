package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	accessTokenSecret string
	userRepo          IMiddlewareUserRepo
	tokenManager      IMiddlewareTokenManager
}

type IMiddlewareUserRepo interface {
	GetUser(filter interface{}) (*User, error)
}

func NewMiddleWare(accessTokenSecret string, userRepo IMiddlewareUserRepo, tokenManager IMiddlewareTokenManager) *Middleware {
	return &Middleware{
		accessTokenSecret: accessTokenSecret,
		userRepo:          userRepo,
		tokenManager:      tokenManager,
	}
}

func (mw *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := mw.extractToken(c.Request)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: token is required"}})
			c.Abort()
			return
		}
		jwtToken, err := mw.tokenManager.ValidateToken(token, mw.accessTokenSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
				c.Abort()
				return
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: token expired"}})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": "internal server error"}})
			c.Abort()
			return
		}
		td, err := mw.tokenManager.ExtractTokenMetadata(jwtToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		accessDetails, err := mw.tokenManager.FindAccessToken(td.AccessUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		if td.UserId != accessDetails.UserId{
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "unauthorized: invalid token"}})
			c.Abort()
			return
		}
		c.Set("user_id", accessDetails.UserId)
		c.Next()
	}

}


func (mw *Middleware)extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	ttoken := strings.Split(token, " ")
	if len(ttoken) != 2 {
		return ""
	}
	return ttoken[1]
}