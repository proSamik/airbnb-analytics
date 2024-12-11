import json
import random
from datetime import datetime, timedelta
import string

def generate_room_id():
    letter = random.choice(string.ascii_uppercase)
    number = random.randint(100, 999)
    return f"{letter}{number}"

def generate_dates():
    start_date = datetime.now()  # Start from today
    # Generate 7 months of data
    end_date = start_date + timedelta(days=7*30)  # Approximate 7 months
    
    dates = []
    current_date = start_date
    
    while current_date <= end_date:
        dates.append(current_date.strftime("%Y-%m-%d"))
        current_date += timedelta(days=1)
    
    return dates

def generate_db_json():
    room_ids = [generate_room_id() for _ in range(10)]
    dates = generate_dates()
    
    db = {"rooms": {}}
    
    for room_id in room_ids:
        base_rate = random.randint(80, 200)
        room_data = []
        
        for date in dates:
            rate_variation = random.uniform(0.8, 1.2)
            daily_rate = round(base_rate * rate_variation, 2)
            is_booked = random.random() < 0.6
            
            booking = {
                "date": date,
                "is_booked": is_booked,
                "rate": daily_rate
            }
            room_data.append(booking)
        
        db["rooms"][room_id] = room_data

    with open('db.json', 'w') as f:
        json.dump(db, f, indent=2)
    
    return room_ids

room_ids = generate_db_json()
print("Generated db.json with the following room IDs:")
for room_id in sorted(room_ids):
    print(f"- {room_id}")
