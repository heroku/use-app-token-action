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

//go:generate mockery --name IAppTokenService --structname AppTokenService --output ../mocks/services
type IAppTokenService interface {
	GetAppToken() (*string, error)
}

//go:generate mockery --name IGitHubApiOperationsProvider --structname GitHubApiOperationsProvider --output ../mocks/services
type IGitHubApiOperationsProvider interface {
	FindRepositoryInstallation() (*github.Installation, *github.Response, error)
	CreateInstallationToken(
		installationId int64,
		tokenOptions *github.InstallationTokenOptions,
	) (*github.InstallationToken, *github.Response, error)
}

type GitHubApiOperationsProvider struct {
	IGitHubApiOperationsProvider
	context      context.Context
	client       *github.Client
	fullRepoName string
	owner        string
	repo         string
}

func NewGitHubApiOperationsProvider() IGitHubApiOperationsProvider {
	appId, appIdOk := os.LookupEnv("APP_ID")
	privateKey, privateKeyOk := os.LookupEnv("PRIVATE_KEY")
	repository, repositoryOk := os.LookupEnv("GITHUB_REPOSITORY")

	if !appIdOk || !privateKeyOk || !repositoryOk {
		panic(errors.New("'APP_ID', 'PRIVATE_KEY', and 'GITHUB_REPOSITORY' env vars are required"))
	}

	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(2 * time.Minute)
	bearerToken, err := GenerateJwtToken(appId, privateKey, iss, exp)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	oauth2Token := &oauth2.Token{AccessToken: bearerToken}
	ts := oauth2.StaticTokenSource(oauth2Token)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repoNameSplit := strings.Split(repository, "/")
	owner := repoNameSplit[0]
	repo := repoNameSplit[1]

	return &GitHubApiOperationsProvider{
		context:      context.Background(),
		client:       client,
		fullRepoName: repository,
		owner:        owner,
		repo:         repo,
	}
}

func (p *GitHubApiOperationsProvider) FindRepositoryInstallation() (*github.Installation, *github.Response, error) {
	return p.client.Apps.FindRepositoryInstallation(p.context, p.owner, p.repo)
}

func (p *GitHubApiOperationsProvider) CreateInstallationToken(
	installationId int64,
	tokenOptions *github.InstallationTokenOptions,
) (*github.InstallationToken, *github.Response, error) {
	return p.client.Apps.CreateInstallationToken(p.context, installationId, tokenOptions)
}

type AppTokenService struct {
	IAppTokenService
	ghApiOpsProvider IGitHubApiOperationsProvider
}

func NewAppTokenService(ghApiOpsProvider IGitHubApiOperationsProvider) IAppTokenService {
	return &AppTokenService{ghApiOpsProvider: ghApiOpsProvider}
}

func (s *AppTokenService) GetAppToken() (*string, error) {
	installation, _, err := s.ghApiOpsProvider.FindRepositoryInstallation()

	if err != nil {
		return nil, err
	}

	appToken, _, err := s.ghApiOpsProvider.CreateInstallationToken(*installation.ID, &github.InstallationTokenOptions{})

	if err != nil {
		return nil, err
	}

	actions.AddMask(*appToken.Token)

	return appToken.Token, nil
}

func GenerateJwtToken(appId string, privateKey string, issuedAt time.Time, expiresAt time.Time) (string, error) {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))

	if err != nil {
		return "", err
	}

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    appId,
	}
	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	bearerToken, err := bearer.SignedString(rsaPrivateKey)

	if err != nil {
		return "", err
	}

	actions.AddMask(bearerToken)

	return bearerToken, nil
}
