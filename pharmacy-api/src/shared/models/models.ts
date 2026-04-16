import * as amqp from 'amqplib';
import * as jwt from 'jsonwebtoken';

export interface Broker {
  connection: amqp.ChannelModel | null;
  channel: amqp.Channel | null;
}

export interface RequestContext {
  query: any;
  context: any;
  claims: any;
}

export interface Claims extends jwt.JwtPayload {
  userId: number;
  username: string;
  role: string;
}

export interface RefreshClaims extends jwt.JwtPayload {
  userId: number;
}

export interface Message {
  exchange: string;
  routingKey: string;
  mandatory: boolean;
  immediate: boolean;
  publishing: amqp.MessageProperties & {
    content?: Buffer;
  };
}

export class RawMessage {
  private data: Buffer;

  constructor(data: Buffer | string | object) {
    if (Buffer.isBuffer(data)) {
      this.data = data;
    } else if (typeof data === 'string') {
      this.data = Buffer.from(data);
    } else {
      this.data = Buffer.from(JSON.stringify(data));
    }
  }

  toString(): string {
    return this.data.toString();
  }

  toBuffer(): Buffer {
    return this.data;
  }

  toJSON<T = any>(): T {
    return JSON.parse(this.data.toString());
  }
}

export interface RequestContextWithRaw {
  query: RawMessage;
  context: RawMessage;
  claims: RawMessage;
}