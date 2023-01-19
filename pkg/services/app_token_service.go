package services

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

//go:generate mockery --name IAppTokenService --structname AppTokenService --output ../_mocks/services
type IAppTokenService interface {
	GetAppToken() (string, error)
}

//go:generate mockery --name IGitHubApiOperationsProvider --structname GitHubApiOperationsProvider --output ../_mocks/services
type IGitHubApiOperationsProvider interface {
	FindRepositoryInstallation() (*github.Installation, error)
	CreateInstallationToken(
		installationId int64,
		tokenOptions *github.InstallationTokenOptions,
	) (*github.InstallationToken, error)
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

func NewGitHubApiOperationsProvider(appId string, privateKey string, repository string) IGitHubApiOperationsProvider {
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

	if err := p.refreshClient(); err != nil {
		return nil, err
	}

	installation, _, err := p.client.Apps.FindRepositoryInstallation(p.context, p.owner, p.repo)
	p.installation = installation

	return p.installation, err
}

func (p *GitHubApiOperationsProvider) CreateInstallationToken(
	installationId int64,
	tokenOptions *github.InstallationTokenOptions,
) (*github.InstallationToken, error) {
	appTokenExpired := p.appToken == nil || time.Now().Unix()+60 > p.appToken.ExpiresAt.Unix()

	if err := p.refreshClient(); err != nil {
		return nil, err
	}

	if appTokenExpired {
		appToken, _, err := p.client.Apps.CreateInstallationToken(p.context, installationId, tokenOptions)
		p.appToken = appToken

		return p.appToken, err
	}

	return p.appToken, nil
}

func (p *GitHubApiOperationsProvider) refreshClient() error {
	if p.bearerTokenExpiresAt != nil && time.Now().Unix()+30 <= p.bearerTokenExpiresAt.Unix() {
		return nil
	}

	iss := time.Now().Add(-30 * time.Second).Truncate(time.Second)
	exp := iss.Add(2 * time.Minute)
	p.bearerTokenExpiresAt = &exp
	bearerToken, err := p.generateJwtToken(iss, exp)

	if err != nil {
		return err
	}

	ctx := context.Background()
	oauth2Token := &oauth2.Token{AccessToken: bearerToken}
	ts := oauth2.StaticTokenSource(oauth2Token)
	tc := oauth2.NewClient(ctx, ts)
	p.client = github.NewClient(tc)

	return nil
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

	return bearerToken, nil
}

type AppTokenService struct {
	IAppTokenService
	ghApiOpsProvider IGitHubApiOperationsProvider
}

func NewAppTokenService(ghApiOpsProvider IGitHubApiOperationsProvider) IAppTokenService {
	return &AppTokenService{ghApiOpsProvider: ghApiOpsProvider}
}

func (s *AppTokenService) GetAppToken() (string, error) {
	installation, err := s.ghApiOpsProvider.FindRepositoryInstallation()

	if err != nil {
		return "", err
	}

	appToken, err := s.ghApiOpsProvider.CreateInstallationToken(*installation.ID, &github.InstallationTokenOptions{})

	if err != nil {
		return "", err
	}

	return *appToken.Token, nil
}
