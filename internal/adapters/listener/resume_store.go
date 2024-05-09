package listener

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WildEgor/cdc-listener/internal/models"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"time"
)

// TODO: move to config
const (
	serviceFolderPath   = "/data/"
	resumeStateFilename = "resume_state.json"
)

var _ ITokenSaver = (*ResumeTokenSaver)(nil)

// TODO: make using adapters: switch store between fs and mongodb
// ResumeTokenSaver
type ResumeTokenSaver struct {
	in   chan *models.ResumeTokenState
	stop chan struct{}
}

func NewResumeStore() *ResumeTokenSaver {
	return &ResumeTokenSaver{
		in:   make(chan *models.ResumeTokenState, 100),
		stop: make(chan struct{}),
	}
}

func (rs *ResumeTokenSaver) SaveResumeToken(data *models.ResumeTokenState) error {
	rs.in <- data

	return nil
}

func (rs *ResumeTokenSaver) GetResumeToken(db, coll string) string {
	stateFilename := fmt.Sprintf("%s_%s_%s", db, coll, resumeStateFilename)

	currDir, _ := os.Getwd()
	dir := path.Join(currDir, serviceFolderPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return ""
		}
	}

	f, err := os.Open(path.Join(dir, stateFilename))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ""
		} else {
			return ""
		}
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return ""
	}

	state := &models.ResumeTokenState{}

	err = json.Unmarshal(data, &state)
	if err != nil {
		return ""
	}

	return state.LastMongoResumeToken
}

func (rs *ResumeTokenSaver) Run() {
	for {
		select {
		case state := <-rs.in:
			slog.Debug("receive resume token", slog.String("value", state.LastMongoResumeToken))

			state.LastMongoProcessedTime = time.Now()

			bytes, err := json.Marshal(state)
			if err != nil {
				slog.Error("failed to marshal robust message", slog.Any("err", err.Error()))
				continue
			}

			stateFilename := fmt.Sprintf("%s_%s_%s", state.Db, state.Coll, resumeStateFilename)

			currDir, _ := os.Getwd()
			f, err := os.Create(path.Join(currDir, serviceFolderPath, stateFilename))
			if err != nil {
				slog.Error("failed to create robust state file", slog.Any("err", err.Error()))
				continue
			}

			_, err = f.Write(bytes)
			if err != nil {
				slog.Error("failed to write robust state file", slog.Any("err", err.Error()))
				continue
			}

			defer f.Close()
		case <-rs.stop:
			slog.Debug("resume token saver closed")
			break
		}
	}
}

func (rs *ResumeTokenSaver) Stop() {
	close(rs.stop)
}
