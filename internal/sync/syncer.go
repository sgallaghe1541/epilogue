package sync

import (
	"log/slog"
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
)

type Syncer struct {
	viewpoint *viewpoint.ViewpointConnection
	epilogue  *db.EpilogueConnection
	logger    *slog.Logger
	stopSync  chan bool
}

func New(viewpoint *viewpoint.ViewpointConnection, epilogue *db.EpilogueConnection, logger *slog.Logger) *Syncer {
	return NewWithSyncInterval(viewpoint, epilogue, logger, 4*time.Hour)
}

func NewWithSyncInterval(viewpoint *viewpoint.ViewpointConnection, epilogue *db.EpilogueConnection, logger *slog.Logger, syncInterval time.Duration) *Syncer {
	s := &Syncer{
		viewpoint: viewpoint,
		epilogue:  epilogue,
		logger:    logger,
	}
	if syncInterval > 0 {
		go s.startSync(syncInterval)
	}
	return s
}

func (s *Syncer) startSync(syncInterval time.Duration) {
	s.stopSync = make(chan bool)
	ticker := time.NewTicker(syncInterval)

	for {
		select {
		case <-ticker.C:
			errCode := s.sync()
			if errCode != 0 {
				s.logger.Error("sync failed")
				s.StopSync()
			}
		case <-s.stopSync:
			ticker.Stop()
			return
		}
	}
}

func (s *Syncer) StopSync() {
	if s.stopSync != nil {
		s.stopSync <- true
	}
}

func (s *Syncer) Sync() {
	err := s.sync()
	if err > 0 {
		s.logger.Error("sync failed")
	}
}
