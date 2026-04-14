export interface Goods {
    ID: number;
    Name: string;
    Image: string;
    Instruction?: string;
    Description?: string;
    IsPrescriptionNeeded: boolean;
	IsInStock: boolean;
    Price: number;
    Producer: string;
	Tags?: string[];
    Type: "Goods"
}

export interface Schedule {
	ID: number;
	Date: string;
    TimeStart: string | null;
    TimeEnd: string | null;
	IsOpened: boolean;
    Type: "Schedule";
}

export interface Order {
    ID: number;
    Name: string;
    Email?: string;
    Phone: string;
    Items: OrderedItem[];
    Type: "Order";
}

export interface OrderedItem {
    ID: number;
	Name: string;
	Image: string;
	Description: string;
	Price: number;
	Quantity: number;
    Type: "OrderedItem";
}

export interface Announce {
    ID: number;
	DateTime: string;
	From: string;
	Announce: string;
    Type: "Announce"
}

export interface User {
    ID: number;
    Login: string;
    UserName: string;
    RoleID: number;
    Role: Role | undefined;
    Password: string;
    Type: "User"
}

export interface Role {
    ID: number;
    Name: string;
    Permissions: Permission[];
    Type: "Role"
}

export interface Permission {
    ID: number;
    Action: string;
    Type: "Permission"
}

export interface Log {
    Name: string;
    Type: "Log"
}