import { Channel } from 'amqplib';
import { Broker } from '@pharmacy/src/shared/models/models';

export async function setupBroker(broker: Broker): Promise<void> {
  try {
    await broker.channel.assertExchange("local_exchange", "direct", {
      durable: true,
      autoDelete: false,
    });

    const queues = [
      { name: "local_goods_queue", routingKey: "goods" },
      { name: "local_schedule_queue", routingKey: "schedule" },
      { name: "local_orders_queue", routingKey: "orders" },
      { name: "local_announces_queue", routingKey: "announces" },
      { name: "local_users_queue", routingKey: "users" },
      { name: "local_roles_queue", routingKey: "roles" },
    ];

    for (const q of queues) {
      const queue = await broker.channel.assertQueue(q.name, {
        durable: true,
        autoDelete: false,
      });

      await broker.channel.bindQueue(queue.queue, "local_exchange", q.routingKey);
    }

    console.log("RabbitMQ setup completed successfully");
    console.log(`Exchange "local_exchange" declared with ${queues.length} queues bound`);
  } catch (err) {
    console.error("Failed to setup RabbitMQ:", err);
    throw err;
  }
}