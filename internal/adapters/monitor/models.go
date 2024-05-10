package monitor

type IMonitor interface {
	IncPublishedEventsCounter(db, coll string)
	IncProblematicEventsCounter(kind string)
	IncFilteredEventsCounter(db, coll string)
}
