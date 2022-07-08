package repository

import (
	"website/internal/domain"
	"website/internal/repository/github"
)

type Repo struct {
	github.GithubDao
}

func NewRepo() domain.IRepository {
	return Repo{
		GithubDao: github.NewGithubDao(),
	}
}
