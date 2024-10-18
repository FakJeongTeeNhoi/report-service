## How to Run

To start the application using Docker Compose, run the following command:

```bash
docker-compose up --build -d
```

## How to run mock data

To run the mock data, run the following command:

```bash
cd test_rabbitmq
python producer.py
```

## get stats api

Get stats api is on the /api/reports/{space_name}

{space_name} is the name of the space you want to get the stats for.
