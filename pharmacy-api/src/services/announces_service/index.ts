import { DataSource } from 'typeorm';
import { Redis } from 'ioredis';
import { Channel, Message } from 'amqplib';
import { randomGen } from '@pharmacy/src/shared/controllers/random_generator';
import { connectDB } from '@pharmacy/src/shared/setup/db_connect';
import { connectRedis, closeRedis } from '@pharmacy/src/shared/setup/redis_connect';
import { connectBroker } from '@pharmacy/src/shared/setup/broker_connect';
import { createAnnounce, deleteAnnounce, getAnnounceByID, getAnnounces, updateAnnounce } from '@pharmacy/src/services/announces_service/controller';
import { AnnouncementResponse } from '@pharmacy/src/shared/models/responses';
import { Claims, Broker } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'
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
    Logger.info('Announces_service: Database connected');

    redisDB = await connectRedis();
    Logger.info('Announces_service: Redis connected');

    broker = await connectBroker();
    Logger.info('Announces_service: RabbitMQ connected');

    if (broker.channel) {
      consumeMessages(broker.channel, 'public_announces_queue', redisDB, dataSource, uname_public);
      consumeMessages(broker.channel, 'local_announces_queue', redisDB, dataSource, uname_local);
    }

    Logger.info('Announces_service: Service is running.');

  } catch (err) {
    process.exit(1);
  }
}

async function consumeMessages(ch: Channel, queueName: string, redisDB: Redis, dataSource: DataSource, consumerName: string): Promise<void> {
  try {
    await ch.consume(queueName, async (msg: Message | null) => {
      if (!msg) return;

      try {
        const taskKey = msg.content.toString();
        Logger.info(`Announces_service: Received task: ${taskKey}`);

        const taskData = await redisDB.hgetall(taskKey);

        if (!taskData || Object.keys(taskData).length === 0) {
          Logger.info(`Announces_service: Task ${taskKey} not found in Redis`);
          ch.ack(msg);
          return;
        }

        if (taskData.status !== 'pending') {
          Logger.info(`Announces_service: Task ${taskKey} is not pending (status: ${taskData.status}), skipping`);
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
              const claims = taskContext.Claims || taskContext.claims;
              
              if (!query.id || query.id === 0) {
                result = await getAnnounces(dataSource, query.q || '', query.page || '1', query.limit || '10',);
              } else {
                result = await getAnnounceByID(dataSource, query.id, claims);
              }
              break;
            }

            case 'post': {
              const taskContext = JSON.parse(taskData.context);
              const context = taskContext.Context || taskContext.context;
              const claims = taskContext.Claims || taskContext.claims;
              
              const createResult = await createAnnounce(dataSource, context as AnnouncementResponse, claims as Claims);
              result = { Response: createResult };
              break;
            }

            case 'patch': {
              const taskContext = JSON.parse(taskData.context);
              const query = taskContext.Query || taskContext.query || {};
              const context = taskContext.Context || taskContext.context;
              const claims = taskContext.Claims || taskContext.claims;
              
              const updateResult = await updateAnnounce(dataSource,query.id, context as AnnouncementResponse, claims as Claims);
              result = { Response: updateResult };
              break;
            }

            case 'delete': {
              const taskContext = JSON.parse(taskData.context);
              const query = taskContext.Query || taskContext.query || {};
              const claims = taskContext.Claims || taskContext.claims;
              
              const deleteResult = await deleteAnnounce(dataSource, query.id, claims as Claims);
              result = { Response: deleteResult };
              break;
            }

            default: {
              Logger.info(`Announces_service: Unknown task type: ${taskData.task}`);
              execError = new Error('Unknown task');
            }
          }
        } catch (err) {
          execError = err as Error;
          Logger.error(`Announces_service: Error:`, err)
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

        Logger.info(`Announces_service: Task ${taskKey} completed successfully with status: ${taskData.status}`);
        ch.ack(msg);

      } catch (err) {
        Logger.error(`Announces_service: Error processing message:`, err);
        ch.ack(msg);
      }
    });
  } catch (err) {
    Logger.error(`Announces_service: Failed to register consumer for ${queueName}:`, err);
    throw err;
  }
}

main();