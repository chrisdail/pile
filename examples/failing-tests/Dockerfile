FROM alpine:latest as base

FROM base AS test

ADD failure.xml /app/build/test-results/

ENTRYPOINT ["/bin/sh", "-c"]
CMD [ "false" ]

FROM base AS release
ENTRYPOINT ["/bin/sh", "-c"]
CMD [ "true" ]
