package main

import (
	"fmt"
	"math/big"
	"time"
)

func ConvertWeiToEther(w *big.Int) float64 {
	bigfloat := new(big.Float).SetInt(w)
	bigfloat = bigfloat.Mul(bigfloat, big.NewFloat(1.0e-18))

	f, _ := bigfloat.Float64()
	return f
}

func ConvertEtherToWei(e float64) *big.Int {
	tf := new(big.Float).SetFloat64(e)
	tf = tf.Mul(tf, big.NewFloat(1.0e+18))

	f, _ := tf.Float64()
	s := fmt.Sprintf("%.0f", f)

	ib, _ := new(big.Int).SetString(s, 10)

	return ib
}

func main() {
	bigint := ConvertEtherToWei(10000000.00000001)
	fmt.Printf("%s\n", bigint.Uint64())

	time.Sleep(3 * time.Second)
}
