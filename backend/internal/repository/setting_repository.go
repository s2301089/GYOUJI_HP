package repository

import (
	"database/sql"
)

// SettingRepository は設定データへのアクセスに関するインターフェースです。
type SettingRepository interface {
	GetSetting(key string) (string, error)
	UpdateSetting(key string, value string) error
}

// settingRepository は SettingRepository の実装です。
type settingRepository struct {
	db *sql.DB
}

// NewSettingRepository は settingRepository の新しいインスタンスを生成します。
func NewSettingRepository(db *sql.DB) SettingRepository {
	return &settingRepository{db: db}
}

// GetSetting はキーで設定値を取得します。
func (r *settingRepository) GetSetting(key string) (string, error) {
	query := "SELECT setting_value FROM scoretable_settings WHERE setting_key = ?"
	row := r.db.QueryRow(query, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// UpdateSetting はキーで設定値を更新します。
func (r *settingRepository) UpdateSetting(key string, value string) error {
	query := "INSERT INTO scoretable_settings (setting_key, setting_value) VALUES (?, ?) ON DUPLICATE KEY UPDATE setting_value = ?"
	_, err := r.db.Exec(query, key, value, value)
	return err
}
