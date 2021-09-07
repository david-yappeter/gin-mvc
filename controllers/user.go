package controllers

import (
	"errors"
	"myapp/entity"
	"myapp/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userControllerInterface interface {
	CheckToken(*gin.Context)
	Register(*gin.Context)
	Login(*gin.Context)
	LoginPage(*gin.Context)
	RegisterPage(*gin.Context)
}

type userController struct{}

var (
	UserController userControllerInterface
)

func init() {
	UserController = new(userController)
}

// @summary Register a New User
// @description Register a New User, return access_token and refresh_token in cookie
// @Tags Auth
// @Accept json
// @Param UserData body entity.UserRegister true "Registration"
// @Success 200
// @Router /api/users/register [post]
func (controller *userController) Register(c *gin.Context) {
	var input entity.UserRegister
	err := c.ShouldBind(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errH(err))
		return
	}

	accToken, rememberToken, err := service.UserRegister(c.Request.Context(), entity.User{
		ID:               "",
		Name:             input.Name,
		Email:            input.Email,
		Password:         input.HashedPass(),
		Address:          input.Address,
		PhoneCountryCode: input.PhoneCountryCode,
		Phone:            input.Phone,
		Balance:          0,
		RememberToken:    new(string),
	})
	if err == service.ErrRecordFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, errH(errors.New("email exist")))
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errH(err))
		return
	}

	// c.JSON(http.StatusOK, gin.H{})
	setToken(c, accToken, rememberToken)
	c.Status(http.StatusOK)
}

// @summary Login User
// @description  return access_token and refresh_token in cookie
// @Tags Auth
// @Accept json
// @Param UserData body entity.UserLogin true "Registration"
// @Success 200
// @Failure 400 {object} dto.ErrorResp{}
// @Router /api/users/login [post]
func (controller *userController) Login(c *gin.Context) {
	var input entity.UserLogin
	err := c.ShouldBind(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errH(err))
		return
	}

	accessToken, rememberToken, err := service.UserLogin(c.Request.Context(), input)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.AbortWithStatusJSON(http.StatusBadRequest, errH(errors.New("wrong password")))
	} else if err == service.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, errH(errors.New("email not exist")))
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errH(err))
		return
	}
	setToken(c, accessToken, rememberToken)
	c.Status(http.StatusOK)
}

func setToken(c *gin.Context, accessToken string, rememberToken string) {
	c.SetCookie("access_token", accessToken, service.AccExpired, "", "", false, true)
	c.SetCookie("remember_token", rememberToken, service.RemExpired, "", "", false, true)
}

// @summary Token Data
// @description  return data of access_token
// @Tags Auth
// @Accept json
// @Success 200 {object} service.JwtCustomClaim{}
// @Failure 401 {object} dto.ErrorResp{}
// @Router /api [get]
func (controller *userController) CheckToken(c *gin.Context) {
	test := service.CtxVal(c.Request.Context())

	c.JSON(http.StatusOK, test)
}

func (controller *userController) LoginPage(c *gin.Context) {
	data, err := userData(c)
	if err == service.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Login Page",
		"user":  data,
	})
}

func (controller *userController) RegisterPage(c *gin.Context) {
	data, err := userData(c)
	if err == service.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register Page",
		"user":  data,
	})
}
