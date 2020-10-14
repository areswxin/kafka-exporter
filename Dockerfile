FROM  centos:latest

ADD kafka-exporter /bin/usr/sbin/kafka-exporter

CMD ["/bin/usr/sbin/kafka-exporter"]