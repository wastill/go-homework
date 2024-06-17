package go_ethcli

import "testing"

func TestQueryBlock(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QueryBlock()
		})
	}
}

func TestQueryTransaction(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QueryTransaction()
		})
	}
}

func TestTransferToken(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TransferToken()
		})
	}
}

func TestPublishBlock(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PublishBlock()
		})
	}
}
