package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type Version struct {
	Version string `json:"version"`
}

func GetActualVersion(addr string) (*string, error) {
	resp, err := http.Get(addr + "/irb/ver")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v Version

	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	return &v.Version, nil
}
func GetClient(addr, verName string) error {
	resp, err := http.Get(addr + "/irb/get")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	os.Mkdir("clients", 0777)
	f, err := os.Create("clients/" + verName + ".zip")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
