package services

import (
	"context"
	"service-book/config"
	"service-book/constants"
	constUserErr "service-book/constants/errors/userErr"
	"service-book/domains/dto"
	"service-book/repositories"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (u *UserService) isUsernameExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) isEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}
func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Role:  strings.ToLower(user.Role.Code),
	}

	claim := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	res := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return res, nil
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	if u.isUsernameExist(ctx, req.Username) {
		return nil, constUserErr.ErrUsernameExists
	}

	if u.isEmailExist(ctx, req.Email) {
		return nil, constUserErr.ErrEmailExists
	}

	if req.Password != req.ConfirmPassword {
		return nil, constUserErr.ErrPaswordDoesNotMatch
	}

	user, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   constants.Admin, // Register admin only
	})
	if err != nil {
		return nil, err
	}

	res := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:     user.UUID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	return res, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:     userLogin.UUID,
		Name:     userLogin.Name,
		Username: userLogin.Username,
		Email:    userLogin.Email,
		Role:     userLogin.Role,
	}

	return &data, nil
}
