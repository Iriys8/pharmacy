import { Role } from "@pharmacy/src/shared/models/database_models";

export interface GoodsResponse {
  ID: number;
  Name: string;
  Image: string;
  Producer: string;
  IsInStock: boolean;
  Tags: string[];
  Instruction: string;
  Description: string;
  IsPrescriptionNeeded: boolean;
  Price: number;
}

export interface ScheduleResponse {
  ID: number;
  Date: string;
  TimeStart: string;
  TimeEnd: string;
  IsOpened: boolean;
}

export interface AnnouncementResponse {
  ID: number;
  DateTime: string;
  From: string;
  Announce: string;
}

export interface OrderedItem {
  ID: number;
  Name: string;
  Image: string;
  Description: string;
  Price: number;
  Quantity: number;
}

export interface OrderResponse {
  ID: number;
  Name: string;
  Email?: string;
  Phone: string;
  Items: OrderedItem[];
}

export interface GoodsUpdateRequest {
  ID: number;
  Name: string;
  Instruction: string;
  Description: string;
  IsPrescriptionNeeded: boolean;
  IsInStock: boolean;
  Price: number;
}

export interface PromoItem {
  id: number;
  name: string;
  description: string;
  price: number;
  image: string;
}

export interface UserUpdateRequest {
  ID: number;
  Login: string;
  UserName: string;
  RoleID: number;
  Role: Role;
  Password: string;
}

export interface UserResponse {
  ID: number;
  Login: string;
  UserName: string;
  RoleID: number;
  Role: RoleResponse;
  Password: string;
}

export interface RoleResponse {
  ID: number;
  Name: string;
  Permissions: PermissionRespons[];
}

export interface PermissionRespons {
  ID: number;
  Action: string;
}

export interface LogsResponse {
  Name: string;
}