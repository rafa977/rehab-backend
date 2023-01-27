# rehab-backend


# installing postgresql
docker run --name postgresql -e POSTGRES_USER=myusername -e POSTGRES_PASSWORD=mypassword -p 5432:5432 -v /data:/var/lib/postgresql/data -d postgres


# roles

1. admin
2. patient
3. doctor
4. employee
