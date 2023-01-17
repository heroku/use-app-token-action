package services

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v48/github"
	actions "github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"time"
)

//go:generate mockery --name IAppTokenService --structname AppTokenService --output ../_mocks/services
type IAppTokenService interface {
	GetAppToken() (*string, error)
}

//go:generate mockery --name IGitHubApiOperationsProvider --structname GitHubApiOperationsProvider --output ../_mocks/services
type IGitHubApiOperationsProvider interface {
	FindRepositoryInstallation() (*github.Installation, error)
	CreateInstallationToken(
		installationId int64,
		tokenOptions *github.InstallationTokenOptions,
	) (*github.InstallationToken, error)
}

type GitHubOperationsProviderOptions struct {
	appId      string
	privateKey string
	repository string
}

type GitHubApiOperationsProvider struct {
	IGitHubApiOperationsProvider
	context              context.Context
	client               *github.Client
	fullRepoName         string
	owner                string
	repo                 string
	appId                string
	privateKey           string
	installation         *github.Installation
	appToken             *github.InstallationToken
	bearerTokenExpiresAt *time.Time
}

func NewGitHubApiOperationsProvider(options *GitHubOperationsProviderOptions) IGitHubApiOperationsProvider {
	var appId, privateKey, repository string

	if options == nil {
		appId, _ = os.LookupEnv("APP_ID")
		privateKey, _ = os.LookupEnv("PRIVATE_KEY")
		repository, _ = os.LookupEnv("GITHUB_REPOSITORY")
	} else {
		appId = options.appId
		privateKey = options.privateKey
		repository = options.repository
	}

	if IsWhitespaceOrEmpty(&appId) || IsWhitespaceOrEmpty(&privateKey) || IsWhitespaceOrEmpty(&repository) {
		panic(errors.New(
			"appId, privateKey, and repository are required. " +
				"Supply `options` for these values, or " +
				"set 'APP_ID', 'PRIVATE_KEY', and 'GITHUB_REPOSITORY' env vars",
		))
	}

	actions.AddMask(privateKey)

	repoNameSplit := strings.Split(repository, "/")
	owner := repoNameSplit[0]
	repo := repoNameSplit[1]

	return &GitHubApiOperationsProvider{
		context:      context.Background(),
		fullRepoName: repository,
		owner:        owner,
		repo:         repo,
		appId:        appId,
		privateKey:   privateKey,
	}
}

func (p *GitHubApiOperationsProvider) FindRepositoryInstallation() (*github.Installation, error) {
	if p.installation != nil {
		return p.installation, nil
	}

	p.refreshClient()

	installation, _, err := p.client.Apps.FindRepositoryInstallation(p.context, p.owner, p.repo)
	p.installation = installation

	return p.installation, err
}

func (p *GitHubApiOperationsProvider) CreateInstallationToken(
	installationId int64,
	tokenOptions *github.InstallationTokenOptions,
) (*github.InstallationToken, error) {
	appTokenExpired := p.appToken == nil || time.Now().Unix()+60 > p.appToken.ExpiresAt.Unix()

	p.refreshClient()

	if appTokenExpired {
		appToken, _, err := p.client.Apps.CreateInstallationToken(p.context, installationId, tokenOptions)
		p.appToken = appToken

		return p.appToken, err
	}

	return p.appToken, nil
}

func (p *GitHubApiOperationsProvider) refreshClient() {
	if p.bearerTokenExpiresAt != nil && time.Now().Unix()+30 <= p.bearerTokenExpiresAt.Unix() {
		return
	}

	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(2 * time.Minute)
	p.bearerTokenExpiresAt = &exp
	bearerToken, err := p.generateJwtToken(iss, exp)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	oauth2Token := &oauth2.Token{AccessToken: bearerToken}
	ts := oauth2.StaticTokenSource(oauth2Token)
	tc := oauth2.NewClient(ctx, ts)
	p.client = github.NewClient(tc)
}

func (p *GitHubApiOperationsProvider) generateJwtToken(issuedAt time.Time, expiresAt time.Time) (string, error) {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(p.privateKey))

	if err != nil {
		return "", err
	}

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    p.appId,
	}
	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	bearerToken, err := bearer.SignedString(rsaPrivateKey)

	if err != nil {
		return "", err
	}

	actions.AddMask(bearerToken)

	return bearerToken, nil
}

type AppTokenService struct {
	IAppTokenService
	ghApiOpsProvider IGitHubApiOperationsProvider
}

func NewAppTokenService(ghApiOpsProvider IGitHubApiOperationsProvider) IAppTokenService {
	return &AppTokenService{ghApiOpsProvider: ghApiOpsProvider}
}

func (s *AppTokenService) GetAppToken() (*string, error) {
	installation, err := s.ghApiOpsProvider.FindRepositoryInstallation()

	if err != nil {
		return nil, err
	}

	appToken, err := s.ghApiOpsProvider.CreateInstallationToken(*installation.ID, &github.InstallationTokenOptions{})

	if err != nil {
		return nil, err
	}

	actions.AddMask(*appToken.Token)

	return appToken.Token, nil
}

func IsWhitespaceOrEmpty(value *string) bool {
	if value == nil {
		return true
	}

	return strings.TrimSpace(*value) == ""
}
