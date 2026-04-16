import { Channel, Options } from 'amqplib';
import { Broker } from '@pharmacy/src/shared/models/models';

export interface BrokerSetup {
  channel: Channel;
}

export async function setupBroker(broker: Broker): Promise<void> {
  try {
    const exchangeOptions: Options.AssertExchange = {
      durable: true,
      autoDelete: false,
      arguments: {}
    };
    
    await broker.channel.assertExchange("public_exchange", "direct", exchangeOptions);

    const queues = [
      { name: "public_goods_queue", routingKey: "goods" },
      { name: "public_schedule_queue", routingKey: "schedule" },
      { name: "public_orders_queue", routingKey: "orders" },
      { name: "public_announces_queue", routingKey: "announces" },
    ];

    for (const q of queues) {
      const queueOptions: Options.AssertQueue = {
        durable: true,
        autoDelete: false,
        exclusive: false,
      };
      
      const queue = await broker.channel.assertQueue(q.name, queueOptions);

      await broker.channel.bindQueue(queue.queue, "public_exchange", q.routingKey);
    }

    console.log("RabbitMQ setup completed successfully");
    console.log(`Exchange "public_exchange" declared with 4 queues bound`);
  } catch (err) {
    console.error("Failed to setup RabbitMQ:", err);
    throw err;
  }
}

export async function teardownBroker(broker: { channel: Channel }): Promise<void> {
  try {
    const queues = [
      "public_goods_queue",
      "public_schedule_queue",
      "public_orders_queue",
      "public_announces_queue",
    ];

    for (const queue of queues) {
      await broker.channel.deleteQueue(queue);
    }

    await broker.channel.deleteExchange("public_exchange");
    console.log("RabbitMQ teardown completed successfully");
  } catch (err) {
    console.error("Failed to teardown RabbitMQ:", err);
    throw err;
  }
}