package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LucasRoesler/openfaas-loki/pkg"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// ConfigHandlerFunc provides a debug endpoint to query for the server configuration
func ConfigHandlerFunc(w http.ResponseWriter, r *http.Request) {
	config := viper.AllSettings()
	config["version"] = pkg.Version
	config["commit"] = pkg.GitCommit

	bytes, err := json.Marshal(config)
	if err != nil {
		http.Error(w, "can not marshal config", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		logrus.WithError(err).Error("can not write config")
	}
}
