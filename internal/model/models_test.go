package model_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/KazikovAP/merch_store/internal/model"
)

func TestAuthRequest_JSON(t *testing.T) {
	orig := model.AuthRequest{
		Username: "testuser",
		Password: "secret",
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	var result model.AuthRequest
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if !reflect.DeepEqual(orig, result) {
		t.Errorf("expected %+v, got %+v", orig, result)
	}
}
