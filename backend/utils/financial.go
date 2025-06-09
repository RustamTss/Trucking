package utils

import (
	"math"
	"time"
)

// CalculateMonthlyPayment рассчитывает месячный платеж по аннуитетной схеме
func CalculateMonthlyPayment(principal float64, annualRate float64, termMonths int) float64 {
	if annualRate == 0 {
		return principal / float64(termMonths)
	}

	monthlyRate := annualRate / 100 / 12
	payment := principal * (monthlyRate * math.Pow(1+monthlyRate, float64(termMonths))) /
		(math.Pow(1+monthlyRate, float64(termMonths)) - 1)

	return math.Round(payment*100) / 100
}

// CalculateDepreciation рассчитывает амортизацию по прямолинейному методу
func CalculateDepreciation(purchasePrice float64, usefulLifeYears int, ageInYears float64) float64 {
	if usefulLifeYears <= 0 {
		return 0
	}

	annualDepreciation := purchasePrice / float64(usefulLifeYears)
	totalDepreciation := annualDepreciation * ageInYears

	if totalDepreciation > purchasePrice {
		return purchasePrice
	}

	return math.Round(totalDepreciation*100) / 100
}

// CalculateRemainingBalance рассчитывает остаток по кредиту после платежа
func CalculateRemainingBalance(currentBalance, principalPaid float64) float64 {
	remaining := currentBalance - principalPaid
	if remaining < 0 {
		return 0
	}
	return math.Round(remaining*100) / 100
}

// CalculateInterestPayment рассчитывает процентную часть платежа
func CalculateInterestPayment(balance float64, annualRate float64) float64 {
	monthlyRate := annualRate / 100 / 12
	return math.Round(balance*monthlyRate*100) / 100
}

// CalculateVehicleAge рассчитывает возраст транспорта в годах
func CalculateVehicleAge(purchaseDate time.Time) float64 {
	now := time.Now()
	duration := now.Sub(purchaseDate)
	years := duration.Hours() / 24 / 365.25
	return math.Round(years*100) / 100
}
