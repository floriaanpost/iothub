package common

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// DigiCert Global G2 Root
// DigiCert Baltimore Root (sha1 fingerprint=d4de20d05e66fc53fe1a50882c78db2852cae474)
// Microsoft RSA TLS CA 01 (sha1 fingerprint=703d7a8f0ebf55aaa59f98eaf4a206004eb2516a)
// Microsoft RSA TLS CA 02 (sha1 fingerprint=b0c2d2d13cdd56cdaa6ab6e2c04440be4a429c75)
// Microsoft Azure TLS Issuing CA 01 (sha1 fingerprint=2f2877c5d778c31e0f29c7e371df5471bd673173)
// Microsoft Azure TLS Issuing CA 02 (sha1 fingerprint=e7eea674ca718e3befd90858e09f8372ad0ae2aa)
// Microsoft Azure TLS Issuing CA 05 (sha1 fingerprint=6c3af02e7f269aa73afd0eff2a88a4a1f04ed1e5)
// Microsoft Azure TLS Issuing CA 06 (sha1 fingerprint=30e01761ab97e59a06b41ef20af6f2de7ef4f7b0)
var caCerts = []byte(`-----BEGIN CERTIFICATE-----
MIIDjjCCAnagAwIBAgIQAzrx5qcRqaC7KGSxHQn65TANBgkqhkiG9w0BAQsFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0xMzA4MDExMjAwMDBaFw0zODAxMTUxMjAwMDBaMGExCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j
b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IEcyMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuzfNNNx7a8myaJCtSnX/RrohCgiN9RlUyfuI
2/Ou8jqJkTx65qsGGmvPrC3oXgkkRLpimn7Wo6h+4FR1IAWsULecYxpsMNzaHxmx
1x7e/dfgy5SDN67sH0NO3Xss0r0upS/kqbitOtSZpLYl6ZtrAGCSYP9PIUkY92eQ
q2EGnI/yuum06ZIya7XzV+hdG82MHauVBJVJ8zUtluNJbd134/tJS7SsVQepj5Wz
tCO7TG1F8PapspUwtP1MVYwnSlcUfIKdzXOS0xZKBgyMUNGPHgm+F6HmIcr9g+UQ
vIOlCsRnKPZzFBQ9RnbDhxSJITRNrw9FDKZJobq7nMWxM4MphQIDAQABo0IwQDAP
BgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNVHQ4EFgQUTiJUIBiV
5uNu5g/6+rkS7QYXjzkwDQYJKoZIhvcNAQELBQADggEBAGBnKJRvDkhj6zHd6mcY
1Yl9PMWLSn/pvtsrF9+wX3N3KjITOYFnQoQj8kVnNeyIv/iPsGEMNKSuIEyExtv4
NeF22d+mQrvHRAiGfzZ0JFrabA0UWTW98kndth/Jsw1HKj2ZL7tcu7XUIOGZX1NG
Fdtom/DzMNU+MeKNhJ7jitralj41E6Vf8PlwUHBHQRFXGU7Aj64GxJUTFy8bJZ91
8rGOmaFvE7FBcf6IKshPECBV1/MUReXgRPTqh5Uykw7+U0b6LJ3/iyK5S9kJRaTe
pLiaWN0bfVKfjllDiIGknibVb63dDcY3fe0Dkhvld1927jyNxF1WW6LZZm6zNTfl
MrY=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDdzCCAl+gAwIBAgIEAgAAuTANBgkqhkiG9w0BAQUFADBaMQswCQYDVQQGEwJJ
RTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJlclRydXN0MSIwIAYD
VQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTAwMDUxMjE4NDYwMFoX
DTI1MDUxMjIzNTkwMFowWjELMAkGA1UEBhMCSUUxEjAQBgNVBAoTCUJhbHRpbW9y
ZTETMBEGA1UECxMKQ3liZXJUcnVzdDEiMCAGA1UEAxMZQmFsdGltb3JlIEN5YmVy
VHJ1c3QgUm9vdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKMEuyKr
mD1X6CZymrV51Cni4eiVgLGw41uOKymaZN+hXe2wCQVt2yguzmKiYv60iNoS6zjr
IZ3AQSsBUnuId9Mcj8e6uYi1agnnc+gRQKfRzMpijS3ljwumUNKoUMMo6vWrJYeK
mpYcqWe4PwzV9/lSEy/CG9VwcPCPwBLKBsua4dnKM3p31vjsufFoREJIE9LAwqSu
XmD+tqYF/LTdB1kC1FkYmGP1pWPgkAx9XbIGevOF6uvUA65ehD5f/xXtabz5OTZy
dc93Uk3zyZAsuT3lySNTPx8kmCFcB5kpvcY67Oduhjprl3RjM71oGDHweI12v/ye
jl0qhqdNkNwnGjkCAwEAAaNFMEMwHQYDVR0OBBYEFOWdWTCCR1jMrPoIVDaGezq1
BE3wMBIGA1UdEwEB/wQIMAYBAf8CAQMwDgYDVR0PAQH/BAQDAgEGMA0GCSqGSIb3
DQEBBQUAA4IBAQCFDF2O5G9RaEIFoN27TyclhAO992T9Ldcw46QQF+vaKSm2eT92
9hkTI7gQCvlYpNRhcL0EYWoSihfVCr3FvDB81ukMJY2GQE/szKN+OMY3EU/t3Wgx
jkzSswF07r51XgdIGn9w/xZchMB5hbgF/X++ZRGjD8ACtPhSNzkE1akxehi/oCr0
Epn3o0WC4zxe9Z2etciefC7IpJ5OCBRLbf1wbWsaY71k5h+3zvDyny67G7fyUIhz
ksLi4xaNmjICq44Y3ekQEe5+NauQrz4wlHrQMz2nZQ/1/I6eYs9HRCwBXbsdtTLS
R9I4LtD+gdwyah617jzV/OeBHRnDJELqYzmp
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFWjCCBEKgAwIBAgIQDxSWXyAgaZlP1ceseIlB4jANBgkqhkiG9w0BAQsFADBa
MQswCQYDVQQGEwJJRTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJl
clRydXN0MSIwIAYDVQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTIw
MDcyMTIzMDAwMFoXDTI0MTAwODA3MDAwMFowTzELMAkGA1UEBhMCVVMxHjAcBgNV
BAoTFU1pY3Jvc29mdCBDb3Jwb3JhdGlvbjEgMB4GA1UEAxMXTWljcm9zb2Z0IFJT
QSBUTFMgQ0EgMDEwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCqYnfP
mmOyBoTzkDb0mfMUUavqlQo7Rgb9EUEf/lsGWMk4bgj8T0RIzTqk970eouKVuL5R
IMW/snBjXXgMQ8ApzWRJCZbar879BV8rKpHoAW4uGJssnNABf2n17j9TiFy6BWy+
IhVnFILyLNK+W2M3zK9gheiWa2uACKhuvgCca5Vw/OQYErEdG7LBEzFnMzTmJcli
W1iCdXby/vI/OxbfqkKD4zJtm45DJvC9Dh+hpzqvLMiK5uo/+aXSJY+SqhoIEpz+
rErHw+uAlKuHFtEjSeeku8eR3+Z5ND9BSqc6JtLqb0bjOHPm5dSRrgt4nnil75bj
c9j3lWXpBb9PXP9Sp/nPCK+nTQmZwHGjUnqlO9ebAVQD47ZisFonnDAmjrZNVqEX
F3p7laEHrFMxttYuD81BdOzxAbL9Rb/8MeFGQjE2Qx65qgVfhH+RsYuuD9dUw/3w
ZAhq05yO6nk07AM9c+AbNtRoEcdZcLCHfMDcbkXKNs5DJncCqXAN6LhXVERCw/us
G2MmCMLSIx9/kwt8bwhUmitOXc6fpT7SmFvRAtvxg84wUkg4Y/Gx++0j0z6StSeN
0EJz150jaHG6WV4HUqaWTb98Tm90IgXAU4AW2GBOlzFPiU5IY9jt+eXC2Q6yC/Zp
TL1LAcnL3Qa/OgLrHN0wiw1KFGD51WRPQ0Sh7QIDAQABo4IBJTCCASEwHQYDVR0O
BBYEFLV2DDARzseSQk1Mx1wsyKkM6AtkMB8GA1UdIwQYMBaAFOWdWTCCR1jMrPoI
VDaGezq1BE3wMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYI
KwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYI
KwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTA6BgNVHR8EMzAxMC+g
LaArhilodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vT21uaXJvb3QyMDI1LmNybDAq
BgNVHSAEIzAhMAgGBmeBDAECATAIBgZngQwBAgIwCwYJKwYBBAGCNyoBMA0GCSqG
SIb3DQEBCwUAA4IBAQCfK76SZ1vae4qt6P+dTQUO7bYNFUHR5hXcA2D59CJWnEj5
na7aKzyowKvQupW4yMH9fGNxtsh6iJswRqOOfZYC4/giBO/gNsBvwr8uDW7t1nYo
DYGHPpvnpxCM2mYfQFHq576/TmeYu1RZY29C4w8xYBlkAA8mDJfRhMCmehk7cN5F
JtyWRj2cZj/hOoI45TYDBChXpOlLZKIYiG1giY16vhCRi6zmPzEwv+tk156N6cGS
Vm44jTQ/rs1sa0JSYjzUaYngoFdZC4OfxnIkQvUIA4TOFmPzNPEFdjcZsgbeEz4T
cGHTBPK4R28F44qIMCtHRV55VMX53ev6P3hRddJb
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFWjCCBEKgAwIBAgIQD6dHIsU9iMgPWJ77H51KOjANBgkqhkiG9w0BAQsFADBa
MQswCQYDVQQGEwJJRTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJl
clRydXN0MSIwIAYDVQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTIw
MDcyMTIzMDAwMFoXDTI0MTAwODA3MDAwMFowTzELMAkGA1UEBhMCVVMxHjAcBgNV
BAoTFU1pY3Jvc29mdCBDb3Jwb3JhdGlvbjEgMB4GA1UEAxMXTWljcm9zb2Z0IFJT
QSBUTFMgQ0EgMDIwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQD0wBlZ
qiokfAYhMdHuEvWBapTj9tFKL+NdsS4pFDi8zJVdKQfR+F039CDXtD9YOnqS7o88
+isKcgOeQNTri472mPnn8N3vPCX0bDOEVk+nkZNIBA3zApvGGg/40Thv78kAlxib
MipsKahdbuoHByOB4ZlYotcBhf/ObUf65kCRfXMRQqOKWkZLkilPPn3zkYM5GHxe
I4MNZ1SoKBEoHa2E/uDwBQVxadY4SRZWFxMd7ARyI4Cz1ik4N2Z6ALD3MfjAgEED
woknyw9TGvr4PubAZdqU511zNLBoavar2OAVTl0Tddj+RAhbnX1/zypqk+ifv+d3
CgiDa8Mbvo1u2Q8nuUBrKVUmR6EjkV/dDrIsUaU643v/Wp/uE7xLDdhC5rplK9si
NlYohMTMKLAkjxVeWBWbQj7REickISpc+yowi3yUrO5lCgNAKrCNYw+wAfAvhFkO
eqPm6kP41IHVXVtGNC/UogcdiKUiR/N59IfYB+o2v54GMW+ubSC3BohLFbho/oZZ
5XyulIZK75pwTHmauCIeE5clU9ivpLwPTx9b0Vno9+ApElrFgdY0/YKZ46GfjOC9
ta4G25VJ1WKsMmWLtzyrfgwbYopquZd724fFdpvsxfIvMG5m3VFkThOqzsOttDcU
fyMTqM2pan4txG58uxNJ0MjR03UCEULRU+qMnwIDAQABo4IBJTCCASEwHQYDVR0O
BBYEFP8vf+EG9DjzLe0ljZjC/g72bPz6MB8GA1UdIwQYMBaAFOWdWTCCR1jMrPoI
VDaGezq1BE3wMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYI
KwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYI
KwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTA6BgNVHR8EMzAxMC+g
LaArhilodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vT21uaXJvb3QyMDI1LmNybDAq
BgNVHSAEIzAhMAgGBmeBDAECATAIBgZngQwBAgIwCwYJKwYBBAGCNyoBMA0GCSqG
SIb3DQEBCwUAA4IBAQCg2d165dQ1tHS0IN83uOi4S5heLhsx+zXIOwtxnvwCWdOJ
3wFLQaFDcgaMtN79UjMIFVIUedDZBsvalKnx+6l2tM/VH4YAyNPx+u1LFR0joPYp
QYLbNYkedkNuhRmEBesPqj4aDz68ZDI6fJ92sj2q18QvJUJ5Qz728AvtFOat+Ajg
K0PFqPYEAviUKr162NB1XZJxf6uyIjUlnG4UEdHfUqdhl0R84mMtrYINksTzQ2sH
YM8fEhqICtTlcRLr/FErUaPUe9648nziSnA0qKH7rUZqP/Ifmbo+WNZSZG1BbgOh
lk+521W+Ncih3HRbvRBE0LWYT8vWKnfjgZKxwHwJ
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQCq+mxcpjxFFB6jvh98dTFzANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwMTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAMedcDrkXufP7pxVm1FHLDNA9IjwHaMoaY8arqqZ4Gff4xyr
RygnavXL7g12MPAx8Q6Dd9hfBzrfWxkF0Br2wIvlvkzW01naNVSkHp+OS3hL3W6n
l/jYvZnVeJXjtsKYcXIf/6WtspcF5awlQ9LZJcjwaH7KoZuK+THpXCMtzD8XNVdm
GW/JI0C/7U/E7evXn9XDio8SYkGSM63aLO5BtLCv092+1d4GGBSQYolRq+7Pd1kR
EkWBPm0ywZ2Vb8GIS5DLrjelEkBnKCyy3B0yQud9dpVsiUeE7F5sY8Me96WVxQcb
OyYdEY/j/9UpDlOG+vA+YgOvBhkKEjiqygVpP8EZoMMijephzg43b5Qi9r5UrvYo
o19oR/8pf4HJNDPF0/FJwFVMW8PmCBLGstin3NE1+NeWTkGt0TzpHjgKyfaDP2tO
4bCk1G7pP2kDFT7SYfc8xbgCkFQ2UCEXsaH/f5YmpLn4YPiNFCeeIida7xnfTvc4
7IxyVccHHq1FzGygOqemrxEETKh8hvDR6eBdrBwmCHVgZrnAqnn93JtGyPLi6+cj
WGVGtMZHwzVvX1HvSFG771sskcEjJxiQNQDQRWHEh3NxvNb7kFlAXnVdRkkvhjpR
GchFhTAzqmwltdWhWDEyCMKC2x/mSZvZtlZGY+g37Y72qHzidwtyW7rBetZJAgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUDyBd16FXlduSzyvQx8J3BM5ygHYwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAlFvNh7QgXVLAZSsNR2XRmIn9iS8OHFCBA
WxKJoi8YYQafpMTkMqeuzoL3HWb1pYEipsDkhiMnrpfeYZEA7Lz7yqEEtfgHcEBs
K9KcStQGGZRfmWU07hPXHnFz+5gTXqzCE2PBMlRgVUYJiA25mJPXfB00gDvGhtYa
+mENwM9Bq1B9YYLyLjRtUz8cyGsdyTIG/bBM/Q9jcV8JGqMU/UjAdh1pFyTnnHEl
Y59Npi7F87ZqYYJEHJM2LGD+le8VsHjgeWX2CJQko7klXvcizuZvUEDTjHaQcs2J
+kPgfyMIOY1DMJ21NxOJ2xPRC/wAh/hzSBRVtoAnyuxtkZ4VjIOh
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQDGrpfM7VmYOGkKAKnqUyFDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwMjCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAOBiO1K6Fk4fHI6t3mJkpg7lxoeUgL8tz9wuI2z0UgY8vFra
3VBo7QznC4K3s9jqKWEyIQY11Le0108bSYa/TK0aioO6itpGiigEG+vH/iqtQXPS
u6D804ri0NFZ1SOP9IzjYuQiK6AWntCqP4WAcZAPtpNrNLPBIyiqmiTDS4dlFg1d
skMuVpT4z0MpgEMmxQnrSZ615rBQ25vnVbBNig04FCsh1V3S8ve5Gzh08oIrL/g5
xq95oRrgEeOBIeiegQpoKrLYyo3R1Tt48HmSJCBYQ52Qc34RgxQdZsLXMUrWuL1J
LAZP6yeo47ySSxKCjhq5/AUWvQBP3N/cP/iJzKKKw23qJ/kkVrE0DSVDiIiXWF0c
9abSGhYl9SPl86IHcIAIzwelJ4SKpHrVbh0/w4YHdFi5QbdAp7O5KxfxBYhQOeHy
is01zkpYn6SqUFGvbK8eZ8y9Aclt8PIUftMG6q5BhdlBZkDDV3n70RlXwYvllzfZ
/nV94l+hYp+GLW7jSmpxZLG/XEz4OXtTtWwLV+IkIOe/EDF79KCazW2SXOIvVInP
oi1PqN4TudNv0GyBF5tRC/aBjUqply1YYfeKwgRVs83z5kuiOicmdGZKH9SqU5bn
Kse7IlyfZLg6yAxYyTNe7A9acJ3/pGmCIkJ/9dfLUFc4hYb3YyIIYGmqm2/3AgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUAKuR/CFiJpeaqHkbYUGQYKliZ/0wHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAzo/KdmWPPTaYLQW7J5DqxEiBT9QyYGUfe
Zd7TR1837H6DSkFa/mGM1kLwi5y9miZKA9k6T9OwTx8CflcvbNO2UkFW0VCldEGH
iyx5421+HpRxMQIRjligePtOtRGXwaNOQ7ySWfJhRhKcPKe2PGFHQI7/3n+T3kXQ
/SLu2lk9Qs5YgSJ3VhxBUznYn1KVKJWPE07M55kuUgCquAV0PksZj7EC4nK6e/UV
bPumlj1nyjlxhvNud4WYmr4ntbBev6cSbK78dpI/3cr7P/WJPYJuL0EsO3MgjS3e
DCX7NXp5ylue3TcpQfRU8BL+yZC1wqX98R4ndw7X4qfGaE7SlF7I
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQDXvt6X2CCZZ6UmMbi90YvTANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwNTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAKplDTmQ9afwVPQelDuu+NkxNJ084CNKnrZ21ABewE+UU4GK
DnwygZdK6agNSMs5UochUEDzz9CpdV5tdPzL14O/GeE2gO5/aUFTUMG9c6neyxk5
tq1WdKsPkitPws6V8MWa5d1L/y4RFhZHUsgxxUySlYlGpNcHhhsyr7EvFecZGA1M
fsitAWVp6hiWANkWKINfRcdt3Z2A23hmMH9MRSGBccHiPuzwrVsSmLwvt3WlRDgO
bJkE40tFYvJ6GXAQiaGHCIWSVObgO3zj6xkdbEFMmJ/zr2Wet5KEcUDtUBhA4dUU
oaPVz69u46V56Vscy3lXu1Ylsk84j5lUPLdsAxtultP4OPQoOTpnY8kxWkH6kgO5
gTKE3HRvoVIjU4xJ0JQ746zy/8GdQA36SaNiz4U3u10zFZg2Rkv2dL1Lv58EXL02
r5q5B/nhVH/M1joTvpRvaeEpAJhkIA9NkpvbGEpSdcA0OrtOOeGtrsiOyMBYkjpB
5nw0cJY1QHOr3nIvJ2OnY+OKJbDSrhFqWsk8/1q6Z1WNvONz7te1pAtHerdPi5pC
HeiXCNpv+fadwP0k8czaf2Vs19nYsgWn5uIyLQL8EehdBzCbOKJy9sl86S4Fqe4H
GyAtmqGlaWOsq2A6O/paMi3BSmWTDbgPLCPBbPte/bsuAEF4ajkPEES3GHP9AgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUx7KcfxzjuFrv6WgaqF2UwSZSamgwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAe+G+G2RFdWtYxLIKMR5H/aVNFjNP7Jdeu
+oZaKaIu7U3NidykFr994jSxMBMV768ukJ5/hLSKsuj/SLjmAfwRAZ+w0RGqi/kO
vPYUlBr/sKOwr3tVkg9ccZBebnBVG+DLKTp2Ox0+jYBCPxla5FO252qpk7/6wt8S
Zk3diSU12Jm7if/jjkhkGB/e8UdfrKoLytDvqVeiwPA5FPzqKoSqN75byLjsIKJE
dNi07SY45hN/RUnsmIoAf93qlaHR/SJWVRhrWt3JmeoBJ2RDK492zF6TGu1moh4a
E6e00YkwTPWreuwvaLB220vWmtgZPs+DSIb2d9hPBdCJgvcho1c7
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQAueRcfuAIek/4tmDg0xQwDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwNjCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBALVGARl56bx3KBUSGuPc4H5uoNFkFH4e7pvTCxRi4j/+z+Xb
wjEz+5CipDOqjx9/jWjskL5dk7PaQkzItidsAAnDCW1leZBOIi68Lff1bjTeZgMY
iwdRd3Y39b/lcGpiuP2d23W95YHkMMT8IlWosYIX0f4kYb62rphyfnAjYb/4Od99
ThnhlAxGtfvSbXcBVIKCYfZgqRvV+5lReUnd1aNjRYVzPOoifgSx2fRyy1+pO1Uz
aMMNnIOE71bVYW0A1hr19w7kOb0KkJXoALTDDj1ukUEDqQuBfBxReL5mXiu1O7WG
0vltg0VZ/SZzctBsdBlx1BkmWYBW261KZgBivrql5ELTKKd8qgtHcLQA5fl6JB0Q
gs5XDaWehN86Gps5JW8ArjGtjcWAIP+X8CQaWfaCnuRm6Bk/03PQWhgdi84qwA0s
sRfFJwHUPTNSnE8EiGVk2frt0u8PG1pwSQsFuNJfcYIHEv1vOzP7uEOuDydsmCjh
lxuoK2n5/2aVR3BMTu+p4+gl8alXoBycyLmj3J/PUgqD8SL5fTCUegGsdia/Sa60
N2oV7vQ17wjMN+LXa2rjj/b4ZlZgXVojDmAjDwIRdDUujQu0RVsJqFLMzSIHpp2C
Zp7mIoLrySay2YYBu7SiNwL95X6He2kS8eefBBHjzwW/9FxGqry57i71c2cDAgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQU1cFnOsKjnfR3UltZEjgp5lVou6UwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQB2oWc93fB8esci/8esixj++N22meiGDjgF
+rA2LUK5IOQOgcUSTGKSqF9lYfAxPjrqPjDCUPHCURv+26ad5P/BYtXtbmtxJWu+
cS5BhMDPPeG3oPZwXRHBJFAkY4O4AF7RIAAUW6EzDflUoDHKv83zOiPfYGcpHc9s
kxAInCedk7QSgXvMARjjOqdakor21DTmNIUotxo8kHv5hwRlGhBJwps6fEVi1Bt0
trpM/3wYxlr473WSPUFZPgP1j519kLpWOJ8z09wxay+Br29irPcBYv0GMXlHqThy
8y4m/HyTQeI2IMvMrQnwqPpY+rLIXyviI2vLoI+4xKE4Rn38ZZ8m
-----END CERTIFICATE-----
`)

// RootCAs root CA certificates pool for connecting to the cloud.
func RootCAs() *x509.CertPool {
	p := x509.NewCertPool()
	if ok := p.AppendCertsFromPEM(caCerts); !ok {
		panic("tls: unable to append certificates")
	}
	return p
}

// TrustBundleResponse aids parsing the response from the edge.
type TrustBundleResponse struct {
	Certificate string `json:"certificate"`
}

// TrustBundle root CA certificates pool for connecting to EdgeHub Gateway.
func TrustBundle(workloadURI string) (*x509.CertPool, error) {
	tbr := TrustBundleResponse{}
	var err error

	// catch unix domain sockets URIs
	if strings.Contains(workloadURI, "unix://") {
		addr, err := net.ResolveUnixAddr("unix", strings.TrimPrefix(workloadURI, "unix://"))
		if err != nil {
			fmt.Printf("Failed to resolve: %v\n", err)
			return nil, err
		}

		setSharedUnixHTTPClient(addr.Name)
		response, err := sharedUnixHTTPClient.Get("http://iotedge" + "/trust-bundle?api-version=2019-11-05")
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}

		err = json.Unmarshal(body, &tbr)
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}
	} else {
		// format uri string
		uri := fmt.Sprintf("%strust-bundle?api-version=2019-11-05", workloadURI)

		// get http response and handle error
		resp, err := http.Get(uri)
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}
		defer resp.Body.Close()

		// read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}

		err = json.Unmarshal(body, &tbr)
		if err != nil {
			return nil, fmt.Errorf("tls: unable to append certificates: %w", err)
		}
	}

	p := x509.NewCertPool()
	if ok := p.AppendCertsFromPEM([]byte(tbr.Certificate)); !ok {
		err = fmt.Errorf("tls: unable to append certificates: count not append certs to pool")
	}

	return p, err
}
