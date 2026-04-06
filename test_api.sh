#!/bin/bash

# Примеры curl команд для тестирования API периодических задач

API_URL="http://localhost:8080/api/v1"

echo "=== Task Service API Testing Examples ==="
echo ""

# 1. Создание ежедневной задачи
echo "1. Create daily task (every day)"
curl -X POST $API_URL/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Daily patient review",
    "description": "Review all active patients",
    "status": "new",
    "recurrence": {
      "recurrence_type": "daily",
      "interval_days": 1
    }
  }' | jq .

echo ""
echo ""

# 2. Создание ежемесячной задачи
echo "2. Create monthly task (on 15th of each month)"
curl -X POST $API_URL/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Monthly report",
    "description": "Generate and submit monthly report",
    "status": "new",
    "recurrence": {
      "recurrence_type": "monthly",
      "monthly_day_of_month": 15
    }
  }' | jq .

echo ""
echo ""

# 3. Создание еженедельной задачи
echo "3. Create weekly task (Mon, Wed, Fri)"
curl -X POST $API_URL/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Weekly inventory check",
    "description": "Check inventory on these days",
    "status": "new",
    "recurrence": {
      "recurrence_type": "weekly",
      "weekly_days": [0, 2, 4]
    }
  }' | jq .

echo ""
echo ""

# 4. Создание задачи на четные дни
echo "4. Create task for even days only"
curl -X POST $API_URL/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Even day task",
    "description": "Execute only on even days",
    "status": "new",
    "recurrence": {
      "recurrence_type": "odd_even",
      "odd_even_type": "even"
    }
  }' | jq .

echo ""
echo ""

# 5. Создание задачи на конкретные даты
echo "5. Create task on specific dates"
curl -X POST $API_URL/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Special events",
    "description": "Perform tasks on specific dates",
    "status": "new",
    "recurrence": {
      "recurrence_type": "specific",
      "specific_dates": ["2024-04-15", "2024-12-25", "2024-12-31"]
    }
  }' | jq .

echo ""
echo ""

# 6. Получение списка всех задач
echo "6. Get all tasks"
curl -X GET $API_URL/tasks | jq .

echo ""
echo ""

# 7. Получение расписания для задачи
echo "7. Get schedule for task ID 1 (next 5 occurrences)"
curl -X GET "$API_URL/tasks/1/schedule?limit=5" | jq .

echo ""
echo ""

# 8. Получение расписания с расширенным лимитом
echo "8. Get schedule for task ID 2 (next 10 occurrences)"
curl -X GET "$API_URL/tasks/2/schedule?limit=10" | jq .

echo ""
echo ""

# 9. Получение задачи по ID
echo "9. Get task details by ID"
curl -X GET $API_URL/tasks/1 | jq .

echo ""
echo ""

# 10. Обновление задачи (изменение периодичности)
echo "10. Update task - change to every 2 days"
curl -X PUT $API_URL/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Daily patient review (updated)",
    "description": "Review all active patients - updated",
    "status": "in_progress",
    "recurrence": {
      "recurrence_type": "daily",
      "interval_days": 2
    }
  }' | jq .

echo ""
echo ""

# 11. Удаление задачи
echo "11. Delete task"
curl -X DELETE $API_URL/tasks/3

echo ""
echo ""

echo "=== End of examples ==="
