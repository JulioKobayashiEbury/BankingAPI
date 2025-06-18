# BankingAPI
A Banking API using Golang, Firestore DB and RestAPI

# Packages Used

- [echo](https://echo.labstack.com/) - Web framework for Golang
- [firestore](https://cloud.google.com/firestore) - NoSQL database by Google
- [zerolog]("github.com/rs/zerolog/log") - Structured logger for golang

# Docker

# Firestore Emulator Spine3
sudo docker run --rm -p=4000:4000 -p=8080:8080 -p=9099:9099 -p=5001:5001 -p=9199:9199 -p=9000:9000 -p=8085:8085 --env "GCP_PROJECT=banking" --name database-fs-emulator spine3/firebase-emulator --import /firebase/data

# Firestore DB Structure
```- users (collection)
  - userID (document)
    - name: string
    - dcoument: string
    - password: string
    - status: bool
    - register_date: timestamp
```
``` - clients (collection)
    - clientID (document)
      - name: string
      - document: string
      - status: bool
      - register_date: timestamp
```
``` - accounts (collection)
    - accountID (document)
      - userID: string
      - clientID: string
      - balance: float64
      - agencyID: uint32
      - status: bool
      - register_date: timestamp
```
``` - transfers (collection)
    - transferID (document)
      - account_id: string
      - account_to: string
      - value: float64
      - register_date: timestamp
```
``` - deposits (collection)
    - depositID (document)
      - account_id: string
      - client_id: string
      - agency_id: uint32
      - deposit: float64
      - register_date: timestamp
```
``` - withdrawals (collection)
    - withdrawalID (document)
      - account_id: string
    - client_id: string
    - agency_id: uint32
    - withdrawal: float64
    - withdrawal_date: timestamp
    - status: string
```
```
    - automaticdebits (collection)
      - automaticDebitID (document)
        - account_id: string
        - client_id: string
        - agency_id: uint32
        - value: float64
        - status: bool
        - expiration_date: timestamp
        - register_date: timestamp
```
# How to run the project
1. Clone the repository
    ```bash
    git clone https://github.com/JulioKobayashiEbury/BankingAPI.git
    ```
2. Change directory to the project folder
    ```bash
    cd BankingAPI
    ```
3. Install dependencies
    ```bash
    go mod tidy
    ```
4. Set up Firestore credentials
    - Create a service account in the Google Cloud Console
    - Download the JSON key file
    - Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the path of the JSON key file
    ```bash
    export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-file.json"
    ```
5. Run the application
    ```bash
    go run cmd/main.go
    ```
6. Access the API
    - Open your browser or use a tool like Postman to access the API at `http://localhost:25565`
# API Endpoints
- **User Endpoints**
    - `POST /users` - Create a new user
    - `GET /users/:id` - Get user by ID
    - `PUT /users/:id` - Update user by ID
    - `DELETE /users/:id` - Delete user by ID
- **Client Endpoints**
    - `POST /clients` - Create a new client
    - `GET /clients/:id` - Get client by ID
    - `PUT /clients/:id` - Update client by ID
    - `DELETE /clients/:id` - Delete client by ID
- **Account Endpoints**
    - `POST /accounts` - Create a new account
    - `GET /accounts/:id` - Get account by ID
    - `PUT /accounts/:id` - Update account by ID
    - `DELETE /accounts/:id` - Delete account by ID
- **Transfer Endpoints**
    - `POST /transfers` - Create a new transfer
    - `GET /transfers/:id` - Get transfer by ID
    - `PUT /transfers/:id` - Update transfer by
    - `DELETE /transfers/:id` - Delete transfer by ID
- **Deposit Endpoints**
    - `POST /deposits` - Create a new deposit
    - `GET /deposits/:id` - Get deposit by ID
    - `PUT /deposits/:id` - Update deposit by ID
    - `DELETE /deposits/:id` - Delete deposit by ID
- **Withdrawal Endpoints**
    - `POST /withdrawals` - Create a new withdrawal
    - `GET /withdrawals/:id` - Get withdrawal by ID
    - `PUT /withdrawals/:id` - Update withdrawal by ID
    - `DELETE /withdrawals/:id` - Delete withdrawal by ID
- **Automatic Debit Endpoints**
    - `POST /automaticdebits` - Create a new automatic debit
    - `GET /automaticdebits/:id` - Get automatic debit by ID
    - `PUT /automaticdebits/:id` - Update automatic debit by ID
    - `DELETE /automaticdebits/:id` - Delete automatic debit by ID
# Testing
- To run tests, use the following command:
    ```bash
    go test ./...
    ```
# Contributing
- Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
