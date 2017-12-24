package decimal_math

import (
	"fmt"
	"testing"
)

func TestDecimalMath(t *testing.T) {
	rt := Calculate_Percent_Change_Float(float64(19649.74000000), float64(19649.74000000))
	fmt.Println(rt)
}
