package keygen

import "fmt"

func GenAccessTokenKey(uid int64) string {
	return fmt.Sprintf("user:token:access:%v", uid)
}

func GenRefreshTokenKey(uid int64) string {
	return fmt.Sprintf("user:token:refresh:%v", uid)
}
