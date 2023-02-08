package middle

import "testing"

func TestJwt(t *testing.T) {
	Permission()
	ShaMiddleWare()
	NoAuthToGetUserId()
}
