package fixed

import (
	"fmt"
	"math"
	"time"

	shared "github.com/awazonek/mortgage-calculator/internal/shared"
)

// CalculateMortgageCost calculates the total cost of the mortgage
func CalculateMortgageCost(principal float64, annualInterestRate float64, amortizationPeriodYears int, mortgageDurationYears int, startDate string) (float64, []shared.MonthlyPaymentBreakdown, error) {
	// Parse the start date (optional for now, included for future use)
	_, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid start date format: %v", err)
	}

	// Convert annual interest rate to a monthly rate
	monthlyInterestRate := annualInterestRate / 12 / 100

	r := annualInterestRate / 12 / 100
	// Total number of payments (months)
	n := amortizationPeriodYears * 12

	fmt.Printf("\nPeriod: %d\n", n)
	// Monthly payment formula
	monthlyPayment := (r * principal) / (1 - math.Pow(1+r, -float64(n)))

	fmt.Printf("\nMonthly payment :%f\n", monthlyPayment)
	var breakdown []shared.MonthlyPaymentBreakdown
	remainingBalance := principal
	totalCost := 0.0
	// Iterate over the total mortgage duration in months
	for month := 1; month <= mortgageDurationYears*12; month++ {
		if remainingBalance <= 0 {
			break
		}

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
			Month:            month,
			InterestPayment:  interestPayment,
			PrincipalPayment: principalPayment,
			RemainingBalance: remainingBalance,
		})

		// Accumulate total cost
		totalCost += monthlyPayment
	}

	return totalCost, breakdown, nil
}
