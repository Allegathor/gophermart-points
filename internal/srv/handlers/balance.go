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
	userId := c.MustGet(USER_ID_KEY).(int)

	balance, err := api.db.GetBalance(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{
			Err: err.Error(),
		})
		return
	}

	withdrawn, err := api.db.GetWithdrawalSum(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{
			Err: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RsBalance{
		Curr:      balance,
		Withdrawn: withdrawn,
	})
}
