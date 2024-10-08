// Модуль генерации рандомных строк
package keygen

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Длина строки при генерации
const KeyLength = 5

// Получение рандомного ключа/строки
func GetRandkey(n uint) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
