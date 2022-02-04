package shortener

import (
	"encoding/json"
	"go-musthave-shortener-tpl/internal/entity"
	"os"
)

func OpenFile(file string) (*os.File, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return fp, nil
}

func Restore(file string, userLinks map[string][]entity.UserLinks) error {
	fp, err := OpenFile(file)
	if err != nil {
		return err
	}
	return json.NewDecoder(fp).Decode(&userLinks)
}

func SaveInFile(file string, userLinks map[string][]entity.UserLinks) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(userLinks)
}
