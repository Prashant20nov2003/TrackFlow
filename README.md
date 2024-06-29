# TrackFlow

### A Scalable and High-Performance Tracking System

---

## Overview

**TrackFlow** is an advanced tracking system designed to handle large volumes of tracking data with high efficiency. The system is built using Go for backend services and features a range of capabilities including geolocation integration, database optimization, and real-time metrics visualization. TrackFlow leverages modern technologies such as Docker for container management and ClickHouse for high-volume data processing, making it an ideal solution for developers seeking robust data tracking and analytics.

---

## Key Features

- **Client-Side Tracker**: JavaScript-based tracking for capturing page views and user interactions.
- **Go Backend API**: High-performance API designed for efficient data handling and processing.
- **Geolocation Services**: Real-time IP address to geographic location mapping.
- **PostgreSQL and ClickHouse Integration**: Efficient data storage and processing with a migration path from PostgreSQL to ClickHouse.
- **Dockerized Development**: Simplified setup and deployment with Docker.
- **Performance Testing and Traffic Generation**: Tools for simulating high-load scenarios and testing system performance.
- **Database Optimization**: Enhanced query performance and data handling with ClickHouse.
- **CLI Dashboard**: Command-line interface for real-time system monitoring and metrics visualization.
- **VM Deployment**: Scalable deployment strategies for virtual machines.

---

## Technology Stack

- **Go**: Core language for backend services.
- **Chi**: Lightweight and idiomatic router for Go applications.
- **Docker SDK**: Comprehensive container management for streamlined development.
- **PostgreSQL**: Initial relational database setup.
- **ClickHouse**: High-performance columnar database for large-scale data handling.
- **BoltDB**: Lightweight key-value store for persistent task management.
- **goprocinfo**: Library for system and process metrics collection in Go.

---

## System Architecture

### Components

1. **Client-Side Tracker**:
   - JavaScript-based solution for tracking page views and user interactions.
   - Handles missing data gracefully and supports asynchronous data collection.

2. **Backend API**:
   - Built with Go and Chi for high performance and low latency.
   - Supports tracking data processing, geolocation lookups, and Docker-based testing.

3. **Data Storage**:
   - Initially stores data in PostgreSQL, with optimized migration to ClickHouse for high-volume processing.
   - Ensures fast and efficient CRUD operations with data normalization techniques.

4. **Docker Integration**:
   - Docker SDK utilized for managing container lifecycle during development and testing.
   - Simplifies setup and ensures consistent environments across different stages.

5. **Real-Time Monitoring**:
   - CLI dashboard provides a real-time view of system metrics and performance.
   - goprocinfo used for detailed process and system metric collection.

6. **Deployment**:
   - Supports scalable deployment to virtual machines, with robust configuration and API authentication.

---

## Installation and Setup

### Prerequisites

- **Go** (version 1.18+)
- **Docker** (version 20.10+)
- **PostgreSQL** and **ClickHouse** databases