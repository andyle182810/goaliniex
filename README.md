# goaliniex

**goaliniex** is a lightweight and idiomatic Go SDK for integrating with the **Aliniex API**.

It is designed for simplicity, correctness, and ease of integration in production systems.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ Features

- 🚀 **Simple & Idiomatic API**
  Clean, readable interfaces that follow Go best practices.

- 🔐 **Built-in Request Signing**
  Secure RSA-based request signing handled internally.

- 🧵 **Context Support**
  Full support for `context.Context` for cancellation and timeouts.

- 🪶 **Zero External Dependencies**
  Uses only the Go standard library.

## 📦 Installation

```bash
go get github.com/andyle182810/goaliniex
```

---

## 🚀 Quick Start

```go
client, err := goaliniex.NewClient(
    "https://api.aliniex.com",
    partnerCode,
    secretKey,
    privateKeyPEM,
    goaliniex.WithDebug(true),
)
if err != nil {
    log.Fatal(err)
}

resp, err := client.GetKycInformation(
    context.Background(),
    &goaliniex.KycInformationRequest{
        UserEmail: "user@example.com",
    },
)
if err != nil {
    log.Fatal(err)
}

fmt.Println("KYC status:", resp.Data.KycStatus)
```

## 🧪 Testing

Integration tests automatically skip when required credentials are missing.

Required environment variables:

```bash
ALIX_PARTNER_CODE=your_partner_code
ALIX_SECRET_KEY=your_secret_key
```

Private key file:

```text
./alix-private-key.pem
```

Run tests:

```bash
go test ./...
```

## 🤝 Contributing

## 📬 Support

For bugs, questions, or feature requests:

- Open an issue on GitHub
  👉 [https://github.com/andyle182810/goaliniex/issues](https://github.com/andyle182810/goaliniex/issues)
- Email: **[andyle182810@gmail.com](mailto:andyle182810@gmail.com)**

## 📄 License

**goaliniex** is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.
