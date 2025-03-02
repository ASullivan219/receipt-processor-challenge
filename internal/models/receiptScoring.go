package models

import (
	"log/slog"
	"strings"
	"unicode"
)

const (
	ROUND_SCORE_BONUS       = 50
	QUARTER_SCORE_BONUS     = 25
	ALPHANUMERIC_BONUS      = 1
	ITEM_COUPLE_BONUS       = 5
	ITEM_DESCRIPTION_FACTOR = 0.2
	ODD_DAY_SCORE           = 6
	PURCHASE_TIME_BONUS     = 10
)

// Score the current receipt
func (vr *ValidReceipt) ScoreReceipt() int {
	score := 0
	score += scoreAlphanumericCharacter(vr.Retailer)
	score += scoreDecimal(vr.Total)
	for _, item := range vr.Items {
		score += scoreItemDescription(item.ShortDescription, item.Price)
	}
	score += scoreTwoItemBonus(len(vr.Items))
	score += scoreDayOfMonth(vr.PurchaseDateTime.Day())
	score += scoreTimeOfDay(vr.PurchaseDateTime.Hour(), vr.PurchaseDateTime.Minute())
	return score
}

// One point for every alphanumeric character in the retailer name.
func scoreAlphanumericCharacter(retailerName string) int {
	totalAlphnumerics := 0
	for _, curr := range retailerName {
		if unicode.IsDigit(curr) || unicode.IsLetter(curr) {
			totalAlphnumerics += 1
		}
	}
	slog.Info("score alphnumerics",
		"total alphas", totalAlphnumerics,
		"bonus", totalAlphnumerics*ALPHANUMERIC_BONUS)
	return totalAlphnumerics * ALPHANUMERIC_BONUS
}

// 50 points if the total is a round dollar amount with no cents.
// 25 points if the total is a multiple of 0.25
func scoreDecimal(total float64) int {
	// Remove Decimal
	t := total * 100
	intT := int(t)

	// Check if round number
	score := 0
	if intT%10 == 0 {
		slog.Info("round score bonus", "bonus", ROUND_SCORE_BONUS)
		score += ROUND_SCORE_BONUS
	}

	if intT%25 == 0 {
		slog.Info("quarter score bonus", "bonus", QUARTER_SCORE_BONUS)
		score += QUARTER_SCORE_BONUS
		return score
	}

	if score != 0 {
		return score
	}

	return 0
}

// 5 points for every two items on the receipt.
func scoreTwoItemBonus(totalItems int) int {
	totalItems /= 2
	slog.Info("two item bonus",
		"total items", totalItems,
		"bonus", totalItems*ITEM_COUPLE_BONUS)
	return totalItems * ITEM_COUPLE_BONUS
}

// if the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer.
// The result is the number of points earned.
func scoreItemDescription(description string, price float64) int {
	trimmedString := strings.TrimSpace(description)
	if len(trimmedString)%3 == 0 {
		score := int(price*ITEM_DESCRIPTION_FACTOR) + 1
		slog.Info("item description bonus", "bonus", score)
		return score
	}
	return 0
}

// 6 points if the day in the purchase date is odd.
func scoreDayOfMonth(dom int) int {
	if dom%2 != 0 {
		slog.Info("odd day bonus", "bonus", ODD_DAY_SCORE)
		return ODD_DAY_SCORE
	}
	return 0

}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func scoreTimeOfDay(hour int, minute int) int {
	slog.Info("time of day",
		"hour", hour,
		"minute", minute)

	switch hour {
	case 14, 15:
		switch minute {
		case 0:
			return 0
		default:
			slog.Info("time of day bonus", "bonus", PURCHASE_TIME_BONUS)
			return PURCHASE_TIME_BONUS
		}
	default:
		return 0
	}
}
