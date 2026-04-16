import { Request, Response } from 'express';
import Redis from 'ioredis';

export function pickup(redisDB: Redis) {
  return async (req: Request, res: Response): Promise<void> => {
    const key = req.query.key as string;

    if (!key) {
      res.status(400).json({
        Error: "key is required",
      });
      return;
    }

    try {
      const val = await redisDB.hgetall(`local_task:${key}`);

      if (!val || Object.keys(val).length === 0) {
        res.status(404).json({
          Status: "not_found",
        });
        return;
      }

      switch (val.status) {
        case "completed":
          // Удаляем ключ после получения результата
          await redisDB.del(`local_task:${key}`).catch((err) => {
            console.error(`Error deleting key local_task:${key}:`, err);
          });
          
          res.status(200).json({
            Status: val.status,
            Value: val.result,
          });
          return;
          
        case "pending":
          res.status(200).json({
            Status: val.status,
          });
          return;
          
        case "error":
          res.status(200).json({
            Status: val.status,
          });
          return;
          
        default:
          res.status(200).json({
            Status: val.status,
          });
          return;
      }
    } catch (err) {
      console.error("Error in pickup controller:", err);
      res.status(500).json({
        Error: "Internal error",
      });
      return;
    }
  };
}