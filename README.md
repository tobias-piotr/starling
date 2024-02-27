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

Final document:

- General description
- Attractions
- Weather
- Prices
- What to take (clothes, documents)
- Commuting

Extra:

- Retry steps

Next:

- README
- Some interface (templ + htmx + tailwind + daisyui)
- Patch endpoint
- Extra field in the result that would address additional requirements

## Installation

Set up the `.env` file:

```
DATABASE_DSN="host=postgres port=5432 user=postgres dbname=starling password=postgres sslmode=disable"
REDIS_ADDR=redis:6379
REDIS_STREAM=trips
REDIS_FAILURE_STREAM=starling-trips-failures
REDIS_CGROUP=starling
REDIS_CNAME=starling
OPENAI_KEY=
```

Build the Docker images:

```
make build
```

Start the Docker containers:

```
make up
```
