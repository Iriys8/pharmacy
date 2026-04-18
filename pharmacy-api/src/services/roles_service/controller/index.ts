import { DataSource, Like, In } from 'typeorm';
import { Role, Permission } from '@pharmacy/src/shared/models/database_models';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller';
import { PermissionRespons, RoleResponse } from '@pharmacy/src/shared/models/responses';

export async function getRoles(dataSource: DataSource, query: string, pageStr: string, limitStr: string, claims: Claims): Promise<{ Items: RoleResponse[]; TotalPages: number; CurrentPage: number }> {
  const roleRepository = dataSource.getRepository(Role);
  
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
  
  let roles: Role[] = [];
  let totalCount = 0;
  
  if (query) {
    roles = await roleRepository.find({
      where: {
        name: Like(`%${query}%`)
      },
      relations: ['permissions']
    });
    
    totalCount = roles.length;
    
    const start = Math.min(offset, roles.length);
    const end = Math.min(offset + limit, roles.length);
    roles = roles.slice(start, end);
  } else {
    [roles, totalCount] = await roleRepository.findAndCount({
      relations: ['permissions'],
      order: { id: 'DESC' },
      skip: offset,
      take: limit,
    });
  }

  const responseRoles: RoleResponse[] = roles.map(role => ({
    ID: role.id,
    Name: role.name,
    Permissions: role.permissions.map(permission => ({
      ID: permission.id,
      Action: permission.action,
    }))
  }))
  
  const totalPages = Math.ceil(totalCount / limit);
  
  Logger.info(`Roles service: Roles GET [${claims.username}]`);
  
  return {
    Items: responseRoles,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getRoleByID(dataSource: DataSource, id: number, claims: Claims): Promise<RoleResponse> {
  const roleRepository = dataSource.getRepository(Role);
  
  const role = await roleRepository.findOne({
    where: { id },
    relations: ['permissions']
  });
  
  if (!role) {
    throw new Error('Role not found');
  }
  
  const responseRoles: RoleResponse = {
    ID: role.id,
    Name: role.name,
    Permissions: role.permissions.map(permission => ({
      ID: permission.id,
      Action: permission.action,
    }))
  }

  Logger.info(`Roles service: Role GET [${claims.username}]`);
  
  return responseRoles;
}

export async function createRole(dataSource: DataSource,roleData: Partial<RoleResponse>,claims: Claims): Promise<string> {
  const roleRepository = dataSource.getRepository(Role);
  const permissionRepository = dataSource.getRepository(Permission);
  
  const newRole = new Role();
  newRole.name = roleData.Name || '';
  
  const savedRole = await roleRepository.save(newRole);
  
  if (roleData.Permissions && roleData.Permissions.length > 0) {
    const permissionIds = roleData.Permissions.map(p => p.ID);
    const permissions = await permissionRepository.findBy({
      id: In(permissionIds)
    });
    
    savedRole.permissions = permissions;
    await roleRepository.save(savedRole);
  }
  
  Logger.info(`Roles service: Role POST [${claims.username}]`);
  
  return 'role created';
}

export async function deleteRole(dataSource: DataSource, id: number, claims: Claims): Promise<string> {
  const roleRepository = dataSource.getRepository(Role);
  
  const role = await roleRepository.findOne({
    where: { id },
    relations: ['permissions']
  });
  
  if (!role) {
    throw new Error('Role not found');
  }
  
  role.permissions = [];
  await roleRepository.save(role);
  
  const result = await roleRepository.delete(id);
  
  if (result.affected === 0) {
    throw new Error('Role not found');
  }
  
  Logger.info(`Roles service: Role DELETE [${claims.username}]`);
  
  return 'role deleted';
}

export async function updateRole(dataSource: DataSource, id: number, roleData: Partial<RoleResponse>, claims: Claims): Promise<string> {
  const roleRepository = dataSource.getRepository(Role);
  const permissionRepository = dataSource.getRepository(Permission);
  
  const existingRole = await roleRepository.findOne({
    where: { id },
    relations: ['permissions']
  });
  
  if (!existingRole) {
    throw new Error('Role not found');
  }
  
  if (roleData.Name) {
    existingRole.name = roleData.Name;
  }
  
  if (roleData.Permissions) {
    const permissionIds = roleData.Permissions.map(p => p.ID);
    const permissions = await permissionRepository.findBy({
      id: In(permissionIds)
    });
    existingRole.permissions = permissions;
  }
  
  await roleRepository.save(existingRole);
  
  Logger.info(`Roles service: Role PATCH [${claims.username}]`);
  
  return 'role updated';
}

export async function getPermissions(dataSource: DataSource, query: string, pageStr: string, limitStr: string, claims: Claims): Promise<{ Items: PermissionRespons[]; TotalPages: number; CurrentPage: number }> {
  const permissionRepository = dataSource.getRepository(Permission);
  
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
  
  let permissions: Permission[] = [];
  let totalCount = 0;
  
  if (query) {
    permissions = await permissionRepository.find({
      where: {
        action: Like(`%${query}%`)
      }
    });
    
    totalCount = permissions.length;
    
    const start = Math.min(offset, permissions.length);
    const end = Math.min(offset + limit, permissions.length);
    permissions = permissions.slice(start, end);
  } else {
    [permissions, totalCount] = await permissionRepository.findAndCount({
      order: { id: 'DESC' },
      skip: offset,
      take: limit,
    });
  }
  
  const totalPages = Math.ceil(totalCount / limit);
  
  const permissionsResponse: PermissionRespons[] = permissions.map(permission => ({
    ID: permission.id,
    Action: permission.action,
  }))

  Logger.info(`Roles service: Permissions GET [${claims.username}]`);
  
  return {
    Items: permissionsResponse,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}