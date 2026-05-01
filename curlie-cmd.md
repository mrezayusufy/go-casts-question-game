```
curlie post localhost:8080/api/auth/register name=ali phone_number=09113334433 password=12341233
``` 
```
curlie post localhost:8080/api/auth/login phone_number=09113334433 password=12341233
```
```
curlie get localhost:8080/api/profile Authorization:"Bearer TOKEN"
```
```
curlie put localhost:8080/api/profile Authorization:"Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODEyNTU5NTUsImlhdCI6MTc3NzY1NTk1NSwidXNlcl9pZCI6Mjh9.ySbSypr4EG9iUKLvPyg8WPfbF7oN1pnYkEtgSfAE19k" name="Reza Yusufy Developer"
```
```
curlie post localhost:8080/api/change-password Authorization:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODEyNTU5NTUsImlhdCI6MTc3NzY1NTk1NSwidXNlcl9pZCI6Mjh9.ySbSypr4EG9iUKLvPyg8WPfbF7oN1pnYkEtgSfAE19k" old_password=12341233 new_password=3232absd@$
```