package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/itouri/fortnite/web/domain"
	"github.com/itouri/fortnite/web/player"
	"github.com/labstack/echo"

	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:message`
}

type HttpPlayerHandler struct {
	PlayerUsecase player.Usecase
}

func NewPlayerHttpHandler(e *echo.Echo, us player.Usecase) {
	handler := &HttpPlayerHandler{
		PlayerUsecase: us,
	}
	//FIXME group化の検討
	e.GET("/players", handler.FetchPlayer)
	e.POST("/players", handler.Store)
	e.GET("/players/:id", handler.GetByID)
	e.DELETE("/players/:id", handler.Delete)
}

func (ph *HttpPlayerHandler) FetchPlayer(c echo.Context) error {

	numStr := c.QueryParam("num")
	num, _ := strconv.Atoi(numStr)
	cursor := c.QueryParam("cursor")
	// LEARN
	ctx := c.Request().Context()
	if ctx == nil {
		// LEARN
		ctx = context.Background()
	}
	listPlayer, nextCursor, err := ph.PlayerUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Curosr`, nextCursor)
	return c.JSON(http.StatusOK, listPlayer)
}

func (ph *HttpPlayerHandler) GetByID(c echo.Context) error {

	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	pl, err := ph.PlayerUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, pl)
}

func isRequestValid(p *domain.Player) (bool, error) {

	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		return false, err
	}
	return true, nil
}

// post って名前にしないのが普通なのだろうか
func (ph *HttpPlayerHandler) Store(c echo.Context) error {
	var pl domain.Player
	err := c.Bind(&pl)
	if err != nil {
		//LEARN なにこのエラー
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&pl); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = ph.PlayerUsecase.Store(ctx, &pl)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, pl)
}

func (ph *HttpPlayerHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = ph.PlayerUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
