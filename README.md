# gorand
threadsafe math/rand NewSource

## Usage

```go
  rnd := rand.New(gorand.NewSource(time.Now().UnixNano()))
  log.Println(rnd.Int63())
```
