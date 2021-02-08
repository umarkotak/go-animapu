package models

// UserData user data represent user entity
type UserData struct {
	Username      string                  `json:"username"`
	Password      string                  `json:"password"`
	LoginToken    string                  `json:"login_token"`
	ReadHistories map[string]*ReadHistory `json:"ReadHistories"`
	MyLibraries   map[string]*MyLibrary   `json:"MyLibraries"`
}
