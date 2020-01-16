package main

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
