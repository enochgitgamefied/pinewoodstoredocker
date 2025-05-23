# 🛒 PinewoodStore Docker – Official Vulnerable App Release

Welcome to the **official release of PinewoodStore**, a deliberately vulnerable Java Spring Boot application containerized for quick deployment and secure experimentation.

This application is designed for:

* ✅ Security enthusiasts
* ✅ Penetration testers
* ✅ Ethical hackers
* ✅ Students and educators in cybersecurity

> ⚠️ **DISCLAIMER:** PinewoodStore is **for educational and ethical use only**. Always run it in isolated environments that you own or control.

---

## 🚀 Getting Started

### ✅ Prerequisites

Ensure you have the following installed:

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)

---

### 🔧 Step-by-Step Setup

1. **Clone the Repository**

```bash
git clone https://github.com/enochgitgamefied/pinewoodstoredocker.git
cd pinewoodstoredocker
```

2. **Build and Start the Application**

```bash
docker-compose up --build
```

* The application will be accessible at:
  📍 `http://localhost:8088`

---

### 🌐 Remote Server Configuration

If you're running PinewoodStore on a **remote server** and want to access it from your local machine or browser:

1. **Edit the `docker-compose.yml` file**

Replace all instances of `localhost` in the `environment` section with the **public IP address** of your remote server:

```yaml
environment:
  - API_BASE_URL=http://<your-server-ip>:8088
  - API_BASE_URL_DIAG=http://<your-server-ip>:84
  - API_BASE_URL_PINEWOODSTORE=http://<your-server-ip>:8088/welcome
```

2. **Rebuild and restart the containers**

```bash
docker-compose down
docker-compose up --build
```

---

## 🔐 Vulnerabilities Included in PinewoodStore

PinewoodStore includes a wide range of common and impactful web vulnerabilities for safe, hands-on learning. Each entry links to a corresponding demonstration video.

| #  | Vulnerability Name                                                                         | Demo Video Link        |
| -- | ------------------------------------------------------------------------------------------ | ---------------------- |
| 1  | [Server Side Request Forgery (SSRF)](https://www.youtube.com/watch?v=PAx_RSrKoVc)          | ✅                      |
| 2  | [Reflected Cross Site Scripting (XSS)](https://www.youtube.com/watch?v=cJTCzxwMVkc&t=527s) | ✅                      |
| 3  | [Stored Cross Site Scripting (XSS)](https://www.youtube.com/watch?v=1A61ZtWqWH8&t=635s)    | ✅                      |
| 4  | [External XML Entity (XXE)](https://www.youtube.com/watch?v=e552T0B4hb0&t=270s)            | ✅                      |
| 5  | [Remote File Inclusion (RFI)](https://www.youtube.com/watch?v=N53OopT6msU)                 | ✅                      |
| 6  | [Local File Inclusion (LFI)](https://www.youtube.com/watch?v=TfgQJXVT_yk)                  | ✅                      |
| 7  | Directory Traversal                                                                        | 📋 (Video Coming Soon) |
| 8  | [Remote Code Execution (RCE)](https://www.youtube.com/watch?v=z4Xf9LYXc9o&t=714s)          | ✅                      |
| 9  | [Server Side JSON Injection](https://www.youtube.com/watch?v=iycTLPaPBD8&t=1079s)          | ✅                      |
| 10 | [JWT Token Tampering](https://www.youtube.com/watch?v=48k7fSGCFAQ&t=1s)                    | ✅                      |
| 11 | [DOM Based XSS](https://www.youtube.com/watch?v=9cUWwr0F0Ng)                               | ✅                      |
| 12 | [SPEL Expression Injection (Spring Boot)](https://www.youtube.com/watch?v=J0t8KNq9c0s)     | ✅                      |
| 13 | [Command Injection](https://www.youtube.com/watch?v=z4Xf9LYXc9o&t=714s)                    | ✅                      |
| 14 | [Request Smuggling](https://www.youtube.com/watch?v=bGwsF3Q3tFs)                           | ✅                      |
| 15 | [File Upload Vulnerability 1](https://www.youtube.com/watch?v=7AbQ9rpKI74)                 | ✅                      |
| 16 | [Command Injection](https://www.youtube.com/watch?v=7AbQ9rpKI74)                 | ✅                      |

> 💡 Use tools like **Burp Suite**, **OWASP ZAP**, or custom payloads to explore and exploit these issues.

---

## 🧷 Endpoints Overview

| Component      | URL                             |
| -------------- | ------------------------------- |
| Main App       | `http://localhost:8088`         |
| Diagnostic API | `http://localhost:84`           |
| Welcome Page   | `http://localhost:8088/welcome` |

---

## 📬 Support & Contribution

* 💬 [Open an Issue](https://github.com/enochgitgamefied/pinewoodstoredocker/issues) for bugs or improvements
* 📦 Submit Pull Requests to contribute new features or enhancements

---

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](./LICENSE) file for details.
