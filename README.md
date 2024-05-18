# Exchange rate notifier



## Set up

Create your .env file to specify SMTP provider and credentials. You can use following template:

```
# SMTP authorization parameters
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your google email>
SMTP_PASS=<your password>

# Database source parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=emaildb
```

A comprehensive guide how to set up google credentials: [link](https://www.youtube.com/watch?v=Ov-5g4lU3NE)