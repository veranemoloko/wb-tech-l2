# ⏱️ NTP Time Fetcher

A simple **Go** project that fetches the **current exact time** from an NTP (Network Time Protocol) server and prints it to the console.  

Built as a **Go module**, using [`github.com/beevik/ntp`](https://pkg.go.dev/github.com/beevik/ntp) to get accurate time.

---

## Features

-  Fetch current time from any NTP server (default: `pool.ntp.org`)  
-  Handles errors gracefully and logs them using Go's structured `slog` logger  
-  Exits with a non-zero status code on error  
- Idiomatic Go code that passes `go vet` and `golint` checks  

---

## ⚙️ Prerequisites

- Go **1.21+**  
- Internet connection to access NTP servers  

---

##  Installation

Clone the repository and set up the Go module:

```bash
git clone https://github.com/veranemoloko/wb-tech-l2.git
cd wb-tech-l2/ntp_time_fetcher
go mod tidy
```
---

##  Usage

```bash
go run main.go --server time.google.com
```
or
```bash
go run main.go
```

## How It Works
- Parse the --server flag (default: pool.ntp.org)
- Call getNTPTime(server) to fetch the current time using github.com/beevik/ntp
- If an error occurs, log it with slog.Error and exit with code 1
- Print the current time to the console
