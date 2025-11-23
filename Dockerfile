FROM debian:latest

RUN apt update
RUN apt install -y iproute2 iptables nftables tcpdump iputils-ping sudo vim curl
RUN apt install -y net-tools

#RUN mkdir -p /dev/net
#RUN mknod /dev/net/tun c 10 200
#RUN chmod 0666 /dev/net/tun

# Включаем форвардинг
#RUN echo 1 > /proc/sys/net/ipv4/ip_forward

CMD ["/bin/bash"]