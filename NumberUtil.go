package util

import (
	"math"
	"strconv"
)

func KeepDecimal(f float64, i int) float64 {
	m := math.Pow10(i)
	x := int(f*m + 0.5)
	return float64(x) / m
}

func Round(f float64, i int) float64 {
	n := math.Pow10(i)
	return math.Trunc((f+0.5/n)*n) / n
}

func FormatFloat(f float64, decimal int) string {
	return strconv.FormatFloat(f, 'g', decimal+1, 64)
}
