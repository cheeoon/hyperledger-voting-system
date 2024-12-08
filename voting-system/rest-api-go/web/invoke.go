package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err), http.StatusBadRequest)
		return
	}

	// Extract parameters
	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"] // Handles multiple args

	fmt.Printf("Channel: %s, Chaincode: %s, Function: %s, Args: %s\n", channelID, chainCodeName, function, args)

	// Connect to network and invoke chaincode
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	txnProposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating transaction proposal: %s", err), http.StatusInternalServerError)
		return
	}

	txnEndorsed, err := txnProposal.Endorse()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error endorsing transaction: %s", err), http.StatusInternalServerError)
		return
	}

	txnCommitted, err := txnEndorsed.Submit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error submitting transaction: %s", err), http.StatusInternalServerError)
		return
	}

	// Prepare JSON response
	response := string(txnEndorsed.Result())
	if response == "" {
		response = "No response from chaincode."
	}
	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txnCommitted.TransactionID(), response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
