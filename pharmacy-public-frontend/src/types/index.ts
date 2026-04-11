export interface CartState {
  items: Record<number, number>;
}

export interface shedule {
	ID: number;
	Date: string;
    TimeStart: string | null;
    TimeEnd: string | null;
	IsOpened: boolean;
}

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

export interface Order {
    Name: string;
    Email?: string;
    Phone: string;
    Items: OrderedItem[];
}

export interface OrderedItem extends Goods {
    Quantity: number;
}

export interface Announce {
    ID: number;
	DateTime: string;
	From: string;
	Announce: string;
    Type: "Announce"
}

export interface CartItemType extends Goods {
    quantity: number;
}

export interface CustomerType {
    name: string;
    email?: string;
    phone: string;
}

export interface PromoItem {
    id: number;
    name: string;
    description: string;
    price: number;
    image: string;
}