package sync

import (
	"time"

	"github.com/sgallaghe1541/epilogue/internal/db"
	"github.com/sgallaghe1541/epilogue/internal/viewpoint"
)

type Syncer struct {
	viewpoint *viewpoint.ViewpointConnection
	epilogue  *db.EpilogueConnection
	stopSync  chan bool
}

func New(viewpoint *viewpoint.ViewpointConnection, epilogue *db.EpilogueConnection) *Syncer {
	return NewWithSyncInterval(viewpoint, epilogue, 4*time.Hour)
}

func NewWithSyncInterval(viewpoint *viewpoint.ViewpointConnection, epilogue *db.EpilogueConnection, syncInterval time.Duration) *Syncer {
	s := &Syncer{
		viewpoint: viewpoint,
		epilogue:  epilogue,
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
			//run syn function s.sync()
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
