# Changelog

## Unreleased

### Added

- TLS support

   `--tls`              verify server with CA certificates installed on this system
   `--tls-insecure`     do NOT verify server (accept any certificate)
   `--tls-cafile value` verify server with CA certificate stored in specified file
   `--tls-capath value` verify server with CA certificates located under specified path

### Changed

- `--timeout` option now affects both dialing to server and RPC call (before that dialing had hard-coded 1s timeout)

## 1.0.0 - 2017-12-13

### Added

- Probing servers and services via gRPC health-checking protocol
