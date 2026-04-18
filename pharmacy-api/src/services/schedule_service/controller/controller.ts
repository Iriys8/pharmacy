import { DataSource, Between, Like, Raw } from 'typeorm';
import { Schedule } from '@pharmacy/src/shared/models/database_models';
import { ScheduleResponse } from '@pharmacy/src/shared/models/responses';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'

function formatTime(time: Date | string): string {
  if (time instanceof Date) {
    const hours = String(time.getHours).padStart(2, '0');
    const minutes = String(time.getMinutes).padStart(2, '0');
    return `${hours}:${minutes}`;
  }
  if (typeof time === 'string') {
    return time.substring(0, 5);
  }
  return '09:00';
}

function parseTimeString(timeStr: string): string {
  if (/^\d{2}:\d{2}$/.test(timeStr)) {
    return `${timeStr}:00`;
  }
  return '09:00:00';
}

export async function getScheduleDated(
  dataSource: DataSource,
  startDate: string,
  endDate: string
): Promise<ScheduleResponse[]> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
  const startParsed = new Date(startDate);
  const endParsed = new Date(endDate);
  
  if (isNaN(startParsed.getTime()) || isNaN(endParsed.getTime())) {
    throw new Error('Invalid date format');
  }
  
  const maxEndDate = new Date(startParsed);
  maxEndDate.setDate(maxEndDate.getDate() + 31);
  
  if (endParsed < startParsed || endParsed > maxEndDate) {
    throw new Error('Invalid date range');
  }
  
  const schedule = await scheduleRepository.find({
    where: {
      date: Between(new Date(startDate), new Date(endDate))
    }
  });
  
  const response: ScheduleResponse[] = [];
  
  const currentDate = new Date(startParsed);
  while (currentDate <= endParsed) {
    const dateStr = currentDate.toISOString().split('T')[0];
    const found = schedule.find(item => {
      if (!item.date) return false;
      const itemDate = new Date(item.date);
      return itemDate.toDateString() === currentDate.toDateString();
    });
    
    if (found) {
      response.push({
        ID: found.id,
        Date: dateStr,
        TimeStart: found.timeStart ? formatTime(found.timeStart) : '09:00',
        TimeEnd: found.timeEnd ? formatTime(found.timeEnd) : '18:00',
        IsOpened: found.isOpened,
      });
    } else {
      response.push({
        ID: 0,
        Date: dateStr,
        TimeStart: '08:00',
        TimeEnd: '22:00',
        IsOpened: true,
      });
    }
    
    currentDate.setDate(currentDate.getDate() + 1);
  }
  
  return response;
}

export async function getSchedule(
  dataSource: DataSource,
  query: string,
  pageStr: string,
  limitStr: string,
  claims: Claims
): Promise<{ Items: ScheduleResponse[]; TotalPages: number; CurrentPage: number }> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
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
  
  let schedule: Schedule[] = [];
  let totalCount = 0;

    if (query) {
    schedule = await scheduleRepository.find({
      where: {
        date: Raw((alias) => `DATE(${alias}) LIKE :query`, { query: `%${query}%` })
      }
    });
    
    totalCount = schedule.length;
    
    const start = Math.min(offset, schedule.length);
    const end = Math.min(offset + limit, schedule.length);
    schedule = schedule.slice(start, end);
  } else {
    [schedule, totalCount] = await scheduleRepository.findAndCount({
      order: { date: 'ASC' },
      skip: offset,
      take: limit,
    });
  }
  
  const response: ScheduleResponse[] = schedule.map(item => ({
    ID: item.id,
    Date: item.date ? `${item.date.toString()}` : '',
    TimeStart: item.timeStart ? formatTime(item.timeStart) : '09:00',
    TimeEnd: item.timeEnd ? formatTime(item.timeEnd) : '18:00',
    IsOpened: item.isOpened,
  }));
  
  const totalPages = Math.ceil(totalCount / limit);
  
  Logger.info(`Schedule service: Schedule GET [${claims.username}]`);
  
  return {
    Items: response,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getScheduleByID(
  dataSource: DataSource,
  id: number,
  claims: Claims
): Promise<ScheduleResponse> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
  const schedule = await scheduleRepository.findOne({
    where: { id }
  });
  
  if (!schedule) {
    throw new Error('Schedule not found');
  }
  
  const response: ScheduleResponse = {
    ID: schedule.id,
    Date: schedule.date ? `${schedule.date.toString()}` : '',
    TimeStart: schedule.timeStart ? formatTime(schedule.timeStart) : '09:00',
    TimeEnd: schedule.timeEnd ? formatTime(schedule.timeEnd) : '18:00',
    IsOpened: schedule.isOpened,
  };
  
  Logger.info(`Schedule service: Schedule GET [${claims.username}]`);
  
  return response;
}

export async function createSchedule(
  dataSource: DataSource,
  scheduleData: ScheduleResponse,
  claims: Claims
): Promise<string> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
  const schedule = new Schedule();
  schedule.date = new Date(scheduleData.Date);
  schedule.timeStart = parseTimeString(scheduleData.TimeStart);
  schedule.timeEnd = parseTimeString(scheduleData.TimeEnd);
  schedule.isOpened = scheduleData.IsOpened;
  
  await scheduleRepository.save(schedule);
  
  Logger.info(`Schedule service: Schedule POST [${claims.username}]`);
  
  return 'schedule created';
}

export async function updateSchedule(
  dataSource: DataSource,
  id: number,
  scheduleData: ScheduleResponse,
  claims: Claims
): Promise<string> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
  const schedule = await scheduleRepository.findOne({
    where: { id }
  });
  
  if (!schedule) {
    throw new Error('Schedule not found');
  }
  
  schedule.date = new Date(scheduleData.Date);
  schedule.timeStart = parseTimeString(scheduleData.TimeStart);
  schedule.timeEnd = parseTimeString(scheduleData.TimeEnd);
  schedule.isOpened = scheduleData.IsOpened;
  
  await scheduleRepository.save(schedule);
  
  Logger.info(`Schedule service: Schedule PATCH [${claims.username}]`);
  
  return 'schedule updated';
}

export async function deleteSchedule(
  dataSource: DataSource,
  id: number,
  claims: Claims
): Promise<string> {
  const scheduleRepository = dataSource.getRepository(Schedule);
  
  const result = await scheduleRepository.delete(id);
  
  if (result.affected === 0) {
    throw new Error('Schedule not found');
  }
  
  Logger.info(`Schedule service: Schedule DELETE [${claims.username}]`);
  
  return 'schedule deleted';
}