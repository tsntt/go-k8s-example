# Go Microservice

This is a simple Go microservice that provides a RESTful API for managing users. It's designed to be deployed with Docker and Kubernetes.

## About the Project

The project is a simple Go web server with the following features:

*   A `/users` endpoint to manage users (GET and POST).
*   Connection to a PostgreSQL database.
*   Health check endpoints (`/healthz` and `/readiness`).
*   Prometheus metrics at the `/metrics` endpoint.

## Getting Started

To get the project running locally, you'll need to have [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-microservices.git
    cd go-microservices
    ```

2.  **Set up the database:**

    The application requires a PostgreSQL database. You can use the provided `compose.yaml` to spin up a database container. Uncomment the `db` service in the `compose.yaml` file and create a `db/password.txt` file with a password for the database.

3.  **Run the application:**

    ```bash
    docker compose up --build
    ```

    The application will be available at `http://localhost:8000`.

## API Endpoints

### `/users`

*   **GET /**: Retrieves a list of all users.
*   **POST /**: Creates a new user. The request body should be a JSON object with the following fields:
    *   `id` (string)
    *   `name` (string)

## Deployment

The project includes Kubernetes configuration files in the `kubernetes` directory for deploying the application to a Kubernetes cluster.

The Kubernetes setup includes:

*   A Deployment for the PostgreSQL database.
*   A Service for the PostgreSQL database.
*   A Deployment for the Go API (v1 and a canary v2).
*   A LoadBalancer Service for the Go API.

Before deploying, you'll need to create a secret called `db-credentials` with the following keys:

*   `username`
*   `password`
*   `dbname`

## Infrastructure as Code (IaC)

The project also includes Terraform configuration files in the `terraform` directory for provisioning infrastructure on AWS.

The Terraform setup will create the following resources:

*   A VPC (Virtual Private Cloud).
*   A Security Group that allows inbound traffic on ports 22 (SSH) and 80 (HTTP).
*   An EC2 instance.

### Terraform Files

*   `main.tf`: The main configuration file that defines the AWS resources to be created.
*   `variables.tf`: Defines variables for the Terraform configuration.
*   `outputs.tf`: Defines the outputs of the Terraform configuration, such as the EC2 instance ID and the security group ID.

To use the Terraform configuration, you'll need to have [Terraform](https://www.terraform.io/) installed and configured with your AWS credentials. Then, you can run the following commands in the `terraform` directory:

```bash
terraform init
terraform plan
terraform apply
```