# TempAPI
## Requirements
* Ubuntu 18.04 LTS

## Database Setup
1. Install postgresql.
```
apt install postgresql
```
2. Setup user and database. Replace "$PASSWORD" with the password.
```
sudo -u postgres psql -d template1 -c "CREATE USER tempapi CREATEDB PASSWORD '$PASSWORD';"
sudo -u postgres psql -d template1 -c "CREATE DATABASE tempapi OWNER tempapi;"
```  
3. Update the configuration with the address, user, name and password for the database. Remember to reconfigure PostgreSQL for remote access if needed.
