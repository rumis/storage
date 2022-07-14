package skafka

import (
	"fmt"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	addr := "39.105.221.194:9093,39.105.204.21:9093,101.200.120.178:9093"
	username := "jiaoyanyun"
	password := "Chengjiukehu123"
	ca := `
-----BEGIN CERTIFICATE-----
MIIDPDCCAqWgAwIBAgIJAMRsb0DLM1fsMA0GCSqGSIb3DQEBBQUAMHIxCzAJBgNV
BAYTAkNOMQswCQYDVQQIEwJIWjELMAkGA1UEBxMCSFoxCzAJBgNVBAoTAkFCMRAw
DgYDVQQDEwdLYWZrYUNBMSowKAYJKoZIhvcNAQkBFht6aGVuZG9uZ2xpdS5semRA
YWxpYmFiYS5jb20wIBcNMTcwMzA5MTI1MDUyWhgPMjEwMTAyMTcxMjUwNTJaMHIx
CzAJBgNVBAYTAkNOMQswCQYDVQQIEwJIWjELMAkGA1UEBxMCSFoxCzAJBgNVBAoT
AkFCMRAwDgYDVQQDEwdLYWZrYUNBMSowKAYJKoZIhvcNAQkBFht6aGVuZG9uZ2xp
dS5semRAYWxpYmFiYS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALZV
bbIO1ULQQN853BTBgRfPiRJaAOWf38u8GC0TNp/E9qtI88A+79ywAP17k5WYJ7XS
wXMOJ3h1qkQT2TYJVetZ6E69CUJq4BsOvNlNRvmnW6eFymh5QZsEz2MTooxJjVjC
JQPlI2XRDjIrTVYEQWUDxj2JhB8VVPEed+6u4KQVAgMBAAGjgdcwgdQwHQYDVR0O
BBYEFHFlOoiqQxXanVi2GUoDiKDD33ujMIGkBgNVHSMEgZwwgZmAFHFlOoiqQxXa
nVi2GUoDiKDD33ujoXakdDByMQswCQYDVQQGEwJDTjELMAkGA1UECBMCSFoxCzAJ
BgNVBAcTAkhaMQswCQYDVQQKEwJBQjEQMA4GA1UEAxMHS2Fma2FDQTEqMCgGCSqG
SIb3DQEJARYbemhlbmRvbmdsaXUubHpkQGFsaWJhYmEuY29tggkAxGxvQMszV+ww
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQUFAAOBgQBTSz04p0AJXKl30sHw+UM/
/k1jGFJzI5p0Z6l2JzKQYPP3PfE/biE8/rmiGYEenNqWNy1ZSniEHwa8L/Ux98ci
4H0ZSpUrMo2+6bfuNW9X35CFPp5vYYJqftilJBKIJX3C3J1ruOuBR28UxE42xx4K
pQ70wChNi914c4B+SxkGUg==
-----END CERTIFICATE-----`

	InitDefaultDialer(D_WithCA(ca), D_WithUserNamePassword(username, password))
	r, c := NewReaderChannel(
		R_WithBrokers(strings.Split(addr, ",")),
		R_WithDialer(DefaultDialer()),
		R_WithGroupID("GID_fz_cservice_thumbnail_page_screenshot"),
		R_WithTopic("topic_fz_cservice_thumbnail_page_task"))

	defer c()
	msg := <-r
	fmt.Println(msg)
}
