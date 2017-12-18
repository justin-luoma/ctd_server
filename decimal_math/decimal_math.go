package decimal_math

import (
	"github.com/shopspring/decimal"
)

func Calculate_Percent_Change(open float64, last float64) float64 {
	openDecimal := decimal.NewFromFloatWithExponent(open, -8)
	lastDecimal := decimal.NewFromFloatWithExponent(last, -8)
	diff := openDecimal.Neg().Add(lastDecimal)
	if checkdiff, _ := diff.Float64(); checkdiff == 0 {
		return float64(0)
	}
	change := diff.Div(openDecimal).Mul(decimal.New(100, 0)).Round(2)
	rtFloat, _ := change.Float64()

	return rtFloat
}
