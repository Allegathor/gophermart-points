package handlers

import (
	"context"
	"errors"
	"gophermart-points/internal/datacrypt"
	"gophermart-points/internal/entity"
	"gophermart-points/internal/repo/pgsql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *UserAPI) UserLogin(c *gin.Context) {
	var u entity.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, RsUser{
			Err:       err.Error(),
			FieldErrs: nil,
		})
		return
	}

	if errors := userValid.Validate(context.TODO(), &u); errors != nil {
		c.JSON(http.StatusBadRequest, RsUser{
			Err:       InvalidInputs,
			FieldErrs: errors,
		})
		return
	}

	dbU, err := api.db.GetUser(c, u.Login)
	if err != nil {
		if errors.Is(err, pgsql.ErrUnexistLogin) {
			c.JSON(http.StatusUnauthorized, RsUser{
				Err: InvalidInputs,
				FieldErrs: map[string][]string{
					"login": {"Wrong login/password"},
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, RsDef{
			Err: InternalSeverErrMsg,
		})
		return
	}

	err = datacrypt.CheckPassword(u.Pwd, dbU.Pwd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, RsUser{
			Err: InvalidInputs,
			FieldErrs: map[string][]string{
				"login": {"Wrong login/password"},
			},
		})
		return
	}

	tn, err := datacrypt.BuildUserJWT(dbU.ID, api.authKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsUser{
			Err: InternalSeverErrMsg,
		})
		return
	}

	c.SetCookie(UserCookieName, tn, int(datacrypt.TokenExp), "api/user/", "localhost", false, true)
	c.JSON(http.StatusOK, RsUser{Err: "", FieldErrs: nil})
}
