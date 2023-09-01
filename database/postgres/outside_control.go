package postgres

// функция отправляет сообщение в postgres для старта синхронизации
func StartSyncPg(connChPg IncomCh, DBalias, Table, Offset string) CommToSync {
	connSyncCh := make(CommToSync)
	messToPg := IncomingMess{
		Table:      Table,
		Database:   DBalias,
		Offset:     Offset,
		ChCommSync: connSyncCh,
	}
	connChPg <- messToPg
	return connSyncCh
}
