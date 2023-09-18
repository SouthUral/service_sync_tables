package postgres

// функция отправляет сообщение в postgres для старта синхронизации
func StartSyncPg(connChPg IncomCh, DBalias, Table, Schema, Offset string, Clean bool) CommToSync {
	connSyncCh := make(CommToSync)
	messToPg := IncomingMess{
		Table:      Table,
		Schema:     Schema,
		Database:   DBalias,
		Offset:     Offset,
		ChCommSync: connSyncCh,
		Clean:      Clean,
	}
	connChPg <- messToPg
	return connSyncCh
}
