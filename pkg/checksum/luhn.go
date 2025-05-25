package checksum

func VerifyLuhn(v string) bool {
	s := []rune(v)
	l := len(v)
	parity := l % 2
	sum := 0

	for i, v := range s {
		d := int(v - '0')
		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}

	return sum%10 == 0
}
