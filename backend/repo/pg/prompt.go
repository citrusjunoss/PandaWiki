package pg

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/pg"
	"gorm.io/gorm"
)

type PromptRepo struct {
	db     *pg.DB
	logger *log.Logger
}

type promptJson struct {
	Content string
}

func NewPromptRepo(db *pg.DB, logger *log.Logger) *PromptRepo {
	return &PromptRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PromptRepo) GetPrompt(ctx context.Context, kbID string) (string, error) {
	var setting domain.Setting
	var prompt promptJson
	err := r.db.WithContext(ctx).Table("settings").
		Where("kb_id = ? AND key = ?", kbID, domain.SettingKeySystemPrompt).
		First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if err := json.Unmarshal(setting.Value, &prompt); err != nil {
		return "", err
	}
	return prompt.Content, nil
}
