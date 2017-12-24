package decimal_math

import (
	"github.com/shopspring/decimal"
	"github.com/golang/glog"
)

func Calculate_Percent_Change_Float(open float64, last float64) float64 {
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

func Calculate_Percent_Change_Decimal(open decimal.Decimal, last decimal.Decimal) (float64) {
	diff := open.Neg().Add(last)
	if checkdiff, _ := diff.Float64(); checkdiff == 0 {
		return float64(0)
	}
	change := diff.Div(open).Mul(decimal.New(100, 0)).Round(2)
	rtFloat, _ := change.Float64()

	return rtFloat
}

func Convert_To_Float64(decimalValue decimal.Decimal) (float64) {
	rtFloat, exact := decimalValue.Float64()
	if exact {
		glog.Warningln("Conversion to float not exact")
	}

	return rtFloat
}

