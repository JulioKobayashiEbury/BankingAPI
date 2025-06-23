# BankingAPI
A Banking API using Golang.

- For Database it uses Firestore emulator called Spine3.
- It uses Echo as the web framework.
- It uses Zerolog for structured logging.
- It uses Swagger for API documentation.
- It uses PlantUML for sequence diagrams and other documentation.
- Conformant with RESTful API principles.

# Packages Used

- [echo](https://echo.labstack.com/) - Web framework for Golang
- [firestore](https://cloud.google.com/firestore) - NoSQL database by Google
- [zerolog]("github.com/rs/zerolog/log") - Structured logger for golang

# Documentation
All documentation is available in the [docs](docs) folder of the project. All sequence diagrams use the authorization diagram.

All documentation (except the swagger) is generated using [PlantUML](https://plantuml.com/).

The swagger is compliant with OpenAPI.

# Docker

# Firestore Emulator (Para rodar localmente)

- Download Firebase CLI
```bash
    npm install -g firebase-tools
```
- No diretório que o projeto firestore estará faça:
    - Crie um firebase.json que contenha um conteúdo simples
        Utilizei o nano para criar o arquivo com o seguinte conteúdo:
        "
            {
                "emulators": {
                    "firestore": {
                    "port": 8080
                    }
                }
            }
        "
    - Crie o diretório firestore-data:
    ```bash
        mkdir firestore-data
    ```
    - Inicie o emulador com um ID de projeto fictício (sempre inicie com o mesmo ID):
    ```bash
        firebase emulators:start --only firestore --import=./firestore-data --export-on-exit --project banking
    ```
    - Utilize o IP e porta dados ao iniciar o firestore e coloque em uma environment variable chamada "FIRESTORE_EMULATOR_HOST".



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
4. Set up Firestore
    - Make sure you have the Firestore emulator running. You can use the Docker command provided above to run the Spine3 emulator.
    - Ensure that the Firestore emulator is running on port the same port and ip as specified in the FIRESTORE_EMULATOR_HOST environment variable.

5. Run the application
    ```bash
    go run cmd/main.go
    ```
6. Access the API
    - Connect via other API or Front-End application or use a tool like Postman to access the API at `http://localhost:25565`
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
