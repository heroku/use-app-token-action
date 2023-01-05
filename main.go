package main

import (
	"github.com/heroku/use-app-token-action/services"
	actions "github.com/sethvargo/go-githubactions"
)

func main() {
	ghApiOpsProvider := services.NewGitHubApiOperationsProvider()
	appTokenSvc := services.NewAppTokenService(ghApiOpsProvider)
	token, err := appTokenSvc.GetAppToken()

	if err != nil {
		actions.Fatalf(err.Error())
	}

	actions.SetOutput("app_token", *token)
	actions.Infof("Token generated successfully: 🔑")
}
