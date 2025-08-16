package handlers

import (
	"errors"
	"github.com/AshkanAbd/arvancloud_sms_gateway/common"
	"net/http"

	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/smsgateway"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	usermodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
)

type HttpHandler struct {
	gateway  *smsgateway.SmsGateway
	validate *validator.Validate
}

func buildResponse(c *fiber.Ctx, status int, resp stdResponse) error {
	return c.Status(status).JSON(resp)
}

func paginateFromQuery(c *fiber.Ctx) (int, int) {
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("pageSize", "10")
	page := common.ParseIntWithFallback(pageStr, 1)
	pageSize := common.ParseIntWithFallback(pageSizeStr, 10)

	return (page - 1) * pageSize, pageSize
}

func (h *HttpHandler) getValidationErrors(v any) validator.ValidationErrors {
	validationRes := h.validate.Struct(v)
	if validationRes != nil {
		var errs validator.ValidationErrors
		if errors.As(validationRes, &errs) {
			return errs
		}
	}

	return nil
}

func NewHttpHandler(gateway *smsgateway.SmsGateway) *HttpHandler {
	validate := validator.New()

	return &HttpHandler{
		gateway:  gateway,
		validate: validate,
	}
}

func (h *HttpHandler) CreateUser(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
	}
	validationErrs := h.getValidationErrors(req)
	if len(validationErrs) > 0 {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(validationErrs.Error()))
	}

	createdUser, err := h.gateway.CreateUser(c.Context(), req.toUser())
	if err != nil {
		if errors.Is(err, usermodels.EmptyNameError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}

		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	return buildResponse(c, http.StatusOK, newObjectResponse(fromUser(createdUser)))
}

func (h *HttpHandler) GetUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse("Invalid user id"))
	}

	user, err := h.gateway.GetUser(c.Context(), userId)
	if err != nil {
		if errors.Is(err, usermodels.UserNotExistError) {
			return buildResponse(c, http.StatusNotFound, newMessageResponse(err.Error()))
		}

		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	return buildResponse(c, http.StatusOK, newObjectResponse(fromUser(user)))
}

func (h *HttpHandler) GetUserMessages(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse("Invalid user id"))
	}
	skip, limit := paginateFromQuery(c)
	order := c.Query("order", "desc")
	if order != "desc" && order != "asc" {
		order = "desc"
	}

	messages, err := h.gateway.GetUserMessages(c.Context(), userId, skip, limit, order == "desc")
	if err != nil {
		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	se := make([]smsResponse, len(messages))
	for i := range messages {
		se[i] = fromSms(messages[i])
	}

	return buildResponse(c, http.StatusOK, newObjectResponse(se))
}

func (h *HttpHandler) SendSingleMessage(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse("Invalid user id"))
	}

	var req smsRequest
	if err := c.BodyParser(&req); err != nil {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
	}
	validationErrs := h.getValidationErrors(req)
	if len(validationErrs) > 0 {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(validationErrs.Error()))
	}

	err := h.gateway.SendSingleMessage(c.Context(), userId, req.toSms())
	if err != nil {
		if errors.Is(err, smsmodels.EmptyContentError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}
		if errors.Is(err, smsmodels.EmptyReceiverError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}

		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	return buildResponse(c, http.StatusOK, newMessageResponse("message scheduled successfully"))
}

func (h *HttpHandler) SendBulkMessage(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse("Invalid user id"))
	}

	var req []smsRequest
	if err := c.BodyParser(&req); err != nil {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
	}
	validationErrs := h.getValidationErrors(req)
	if len(validationErrs) > 0 {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(validationErrs.Error()))
	}

	ss := make([]smsmodels.Sms, len(req))
	for i := range req {
		ss[i] = req[i].toSms()
	}
	err := h.gateway.SendBulkMessage(c.Context(), userId, ss)
	if err != nil {
		if errors.Is(err, smsmodels.EmptyContentError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}
		if errors.Is(err, smsmodels.EmptyReceiverError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}

		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	return buildResponse(c, http.StatusOK, newMessageResponse("messages scheduled successfully"))
}

func (h *HttpHandler) IncreaseUserBalance(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse("Invalid user id"))
	}

	var req increaseBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
	}
	validationErrs := h.getValidationErrors(req)
	if len(validationErrs) > 0 {
		return buildResponse(c, http.StatusBadRequest, newMessageResponse(validationErrs.Error()))
	}

	newBalance, err := h.gateway.IncreaseUserBalance(c.Context(), userId, req.Amount)
	if err != nil {
		if errors.Is(err, usermodels.InvalidBalanceError) {
			return buildResponse(c, http.StatusBadRequest, newMessageResponse(err.Error()))
		}
		if errors.Is(err, usermodels.UserNotExistError) {
			return buildResponse(c, http.StatusNotFound, newMessageResponse(err.Error()))
		}

		return buildResponse(c, http.StatusInternalServerError, newMessageResponse(err.Error()))
	}

	return buildResponse(c, http.StatusOK, newObjectResponse(fiber.Map{
		"balance": newBalance,
	}))
}
