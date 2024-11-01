package server

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
)

//go:embed views/*
var viewsFS embed.FS

// Data represents the server of the application
type Data struct {
	mux   *http.ServeMux
	db    *sql.DB
	pool  *pgxpool.Pool
	pages *template.Template
}

// New creates a new web server
func New(ctx context.Context) (*Data, error) {
	var err error

	pages, err := template.New("").ParseFS(viewsFS, "views/*.html")
	if err != nil {
		return nil, fmt.Errorf("while parsing embedded pages: %w", err)
	}

	pool, err := pgxpool.New(ctx, viper.GetString("connection-string"))
	if err != nil {
		return nil, fmt.Errorf("while connecting to PostgreSQL: %w", err)
	}

	mux := http.NewServeMux()
	result := &Data{
		mux:   mux,
		pool:  pool,
		pages: pages,
	}

	mux.HandleFunc("/", result.indexPage)
	mux.HandleFunc("/readyz", result.readinessProbe)
	return result, nil
}

// Start starts the web server
func (d *Data) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    viper.GetString("listen"),
		Handler: d.mux,
	}

	pool := pool.
		New().
		WithErrors().
		WithContext(ctx).
		WithCancelOnError()
	pool.Go(func(_ context.Context) error {
		log.Println("Starting web server on", server.Addr)
		err := server.ListenAndServe()
		log.Printf("Shut down web server (%s)\n", err)

		if err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	pool.Go(func(ctx context.Context) error {
		<-ctx.Done()
		server.Shutdown(ctx)
		return nil
	})
	return pool.Wait()
}

func (d *Data) indexPage(w http.ResponseWriter, r *http.Request) {
	tx, err := d.pool.Begin(r.Context())
	if err != nil {
		log.Println("Error connecting to the database", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		err := tx.Rollback(r.Context())
		if err != nil {
			log.Println("Error while rolling back a transaction", err)
		}
	}()

	rows, err := tx.Query(r.Context(), "SELECT * FROM schema_migrations")
	if err != nil {
		log.Println("Error getting the migrations", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	type Migration struct {
		Version int64
		Dirty   bool
	}

	migrations, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Migration])
	if err != nil {
		log.Println("Error collecting the migrations", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := d.pages.ExecuteTemplate(w, "index.html", migrations); err != nil {
		log.Println("Error executing the page template", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (d *Data) readinessProbe(w http.ResponseWriter, r *http.Request) {
	tx, err := d.pool.Begin(r.Context())
	if err != nil {
		log.Println("Error connecting to the database", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		err := tx.Rollback(r.Context())
		if err != nil {
			log.Println("Error while rolling back the transaction", err)
		}
	}()

	w.Write([]byte("ok"))
}
