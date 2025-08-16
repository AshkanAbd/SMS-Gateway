package handlers

import (
	"time"

	smsmodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/sms/models"
	usermodels "github.com/AshkanAbd/arvancloud_sms_gateway/internal/modules/user/models"
)

type stdResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func newMessageResponse(message string) stdResponse {
	return stdResponse{
		Message: message,
		Data:    nil,
	}
}

func newObjectResponse(v any) stdResponse {
	return stdResponse{
		Data:    v,
		Message: "",
	}
}

type createUserRequest struct {
	Name string `json:"name"`
}

func (c createUserRequest) toUser() usermodels.User {
	return usermodels.User{
		Name: c.Name,
	}
}

type userResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func fromUser(user usermodels.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
	}
}

type smsResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Receiver  string    `json:"receiver"`
	Status    string    `json:"status"`
	Cost      int       `json:"cost"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func fromSms(sms smsmodels.Sms) smsResponse {
	return smsResponse{
		ID:        sms.ID,
		Content:   sms.Content,
		Receiver:  sms.Receiver,
		Status:    fromSmsStatus(sms.Status),
		Cost:      sms.Cost,
		CreatedAt: sms.CreatedAt,
		UpdatedAt: sms.UpdatedAt,
	}
}

func fromSmsStatus(status smsmodels.SmsStatus) string {
	switch status {
	case smsmodels.StatusScheduled:
		return "Scheduled"
	case smsmodels.StatusEnqueued:
		return "Enqueued"
	case smsmodels.StatusSent:
		return "Sent"
	case smsmodels.StatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

type smsRequest struct {
	Content  string `json:"content"`
	Receiver string `json:"receiver"`
}

func (r smsRequest) toSms() smsmodels.Sms {
	return smsmodels.Sms{
		Content:  r.Content,
		Receiver: r.Receiver,
	}
}

type increaseBalanceRequest struct {
	Amount int64 `json:"balance"`
}
