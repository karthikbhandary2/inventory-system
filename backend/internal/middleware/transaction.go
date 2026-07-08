package middleware

import (
    "context"
    "net/http"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string
const TxKey contextKey = "db_tx"

// TxFromContext retrieves the transaction from context.
// Handlers call this instead of using the pool directly.
func TxFromContext(ctx context.Context) pgx.Tx {
    tx, _ := ctx.Value(TxKey).(pgx.Tx)
    return tx
}

// Transaction wraps the handler in a DB transaction.
// Commits on success, rolls back on any error or panic.
func Transaction(pool *pgxpool.Pool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tx, err := pool.Begin(r.Context())
            if err != nil {
                http.Error(w, "could not begin transaction", http.StatusInternalServerError)
                return
            }

            // Always rollback if we exit early (panic, error response)
            defer func() {
                if p := recover(); p != nil {
                    _ = tx.Rollback(r.Context())
                    panic(p) // re-panic after rollback
                }
            }()

            // Inject tx into context so handlers can use it
            ctx := context.WithValue(r.Context(), TxKey, tx)
            
            // Use a response recorder to detect error status codes
            rw := &responseWriter{ResponseWriter: w}
            next.ServeHTTP(rw, r.WithContext(ctx))

            if rw.status >= 400 {
                _ = tx.Rollback(r.Context())
            } else {
                if err := tx.Commit(r.Context()); err != nil {
                    http.Error(w, "commit failed", http.StatusInternalServerError)
                }
            }
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}