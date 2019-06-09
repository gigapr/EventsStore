FROM library/postgres

ENV POSTGRES_USER postgressuperuser
ENV POSTGRES_PASSWORD mysecretpassword
ENV POSTGRES_DB eventsStore

COPY sql/* /docker-entrypoint-initdb.d/

RUN ls /docker-entrypoint-initdb.d/
