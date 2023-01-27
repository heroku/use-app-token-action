[![Use GitHub App Token](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml/badge.svg)](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml)

# Use GitHub App Token Action

This was created to generate a GitHub application token for use in workflows and applications by supplying the app ID,
the app's RSA private key, and name of the repository where the app is installed. It is intended to be used in 3
potential ways:

1. As a GitHub action (primary use case)
2. As a command line application
3. As a library within other Go apps

## Usage

### As a GitHub Action

```yaml
- uses: heroku/use-app-token-action@main
  with:
    # GitHub App ID
    # required: true
    app_id: ""
    # GitHub App private key
    # required: true
    private_key: ""
    # GitHub repository where the GH App is installed and authorized
    # Use if the generated token is being used to target another repository, i.e. target is NOT the current repo.
    # For example, if the token will be used to check out code from another repository, then the app must have
    # repository contents permissions to that target repository. Otherwise, if not specified (default), it's assumed
    # that the app is authorized to perform the required actions for the current repository.
    # https://docs.github.com/en/rest/overview/permissions-required-for-github-apps?apiVersion=2022-11-28
    # required: false
    # default: ${{ github.repository }}
    repository: ""
```

Returns: `steps.<step_id>.outputs.app_token`

In your workflow YAML file, include this action similar to the following: \
  
```yaml
job:
  name: My Job
  runs_on: sfdc_hk_ubuntu_latest
  steps:
   - name: Generate access token
     id: generate_access_token
     uses: heroku/use-app-token-action@main
     with:
        app_id: ${{ secrets.GH_APP_ID }}
        private_key: ${{ secrets.GH_APP_PRIVATE_KEY }}
        repository: heroku/some-other-repository
   - name: Task that needs a token
     uses: actions/checkout@v3
     with:
       repo: heroku/some-other-repository
       path: ./other-repo
       token: ${{ steps.generate_access_token.outputs.app_token }}
```

### As a command line application

1. Install the binary to your machine:
   ```bash
   go install github.com/heroku/use-app-token-action/cmd/get-app-token@v0.0.1
   ```
2. The application will be installed in `${GOPATH}/bin/get-app-token`
3. Run the command (assuming your `${GOPATH}\bin` folder is on your `${PATH}`):
   ```bash
   get-app-token --app-id <APP_ID> --private-key-file <PATH_TO_RSA_PRIVATE_KEY> --repository <GH_REPO_WITH INSTALLED_APP>
   ```

### AS a library withing other Go apps

1. Add the dependency to your application
   ```bash
    go get github.com/heroku/use-app-token-action/pkg/...@v0.0.1
   ```
2. Use it in your application. For example:
   ```go
   package main
   
   import (
       "fmt"
       "os"
   
       "github.com/heroku/use-app-token-action/pkg/services"
   )
   
   func main() {
       appId := "123456"
       privateKeyFile := "/path/to/the/gh_app_rsa_private_key.pem"
       privateKeyBytes, _ := os.ReadFile(privateKeyFile)
       privateKey := string(privateKeyBytes)
       repository := "heroku/my-go-application"
   
       appTokenSvc := services.NewAppTokenService(services.NewGitHubApiOperationsProvider(
           appId,
           privateKey,
           repository,
       ))
   
       token, _ := appTokenSvc.GetAppToken()
   
       fmt.Println(token)
   }
   ```
3. NOTE: If you're using this in a long-running application, the token will expire after 10 minutes. However, for a
   single instance of the `AppTokenService`, subsequent calls to `GetAppToken()` will return the same token as long as
   it hasn't expired, and returns a new token if it has expired (auto-refresh).

## Development

Modifications to this project require that the version number, `VERSION`, in the [Makefile](./Makefile) is updated to
reflect the scope of the change being applied. Follow semantic versioning best practice.

For release and feature branches, a tag will be generated with a `release` or `beta` suffix followed by the UTC date of
the build. This is done performed automatically as a part of the CI workflow upon push to the remote, e.g.
`v0.0.1-beta-20230127.185110`. The generated tag can be used for debugging. Remove release/beta tags when they're no 
longer needed.

When a PR is merged to the `main` branch, binaries are regenerated and tagged as a part of the CI workflow. The tag will
be the `VERSION` as specified in the [Makefile](./Makefile) (without a suffix, e.g. `v0.0.1`).

NOTE: Perform a `git pull` after pushing changes to the remote to keep your development branch up to date with the
remote branch. This is required as the CI process for all branches (except `main`) uses the current datetime, which will
cause the process to commit changes and tag the development branch accordingly.

[![Use GitHub App Token](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml/badge.svg)](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml)
