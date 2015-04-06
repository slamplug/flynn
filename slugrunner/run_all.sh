#!/bin/bash
(/runner/init "$@" &) && \
(/opt/logstash-forwarder/logstash-forwarder -config /etc/logstash-forwarder/logstash-forwarder.json -spool-size 100)
