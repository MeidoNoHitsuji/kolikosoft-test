package shutdowner

const (
	// PriorityLayerWeb: закрытие web-сервиса
	PriorityLayerWeb = 0
	// PriorityLayerStorage: закрытие persistent-хранилищ/кэша
	PriorityLayerStorage = 1
)
