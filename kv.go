package kv

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"
)

type DB struct {
	mtx   sync.RWMutex
	file  *os.File
	store map[string][]byte
	close chan struct{}
}

func Open(path string) (*DB, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return &DB{}, fmt.Errorf("can't open database: %w", err)
	}
	db := DB{
		file:  file,
		store: map[string][]byte{},
		close: make(chan struct{}),
		mtx:   sync.RWMutex{},
	}
	stat, err := file.Stat()
	if err != nil {
		return &db, fmt.Errorf("can't statistics of opened database: %w", err)
	}
	if stat.Size() == 0 {
		err = db.encode()
		if err != nil {
			return &db, fmt.Errorf("can't init database: %w", err)
		}
	} else {
		err = db.decode()
		if err != nil {
			return &db, fmt.Errorf("can't decode opened database: %w", err)
		}
	}
	go db.sync()
	return &db, nil
}

func (db *DB) Close() error {
	db.close <- struct{}{}
	return db.file.Close()
}

func (db *DB) sync() {
	for {
		select {
		case <-db.close:
		default:
			db.file.Sync()
			db.decode()
			time.Sleep(time.Second * 10)
		}
	}
}

func (db *DB) decode() error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	return gob.NewDecoder(db.file).Decode(&db.store)
}

func (db *DB) encode() error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	_, err := db.file.Seek(0, 0)
	if err != nil {
		return err
	}
	return gob.NewEncoder(db.file).Encode(db.store)
}

func (db *DB) getRaw(key string) ([]byte, error) {
	if data, ok := db.store[key]; ok {
		return data, nil
	}
	return []byte{}, fmt.Errorf("can't find data for key '%s'", key)
}

func (db *DB) setRaw(key string, value []byte) error {
	db.store[key] = value
	return db.encode()
}

func Get[T any](db *DB, key string) (T, error) {
	var value T
	data, err := db.getRaw(key)
	if err != nil {
		return value, fmt.Errorf("can't get value: %w", err)
	}
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&value)
	if err != nil {
		return value, fmt.Errorf("can't decode data to %s: %w", reflect.TypeOf(value).String(), err)
	}
	return value, nil
}

func Post[T any](db *DB, key string, value T) error {
	_, err := db.getRaw(key)
	if err == nil {
		return fmt.Errorf("key '%s' already exists in store", key)
	}
	return Put(db, key, value)
}

func Put[T any](db *DB, key string, value T) error {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return fmt.Errorf("can't encode data: %w", err)
	}
	err = db.setRaw(key, buf.Bytes())
	if err != nil {
		return fmt.Errorf("can't set data to store: %w", err)
	}
	return nil
}

func Delete(db *DB, key string) error {
	delete(db.store, key)
	return db.encode()
}

func Head(db *DB) []string {
	keys := []string{}
	for key := range db.store {
		keys = append(keys, key)
	}
	return keys
}
