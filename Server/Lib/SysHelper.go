package Lib

import (
	"os"
)

func GetEnv(Key string) string {
	return os.Getenv(Key)
}

// func SetEnv(Key string, Value string) error {
// 	return os.Setenv(Key, Value)
// }

// func UnsetEnv(Key string) error {
// 	return os.Unsetenv(Key)
// }
