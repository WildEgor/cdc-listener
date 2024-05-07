package db

import (
	"github.com/WildEgor/cdc-listener/internal/db/mongodb"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	mongodb.NewMongoConnection,
)
