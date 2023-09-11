package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jenniekibiri/jwt-go/initializers"
	"github.com/jenniekibiri/jwt-go/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//Get the email and password from the request body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are empty"})
		return
	}
	// check if the user already exists
	var user models.User
	results := initializers.DB.First(&user, "email = ?", body.Email)

	if results.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// hash the paswword
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// create a new user
	newUser := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	//Respond with the user
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})

}

func Login(c *gin.Context) {
	//get the email and password from the request body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields are empty"})
		return
	}
	// look up the the requested user
	var user models.User
	results := initializers.DB.First(&user, "email = ?", body.Email)

	if results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return

	}

	//compare the password with the hash password in the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}
	//if the password is correct, generate a token and send it to the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})



	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "localhost", false, true)


	//Respond with the token
	c.JSON(http.StatusOK, gin.H{})

}


func Validate(c *gin.Context){
	user, _ := c.Get("user")


	c.JSON(http.StatusOK, gin.H{"message": user})

}