package lib

var charMap = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
	'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '-', '_',
}

func b64encode(i int) string {
	tmp := make([]byte, 0, 11)
	first := true
	for i > 0 || first {
		low := i & 0x3f
		tmp = append(tmp, charMap[low])
		i = i >> 6
		first = false
	}
	return string(tmp)
}
