package postgres

const (
	GorReadData      = "GorReadData"
	GorWriteData     = "GorWriteData"
	Waiting          = "Waiting"
	First            = "first"
	Last             = "last"
	Stop             = "stop"
	Continue         = "continue"
	StartSync        = "start_sync"
	ErrorSync        = "error_sync"
	StopSync         = "stop_sync"
	RegularSync      = "regular_sinc"
	QueryTableStruct = "SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s';"
	// QueryReadData = "SELECT * FROM %s.%s WHERE id > 500 ORDER BY id limit 10;"
)
