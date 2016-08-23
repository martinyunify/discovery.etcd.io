package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var cfg = client.Config{
	Endpoints: []string{"http://127.0.0.1:2379"},
	Transport: client.DefaultTransport,
	// set timeout per request to fail fast when the target endpoint is unavailable
	HeaderTimeoutPerRequest: time.Second,
}

var baseURI = flag.String("host", "https://discovery.etcd.io", "base location for computed token URI")

func generateCluster() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func setupToken(size int) (string, error) {
	token := generateCluster()
	if token == "" {
		return "", errors.New("Couldn't generate a token")
	}

	c, _ := client.New(cfg)
	kapi := client.NewKeysAPI(c)

	key := path.Join("_etcd", "registry", token)

	resp, err := kapi.Create(context.Background(), path.Join(key, "_config", "size"), strconv.Itoa(size))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Couldn't setup state %v %v", resp, err))
	}

	return token, nil
}

func deleteToken(token string) error {
	c, _ := client.New(cfg)
	kapi := client.NewKeysAPI(c)

	if token == "" {
		return errors.New("No token given")
	}

	_, err := kapi.Delete(
		context.Background(),
		path.Join("_etcd", "registry", token),
		&client.DeleteOptions{Recursive: true},
	)

	return err
}

func NewTokenHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	size := 3
	s := r.FormValue("size")
	if s != "" {
		size, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	token, err := setupToken(size)

	if err != nil {
		log.Printf("setupToken returned: %v", err)
		http.Error(w, "Unable to generate token", 400)
		return
	}

	log.Println("New cluster created", token)

	fmt.Fprintf(w, "%s/%s", bytes.TrimRight([]byte(*baseURI), "/"), token)
}
