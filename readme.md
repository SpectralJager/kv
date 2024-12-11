# KVDB is simple persistent key-value store 

# Install
```
go get github.com/SpectralJager/kv
```

# Supported types
Because of gob serialization, values could be any builtin golang types.

# Api
- Open -- open existed or create a new database file:
```go
db, err := kv.Open("example/example.kv")
```
- Post -- create a new value for specific key. Doesn't allowed to create value for already existed key:
```go
err := kv.Post(db, "name", "Daniel") // "name" -> "Daniel"
err := kv.Post(db, "name", "Alex") // error: value for key 'name' already exist
```
- Get -- get value for specific key:
```go
name, err := kv.Get[string](db, "name") // name = "Daniel"
```
- Put -- update value for existed key:
```go
err := kv.Put(db, "name", "Alex") // old: "name" -> "Daniel; new: "name" -> "Alex"
```
- Head -- Get all available store's keys:
```go
keys := kv.Head(db) // keys = ["name"]
```
- Delete -- remove specific key and it's value from database:
```go
err := kv.Delete(db, "name") 
keys := kv.Head(db) // keys = []
```