package repositories

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/listener"
	"github.com/WildEgor/cdc-listener/internal/db"
	"github.com/google/wire"
)

var RepositoriesSet = wire.NewSet(
	db.DbSet,
	NewCDCRepository,
	wire.Bind(new(listener.ICDCRepository), new(*CDCRepository)),
)
