[![Use GitHub App Token](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml/badge.svg)](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml)

# Use GitHub App Token Action

This action is intended to be used to create and return a GitHub installation access token given a GitHub Apps `app_id`
and RSA `private_key`.

## Usage

```yaml
- uses: heroku/use-app-token-action@main
  with:
    # GitHub App ID
    # required: true
    app_id: ""
    # GitHub App private key
    # required: true
    private_key: ""
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
   - name: Task that needs a token
     uses: actions/checkout@v3
     with:
       token: ${{ steps.generate_access_token.outputs.app_token }}
```

## Development

Modifications to this project require that the binaries located in the [/bin](bin) directory be compiled and checked in.
To generate the binaries, run `make clean build`, and check in all changes.

[![Use GitHub App Token](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml/badge.svg)](https://github.com/heroku/use-app-token-action/actions/workflows/ci.yaml)
