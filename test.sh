#!/bin/bash

echo "🏥 NAVBAT CHIPTALARINI JO'NATISH"
echo "================================"
echo ""

# Xonalar ro'yxati
rooms_200=(201 202 203 204 205)
rooms_300=(301 302 303 304 305)

counter=20
total=11

for i in $(seq 1 $total); do
  # Xonani avtomatik tanlash
  if [ $counter -le 24 ]; then
    # 200-seriya xonalar
    idx=$(( (counter - 20) % 5 ))
    room="${rooms_200[$idx]}"
  else
    # 300-seriya xonalar
    idx=$(( (counter - 25) % 5 ))
    room="${rooms_300[$idx]}"
  fi
  
  echo "📤 [$i/$total] Sending A-$counter → Xona $room"
  
  curl -s -X POST http://192.168.100.86:8080/api/audio/announcement \
    -H "Content-Type: application/json" \
    -d "{\"ticket_id\":\"$counter\",\"room_number\":\"$room\",\"queue_number\":\"A-$counter\"}" \
    | jq -r '.status' > /dev/null
  
  if [ $? -eq 0 ]; then
    echo "   ✅ Success"
  else
    echo "   ❌ Failed"
  fi
  
  counter=$((counter + 1))
  sleep 0.1
done

echo ""
echo "✅ TUGADI! $total ta request jo'natildi"