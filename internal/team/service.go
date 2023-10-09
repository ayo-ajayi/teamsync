package team

import (
	"time"
)

type TeamService struct {
	repo ITeamRepo
}

func NewTeamService(repo ITeamRepo) *TeamService {
	return &TeamService{repo: repo}
}

func (ts *TeamService) CreateTeam(team *Team) error {
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()
	return ts.repo.CreateTeam(team)
}

