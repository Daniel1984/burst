package burstcfg

import (
	"github.com/bitburst/burstconsumer/pkg/config"
	"github.com/bitburst/burstconsumer/pkg/db"
	"github.com/bitburst/burstconsumer/pkg/logger"
)

// BurstCfg will hold program wide dependencies for easy DI
type BurstCfg struct {
	*db.DB
	*config.Config
	*logger.Logger
}

// New returns pointer to constructed BurstCfg
func New() (*BurstCfg, error) {
	cfg := config.New()
	log := logger.New()
	db, err := db.New(cfg.GetDBConnStr())
	if err != nil {
		return nil, err
	}

	return &BurstCfg{db, cfg, log}, nil
}
