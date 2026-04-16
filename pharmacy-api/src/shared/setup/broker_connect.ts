import * as amqp from 'amqplib';
import { Broker } from '@pharmacy/src/shared/models/models';

export async function connectBroker(): Promise<Broker> {
  const connString = `amqp://${process.env.RABBITMQ_USER}:${process.env.RABBITMQ_PASSWORD}@${process.env.RABBITMQ_HOST}:5672/`;
  
  try {
    const connection = await amqp.connect(connString);
    console.log('Connected to RabbitMQ successfully');
    
    const channel = await connection.createChannel();
    console.log('RabbitMQ channel created successfully');
    
    return {
      connection: connection,
      channel: channel,
    };
  } catch (error) {
    console.error('Failed to connect to broker:', error);
    throw error;
  }
}

export async function closeBroker(broker: Broker): Promise<void> {
  try {
    if (broker.channel) {
      await broker.channel.close();
      console.log('RabbitMQ channel closed');
    }
    if (broker.connection) {
      await broker.connection.close();
      console.log('RabbitMQ connection closed');
    }
  } catch (error) {
    console.error('Error while closing broker connection:', error);
    throw error;
  }
}