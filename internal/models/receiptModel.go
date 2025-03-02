package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/asullivan219/receiptProcessor/common"
)

const (
	RECEIPT_DATE_LAYOUT = "2006-01-02 15:04"
)

var MISMATCHED_TOTAL_ERROR = errors.New("MISMATCHED_TOTAL")

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type ValidReceipt struct {
	Retailer         string
	PurchaseDateTime time.Time
	Total            float64
	Items            []ValidItem
}

type ValidItem struct {
	ShortDescription string
	Price            float64
}

func (r *Receipt) ValidateReceipt() (ValidReceipt, error) {
	var vr ValidReceipt

	purchaseDateTime, err := time.Parse(
		RECEIPT_DATE_LAYOUT,
		fmt.Sprintf("%s %s", r.PurchaseDate, r.PurchaseTime))

	if err != nil {
		return vr, err
	}

	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return vr, err
	}

	itemList := make([]ValidItem, 0, len(r.Items))

	itemTotal := 0.0
	for _, item := range r.Items {
		validItem, err := item.ValidateItem()
		if err != nil {
			return vr, err
		}
		itemTotal += validItem.Price
		itemList = append(itemList, validItem)
	}

	if !common.FloatsEqual(itemTotal, total) {
		return vr, MISMATCHED_TOTAL_ERROR
	}

	vr.Retailer = r.Retailer
	vr.PurchaseDateTime = purchaseDateTime
	vr.Total = total
	vr.Items = itemList

	return vr, nil
}

func (i *Item) ValidateItem() (ValidItem, error) {
	var vi ValidItem

	price, err := strconv.ParseFloat(i.Price, 64)
	if err != nil {
		return vi, err
	}

	vi.Price = price
	vi.ShortDescription = i.ShortDescription
	return vi, nil
}
