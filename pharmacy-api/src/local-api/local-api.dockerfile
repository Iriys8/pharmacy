FROM oven/bun:1.2
WORKDIR /app

COPY ./pharmacy-api/package.json ./pharmacy-api/bun.lock ./pharmacy-api/tsconfig.json ./

COPY ./pharmacy-api/src/local-api ./src/local-api
COPY ./pharmacy-api/src/shared ./src/shared
COPY ./pharmacy-api/src/migration ./src/migration

RUN bun install
RUN bun add -g typescript
RUN tsc

CMD ["node", "/app/dist/src/local-api/index.js"] 