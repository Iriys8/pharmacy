import express from 'express';
import { DataSource } from 'typeorm';
import { Redis } from 'ioredis';
import { connectDB } from '@pharmacy/src/shared/setup/db_connect';
import { connectRedis } from '@pharmacy/src/shared/setup/redis_connect';
import { connectBroker } from '@pharmacy/src/shared/setup/broker_connect';
import { setupDB } from '@pharmacy/src/shared/setup/db_setup';
import { setupRoutes } from '@pharmacy/src/local-api/routes/routes';
import { setupBroker } from '@pharmacy/src/local-api/broker_setup';
import { Broker } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'
import cookieParser from 'cookie-parser';

let dataSource: DataSource | null = null;
let redisDB: Redis | null = null;
let broker: Broker | null = null;

async function main(): Promise<void> {
  try {
    dataSource = await connectDB();
    Logger.info(`Local-api: Database connected`);
    
    await setupDB(dataSource);
    Logger.info(`Local-api: Database initialized`);

    redisDB = await connectRedis();
    Logger.info(`Local-api: Redis connected`);

    broker = await connectBroker();
    Logger.info(`Local-api: RabbitMQ connected`);

    await setupBroker(broker);
    Logger.info(`Local-api: RabbitMQ setup completed`);

    const app = express();
    
    app.use(express.json());
    app.use(express.urlencoded({ extended: true }));
    app.use(cookieParser())

    setupRoutes(app, dataSource, redisDB, broker.channel);
    Logger.info(`Local-api: Routers initialized`);

    const PORT = 8080;
    app.listen(PORT, () => {
      Logger.info(`Local-api: Local API server running on port ${PORT}`);
      Logger.info(`Local-api: Initialized`);
    });
  } catch (err) {
    Logger.error(`Local-api: Failed to start server:`, err);
    process.exit(1);
  }
}

main();