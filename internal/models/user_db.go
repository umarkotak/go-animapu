package models

// UserDB user database representative
type UserDB struct {
	UserDatas map[string]*UserData `json:"user_db"`
}
