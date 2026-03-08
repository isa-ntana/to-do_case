package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/isa-ntana/to-do_case/pkg/errors"
)

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Data interface{} `json:"data"`
}

func respondError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		c.JSON(appErr.Code, errorResponse{Error: appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, errorResponse{
		Error: "erro interno do servidor",
	})
}

func respondSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, successResponse{Data: data})
}
