package main

import (
	"fmt"

	"github.com/awazonek/mortgage-calculator/internal/fixed"
	variable "github.com/awazonek/mortgage-calculator/internal/variable"
)

func main() {
	dt, err := variable.ReadData("data.json")
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	dateToRate, err := variable.MapDatesToRates(dt)
	if err != nil {
		fmt.Println("Error mapping dates to rates:", err)
		return
	}

	// Example of printing the mapped dates to rates
	for dateStr, rate := range dateToRate {
		fmt.Printf("Date: %s, Rate: %.2f%%\n", dateStr, rate)
	}

	// Example input values
	principal := 680000.0                    // Example principal amount
	annualInterestRate := 1.4                // Annual interest rate in percentage
	amortizationPeriodYears := 25            // Amortization period in years
	mortgageDurationYears := 5               // Total mortgage duration in years
	startDate := "2021-02-01T00:00:00-04:00" // Start date in RFC3339 format

	// Calculate the total mortgage cost
	totalCost, monthly, err := fixed.CalculateMortgageCost(principal, annualInterestRate, amortizationPeriodYears, mortgageDurationYears, startDate)
	if err != nil {
		fmt.Println("Error calculating mortgage cost:", err)
		return
	}

	fmt.Printf("The total cost of the fixed rate mortgage is: $%.2f\n", totalCost)

	fmt.Printf("Month | Interest Payment | Principal Payment | Remaining Balance\n")
	// Output the monthly breakdown for variable rates
	for _, b := range monthly {
		fmt.Printf("%d\t$%.2f\t\t\t$%.2f\t\t\t$%.2f\n", b.Month, b.InterestPayment, b.PrincipalPayment, b.RemainingBalance)
	}

	fmt.Printf("\nVariable:\n")

	// Calculate the total mortgage cost and breakdown for the variable rate using the date-to-rate map
	totalCostVariable, monthlyVariable, err := variable.CalculateVariableMortgageCostWithMap(principal, dateToRate, startDate, amortizationPeriodYears, mortgageDurationYears)
	if err != nil {
		fmt.Println("Error calculating variable mortgage cost:", err)
		return
	}

	fmt.Printf("The total cost of the variable mortgage is: $%.2f\n", totalCostVariable)

	fmt.Printf("Month | Interest Payment | Principal Payment | Remaining Balance\n")
	// Output the monthly breakdown for variable rates
	for _, b := range monthlyVariable {
		fmt.Printf("%d\t$%.2f\t\t\t$%.2f\t\t\t$%.2f\n", b.Month, b.InterestPayment, b.PrincipalPayment, b.RemainingBalance)
	}

}
