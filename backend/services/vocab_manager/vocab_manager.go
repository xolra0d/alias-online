package main

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"
)

type Vocabulary struct {
	PrimaryWords []string
	RudeWords    []string
}

type VocabManager struct {
	vocabs map[string]Vocabulary
	lock   sync.RWMutex

	db     *Postgres
	logger *slog.Logger

	done       chan struct{}
	doneCtx    context.Context
	doneCancel context.CancelFunc

	LoadVocabsTimeout        time.Duration
	PollInterval             time.Duration
	ClosePostgresConnTimeout time.Duration
}

func NewVocabManager(db *Postgres, logger *slog.Logger, loadVocabsTimeout, pollInterval, closePostgresConnTimeout time.Duration) (*VocabManager, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), loadVocabsTimeout)
	vocabs, ok := db.LoadVocabs(ctx)
	cancel()
	if !ok {
		return nil, false
	}

	ctx, cancel = context.WithCancel(context.Background())

	return &VocabManager{
		vocabs: vocabs,
		db:     db,
		logger: logger,

		done:       make(chan struct{}),
		doneCtx:    ctx,
		doneCancel: cancel,

		LoadVocabsTimeout:        loadVocabsTimeout,
		PollInterval:             pollInterval,
		ClosePostgresConnTimeout: closePostgresConnTimeout,
	}, true
}

func (v *VocabManager) StartObservation() {
	v.logger.Info("starting observation loop")

	conn, err := v.db.db.Acquire(v.doneCtx)
	if err != nil {
		if v.doneCtx.Err() == nil {
			v.logger.Error("could not acquire database connection", "err", err)
		}
		v.done <- struct{}{}
		return
	}
	_, err = conn.Exec(v.doneCtx, "LISTEN vocab_updates")
	if err != nil {
		if v.doneCtx.Err() == nil {
			v.logger.Error("could not listen vocab_updates", "err", err)
		}
		conn.Release()
		v.done <- struct{}{}
		return
	}

	for {
		n, err := conn.Conn().WaitForNotification(v.doneCtx)
		if err != nil {
			if v.doneCtx.Err() == nil {
				v.logger.Error("could not wait for notification", "err", err)
			}
			conn.Release()
			v.done <- struct{}{}
			return
		}

		ctx, cancel := context.WithTimeout(v.doneCtx, v.LoadVocabsTimeout)
		vocabs, ok := v.db.LoadVocabs(ctx)
		cancel()
		if ok {
			v.logger.Info("updating vocab_manager", "n", n)
			v.lock.Lock()
			v.vocabs = vocabs
			v.lock.Unlock()
		}

		timer := time.NewTimer(v.PollInterval)

		select {
		case <-v.doneCtx.Done():
			timer.Stop()
			conn.Release()
			v.done <- struct{}{}
			return
		case <-timer.C:
		}
	}
}

func (v *VocabManager) StopObservation() {
	v.logger.Info("stopping observation loop")
	v.doneCancel()
	<-v.done
	v.logger.Info("stopped observation loop")
}

func (v *VocabManager) AvailableVocabs() []string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return slices.Collect(maps.Keys(v.vocabs))
}

func (v *VocabManager) Vocab(name string) (primaryWords, RudeWords []string) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	vocab := v.vocabs[name]

	primary := make([]string, len(vocab.PrimaryWords))
	copy(primary, vocab.PrimaryWords)

	rude := make([]string, len(vocab.RudeWords))
	copy(rude, vocab.RudeWords)

	return primary, rude
}
