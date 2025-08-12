package pgsql

import (
	"fmt"
	"strings"
	"time"

	"github.com/AshkanAbd/arvancloud_sms_gateway/common"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/shared"
	"github.com/AshkanAbd/arvancloud_sms_gateway/internal/user/models"
)

type UserEntity struct {
	ID        uint
	Name      string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserEntity) TableName() string {
	return "users"
}

func fromUser(u models.User) UserEntity {
	ue := UserEntity{
		Name:    strings.Trim(u.Name, " "),
		Balance: u.Balance,
	}

	if u.Entity != nil {
		ue.ID = uint(common.AtoiWithFallback(u.ID, 0))
	}
	if u.CreateDate != nil {
		ue.CreatedAt = u.CreatedAt
	}
	if u.UpdateDate != nil {
		ue.UpdatedAt = u.UpdatedAt
	}

	return ue
}

func toUser(ue UserEntity) models.User {
	return models.User{
		Entity: &shared.Entity{
			ID: fmt.Sprintf("%d", ue.ID),
		},
		CreateDate: &shared.CreateDate{
			CreatedAt: ue.CreatedAt,
		},
		UpdateDate: &shared.UpdateDate{
			UpdatedAt: ue.UpdatedAt,
		},
		Name:    ue.Name,
		Balance: ue.Balance,
	}
}
