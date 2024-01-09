package utils

import "math/rand"

func GenerateOtp(len int) string {
	pool := "0123456789"

	return generate(len, pool)
}

func generate(length int, pool string) string {
	var otp string

	for i := 0; i < length; i++ {
		otp += string(pool[rand.Intn(len(pool))])
	}

	return otp
}
