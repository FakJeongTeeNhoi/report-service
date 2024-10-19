#!/usr/bin/env python
import pika
import json
import uuid
from datetime import datetime, timedelta

# Create a new report
report = {
    "id": str(uuid.uuid4()),  # Convert UUID to string
    "reservation_id": str(uuid.uuid4()),  # Convert UUID to string
    "room_id": str(uuid.uuid4()),  # Convert UUID to string
    "space_name": "Conference Room",
    "space_id" : str(uuid.uuid4()),  # Convert UUID to string
    "status": "Confirmed",
    "start_datetime": datetime.now().astimezone().isoformat(),  # Ensure timezone is included
    "end_datetime": (datetime.now() + timedelta(hours=1)).astimezone().isoformat(),  # Ensure timezone is included
    "participant": [
        {"type": "Staff", "faculty": "Engineering"}
    ]
}


# Convert report to JSON
report_json = json.dumps(report)

connection = pika.BlockingConnection(
    pika.ConnectionParameters(host='localhost',port=5672))
channel = connection.channel()

channel.exchange_declare(exchange='Receiver', exchange_type='topic', durable=True)

routing_key = "reservation.*"

print(f" [x] Sent {routing_key}:{report_json}")
for i in range(1):
    channel.basic_publish(
        exchange='Receiver', routing_key=routing_key, body=report_json)
    print(f" [x] Sent {routing_key}:{report_json}")
connection.close()