package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/itouri/fortnite/web/domain"
	playerHttp "github.com/itouri/fortnite/web/player/delivery/http"
)

func TestFetch(t *testing.T) {
	var mockPlayer domain.Player
	err := faker.FakeData(&mockPlayer)
	//FIXME このエラーハンドリングはテストとは assert でいいのか？
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	mockListPlayer := make([]*domain.Player, 0)
	mockListPlayer = append(mockListPlayer, &mockPlayer)

	num := 1
	cursor := "2"
	//LEARN
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListPlayer, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/player?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := playerHttp.HttpPlayerHandler{
		PlayerUsecase: mockUCase,
	}
	handler.FetchPlayer(c)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}


