package utils

import (
	"gonum.org/v1/gonum/floats"
	"github.com/markcheno/go-talib"
)

func KDJ_CN(high, low, close []float64, fastk_period, slowk_period, fastd_period int) ([]float64, []float64, []float64) {
	var periodLowArr, periodHighArr []float64
	length := len(high)
	var rsv []float64 = make([]float64, length)
	var k []float64 = make([]float64, length)
	var d []float64 = make([]float64, length)
	var j []float64 = make([]float64, length)
	for i := 0; i < length; i++ {
		periodLowArr = append(periodLowArr, low[i])
		periodHighArr = append(periodHighArr, high[i])

		if fastk_period == len(periodLowArr) {
			lowest := floats.Min(periodLowArr)
			highest := floats.Max(periodHighArr)
			if highest-lowest < 0.000001 {
				rsv[i] = 100
			} else {
				rsv[i] = 100 * (close[i] - lowest) / (highest - lowest)
			}
			k[i] = (2.0/float64(slowk_period))*k[i-1] + 1.0/float64(slowk_period)*rsv[i]
			d[i] = (2.0/float64(fastd_period))*d[i-1] + 1.0/float64(fastd_period)*k[i]
			j[i] = 3*k[i] - 2*d[i]
			periodLowArr = periodLowArr[1:]
			periodHighArr = periodHighArr[1:]
		} else {
			k[i] = 50
			d[i] = 50
			rsv[i] = 0
			j[i] = 3*k[i] - 2*d[i]
		}
	}
	return k, d, j
}

func SMA_CN(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var smLine []float64 = make([]float64, length)
	for i := 0; i < length; i++ {
		periodArr = append(periodArr, Price[i])
		if periods == len(periodArr) {
			smLine[i] = floats.Sum(periodArr) / (float64)(len(periodArr))
			periodArr = periodArr[1:]
		} else {
			smLine[i] = 0
		}
	}

	return smLine
}

func MACD_CN(inReal []float64, inFastPeriod int, inSlowPeriod int, inSignalPeriod int) ([]float64, []float64, []float64) {
	macdDIFF, macdDEA, m := talib.MacdExt(inReal, inFastPeriod, 1, inSlowPeriod, 1, inSignalPeriod, 1)
	length := len(m)
	var macd []float64 = make([]float64, length)
	for i := 0; i < length; i++ {
		macd[i] = 2 * m[i]
	}
	return macdDIFF, macdDEA, macd
}
