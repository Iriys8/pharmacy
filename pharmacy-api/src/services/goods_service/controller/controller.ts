import { DataSource } from 'typeorm';
import { Goods } from '@pharmacy/src/shared/models/database_models';
import { GoodsResponse, PromoItem, GoodsUpdateRequest } from '@pharmacy/src/shared/models/responses';
import { Claims } from "@pharmacy/src/shared/models/models"

export async function getGoods(
  dataSource: DataSource,
  query: string,
  pageStr: string,
  limitStr: string
): Promise<{ Items: GoodsResponse[]; TotalPages: number; CurrentPage: number }> {
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

  const goodsRepository = dataSource.getRepository(Goods);
  let goods: Goods[] = [];
  let totalCount = 0;

  if (query) {
    const goodsByTags = await goodsRepository
      .createQueryBuilder('goods')
      .leftJoinAndSelect('goods.producer', 'producer')
      .leftJoin('goods.goodsTags', 'goodsTags')
      .leftJoin('goodsTags.tag', 'tag')
      .where('tag.tagName LIKE :query', { query: `%${query}%` })
      .orderBy('goods.name', 'ASC')
      .getMany();

    const goodsByProducers = await goodsRepository
      .createQueryBuilder('goods')
      .leftJoinAndSelect('goods.producer', 'producer')
      .where('producer.producerName LIKE :query', { query: `%${query}%` })
      .orderBy('goods.name', 'ASC')
      .getMany();

    const goodsByName = await goodsRepository
      .createQueryBuilder('goods')
      .leftJoinAndSelect('goods.producer', 'producer')
      .where('goods.name LIKE :query', { query: `%${query}%` })
      .orderBy('goods.name', 'ASC')
      .getMany();

    const seenIDs = new Set<number>();
    const allGoods = [...goodsByTags, ...goodsByProducers, ...goodsByName];
    
    for (const item of allGoods) {
      if (!seenIDs.has(item.id)) {
        goods.push(item);
        seenIDs.add(item.id);
      }
    }

    totalCount = goods.length;
    goods = goods.slice(offset, offset + limit);
  } else {
    [goods, totalCount] = await goodsRepository
      .createQueryBuilder('goods')
      .leftJoinAndSelect('goods.producer', 'producer')
      .orderBy('goods.name', 'ASC')
      .skip(offset)
      .take(limit)
      .getManyAndCount();
  }

  const totalPages = Math.ceil(totalCount / limit);

  const response: GoodsResponse[] = goods.map(goodsItem => ({
    ID: goodsItem.id,
    Name: goodsItem.name,
    Image: goodsItem.image,
    IsInStock: goodsItem.isInStock,
    Description: goodsItem.description,
    Price: goodsItem.price,
    Producer: null,
    Tags: null,
    IsPrescriptionNeeded: null,
    Instruction: null
  }));

  return {
    Items: response,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getGoodsByID(
  dataSource: DataSource,
  id: number
): Promise<GoodsResponse> {
  const goodsRepository = dataSource.getRepository(Goods);
  
  const good = await goodsRepository.findOne({
    where: { id },
    relations: ['producer', 'goodsTags', 'goodsTags.tag'],
  });

  if (!good) {
    throw new Error('Good not found');
  }

  const tagNames = good.goodsTags?.map(goodsTag => goodsTag.tag?.tagName).filter(Boolean) || [];

  const response: GoodsResponse = {
    ID: good.id,
    Name: good.name,
    Image: good.image,
    Producer: good.producer?.producerName || '',
    IsInStock: good.isInStock,
    Tags: tagNames,
    Instruction: good.instruction,
    Description: good.description,
    IsPrescriptionNeeded: good.isPrescriptionNeeded,
    Price: good.price,
  };

  return response;
}

export async function getPromoItems(dataSource: DataSource): Promise<PromoItem[]> {
  const goodsRepository = dataSource.getRepository(Goods);
  
  const goods = await goodsRepository.find({
    where: { isInStock: true },
  });

  if (goods.length === 0) {
    return [];
  }

  const shuffled = [...goods];
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
  }
  
  const promoCount = Math.min(5, shuffled.length);
  const promoItems: PromoItem[] = shuffled.slice(0, promoCount).map(item => ({
    id: item.id,
    name: item.name,
    description: item.description,
    price: item.price,
    image: item.image,
  }));

  return promoItems;
}

export async function updateGoods(
  dataSource: DataSource,
  id: number,
  updateData: GoodsUpdateRequest,
  claims: Claims
): Promise<string> {
  const goodsRepository = dataSource.getRepository(Goods);
  
  const existingGood = await goodsRepository.findOne({
    where: { id },
  });

  if (!existingGood) {
    console.log(`Good error [${claims.username}] Good not found with id: ${id}`);
    throw new Error('Good not found');
  }

  existingGood.name = updateData.Name;
  existingGood.instruction = updateData.Instruction;
  existingGood.description = updateData.Description;
  existingGood.isPrescriptionNeeded = updateData.IsPrescriptionNeeded;
  existingGood.isInStock = updateData.IsInStock;
  existingGood.price = updateData.Price;

  await goodsRepository.save(existingGood);

  console.log(`Good PATH [${claims.username}] - Updated good ID: ${id}`);

  return "Good updated";
}