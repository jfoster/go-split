package split

// benchmark these
func countDigits(number uint64) uint64 {
	var count uint64 = 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count
}

func countDigitsR(number uint64) uint64 {
	if number < 10 {
		return 1
	} else {
		return 1 + countDigitsR(number/10)
	}
}
