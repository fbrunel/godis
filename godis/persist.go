package godis

import (
	"encoding/gob"
	"io"
	"os"
)

func LoadStoreFromFile(name string) (*StandardStore, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadStore(file)
}

func LoadStore(reader io.Reader) (*StandardStore, error) {
	decoder := gob.NewDecoder(reader)
	var hmap map[string]string
	err := decoder.Decode(&hmap)
	if err != nil {
		return nil, err
	}
	return &StandardStore{hmap: hmap}, nil
}

func SaveStoreToFile(store *StandardStore, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return SaveStore(store, file)
}

func SaveStore(store *StandardStore, writer io.Writer) error {
	encoder := gob.NewEncoder(writer)
	return encoder.Encode(store.hmap)
}
