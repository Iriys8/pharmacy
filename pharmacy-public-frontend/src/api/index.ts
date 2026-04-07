import axios from "axios";
import type { GoodsResponse, AnnouncesResponse } from "@/types/api";
import type { Goods, Order, OrderedItem, PromoItem, WorkTime } from "@/types";

export const api = axios.create({
	  baseURL: 'http://localhost:3000/api',
	  withCredentials: true,
});

export const goodsAPI = {
	getGoods: async (searchQuery: string, page: string, limit: string): Promise<GoodsResponse> => {
		const response = (await api.get<GoodsResponse>(`/goods?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data;

		if (response.Items !== null) {
    		response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Goods"
    		}));
		}

    	return response;
	},

	getGood: async (id: number): Promise<Goods> => {
		const response = (await api.get<Goods>(`/goods/${id}`)).data;
		response.Type = 'Goods';
		return response;
	},

	getPromoItem: async (): Promise<PromoItem[]> => {
		return (await api.get<PromoItem[]>(`/goods/advert`)).data
	},
}

export const scheduleAPI = {
	getShedule: async (startDate: string, endDate: string): Promise<WorkTime[]> => {
		return (await api.get<WorkTime[]>(`/schedule?start=${startDate}&end=${endDate}`)).data
	}
}

export const orderAPI = {
	createOrder: async (order: Order) => {
    	await api.post(`/order/${order.ID}`, order)
	},
}

export const announcesAPI = {
	getAnnounces: async (searchQuery: string, page: string, limit: string): Promise<AnnouncesResponse> => {
		const response = (await api.get<AnnouncesResponse>(`/announce?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Announce"
    		}));
		}

    	return response;
	},
}

export const imagesAPI = {
	getImageSRC: (): string => {
		return "http://localhost:3000/api/image?name=";
	}
}