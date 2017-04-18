#!/bin/bash


# default interfaces
EXTERNAL_INTERFACE="eth1"
INTERNAL_INTERFACE="docker0"

#default port

EXTERNAL_PORT="3306"
INTERNAL_PORT="3306"

#DOCKERIP

DOCKER_IP="172.17.0.2"

while [[ $# -gt 1 ]]
do
key="$1"

case $key in
    -e|--externalinterface)
    EXTERNAL_INTERFACE="$2"
    shift # past argument
    ;;

    -i|--internalinterface)
    INTERNAL_INTERFACE="$2"
    shift # past argument
    ;;

    -p|--internalport)
    INTERNAL_PORT="$2"
    shift # past argument
    ;;

    -x|--externalport)
    EXTERNAL_PORT="$2"
    shift # past argument
    ;;

    -d|--dockerip)
    DOCKER_IP="$2"
    shift # past argument
    ;;

    --default)
    echo "Default Do nothing"
    ;;
    *)
            # unknown option

    ;;
esac
shift # past argument or value
done



#Finding IP address of interfaces
EXTERNAL_IP=$(ip addr | grep inet | grep $EXTERNAL_INTERFACE | awk -F" " '{print $2}'| sed -e 's/\/.*$//')
INTERNAL_IP=$(ip addr | grep inet | grep $INTERNAL_INTERFACE | awk -F" " '{print $2}'| sed -e 's/\/.*$//')

#Print

echo "External Interface - $EXTERNAL_INTERFACE"
echo "Internal Interface - $INTERNAL_INTERFACE"

echo "External Port - $EXTERNAL_PORT"
echo "Internal Port - $INTERNAL_PORT"

echo "External IP - $EXTERNAL_IP"
echo "Internal IP - $INTERNAL_IP"

echo "Docker IP - $DOCKER_IP"




#Forward Status
echo "Port Forward Status (1 if enabled)"
cat /proc/sys/net/ipv4/ip_forward

#echo "sudo iptables -A FORWARD -i $EXTERNAL_INTERFACE -o INTERNAL_INTERFACE -p tcp --syn --dport $EXTERNAL_PORT -m conntrack --ctstate NEW -j ACCEPT"

echo "sudo iptables -A FORWARD -i $EXTERNAL_INTERFACE -o $INTERNAL_INTERFACE -p tcp --syn --dport $EXTERNAL_PORT -m conntrack --ctstate NEW -j ACCEPT"
sudo iptables -A FORWARD -i $EXTERNAL_INTERFACE -o $INTERNAL_INTERFACE -p tcp --syn --dport $EXTERNAL_PORT -m conntrack --ctstate NEW -j ACCEPT

echo "sudo iptables -A FORWARD -i $EXTERNAL_INTERFACE -o $INTERNAL_INTERFACE -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT"
sudo iptables -A FORWARD -i $EXTERNAL_INTERFACE -o $INTERNAL_INTERFACE -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT


echo "sudo iptables -A FORWARD -i $INTERNAL_INTERFACE -o $EXTERNAL_INTERFACE -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT"
sudo iptables -A FORWARD -i $INTERNAL_INTERFACE -o $EXTERNAL_INTERFACE -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT


echo "sudo iptables -t nat -A PREROUTING -i $EXTERNAL_INTERFACE -p tcp --dport $INTERNAL_PORT -j DNAT --to-destination $DOCKER_IP"
sudo iptables -t nat -A PREROUTING -i $EXTERNAL_INTERFACE -p tcp --dport $INTERNAL_PORT -j DNAT --to-destination $DOCKER_IP

echo "sudo iptables -t nat -A POSTROUTING -o $INTERNAL_INTERFACE -p tcp --dport $INTERNAL_PORT -d $DOCKER_IP -j SNAT --to-source $INTERNAL_IP"
sudo iptables -t nat -A POSTROUTING -o $INTERNAL_INTERFACE -p tcp --dport $INTERNAL_PORT -d $DOCKER_IP -j SNAT --to-source $INTERNAL_IP
