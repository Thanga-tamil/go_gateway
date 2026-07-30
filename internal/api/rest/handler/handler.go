package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Thanga-tamil/noway_service/internal/dto"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/Thanga-tamil/noway_service/internal/response"
	"github.com/Thanga-tamil/noway_service/internal/service"
	"github.com/Thanga-tamil/noway_service/internal/utils"
)

func HandleUserRegister(c *gin.Context) {
	tenantId := c.Request.Header.Get("tenant-x")
	tenantDB, exists := c.Get(tenantId)
	if !exists {
		logger.Errorf("tenant-x '%s' not found in context", tenantId)
		c.JSON(http.StatusBadRequest, response.Error("invalid tenant", utils.STATUS_BAD_REQUEST))
		return
	}

	db, ok := tenantDB.(*sqlx.DB)
	if !ok {
		logger.Errorf("tenant DB for '%s' has unexpected type %T", tenantId, tenantDB)
		c.JSON(http.StatusInternalServerError, response.Error("internal server error", http.StatusInternalServerError))
		return
	}

	user, err := parseInputFromReq(c)
	if err != nil {
		logger.Infof("invalid register input: %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof("parsed register user input: %#v", user)

	userId := uuid.New().String()

	service.RegisterService(db, userId, user)

	jwt, err := service.ServeJwt(userId)
	if err != nil {
		logger.Errorf("failed to generate JWT for user %s: %s", userId, err.Error())
		c.JSON(http.StatusInternalServerError, response.Error("failed to generate token", http.StatusInternalServerError))
		return
	}

	if err := service.StoreJwtInRedis(userId, jwt); err != nil {
		logger.Errorf("failed to store JWT in redis for user %s: %s", userId, err.Error())
		c.JSON(http.StatusInternalServerError, response.Error("failed to store token", http.StatusInternalServerError))
		return
	}

	res := response.Success("Onboarding completed successfully",
		utils.STATUS_OK,
		map[string]any{"accessToken": jwt})

	c.JSON(http.StatusOK, res)
}

func parseInputFromReq(c *gin.Context) (dto.UserRegisterReqPayload, error) {

	var user dto.UserRegisterReqPayload
	if err := c.ShouldBindBodyWith(&user, binding.JSON); err != nil {

		logger.Info("Error while Decode input payload: ", err)
		resp := map[string]any{
			"status": 400,
			"message": "Request body must not be null"}
			val, _ := json.Marshal(resp)
			return dto.UserRegisterReqPayload{}, errors.New(string(val))
	}

	if err := service.ValidateInput(user); err != nil {
		logger.Error("input payload validation error: ", err)

		msg := map[string]any{"status": 400, "message": err.Error()}
		val, _ := json.Marshal(msg)
		return dto.UserRegisterReqPayload{}, errors.New(string(val))
	}

	return user, nil
}


func GenerateJwtToken(c *gin.Context) {
	userId := c.Request.Header.Get("userId")
	if len(userId) == 0 {
		logger.Error("GenerateJwtToken called with empty userId header")
		c.JSON(http.StatusBadRequest, response.Error("userId cannot be empty", utils.STATUS_BAD_REQUEST))
		return
	}

	jwt, err := service.ServeJwt(userId)
	if err != nil {
		logger.Errorf("failed to generate JWT for userId '%s': %s", userId, err.Error())
		c.JSON(http.StatusInternalServerError, response.Error("failed to generate token", http.StatusInternalServerError))
		return
	}

	if err := service.StoreJwtInRedis(userId, jwt); err != nil {
		logger.Errorf("failed to store JWT in redis for userId '%s': %s", userId, err.Error())
		c.JSON(http.StatusInternalServerError, response.Error("failed to store token", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.Success("token generated successfully", utils.STATUS_OK, map[string]any{"jwtToken": jwt}))
}

func handleErr(c *gin.Context, err error) {

	resp := response.Error(err.Error(), utils.STATUS_BAD_REQUEST)

	c.JSON(http.StatusBadRequest, resp)
}
