FROM alpine

COPY bin/api /bin/velatemplatetesterapi

ENTRYPOINT ["/bin/velatemplatetesterapi"]