package errors

import (
	"errors"
)

// Application domain level errors
var (
	ErrFailListCollections    = errors.New("could not list mongo collection names")
	ErrFailCreateCollection   = errors.New("could not create mongo collection")
	ErrFailFindResumeToken    = errors.New("could not fetch or decode resume token")
	ErrFailFindChangeStream   = errors.New("could not watch mongo collection")
	ErrFailCloseStream        = errors.New("could not close change stream")
	ErrFailMarshalStreamData  = errors.New("could not marshal mongo change event from bson")
	ErrFailMarshalResumeToken = errors.New("could not marshal mongo resume token")
	ErrMongoConnection        = errors.New("mongo connection not alive")
	ErrReadResumeFromDisk     = errors.New("failed to process resume token persistence")
)
