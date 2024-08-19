package http

import (
	"net/http"
	"pragusga/internal/usecase"
	ev "pragusga/pkg/env"

	sp "pragusga/pkg/supertokens"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	cfg         *ev.Config
}

func NewAuthHandler(router *gin.Engine, authUseCase *usecase.AuthUseCase, cfg *ev.Config) {
	handler := &AuthHandler{
		authUseCase: authUseCase,
		cfg:         cfg,
	}

	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", handler.SignUp)
		auth.POST("/signin", handler.SignIn)
		auth.POST("/signout", protectedMiddleware(cfg), handler.SignOut)
		auth.GET("/me", protectedMiddleware(cfg), handler.Me)
	}
}

func protectedMiddleware(cfg *ev.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString, err := c.Cookie("access_token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			c.Abort()
			return
		}

		if jwtString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			c.Abort()
			return
		}

		jwks, err := sp.GetJWKS(cfg.SuperTokensConnectionURI)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		parsedToken, parseError := jwt.Parse(jwtString, jwks.Keyfunc)

		if parseError != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": parseError.Error()})
			return
		}

		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			return
		}

		claims := parsedToken.Claims.(jwt.MapClaims)

		// Convert the claims to a key-value pair
		claimsMap := make(map[string]interface{})
		for key, value := range claims {
			claimsMap[key] = value
		}

		c.Set("uid", claimsMap["sub"])

		c.Next()

	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authUseCase.SignUp(c, input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.authUseCase.SignIn(c, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tenantId := "public"
	tokens, err := session.CreateNewSessionWithoutRequestResponse(tenantId, u.ID, nil, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exp, err := tokens.GetExpiry()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", tokens.GetAccessToken(), int(exp), "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Signed in successfully",
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetString("uid")

	user, err := h.authUseCase.GetUserInfo(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
