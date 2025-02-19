package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/docs"
	wordfilter "github.com/WhileCodingDoLearn/bootdev_server_tut/filter"
	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"

/*
go build -o out && ./out

You may be wondering how the fileserver knew to serve the index.html file to the root of the server.
It’s such a common convention on the web to use a file called index.html to serve the webpage for a given path,
that the Go standard library’s FileServer does it automatically.


Postgresql:

1.
sudo service postgresql start
2.
sudo -u postgres psql
3:
psql \c dbname

goose:

goose postgres postgres://postgres:postgres@localhost:5432/chirpy up/down

Secret generator

openssl rand -base64 64
*/

type apiConfig struct {
	filserverHits  atomic.Int32
	dbQueries      *database.Queries
	plattform      string
	filter         wordfilter.MessageFilter
	secret         string
	validTokenTime time.Duration
	webhookKey     string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	portFromEnv := os.Getenv("PORT")

	dbUrl := os.Getenv("DB_URL")
	db, errDB := sql.Open("postgres", dbUrl)
	if errDB != nil {
		log.Println(errDB)
	}
	plattformType := os.Getenv("PLATFORM")

	secretFromEnv := os.Getenv("SECRET")

	WebhookKeyFromEnv := os.Getenv("POLKA_KEY")

	apiMiddleware := apiConfig{
		filserverHits: atomic.Int32{},
		dbQueries:     database.New(db),
		plattform:     plattformType,
		filter: wordfilter.MessageFilter{
			WordFilter: map[string]bool{
				"kerfuffle": true,
				"sharbert":  true,
				"fornax":    true,
			},
		},
		secret:         secretFromEnv,
		validTokenTime: 1 * time.Hour,
		webhookKey:     WebhookKeyFromEnv,
	}
	smux := http.NewServeMux()

	/*App________________________________________________________________________________________________*/

	smux.Handle("/app/", http.StripPrefix("/app", apiMiddleware.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	/*Api________________________________________________________________________________________________*/

	smux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	smux.HandleFunc("POST /api/users", apiMiddleware.createUser)
	smux.HandleFunc("PUT /api/users", apiMiddleware.isAutneticated(apiMiddleware.updateUser))
	smux.HandleFunc("DELETE /api/users", apiMiddleware.isAutneticated(apiMiddleware.deleteUser))

	smux.HandleFunc("POST /api/polka/webhooks", apiMiddleware.handleEvent)

	smux.HandleFunc("POST /api/login", apiMiddleware.loginUser)
	smux.HandleFunc("POST /api/refresh", apiMiddleware.RefreshToken)
	smux.HandleFunc("POST /api/revoke", apiMiddleware.RevokeToken)

	smux.HandleFunc("GET /api/chirps/", apiMiddleware.GetAllChirps)
	smux.HandleFunc("GET /api/chirps/{chirpID}", apiMiddleware.GetChirp)
	smux.HandleFunc("DELETE /api/chirps/{chirpID}", apiMiddleware.isAutneticated(apiMiddleware.DeleteChirp))

	smux.HandleFunc("POST /api/chirps", apiMiddleware.isAutneticated(apiMiddleware.ChirpHandler))

	/*Admin________________________________________________________________________________________________*/

	smux.HandleFunc("GET /admin/metrics", apiMiddleware.handleMetrics)

	smux.HandleFunc("POST /admin/reset", apiMiddleware.reset)

	var docs = docs.GetDocsGenerator().GenerateDocs()

	smux.HandleFunc("GET /admin/doc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(docs))
	})

	server := http.Server{Handler: smux}
	server.Addr = ":" + portFromEnv
	errServer := server.ListenAndServe()
	if errServer != nil {
		log.Fatal(errServer)
	}

}
