// Package utils
package utils

func Uint8ToString(input [32]uint8) string {
	var name []byte
	for _, b := range input {
		if b != 0 {
			name = append(name, b)
		}
	}

	return string(name)
}
