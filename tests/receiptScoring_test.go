package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/asullivan219/receiptProcessor/internal/models"
)

var RECEIPT_TEST_DATA = []ReceiptTestData{{
	jsonBytes: []byte(`
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
    `),
	score: 28,
}, {
	jsonBytes: []byte(`{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}`),
	score: 109,
}}

func TestScoreReceipt(t *testing.T) {
	for _, td := range RECEIPT_TEST_DATA {
		var receipt models.Receipt
		json.Unmarshal(td.jsonBytes, &receipt)
		validReceipt, _ := receipt.ValidateReceipt()
		fmt.Println(validReceipt)
		score := validReceipt.ScoreReceipt()
		if score != td.score {
			t.Fatalf("Error score was %d, expected %d", score, td.score)
		}
	}
}
