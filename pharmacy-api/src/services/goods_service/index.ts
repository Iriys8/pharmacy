import { DataSource } from 'typeorm';
import { Redis } from 'ioredis';
import { Channel, Message } from 'amqplib';
import { randomGen } from '@pharmacy/src/shared/controllers/random_generator';
import { connectDB } from '@pharmacy/src/shared/setup/db_connect';
import { connectRedis, closeRedis } from '@pharmacy/src/shared/setup/redis_connect';
import { connectBroker } from '@pharmacy/src/shared/setup/broker_connect';
import { getGoods, getGoodsByID, getPromoItems, updateGoods } from './controller/controller';
import { GoodsUpdateRequest } from '@pharmacy/src/shared/models/responses';
import { Claims, Broker } from '@pharmacy/src/shared/models/models';
import * as fs from 'fs';

let logFile: fs.WriteStream | null = null;

async function main(): Promise<void> {
  const uname_public = randomGen();
  const uname_local = randomGen();

  console.log(`Service public name: ${uname_public}, local name: ${uname_local}`);

  let dataSource: DataSource | null = null;
  let redisDB: Redis | null = null;
  let broker: Broker = null;
  try {

    dataSource = await connectDB();
    console.log('Database connected');

    redisDB = await connectRedis();
    console.log('Redis connected');

    broker = await connectBroker();
    console.log('RabbitMQ connected');

    if (broker.channel) {
      consumeMessages(broker.channel, 'public_goods_queue', redisDB, dataSource, uname_public);
      consumeMessages(broker.channel, 'local_goods_queue', redisDB, dataSource, uname_local);
    }

    console.log('Service is running.');

  } catch (err) {
    process.exit(1);
  }
}

async function consumeMessages(
  ch: Channel,
  queueName: string,
  redisDB: Redis,
  dataSource: DataSource,
  consumerName: string
): Promise<void> {
  try {
    await ch.consume(queueName, async (msg: Message | null) => {
      if (!msg) return;

      try {
        const taskKey = msg.content.toString();
        console.log(`Received task: ${taskKey}`);

        const taskData = await redisDB.hgetall(taskKey);

        if (!taskData || Object.keys(taskData).length === 0) {
          console.log(`Task ${taskKey} not found in Redis`);
          ch.ack(msg);
          return;
        }

        if (taskData.status !== 'pending') {
          console.log(`Task ${taskKey} is not pending (status: ${taskData.status}), skipping`);
          ch.ack(msg);
          return;
        }

        let result: any = {};
        let execError: Error | null = null;

        try {
          switch (taskData.task) {
            case 'get': {
              const taskContext = JSON.parse(taskData.context);
              const query = taskContext.Query || taskContext.query || {};
              
              if (!query.id || query.id === 0) {
                result = await getGoods(
                  dataSource,
                  query.q || '',
                  query.page || '1',
                  query.limit || '10'
                );
              } else {
                result = await getGoodsByID(dataSource, query.id);
              }
              break;
            }

            case 'advert': {
              result = await getPromoItems(dataSource);
              break;
            }

            case 'patch': {
              const taskContext = JSON.parse(taskData.context);
              const query = taskContext.Query || taskContext.query || {};
              const context = taskContext.Context || taskContext.context;
              const claims = taskContext.Claims || taskContext.claims;
              
              const updateResult = await updateGoods(
                dataSource,
                query.id,
                context as GoodsUpdateRequest,
                claims as Claims
              );
              result = { Response: updateResult };
              break;
            }

            default: {
              console.log(`Unknown task type: ${taskData.task}`);
              execError = new Error('Unknown task');
            }
          }
        } catch (err) {
          execError = err as Error;
          console.error(`Error:`, err)
        }

        let jsonResult: string;
        if (execError) {
          taskData.status = 'error';
          jsonResult = JSON.stringify({ error: execError.message });
        } else {
          taskData.status = 'completed';
          const responseData = result.Response !== undefined ? result.Response : result;
          jsonResult = JSON.stringify(responseData);
        }

        taskData.result = jsonResult;

        await redisDB.hset(taskKey, taskData);
        await redisDB.expire(taskKey, 20);

        console.log(`Task ${taskKey} completed successfully with status: ${taskData.status}`);
        ch.ack(msg);

      } catch (err) {
        console.error(`Error processing message:`, err);
        ch.ack(msg);
      }
    });
  } catch (err) {
    console.error(`Failed to register consumer for ${queueName}:`, err);
    throw err;
  }
}

main();