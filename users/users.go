package users

import (
	"encoding/json"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func LoadUsers() ([]User, error) {
	file, err := os.Open("./static/users.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, err
}

func LoadUsersByRole() (map[string][]User, error) {
	users, err := LoadUsers()
	if err != nil {
		return nil, err
	}
	segmentedUsers := make(map[string][]User)
	for _, user := range users {
		segmentedUsers[user.Role] = append(segmentedUsers[user.Role], user)
	}
	return segmentedUsers, nil
}
