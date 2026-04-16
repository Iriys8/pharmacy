import express from 'express';
import { DataSource } from 'typeorm';
import { Redis } from 'ioredis';
import { randomGen } from '@pharmacy/src/shared/controllers/random_generator';
import { connectDB } from '@pharmacy/src/shared/setup/db_connect';
import { connectRedis } from '@pharmacy/src/shared/setup/redis_connect';
import { connectBroker } from '@pharmacy/src/shared/setup/broker_connect';
import { setupDB } from '@pharmacy/src/shared/setup/db_setup';
import { setupRoutes } from '@pharmacy/src/local-api/routes/routes';
import { setupBroker } from '@pharmacy/src/local-api/broker_setup';
import { Broker } from '@pharmacy/src/shared/models/models';

let dataSource: DataSource | null = null;
let redisDB: Redis | null = null;
let broker: Broker | null = null;

async function main(): Promise<void> {
  const name = randomGen();

  try {
    dataSource = await connectDB();
    console.log('Database connected');
    
    await setupDB(dataSource);
    console.log('Database initialized');

    redisDB = await connectRedis();
    console.log('Redis connected');

    broker = await connectBroker();
    console.log('RabbitMQ connected');

    await setupBroker(broker);
    console.log('RabbitMQ setup completed');

    const app = express();
    
    app.use(express.json());
    app.use(express.urlencoded({ extended: true }));
    
    setupRoutes(app, dataSource, redisDB, broker.channel);
    console.log('Routers initialized');

    const PORT = 8080;
    const server = app.listen(PORT, () => {
      console.log(`Local API server running on port ${PORT}`);
      console.log('Initialized');
    });

  } catch (err) {
    console.error('Failed to start server:', err);
    process.exit(1);
  }
}

main();