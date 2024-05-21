package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"storage"
	"users"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Register user
// @Description  Register a user
// @Tags         registration
// @Param        note body users.User true "user JSON"
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /register [post]
func Register(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var user users.User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		if user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty password"})
			return
		}
		if user.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty usename"})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		user.Password = string(hashedPassword)
		err = dbc.AddUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't add a user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successful registration"})
	}
	return gin.HandlerFunc(fn)
}

// @Summary      Authorisation
// @Description  Authorise a user
// @Tags         authorisation
// @Param        note body users.User true "user JSON"
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /authorise [post]
func Authorise(dbc *storage.DBConnected) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		secretKey := []byte(os.Getenv("AUTH_KEY"))
		var user users.User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}
		existingUser, err := dbc.GetUser(user.Username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No such user"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{"id": existingUser.ID, "username": existingUser.Username},
		)
		jwtToken, err := token.SignedString(secretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successful authorisation", "token": jwtToken})
	}
	return gin.HandlerFunc(fn)
}

func AuthMiddleware(secretKey []byte) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorisation")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorisation header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorisation header"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}
		c.Next()
	}
	return gin.HandlerFunc(fn)
}

func ExtractToken(c *gin.Context) (string, error) {
	token := c.Query("token")
	if token != "" {
		return token, nil
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", fmt.Errorf("No token")
}

func ExtractTokenUsername(c *gin.Context) (string, error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return "", err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_KEY")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["username"].(string), nil
	}
	return "", fmt.Errorf("Error occured along the way")
}
