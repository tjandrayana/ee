package secure

func Encrypt(input string) string {
	result := []byte{}

	// ini optional cuma supaya hasil transform bisa dibaca manusia
	modifier := byte(100)

	for i := 1; i < len(input); i++ {
		result = append(result, input[i-1]-input[i]+modifier)
	}
	return string(result)
}
