# goaliniex

**goaliniex** is a lightweight and idiomatic Go SDK for integrating with the **Aliniex API**.

It is designed for simplicity, correctness, and ease of integration in production systems.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ Features

- 🚀 **Simple & Idiomatic API**
  Clean, readable interfaces that follow Go best practices.

- 🔐 **Secure Bearer-Token Auth**
  RSA-signed JWT client assertions exchanged for short-lived bearer access tokens, cached and auto-renewed per user.

- 🧵 **Context Support**
  Full support for `context.Context` for cancellation and timeouts.

- 🪶 **Minimal Dependencies**
  Only `github.com/golang-jwt/jwt/v5` for RS256 client-assertion signing; otherwise standard library.

## 📦 Installation

```bash
go get github.com/andyle182810/goaliniex/v2
```

## 🚀 Quick Start

```go
client, err := goaliniex.NewClient(
    "https://api.alixpay.com",
    partnerCode,
    privateKeyPEM,
    goaliniex.WithDebug(true),
)
if err != nil {
    log.Fatal(err)
}

resp, err := client.GetKycStatus(context.Background(), "user@example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Println("KYC status:", resp.Data.Status)
```

## 🧪 Testing

Integration tests automatically skip when required credentials are missing.

Required environment variables:

```bash
ALIX_PARTNER_CODE=your_partner_code
ALIX_TEST_EMAIL=your_test_email@example.com
ALIX_TEST_EMAIL_2=your_secondary_test_email@example.com  # optional
```

Private key file:

```text
./alix-private-key.pem
```

Run tests:

```bash
go test ./...
```

## 📬 Support

For bugs, questions, or feature requests:

- Open an issue on GitHub
  👉 [https://github.com/andyle182810/goaliniex/issues](https://github.com/andyle182810/goaliniex/issues)
- Email: **[andyle182810@gmail.com](mailto:andyle182810@gmail.com)**

## 📄 License

**goaliniex** is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.
