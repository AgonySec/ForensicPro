package utils

// Decrypt function to decrypt the data using the given key
func Decrypt(key []byte, data []byte) []byte {
	return encryptOutput(key, data)
}

// EncryptInitalize initializes the RC4 state with the key
func encryptInitalize(key []byte) []byte {
	s := make([]byte, 256)
	for i := 0; i < 256; i++ {
		s[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) & 0xFF
		swap(s, i, j)
	}
	return s
}

// EncryptOutput performs the RC4 encryption or decryption on the data
func encryptOutput(key []byte, data []byte) []byte {
	s := encryptInitalize(key)
	i, j := 0, 0
	result := make([]byte, len(data))

	for k := 0; k < len(data); k++ {
		i = (i + 1) & 0xFF
		j = (j + int(s[i])) & 0xFF
		swap(s, i, j)
		result[k] = data[k] ^ s[(int(s[i])+int(s[j]))&0xFF]
	}
	return result
}

// Swap swaps two elements in the slice
func swap(s []byte, i, j int) {
	s[i], s[j] = s[j], s[i]
}
