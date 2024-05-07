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
	Actions    []ChangedData
}

func NewCDCStore(pool *sync.Pool) *CDCStore {
	return &CDCStore{
		eventsPool: pool,
		Actions:    make([]ChangedData, 0),
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

	s.Actions = append(s.Actions, ad)

	return &ad, nil
}

// CreateEventsWithFilter filter db events
func (s *CDCStore) CreateEventsWithFilter(ctx context.Context, tableMap map[string][]string) <-chan *publisher.Event {
	output := make(chan *publisher.Event)

	go func(ctx context.Context) {
		for _, item := range s.Actions {
			if err := ctx.Err(); err != nil {
				slog.Debug("create events with filter: context canceled")
				break
			}

			event := s.eventsPool.Get().(*publisher.Event)
			event.ID = item.ID
			event.Collection = item.Coll
			event.Data = item.NewDocument
			event.Action = item.Kind.string()

			actions, validTable := tableMap[fmt.Sprintf("%s.%s", item.Db, item.Coll)]
			validAction := inArray(actions, item.Kind.string())
			if validTable && validAction {
				output <- event
				continue
			}

			// TODO: add prom metric counter

			slog.Debug(
				"cdc-message was skipped by filter",
				slog.String("collection", item.Coll),
				slog.String("action", string(item.Kind)),
			)
		}
		close(output)
	}(ctx)

	return output
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
