package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection
	uploadDir  string
	baseURL    string
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getenv("DB_NAME", "traktors")
	collName := getenv("COLLECTION_NAME", "tractors")
	port := getenv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	log.Println("connected to MongoDB")

	collection = client.Database(dbName).Collection(collName)

	uploadDir = getenv("UPLOAD_DIR", "./uploads")
	baseURL = getenv("BASE_URL", "http://localhost:8080")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("create upload dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tractors", getAllTractors)
	mux.HandleFunc("GET /tractors/{id}", getByID)
	mux.HandleFunc("POST /tractors", createTractor)
	mux.HandleFunc("PUT /tractors/{id}", updateTractor)
	mux.HandleFunc("DELETE /tractors/{id}", deleteTractor)

	mux.HandleFunc("POST /media", uploadImage)
	mux.Handle("GET /media/", http.StripPrefix("/media/", http.FileServer(http.Dir(uploadDir))))

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
