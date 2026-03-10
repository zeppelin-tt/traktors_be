package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const featureDocID = "feature_flag"

var featureCollection *mongo.Collection

type featureDoc struct {
	ID    string `bson:"_id"`
	Value bool   `bson:"value"`
}

// GET /feature_check
func featureCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var doc featureDoc
	err := featureCollection.FindOne(ctx, bson.M{"_id": featureDocID}).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		writeError(w, http.StatusInternalServerError, "failed to fetch feature flag")
		return
	}

	writeJSON(w, http.StatusOK, doc.Value)
}

// POST /feature_set
func featureSet(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Value bool `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	_, err := featureCollection.UpdateOne(
		ctx,
		bson.M{"_id": featureDocID},
		bson.M{"$set": bson.M{"value": body.Value}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update feature flag")
		return
	}

	writeJSON(w, http.StatusOK, body.Value)
}
