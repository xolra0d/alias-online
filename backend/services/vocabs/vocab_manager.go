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
}

func NewVocabManager(db *Postgres, logger *slog.Logger) *VocabManager {
	return &VocabManager{
		db:     db,
		logger: logger,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func (v *VocabManager) StartObservation(loadVocabsTimeout, pollInterval time.Duration) {
	v.logger.Info("starting observation loop")
	v.runObservation.Store(true)

	for {
		if !v.runObservation.Load() {
			break
		}
		ctx, cancel := context.WithTimeout(context.Background(), loadVocabsTimeout)
		vocabs, ok := v.db.LoadVocabs(ctx)
		cancel()

		if ok {
			v.lock.RLock()
			eq := mapsEqual(vocabs, v.vocabs)
			v.lock.RUnlock()

			if !eq {
				v.lock.Lock()
				v.vocabs = vocabs
				v.lock.Unlock()
			}
		}

		select {
		case <-v.stop:
			break
		case <-time.After(pollInterval):
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
