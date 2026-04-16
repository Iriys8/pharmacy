import { DataSource } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { Goods, User, Role, Permission } from '@pharmacy/src/shared/models/database_models';
import { testData } from '@pharmacy/src/shared/setup/db_data';

export async function setupDB(dataSource: DataSource): Promise<void> {
  const queryRunner = dataSource.createQueryRunner();
  try {
    await queryRunner.connect();
    
    let goodsCount: Goods[] = null
    try {
      goodsCount = await dataSource.getRepository(Goods).createQueryBuilder('goods').limit(1).getMany();
      if (goodsCount.length === 0) {
        await testData(dataSource);
        console.log("TEST DATA USED!");
      }
    }
    catch {
      await dataSource.runMigrations();
      await testData(dataSource);
      console.log("TEST DATA USED!");
    }

    const permissionsCount = await dataSource.getRepository(Permission).count();
    
    if (permissionsCount === 0) {
      const permissions: Partial<Permission>[] = [
        { action: "Update_Goods" },
        { action: "Read_Orders" },
        { action: "Update_Orders" },
        { action: "Create_Orders" },
        { action: "Delete_Orders" },
        { action: "Update_Schedule" },
        { action: "Create_Schedule" },
        { action: "Delete_Schedule" },
        { action: "Update_Announces" },
        { action: "Create_Announces" },
        { action: "Delete_Announces" },
        { action: "Change_Users" },
        { action: "Change_Roles" },
        { action: "Download_Logs" },
      ];
      await dataSource.getRepository(Permission).save(permissions);
      console.log("Permissions initialized");
    }

    const usersCount = await dataSource.getRepository(User).count();

    if (usersCount === 0 && process.env.CONTROL_PANEL_ADMIN_PASSWORD && process.env.CONTROL_PANEL_ADMIN_LOGIN) {
      await setupAdmin(dataSource);
    }

  } catch (error) {
    console.error("Failed to setup database:", error);
    throw error;
  } finally {
    await queryRunner.release();
  }
}

export async function setupAdmin(dataSource: DataSource): Promise<void> {
  const saltRounds = 10;

  const passwordHash = await bcrypt.hash(process.env.CONTROL_PANEL_ADMIN_PASSWORD, saltRounds);
  
  const permissions = await dataSource.getRepository(Permission).find();
  
  const adminRole = new Role();
  adminRole.name = "Admin";
  adminRole.permissions = permissions;
  
  const savedRole = await dataSource.getRepository(Role).save(adminRole);
  
  const user = new User();
  user.login = process.env.CONTROL_PANEL_ADMIN_LOGIN;
  user.userName = "Admin";
  user.passwordHash = passwordHash;
  user.roleId = savedRole.id;
  user.role = savedRole;
  
  await dataSource.getRepository(User).save(user);
  console.log("Admin user created successfully");
}