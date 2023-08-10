package state

func CopyMap(data StateStorage) StateStorage {
	copyMap := make(StateStorage)
	for key, value := range data {
		copyMap[key] = value
	}
	return copyMap
}
