package urlshort

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type pathUrl struct{
	gorm.Model
	Path string `yaml:"path" json:"path"`
	URL string `yaml:"url" json:"url"`
}

func (user *pathUrl) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now())
	scope.SetColumn("updated_at", time.Now())
	scope.SetColumn("id", uuid.New())
	return nil
}