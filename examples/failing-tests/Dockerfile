FROM alpine:latest as base

FROM base AS test
ENTRYPOINT ["/bin/sh", "-c"]
CMD [ "false" ]

FROM base AS release
ENTRYPOINT ["/bin/sh", "-c"]
CMD [ "true" ]
