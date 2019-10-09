package main

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIC/DCCAeSgAwIBAgIRALPaIawoYO+I3fv3BvMxOLAwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xOTEwMDkyMDExMzlaFw0yMDEwMDgyMDEx
MzlaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDuAzjcfyOm5G1+rpPyAqzPeAii7Jnk2rnszQUBESkLA8BgsLEN0RUm
51nSepoqo2JjD/FGeIfgkDfwLd10YKHgjmqC2zX46VPTmS7bhmUQJVHZZTnp2lzb
qUjtfDFvPE5vrM5BIY/8Ono7DIf5Yg5Dr0soj1eTp7IYcUIy3eZY1JG4hz/TaGwr
yAn6EsB4arslFmTrQD4jj6pxDWAcp3y2z31ezRUOzWIKXNRNmmfqiEEWJtMvFOSI
MiGF9BNVu4h0g7Itbf/9FSJKQfvlocAncX7K6kac6hWXtcYFG/DNGWlLKCm3dC84
GdTSZ4dEvsEyrKjyM8m40VjiV/uyj8+5AgMBAAGjTTBLMA4GA1UdDwEB/wQEAwIC
pDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MBMGA1UdEQQM
MAqCCHJlZGlyZWN0MA0GCSqGSIb3DQEBCwUAA4IBAQDfBW9j7gZeYGTRrxvdkTG+
ULOwBAKbKqjQUCGuQrpAX2qjkeS2ttegMl4k/bCBLyo7SVVjjIQE6PZ53MZppSZT
uvCb6EgW4wR2PpN4vkh4Yh34dmMNpzUC0V5FepKACX+DQR4oGaVecoAhmJAbTKGe
HFVORJ6XZBbY3m4AHd8qeg0lZzmdxWhhBjLCoDw9EUjSLZV+zl2FOOVfhLMSqr9E
sQK/iorhXFn2zJXAgukrWBWK4seibSjYovhulpt9x5hzeMTNoNHpRNeZOxp4Tukc
9dR3me9EZ7kyQPxJA2/fp68wMgIXo57Gw9e+30mwDfVUaXzY1HjfZ2IctNTgZXWk
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDuAzjcfyOm5G1+
rpPyAqzPeAii7Jnk2rnszQUBESkLA8BgsLEN0RUm51nSepoqo2JjD/FGeIfgkDfw
Ld10YKHgjmqC2zX46VPTmS7bhmUQJVHZZTnp2lzbqUjtfDFvPE5vrM5BIY/8Ono7
DIf5Yg5Dr0soj1eTp7IYcUIy3eZY1JG4hz/TaGwryAn6EsB4arslFmTrQD4jj6px
DWAcp3y2z31ezRUOzWIKXNRNmmfqiEEWJtMvFOSIMiGF9BNVu4h0g7Itbf/9FSJK
QfvlocAncX7K6kac6hWXtcYFG/DNGWlLKCm3dC84GdTSZ4dEvsEyrKjyM8m40Vji
V/uyj8+5AgMBAAECggEASFJ+UrHnWW0LwHS3y8/4RsqIhUkzshCsckISBTL7r7ci
G79U7Yfcz4d5CbXrZo1i9gsAG5PAZgIsnTSymAxM4/kicES/77SmniNr05TQ3Mka
R960bFTH5o9X86HLO6utgc2Wlr/mCpSSU6MJJkQfZX28bsSvrdRFD5xKqz42IkNP
Qe4k2ytBFAy3G4gWHFY56RHGN30shXgGP5My68kABPuywqfnrjocx1XD2+xYbqAX
Zf7U3wZgxtLPFxLSZeePgmMh2fsckQKZTE966SmotR51sEY94+AUZaZfLWIqUMCv
pblNkyL4nKyF1tIkaRlgPbmo4QlPLmn2MHle9bE7oQKBgQDyMLXZs8jF9nMa9rbS
AcQ2xgiCOePjFcL05gz5V5qwUcKbYrsxd1hwNaSCRSORw9cBmLDcSKOLsxZiiVR2
OpCtria2ZlPI1+Dr3+eUEr9VQLj3RZBdpwQlll9vZGHgDcWwbAxtWyUPQvVZJIwa
ZDlnztQ9vBQeio9rl19U9WMElQKBgQD7lYeFuwrigi7rJ0LrcyzZ/xtDoY59uquf
kMsmBJTwNbHyo/oaOZCE5wTYbRe6hojbPlajSpG+rqaVivvlSsIXc9L3V1oZXvXp
c7o9Is+tHyP2oeqh8kySD/TpFAwPbKEHVbmBOHy9FFCYiOlLDTv6z/w5z/+7venV
te7Pu89RlQKBgGrKs8UVE3jHHSZMl3yurriARgw2PphJZjfaoOnpiRoqUyd1N5mu
SF7iKHIQzohd1Jatn37iwMq+4yX77DRdyqHq4sMXB+bN2i3oAxM12Qxch7LxB6Fk
Hd39GoPhvY6wQ/VxD2HBCOxb2BfAl86jVvTBLLE0F6MH8gm9K5oowcqpAoGALGo/
nLpit45oHhe2Vr7koi/Jbm0tLMEx31++nZ2ddbLlEYMlek/DVdM7JcJMuB9cNeiR
fw6BIHrQ6gG5aseB8IYALq57N5NuMqK9tGFa7KNcxAPd2m1eW0L559QkNOzmmNbn
gwqn6vGVMPiqxxc1CZiCXOp9qXVjvNj7qizr8ukCgYEAkBgoKE+g5KKRkpnXw1y6
crNARUasUsIBRUuigjZPRu9GKiYoGPwWEXL5ldof/Y5isnoYlUaCGV4TuCFHdU1e
IjMgEuvnUSumJ6OQI4ylw8HZWMqDHd3aRiuSgJAfco2SzluZeaPNJaaAQy+PfrOx
C8bCW7tJTH9joFHtVZZ6q/M=
-----END PRIVATE KEY-----`)

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}
