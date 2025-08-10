import { ApiResponse } from "./api"

export type GarmentResponse = ApiResponse<{ message: string, data: Garment}>
export type GarmentListResponse = ApiResponse<{"closet_name": string, "garments": Garment[]}>

export type Garment = {
    id: string
    user_id: string
    category: string
    color: string
    labels: string[]
    image_url: string
    is_verified: string
    boundingPoly?: {x: number, y:number}[]
}

export type GarmentOptimized = {
    id: string
    image_url: string
}