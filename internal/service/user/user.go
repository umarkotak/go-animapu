package user

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/umarkotak/go-animapu/internal/models"
	rUser "github.com/umarkotak/go-animapu/internal/repository/user"
)

// RegisterService register user to firebase
func RegisterService(userData models.UserData) (models.UserData, error) {
	existingUser := rUser.GetUserByUsernameFromFirebase(userData.Username)

	if existingUser.Username != "" {
		return userData, errors.New("username is taken")
	}

	userData.LoginToken = loginTokenGenerator(userData)
	newUser := rUser.SetUserToFirebase(userData)

	return newUser, nil
}

// LoginService register user to firebase
func LoginService(userData models.UserData) (models.UserData, error) {
	existingUser := rUser.GetUserByUsernameFromFirebase(userData.Username)

	if existingUser.Password != userData.Password {
		return userData, errors.New("username or password is wrong")
	}

	loginTokenDecoder(existingUser.LoginToken)

	return existingUser, nil
}

// DetailService get user detail from firebase
func DetailService(auth string) (models.UserData, error) {
	username, err := loginTokenDecoder(auth)

	var userData models.UserData
	if err != nil {
		return userData, err
	}

	userData = rUser.GetUserByUsernameFromFirebase(username)

	return userData, nil
}

// RecordLastReadHistory set user data to firebase
func RecordLastReadHistory(userData models.UserData, readHistory models.ReadHistory) (models.UserData, error) {
	rUser.SetMangahubHistory(userData, readHistory)

	return userData, nil
}

func loginTokenGenerator(userData models.UserData) string {
	goAnimapuJwtSecret := os.Getenv("GO_ANIMAPU_JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userData.Username,
	})

	tokenString, err := token.SignedString([]byte(goAnimapuJwtSecret))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func loginTokenDecoder(loginToken string) (string, error) {
	tokenString := loginToken
	goAnimapuJwtSecret := os.Getenv("GO_ANIMAPU_JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(goAnimapuJwtSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	username := token.Claims.(jwt.MapClaims)["username"]
	usernameStr := fmt.Sprintf("%v", username)

	return usernameStr, nil
}

func StoreMangaToMyLibrary(userData models.UserData, myLibrary models.MyLibrary) (string, error) {
	rUser.SetMangaToMyLibrary(userData, myLibrary)

	return "", nil
}

func RemoveMangaFromMyLibrary(userData models.UserData, myLibrary models.MyLibrary) (string, error) {
	rUser.RemoveMangaFromMyLibrary(userData, myLibrary)

	return "", nil
}
