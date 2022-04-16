# sqls

SQL helper functions for queries and transactions to remove common boilerplate code.

```go
import (
	"context"
	"database/sql"
	"github.com/tomcz/gotools/sqls"
)

func main() {
	ctx := context.Background()

	var db *sql.DB // initialisation omitted

	results := make(map[string]string)

	selectLeadersSQL := "select leader, node from current_leaders"
	err := sqls.QueryRowsContext(ctx, db, selectLeadersSQL)(func(row sqls.ScanFunc) error {
		var leader, node string
		if err := row(&leader, &node); err != nil {
			return nil
		}
		results[leader] = node
		return nil
	})
	// error handling omitted

	insertLeaderSQL := "insert into old_leaders (leader, node, created_at) values (?, ?, ?)"
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
}
```
