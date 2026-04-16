import { DataSource } from 'typeorm';
import { AppDataSource } from '@pharmacy/src/shared/models/database_models';
import { setupDB } from '@pharmacy/src/shared/setup/db_setup';

export async function connectDB(): Promise<DataSource> {
  try {
    if (!AppDataSource.isInitialized) {
      await AppDataSource.initialize();
      console.log('Database connected successfully');
    }
    await setupDB(AppDataSource)
    return AppDataSource;
  } catch (error) {
    console.error('Failed to connect to database:', error);
    throw error;
  }
}