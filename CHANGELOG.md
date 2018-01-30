# Changelog

## 1.1.0 - 2018-01-30

### Added

- TLS support

   `--tls`              verify server with CA certificates installed on this system
   `--tls-insecure`     do NOT verify server (accept any certificate)
   `--tls-cafile value` verify server with CA certificate stored in specified file
   `--tls-capath value` verify server with CA certificates located under specified path
- Windows build

### Changed

- `--timeout` option now affects both dialing to server and RPC call (before that dialing had hard-coded 1s timeout)
- more descriptive messages are printed out in case of failure

## 1.0.0 - 2017-12-13

### Added

- Probing servers and services via gRPC health-checking protocol
