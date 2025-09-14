package service

import (
	"github.com/saku0512/GYOUJI_HP/backend/internal/model"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

type RelayService interface {
	RegisterRelayResults(relayType string, classIDs []int) error
	GetRelayResults(relayType string) ([]model.RelayResult, error)
}

type RelayServiceImpl struct {
	Repo *repository.RelayRepository
}

func NewRelayService(repo *repository.RelayRepository) RelayService {
	return &RelayServiceImpl{Repo: repo}
}

func (s *RelayServiceImpl) RegisterRelayResults(relayType string, classIDs []int) error {
	return s.Repo.RegisterRelayResults(relayType, classIDs)
}

func (s *RelayServiceImpl) GetRelayResults(relayType string) ([]model.RelayResult, error) {
	return s.Repo.GetRelayResults(relayType)
}
