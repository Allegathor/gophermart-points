package handlers

import (
	"errors"
	"gophermart-points/internal/entity"
	"gophermart-points/internal/repo/pgsql"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RsWithdrawal struct {
	Err       string
	FieldErrs any `json:"fieldErrors"`
}

func (api *API) Withdraw(c *gin.Context) {
	var w entity.Withdrawal
	if err := c.BindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, RsWithdrawal{
			Err:       err.Error(),
			FieldErrs: nil,
		})
		return
	}

	fieldErrs := orderNumValid.Validate(c, &w.Num)
	if fieldErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, RsWithdrawal{
			Err:       InvalidInputs,
			FieldErrs: fieldErrs,
		})
		return
	}

	userId := c.MustGet(USER_ID_KEY).(int)
	w.UserId = userId
	w.Amount = -w.AbsAmount

	err := api.db.Withdraw(c, w)
	if err != nil {
		if errors.Is(err, pgsql.ErrInsufficentFunds) {
			c.JSON(http.StatusPaymentRequired, RsDef{Err: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}

	c.JSON(http.StatusOK, RsDef{Err: ""})
}

type WithdrawalRec struct {
	Num    string  `json:"order"`
	Amount float64 `json:"sum"`
	ProcAt string  `json:"processed_at"`
}

func (api *API) Withdrawals(c *gin.Context) {
	userId := c.MustGet(USER_ID_KEY).(int)
	wls, err := api.db.GetWithdrawals(c, userId)
	var sl []WithdrawalRec
	for _, w := range wls {
		rec := WithdrawalRec{}
		rec.Num = w.Num
		rec.Amount = math.Abs(w.Amount)
		rec.ProcAt = w.ProcAt.Format(time.RFC3339)
		sl = append(sl, rec)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sl)
}
