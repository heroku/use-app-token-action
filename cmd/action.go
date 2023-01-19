/*
Copyright © 2023 NAME HERE heroku-production-services@salesforce.com
*/

package cmd

import (
	actions "github.com/sethvargo/go-githubactions"
	"os"

	"github.com/spf13/cobra"
)

// forGhActionCmd represents the action command
var forGhActionCmd = &cobra.Command{
	Use:   "for-gh-action",
	Short: "Generate a GitHub app token for use in a GitHub action",
	Long: `Generate a GitHub app token for use in a GitHub action and set the GitHub action output,
"app_token", with the generated value`,
	RunE: runAction,
}

func init() {
	rootCmd.AddCommand(forGhActionCmd)
}

func runAction(cmd *cobra.Command, _ []string) (err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	var appToken string

	if appToken, err = getAppToken(cmd.Root()); err != nil {
		return err
	}

	actions.AddMask(appToken)

	if _, ok := os.LookupEnv("GITHUB_OUTPUT"); ok {
		actions.SetOutput("app_token", appToken)
	}

	actions.Infof("Token generated successfully: 🔑")

	return err
}
