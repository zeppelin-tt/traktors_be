package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// normalizeTractor гарантирует images == [] а не null.
func normalizeTractor(t *Tractor) {
	if t.Images == nil {
		t.Images = []string{}
	}
}

// GET /tractors
func getAllTractors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch tractors")
		return
	}
	defer cursor.Close(ctx)

	var list []Tractor
	if err := cursor.All(ctx, &list); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to decode tractors")
		return
	}

	if list == nil {
		list = []Tractor{}
	}
	for i := range list {
		normalizeTractor(&list[i])
	}

	writeJSON(w, http.StatusOK, list)
}

// GET /tractors/{id}
func getByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var tractor Tractor
	err := collection.FindOne(ctx, bson.M{"_id": r.PathValue("id")}).Decode(&tractor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeError(w, http.StatusNotFound, "tractor not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch tractor")
		return
	}

	normalizeTractor(&tractor)
	writeJSON(w, http.StatusOK, tractor)
}

// POST /tractors
func createTractor(w http.ResponseWriter, r *http.Request) {
	var tractor Tractor
	if err := json.NewDecoder(r.Body).Decode(&tractor); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	if tractor.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	if tractor.ID == "" {
		tractor.ID = primitive.NewObjectID().Hex()
	}
	normalizeTractor(&tractor)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if _, err := collection.InsertOne(ctx, tractor); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			writeError(w, http.StatusConflict, "ID already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create tractor")
		return
	}

	writeJSON(w, http.StatusCreated, tractor)
}

// PUT /tractors/{id} — полная замена документа.
func updateTractor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var tractor Tractor
	if err := json.NewDecoder(r.Body).Decode(&tractor); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	if tractor.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tractor.ID = id
	normalizeTractor(&tractor)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, tractor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update tractor")
		return
	}
	if result.MatchedCount == 0 {
		writeError(w, http.StatusNotFound, "tractor not found")
		return
	}

	writeJSON(w, http.StatusOK, tractor)
}

// DELETE /tractors/{id}
func deleteTractor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": r.PathValue("id")})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete tractor")
		return
	}
	if result.DeletedCount == 0 {
		writeError(w, http.StatusNotFound, "tractor not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
