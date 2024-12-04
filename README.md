# Fling Service

A backend system for a Dating Mobile App called **Fling**.


## Tech Stacks

- **Go** - <https://go.dev/>
- **Gin Web Framework** - <https://gin-gonic.com/>
- **SQLC** - <https://sqlc.dev/>
- **PostgreSQL** - <https://www.postgresql.org/>
- **Docker** - <https://www.docker.com/>
- **Kubernetes** - <https://kubernetes.io/>


## Prerequisites

- Install `Docker Desktop` - https://www.docker.com/products/docker-desktop/
- Install `Makefile` -  https://www.gnu.org/software/make/
- Install `Go` - https://go.dev/dl/
- Install `Postman` - https://www.postman.com/downloads/

## Getting Started

Clone the project:

```bash
$ git clone https://github.com/rhmnmbr/fling-service
$ cd fling-service
$ cp app.env.example app.env
```

### How to Run Locally

1. Make sure docker is up and running.
2. Go to `src` folder.
3. Execute `make network` to create network inside docker.
4. Execute `make postgres` to create postgres container.
5. Execute `make createdb` to create the database.
6. Execute `make server` to run the app.
7. Import postman collections from the repo and test the endpoints.

### Using Docker Compose

1. Make sure docker is up and running.
2. Go to `src` folder.
3. Execute `docker compose up -d --build` to create network inside docker.

### How to Deploy to Local Kubernetes

1. Make sure docker is up and running.
2. Enable `kubernetes` in the setting inside docker desktop.
3. Execute `docker build -f Dockerfile -t fling-svc .` inside the `src` folder to build the image of the app.
4. Execute `kubectl create secret tls fling-tls --key server.key --cert server.crt` inside the `infra/devcerts` folder to register the ssl certificate secrite in the kubernetes.
5. Execute `kubectl apply -f .\ingress\ingress-depl.yml` inside the `infra` folder to deploy the ingress.
6. Execute `kubectl apply -f .\K8S\config.yml` to deploy the config.
7. Execute `kubectl apply -f .\K8S\local-pvc.yml` to deploy the persistent volume claim.
8. Execute `kubectl apply -f .\K8S\postgres-depl.yml` to deploy the postgres service.
9. Execute `kubectl apply -f .\K8S\fling-depl.yml` to deploy the backend service.
10. Execute `kubectl apply -f .\K8S\ingress-svc.yml` to deploy the ingress service.

## Endpoints
1. Sign up
   endpoint: `POST` `{{base_url}}/users/sign-up`
   request body:
   ```json
    {
      "email": "aman@yopmail.com",
      "password": "password",
      "phone": "+618123456789",
      "first_name": "Aman",
      "birth_date": "1992-12-04",
      "gender": "male",
      "location_info": "Bandung, Indonesia",
      "bio": "Hello."
    }
   ```
3. Login
   endpoint: `POST` `{{base_url}}/users/login`
   request body:
   ```json
    {
      "email": "aman@yopmail.com",
      "password": "password"
    }
   ```

## How to Test and Lint

Make sure you are in the `src` folder

- Test

  ```bash
  $ make test
  ```

- Lint

  ```bash
  $ make lint
  ```


## Repo Structure

```
├── infra/
  └── devcerts/         * ssl certificate.
  └── ingress/          * ingress deployment file.
  └── K8S/              * kubernetes related files.
├── postman/    * postman collections.
├── src/
  └── api/          * HTTP server init, routing and handlers.
  └── db/
    └── migration/      * database migrations.
    └── mock/           * mock for the store.
    └── query/          * necessary queries for SQLC generation.
    └── sqlc/           * SQLC generated files & store.
  └── token/        * token maker & payload.
  └── util/         * utility files.
└── ...
```
