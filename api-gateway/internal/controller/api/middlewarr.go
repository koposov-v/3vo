package api

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

// AuthMiddleware Middleware для проверки токена
func (g *Gateway) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth header")
			}

			g.logger.Info(g.authServiceURL + "/api/v1/validate")
			req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, g.authServiceURL+"/api/v1/validate", nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			req.Header.Set("Authorization", authHeader)

			resp, err := g.client.Do(req)
			if err != nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "auth service unavailable")
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				return echo.NewHTTPError(http.StatusUnauthorized, string(bodyBytes))
			}

			return next(c)
		}
	}
}
