# Access Token / Refresh Token generator

## tl;dr
```
cd SOME_FOLDER

mkdir cert
cd cert
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

cd ..
tokenGen
```

Browse 'https://localhost/tokenGen'

### 1. Prepare self signed certificate.

#### command sample(openssl)
`openssl genrsa -out server.key 2048`

`openssl ecparam -genkey -name secp384r1 -out server.key`

`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`

### 2. Deploy certifications.
`server.crt` and `server.key` file in `cert` directory.

### 3. Set  RedirectUri of your Box application as 'https://localhost'

### 4. Run tokenGen.
