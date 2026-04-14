import type { Announce, Goods, Order, Schedule, User, Role, Permission, Log } from "@/types"

export interface GoodsResponse {
	Items: Goods[];
	TotalPages: number;
	CurrentPage: number;
}

export interface GoodsUpdateRequest {
	ID: number;
    Name: string;
    Instruction: string | undefined;
    Description: string | undefined;
    Prescription: boolean;
    IsInStock: boolean;
    Price: number;
}

export interface ScheduleResponse {
	Items: Schedule[];
	TotalPages: number;
	CurrentPage: number;
}

export interface OrdersResponse {
	Items: Order[];
	TotalPages: number;
	CurrentPage: number;
}

export interface AnnouncesResponse {
	Items: Announce[];
	TotalPages: number;
	CurrentPage: number;
}

export interface UsersResponse {
	Items: User[];
	TotalPages: number;
	CurrentPage: number;
}

export interface UserResponse {
	User: User;
	Roles: Role[] | undefined;
}

export interface RolesResponse {
	Items: Role[];
	TotalPages: number;
	CurrentPage: number;
}

export interface PermissionsResponse {
	Items: Permission[];
	TotalPages: number;
	CurrentPage: number;
}

export interface LogsResponse {
	Items: Log[];
	TotalPages: number;
	CurrentPage: number;
}

export interface KeyResponse {
  TaskID: string;
}

type TaskStatus = 'completed' | 'pending' | 'error' | 'not_found';


export interface TaskResponse {
  Status: TaskStatus;
  Value: string;
  Error?: string;
}