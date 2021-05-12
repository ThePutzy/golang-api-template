package main

import (
	"net/http"
)

type User struct {
	BaseModel
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	_ = srv.Gorm.Find(&users)

	RespondWithJSON(w, http.StatusOK, users)
}
