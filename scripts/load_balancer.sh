#!/bin/bash



# Install Nginx

#sudo apt-get update
#sudo apt-get install nginx

echo " Install nginx plus manually -"
echo -n "Enter number of agent hosts to loadbalance [ENTER]: "
read number

if [ $number -eq $number 2>/dev/null ]
then



# Configure number of hosts =2 max

NGINX_FILE="/etc/nginx/sites-available/default"

echo "upstream webserver  {" > $NGINX_FILE

for (( c=1; c<=$number; c++ ))
do
	echo "Enter loadbalance agent host $c public ip - eg 152.1.13.183, 152.46.18.63"
	read ip
  echo "server $ip;">>$NGINX_FILE
done

echo "}

server {

  listen 80;

  location / {
    proxy_pass  http://webserver;
  }
}" >> $NGINX_FILE

echo "[INFO] Config saved - /etc/nginx/sites-available/default "

cat /etc/nginx/sites-available/default


echo "[INFO] Restart nginx"

service nginx restart


else
    echo "exit , number of agent hosts needs to be integer....$input is not an integer"
fi
