package MFei

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	LOGGER *zap.SugaredLogger
)
