# 📡 Chat-WebSocket System (Built with Golang, Gin, WebSocket, Redis)

## Overview

This is a high-performance WebSocket chat room system built with Golang, Gin framework, WebSocket protocol, and Redis Pub/Sub. The project demonstrates:
- Multi-server cross-room real-time communication.
- Cross-server message broadcasting using Redis Pub/Sub.
- Real-time interactive features such as user join/leave rooms and sending messages.
- Scalable architecture design suitable for high-concurrency chat applications.
- Integration with Prometheus metrics monitoring for easy tracking of system operation status.

---

## 🏗 Project Structure

```
chat-websocket/
├── cmd/
│   └── server/
│       └── main.go           # API server main entry point
├── config/
│   └── config.go             # Loads environment variables and system configurations (Singleton)
├── db/
│   └── mysql.go              # GORM initializes MySQL connection
├── api/
│   ├── router.go             # Gin router setup
│   └── websocket_handler.go  # WebSocket connection handling logic
├── model/
│   ├── client.go             # Client data model
│   ├── message.go            # Message data model
│   ├── room.go               # Room data model
├── pkg/
│   └── metrics/            # Prometheus metrics definitions and initialization
│       └── metrics.go
├── redis/
│   ├── lock.go               # Redis distributed lock (Simplified RedLock)
│   ├── pubsub.go             # Redis Pub/Sub functionality encapsulation
│   ├── redis.go              # Redis Client initialization and connection management
├── repository/
│   ├── client_repository.go  # Client database operation encapsulation
│   └── message_repository.go # Message database operation encapsulation
├── service/
│   ├── message_service.go  # Message-related business logic
│   └── room_service.go     # Room-related business logic
├── usecase/                  # Application scenario Use Cases (consider moving to service or api handler)
│   ├── message_usecase.go  # Message processing Use Case (adjustable)
│   └── room_usecase.go     # Room management Use Case (adjustable)
├── .env                      # Environment variable settings
├── Dockerfile                # Dockerfile configuration
├── docker-compose.yml        # Docker Compose configuration (including MySQL, Redis, Prometheus, Grafana)
├── go.mod                    # Go Modules dependency management
├── go.sum                    # Go Modules dependency version verification
├── prometheus.yml            # Prometheus configuration file
├── README.md                 # Project documentation (this file)
```

---

## 🚀 Distributed Architecture

```
+--------------------------------+
|  Client (APP / Web Frontend)   |
+--------------------------------+
            ↓ WebSocket Connection
+--------------------------------+
|   Gin API Server (Golang)      |
+--------------------------------+
            ↓ Redis Pub/Sub
+--------------------------------+
|       Redis Cluster            |
+--------------------------------+
            ↓ MySQL (Message Persistence)
+--------------------------------+
|       MySQL Database           |
+--------------------------------+
            ↓ Prometheus (Monitoring)
+--------------------------------+
|       Prometheus Server        |
+--------------------------------+
            ↓ Grafana (Dashboard)
+--------------------------------+
|        Grafana Dashboard       |
+--------------------------------+
```

--- 
## Communication Flow

```
                                 Redis Pub/Sub (room1)
                                      ^       ^
                                     /         \
[server1:8080]         +-----------------------+         [server2:8080]
  User A            Consume Message <---------+           User B
    |                       |                       |
    |     Join Room "room1"       Join Room "room1"     |
    V                       V                       V
 Publish Message                                 (Wait for message)
    |                       |                       |
    +-------------------------------------------------+
                                Broadcast
    |                       |                       |
Broadcast to A                                Broadcast to B
```

## 🔑 Key Techniques & Features

### **1. WebSocket Protocol**
- Two-way real-time communication: Enables real-time message exchange between server and clients.
- Persistent connection: Reduces HTTP handshake overhead and improves communication efficiency.

### **2. Redis Pub/Sub**
- Cross-server broadcast: Utilizes Redis Pub/Sub to achieve message synchronization and distribution across multiple servers.
- Channel Subscription: Each chat room corresponds to a Redis channel for easy message routing.

### **3. Golang + Gin Framework**
- High Performance: Golang language and Gin framework provide excellent performance and concurrency processing capabilities.
- Lightweight Framework: Gin framework is simple and efficient, suitable for building API services.

### **4. MySQL Database**
- Message Persistence: Uses MySQL to store chat messages, ensuring data is not lost.
- GORM ORM: Uses GORM for database operations, simplifying the development process.

### **5. Prometheus Metrics Monitoring**
- Real-time Monitoring: Integrates Prometheus metrics to monitor key indicators such as WebSocket connection count, message read rate, etc.
- Grafana Dashboard:  Paired with Grafana to visualize monitoring data, making it easy to understand system operation status.

### **6. Multi-server Scalability**
- Stateless Servers: API servers are designed to be stateless, facilitating horizontal scaling.
- Redis Centralized State Management: Room lists and cross-server message broadcasting rely on Redis to achieve state sharing.

### **7. Simple Authentication**
- `sender_id` Parameter: User identity is passed through the `sender_id` URL parameter during WebSocket connection for basic authentication.

### **8. Graceful Shutdown**
- Signal Handling: Listens for signals like `SIGTERM` to safely shut down the HTTP server and Redis connections.

### **9. Docker & Docker Compose**
- Containerized Deployment: Uses Dockerfile and docker-compose.yml to achieve one-click deployment of MySQL, Redis, Prometheus, Grafana, and the application.
- Environment Consistency: Ensures consistency across development, testing, and production environments.

---
## 🚀 Quick Start

### **1. Prerequisites**
Ensure you have the following software installed:
- **Docker** & **Docker Compose**
- **Golang 1.22+**

### **2. Setup & Installation**
```
git clone **this repo**
cd chat-websocket
docker-compose up -d --build
```

### **3. Database Migration**
Run database migrations to create tables:
```
docker exec -it server-api /app/migrate -path=db/migrations \
    -database "mysql://root:123456@tcp(mysql:3306)/chat_websocket?charset=utf8mb4&parseTime=True&loc=Local" up
```

### **4. Prometheus & Grafana Monitoring**
- **Prometheus Targets Verification**:
    1. Access the Prometheus targets dashboard at `http://localhost:9090/targets?search=`.
    2. In the "Targets" page, under "State", you should see an "UP" status for your `chat-websocket` job, confirming Prometheus is successfully scraping metrics from your application.

- **Prometheus Metrics Endpoint Verification**:
    1. Open your browser and access the metrics endpoint of your application at `http://localhost:8080/metrics`.
    2. You should see a long list of Prometheus metrics in plain text format. This confirms that your application is correctly exposing Prometheus metrics.

- **Grafana Dashboard Setup**:
    1. Access Grafana at `http://localhost:3000` (default credentials are admin/admin).
    2. Import the provided Grafana dashboard JSON file (`assets/websocket_rev1.json`). You can do this by:
        - Going to "Dashboards" -> "Import".
        - Click "Upload JSON file" and select the `websocket_rev1.json` file from your `assets` folder (or wherever you saved it).
    3. After importing, you should see the "Chat WebSocket Metrics" dashboard with pre-configured panels for monitoring your WebSocket application.
    4. **Restart Prometheus and Grafana**: After importing the dashboard, it's recommended to restart Prometheus and Grafana containers to ensure the new dashboard is correctly loaded and data is being displayed. Run the following command in your terminal:
       ```
       docker restart prometheus grafana
       ```
- **Grafana Dashboard Image**:
  ![Grafana Dashboard](./assets/Grafana-websocket.png)

       


### **5. WebSocket Connection Testing**
Use a WebSocket client tool (such as Chrome browser developer tools Console, wscat, etc.) to connect to the WebSocket server:

- WebSocket URL: ws://localhost:8080/chat?sender_id=user123 (Replace sender_id parameter value)
```js
// Establish WebSocket connection (replace sender_id parameter value)
ws = new WebSocket("ws://localhost:8080/chat?sender_id=test_user");

// Listen for connection open event
ws.onopen = function() {
    console.log("WebSocket connection opened");
};

// Listen for message receive event
ws.onmessage = function(event) {
    console.log("Received message: " + event.data);
};

// Listen for connection close event
ws.onclose = function() {
    console.log("WebSocket connection closed");
};

// Listen for error event
ws.onerror = function(error) {
    console.error("WebSocket Error: ", error);
};

// Send message (join room)
ws.send(JSON.stringify({action: "join", room_id: "room101"}));

// Send message (chat message)
ws.send(JSON.stringify({action: "message", room_id: "room101", content: "Hello, Room 101!"}));

// Send message (leave room)
ws.send(JSON.stringify({action: "leave", room_id: "room101"}));

// Close connection
ws.close();
```

### **6. API Server Health Check**
```
curl -X GET "http://localhost:8080"
```

### **7. View Logs**
```
# API Server Logs
docker logs -f server-api
```