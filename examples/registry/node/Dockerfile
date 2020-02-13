FROM node:12-alpine as base

WORKDIR /app
COPY . .

RUN npm install

FROM base AS test
ENTRYPOINT ["/bin/sh", "-c", "npm run test"]

FROM base AS build

EXPOSE 8000

USER node
CMD ["node", "index.js"]
