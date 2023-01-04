## OneNUS Backend [22/23 CVWO Winter Assignment]

Backend for `OneNUS`.

## 🎮 Demo

You can find the **live version** of this project [here](https://app.onenus.link).

### Frontend

You can find the frontend that consumes this project [here](https://github.com/mfjkri/OneNUS).

## 💻 Project Overview

| Project aspect          | Technologies used                                                                                                                                                                                                                                                                                                                             |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Frontend                | Typescript + React                                                                                                                                                                                                                                                                                                                            |
| Backend (**this repo**) | Golang + Gin + Gorm<br>See [⚡️Technologies](docs/technologies-used.md#%EF%B8%8Ftechnologies) for more details.                                                                                                                                                                                                                               |
| Database                | MySQL                                                                                                                                                                                                                                                                                                                                         |
| Deployment plan         | Frontend: AWS S3 Bucket + Cloudfront<br>See [Frontend Deployment](https://github.com/mfjkri/OneNUS/blob/master/docs/deployment.md#-deployment) for more details.<br><br>Backend: AWS EC2 + Nginx (reverse proxy)<br>See [🌐 Deployment](docs/deployment.md#-deployment) for more details.<br><br>Database: AWS EC2 (same instance as backend) |

## 🛠 Building the project

### Prerequisites

1. `Go`

   Install [Go](https://go.dev/doc/install) if you have not done so yet.

### Installation

1. Clone this repo.
   ```
   $ git clone https://github.com/mfjkri/OneNUS-Backend.git
   ```
2. Change into the repo directory.

   ```
   $ cd OneNUS-Backend
   ```

3. Copy the template `.env` file.

   ```
   $ cp .env.example .env
   ```

   Modify the following environment variables in the new `.env` file accordingly:

   ```python
   PORT=8080 # Port number that the project  will be listening to
   DB="USERNAME:PASSWORD@tcp(HOSTNAME:PORT_NUMBER)/DATABASE_NAME?charset=utf8mb4&parseTime=True&loc=Local" # Credentials to connect to database
   JWT_SECRET=JWT_SECRET # Random string that is used to generate JWT tokens
   GIN_MODE="debug" # Set to either "debug" or "release" accordingly
   ```

4. All set!

   ```
   $ go run main.go
   ```

   This command will install all dependencies automatically when first ran.

## 📚 Table of Contents

- [⚡️Technologies](docs/technologies-used.md#%EF%B8%8Ftechnologies)
- [📦 Models](docs/project-details.md#-models)
- [🛣️ API Endpoints](docs/project-details.md#%EF%B8%8F-api-endpoints)
- [🗒️ Database](docs/database.md#%EF%B8%8F-database)
- [🌐 Deployment](docs/deployment.md#-deployment)
