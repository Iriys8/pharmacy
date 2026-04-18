import { Request, Response } from 'express';
import fs from 'fs';
import path from 'path';
import { Claims } from '@pharmacy/src/shared/models/models';
import { LogsResponse } from '@pharmacy/src/shared/models/responses';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'

const LOGS_FOLDER = './pharmacy-content/logs';

export function getLogs(req: Request, res: Response): void {
  const query = (req.query.q as string) || '';
  const pageStr = req.query.page as string;
  const limitStr = req.query.limit as string;
  
  const claims = req.user as Claims;

  let limit = parseInt(limitStr);
  if (isNaN(limit) || limit < 1) {
    limit = 10;
  } else if (limit > 40) {
    limit = 40;
  }

  let page = parseInt(pageStr);
  if (isNaN(page) || page < 1) {
    page = 1;
  }
  const offset = (page - 1) * limit;

  try {
    const files = fs.readdirSync(LOGS_FOLDER);
    
    let allLogs: LogsResponse[] = [];
    
    for (const file of files) {
      const filePath = path.join(LOGS_FOLDER, file);
      const stat = fs.statSync(filePath);
      
      if (stat.isFile()) {
        const name = file;
        if (query === '' || name.toLowerCase().includes(query.toLowerCase())) {
          allLogs.push({ Name: name });
        }
      }
    }

    const totalCount = allLogs.length;

    const start = Math.min(offset, totalCount);
    const end = Math.min(offset + limit, totalCount);
    const pagedLogs = allLogs.slice(start, end);

    const totalPages = Math.ceil(totalCount / limit);

    Logger.info(`Local-api: Logs GET [${claims?.username}]`);
    
    res.status(200).json({
      Items: pagedLogs,
      TotalPages: totalPages,
      CurrentPage: page,
    });
    
  } catch (err) {
    Logger.error(`Local-api: Logs GET error [${claims?.username}]:`, err);
    res.status(404).json({ error: "not found" });
  }
}

export function getLog(req: Request, res: Response): void {
  const claims = req.user as Claims;
  const logFile = req.query.name as string;

  if (!logFile) {
    res.status(404).json({ error: "" });
    return;
  }
  
  const fullPath = path.join(LOGS_FOLDER, logFile);

  try {
    if (!fs.existsSync(fullPath)) {
      res.status(404).json({ error: "file not found" });
      return;
    }
    
    Logger.info(`Local-api: Log GET [${claims?.username}]`);
    
    res.download(fullPath, logFile, (err) => {
      if (err) {
        Logger.error(`Local-api: Error sending log file ${fullPath}:`, err);
        if (!res.headersSent) {
          res.status(404).json({ error: "file not found" });
        }
      }
    });
    
  } catch (err) {
    Logger.error(`Local-api: Log GET error [${claims?.username}]:`, err);
    res.status(404).json({ error: "file not found" });
  }
}