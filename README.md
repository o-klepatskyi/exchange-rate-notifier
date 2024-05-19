# Exchange rate notifier

Small service for getting updates of current USD to UAH rate.

## Features

- service automatically fetches current USD-UAH rate from [Monobank API](https://api.monobank.ua/docs/index.html)
- service automatically sends regular updates to all subscribers every day
- get current rate using `GET /rate`
- subscribe any email using `POST /subscribe`
- send manual updated for all subscribed emails using `POST /sendEmails`

## Setup
### Define parameters
Create your .env file in the project root folder to specify SMTP and database parameters. You can use the following template:

```
# SMTP authorization parameters
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your google email>
SMTP_PASS=<your password>

# Database source parameters
DB_HOST=localhost # ignored by docker compose!
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=emaildb
```

A comprehensive guide how to set up google credentials to use in apps: [link](https://www.youtube.com/watch?v=Ov-5g4lU3NE)

### Build and run
You can use following commands to get the service up and running:
```
docker compose build
docker compose run
```

You can do this manually by running `go run .` but make sure to change default database parameters and have PostgreSQL Server running locally or remotely.

## How it works

Application consists of three modules: `database` for connecting to database and to add and get subscription data, `mailsender` for working with SMTP to send current rate over email, and `ratefetcher` to fetch current rate from open API. <br>

First, database connection is initialized and two background processes are started using goroutines for 1. automatic fetching of current rate, 2. automatic email sends to all subscribed emails. <br>

Current rate is automatically fetched from remote API using HTTPS. It is done in reasonable periods of time so remote API is not overwhelmed. Rate value is then cached in memory on each update. Current rate can be fetched from `/rate`. <br>

Email updates can also be manually sent from `/sendEmails` and it does not conflict with automatic updates. New email can be subscribed to updates using `/subscribe`, if email is invalid or it is already subscribed, HTTP status 400 is returned. Email unique contraint is checked by database.

Service uses PostgreSQL for storing subscribers data (in this case only email address). Docker Compose is set up to run database on the separate container but app can use any local or remote database for storing data. Database schema migration is relatively small and is done on each service restart from the service itself.