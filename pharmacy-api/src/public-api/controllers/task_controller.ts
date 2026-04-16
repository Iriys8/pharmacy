import { Request, Response } from 'express';
import Redis from 'ioredis';
import { Channel } from 'amqplib';
import { randomBytes } from 'crypto';
import { RequestContext } from '@pharmacy/src/shared/models/models';

function randomGen(): string {
  return randomBytes(8).toString('hex');
}

export function makeTask(route: string, task: string, redisDB: Redis, channel: Channel) {
  return async (req: Request, res: Response): Promise<void> => {
    let taskContext: RequestContext = {
      query: {},
      context: {},
      claims: {}
    };

    let query: Record<string, any> = {};

    try {
      switch (task) {
        case "get":
          const id = req.query.id as string;
          if (id && isNaN(parseInt(id))) {
            res.status(400).json({
              error: "Invalid request body",
            });
            return;
          }
          
          query["id"] = id ? parseInt(id) : undefined;
          query["q"] = req.query.q || "";
          query["page"] = req.query.page || "";
          query["limit"] = req.query.limit || "";
          
          taskContext.query = query;
          break;
          
        case "schedule_dated":
          const start = req.query.start as string;
          const end = req.query.end as string;
          
          if (!start || !end) {
            res.status(400).json({ error: "Start and end dates are required" });
            return;
          }
          
          if (!/^\d{4}-\d{2}-\d{2}$/.test(start) || !/^\d{4}-\d{2}-\d{2}$/.test(end)) {
            res.status(400).json({ error: "Invalid date format" });
            return;
          }
          
          taskContext.query = { start, end };
          break;
          
        case "post":
          taskContext.context = req.body;
          break;
      }
      
      const taskID = randomGen();
      const key = `public_task:${taskID}`;
      
      await redisDB.hset(key, {
        status: "pending",
        task: task,
        route: route,
        result: "",
        context: JSON.stringify(taskContext),
      });
      
      await redisDB.expire(key, 300);
      
      const published = channel.publish(
        "public_exchange",
        route,
        Buffer.from(key),
        {
          contentType: "application/json",
          timestamp: Date.now(),
        }
      );
      
      if (!published) {
        console.error("Failed to publish message to RabbitMQ");
        res.status(500).json({
          error: "Internal error",
        });
        return;
      }
      
      res.status(200).json({
        TaskID: taskID,
      });
      
    } catch (err) {
      console.error("Error in makeTask:", err);
      res.status(500).json({
        error: "Internal error",
      });
    }
  };
}