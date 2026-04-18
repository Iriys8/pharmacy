import { DataSource, Like, Raw } from 'typeorm';
import { Announcement } from '@pharmacy/src/shared/models/database_models';
import { AnnouncementResponse } from '@pharmacy/src/shared/models/responses';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'
 
export async function getAnnounces(dataSource: DataSource, query: string, pageStr: string, limitStr: string): Promise<{ Items: AnnouncementResponse[]; TotalPages: number; CurrentPage: number }> {
  const announceRepository = dataSource.getRepository(Announcement);
  
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
  
  let announces: Announcement[] = [];
  let totalCount = 0;
  
  if (query) {
    announces = await announceRepository.find({
      where: {
        dateTime: Raw((alias) => `DATE(${alias}) LIKE :query`, { query: `%${query}%` })
      }
    });
    
    totalCount = announces.length;
    
    const start = Math.min(offset, announces.length);
    const end = Math.min(offset + limit, announces.length);
    announces = announces.slice(start, end);
  } else {
    [announces, totalCount] = await announceRepository.findAndCount({
      order: { dateTime: 'ASC' },
      skip: offset,
      take: limit,
    });
  }
  
  const response: AnnouncementResponse[] = announces.map(announce => ({
    ID: announce.id,
    DateTime: `${announce.dateTime.toDateString()}`,
    From: announce.from,
    Announce: announce.announce,
  }));
  
  const totalPages = Math.ceil(totalCount / limit);
  
  return {
    Items: response,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getAnnounceByID(dataSource: DataSource,id: number,claims: Claims): Promise<AnnouncementResponse> {
  const announceRepository = dataSource.getRepository(Announcement);
  
  const announce = await announceRepository.findOne({
    where: { id }
  });
  
  if (!announce) {
    throw new Error('Announcement not found');
  }
  
  const response: AnnouncementResponse = {
    ID: announce.id,
    DateTime: `${announce.dateTime.toDateString()}`,
    From: announce.from,
    Announce: announce.announce,
  };
  
  Logger.info(`Announces service: Announce GET [${claims.username}]`);
  
  return response;
}

export async function createAnnounce(dataSource: DataSource, announcementData: AnnouncementResponse, claims: Claims): Promise<string> {
  const announceRepository = dataSource.getRepository(Announcement);
  
  const now = new Date();
  const dateTime = new Date(
    now.getFullYear(),
    now.getMonth(),
    now.getDate(),
    now.getHours(),
    now.getMinutes(),
    0,
    0
  );
  
  const announce = new Announcement();
  announce.dateTime = dateTime;
  announce.from = claims.username;
  announce.announce = announcementData.Announce;
  
  await announceRepository.save(announce);
  
  Logger.info(`Announces service: Announce POST [${claims.username}]`);
  
  return 'announce created';
}

export async function updateAnnounce(dataSource: DataSource, id: number, announcementData: AnnouncementResponse, claims: Claims): Promise<string> {
  const announceRepository = dataSource.getRepository(Announcement);
  
  const announce = await announceRepository.findOne({
    where: { id }
  });
  
  if (!announce) {
    throw new Error('Announcement not found');
  }
  
  announce.announce = announcementData.Announce;
  
  await announceRepository.save(announce);
  
  Logger.info(`Announces service: Announce PATCH [${claims.username}]`);
  
  return 'announce updated';
}

export async function deleteAnnounce(dataSource: DataSource, id: number, claims: Claims): Promise<string> {
  const announceRepository = dataSource.getRepository(Announcement);
  
  const result = await announceRepository.delete(id);
  
  if (result.affected === 0) {
    throw new Error('Announcement not found');
  }
  
  Logger.info(`Announces service: Announce DELETE [${claims.username}]`);
  
  return 'announce deleted';
}