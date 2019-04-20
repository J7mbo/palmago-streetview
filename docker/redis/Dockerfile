FROM redis:3.0.7

# Defaults to 6379 but is set from the env var.
ARG REDIS_PORT=6379
ENV REDIS_PORT=$REDIS_PORT

COPY ./redis/redis.conf  /etc/redis/redis.conf
COPY /redis/run /usr/local/bin/run

RUN ["chmod", "+x", "/usr/local/bin/run"]

CMD ["/usr/local/bin/run"]