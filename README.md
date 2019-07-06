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

## Use TempAPI
Consider this example URL:
```
https://10.125.135.142:8443/?room=fifthroom&time=792374&temp=63.3&token=kfkkofkdokfdo
```
* `room` = The room where the measurement was taken
* `time` = The time when the measurement was taken as a UNIX timestamp   
* `temp` = The temperature
* `token` = The specified token in the server

Note that the room will be created automatically when it doesn't exist.
