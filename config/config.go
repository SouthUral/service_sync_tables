package config

import (
	tools "github.com/SouthUral/service_sync_tables/tools"
	log "github.com/sirupsen/logrus"
)

// Метод для получения данных из конфига
func GetConf() (ConfEnum, interface{}) {
	conf := ConfEnum{}
	err := tools.JsonRead(&conf, "conf.json")
	if err != nil {
		log.Error("Не загружены данные из конфига")
		return conf, err
	}
	return conf, nil
}
