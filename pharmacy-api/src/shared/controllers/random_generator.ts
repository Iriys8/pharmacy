import { randomBytes } from 'crypto';

export function randomGen(): string {
  const bytes = randomBytes(8);
  return bytes.toString('hex');
}