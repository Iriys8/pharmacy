import axios, { type AxiosResponse } from "axios";
import type { AuthResponse, LoginData } from "@/types/auth";
import type { GoodsResponse, GoodsUpdateRequest, ScheduleResponse, OrdersResponse, AnnouncesResponse, UsersResponse, UserResponse, RolesResponse, PermissionsResponse, LogsResponse, TaskResponse, KeyResponse } from "@/types/api";
import type { Announce, Goods, Order, Role, User, Schedule } from "@/types";

export const api = axios.create({
	  baseURL: 'http://localhost:3001/api',
	  withCredentials: true,
});

const delays = [0, 500, 500, 2000, 2000, 2000, 10000, 10000, 10000]

const delay = (ms: number | undefined) => new Promise(resolve => setTimeout(resolve, ms));

async function executeWithPickup<T>(axiosRequest: Promise<AxiosResponse<KeyResponse>>): Promise<T> {
  try {
    const axiosResponse = await axiosRequest;
    let taskKey = axiosResponse.data;
    
    if (!taskKey.TaskID) {
      throw new Error('No taskID');
    }

	const pickupResponse = await api.get<TaskResponse>(`/pickup?key=${taskKey.TaskID}`);
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

        const pickupResponse = await api.get<TaskResponse>(`/pickup?key=${taskKey.TaskID}`);
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
		const response = await executeWithPickup<GoodsResponse>(api.get<KeyResponse>(`/goods?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`))

		if (response.Items !== null) {
    		response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Goods"
    		}));
		}

    	return response;
	},

	getGood: async (id: number): Promise<Goods> => {
		return await executeWithPickup<Goods>(api.get<KeyResponse>((`/goods?id=${id}`)))
	},

	updateGood: async (good: GoodsUpdateRequest) => {
    	await api.patch(`/goods?id=${good.ID}`, good)
	}
}

export const scheduleAPI = {
	getSchedule: async (searchQuery: string, page: string, limit: string): Promise<ScheduleResponse> => {
		const response = await executeWithPickup<ScheduleResponse>(api.get<KeyResponse>(`/schedule?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`))
		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Schedule"
    		}));
		}
    	return response;
	},

	getScheduleByID: async (id: number): Promise<Schedule> => {
		return await executeWithPickup<Schedule>(api.get<KeyResponse>(`/schedule?id=${id}`))
	},

	updateScheduleByID: async (sheduleDate: Schedule) => {
    	await api.patch(`/schedule?id=${sheduleDate.ID}`, sheduleDate)
	},

	createSchedule: async (sheduleDate: Schedule) => {
    	await api.post(`/schedule`, sheduleDate)
	},

	deleteScheduleByID: async (sheduleDate: Schedule) => {
    	await api.delete(`/schedule?id=${sheduleDate.ID}`)
	}
}

export const orderAPI = {
	getOrders: async (searchQuery: string, page: string, limit: string): Promise<OrdersResponse> => {

		const response = await executeWithPickup<OrdersResponse>(api.get<KeyResponse>(`/order?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`))
		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Order"
    		}));
		}
    	return response;
	},

	getOrderByID: async (id: number): Promise<Order> => {
		const response = await executeWithPickup<Order>(api.get<KeyResponse>(`/order/${id}`))
		response.Items = response.Items.map(item => ({
    	  ...item,
    	  Type: "OrderedItem"
    	}));
		return response
	},

	updateOrderByID: async (order: Order) => {
    	await api.patch(`/order?id=${order.ID}`, order)
	},

	createOrder: async (order: Order) => {
    	await api.post(`/order`, order)
	},

	deleteOrderByID: async (order: Order) => {
    	await api.delete(`/order?id=${order.ID}`)
	}
}

export const announcesAPI = {
	getAnnounces: async (searchQuery: string, page: string, limit: string): Promise<AnnouncesResponse> => {
		const response = await executeWithPickup<AnnouncesResponse>(api.get<KeyResponse>(`/announce?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`))

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Announce"
    		}));
		}

    	return response;
	},

	getAnnounceByID: async (id: number): Promise<Announce> => {
		return (await executeWithPickup<Announce>(api.get<KeyResponse>(`/announce?id=${id}`)))
	},

	updateAnnounceByID: async (announce: Announce) => {
    	await api.patch(`/announce?id=${announce.ID}`, announce)
	},

	createAnnounce: async (announce: Announce) => {
    	await api.post(`/announce`, announce)
	},

	deleteAnnounceByID: async (announce: Announce) => {
    	await api.delete(`/announce?id=${announce.ID}`)
	}
}

export const usersAPI = {
	getUsers: async (searchQuery: string, page: string, limit: string): Promise<UsersResponse> => {
		const response = (await executeWithPickup<UsersResponse>(api.get<KeyResponse>(`/user?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)))

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "User"
    		}));
		}

    	return response;
	},

	getUserByID: async (id: number): Promise<UserResponse> => {
		return (await executeWithPickup<UserResponse>(api.get<KeyResponse>(`/user?id=${id}`)))
	},

	updateUserByID: async (user: User) => {
    	await api.patch(`/user?id=${user.ID}`, user)
	},

	createUser: async (user: User) => {
    	await api.post(`/user`, user)
	},

	deleteUserByID: async (user: User) => {
    	await api.delete(`/user?id=${user.ID}`)
	}
}

export const rolesAPI = {
	getRoles: async (searchQuery: string, page: string, limit: string): Promise<RolesResponse> => {
		const response = (await executeWithPickup<RolesResponse>(api.get<KeyResponse>(`/role?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)))

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Role"
    		}));
		}

    	return response;
	},

	getRoleByID: async (id: number): Promise<Role> => {
		const response = (await executeWithPickup<Role>(api.get<KeyResponse>(`/role?id=${id}`)))

		if (response.Permissions !== null) {
			response.Permissions = response.Permissions.map(item => ({
    		  ...item,
    		  Type: "Permission"
    		}));
		}

		return response
	},

	updateRoleByID: async (role: Role) => {
    	await api.patch(`/role?id=${role.ID}`, role)
	},

	createRole: async (role: Role) => {
    	await api.post(`/role`, role)
	},

	deleteRoleByID: async (role: Role) => {
    	await api.delete(`/role?id=${role.ID}`)
	},

	getPermissions: async (searchQuery: string, page: string, limit: string): Promise<PermissionsResponse> => {
		const response = (await executeWithPickup<PermissionsResponse>(api.get<KeyResponse>(`/permission?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)))

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Permission"
    		}));
		}

    	return response;
	},
}

export const imagesAPI = {
	getImageSRC: (): string => {
		return "http://localhost:3001/api/image?name=";
	}
}

export const logsAPI = {
	getLogs: async (searchQuery: string, page: string, limit: string): Promise<LogsResponse> => {
    	const response = (await api.get<LogsResponse>(`/logs?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data;
		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Log"
    		}));
		}

    	return response;
	},

	getLog: async (name: string) => {
    	await api.get(`/log?name=${name}`).then((res) => {
    		const url = window.URL.createObjectURL(new Blob([res.data]));
    		const link = document.createElement("a");
    		link.href = url;
    		link.setAttribute("download", name);
    		document.body.appendChild(link);
    		link.click();
    		link.remove();
  		})
	},
}

export const authAPI = {
	login: async (loginData: LoginData): Promise<AuthResponse> => {
		return (await axios.create({
	  		baseURL: 'http://localhost:3001/api',
	  		withCredentials: true,
		}).post<AuthResponse>('http://localhost:3001/api/login', loginData)).data
	},
	
	logout: async () => {
		return await api.post('/logout');
	},

	refreshToken: async (): Promise<AuthResponse> => {
    	const response = await axios.post<AuthResponse>( '/api/refresh',
    	  {},
    	  { 
    	    baseURL: 'http://localhost:3001',
    	    withCredentials: true 
    	  }
    	)
		return response.data
	}
};