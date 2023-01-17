// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	github "github.com/google/go-github/v48/github"
	mock "github.com/stretchr/testify/mock"
)

// GitHubApiOperationsProvider is an autogenerated mock type for the IGitHubApiOperationsProvider type
type GitHubApiOperationsProvider struct {
	mock.Mock
}

// CreateInstallationToken provides a mock function with given fields: installationId, tokenOptions
func (_m *GitHubApiOperationsProvider) CreateInstallationToken(installationId int64, tokenOptions *github.InstallationTokenOptions) (*github.InstallationToken, error) {
	ret := _m.Called(installationId, tokenOptions)

	var r0 *github.InstallationToken
	if rf, ok := ret.Get(0).(func(int64, *github.InstallationTokenOptions) *github.InstallationToken); ok {
		r0 = rf(installationId, tokenOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.InstallationToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, *github.InstallationTokenOptions) error); ok {
		r1 = rf(installationId, tokenOptions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindRepositoryInstallation provides a mock function with given fields:
func (_m *GitHubApiOperationsProvider) FindRepositoryInstallation() (*github.Installation, error) {
	ret := _m.Called()

	var r0 *github.Installation
	if rf, ok := ret.Get(0).(func() *github.Installation); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.Installation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGitHubApiOperationsProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewGitHubApiOperationsProvider creates a new instance of GitHubApiOperationsProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGitHubApiOperationsProvider(t mockConstructorTestingTNewGitHubApiOperationsProvider) *GitHubApiOperationsProvider {
	mock := &GitHubApiOperationsProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}