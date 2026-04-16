FROM oven/bun:1.2
WORKDIR /public-api

COPY ./pharmacy-api/package.json ./pharmacy-api/bun.lock ./pharmacy-api/tsconfig.json ./

COPY ./pharmacy-api/src/public-api ./src/public-api
COPY ./pharmacy-api/src/shared ./src/shared
COPY ./pharmacy-api/src/migration ./src/migration

RUN bun install
RUN bun add -g typescript
RUN tsc

CMD ["node", "/public-api/dist/src/public-api/index.js"] 