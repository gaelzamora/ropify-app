import { ApiResponse } from "./api";

export type ClosetResponse = ApiResponse<{ data: Closet }>
export type ClosetListResponse = ApiResponse<Closet[]>

export type Closet = {
    id: string
    user_id: string
    name: string
    image_url: string
    outfits_ids: string[]
    created_at: string
    updated_at: string
}