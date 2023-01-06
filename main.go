package main

import (
	"github.com/heroku/use-app-token-action/services"
	actions "github.com/sethvargo/go-githubactions"
)

func main() {
	appTokenSvc, err := getAppTokenSvc()

	if err != nil {
		actions.Fatalf(err.Error())
	}

	appToken, err := generateAppToken(appTokenSvc)

	if err != nil {
		actions.Fatalf(err.Error())
	}

	actions.SetOutput("app_token", *appToken)
	actions.Infof("Token generated successfully: 🔑")
}

func getAppTokenSvc() (appTokenSvc services.IAppTokenService, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	ghApiOpsProvider := services.NewGitHubApiOperationsProvider()
	appTokenSvc = services.NewAppTokenService(ghApiOpsProvider)

	return appTokenSvc, err
}

func generateAppToken(appTokenSvc services.IAppTokenService) (appToken *string, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	appToken, err = appTokenSvc.GetAppToken()

	return appToken, err
}
