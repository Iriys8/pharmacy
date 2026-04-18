import { DataSource, Like } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { User, Role, Permission } from '@pharmacy/src/shared/models/database_models';
import { RoleResponse, UserResponse, UserUpdateRequest } from '@pharmacy/src/shared/models/responses';
import { Claims } from '@pharmacy/src/shared/models/models';
import { Logger } from '@pharmacy/src/shared/controllers/logs_controller'

export async function getUsers(dataSource: DataSource, query: string, pageStr: string, limitStr: string, claims: Claims): Promise<{ Items: UserResponse[]; TotalPages: number; CurrentPage: number }> {
  const userRepository = dataSource.getRepository(User);
  
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
  
  let users: User[] = [];
  let totalCount = 0;
  
  if (query) {
    users = await userRepository.find({
      where: {
        userName: Like(`%${query}%`)
      },
      relations: ['role']
    });
    
    totalCount = users.length;
    
    const start = Math.min(offset, users.length);
    const end = Math.min(offset + limit, users.length);
    users = users.slice(start, end);
  } else {
    [users, totalCount] = await userRepository.findAndCount({
      relations: ['role'],
      order: { id: 'DESC' },
      skip: offset,
      take: limit,
    });
  }
  
  const usersResponse: UserResponse[] = users.map(user => ({
    ID: user.id,
    Login: user.login,
    UserName: user.userName,
    RoleID: user.roleId,
    Role: {
      ID: user.role.id,
      Name: user.role.name,
      Permissions: []
    },
    Password: ''
  }));
  
  const totalPages = Math.ceil(totalCount / limit);
  
  Logger.info(`Users service: Users GET [${claims.username}]`);
  
  return {
    Items: usersResponse,
    TotalPages: totalPages,
    CurrentPage: page,
  };
}

export async function getUserByID(dataSource: DataSource, id: number, claims: Claims): Promise<{ User: UserResponse; Roles: RoleResponse[] }> {
  const userRepository = dataSource.getRepository(User);
  const roleRepository = dataSource.getRepository(Role);
  
  let requestedUser: User | null = null;
  let response: UserResponse | null = null;

  if (id !== 0) {
    requestedUser = await userRepository.findOne({
      where: { id },
      relations: ['role']
    });

    response = {
      ID: requestedUser.id,
      Login: requestedUser.login,
      UserName: requestedUser.userName,
      RoleID: requestedUser.roleId,
      Role: {
        ID: requestedUser.role.id,
        Name: requestedUser.role.name,
        Permissions: [],
      },
      Password: ""
    }
  }

  let roles: Role[] = [];
  
  const requesterUser = await userRepository.findOne({
    where: { id: claims.userId },
    relations: ['role', 'role.permissions']
  });
  
  if (!requesterUser) {
    throw new Error('User not found');
  }
  
  const hasChangeRolesPermission = requesterUser.role?.permissions?.some(
    (permission: Permission) => permission.action === 'Change_Roles'
  );
  
  if (hasChangeRolesPermission) {
    roles = await roleRepository.find();
  }

  const responseRoles: RoleResponse[] = roles.map(role => ({
    ID: role.id,
    Name: role.name,
    Permissions: [],
  }))
  
  Logger.info(`Users service: User GET [${claims.username}]`);
  
  return {
    User: response,
    Roles: responseRoles,
  };
}

export async function createUser(dataSource: DataSource, newUserRequest: UserUpdateRequest, claims: Claims): Promise<string> {
  const userRepository = dataSource.getRepository(User);
  
  const saltRounds = 10;
  const passwordHash = await bcrypt.hash(newUserRequest.Password, saltRounds);
  
  const newUser = new User();
  newUser.login = newUserRequest.Login;
  newUser.userName = newUserRequest.UserName;
  newUser.roleId = newUserRequest.RoleID;
  newUser.passwordHash = passwordHash;
  
  await userRepository.save(newUser);
  
  Logger.info(`Users service: User POST [${claims.username}]`);
  
  return 'user created';
}

export async function deleteUser(dataSource: DataSource, id: number, claims: Claims): Promise<string> {
  const userRepository = dataSource.getRepository(User);
  
  const result = await userRepository.delete(id);
  
  if (result.affected === 0) {
    throw new Error('User not found');
  }
  
  Logger.info(`Users service: User DELETE [${claims.username}]`);
  
  return 'user deleted';
}

export async function updateUser(dataSource: DataSource, id: number, userUpdateRequest: UserUpdateRequest, claims: Claims): Promise<string> {
  const userRepository = dataSource.getRepository(User);
  
  const existingUser = await userRepository.findOne({
    where: { id }
  });
  
  if (!existingUser) {
    throw new Error('User not found');
  }

  existingUser.login = userUpdateRequest.Login;
  existingUser.userName = userUpdateRequest.UserName;
  existingUser.roleId = userUpdateRequest.RoleID;
  
  if (userUpdateRequest.Password && userUpdateRequest.Password !== '') {
    const saltRounds = 10;
    existingUser.passwordHash = await bcrypt.hash(userUpdateRequest.Password, saltRounds);
  }
  
  await userRepository.save(existingUser);
  
  Logger.info(`Users service: User PATCH [${claims.username}]`);
  
  return 'user updated';
}