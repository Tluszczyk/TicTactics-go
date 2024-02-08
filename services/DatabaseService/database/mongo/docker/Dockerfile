FROM mongo:latest

COPY init /docker-entrypoint-initdb.d
COPY init/data /tmp/data

ENV MONGO_INITDB_DATABASE=TicTactics
ENV SESSION_TTL_SECONDS=18000

EXPOSE 27017
