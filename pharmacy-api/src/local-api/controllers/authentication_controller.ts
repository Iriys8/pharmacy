import { Request, Response } from 'express';
import { DataSource } from 'typeorm';
import * as bcrypt from 'bcrypt';
import * as jwt from 'jsonwebtoken';
import { User } from '@pharmacy/src/shared/models/database_models';
import { Claims, RefreshClaims } from '@pharmacy/src/shared/models/models';

const jwtSecret = process.env.ACCESSTOKEN_SECRET || 'default_access_secret';
const refreshSecret = process.env.REFRESHTOKEN_SECRET || 'default_refresh_secret';

interface LoginData {
  login: string;
  password: string;
}

interface TokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: {
    username: string;
    permissions: string[];
  };
}

function generateAccessToken(user: User): string {
  const expirationTime = Math.floor(Date.now() / 1000) + 15 * 60; // 15 минут

  const claims: Claims = {
    userId: user.id,
    username: user.login,
    role: user.role?.name || '',
    exp: expirationTime,
    iat: Math.floor(Date.now() / 1000),
  };

  return jwt.sign(claims, jwtSecret);
}

function generateRefreshToken(user: User): string {
  const expirationTime = Math.floor(Date.now() / 1000) + 8 * 60 * 60; // 8 часов

  const claims: RefreshClaims = {
    userId: user.id,
    exp: expirationTime,
    iat: Math.floor(Date.now() / 1000),
  };

  return jwt.sign(claims, refreshSecret);
}

export function login(dataSource: DataSource) {
  return async (req: Request, res: Response): Promise<void> => {
    const loginData: LoginData = req.body;

    if (!loginData.login || !loginData.password) {
      res.status(400).json({ error: "Invalid input" });
      return;
    }

    try {
      const userRepository = dataSource.getRepository(User);
      const user = await userRepository.findOne({
        where: { login: loginData.login },
        relations: ['role', 'role.permissions'],
      });

      if (!user) {
        res.status(401).json({ error: "Invalid credentials" });
        return;
      }

      const userPermissions = user.role?.permissions?.map(p => p.action) || [];

      const passwordHash = typeof user.passwordHash === 'string' 
        ? user.passwordHash 
        : Buffer.from(user.passwordHash).toString();
      
      const isPasswordValid = await bcrypt.compare(loginData.password, passwordHash);
      
      if (!isPasswordValid) {
        res.status(401).json({ error: "Invalid credentials" });
        return;
      }

      const accessToken = generateAccessToken(user);
      const refreshToken = generateRefreshToken(user);

      res.cookie('pharmacy_refresh_token', refreshToken, {
        httpOnly: true,
        secure: true,
        sameSite: 'strict',
        maxAge: 8 * 60 * 60 * 1000,
        path: '/',
        domain: 'localhost',
      });

      console.log(`Login [${user.userName}]`);

      res.status(200).json({
        access_token: accessToken,
        token_type: "Bearer",
        expires_in: 15 * 60,
        user: {
          username: user.userName,
          permissions: userPermissions,
        },
      });
    } catch (err) {
      console.error("Login error:", err);
      res.status(500).json({ error: "Internal server error" });
    }
  };
}

export function logout() {
  return async (req: Request, res: Response): Promise<void> => {
    res.cookie('pharmacy_refresh_token', '', {
      httpOnly: true,
      secure: true,
      sameSite: 'strict',
      maxAge: -1,
      path: '/',
      domain: 'localhost',
    });
    
    res.status(200).json({ message: "Logged out successfully" });
  };
}

export function refreshToken(dataSource: DataSource) {
  return async (req: Request, res: Response): Promise<void> => {
    console.log(req.cookies)
    const refreshToken = req.cookies?.pharmacy_refresh_token;

    if (!refreshToken) {
      res.status(401).json({ error: "Refresh token required" });
      return;
    }

    try {
      const decoded = jwt.verify(refreshToken, refreshSecret) as RefreshClaims;

      if (!decoded.userId) {
        res.status(401).json({ error: "Invalid refresh token" });
        return;
      }

      const userRepository = dataSource.getRepository(User);
      const user = await userRepository.findOne({
        where: { id: decoded.userId },
        relations: ['role', 'role.permissions'],
      });

      if (!user) {
        res.status(401).json({ error: "User not found" });
        return;
      }

      const userPermissions = user.role?.permissions?.map(p => p.action) || [];

      const newAccessToken = generateAccessToken(user);
      const newRefreshToken = generateRefreshToken(user);

      res.cookie('pharmacy_refresh_token', newRefreshToken, {
        httpOnly: true,
        secure: true,
        sameSite: 'strict',
        maxAge: 8 * 60 * 60 * 1000,
        path: '/',
        domain: 'localhost',
      });

      res.status(200).json({
        access_token: newAccessToken,
        token_type: "Bearer",
        expires_in: 15 * 60,
        user: {
          username: user.userName,
          permissions: userPermissions,
        },
      });
    } catch (err) {
      if (err instanceof jwt.TokenExpiredError) {
        res.status(401).json({ error: "Refresh token expired" });
      } else if (err instanceof jwt.JsonWebTokenError) {
        res.status(401).json({ error: "Invalid refresh token" });
      } else {
        console.error("Refresh error:", err);
        res.status(500).json({ error: "Could not generate token" });
      }
    }
  };
}