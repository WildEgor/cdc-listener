package listener

import (
	"context"
	"fmt"
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"log/slog"
	"strings"
	"sync"
)

// ChangedData is kind of CDC message data.
type ChangedData struct {
	// ID document id
	ID string
	// Db name
	Db string
	// Coll collection
	Coll string
	// Kind operation type
	Kind OperationType
	// TODO: need somehow get it
	OldDocument []byte
	// NewDocument updated data
	NewDocument []byte
}

// CDCStore collects changes
type CDCStore struct {
	eventsPool *sync.Pool
	data       *ChangedData
}

func NewCDCStore(pool *sync.Pool) *CDCStore {
	return &CDCStore{
		eventsPool: pool,
	}
}

// AssertData add changes to store
func (s *CDCStore) AssertData(_id string, subj string, kind OperationType, oldDocument []byte, newDocument []byte) (a *ChangedData, err error) {
	subjects := strings.Split(subj, ".")

	ad := ChangedData{
		ID:          _id,
		Db:          subjects[0],
		Coll:        subjects[1],
		Kind:        kind,
		OldDocument: []byte{},
		NewDocument: newDocument,
	}

	s.data = &ad

	return s.data, nil
}

// CreateEventsWithFilter filter db events
func (s *CDCStore) CreateEventsWithFilter(ctx context.Context, tableMap map[string][]string) *publisher.Event {
	if s.data == nil {
		return nil
	}

	if err := ctx.Err(); err != nil {
		slog.Debug("create events with filter: context canceled")
		return nil
	}

	event := s.eventsPool.Get().(*publisher.Event)
	event.ID = s.data.ID
	event.Collection = s.data.Coll
	event.Data = s.data.NewDocument
	event.Action = s.data.Kind.string()

	s.eventsPool.Put(event)

	actions, validTable := tableMap[fmt.Sprintf("%s.%s", s.data.Db, s.data.Coll)]
	validAction := inArray(actions, s.data.Kind.string())
	if validTable && validAction {
		return event
	}

	// TODO: add prom metric counter

	slog.Debug(
		"cdc-message was skipped by filter",
		slog.String("collection", s.data.Coll),
		slog.String("action", string(s.data.Kind)),
	)

	return nil
}

// inArray checks whether the value is in an array.
func inArray(arr []string, value string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}
