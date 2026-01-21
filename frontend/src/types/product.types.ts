export interface Product {
  id: string;
  sku: string;
  name: string;
  brand: string;
  model: string;
  color: string;
  description: string;
  category: string;
  releaseYear: string;
  retailPrice: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  images?: ProductImage[];
  sizes?: Size[];
}

export interface ProductImage {
  id: string;
  productId: string;
  imageUrl: string;
  displayOrder: number;
  isPrimary: boolean;
  createdAt: string;
}

export interface Size {
  id: string;
  productId: string;
  size: string;
  quantity: number;
  reserved: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductsResponse {
  products: Product[];
  total: string;
  page: number;
  pageSize: number;
}

export interface ProductsRequest {
  page?: number;
  pageSize?: number;
  category?: string;
  brand?: string;
}

export interface SearchProductsRequest {
  query: string;
  page?: number;
  pageSize?: number;
}
