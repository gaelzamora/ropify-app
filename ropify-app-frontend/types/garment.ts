import { ApiResponse } from "./api"

export type GarmentResponse = ApiResponse<{ message: string, data: Garment}>
export type GarmentListResponse = ApiResponse<Garment[]>

export type Garment = {
    id: string
    user_id: string
    name: string
    category: string
    color: string
    brand: string
    size: string
    image_url: string
    barcode: string
    is_verified: string
}