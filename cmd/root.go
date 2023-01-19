/*
Copyright © 2023 Heroku PE Developer Experience

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"errors"
	"github.com/heroku/get-app-token/pkg/services"
	"github.com/heroku/get-app-token/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "get-app-token app",
	Short:             "Generate a GitHub app token",
	Long:              `Generate a GitHub app token and print the value`,
	PersistentPreRunE: persistentPreRun,
	RunE:              runLocal,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().StringP(
		"app-id",
		"a",
		"",
		"GitHub app ID. Or set 'APP_ID' environment variable",
	)
	markPersistentRequired("app-id")
	setPersistentWithEnvVar("app-id", "APP_ID")

	rootCmd.PersistentFlags().StringP(
		"private-key",
		"p",
		"",
		"GitHub app private key (mutually exclusive with 'private-key-file'). "+
			"Or set 'PRIVATE_KEY' environment variable",
	)
	setPersistentWithEnvVar("private-key", "PRIVATE_KEY")

	rootCmd.PersistentFlags().StringP(
		"private-key-file",
		"f",
		"",
		"Path to GitHub app private key file (mutually exclusive with 'private-key'). "+
			"Or set 'PRIVATE_KEY_FILE' environment variable",
	)
	setPersistentWithEnvVar("private-key-file", "PRIVATE_KEY_FILE")

	rootCmd.MarkFlagsMutuallyExclusive("private-key", "private-key-file")

	rootCmd.PersistentFlags().StringP(
		"repository",
		"r",
		"",
		"GitHub repository where app is installed. Or set 'GITHUB_REPOSITORY' environment variable",
	)
	markPersistentRequired("repository")
	setPersistentWithEnvVar("repository", "GITHUB_REPOSITORY")
}

func markPersistentRequired(flagName string) {
	if err := rootCmd.MarkPersistentFlagRequired(flagName); err != nil {
		panic(err)
	}
}

func setPersistentWithEnvVar(flagName string, envKey string) {
	if value, ok := os.LookupEnv(envKey); ok {
		if err := rootCmd.PersistentFlags().Set(flagName, value); err != nil {
			panic(err)
		}
	}
}

func persistentPreRun(cmd *cobra.Command, _ []string) error {
	privateKey, err := cmd.Root().PersistentFlags().GetString("private-key")

	if err != nil {
		return err
	}

	privateKeyFile, err := cmd.Root().PersistentFlags().GetString("private-key-file")

	if err != nil {
		return err
	}

	if utils.IsWhitespaceOrEmpty(&privateKey) && utils.IsWhitespaceOrEmpty(&privateKeyFile) {
		if requiredErr := cmd.Root().ValidateRequiredFlags(); requiredErr != nil {
			errMsg := strings.TrimRight(requiredErr.Error(), "not set")

			return errors.New(errMsg + `, "private-key" (or "private-key-file") not set`)
		}

		return errors.New(`required flag(s) "private-key" or "private-key-file" not set`)
	}

	return nil
}

func runLocal(cmd *cobra.Command, _ []string) (err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	var appToken string

	if appToken, err = getAppToken(cmd.Root()); err != nil {
		return err
	}

	println(appToken)

	return err
}

//region Global helper methods

func getAppTokenSvc(cmd *cobra.Command) (services.IAppTokenService, error) {
	appId, _ := cmd.PersistentFlags().GetString("app-id")
	privateKey, _ := cmd.PersistentFlags().GetString("private-key")
	privateKeyFile, _ := cmd.PersistentFlags().GetString("private-key-file")
	repository, _ := cmd.PersistentFlags().GetString("repository")

	if utils.IsWhitespaceOrEmpty(&privateKey) && !utils.IsWhitespaceOrEmpty(&privateKeyFile) {
		var pkBytes []byte
		var err error

		if pkBytes, err = os.ReadFile(privateKeyFile); err != nil {
			return nil, err
		}

		privateKey = string(pkBytes)
	}

	ghApiOpsProvider := services.NewGitHubApiOperationsProvider(appId, privateKey, repository)
	appTokenSvc := services.NewAppTokenService(ghApiOpsProvider)

	return appTokenSvc, nil
}

func getAppToken(cmd *cobra.Command) (appToken string, err error) {
	var appTokenSvc services.IAppTokenService

	if appTokenSvc, err = getAppTokenSvc(cmd); err != nil {
		return "", err
	}

	if appToken, err = appTokenSvc.GetAppToken(); err != nil {
		return "", err
	}

	return appToken, err
}

//endregion
