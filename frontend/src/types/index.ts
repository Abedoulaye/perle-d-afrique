export interface User {
  id: number;
  email: string;
  role: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price_cents: number;
  stock: number;
  image: string;
}

export interface CartItem {
  id: number;
  cart_id: number;
  product_id: number;
  quantity: number;
  name?: string;
  price_cents?: number;
  description?: string;
  stock?: number;
  image?: string;
}

export interface Order {
  id: number;
  user_id: number;
  status: string;
  total_cents: number;
  created_at: string;
  shipping_name?: string;
  shipping_address?: string;
  shipping_phone?: string;
}

// Everything above this is from the go structs, as for this thing below it, its for pages/Cart.tsx
export interface CartDetail {
  id: number;
  cart_id: number;
  product_id: number;
  quantity: number;
  name?: string;
  price_cents?: number;
  description?: string;
  stock?: number;
  image?: string;
}
