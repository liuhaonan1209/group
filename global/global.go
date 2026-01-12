package global

import (
	"group/config"

	"gorm.io/gorm"
)

var (
	AppConfig config.Config
	DB        *gorm.DB
)
