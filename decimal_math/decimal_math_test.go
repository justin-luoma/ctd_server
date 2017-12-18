package decimal_math

import (
	"testing"
	"fmt"
)

func TestDecimalMath(t *testing.T) {
	rt :=Calculate_Percent_Change(float64(19649.74000000), float64(19649.74000000))
	fmt.Println(rt)
}