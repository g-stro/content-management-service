# Content Service

This repository implements a **content management microservice** in Go, built for managing content via a REST API.

---

## **Features**
- REST API to manage content:
    - `GET /content`: Retrieve all content.
    - `POST /content`: Create new content with associated details.
- Built with **clean architecture principles**.
- PostgreSQL for database management.
- Fully containerized with Docker and Docker Compose.
- Implements middleware for CORS support.
- Includes CI pipeline for testing and building.
- Unit and integration tests using Go's standard testing package.

---

## **Quick Start**

### **Prerequisites**
- [Go 1.24+](https://go.dev/dl/) installed.
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).
- PostgreSQL database (bundled with Docker Compose).

---

### **1. Clone the Repository**
```bash
git clone https://github.com/g-stro/content-service.git
cd content-service
```

---

### **2. Setup Environment**
Modify the .env.example as required.

Example environment variables:
```dotenv
SERVICE_PORT=8080

DB_USERNAME=test_user
DB_PASSWORD=test_password
DB_NAME=test_db
DB_HOST=postgres
DB_PORT=5432
DB_SSL_MODE=disable
DB_TIMEZONE=UTC
```

---

### **3. Run the Application Locally**

#### Using Docker Compose
1. **Build and Start Services**:
   ```bash
   make build
   ```
2. **Run the Services**:
   ```bash
   make up
   ```

3. **Access the API**:  
   The service will be available at `http://localhost:8080`.

4. **Stop Docker Containers**:
   ```bash
   make down
   ```

---

## **Using the API**
### Endpoints:
1. **`GET /content`**  
   Retrieve all content.  
   Example response:
   ```json
   {
     "status": "success",
     "data": {
       "content": [
         {
           "id": 1,
           "name": "Sample Name",
           "description": "Sample Description",
           "details": [
             {
               "content_type": "text",
               "value": "Sample Text"
             }
           ]
         }
       ]
     }
   }
   ```

2. **`POST /content`**  
   Create new content:  
   Example request body:
   ```json
   {
     "name": "New Content",
     "description": "Details about new content",
     "details": [
       {
         "content_type": "text",
         "value": "Sample Value"
       }
     ]
   }
   ```
   Example response:
   ```json
   {
     "status": "success",
     "data": {
       "id": 1,
       "name": "New Content",
       "created_at": "2025-05-13 10:52:07"
     }
   }
   ``` 

---

## **Database Schema**
The PostgreSQL schema is initialized with the following tables:

- `content`: Stores basic content data.
- `content_details`: Stores additional details associated with content.
- `content_type`: Stores types of content (e.g. text, image, video).

### Schema Setup
If you're running the service via Docker Compose, the schema is automatically initialized using `sql.sql`. To apply it manually:
```sql
psql -U <username> -d <database> -f sql.sql
```

---

## **Running Tests**

### **Unit Tests**
```bash
make unit-tests
```

### **Integration Tests**
Integration tests require both the service and the database to be running. If using Docker Compose:
```bash
make integration-tests
```

### **All Tests**
```bash
make tests
```

---

## **CI Pipeline**
The repository includes a GitHub Actions pipeline (`ci.yml`) that:
- Builds the Go application.
- Executes both unit and integration tests.
- Ensures any changes do not break the existing functionality.

---

## **Makefile Commands**
### **Basic Commands**:
- `make help`: Display available `make` targets.
- `make build`: Build and start Docker containers.
- `make up`: Start the service containers.
- `make clean`: Remove containers, images, volumes, and orphans for a clean environment.
- `make restart`: Restart the entire service stack.
- `make logs`: Tail logs from the running containers.

### **Testing Commands**:
- `make unit-tests`: Run unit tests.
- `make integration-tests`: Run integration tests.
- `make tests`: Run all tests.

---

## **License**
This project is open-source and is distributed under the [MIT License](LICENSE).
