# 🗒️ Database

This project uses a MySQL database that is running on the same EC2 instance as the API backend (see [🌐 Deployment](deployment.md#-deployment)).

This allows the database to not be exposed to the internet except in the early phases of development for ease of testing and debugging.

See the models in use at [📦 Models](project-details.md#-models).

A [cron](https://en.wikipedia.org/wiki/Cron) job has been scheduled to run twice daily to backup the database (it dumps the database to a local password-protected file on the EC2 instance).
