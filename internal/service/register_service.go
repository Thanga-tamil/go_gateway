package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Thanga-tamil/noway_service/internal/dto"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/Thanga-tamil/noway_service/internal/utils"
	"github.com/Thanga-tamil/noway_service/internal/config"
	"github.com/Thanga-tamil/noway_service/internal/repository"
)

func RegisterService(db *sqlx.DB, userId string, user dto.UserRegisterReqPayload) {

	logger.Infof("generated uuid: %v for email_id: %s", userId, user.EmailID)

	if err := repository.SaveRegisterUser(db, userId, user.Username, 
				  user.MobileNumber, user.EmailID, false, time.Now()); err != nil {
		logger.Errorf("Error while saving register user: %s", err.Error())
		return
	}

}

func ValidateInput(user dto.UserRegisterReqPayload) error {
	if len(strings.TrimSpace(user.Username)) == 0 { 
		return errors.New(utils.USERNAME_EMPTY)
	} else if len(strings.TrimSpace(user.EmailID)) == 0 { 
		return errors.New(utils.EMAILID_EMPTY)
	} else if len(strings.TrimSpace(user.MobileNumber)) == 0 { 
		return errors.New(utils.MOBILENUMBER_EMPTY)
	} 
	return nil
}

func StoreJwtInRedis(userId, token string) error {
	logger.Infof("store jwt token in redis cache for userId:: %s, with token=> %s", userId, token)
	ctx := context.Background()
	key := "user:" + userId + ":access"
	return config.GoRedis.Set(ctx, key, token, 0).Err()
}
