FROM amacneil/dbmate:2

ADD migrations /opt/migrations
ENV DBMATE_MIGRATIONS_DIR=/opt/migrations
ENV DBMATE_NO_DUMP_SCHEMA=true

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Use our custom entrypoint.
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
