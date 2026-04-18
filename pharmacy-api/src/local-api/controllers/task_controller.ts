import { Request, Response } from 'express';
import Redis from 'ioredis';
import { Channel } from 'amqplib';
import { randomBytes } from 'crypto';
import { RequestContext, Claims } from '@pharmacy/src/shared/models/models';

function randomGen(): string {
  return randomBytes(8).toString('hex');
}

declare global {
  namespace Express {
    interface Request {
      user?: Claims;
    }
  }
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
      const claims = req.user;
      
      if (claims) {
        taskContext.claims = claims;
      }

      switch (task) {
        case "get":
          const idParam = req.query.id as string;
          
          if (route !== "users") {
            const idInt = parseInt(idParam);
            if (isNaN(idInt) && idParam) {
              res.status(400).json({
                error: "Invalid request body",
              });
              return;
            }
            query["id"] = isNaN(idInt) ? undefined : idInt;
          } else {
            query["id"] = idParam || undefined;
          }
          
          query["q"] = req.query.q || "";
          query["page"] = req.query.page || "";
          query["limit"] = req.query.limit || "";
          
          taskContext.query = query;
          break;
          
        case "post":
          taskContext.context = req.body;
          break;
          
        case "patch":
          const patchId = parseInt(req.query.id as string);
          if (isNaN(patchId) && req.query.id) {
            res.status(400).json({
              error: "Invalid request body",
            });
            return;
          }
          query["id"] = isNaN(patchId) ? undefined : patchId;
          taskContext.query = query;
          taskContext.context = req.body;
          break;
          
        case "delete":
          const deleteId = parseInt(req.query.id as string);
          if (isNaN(deleteId) && req.query.id) {
            res.status(400).json({
              error: "Invalid request body",
            });
            return;
          }
          query["id"] = isNaN(deleteId) ? undefined : deleteId;
          taskContext.query = query;
          break;

        case "permissions":
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
      }
      
      const taskID = randomGen();
      const key = `local_task:${taskID}`;
      
      await redisDB.hset(key, {
        status: "pending",
        task: task,
        route: route,
        result: "",
        context: JSON.stringify(taskContext),
      });
      
      await redisDB.expire(key, 300);
      
      const published = channel.publish(
        "local_exchange",
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