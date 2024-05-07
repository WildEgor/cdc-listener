package listener

import (
	"context"
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/WildEgor/cdc-listener/internal/configs"
	"github.com/WildEgor/cdc-listener/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sync"
	"time"
)

var _ IListener = (*Listener)(nil)

// Listener wrapper for collections listen logic
type Listener struct {
	eventPubAdapter *publisher.EventPublisherAdapter
	repo            ICDCRepository
	store           *CDCStore
	cfg             *configs.ListenerConfig
}

// NewListener create new listener
func NewListener(pub *publisher.EventPublisherAdapter, repo ICDCRepository, cfg *configs.ListenerConfig) *Listener {
	pool := &sync.Pool{
		New: func() any {
			return &publisher.Event{}
		},
	}

	return &Listener{
		eventPubAdapter: pub,
		repo:            repo,
		store:           NewCDCStore(pool),
		cfg:             cfg,
	}
}

// WatchCollection get updates from collections
func (l *Listener) WatchCollection(ctx context.Context, opts *WatchCollectionOptions) error {
	resume := true
	for resume {
		time.Sleep(time.Second)

		token, err := l.repo.GetResumeToken(opts.ResumeTokensCollCapped)
		if err != nil {
			return err
		}

		changeStreamOpts := options.ChangeStream().
			SetFullDocument(options.UpdateLookup).
			SetFullDocumentBeforeChange(options.WhenAvailable)

		if len(token) != 0 {
			m := bson.M{"_data": token}
			data, err := bson.Marshal(m)
			if err != nil {
				slog.Error("failed to marshal bson resume token: ", err.Error())
				return errors.ErrFailMarshalResumeToken
			}

			changeStreamOpts.SetResumeAfter(data)
		}

		cs, err := l.repo.GetWatchStream(changeStreamOpts)
		if err != nil {
			return err
		}

		for cs.Next(ctx) {
			currentResumeToken := cs.Current.Lookup("_id", "_data").StringValue()
			operationType := cs.Current.Lookup("operationType").StringValue()

			json, err := bson.MarshalExtJSON(cs.Current, false, false)
			if err != nil {
				return errors.ErrFailMarshalStreamData
			}

			if _, ok := publishableOperationTypes[operationType]; !ok {
				if operationType == invalidateOperationType {
					resume = false
					slog.Warn("invalid operation type")
					break
				}
				continue
			}

			if err = opts.ChangeEventHandler(ctx, opts.WatchedSubj, currentResumeToken, OperationType(operationType), json); err != nil {
				slog.Error("could not publish change event", "err", err)
				break
			}

			// FIXME
			//if err = l.repo.SaveResumeToken(currentResumeToken); err != nil {
			//	slog.Error("could not insert resume token", "err", err)
			//	break
			//}
		}

		slog.Info("stopped watching mongodb collection")

		if err = cs.Close(context.Background()); err != nil {
			return errors.ErrFailCloseStream
		}
	}

	return nil
}

// Run start listen collections from db
func (l *Listener) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for subj, _ := range l.cfg.MappedFilter {
		group.Go(func() error {
			watchCollOpts := &WatchCollectionOptions{
				WatchedSubj:            subj,
				ResumeTokensCollCapped: false,
				StreamName:             "", // TODO: use as topic prefix
				ChangeEventHandler: func(ctx context.Context, subj, msgId string, operationType OperationType, data []byte) error {
					_, err := l.store.AssertData(msgId, subj, operationType, []byte{}, data)
					if err != nil {
						return err
					}

					event := l.store.CreateEventsWithFilter(groupCtx, l.cfg.MappedFilter)
					if event == nil {
						return nil
					}
					slog.Debug("publish")

					if err := l.eventPubAdapter.Publisher.Publish(groupCtx, "notifier", event); err != nil {
						return err
					}

					return nil
				},
			}

			return l.WatchCollection(groupCtx, watchCollOpts)
		})
	}

	return group.Wait()
}

// Stop stop listener and disconnect db
func (l *Listener) Stop() error {
	return l.repo.Close()
}
