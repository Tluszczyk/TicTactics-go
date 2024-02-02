FROM mongo:latest

COPY init /docker-entrypoint-initdb.d
COPY init/data /tmp/data

ENV MONGO_INITDB_DATABASE=TicTactics

EXPOSE 27017
