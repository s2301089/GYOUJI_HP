package service

import (
	"strconv"

	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
)

// SettingService は設定に関するビジネスロジックのインターフェースです。
type SettingService interface {
	GetVisibility() (bool, error)
	UpdateVisibility(value bool) error
}

// settingService は SettingService の実装です。
type settingService struct {
	repo repository.SettingRepository
}

// NewSettingService は settingService の新しいインスタンスを生成します。
func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

// GetVisibility はスコア表示設定を取得します。
func (s *settingService) GetVisibility() (bool, error) {
	value, err := s.repo.GetSetting("showTotalScores")
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

// UpdateVisibility はスコア表示設定を更新します。
func (s *settingService) UpdateVisibility(value bool) error {
	return s.repo.UpdateSetting("showTotalScores", strconv.FormatBool(value))
}
