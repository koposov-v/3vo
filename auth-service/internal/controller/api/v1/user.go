package v1

import (
	"authjwt/internal/controller"
	"authjwt/internal/controller/dto"
	"authjwt/internal/domain"
	"authjwt/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserRoutes struct {
	userUC controller.UserUseCase
}

type UserResponse struct {
	Response
	User *dto.UserResponse `json:"user"`
}

func NewUserRoutes(uc usecase.UserUseCase) *UserRoutes {
	return &UserRoutes{
		userUC: uc,
	}
}

func (r *UserRoutes) Register(e *echo.Group) {
	e.POST("/register", r.CreateUser)
}

func (r *UserRoutes) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return ErrorHTTP(c, http.StatusBadRequest, err)
	}

	user := &domain.User{
		ID:        uuid.NewString(),
		Username:  req.Username,
		CreatedAt: time.Now(),
	}

	if err := r.userUC.CreateUser(ctx, user, req.Password); err != nil {
		return ErrorHTTP(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, UserResponse{
		User: &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		Response: NewResponse("User created"),
	})
}
