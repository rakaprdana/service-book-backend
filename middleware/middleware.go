package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"service-book/common/responses"
	"service-book/config"
	"service-book/constants"
	constErr "service-book/constants/errors"
	"service-book/services"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recoverd from panic error: %v", r)
				ctx.JSON(http.StatusInternalServerError, responses.Response{
					Status:  constants.Error,
					Message: constErr.ErrInternalServerError,
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusTooManyRequests, responses.Response{
				Status:  constants.Error,
				Message: constErr.ErrTooManyRequests.Error(),
			})
			ctx.Abort()
		}
		ctx.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")

	if len(arrayToken) == 2 {
		return arrayToken[1]
	}

	return ""
}

func responseUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, responses.Response{
		Status:  constants.Error,
		Message: message,
	})
	ctx.Abort()
}

func validateAPIKey(ctx *gin.Context) error {
	apiKey := ctx.GetHeader(constants.XApiKey)
	requestAt := ctx.GetHeader(constants.XRequestAt)
	serviceName := ctx.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))

	resultHash := hex.EncodeToString(hash.Sum(nil))
	if apiKey != resultHash {
		return constErr.ErrUnauthorized
	}

	return nil
}

func validateBearerToken(ctx *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		return constErr.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		return constErr.ErrUnauthorized
	}

	claims := &services.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, constErr.ErrInvalidToken
		}

		jwtSecret := []byte(config.Config.JwtSecretKey)
		return jwtSecret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return constErr.ErrUnauthorized
	}

	userLogin := ctx.Request.WithContext(context.WithValue(
		ctx.Request.Context(),
		constants.UserLogin,
		claims.User,
	))

	ctx.Request = userLogin
	ctx.Set(constants.Token, token)
	return nil
}

func Authenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		token := ctx.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(ctx, constErr.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(ctx, token)
		if err != nil {
			responseUnauthorized(ctx, err.Error())
			return
		}
		ctx.Next()
	}
}
