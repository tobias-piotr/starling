# 🐦‍⬛ starling

Smart traveling assistant.

The idea behind Starling is simple: it takes a set of criteria for your trip
and it generates a result that contains all the necessary information. The
result is a document that you can use to plan your trip.

From a technical POV, Starling has two components: the API and the worker. The
API is obviously the entry point for the user. When the user is happy with the
criteria, the next step is to 'request' the trip result. The worker receives the
event via Redis Stream and then processes the task asynchronously. The result
is updated in the database and the user can fetch it later.

Initially, the app is just using AI to generate the results. Later iterations
could include calls to external APIs to get the offers, prices, flights, etc.

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

Next:

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
