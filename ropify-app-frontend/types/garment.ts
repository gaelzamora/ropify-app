import { ApiResponse } from "./api"

export type GarmentResponse = ApiResponse<{ message: string, data: Garment}>
export type GarmentListResponse = ApiResponse<Garment[]>

export type Garment = {
    id: number
    user_id: number
    name: string
    category: string
    color: string
    brand: string
    size: string
    image_url: string
    barcode: string
    is_verified: string
}