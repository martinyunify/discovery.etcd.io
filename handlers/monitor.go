package handlers

import (
	"net/http"
)

var monitorUrl string

func SetupMonitorAddr(monitor string) {
	monitorUrl = monitor
}

func MonitorHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, monitorUrl,
		http.StatusMovedPermanently,
	)
}
