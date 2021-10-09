package conf

import "os"

const (
	QWEATHER_PUBLIC_ID = "QWEATHER_PUBLIC_ID"
	QWEATHER_KEY       = "QWEATHER_KEY"
)

func GetPublicId() string {
	publicId := os.Getenv(QWEATHER_PUBLIC_ID)
	return publicId
}

func GetKey() string {
	key := os.Getenv(QWEATHER_KEY)
	return key
}
