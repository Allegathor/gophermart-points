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

	userId := c.MustGet(USER_ID_KEY).(int)
	order := entity.NewOrder(userId, num, 0)
	exOrder, err := api.db.GetOrder(c, order.Num)
	if err != nil && !errors.Is(err, pgsql.ErrNoOrder) {
		c.JSON(http.StatusInternalServerError, RsDef{Err: err.Error()})
		return
	}

	if order.Num == exOrder.Num {
		if order.UserId == exOrder.UserId {
			c.JSON(http.StatusOK, RsDef{Err: ""})
			return
		}

		c.JSON(http.StatusConflict, RsDef{Err: "Has been arleady uploaded"})
		return
	}

	newId, err := api.db.AddOrder(c, *order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}

	order.OrderId = newId
	api.Queue.Push(*order)

	c.JSON(http.StatusAccepted, RsDef{Err: ""})
}

func (api *OrderAPI) Orders(c *gin.Context) {
	userId := c.MustGet(USER_ID_KEY).(int)
	orders, err := api.db.GetOrders(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsDef{Err: InternalSeverErrMsg})
		return
	}

	var sl []OrderRec
	for _, o := range orders {
		rec := OrderRec{}
		rec.Num = o.Num
		rec.Amount = o.Amount
		rec.Status = o.PntsEvalStatus
		rec.UploadAt = o.UploadAt.Format(time.RFC3339)
		sl = append(sl, rec)
	}

	c.JSON(http.StatusOK, sl)
}
