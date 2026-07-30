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
	tenantDB, _ := c.Get(tenantId)

	user, err := parseInputFromReq(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error()); return
	}

	logger.Infof("parsed register user input: %#v\n", user)

	userId := uuid.New().String()

	service.RegisterService(tenantDB.(*sqlx.DB), userId, user)

	jwt, err := service.ServeJwt(userId)

	if err != nil {
		handleErr(c, err); return
	}

	if err := service.StoreJwtInRedis(userId, jwt); err != nil {
		handleErr(c, err); return
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
		resp := map[string]any{"message": "userId can not be null or empty",
		"status": 400}

		val, _ := json.Marshal(resp)
		c.JSON(http.StatusOK, val)
		return
	}


	jwt, err := service.ServeJwt(userId)

	if err != nil {
		logger.Infof("Error on return of ServeJwt: %s", err.Error())
		panic(err)
	}


	if err := service.StoreJwtInRedis(userId, jwt); err != nil {
		resp := map[string]any{"status": 400, "message": err.Error()}
		val, _ := json.Marshal(resp)
		c.JSON(http.StatusOK, val)
		return
	}

	resp := map[string]any{ "status": 200, "jwtToken": jwt}
	val, _ := json.Marshal(resp)
	c.JSON(http.StatusOK, val)
}

func handleErr(c *gin.Context, err error) {

	resp := response.Error(err.Error(), utils.STATUS_BAD_REQUEST)

	c.JSON(http.StatusBadRequest, resp)
}
