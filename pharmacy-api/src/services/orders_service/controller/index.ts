import { DataSource, Like } from 'typeorm';
import { Order, GoodsOrders } from '@pharmacy/src/shared/models/database_models';
import { OrderResponse, OrderedItem } from '@pharmacy/src/shared/models/responses';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'

export async function getOrders(dataSource: DataSource, query: string, pageStr: string, limitStr: string, claims: Claims): Promise<{ Items: OrderResponse[]; TotalPages: number; CurrentPage: number }> {
  const orderRepository = dataSource.getRepository(Order);
  
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
  
  let orders: Order[] = [];
  let totalCount = 0;
  
  if (query) {
    orders = await orderRepository.find({
      where: {
        clientFIO: Like(`%${query}%`)
      }
    });
    
    totalCount = orders.length;
    
    const start = Math.min(offset, orders.length);
    const end = Math.min(offset + limit, orders.length);
    orders = orders.slice(start, end);
  } else {
    [orders, totalCount] = await orderRepository.findAndCount({
      order: { id: 'DESC' },
      skip: offset,
      take: limit,
    });
  }
  
  const response: OrderResponse[] = orders.map(order => ({
    ID: order.id,
    Name: order.clientFIO,
    Phone: order.clientPhone,
    Email: order.clientEmail,
    Items: [],
  }));
  
  const totalPages = Math.ceil(totalCount / limit);
  
  Logger.info(`Orders service: Orders GET [${claims.username}]`);
  
  return {
    Items: response,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getOrderByID(dataSource: DataSource, id: number, claims: Claims): Promise<OrderResponse> {
  const order = await dataSource
    .createQueryBuilder()
    .select('order')
    .from(Order, 'order')
    .leftJoinAndSelect('order.goodsOrders', 'goodsOrders')
    .leftJoinAndSelect('goodsOrders.goods', 'goods')
    .where('order.id = :id', { id })
    .getOne();
  
  if (!order) {
    throw new Error('Order not found');
  }
  
  const response: OrderResponse = {
    ID: order.id,
    Name: order.clientFIO,
    Email: order.clientEmail,
    Phone: order.clientPhone,
    Items: [],
  };
  
  if (order.goodsOrders) {
    for (const goodsOrder of order.goodsOrders) {
      if (goodsOrder.goods) {
        const item: OrderedItem = {
          ID: goodsOrder.goods.id,
          Name: goodsOrder.goods.name,
          Image: goodsOrder.goods.image,
          Description: goodsOrder.goods.description,
          Price: goodsOrder.goods.price,
          Quantity: goodsOrder.quantity,
        };
        response.Items.push(item);
      }
    }
  }
  
  Logger.info(`Orders service: Order GET [${claims.username}]`);
  
  return response;
}

export async function createOrder(dataSource: DataSource, requestData: OrderResponse): Promise<string> {
  const orderRepository = dataSource.getRepository(Order);
  const goodsOrdersRepository = dataSource.getRepository(GoodsOrders);
  
  const order = new Order();
  order.clientFIO = requestData.Name;
  order.clientEmail = requestData.Email || '';
  order.clientPhone = requestData.Phone;
  
  const savedOrder = await orderRepository.save(order);
  
  for (const item of requestData.Items) {
    const goodsOrder = new GoodsOrders();
    goodsOrder.orderId = savedOrder.id;
    goodsOrder.goodsId = item.ID;
    goodsOrder.quantity = item.Quantity;
    await goodsOrdersRepository.save(goodsOrder);
  }
  
  return 'Order created successfully';
}

export async function updateOrder(dataSource: DataSource, id: number, input: OrderResponse, claims: Claims): Promise<string> {
  const orderRepository = dataSource.getRepository(Order);
  const goodsOrdersRepository = dataSource.getRepository(GoodsOrders);
  
  const order = await orderRepository.findOne({
    where: { id }
  });
  
  if (!order) {
    throw new Error('Order not found');
  }
  
  order.clientFIO = input.Name;
  order.clientEmail = input.Email || '';
  order.clientPhone = input.Phone;
  
  await orderRepository.save(order);
  
  await goodsOrdersRepository.delete({ orderId: id });
  
  for (const item of input.Items) {
    const goodsOrder = new GoodsOrders();
    goodsOrder.orderId = input.ID;
    goodsOrder.goodsId = item.ID;
    goodsOrder.quantity = item.Quantity;
    await goodsOrdersRepository.save(goodsOrder);
  }
  
  Logger.info(`Orders service: Order PATCH [${claims.username}]`);
  
  return 'order updated';
}

export async function deleteOrder(dataSource: DataSource, id: number, claims: Claims): Promise<string> {
  const orderRepository = dataSource.getRepository(Order);
  
  const result = await orderRepository.delete(id);
  
  if (result.affected === 0) {
    throw new Error('Order not found');
  }
  
  Logger.info(`Orders service: Order DELETE [${claims.username}]`);
  
  return 'order deleted';
}