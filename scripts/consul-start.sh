docker run -d \
  --name=consul0 \
  -p 8500:8500 \
  -p 8600:8600/udp \
  hashicorp/consul agent -dev -client=0.0.0.0