package variable

import (
	"fmt"
	"math"
	"time"

	shared "github.com/awazonek/mortgage-calculator/internal/shared"
)

// CalculateVariableMortgageCostWithMap calculates the total cost of the mortgage with variable rates using a date-to-rate map
func CalculateVariableMortgageCostWithMap(principal float64, dateToRate map[string]float64, startDate string, amortizationPeriodYears int, mortgageDurationYears int) (float64, []shared.MonthlyPaymentBreakdown, error) {
	// Parse the start date
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid start date format: %v", err)
	}

	// Monthly breakdown
	var breakdown []shared.MonthlyPaymentBreakdown
	remainingBalance := principal
	totalCost := 0.0

	startingRate, err := findStartingRate(dateToRate, startDate)
	if err != nil {
		return 0, nil, err
	}

	lastRate := startingRate
	// Iterate over each month of the mortgage duration
	for month := 0; month < mortgageDurationYears*12; month++ {
		// Determine the current date for the month
		currentDate := start.AddDate(0, month, 0).Format("2006-01-02T15:04:05-07:00")

		// Get the current month's variable rate from the map
		rate, exists := dateToRate[currentDate]
		if !exists {
			if lastRate < 0 {
				return totalCost, breakdown, fmt.Errorf("missing variable rate for date: %s", currentDate)
			}
			rate = lastRate
		}
		lastRate = rate

		// Convert the current annual interest rate to a monthly rate
		monthlyInterestRate := rate / 12 / 100

		// Calculate the monthly payment using the mortgage formula (assuming full amortization period payments)
		totalAmortizationPayments := amortizationPeriodYears * 12 // This would ideally be different if using actual amortization period
		monthlyPayment := principal * (monthlyInterestRate * math.Pow(1+monthlyInterestRate, float64(totalAmortizationPayments))) /
			(math.Pow(1+monthlyInterestRate, float64(totalAmortizationPayments)) - 1)

		// Calculate the monthly interest payment
		interestPayment := remainingBalance * monthlyInterestRate

		// Calculate the principal payment
		principalPayment := monthlyPayment - interestPayment

		// Reduce the remaining balance
		remainingBalance -= principalPayment

		// Ensure remaining balance does not go negative
		if remainingBalance < 0 {
			remainingBalance = 0
		}

		// Add to the breakdown
		breakdown = append(breakdown, shared.MonthlyPaymentBreakdown{
			Month:            month + 1,
			InterestPayment:  interestPayment,
			PrincipalPayment: principalPayment,
			RemainingBalance: math.Max(0, remainingBalance), // Ensure balance doesn't go negative
		})

		// Accumulate total cost
		totalCost += principalPayment + interestPayment
	}

	return totalCost, breakdown, nil
}

func findStartingRate(dateToRate map[string]float64, startDate string) (float64, error) {
	// Parse the start date
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format: %v", err)
	}

	// Check if the exact start date exists in the map
	if rate, exists := dateToRate[startDate]; exists {
		return rate, nil
	}

	// If no exact match, find the most recent date before the start date
	var latestRate float64
	var latestDate time.Time
	found := false

	for dateStr, rate := range dateToRate {
		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			continue // Skip invalid dates
		}

		// Find the most recent date before the start date
		if date.Before(start) && (latestDate.IsZero() || date.After(latestDate)) {
			latestDate = date
			latestRate = rate
			found = true
		}
	}

	if !found {
		return 0, fmt.Errorf("no valid rate found for the start date or before")
	}

	return latestRate, nil
}
