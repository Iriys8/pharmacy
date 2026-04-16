import Redis from 'ioredis';

export async function connectRedis(): Promise<Redis> {
  const redisClient = new Redis({
    host: process.env.REDIS_HOST,
    port: parseInt(process.env.REDIS_PORT || '6379'),
    password: process.env.REDIS_PASSWORD,
    db: 0,
  });

  try {
    await redisClient.ping();
    console.log('Connected to Redis successfully');
    return redisClient;
  } catch (error) {
    console.error('Failed to connect to Redis:', error);
    throw error;
  }
}

export async function closeRedis(redisClient: Redis): Promise<void> {
  try {
    await redisClient.quit();
    console.log('Redis connection closed');
  } catch (error) {
    console.error('Error while closing Redis connection:', error);
    throw error;
  }
}