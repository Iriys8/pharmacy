import { Express, Router } from 'express';
import { DataSource } from 'typeorm';
import { Redis } from 'ioredis';
import { Channel } from 'amqplib';
import cors from 'cors';
import { authMiddleware } from '@pharmacy/src/local-api/middleware/middleware';
import { login, logout, refreshToken } from '@pharmacy/src/local-api/controllers/authentication_controller';
import { pickup } from '@pharmacy/src/local-api/controllers/pickup_controller';
import { makeTask } from '@pharmacy/src/local-api/controllers/task_controller';
import { getImage } from '@pharmacy/src/shared/controllers/image_controller';
import { getLog, getLogs } from '@pharmacy/src/local-api/controllers/log_conroller';

function setupGoodsRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/goods', authMiddleware(dataSource, ''), makeTask('goods', 'get', redisDB, channel));
  apiGroup.patch('/goods', authMiddleware(dataSource, 'Update_Goods'), makeTask('goods', 'patch', redisDB, channel));
}

function setupScheduleRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/schedule', authMiddleware(dataSource, ''), makeTask('schedule', 'get', redisDB, channel));
  apiGroup.post('/schedule', authMiddleware(dataSource, 'Create_Schedule'), makeTask('schedule', 'post', redisDB, channel));
  apiGroup.patch('/schedule', authMiddleware(dataSource, 'Update_Schedule'), makeTask('schedule', 'patch', redisDB, channel));
  apiGroup.delete('/schedule', authMiddleware(dataSource, 'Delete_Schedule'), makeTask('schedule', 'delete', redisDB, channel));
}

function setupAnnounceRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/announce', authMiddleware(dataSource, ''), makeTask('announces', 'get', redisDB, channel));
  apiGroup.post('/announce', authMiddleware(dataSource, 'Create_Announces'), makeTask('announces', 'post', redisDB, channel));
  apiGroup.patch('/announce', authMiddleware(dataSource, 'Update_Announces'), makeTask('announces', 'patch', redisDB, channel));
  apiGroup.delete('/announce', authMiddleware(dataSource, 'Delete_Announces'), makeTask('announces', 'delete', redisDB, channel));
}

function setupOrderRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/order', authMiddleware(dataSource, 'Read_Orders'), makeTask('orders', 'get', redisDB, channel));
  apiGroup.post('/order', authMiddleware(dataSource, 'Create_Orders'), makeTask('orders', 'post', redisDB, channel));
  apiGroup.patch('/order', authMiddleware(dataSource, 'Update_Orders'), makeTask('orders', 'patch', redisDB, channel));
  apiGroup.delete('/order', authMiddleware(dataSource, 'Delete_Orders'), makeTask('orders', 'delete', redisDB, channel));
}

function setupUserRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/user', authMiddleware(dataSource, 'Change_Users'), makeTask('users', 'get', redisDB, channel));
  apiGroup.post('/user', authMiddleware(dataSource, 'Change_Users'), makeTask('users', 'post', redisDB, channel));
  apiGroup.patch('/user', authMiddleware(dataSource, 'Change_Users'), makeTask('users', 'patch', redisDB, channel));
  apiGroup.delete('/user', authMiddleware(dataSource, 'Change_Users'), makeTask('users', 'delete', redisDB, channel));
}

function setupRoleRouter(apiGroup: Router, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  apiGroup.get('/role', authMiddleware(dataSource, 'Change_Roles'), makeTask('roles', 'get', redisDB, channel));
  apiGroup.post('/role', authMiddleware(dataSource, 'Change_Roles'), makeTask('roles', 'post', redisDB, channel));
  apiGroup.patch('/role', authMiddleware(dataSource, 'Change_Roles'), makeTask('roles', 'patch', redisDB, channel));
  apiGroup.delete('/role', authMiddleware(dataSource, 'Change_Roles'), makeTask('roles', 'delete', redisDB, channel));
  apiGroup.get('/permission', authMiddleware(dataSource, 'Change_Roles'), makeTask('roles', 'permissions', redisDB, channel));
}

function setupOtherRouter(apiGroup: Router, dataSource: DataSource): void {
  apiGroup.get('/logs', authMiddleware(dataSource, 'Download_Logs'), getLogs);
  apiGroup.get('/log', authMiddleware(dataSource, 'Download_Logs'), getLog);
}

function setupPickupRouter(apiGroup: Router, redisDB: Redis ): void {
  apiGroup.get('/pickup', pickup(redisDB));
}

export function setupRoutes(app: Express, dataSource: DataSource, redisDB: Redis, channel: Channel): void {
  app.use(cors({
    origin: ['http://localhost:5001', 'http://127.0.0.1:5001', 'http://localhost:5174', 'http://localhost:5173'],
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'],
    allowedHeaders: ['Origin', 'Content-Length', 'Content-Type', 'Authorization'],
    credentials: true,
    maxAge: 8 * 3600,
  }));

  const apiGroup = Router();

  apiGroup.post('/login', login(dataSource));
  apiGroup.post('/refresh', refreshToken(dataSource));
  apiGroup.post('/logout', logout());
  apiGroup.get('/image', getImage);

  setupGoodsRouter(apiGroup, dataSource, redisDB, channel);
  setupScheduleRouter(apiGroup, dataSource, redisDB, channel);
  setupAnnounceRouter(apiGroup, dataSource, redisDB, channel);
  setupOrderRouter(apiGroup, dataSource, redisDB, channel);
  setupUserRouter(apiGroup, dataSource, redisDB, channel);
  setupRoleRouter(apiGroup, dataSource, redisDB, channel);
  setupOtherRouter(apiGroup, dataSource);
  setupPickupRouter(apiGroup, redisDB);

  app.use('/api', apiGroup);
}