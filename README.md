# CVWO Assignment Project

## TODO:

1. Add starring of posts functionality

Last updated: 19/12/22

<br/>

# Table of Contents

- [Demo](#demo)
- [Getting Started](#getting-started)
  - [Pre-Requistes](#pre-requistes)
  - [Installation](#installation)
- [Deployment](#deployment)

<br/>

# Demo

You can find the live version of the frontend that this project serves [here](https://app.onenus.link).

<br/>

# Getting Started

## Pre-Requistes

1. `Go`

   Install [Go](https://go.dev/doc/install) if you have not done so yet.

<br/>

## Installation

1. Clone this repo.
   ```
   $ git clone https://github.com/mfjkri/One-NUS-Backend.git
   ```
2. Change into the repo directory.

   ```
   $ cd One-NUS-Backend
   ```

3. Required config files.

   Create a dotenv file `.env` under the root directory with the following variables:

   ```python
   PORT = 8080 # Port number that the backend will be listening to
   DB = "host=$HOSTNAME user=$USERNAME password=$PASSWORD dbname=$DATABASE_NAME port=$PORT_NUMBER sslmode=disable" # Credentials to connect to database
   JWT_SECRET: JWT_SECRET # Random string that is used to generate JWT tokens
   ```

4. All set!

   ```
   $ go run main.go
   ```

   This will automatically install all dependecies when first executed.

<br/>

# Deployment

This app is deployed in an [AWS EC2 instance](https://aws.amazon.com/ec2/) with a reverse-proxy using [nginx](https://www.nginx.com).

The EC2 instance is allocated an elastic IP that is routed to by [Route 53](https://aws.amazon.com/route53/).

Signed SSL certificate for the subdomain is provided by [Let's Encrypt](https://letsencrypt.org/).
