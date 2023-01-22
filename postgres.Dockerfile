FROM postgres:15

RUN echo "CREATE EXTENSION pg_stat_statements;" >> /docker-entrypoint-initdb.d/init.sql