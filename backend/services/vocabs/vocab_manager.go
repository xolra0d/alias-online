package main

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"sync/atomic"
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

	runObservation atomic.Bool
	stop           chan struct{}
	done           chan struct{}

	LoadVocabsTimeout time.Duration
	PollInterval      time.Duration
}

func NewVocabManager(db *Postgres, logger *slog.Logger, loadVocabsTimeout, pollInterval time.Duration) (*VocabManager, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), loadVocabsTimeout)
	vocabs, ok := db.LoadVocabs(ctx)
	cancel()
	if !ok {
		return nil, false
	}

	return &VocabManager{
		vocabs:            vocabs,
		db:                db,
		logger:            logger,
		stop:              make(chan struct{}),
		done:              make(chan struct{}),
		LoadVocabsTimeout: loadVocabsTimeout,
		PollInterval:      pollInterval,
	}, true
}

func (v *VocabManager) StartObservation() {
	v.logger.Info("starting observation loop")
	v.runObservation.Store(true)

out:
	for {
		if !v.runObservation.Load() {
			break out
		}
		ctx, cancel := context.WithTimeout(context.Background(), v.LoadVocabsTimeout)
		vocabs, ok := v.db.LoadVocabs(ctx)
		cancel()

		if ok {
			v.lock.RLock()
			eq := mapsEqual(vocabs, v.vocabs)
			v.lock.RUnlock()

			if !eq {
				v.logger.Info("updating vocabs")
				v.lock.Lock()
				v.vocabs = vocabs
				v.lock.Unlock()
			}
		}

		select {
		case <-v.stop:
			break out
		case <-time.After(v.PollInterval):
		}
	}

	v.done <- struct{}{}
}

func mapsEqual(a, b map[string]Vocabulary) bool {
	if len(a) != len(b) {
		return false
	}

	for key, aVal := range a {
		bVal, ok := b[key]
		if !ok {
			return false
		}
		if !slices.Equal(aVal.PrimaryWords, bVal.PrimaryWords) ||
			!slices.Equal(aVal.RudeWords, bVal.RudeWords) {
			return false
		}
	}

	return true
}

func (v *VocabManager) StopObservation() {
	if v.runObservation.CompareAndSwap(true, false) {
		v.logger.Info("stopping observation loop")
		v.stop <- struct{}{}
		<-v.done
	}
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
