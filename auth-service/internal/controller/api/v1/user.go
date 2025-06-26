package v1

import (
	"authjwt/internal/controller"
	"authjwt/internal/controller/dto"
	"authjwt/internal/domain"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// @title Auth Service API
// @version 1.0
// @description This is the API for the Auth Service
// @host localhost:8080
// @BasePath /api/v1

type UserRoutes struct {
	userUC controller.UserUseCase
	v10    *validator.Validate
}

type UserResponse struct {
	Response
	User *dto.UserResponse `json:"user"`
}

func NewUserRoutes(uc controller.UserUseCase, v10 *validator.Validate) *UserRoutes {
	return &UserRoutes{
		userUC: uc,
		v10:    v10,
	}
}

func (r *UserRoutes) Register(e *echo.Group) {
	e.POST("/register", r.CreateUser)
	e.POST("/login", r.Login)
	e.GET("/validate", r.ValidateToken)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /register [post]
func (r *UserRoutes) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return ErrorHTTP(c, http.StatusBadRequest, err)
	}

	if err := validator.New().Struct(req); err != nil {
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

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /login [post]
func (r *UserRoutes) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return ErrorHTTP(c, http.StatusBadRequest, err)
	}
	if err := r.v10.Struct(req); err != nil {
		return ErrorHTTP(c, http.StatusBadRequest, err)
	}

	user, err := r.userUC.FindByUsername(ctx, req.Username)
	if err != nil || !user.VerifyPassword(req.Password) {
		return ErrorHTTP(c, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
	}

	token, err := r.userUC.GenerateToken(user)
	if err != nil {
		return ErrorHTTP(c, http.StatusInternalServerError, fmt.Errorf("system error"))
	}

	return c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: token})
}

// ValidateToken godoc
// @Summary Validate JWT token
// @Description Validate a provided JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /validate [get]
func (r *UserRoutes) ValidateToken(c echo.Context) error {
	h := c.Request().Header.Get("Authorization")
	valid, err := r.userUC.ValidateToken(h)
	if err != nil || !valid {
		return ErrorHTTP(c, http.StatusUnauthorized, err)
	}

	return SuccessHTTP(c, "User validated")
}
