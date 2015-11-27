package util

import (

)

/**
 * KeyExists(map[string]string, string) bool
 */
func KeyExists(data map[string]string, key string) bool {
	_, ok := data[key]
	return ok
}
