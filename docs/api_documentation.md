### Authentication
- POST /register
- POST /login (returns JWT)

### Authorization Rules
- First user → admin
- Admin-only: create/update/delete tasks, promote users
- Authenticated users: read tasks

### Header Format
Authorization: Bearer <token>
