package utils

import "math/rand"

// GenerateOtp generates a random number of length len
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

// GenerateRandomString generates a random string of length len
// using the pool of characters provided
func GenerateRandomString(len int) string {
	pool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return generate(len, pool)
}
