# How to run Currency Exchange App
## Install GoLang
``` brew install golang```
## Init Packages
``` go mod init ```

## Install all dependencies
``` go mod tidy ```


## Run Server
``` go run cmd/main.go ```

## Install Redis
``` brew install redis ```

## Endpoints

### Home

**URL:** `/`  
**Method:** `GET`  
**Description:** Home handler.

**Curl Command:**
```sh
curl -X GET http://localhost:8080/

Login

URL: /login
Method: POST
Description: Authenticate and receive a JWT token.

Curl Command: 
curl -X POST http://localhost:8080/login -d '{"username":"admin","password":"admin"}' -H "Content-Type: application/json"
```

### WebSocket Connection

**URL:** `/ws`  
**Method:** `GET`  
**Description:** Connect to WebSocket to receive real-time exchange rate updates.

**WebSocket Connection:**
```
wscat -c ws://localhost:8080/ws -H "Authorization: Bearer $TOKEN"
```

### Get Historical Exchange Rates

**URL:** `/historical-rates/{currency}`  
**Method:** `GET`  
**Description:** Retrieve historical exchange rates for a specific currency.

**Curl Command:**  
```
curl -X GET http://localhost:8080/historical-rates/USD -H "Authorization: Bearer $TOKEN"
```

### Update Exchange Rate

**URL:** `/updateExchangeRate`  
**Method:** `POST`  
**Description:** Update the exchange rate for a specific currency. Requires JWT token and is rate-limited.

**Curl Command:**  
```
curl -X POST http://localhost:8080/updateExchangeRate -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"currency": "USD", "rate": 1.25}'
```


### WebSocket Client

To test WebSocket connection, you can use wscat:
```
npm install -g wscat
wscat -c ws://localhost:8080/ws -H "Authorization: Bearer $TOKEN"
```
