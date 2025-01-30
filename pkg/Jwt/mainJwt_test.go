package Jwt

import (
	"testing"
)

func TestJWT_Create(t *testing.T) {
	const email = "test@test.com"

	jwtService := NewJWT("RxbxgRcFCFes0enila83XSdWzejBmKuw4cHiPuMgiU8")
	token, err := jwtService.Create(JwtDate{
		Email: "test@test.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("token is not valid")
	}
	if data.Email != email {
		t.Fatalf("data:%s is not equal to %s", data.Email, email)
	}
}
