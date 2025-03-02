package tests

import (
	"encoding/json"
	"testing"

	"github.com/asullivan219/receiptProcessor/internal/models"
	"github.com/asullivan219/receiptProcessor/internal/store"
	"github.com/google/uuid"
)

var s = store.NewStore("../testResources/database.db")

func TestStore(t *testing.T) {

	for _, td := range RECEIPT_TEST_DATA {
		var receipt models.Receipt
		json.Unmarshal(td.jsonBytes, &receipt)
		validReceipt, _ := receipt.ValidateReceipt()
		score := validReceipt.ScoreReceipt()

		id := uuid.New()

		dbReceipt := store.NewDbReceipt(
			id.String(),
			string(td.jsonBytes),
			score,
		)

		err := s.PutReceipt(dbReceipt)
		if err != nil {
			t.Fatalf("Error putting receipt in DB: %s", err.Error())
		}

		_, err = s.GetReceipt(id.String())
		if err != nil {
			t.Fatalf("Error getting receipt from DB: %s", err.Error())
		}

	}
}
