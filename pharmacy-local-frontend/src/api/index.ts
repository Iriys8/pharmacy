import axios from "axios";
import type { AuthResponse, LoginData } from "@/types/auth";
import type { GoodsResponse, GoodsUpdateRequest, WorkTimesResponse, OrdersResponse, AnnouncesResponse, UsersResponse, UserResponse, RolesResponse, PermissionsResponse, LogsResponse } from "@/types/api";
import type { Announce, Goods, Order, Role, User, WorkTime } from "@/types";

export const api = axios.create({
	  baseURL: 'http://localhost:5001/api',
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
		return (await api.get<Goods>(`/goods/${id}`)).data
	},

	updateGood: async (good: GoodsUpdateRequest) => {
    	await api.patch(`/goods/${good.ID}`, good)
	}
}

export const scheduleAPI = {
	getSchedule: async (searchQuery: string, page: string, limit: string): Promise<WorkTimesResponse> => {
		const response = (await api.get<WorkTimesResponse>(`/schedule?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data
		
		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "WorkTime"
    		}));
		}

    	return response;
	},

	getScheduleByID: async (id: number): Promise<WorkTime> => {
		return (await api.get<WorkTime>(`/schedule/${id}`)).data
	},

	updateScheduleByID: async (sheduleDate: WorkTime) => {
    	await api.patch(`/schedule/${sheduleDate.ID}`, sheduleDate)
	},

	createSchedule: async (sheduleDate: WorkTime) => {
    	await api.post(`/schedule/${sheduleDate.ID}`, sheduleDate)
	},

	deleteScheduleByID: async (sheduleDate: WorkTime) => {
    	await api.delete(`/schedule/${sheduleDate.ID}`)
	}
}

export const orderAPI = {
	getOrders: async (searchQuery: string, page: string, limit: string): Promise<OrdersResponse> => {

		const response = (await api.get<OrdersResponse>(`/order?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data
		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Order"
    		}));
		}
    	return response;
	},

	getOrderByID: async (id: number): Promise<Order> => {
		const response = (await api.get<Order>(`/order/${id}`)).data
		response.Items = response.Items.map(item => ({
    	  ...item,
    	  Type: "OrderedItem"
    	}));
		return response
	},

	updateOrderByID: async (order: Order) => {
    	await api.patch(`/order/${order.ID}`, order)
	},

	createOrder: async (order: Order) => {
    	await api.post(`/order/${order.ID}`, order)
	},

	deleteOrderByID: async (order: Order) => {
    	await api.delete(`/order/${order.ID}`)
	}
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

	getAnnounceByID: async (id: number): Promise<Announce> => {
		return (await api.get<Announce>(`/announce/${id}`)).data
	},

	updateAnnounceByID: async (announce: Announce) => {
    	await api.patch(`/announce/${announce.ID}`, announce)
	},

	createAnnounce: async (announce: Announce) => {
    	await api.post(`/announce/${announce.ID}`, announce)
	},

	deleteAnnounceByID: async (announce: Announce) => {
    	await api.delete(`/announce/${announce.ID}`)
	}
}

export const usersAPI = {
	getUsers: async (searchQuery: string, page: string, limit: string): Promise<UsersResponse> => {
		const response = (await api.get<UsersResponse>(`/user?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "User"
    		}));
		}

    	return response;
	},

	getUserByID: async (id: number): Promise<UserResponse> => {
		return (await api.get<UserResponse>(`/user/${id}`)).data
	},

	updateUserByID: async (user: User) => {
    	await api.patch(`/user/${user.ID}`, user)
	},

	createUser: async (user: User) => {
    	await api.post(`/user/${user.ID}`, user)
	},

	deleteUserByID: async (user: User) => {
    	await api.delete(`/user/${user.ID}`)
	}
}

export const rolesAPI = {
	getRoles: async (searchQuery: string, page: string, limit: string): Promise<RolesResponse> => {
		const response = (await api.get<RolesResponse>(`/role?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data

		if (response.Items !== null) {
			response.Items = response.Items.map(item => ({
    		  ...item,
    		  Type: "Role"
    		}));
		}

    	return response;
	},

	getRoleByID: async (id: number): Promise<Role> => {
		const response = (await api.get<Role>(`/role/${id}`)).data

		if (response.Permissions !== null) {
			response.Permissions = response.Permissions.map(item => ({
    		  ...item,
    		  Type: "Permission"
    		}));
		}

		return response
	},

	updateRoleByID: async (role: Role) => {
    	await api.patch(`/role/${role.ID}`, role)
	},

	createRole: async (role: Role) => {
    	await api.post(`/role/${role.ID}`, role)
	},

	deleteRoleByID: async (role: Role) => {
    	await api.delete(`/role/${role.ID}`)
	},

	getPermissions: async (searchQuery: string, page: string, limit: string): Promise<PermissionsResponse> => {
		const response = (await api.get<PermissionsResponse>(`/permission?q=${searchQuery !== undefined ? searchQuery : ""}&page=${page}&limit=${limit}`)).data

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
		return "http://localhost:5001/api/image?name=";
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
	  		baseURL: 'http://localhost:5001/api',
	  		withCredentials: true,
		}).post<AuthResponse>('http://localhost:5001/api/login', loginData)).data
	},
	
	logout: async () => {
		return await api.post('/logout');
	},

	refreshToken: async (): Promise<AuthResponse> => {
    	const response = await axios.post<AuthResponse>( '/api/refresh',
    	  {},
    	  { 
    	    baseURL: 'http://localhost:5001',
    	    withCredentials: true 
    	  }
    	)
		return response.data
	}
};