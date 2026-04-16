import express from 'express';
import { connectRedis } from '@pharmacy/src/shared/setup/redis_connect';
import { connectBroker } from '@pharmacy/src/shared/setup/broker_connect';
import { setupRoutes } from '@pharmacy/src/public-api/routes/routes';
import { setupBroker } from '@pharmacy/src/public-api/broker_setup';

async function main() {
  let redisDB = null;
  let broker = null;

  try {
    redisDB = await connectRedis();
    console.log('Redis connected successfully');

    broker = await connectBroker();
    console.log('RabbitMQ connected successfully');

    await setupBroker(broker);
    console.log('RabbitMQ setup completed');

    const app = express();
    
    app.use(express.json());
    app.use(express.urlencoded({ extended: true }));

    setupRoutes(app, redisDB, broker.channel);

    const PORT = process.env.PUBLIC_API_PORT || 8080;
    app.listen(PORT, () => {
      console.log(`Public API server running on port ${PORT}`);
    });

    const shutdown = async () => {
      console.log('Shutting down gracefully...');
      
      if (broker.channel) {
        await broker.channel.close();
      }
      if (broker.connection) {
        await broker.connection.close();
      }
      if (redisDB) {
        await redisDB.quit();
      }
      
      console.log('All connections closed');
      process.exit(0);
    };

    process.on('SIGINT', shutdown);
    process.on('SIGTERM', shutdown);

  } catch (err) {
    console.error('Failed to start server:', err);
    
    if (broker.channel) await broker.channel.close().catch(console.error);
    if (broker.connection) await broker.connection.close().catch(console.error);
    if (redisDB) await redisDB.quit().catch(console.error);
    
    process.exit(1);
  }
}

main();