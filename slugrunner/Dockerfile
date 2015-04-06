FROM flynn/cedarish

RUN apt-get update

# install deps
RUN apt-get install -y git golang

# clone logstash-forwarder
RUN git clone git://github.com/elasticsearch/logstash-forwarder.git /tmp/logstash-forwarder
RUN cd /tmp/logstash-forwarder && go build

RUN mkdir /opt/logstash-forwarder && cp /tmp/logstash-forwarder/logstash-forwarder /opt/logstash-forwarder/logstash-forwarder

RUN rm -rf /tmp/*

RUN mkdir -p /etc/logstash-forwarder && mkdir -p /opt/logstash-forwarder/certs

ADD logstash-forwarder/logstash-forwarder.json /etc/logstash-forwarder/logstash-forwarder.json

ADD logstash-forwarder/lumberjack.crt /opt/logstash-forwarder/certs/lumberjack.crt

ADD logstash-forwarder/lumberjack.key /opt/logstash-forwarder/certs/lumberjack.key

ADD ./runner/ /runner

ADD run_all.sh /opt/run_all.sh 

RUN chmod +x /opt/run_all.sh

ENTRYPOINT ["/opt/run_all.sh"]
