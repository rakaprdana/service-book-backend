package controllers

import (
	"net/http"
	commErr "service-book/common/errors"
	"service-book/common/responses"
	"service-book/domains/dto"
	"service-book/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type RepairJobController struct {
	service services.IServiceRegistry
}

type IRepairJobController interface {
	GetAllWithoutPagination(*gin.Context)
	Create(*gin.Context)
}

func NewRepairJobController(service services.IServiceRegistry) IRepairJobController {
	return &RepairJobController{
		service: service,
	}
}

func (r *RepairJobController) GetAllWithoutPagination(ctx *gin.Context) {
	result, err := r.service.GetRepairJob().GetAllWithoutPagination(ctx)
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
		Data: result,
		Gin:  ctx,
	})
}

func (r *RepairJobController) Create(ctx *gin.Context) {
	req := &dto.RepairJobRequest{}
	err := ctx.ShouldBindWith(req, binding.FormMultipart)
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
			Err:     err,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})
		return
	}

	result, err := r.service.GetRepairJob().Create(ctx, req)
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
		Err:  err,
		Data: result,
		Gin:  ctx,
	})
}
