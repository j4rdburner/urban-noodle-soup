package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"os"
)

// Struct to hold the response data
type CheckRequirementResponse struct {
	Address string `json:"address"`
	Passed  bool   `json:"passed"`
}

func queryBigQuery(ctx context.Context, client *bigquery.Client, address string) (float64, error) {
	// Query definition
	query := fmt.Sprintf(`
        SELECT
            MAX(balance) AS max_balance
        FROM
            `+"`bigquery-public-data.crypto_ethereum.balances`"+`
        WHERE
            address = '%s'
    `, address)

	// Run the query
	q := client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return 0, err
	}

	// Process the results
	var maxBalance float64
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		if len(values) > 0 {
			maxBalance = values[0].(float64)
		}
	}

	return maxBalance, nil
}

func requirementHandler(w http.ResponseWriter, r *http.Request) {
	// Get the wallet address from the query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	// Get the threshold from the query parameters
	thresholdString := r.URL.Query().Get("threshold")
	if threshold == "" {
		http.Error(w, "Missing threshold parameter", http.StatusBadRequest)
		return
	}

	// Get the modifier from the query parameters
	modifierString := r.URL.Query().Get("modifier")
	if modifier == "" {
		http.Error(w, "Missing modifier parameter", http.StatusBadRequest)
		return
	}

	// Parse numerical arguments
	threshold, err := strconv.ParseFloat(thresholdString, 64)
	if err != nil {
		http.Error(w, "Failed to parse modifier threshold as numeric value:", http.StatusBadRequest)
		return
	}

	modifier, err := strconv.ParseFloat(modifierString, 64)
	if err != nil {
		http.Error(w, "Failed to parse modifier parameter as numeric value:", http.StatusBadRequest)
		return
	}

	// Create a BigQuery client
	ctx := context.Background()
	projectID := "j4rdburner"
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		http.Error(w, "Failed to create BigQuery client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Query BigQuery for the maximum balance
	maxBalance, err := queryBigQuery(ctx, client, address)
	if err != nil {
		http.Error(w, "Failed to query BigQuery", http.StatusInternalServerError)
		return
	}

	// Check if the address meets the required maxBalance
	modifiedThreshold, modifiedMaxBalance := threshold*modifier, maxBalance*modifier
	passed := modifiedMaxBalance > modifiedThreshold

	// Create the response
	response := CheckRequirementResponse{
		Address: address,
		Passed:  passed,
	}

	// Write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/checkRequirement", requirementHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
