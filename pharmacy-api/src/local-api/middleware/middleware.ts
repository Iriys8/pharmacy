import { Request, Response, NextFunction } from 'express';
import { DataSource } from 'typeorm';
import * as jwt from 'jsonwebtoken';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Role, Permission } from '@pharmacy/src/shared/models/database_models';

const jwtSecret = process.env.ACCESSTOKEN_SECRET;

declare global {
  namespace Express {
    interface Request {
      user?: Claims;
    }
  }
}

export function authMiddleware(dataSource: DataSource, requiredPermission: string = '') {
  return async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    const authHeader = req.headers['authorization'];
    
    if (!authHeader) {
      res.status(401).json({ error: "Authorization header required" });
      return;
    }

    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') {
      res.status(401).json({ error: "Invalid authorization format" });
      return;
    }

    const tokenString = parts[1];

    try {
      const decoded = jwt.verify(tokenString, jwtSecret) as Claims;
      const claims = decoded;

      if (!requiredPermission) {
        req.user = claims;
        next();
        return;
      }

      const roleRepository = dataSource.getRepository(Role);
      const role = await roleRepository.findOne({
        where: { name: claims.role },
        relations: ['permissions'],
      });

      if (!role) {
        console.log(`Middleware error [${claims.username}] Role not found: ${claims.role}`);
        res.status(401).json({ error: "Role not found" });
        return;
      }

      const hasPermission = role.permissions?.some(
        (permission: Permission) => permission.action === requiredPermission
      );

      if (!hasPermission) {
        console.log(`Middleware error [${claims.username}] Permission denied: ${requiredPermission}`);
        res.status(403).json({ error: "Permission denied" });
        return;
      }

      req.user = claims;
      next();
      
    } catch (err) {
      if (err instanceof jwt.TokenExpiredError) {
        res.status(401).json({ error: "token_expired" });
      } else if (err instanceof jwt.JsonWebTokenError) {
        res.status(401).json({ error: "Invalid token" });
      } else {
        console.error(`Middleware error:`, err);
        res.status(500).json({ error: "Internal server error" });
      }
      return;
    }
  };
}