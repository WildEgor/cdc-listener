package listener

import (
	"context"
	"fmt"
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/WildEgor/cdc-listener/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"sync"
	"time"
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
	OldDocument bson.M
	// NewDocument updated data
	NewDocument bson.M
}

// CDCData collects changes
type CDCData struct {
	eventsPool *sync.Pool
	data       *ChangedData
}

func NewCDCData(pool *sync.Pool) *CDCData {
	return &CDCData{
		eventsPool: pool,
	}
}

// Assert add changes to store
func (s *CDCData) Assert(rawEvent *ChangeEventRaw) (a *ChangedData, err error) {

	convNewDocument := s.convertBsonToJson(rawEvent.FullDocument)
	convOldDocument := convNewDocument
	if rawEvent.FullDocumentBeforeChanges != nil {
		convOldDocument = s.convertBsonToJson(rawEvent.FullDocumentBeforeChanges)
	}

	ad := ChangedData{
		ID:          rawEvent.ID,
		Db:          rawEvent.Db,
		Coll:        rawEvent.Coll,
		Kind:        rawEvent.Kind,
		OldDocument: convOldDocument,
		NewDocument: convNewDocument,
	}

	s.data = &ad

	return s.data, nil
}

// FilterEvent filter db events
func (s *CDCData) FilterEvent(ctx context.Context, tableMap map[string][]string) *publisher.Event {
	if s.data == nil {
		slog.Warn("call Assert before filter first")
		return nil
	}

	if err := ctx.Err(); err != nil {
		slog.Debug("create events with filter: context canceled")
		return nil
	}

	event := s.eventsPool.Get().(*publisher.Event)
	event.ID = s.data.ID
	event.Db = s.data.Db
	event.Collection = s.data.Coll
	event.Data = s.data.NewDocument
	event.Action = s.data.Kind.string()
	event.EventTime = time.Now()

	s.eventsPool.Put(event)

	actions, validTable := tableMap[fmt.Sprintf("%s.%s", s.data.Db, s.data.Coll)]
	validAction := utils.InArray(actions, s.data.Kind.string())
	if validTable && validAction {
		return event
	}

	// TODO: add prom metric counter

	slog.Debug(
		"cdc-message was skipped by filter",
		slog.String("db", s.data.Db),
		slog.String("collection", s.data.Coll),
		slog.String("action", string(s.data.Kind)),
	)

	return nil
}

func (s *CDCData) convertBsonToJson(doc bson.D) map[string]any {
	convertedItem := make(map[string]any)
	for _, val := range doc {

		switch v := val.Value.(type) {
		case primitive.ObjectID:
			convertedItem[val.Key] = v.Hex()
		case bson.M:
			subMap := make(map[string]interface{})
			for subKey, subVal := range v {
				subMap[subKey] = subVal
			}
			convertedItem[val.Key] = subMap
		case bson.D:
			convertedItem[val.Key] = s.convertBsonToJson(v)
		case bson.A:
			subArray := make([]interface{}, len(v))
			for i, subVal := range v {
				subArray[i] = subVal
			}
			convertedItem[val.Key] = subArray
		default:
			convertedItem[val.Key] = val.Value
		}
	}

	return convertedItem
}
