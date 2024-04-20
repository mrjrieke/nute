#!/bin/bash

# openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout mashup.key -out mashup.crt
openssl req -new -config san.cnf -keyout mashup.key -out mashup.crt

mkdir -p examples/helloworld/hellofyne/tls
cp mashup.key examples/helloworld/hellofyne/tls
cp mashup.crt examples/helloworld/hellofyne/tls

mkdir -p g3nd/worldg3n/tls
cp mashup.key g3nd/worldg3n/tls
cp mashup.crt g3nd/worldg3n/tls
