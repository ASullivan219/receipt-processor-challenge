package tests

import (
	"encoding/json"
	"testing"

	"github.com/asullivan219/receiptProcessor/internal/models"
)

type ReceiptTestData struct {
	jsonBytes []byte
	total     float64
	score     int
}

var VALID_RECEIPT_TEST_DATA = []ReceiptTestData{{
	jsonBytes: []byte(`{
    "retailer": "Walgreens",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "08:13",
    "total": "2.65",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
        {"shortDescription": "Dasani", "price": "1.40"}
    ]
}`),
	total: 2.65,
}, {
	jsonBytes: []byte(`{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
}`),
	total: 1.25,
}}

func TestIngestJSON(t *testing.T) {

	for _, td := range VALID_RECEIPT_TEST_DATA {
		var receipt models.Receipt
		json.Unmarshal(td.jsonBytes, &receipt)
	}

}

func TestValidateReceipt(t *testing.T) {
	for _, td := range VALID_RECEIPT_TEST_DATA {
		var receipt models.Receipt
		json.Unmarshal(td.jsonBytes, &receipt)

		vr, err := receipt.ValidateReceipt()
		if err != nil {
			t.Fatal("Error validating a valid receipt")
		}

		if vr.Total != td.total {
			t.Fatalf("Expected total of: %f got %f", td.total, vr.Total)
		}
	}
}
