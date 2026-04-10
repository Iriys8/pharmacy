import type { Announce, Goods, Order, PromoItem, WorkTime } from "@/types"

export interface GoodsResponse {
	Items: Goods[];
	TotalPages: number;
	CurrentPage: number;
}

export interface AnnouncesResponse {
	Items: Announce[];
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