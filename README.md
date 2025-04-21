# 🎓 The Grid: A Research Collaboration Platform

## 🚀 Project Overview

This project is a research collaboration platform designed to help **graduate students, postdocs, and faculty** connect across universities.

## 📌 Key Features

- **Project Collaboration**: Users can create, join, and manage collaborative research projects.
- **Search & Discovery**: Find collaborators and projects based on expertise and research interests.
- and more...

## Development Setup

#### 1. Prerequisites

Ensure you have the following installed:

- **Git**
- **Node.js & npm/yarn** (for frontend)
- **Go** (for backend)

#### 2. Clone the repository

```bash
git clone https://github.com/swecrocs/Crocs_SE_Project.git
cd Crocs_SE_Project
```

#### 3. Frontend Setup

- Navigate to the frontend folder
  ```bash
  cd frontend/SE_Project
  ```
- Install dependencies
  ```bash
  npm install   # or `npm ci` if you want a clean, lockfile‑exact install
  ```
- Run the dev server
  ```bash
  npm start     # should invoke `ng serve --open` and open http://localhost:4200
  ```
- After that you should see your app live at http://localhost:4200.

#### 4. Backend Setup

- Navigate to the `backend` directory
  ```bash
  cd backend
  ```
- Install dependencies
  ```bash
  go mod tidy
  ```
- Run the backend server
  ```bash
  go run main.go
  ```
- Need access to the backend API docs? Visit
  ```
  http://localhost:8080/swagger/index.html#/
  ```

## Team Members and Roles

| Name                     | Role      |
| ------------------------ | --------- |
| Gianfranco Cortes-Arroyo | Back-end  |
| Sri Vaishnavi Borusu     | Back-end  |
| Bo-Hao Wang              | Front-end |
| Sanket Jadhao            | Front-end |
