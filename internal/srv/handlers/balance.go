package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RsBalance struct {
	Curr      float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (api *API) GetBalance(c *gin.Context) {
	userID := c.MustGet(UserIDKey).(int)

	balance, err := api.db.GetBalance(c, userID)
	if err != nil {
		api.logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, RsDef{
			Err: InternalSeverErrMsg,
		})
		return
	}

	withdrawn, err := api.db.GetWithdrawalSum(c, userID)
	if err != nil {
		api.logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, RsDef{
			Err: InternalSeverErrMsg,
		})
		return
	}

	c.JSON(http.StatusOK, RsBalance{
		Curr:      balance,
		Withdrawn: withdrawn,
	})
}
