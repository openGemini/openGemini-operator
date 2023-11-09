package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"strings"
)

func CalcMd5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbols      = "!@#$%*"
	digits       = "0123456789"
)

func GenPassword(length, numSymbols, numDigits, numUpperLetters int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < numSymbols; i++ {
		random := rand.Intn(len(symbols))
		password.WriteString(string(symbols[random]))
	}

	//Set number
	for i := 0; i < numDigits; i++ {
		random := rand.Intn(len(digits))
		password.WriteString(string(digits[random]))
	}

	//Set uppercase
	for i := 0; i < numUpperLetters; i++ {
		random := rand.Intn(len(upperLetters))
		password.WriteString(string(upperLetters[random]))
	}

	remainingLength := length - numSymbols - numDigits - numUpperLetters
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(lowerLetters))
		password.WriteString(string(lowerLetters[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func GetClusterDomain() string {
	return getEnvWithDefault("KUBERNETES_CLUSTER_DOMAIN", "cluster.local")
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
