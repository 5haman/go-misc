package main

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIQdhJJfbhjS6sg9pxzHtMqnTANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTE5MTAxMDIyNDAzNFoXDTIwMTAwOTIyNDAz
NFowEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBANiN3/9DVgtLjt7vF6Jp0ktC7cpB4m5sl36DXQSenLpgU0nN8QH23erT
EHBjRnBTZqe/T0UYi7SiiYjKTeqYEcfGi5cxd3EdUdt1h3Vw4tAhPotif+YafSKo
TnAsEPhDufexk+uRmsZRiAArJKfk7e0FZYwaT7qGO99KcauTYp/BhvFNzIMn8L9m
M+oFwxXCWO7v+ur6gxAFaksLTRaQUjvYM8Cw7TuW1WTamAJN9hwxyfeLJzIggPED
WcSyP7Xaq18VHh6RfRpSSkZ/zC8AEAkQM6Tg4UJ6giIHg2yJEDJfLavW80+DpbC2
hqqSwvvuSQhLo4/ZivpDULzTQ3fkmJMCAwEAAaNNMEswDgYDVR0PAQH/BAQDAgKk
MBMGA1UdJQQMMAoGCCsGAQUFBwMBMA8GA1UdEwEB/wQFMAMBAf8wEwYDVR0RBAww
CoIIYWNtZS5jb20wDQYJKoZIhvcNAQELBQADggEBAFObh+MZAoWgWdU9w4c2yrQw
QbR4Eo9xCOJXDD35P+rK/F0HU/vNl7sIzRC5+fu6eSsFrFcMVOu6+6DHVwAd1Mr1
fPHb3Yy/ce0RWKQ3S6uFPafM+SCCJyjwLR5tPvvqm3HbtWEzs7dtFifLJfVs7sGz
D5IUBsXZM7d6HpmZmnmkXyrU21ikVoZ9QuWEd6iDk9FiYAVt6Zt8JmsQgdCD4XV4
kIL8dCbmwxLP9S2kUrNpuEFLdwayQKu5UXJBYjzBEGBYCvIilrVzWp3K1cuhymJT
RxkJRyzq6GHd73/ovyBGzInrRZDfIIVkZYVfYV85DKFUQlK4ZG/KcR4/ykZZ4sE=
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDYjd//Q1YLS47e
7xeiadJLQu3KQeJubJd+g10Enpy6YFNJzfEB9t3q0xBwY0ZwU2anv09FGIu0oomI
yk3qmBHHxouXMXdxHVHbdYd1cOLQIT6LYn/mGn0iqE5wLBD4Q7n3sZPrkZrGUYgA
KySn5O3tBWWMGk+6hjvfSnGrk2KfwYbxTcyDJ/C/ZjPqBcMVwlju7/rq+oMQBWpL
C00WkFI72DPAsO07ltVk2pgCTfYcMcn3iycyIIDxA1nEsj+12qtfFR4ekX0aUkpG
f8wvABAJEDOk4OFCeoIiB4NsiRAyXy2r1vNPg6WwtoaqksL77kkIS6OP2Yr6Q1C8
00N35JiTAgMBAAECggEAEweMmoLRSdbO7Do4anZAG4r7GF3nxupV+bETeHdzsFEM
oJyvRAvsflkjxayDoRVDHRSLo7e/dxUdXt7gL/BDB7ojxBp9s3vvGIjgaWqNE9sI
AvmZ4Z+MRYJiuiq1JwvUiLabGAVIg4rgl5sy8moEmmJyBPi+7tYi5sFE8d3WudD/
XtxXfFngT4hTNdqmZnzfO+ziva2pAeRsJ5OVFqbQOOSs8mawUBmoTLUXqRL8oBBM
48l0BqOwShDAC/ZtlpA8nXISWgd9bJazL7n99Ou1drwLvnVwP7GCBvzJDfaczDoy
Ylpf5FLZMwltp1xmMnjlPYchGI+s6tCZw8bzZ4VpAQKBgQDhynV388YUQ7TGUU22
ZFNXAK+zUXb3HIXpgyJsRkTC3OwJip43xhd8EavBAaW6C+pOUG9TrZDbcWjYSdz4
rvS6cVl3Z0ZGDoHGgHffSXG/tBeYdaZjMg2QzQG1l6CcxjVTaxjWM2AXy10bsLb8
EDUtELqO/W4EjCJLSGybuHXyEwKBgQD1hw1kCwg3XqM3AFfYAtud1VpBLGkRMbZE
Go7VPuWA5HYbZ74KSHf/HYMEEdTVYHlnpTn5HsdCYY1Xy+SLZQ5AVZZSB7b9ycwr
HfMMBGBJ6T2p4C8271HLezck4VMw7eAI5zoJE03wOqtmwy4T34oAsvXqIkcCNSmp
6TCBJG+PgQKBgFzFIZSibV1AIFNnbmWlPPS/THGB5D5N0tuJzKfuCyyBNt4IvU8v
LdEFNat8cMpLQP7iX4tjAeSX6TsMxiTLRbQhBGBh52a7aSjU+eudMoZQiW1T0YRq
OVaoVK522T/w1FIs66x+uVmtbdkFt3lDc4XLnMtJZ12o8iI6ZJ1qodNPAoGBANit
MpgTVFDo58jmOJ+dBgsn+dqCQsa1xFAdz+dI9mjlNYXB6+hPQ/aUKMcypU0ZMorR
OXQsQVTHmmDcwvhxWj5USbBito8Jw3BZod/9DKytdYmxGnm0gc69ElEtuKj5hDjX
NlREAQf7/r9ViBhpsfQj+vmA/oFoQTh9XhzZ9soBAoGAEfps+Fnejy6/DmKELx+g
aQukxuG97x45Hxy+RevzEYfQitqMZ/Lz79eTMdgyjW1cRtujLw0EjgWotahm0i2q
OnON96kHS4OgWPulWtszCXZ2HLrFditNcXWWZSiGu7VDgY3Blyf5qxYCHQDHJTvY
6H7I8Ci1GlY08Xkx09l0hH0=
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
