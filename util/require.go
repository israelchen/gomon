package util

func Require(condition bool, panicReason string) {
	if condition == false {
		panic(panicReason)
	}
}
