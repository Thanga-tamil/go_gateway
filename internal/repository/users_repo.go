package repository

import (
	"time"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/jmoiron/sqlx"
)

func SaveRegisterUser(db *sqlx.DB, useId string, username, mobileNumber, 
			emailID string, isDeleted bool, now time.Time) error  {

	stmt := `insert into 
				users (user_id, username, email_id, mobile_number, created_at, is_deleted)
		     values ($1, $2, $3, $4, $5, $6)`

	result, err := db.Exec(stmt, useId, username, emailID, mobileNumber, now, isDeleted)

	if err != nil {
		logger.Error("error while inserting into table user:: ", err.Error())
		return err
	}

	logger.Info("user register db insert result:: ", result)

	return nil
}
