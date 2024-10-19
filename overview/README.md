# Billing-Engine

## Overview
We are building a billing system for our Loan Engine. Basically the job of a billing engine is to provide the
- Loan schedule for a given loan( when am i supposed to pay how much)
- Outstanding Amount for a given loan
- Status of weather the customer is Delinquent or not

We offer loans to our customers a 50 week loan for Rp 5,000,000/-, and a flat interest rate of 10% per annum.

This means that when we create a new loan for the customer(say loan id 100) then it needs to provide the billing schedule for the loan as

W1: 110000<br/>
W2: 110000<br/>
W3: 110000<br/>
...<br/>
W50: 110000

The Borrower repays the Amount every week. (assume that borrower can only pay the exact amount of payable that week or not pay at all)

We need the ability to track the Outstanding balance of the loan (defined as pending amount the borrower needs to pay at any point) eg. at the beginning of the week it is 5,500,000/- and it decreases as the borrower continues to make repayment, at the end of the loan it should be 0/-

Some customers may miss repayments, If they miss 2 continuous repayments they are delinquent borrowers.

To cover up for missed payments customers will need to make repayments for the remaining amounts. ie if there are 2 pending payments they need to make 2 repayments(of the exact amount).

We need to track the borrower if the borrower is Delinquent( any borrower that who’s not paid for last 2 repayments)

## Database Architecture
![alt text](https://github.com/dimassantoso/be-assignment/blob/main/overview/UMLDiagram.png?raw=true)

## Run Biling Engine
### Prepare dependencies (PostgreSQL, Redis and OpenTracing)
```
$ docker-compose up -d
```

### Create database at PostgreSQL
```
create database "billing-engine";
```

### Go to `billing-engine` directory
```
$ cd ..
$ cd billing-engine
```

### Setup everything in `.env` file

### Run migration and seeding data
```
$ make migration
```
### Build and run service
```
$ make run
```
### Feature
- Borrower
1. Create Borrower
2. Update Borrower
3. Get All Borrower
4. Get Detail Borrower (contains active loan and delinquent status)
5. Check Delinquent
- Billing
1. Make Payment
- Loan
1. Create Payment (include generate payment schedule)
2. Installment Simulation
3. Get Outstanding

### Architecture
1. Golang
2. Clean Architecture
3. Restful API
4. Singleton Design Pattern

### Main Endpoint
#### Login
We used Bearer token as security authentication for every main business and Basic auth for login
```
```

```
curl --location 'localhost:8000/v1/auth/login' \
--header 'Authorization: Basic YmlsbGluZy1lbmdpbmU6YmlsbGluZy1lbmdpbmU=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "admin@example.com",
    "password": "P@SSword1234",
    "keep_sign_in": true
}'
```

#### Get Outstanding
```
curl --location 'localhost:8000/v1/loan/outstanding/borrower/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk3OTg3MzIsInN1YiI6IjAiLCJyb2xlIjoiIiwiYWRkaXRpb25hbCI6bnVsbH0.xfH2r5n-qWbeM7jGrOOIXzGWDkp_Px9f-mHvccYT5tU'
```

### Check IsDelinquent
```
curl --location 'localhost:8000/v1/borrower/1/delinquent-check' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk3OTg3MzIsInN1YiI6IjAiLCJyb2xlIjoiIiwiYWRkaXRpb25hbCI6bnVsbH0.xfH2r5n-qWbeM7jGrOOIXzGWDkp_Px9f-mHvccYT5tU'
```
also could be check in borrower detail data
```
curl --location 'localhost:8000/v1/borrower/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk3OTg3MzIsInN1YiI6IjAiLCJyb2xlIjoiIiwiYWRkaXRpb25hbCI6bnVsbH0.xfH2r5n-qWbeM7jGrOOIXzGWDkp_Px9f-mHvccYT5tU'
```

### Make Payment
```
curl --location 'localhost:8000/v1/billing/repayment' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk3OTg3MzIsInN1YiI6IjAiLCJyb2xlIjoiIiwiYWRkaXRpb25hbCI6bnVsbH0.xfH2r5n-qWbeM7jGrOOIXzGWDkp_Px9f-mHvccYT5tU' \
--header 'Content-Type: application/json' \
--data '{
    "billing_id": 1,
    "payment_method_id": 1
}'
```
