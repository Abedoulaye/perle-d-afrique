import { apiFetch } from "./api";
import type { Product } from "../types";

export async function getProducts(): Promise<Product[]> {
  return apiFetch("/products");
}

export async function getProduct(id: number): Promise<Product> {
  return apiFetch(`/products/${id}`);
}

export async function createProduct(
  product: Omit<Product, "id">, // this paremeter is Product type with "id" from the interface removed, since we are creating a new product we don't mention id, the backend does that
): Promise<Product> {
  return apiFetch("/products", {
    method: "POST",
    body: JSON.stringify(product),
  });
}

export async function updateProduct(
  id: number,
  product: Omit<Product, "id">,
): Promise<Product> {
  return apiFetch(`/products/${id}`, {
    method: "PUT",
    body: JSON.stringify(product),
  });
}

export async function deleteProduct(id: number): Promise<void> {
  return apiFetch(`/products/${id}`, {
    method: "DELETE",
  });
}
