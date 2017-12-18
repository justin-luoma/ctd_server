package decimal_math

import (
	"github.com/shopspring/decimal"
)

func Calculate_Percent_Change(open float64, last float64) float64 {
	openDecimal := decimal.NewFromFloatWithExponent(open, -8)
	lastDecimal := decimal.NewFromFloatWithExponent(last, -8)
	rt := openDecimal.Neg().Add(lastDecimal).Div(openDecimal).Mul(decimal.New(100, 0)).Round(2)
	rtFloat, _ := rt.Float64()

	return rtFloat
}
