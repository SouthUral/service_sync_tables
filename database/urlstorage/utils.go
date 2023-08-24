package urlstorage

func CopyMap(data StorageConnDB) StorageConnDB {
	copyMap := make(StorageConnDB)
	for key, value := range data {
		copyMap[key] = value
	}
	return copyMap
}
