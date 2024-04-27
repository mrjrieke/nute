#!/bin/bash

# openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout mashup.key -out mashup.crt
openssl req -new -nodes -newkey rsa:2048 -config san.cnf -reqexts v3_req -keyout mashup.key -out mashup.csr

openssl x509 -req -in mashup.csr -extfile san.cnf -extensions v3_req -signkey mashup.key -days 365 -out mashup.crt

mkdir -p examples/helloworld/hellofyne/tls
cp mashup.key examples/helloworld/hellofyne/tls
cp mashup.crt examples/helloworld/hellofyne/tls

mkdir -p g3nd/worldg3n/tls
cp mashup.key g3nd/worldg3n/tls
cp mashup.crt g3nd/worldg3n/tls
