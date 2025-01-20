# bestwing-assessment
# Event-Driven Real-Time Messaging Application

This project is an event-driven, real-time messaging application built with Go, PostgreSQL, Apache Kafka, and Server-Sent Events (SSE). The application follows a microservices architecture and utilizes Docker for containerization and Kubernetes for orchestration.

## Features
- **Authentication Service**: Secure login with stateless authentication.
- **Event Service**: Real-time event streaming using SSE.
- **Persistence Service**: Consumes events from Kafka and stores them in PostgreSQL.
- **Ingestion Service (Optional)**: Publishes events from a REST API to Kafka.
- **Kafka Integration**: Consumes and produces messages to/from Apache Kafka.
- **Docker and Kubernetes**: Containerized application with Kubernetes manifests for orchestration.

## Setup Instructions

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- Kubernetes (Minikube/Cluster)
- PostgreSQL
- Kafka

### Environment Variables
Create a `.env` file from the `.env.example` template and configure the following:
```
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=events-app
DB_USERNAME=app_user
DB_PASSWORD=password

KAFKA_BROKERS=kafka:9092
KAFKA_CONSUMER_GROUP_ID=events-app-console-consumer
KAFKA_OFFSET_RESET=earliest # it should be latest "as a real time" but for ease of testing we are using earliest for now!
KAFKA_EVENTS_TOPIC=events-topic
```

### Running the Application

#### Docker
1. Build and run the containers using Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. Access the application at [http://localhost:8000](http://localhost:8000).

#### Kubernetes
1. Apply the Kubernetes manifests:
   ```bash
   kubectl apply -f k8s/
   ```

2. Check the services running in your cluster:
   ```bash
   kubectl get pods
   ```

3. Access the application via the appropriate service URL.

#### Running Locally (Without Docker/Kubernetes) (preferred)
1. Make sure you have Go and PostgreSQL running locally.
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start the application:
   ```bash
   go run main.go
   ```

4. Access the application endpoints at [http://localhost:8000](http://localhost:8000).

5. Access the index.html which contians the SSE at resources/views/index.html

### API Endpoints

- **POST** `/login`: Login to the application (requires credentials `username: islacks, password: password`).
- **GET** `/events`: List all events (requires authentication).
- **POST** `/events`: Create a new event (requires authentication).
- **GET** `/events/listen`: Listen for real-time events (SSE/WebSocket).

### Docker Compose Setup

This project includes a `docker-compose.yml` for easy local development. It includes:
- Go application
- PostgreSQL
- Kafka

### Kubernetes Setup
#### Note: Im not the best at it but I tried to make it work, I will be happy to learn more about it and improve it.

The Kubernetes manifests in the `k8s/` directory help deploy the application in a Kubernetes cluster. This includes deployments for:
- Go application
- PostgreSQL
- Kafka

### File Structure

```
├── app
│   ├── console
│   ├── http
│   ├── models
│   ├── providers
│   └── services
├── bootstrap
│   └── app.go
├── configs
├── database
│   ├── migrations
│   └── seeders
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── main.go
├── README.md
├── routes
└── storage
```

## Troubleshooting

- Ensure that PostgreSQL and Kafka are running and accessible.
- Make sure the `.env` file is properly configured with the correct credentials for PostgreSQL and Kafka.

```yaml
### Kubernetes Manifests (Basic Structure)
## Namespace
apiVersion: v1
kind: Namespace
metadata:
  name: microservices
---
## Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: event-app
  namespace: microservices
spec:
  replicas: 2
  selector:
    matchLabels:
      app: event-app
  template:
    metadata:
      labels:
        app: event-app
    spec:
      containers:
        - name: event-app
          image: web_app/go
          ports:
            - containerPort: 8000
          env:
            - name: SUPERVISOR_GO_USER
              value: "app"
            - name: WWWUSER
              value: "1000"
            - name: GO_APP
              value: "1"
            - name: XDEBUG_MODE
              value: "off"
            - name: XDEBUG_CONFIG
              value: "client_host=host.docker.internal"
          volumeMounts:
            - name: app-volume
              mountPath: /var/www/html
            - name: go-cache
              mountPath: /var/tmp/go-cache
      volumes:
        - name: app-volume
          hostPath:
            path: /path/to/your/project
        - name: go-cache
          emptyDir: {}
---
## Service
apiVersion: v1
kind: Service
metadata:
  name: event-app
  namespace: microservices
spec:
  selector:
    app: event-app
  ports:
    - protocol: TCP
      port: 8000
      targetPort: 8000
  type: ClusterIP
```
