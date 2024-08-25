# billing-engine

## Architecture Design
![Architecture Design](images/architecture_design.jpeg)

## How To Run
run docker : 
```
docker-compose up -d
```

migrate up : 
```
migrate -path ./migrations -database "mysql://root:root@tcp(localhost:3306)/billing-engine" up
```

run service : 
```
go run main.go
```

## Specification
### Schedule  
![Sequence Diagram Schedule Task](images/Sequence%20Diagram%20-%20Schedule%20Task.png)
### Get Outstanding Balance
![Sequence Diagram Get Outsanding Balance](images/Sequence%20Diagram_Get%20outstanding%20balance.png)
### Is Delinquent
![Sequence Diagram Check User Is Delinquent](images/Sequence%20Diagram_check%20is%20delinquent%20user.png)
### MakePayment
![Sequence Diagram Make Payment](images/Sequence%20Diagram%20-%20Make%20Payment.png)

## Migrate
### Migrate UP
migrate -path ./migrations -database "mysql://root:root@tcp(localhost:3306)/billing-engine" up

### Migrate Down
migrate -path ./migrations -database "mysql://root:root@tcp(localhost:3306)/billing-engine" down


## Curl
### Create Loan
```curl --location 'localhost:9005/api/v1/create-loan' \
--header 'Content-Type: application/json' \
--data '{
    "username": "bambang",
    "amount": 50000000
}'
```

### Get Outstanding Balance
```curl --location --request GET 'localhost:9005/api/v1/get-outstanding' \
--header 'Content-Type: application/json' \
--data '{
    "username":"bambang"
}'
```

### Is Delinquent User
```curl --location --request GET 'localhost:9005/api/v1/is-delinquent' \
--header 'Content-Type: application/json' \
--data '{
    "username":"tapidah"
}'
```

### Make Payment
```
curl --location 'localhost:9005/api/v1/make-payment' \
--header 'Content-Type: application/json' \
--data '{
    "username": "bambang",
    "amount": 1100000
}'
```