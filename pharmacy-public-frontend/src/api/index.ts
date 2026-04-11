import axios, { type AxiosResponse } from "axios";
import type { GoodsResponse, AnnouncesResponse, KeyResponse, TaskResponse } from "@/types/api";
import type { Goods, Order, PromoItem, shedule } from "@/types";

export const api = axios.create({
	  baseURL: 'http://localhost:3000/api',
	  withCredentials: true,
});

const delays = [0, 500, 500, 2000, 2000, 2000, 10000, 10000, 10000]

const delay = (ms: number | undefined) => new Promise(resolve => setTimeout(resolve, ms));

async function executeWithPickup<T>(
  axiosRequest: Promise<AxiosResponse<KeyResponse>>
): Promise<T> {
  try {
    const axiosResponse = await axiosRequest;
    let taskKey = axiosResponse.data;
    
    if (!taskKey.TaskID) {
      throw new Error('No taskID');
    }

	const pickupResponse = await api.get<TaskResponse>(`/pickup?id=${taskKey.TaskID}`);
	let task = pickupResponse.data;	
	

	if (task.Status === 'completed') {
    	return JSON.parse(task.Value);
    }

    if (task.Status === 'error') {
    	throw new Error(task.Error);
    }

    if (task.Status === 'not_found') {
    	throw new Error(task.Error);
    }

	let attempt: number = 0
    while (task.Status === 'pending') {
        const waitTime = attempt <= delays.length ? delays[attempt] : delays[delays.length - 1];

        await delay(waitTime);

        const pickupResponse = await api.get<TaskResponse>(`/pickup?id=${taskKey.TaskID}`);
        task = pickupResponse.data;

        if (task.Status === 'completed') {
          return JSON.parse(task.Value);
        }

        if (task.Status === 'error') {
          throw new Error(task.Error);
        }

        if (task.Status === 'not_found') {
          throw new Error(task.Error);
        }

        attempt++;
      }
        throw new Error(`Unexpected task status after polling: ${task.Status}`);
	} catch (error) {
    	console.error('Error:', error);
    	throw error;
  	}
}
export const goodsAPI = {
  getGoods: async (searchQuery: string, page: string, limit: string): Promise<GoodsResponse> => {
    const query = searchQuery !== undefined ? searchQuery : "";
    const response = await executeWithPickup<GoodsResponse>(
      api.get<KeyResponse>(`/goods?q=${query}&page=${page}&limit=${limit}`)
    );

    if (response.Items !== null) {
      response.Items = response.Items.map(item => ({
        ...item,
        Type: "Goods"
      }));
    }

    return response;
  },

  getGood: async (id: number): Promise<Goods> => {
    const response = await executeWithPickup<Goods>(
      api.get<KeyResponse>(`/goods/?id=${id}`)
    );
    response.Type = 'Goods';
    return response;
  },

  getPromoItem: async (): Promise<PromoItem[]> => {
    return await executeWithPickup<PromoItem[]>(
      api.get<KeyResponse>(`/goods/advert`)
    );
  },
}

export const scheduleAPI = {
  getShedule: async (startDate: string, endDate: string): Promise<shedule[]> => {
    return await executeWithPickup<shedule[]>(
      api.get<KeyResponse>(`/schedule?start=${startDate}&end=${endDate}`)
    );
  }
}

export const orderAPI = {
  createOrder: async (order: Order): Promise<void> => {
    await executeWithPickup<void>(
      api.post<KeyResponse>(`/order`, order)
    );
  },
}

export const announcesAPI = {
  getAnnounces: async (searchQuery: string, page: string, limit: string): Promise<AnnouncesResponse> => {
    const query = searchQuery !== undefined ? searchQuery : "";
    const response = await executeWithPickup<AnnouncesResponse>(
      api.get<KeyResponse>(`/announces?q=${query}&page=${page}&limit=${limit}`)
    );

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