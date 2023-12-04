package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) loginUserHandler(c echo.Context) error {
	var loginUserRequest struct {
		Username string `json:"username" validate:"required,lte=100"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginUserRequest); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "json input is invalid")
	}

	validate := validator.New()

	err := validate.Struct(loginUserRequest)
	if err != nil {
		errors := strings.Split(err.Error(), "\n")

		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"errors": errors})
	}

	user, err := s.Repository.FindUserByUsername(loginUserRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect username or password")
	}

	ok, err := argon2id.ComparePasswordAndHash(loginUserRequest.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect username or password")
	}

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect username or password")
	}

	isAdmin := false
	if user.Admin == 1 {
		isAdmin = true
	}

	token, err := s.generateNewToken(user.ID, user.Username, isAdmin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, echo.Map{"accessToken": token, "refreshToken": ""})
}

func (s *Server) getLoggedInUser(c echo.Context) error {
	userToken := s.getUserToken(c)

	user, err := s.Repository.FindUserByID(userToken.Id)
	if err != nil {
		s.Logger.Error("could not find user by id", zap.Error(err), zap.Int("userId", userToken.Id))
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	type getLoggedInUserResponse struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Admin     bool   `json:"admin"`
		AvatarURL string `json:"avatarURL"`
	}

	isAdmin := false
	if user.Admin == 1 {
		isAdmin = true
	}

	response := getLoggedInUserResponse{
		ID: user.ID, Username: user.Username, Admin: isAdmin, AvatarURL: "",
	}

	return c.JSON(http.StatusOK, response)
}
