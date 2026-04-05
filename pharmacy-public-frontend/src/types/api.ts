import type { Announce, Goods, Order } from "@/types"

export interface GoodsResponse {
	Items: Goods[];
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