package seeders

import (
	"log"
	"web-app/app/models/user"
	"web-app/app/services"
)

type UserSeeder struct{}

func (u *UserSeeder) Run() {
	userModel := user.NewUserModel()

	userModel.Username = "John-Doe"
	hashedPassword, err := services.HashPassword("password")
	if err != nil {
		log.Fatalf("error hashing password: %v", err)
	}
	userModel.Password = hashedPassword

	err = userModel.Create()
	if err != nil {
		log.Fatalf("error creating user: %v", err)
	}
}
