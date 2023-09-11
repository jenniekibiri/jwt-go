package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jenniekibiri/jwt-go/initializers"
	"github.com/jenniekibiri/jwt-go/models"
)

func RequireAuth(c *gin.Context) {
	// get the cookie
	tokenString, err := c.Cookie(("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
	}
	// decode /validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		//check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		}
		
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		}
		c.Set("user", user)




	//find the user with token sub
	//attach to request
	// continue
	c.Next()
	} else {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
	}

	

}
