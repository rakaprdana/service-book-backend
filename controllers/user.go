package controllers

import (
	"net/http"
	commErr "service-book/common/errors"
	"service-book/common/responses"
	"service-book/domains/dto"
	"service-book/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	GetUserLogin(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{
		service: service,
	}
}

func (u *UserController) Login(ctx *gin.Context) {
	req := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(req)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := commErr.ErrValidationResponse(err)
		responses.HttpResponse(responses.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Login(ctx, req)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	responses.HttpResponse(responses.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

func (u *UserController) Register(ctx *gin.Context) {
	req := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(req)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := commErr.ErrValidationResponse(err)
		responses.HttpResponse(responses.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Register(ctx, req)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	responses.HttpResponse(responses.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResp{
			Code: http.StatusOK,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	responses.HttpResponse(responses.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})

}
