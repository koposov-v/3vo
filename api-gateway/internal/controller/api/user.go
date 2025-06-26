package api

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func (g *Gateway) Login(c echo.Context) error {
	g.logger.Info("Заход на маршрут /login")
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodPost, g.authServiceURL+"/api/v1/login", c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req.Header = c.Request().Header.Clone()

	resp, err := g.client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "auth service unavailable")
	}
	defer resp.Body.Close()

	c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get(echo.HeaderContentType))
	c.Response().WriteHeader(resp.StatusCode)

	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
