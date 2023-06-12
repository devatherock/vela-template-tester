FROM alpine

COPY bin/app /bin/velatemplatetesterapi

ENTRYPOINT ["/bin/velatemplatetesterapi"]