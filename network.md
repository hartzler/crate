
#######
## NAT
#######

# geard example
iptables -t nat -A PREROUTING -d ${local_ip}/32 -p tcp -m tcp --dport ${local_port} -j DNAT --to-destination ${remote_ip}:${remote_port}
iptables -t nat -A OUTPUT -d ${local_ip}/32 -p tcp -m tcp --dport ${local_port} -j DNAT --to-destination ${remote_ip}:${remote_port}
iptables -t nat -A POSTROUTING -o eth0 -j SNAT --to-source ${container_ip}

# armada test
iptables -t nat -F
iptables -t nat -A OUTPUT -d localhost -p tcp --dport 8000 -j DNAT --to-destination 192.168.33.102:22
iptables -t nat -A POSTROUTING -j MASQUERADE
ssh localhost -p 8000


##########
## TUNNEL
##########

# bridge
ip link add name armada0 type bridge
ip addr add 10.4.0.1/8 dev aramda0
ip link set armada0 up

# setup tunnel to from node1 -> node2
ip tunnel add ${name} mode ipip local ${local} remote ${remote} ttl 64
ip addr add 10.4.0.${nodeNum}/8 dev ${name}
ip link set ${name} up
ip route add 10.4.0.0/16 dev ${name}

# setup namespace for component
ip netns add app.1

# setup veth pair for component
ip link add app.1.veth0 type veth peer name app.1.veth1
ip link set app.1.veth1 netns app.1

# route out components
