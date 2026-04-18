import winston from 'winston';
import path from 'path';

const LOGS_FOLDER = path.resolve(process.cwd(), './pharmacy-content/logs');

export const Logger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json(),
    winston.format.printf(({ timestamp, level, message }) => {
      return `${timestamp} [${level.toUpperCase()}]: ${message}`;
    })
  ),
  transports: [
    new winston.transports.File({ 
      filename: path.join(LOGS_FOLDER, 'error.log'), 
      level: 'error' 
    }),
    new winston.transports.File({ 
      filename: path.join(LOGS_FOLDER, 'combined.log') 
    }),
    new winston.transports.Console({
      format: winston.format.simple()
    })
  ]
});