```
curlie post localhost:8080/api/auth/register name=ali phone_number=09113334433 password=12341233
curlie post localhost:8080/api/auth/login phone_number=09113334433 password=12341233
curlie get localhost:8080/api/profile Authorization:"Bearer TOKEN"
curlie put localhost:8080/api/profile Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODEyMjMxOTAsImlhdCI6MTc3NzYyMzE5MCwidXNlcl9pZCI6MjZ9.7aP_E4g4h8O6r4GlhAGpyxGwhGCIjwqt6h0lj9zpbdA" name="Reza Yusufy Developer"
```