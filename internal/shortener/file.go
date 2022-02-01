package shortener

import (
	"encoding/json"
	"os"
)

func (s *Shortener) Restore(file string) error {
	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewDecoder(fp).Decode(&s.userLinks)
}

func (s *Shortener) Save(file string) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(s.userLinks)
}
