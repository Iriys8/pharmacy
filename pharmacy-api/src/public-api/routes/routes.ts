import { Router, Express } from 'express';
import cors from 'cors';
import Redis from 'ioredis';
import { Channel } from 'amqplib';
import { makeTask } from '@pharmacy/src/public-api/controllers/task_controller';
import { pickup } from '@pharmacy/src/public-api/controllers/pickup_controller';
import { getImage } from '@pharmacy/src/shared/controllers/image_controller';

function setupGoodsRouter(apiGroup: Router, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/goods', makeTask('goods', 'get', redisDB, channel));
  apiGroup.get('/goods/advert', makeTask('goods', 'advert', redisDB, channel));
}

function setupScheduleRouter(apiGroup: Router, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/schedule', makeTask('schedule', 'schedule_dated', redisDB, channel));
}

function setupOrderRouter(apiGroup: Router, redisDB: Redis, channel: Channel): void {
  apiGroup.post('/order', makeTask('orders', 'post', redisDB, channel));
}

function setupAnnounceRouter(apiGroup: Router, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/announces', makeTask('announces', 'get', redisDB, channel));
}

function setupPickupRouter(apiGroup: Router, redisDB: Redis): void {
  apiGroup.get('/pickup', pickup(redisDB));
}

function setupOtherRouter(apiGroup: Router): void {
  apiGroup.get('/image', getImage)
}

export function setupRoutes(app: Express, redisDB: Redis, channel: Channel): void {

  app.use(cors({
    origin: ['http://localhost:5000', 'http://localhost:5000', 'http://localhost:5173'],
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'],
    allowedHeaders: ['Origin', 'Content-Length', 'Content-Type', 'Authorization'],
    credentials: true,
    maxAge: 8 * 3600,
  }));

  const apiGroup = Router();
  
  setupGoodsRouter(apiGroup, redisDB, channel);
  setupScheduleRouter(apiGroup, redisDB, channel);
  setupOrderRouter(apiGroup, redisDB, channel);
  setupAnnounceRouter(apiGroup, redisDB, channel);
  setupPickupRouter(apiGroup, redisDB);
  setupOtherRouter(apiGroup)

  app.use('/api', apiGroup);
}