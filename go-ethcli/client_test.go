package go_ethcli

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewClient()
		})
	}
}

func TestAccount(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"TestAccount"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Account()
		})
	}
}

func TestGolemErc20(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GolemErc20()
		})
	}
}

func TestNewWallet(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewWallet()
		})
	}
}

func TestGenKeyStore(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenKeyStore()
		})
	}
}

func TestImportKeyStore(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ImportKeyStore()
		})
	}
}
