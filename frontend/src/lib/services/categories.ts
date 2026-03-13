import { apiJson, apiFetch } from './api';

export interface CategoryType {
    id: number;
    name: string;
    code: string;
    description: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface Category {
    id: number;
    type_id: number;
    parent_id?: number;
    type?: CategoryType;
    parent?: Category;
    name: string;
    color: string;
    icon: string;
    description: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    limit: number;
}

export interface CategoryTypeCreate {
    name: string;
    code: string;
    description?: string;
    is_active?: boolean;
}

export interface CategoryCreate {
    type_id: number;
    parent_id?: number;
    name: string;
    color?: string;
    icon?: string;
    description?: string;
    is_active?: boolean;
}

class CategoryService {
    private typeBaseUrl = '/api/v1/category-types';
    private catBaseUrl = '/api/v1/categories';

    // --- Category Types ---

    async listTypes(page: number = 1, limit: number = 10, search?: string): Promise<PaginatedResponse<CategoryType>> {
        let url = `${this.typeBaseUrl}?page=${page}&limit=${limit}`;
        if (search) {
            url += `&search=${encodeURIComponent(search)}`;
        }
        return await apiJson<PaginatedResponse<CategoryType>>(url) || { data: [], total: 0, page, limit };
    }

    async createType(data: CategoryTypeCreate): Promise<CategoryType> {
        return await apiJson<CategoryType>(this.typeBaseUrl, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async updateType(id: number, data: Partial<CategoryTypeCreate>): Promise<CategoryType> {
        return await apiJson<CategoryType>(`${this.typeBaseUrl}/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async deleteType(id: number): Promise<void> {
        await apiFetch(`${this.typeBaseUrl}/${id}`, { method: 'DELETE' });
    }

    // --- Categories ---

    async listCategories(page: number = 1, limit: number = 10, search?: string, typeId?: number): Promise<PaginatedResponse<Category>> {
        let url = `${this.catBaseUrl}?page=${page}&limit=${limit}`;
        if (search) {
            url += `&search=${encodeURIComponent(search)}`;
        }
        if (typeId) {
            url += `&type_id=${typeId}`;
        }
        return await apiJson<PaginatedResponse<Category>>(url) || { data: [], total: 0, page, limit };
    }

    async getCategoryById(id: number): Promise<Category> {
        return await apiJson<Category>(`${this.catBaseUrl}/${id}`);
    }

    async createCategory(data: CategoryCreate): Promise<Category> {
        return await apiJson<Category>(this.catBaseUrl, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async updateCategory(id: number, data: Partial<CategoryCreate>): Promise<Category> {
        return await apiJson<Category>(`${this.catBaseUrl}/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }

    async deleteCategory(id: number): Promise<void> {
        await apiFetch(`${this.catBaseUrl}/${id}`, { method: 'DELETE' });
    }
}

export const categoryService = new CategoryService();
