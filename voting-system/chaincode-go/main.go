package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for a simple voting system
type SmartContract struct {
	contractapi.Contract
}

// Election contains election metadata
type Election struct {
	Candidates []string `json:"candidates"`
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`
}

// Voter represents a voter
type Voter struct {
	HasVoted bool `json:"hasVoted"`
}

// VoteCount keeps track of votes for a candidate
type VoteCount struct {
	Count int `json:"count"`
}

// InitElection initializes the election with candidates and timing
func (s *SmartContract) InitElection(ctx contractapi.TransactionContextInterface, candidatesJSON string, startTime string, endTime string) error {
	// Parse the candidate list
	var candidates []string
	err := json.Unmarshal([]byte(candidatesJSON), &candidates)
	if err != nil {
		return fmt.Errorf("failed to parse candidates JSON: %v", err)
	}

	// Create an election object
	election := Election{

		Candidates: candidates,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	electionJSON, err := json.Marshal(election)
	if err != nil {
		return fmt.Errorf("failed to marshal election data: %v", err)
	}

	// Save the election to the ledger
	err = ctx.GetStub().PutState("election", electionJSON)
	if err != nil {
		return fmt.Errorf("failed to save election to ledger: %v", err)
	}

	// Initialize vote counts for each candidate
	for _, candidate := range candidates {
		voteCount := VoteCount{Count: 0}
		voteCountJSON, _ := json.Marshal(voteCount)
		err = ctx.GetStub().PutState(candidate, voteCountJSON)
		if err != nil {
			return fmt.Errorf("failed to initialize vote count for candidate %s: %v", candidate, err)
		}
	}

	return nil
}

// RegisterVoter registers a voter
func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, voterID string) error {
	// Check if the voter is already registered
	existingVoter, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("failed to read voter data: %v", err)
	}
	if existingVoter != nil {
		return fmt.Errorf("voter %s is already registered", voterID)
	}

	// Register the voter
	voter := Voter{HasVoted: false}
	voterJSON, _ := json.Marshal(voter)
	return ctx.GetStub().PutState(voterID, voterJSON)
}

// CastVote allows a voter to cast their vote
func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, candidateID string) error {
	// Check if the voter exists and hasn't voted yet
	voterJSON, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("failed to read voter data: %v", err)
	}
	if voterJSON == nil {
		return fmt.Errorf("voter %s is not registered", voterID)
	}

	var voter Voter
	err = json.Unmarshal(voterJSON, &voter)
	if err != nil {
		return fmt.Errorf("failed to parse voter data: %v", err)
	}
	if voter.HasVoted {
		return fmt.Errorf("voter %s has already voted", voterID)
	}

	// Check if the candidate exists
	candidateJSON, err := ctx.GetStub().GetState(candidateID)
	if err != nil {
		return fmt.Errorf("failed to read candidate data: %v", err)
	}
	if candidateJSON == nil {
		return fmt.Errorf("candidate %s does not exist", candidateID)
	}

	// Increment the candidate's vote count
	var voteCount VoteCount
	err = json.Unmarshal(candidateJSON, &voteCount)
	if err != nil {
		return fmt.Errorf("failed to parse candidate vote data: %v", err)
	}
	voteCount.Count++

	// Save the updated vote count
	voteCountJSON, _ := json.Marshal(voteCount)
	err = ctx.GetStub().PutState(candidateID, voteCountJSON)
	if err != nil {
		return fmt.Errorf("failed to update vote count for candidate %s: %v", candidateID, err)
	}

	// Mark the voter as having voted
	voter.HasVoted = true
	voterJSON, _ = json.Marshal(voter)
	return ctx.GetStub().PutState(voterID, voterJSON)
}

// GetResults retrieves the current vote count for all candidates
func (s *SmartContract) GetResults(ctx contractapi.TransactionContextInterface) (map[string]int, error) {
	// Get the election metadata
	electionJSON, err := ctx.GetStub().GetState("election")
	if err != nil {
		return nil, fmt.Errorf("failed to read election data: %v", err)
	}
	if electionJSON == nil {
		return nil, fmt.Errorf("no election data found")
	}

	var election Election
	err = json.Unmarshal(electionJSON, &election)
	if err != nil {
		return nil, fmt.Errorf("failed to parse election data: %v", err)
	}

	// Retrieve vote counts for all candidates
	results := make(map[string]int)
	for _, candidate := range election.Candidates {
		candidateJSON, err := ctx.GetStub().GetState(candidate)
		if err != nil {
			return nil, fmt.Errorf("failed to read candidate data: %v", err)
		}

		var voteCount VoteCount
		err = json.Unmarshal(candidateJSON, &voteCount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse candidate vote data: %v", err)
		}
		results[candidate] = voteCount.Count
	}

	return results, nil
}

// GetVoterStatus retrieves the voting status of a voter
func (s *SmartContract) GetVoterStatus(ctx contractapi.TransactionContextInterface, voterID string) (bool, error) {
	voterJSON, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return false, fmt.Errorf("failed to read voter data: %v", err)
	}
	if voterJSON == nil {
		return false, fmt.Errorf("voter %s is not registered", voterID)
	}

	var voter Voter
	err = json.Unmarshal(voterJSON, &voter)
	if err != nil {
		return false, fmt.Errorf("failed to parse voter data: %v", err)
	}

	return voter.HasVoted, nil
}

// GetElectionInfo retrieves the election metadata
func (s *SmartContract) GetElectionInfo(ctx contractapi.TransactionContextInterface) (*Election, error) {
	electionJSON, err := ctx.GetStub().GetState("election")
	if err != nil {
		return nil, fmt.Errorf("failed to read election data: %v", err)
	}
	if electionJSON == nil {
		return nil, fmt.Errorf("no election data found")
	}

	var election Election
	err = json.Unmarshal(electionJSON, &election)
	if err != nil {
		return nil, fmt.Errorf("failed to parse election data: %v", err)
	}

	return &election, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
