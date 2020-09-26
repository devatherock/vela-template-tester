FROM alpine

COPY velatemplatetesterapi /bin/velatemplatetesterapi

ENTRYPOINT ["/bin/velatemplatetesterapi"]