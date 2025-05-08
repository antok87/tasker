FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY tasker /bin/tasker
ENTRYPOINT [ "/bin/tasker" ]
