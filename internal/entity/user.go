package entity

type User struct {
	ID             uint64 `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type UserWithAccessToken struct {
	ID             uint64 `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	AccessToken    string `json:"access_token"`
}
