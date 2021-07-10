# gotools

Miscellaneous tools for golang applications.

## env

Populate structs from environment variables using [mapstructure](https://github.com/mitchellh/mapstructure).

```go
import "github.com/tomcz/gotools/env"

type testCfg struct {
    User  string `mapstructure:"test_user_name"`
    Age   int    `mapstructure:"test_user_age"`
    Admin bool   `mapstructure:"test_is_admin"`
    Port  int    `mapstructure:"port"`
}

os.Setenv("test_user_name", "Homer")
os.Setenv("test_user_age", "42")
os.Setenv("test_is_admin", "true")

cfg := &testCfg{
    Port: 8080,
}

err := env.PopulateFromEnv(cfg)

if assert.NoError(t, err) {
    assert.Equal(t, "Homer", cfg.User)
    assert.Equal(t, 42, cfg.Age)
    assert.Equal(t, true, cfg.Admin)
    assert.Equal(t, 8080, cfg.Port)
}
```

## leader

Good-enough [leader election](https://aws.amazon.com/builders-library/leader-election-in-distributed-systems/)
pattern implementation using MySQL as the coordinator.

Builds on the ideas found in [this gist](https://gist.github.com/ljjjustin/f2213ac9b9b8c31df746f8b56095ea32).

```go
import (
    "context"
    "database/sql"
    "errors"
    "log"
    "sync"
    "time"

    "github.com/tomcz/gotools/leader"
)

func main() {
    var db *sql.DB // initialisation omitted

    check := leader.NewMysqlLeader(db, "app_leader")

    ctx, cancelElections := context.WithCancel(context.Background())

    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        for {
            err := check.StartElections(ctx)
            if errors.Is(err, context.Canceled) {
                log.Println("shutting down elections")
                wg.Done()
                return
            }
            log.Println("elections will continue despite errors:", err)
        }
    }()

    go func() {
        for i := 0; i < 100; i++ {
            if i > 0 {
                time.Sleep(10 * time.Second)
            }
            isLeader, err := check.IsLeader(ctx)
            if err != nil {
                log.Println("leader check failed:", err)
                continue
            }
            if isLeader {
                log.Println("I am the Leader!")
                continue
            }
            log.Println("I am NOT the Leader.")
        }
        cancelElections()
    }()

    wg.Wait()
}
```

## sets

Generated code of sets based on a map of type-to-bools for a range of standard go types.

Generics in Go would be nice, but here we are.

* Contains
* ContainsAny
* ContainsAll
* SubsetOf
* Union
* Intersection
* Difference
* Ordered
* MarshalJSON
* UnmarshalJSON

## slices

Generated code of common slice operations based on standard go types.

Generics in Go would be nice, but here we are.

* `Split[X]` splits a slice into parts of a given length, with a remainder if necessary.
* `Convert[X]ToInterface` converts a slice of one type to a slice of `interface{}`.
* `Append[X]ToInterface` appends the contents of a slice to a slice of `interface{}`.

## sqls

SQL helper functions for queries and transactions to remove common boilerplate code.

```go
import "github.com/tomcz/gotools/sqls"

selectLeadersSQL = "select leader, node from current_leaders"
insertLeaderSQL = "insert into old_leaders (leader, node, created_at) values (?, ?, ?)"

ctx := context.Background()

var db *sql.DB // initialisation omitted

results := make(map[string]string)
err := sqls.QueryRowsContext(ctx, db, selectLeadersSQL)(func(row sqls.ScanFunc) error {
    var leader, node string
    if err := row(&leader, &node); err != nil {
        return nil
    }
    results[leader] = node
    return nil
})
// error handling omitted

err = sqls.InTxContext(ctx, db, func(tx *sql.Tx) error {
    for leader, node := range results {
        _, err := tx.ExecContext(ctx, insertLeaderSQL, leader, node, time.Now())
        if err != nil {
            return err // tx will be rolled-back
        }
    }
    return nil // tx will be committed
})
// error handling omitted
```
