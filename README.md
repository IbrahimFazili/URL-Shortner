# README

## How to run

To run, make sure the main system has Docker and each of the systems has been logged into Docker.

Restart on the main vm node
```
sudo shutdown -r now
```

Then run in the root level of the folder
```
python3 start_script.py
```

All curl requests must be done on another terminal

To execute Q1, without the writer service so the writes are done by the url-shortner, in the `docker-compose.yml` file
for the environment variable under the services `web`, change the the value of `EXTERNAL_WRITER` to false 

To execute Q2, leave that variables value to true

## Machines assumptions
I assumed that this will be running on my set of vm's, 10.11.*.107. To run on your own vms (that have the password hhhhiotwwg),
change the values of `CASSANDRA_CONNECT_POINT` in the `docker-compose.yml` and the values in the `hosts` variables at the top of `start_script.py`