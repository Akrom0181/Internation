package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user_api_gateway/genproto/user_service"
)

// @Security ApiKeyAuth
// @Router        /LoginAdmin [post]
// @Summary       Admin Login
// @Description   API for admin login
// @Tags          login
// @Accept        json
// @Produce       json
// @Param         login body user_service.LoginPasswors true "Admin credentials"
// @Success       200 {object} user_service.Token
// @Failure       400 {object} models.ResponseError "Invalid request body"
// @Failure       500 {object} models.ResponseError "Internal server error"
func (h *handler) AdministarationLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().AdministarationLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /LoginManager [post]
// @Summary       Manager Login
// @Description   API for manager login
// @Tags          login
// @Accept        json
// @Produce       json
// @Param         login body user_service.LoginPasswors true "Manager credentials"
// @Success       200 {object} user_service.Token
// @Failure       400 {object} models.ResponseError "Invalid request body"
// @Failure       500 {object} models.ResponseError "Internal server error"
func (h *handler) ManagerLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().ManagerLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /LoginStudent [post]
// @Summary        Student Login
// @Description    API for student login
// @Tags           login
// @Accept         json
// @Produce        json
// @Param          login body user_service.LoginPasswors true "Student credentials"
// @Success        200 {object} user_service.Token
// @Failure        400 {object} models.ResponseError "Invalid request body"
// @Failure        500 {object} models.ResponseError "Internal server error"
func (h *handler) StudentLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().StudentLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /LoginSupportTeacher [post]
// @Summary        Support Teacher Login
// @Description    API for support teacher login
// @Tags           login
// @Accept         json
// @Produce        json
// @Param          login body user_service.LoginPasswors true "Support Teacher credentials"
// @Success        200 {object} user_service.Token
// @Failure        400 {object} models.ResponseError "Invalid request body"
// @Failure        500 {object} models.ResponseError "Internal server error"
func (h *handler) SupportTeacherLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().SupportTeacherLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /LoginTeacher [post]
// @Summary        Teacher Login
// @Description    API for teacher login
// @Tags           login
// @Accept         json
// @Produce        json
// @Param          login body user_service.LoginPasswors true "Teacher credentials"
// @Success        200 {object} user_service.Token
// @Failure        400 {object} models.ResponseError "Invalid request body"
// @Failure        500 {object} models.ResponseError "Internal server error"
func (h *handler) TeacherLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().TeacherLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router         /LoginSuperAdmin [post]
// @Summary        Super Admin Login
// @Description    API for super admin login
// @Tags           login
// @Accept         json
// @Produce        json
// @Param          login body user_service.LoginPasswors true "Super Admin credentials"
// @Success 200    {object} user_service.Token
// @Failure 400    {object} models.ResponseError "Invalid request body"
// @Failure 500    {object} models.ResponseError "Internal server error"
func (h *handler) SuperAdminLogin(c *gin.Context) {
	var req user_service.LoginPasswors

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err := h.grpcClient.LoginService().SuperAdminLogin(c, &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}
