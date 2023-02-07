package main

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/rehab-backend/api/accounts"
	"github.com/rehab-backend/internal/pkg/handlers"
)

func main() {

	client := gocloak.NewClient("http://localhost:8080/auth/")
	ctx := context.Background()

	token, err := client.LoginClient(ctx, "rehab-go", "h5rQCqL7dgofp4OdCBZOFUVBxWJXRNBC", "Rehab")
	if err != nil {
		panic("Login failed:" + err.Error())
	}

	// token, err := client.LoginAdmin(ctx, "admin", "admin", "Rehab")
	// if err != nil {
	// 	panic("Something wrong with the credentials or url")
	// }

	// var params gocloak.GetUsersParams
	// params.Username = "admin"

	// _, err = client.GetUsers(ctx, token.AccessToken, "Rehab", params)
	// if err != nil {
	// 	panic(err.Error())
	// }

	user := gocloak.User{
		FirstName: gocloak.StringP("Bob"),
		LastName:  gocloak.StringP("Uncle"),
		Email:     gocloak.StringP("something@really.wrong"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("CoolGuy"),
	}

	_, err = client.CreateUser(ctx, token.AccessToken, "Rehab", user)
	if err != nil {
		panic(err.Error())
	}

	// dbConnection := config.ConnectDB()

	router := handlers.NewRouter()

	// queries := database.New(postgres.DB)
	// authorService := accounts.NewService()

	// authorService.RegisterHandlers()

	// queries := database.New(postgres.DB)
	// patientService := patients.NewService()
	// patientService.RegisterHandlers(router)

	accountService := accounts.NewService()
	accountService.RegisterHandlers(router)

	handlers.ListenRoute(router)
}
