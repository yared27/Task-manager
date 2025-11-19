Task Manager API Documentation

Base URL:

http://localhost:8080

Endpoints Overview
Method	Endpoint	Description
GET	/tasks/	Get all tasks
GET	/tasks/:id	Get a task by ID
POST	/tasks/	Create a new task
PUT	/tasks/:id	Update a task
DELETE	/tasks/:id	Delete a task
1. GET /tasks/
Description

Fetch all tasks stored in memory.

Response 200 OK
{
  "tasks": [
    {
      "id": 1,
      "title": "Example Task",
      "description": "This is a sample task",
      "due_date": "2025-12-31T23:59:59Z",
      "status": "pending"
    }
  ]
}

2. GET /tasks/:id
Description

Retrieve a single task by its ID.

Example Request
GET http://localhost:8080/tasks/1

Response 200 OK
{
  "id": 1,
  "title": "Example Task",
  "description": "This is a sample task",
  "due_date": "2025-12-31T23:59:59Z",
  "status": "pending"
}

Response 400 (Invalid ID)
{
  "error": "invalid id"
}

Response 404 (Not Found)
{
  "error": "task not found"
}

3. POST /tasks/
Description

Create a new task.

Request Body (JSON)
{
  "title": "Do laundry",
  "description": "Wash and fold clothes",
  "due_date": "2025-11-30T18:00:00Z",
  "status": "pending"
}

Response 201 Created
{
  "id": 1,
  "title": "Do laundry",
  "description": "Wash and fold clothes",
  "due_date": "2025-11-30T18:00:00Z",
  "status": "pending"
}

Response 400 (Invalid Input)
{
  "error": "invalid input: ..."
}

4. PUT /tasks/:id
Description

Update a specific task.

Request Body
{
  "title": "Updated Task",
  "description": "Updated description",
  "due_date": "2025-12-02T12:00:00Z",
  "status": "in-progress"
}

Response 200 OK
{
  "id": 1,
  "title": "Updated Task",
  "description": "Updated description",
  "due_date": "2025-12-02T12:00:00Z",
  "status": "in-progress"
}

Response 404 (Not Found)
{
  "error": "task not found"
}

5. DELETE /tasks/:id
Description

Delete a task by ID.

Response 204 No Content
(no body)

Response 404 Not Found
{
  "error": "task not found"
}

Testing with CURL
Create a task:
curl -X POST http://localhost:8080/tasks/ \
-H "Content-Type: application/json" \
-d '{
  "title": "Test Task",
  "description": "Example",
  "due_date": "2025-12-01T12:00:00Z",
  "status": "pending"
}'

Get all tasks:
curl http://localhost:8080/tasks/

Update task:
curl -X PUT http://localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-d '{
  "title": "Updated",
  "description": "Updated desc",
  "due_date": "2025-12-02T12:00:00Z",
  "status": "done"
}'

Delete task:
curl -X DELETE http://localhost:8080/tasks/1