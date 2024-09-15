package shared

// MonthlyPaymentBreakdown represents the details of each monthly payment
type MonthlyPaymentBreakdown struct {
	Month            int
	InterestPayment  float64
	PrincipalPayment float64
	RemainingBalance float64
}
