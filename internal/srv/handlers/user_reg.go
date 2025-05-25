package handlers

import (
	"context"
	"gophermart-points/internal/datacrypt"
	"gophermart-points/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *UserAPI) UserReg(c *gin.Context) {
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

	hasUser, err := api.db.HasUser(c, u.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsUser{
			Err: InternalSeverErrMsg,
		})
		return
	}

	if hasUser {
		c.JSON(http.StatusConflict, RsUser{
			Err: InvalidInputs,
			FieldErrs: map[string][]string{
				"login": {"This login is taken. Please, try another"},
			},
		})
		return
	}

	hashedPwd, err := datacrypt.HashPassword(u.Pwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsUser{
			Err: InternalSeverErrMsg,
		})
		return
	}

	id, err := api.db.CreateUser(c, u.Login, hashedPwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsUser{
			Err: InternalSeverErrMsg,
		})
		return
	}

	tn, err := datacrypt.BuildUserJWT(id, api.authKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RsUser{
			Err: InternalSeverErrMsg,
		})
		return
	}

	c.SetCookie(USER_COOKIE_NAME, tn, int(datacrypt.TOKEN_EXP), "api/user/", "localhost", false, true)
	c.JSON(http.StatusOK, RsUser{Err: "", FieldErrs: nil})
}
