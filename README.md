# 🐦‍⬛ starling

Smart traveling assistant.

Initial plan:

- API to create requirements
  - Destination
  - Departure location (with km threshold)
  - Date
  - Budget
  - Additional requirements (how to spend the time)
- Event-driven flow to create the final document
- AI-based responses
- Requests to external APIs to get offers and prices
- Export functionality
- Discord bot as a UI

## Installation

Set up the `.env` file:

```
DATABASE_DSN="host=postgres port=5432 user=postgres dbname=starling password=postgres sslmode=disable"
```

Build the Docker images:

```
make build
```

Start the Docker containers:

```
make up
```
