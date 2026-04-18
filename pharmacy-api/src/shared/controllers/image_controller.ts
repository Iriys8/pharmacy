import { Request, Response } from 'express';
import path from 'path';
import fs from 'fs';

const IMAGE_FOLDER = path.resolve(process.cwd(), './pharmacy-content/images');

export async function getImage(req: Request, res: Response): Promise<void> {
  const imageName = req.query.name as string;
  
  if (!imageName) {
    res.status(404).json({ error: "" });
    return;
  }
  
  const filePath = path.join(IMAGE_FOLDER, imageName);
  
  if (!fs.existsSync(filePath)) {
    res.status(404).json({ error: "" });
    return;
  }
  
  res.sendFile(filePath, (err) => {
    if (err) {
      console.error(`Error sending file ${filePath}:`, err);
      if (!res.headersSent) {
        res.status(404).json({ error: "" });
      }
    }
  });
}