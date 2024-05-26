# Email Reminder Service

This project is an email reminder service that allows users to create reminders and receive emails when the reminders are due.


### Prerequisites

- Docker
- Docker Compose

### Running the Project

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Copy the `.env.Example` file and rename the copy to `.env`. Fill in the necessary environment variables.
4. Run the following command to start the services:

```bash
docker-compose up
```

## API Endpoints

- `POST /reminders`: Create a new reminder.
- `GET /reminders`: Get a list of all reminders.
- `GET /reminders/{id}`: Get a specific reminder by its ID.
- `PUT /reminders/{id}`: Update a specific reminder by its ID.
- `DELETE /reminders/{id}`: Delete a specific reminder by its ID.

## Sending Reminders via Email

The service uses a background worker to periodically check for reminders that are due. When a reminder is due, the service sends an email to the user with the reminder's message. The email sending functionality is implemented using the MailerSend API. The API key for MailerSend is specified in the `.env` file.