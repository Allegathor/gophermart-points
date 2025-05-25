package handlers

import (
	"errors"
	"gophermart-points/internal/entity"
	"gophermart-points/internal/repo/pgsql"
	"gophermart-points/pkg/checksum"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (api *OrderAPI) RegOrder(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}
	num := string(body)

	fieldErrs := orderNumValid.Validate(c, &num)
	if fieldErrs != nil {
		c.JSON(http.StatusUnprocessableEntity, RsOrder{
			Err:       InvalidInputs,
			FieldErrs: fieldErrs,
		})
		return
	}

	if !checksum.VerifyLuhn(num) {
		c.JSON(http.StatusUnprocessableEntity, RsOrder{
			Err: InvalidInputs,
			FieldErrs: map[string][]string{
				"ordernum": {"Incorrect checksum"},
			},
		})
		return
	}

	userID := c.MustGet(UserIDKey).(int)
	order := entity.NewOrder(userID, num, 0)
	exOrder, err := api.db.GetOrder(c, order.Num)
	if err != nil && !errors.Is(err, pgsql.ErrNoOrder) {
		api.logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, RsDef{Err: err.Error()})
		return
	}

	if order.Num == exOrder.Num {
		if order.UserID == exOrder.UserID {
			c.JSON(http.StatusOK, RsDef{Err: ""})
			return
		}

		c.JSON(http.StatusConflict, RsDef{Err: "Has been arleady uploaded"})
		return
	}

	newID, err := api.db.AddOrder(c, *order)
	if err != nil {
		api.logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}

	order.OrderID = newID
	api.Queue.Push(*order)

	c.JSON(http.StatusAccepted, RsDef{Err: ""})
}

func (api *OrderAPI) Orders(c *gin.Context) {
	userID := c.MustGet(UserIDKey).(int)

	orders, err := api.db.GetOrders(c, userID)
	if err != nil {
		api.logger.Errorln(err)
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}

	sl := make([]OrderRec, 0)
	for _, o := range orders {
		var rec OrderRec
		rec.Num = o.Num
		rec.Amount = o.Amount
		rec.Status = o.PntsEvalStatus
		rec.UploadAt = o.UploadAt.Format(time.RFC3339)
		sl = append(sl, rec)
	}

	c.JSON(http.StatusOK, sl)
}
