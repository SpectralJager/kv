package kv

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPost(t *testing.T) {
	db, err := Open("test/test.kv")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = Post(db, "stringKey", "test string value")
	if err != nil {
		t.Fatal(err)
	}
}
func TestGet(t *testing.T) {
	db, err := Open("test/test.kv")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	value, err := Get[string](db, "stringKey")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, "test string value") {
		t.Fail()
	}
}
func TestPut(t *testing.T) {
	db, err := Open("test/test.kv")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = Put(db, "stringKey", "updates test string value")
	if err != nil {
		t.Fatal(err)
	}
}
func TestDelete(t *testing.T) {
	db, err := Open("test/test.kv")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = Delete(db, "stringKey")
	if err != nil {
		t.Fatal(err)
	}
}
func TestHead(t *testing.T) {
	db, err := Open("test/test.kv")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	keys := Head(db)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(keys)
}
