package handler

import (
	"errors"
	"net/http"
	"strconv"
	"user_api_gateway/genproto/schedule_service"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router        /CreateJournal [post]
// @Summary       Create journal
// @Description   API for creating journal
// @Tags          journal
// @Accept        json
// @Produce       json
// @Param         journal body schedule_service.CreateJournal true "Journal"
// @Success       200 {object} schedule_service.GetJournal
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) CreateJournal(c *gin.Context) {
	var (
		req  schedule_service.CreateJournal
		resp *schedule_service.GetJournal
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	resp, err = h.grpcClient.JournalService().Create(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to create journal")
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /GetJournal/{id} [get]
// @Summary       Get a journal by ID
// @Description   API for getting a journal by ID
// @Tags          journal
// @Accept        json
// @Produce       json
// @Param         id path string true "Journal ID"
// @Success       200 {object} schedule_service.GetJournal
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) GetJournalByID(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp *schedule_service.GetJournal
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req := &schedule_service.JournalPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.JournalService().GetByID(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get journal")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /GetListJournal [get]
// @Summary       Get list of journals
// @Description   API for getting list of journals
// @Tags          journal
// @Accept        json
// @Produce       json
// @Param         search query string false "Search"
// @Param         page query int false "Page"
// @Param         limit query int false "Limit"
// @Success       200 {object} schedule_service.GetListJournalResponse
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) GetListJournal(c *gin.Context) {
	var (
		req  schedule_service.GetListJournalRequest
		resp *schedule_service.GetListJournalResponse
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing page")
		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while parsing limit")
		return
	}

	req.Page = page
	req.Limit = limit

	resp, err = h.grpcClient.JournalService().GetList(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to get list of journals")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /UpdateJournal/{id} [put]
// @Summary       Update a journal by ID
// @Description   API for updating a journal by ID
// @Tags          journal
// @Accept        json
// @Produce       json
// @Param         id path string true "Journal ID"
// @Param         journal body schedule_service.UpdateJournal true "Journal"
// @Success       200 {object} schedule_service.GetJournal
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) UpdateJournal(c *gin.Context) {
	var (
		id   = c.Param("id")
		req  schedule_service.UpdateJournal
		resp *schedule_service.GetJournal
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "invalid request body")
		return
	}

	req.Id = id
	resp, err = h.grpcClient.JournalService().Update(c.Request.Context(), &req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to update journal")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router        /DeleteJournal/{id} [delete]
// @Summary       Delete a journal by ID
// @Description   API for deleting a journal by ID
// @Tags          journal
// @Accept        json
// @Produce       json
// @Param         id path string true "Journal ID"
// @Success       200 {object} schedule_service.EmptyJournal
// @Failure       404 {object} models.ResponseError
// @Failure       500 {object} models.ResponseError
func (h *handler) DeleteJournal(c *gin.Context) {
	var (
		id   = c.Param("id")
		resp = &schedule_service.EmptyJournal{}
		err  error
	)

	data, err := getAuthInfo(c)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "error while getting auth")
		return
	}

	if data.UserRole != "SuperAdmin" && data.UserRole != "Manager" && data.UserRole != "Administration" {
		handleGrpcErrWithDescription(c, h.log, errors.New("Unauthorized"), "You are not a SuperAdmin or a manager or an admin")
		return
	}

	req := &schedule_service.JournalPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.JournalService().Delete(c.Request.Context(), req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "failed to delete journal")
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router       /GetJurnalsStudent/{id} [get]
// @Summary      Get a Jurnal by Student Group ID
// @Description  Get a Jurnal entry by Student Group ID
// @Produce      json
// @Tags         journal
// @Param        id path string true "Student Group ID"
// @Success      200 {object} schedule_service.GetJournal
// @Failure      400 {object} models.ResponseError "Invalid request body"
// @Failure      500 {object} models.ResponseError "Internal server error"
func (h *handler) GetJurnalByIDStudent(c *gin.Context) {
	id := c.Param("id")
	req := &schedule_service.JournalPrimaryKey{Id: id}

	resp, err := h.grpcClient.JournalService().GetByIDStudent(c, req)
	if err != nil {
		handleGrpcErrWithDescription(c, h.log, err, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)
}