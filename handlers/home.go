package handlers

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r,
		"https://github.com/martinyunify/discovery.etcd.io/wiki",
		http.StatusMovedPermanently,
	)
}
