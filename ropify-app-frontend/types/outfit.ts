import { ApiResponse } from "./api";
import { Garment } from "./garment"; 

export type OutfitResponse = ApiResponse<{ data: Outfit }>
export type OutfitListResponse = ApiResponse<Outfit[]>

export type GeneratedOutfitResponse = ApiResponse<{
    outfit: Outfit;
    garments: Garment[];
}>

export type Outfit = {
    id: string
    user_id: string
    name: string
    garment_ids: string[]
    tags: string[] | null
    occasion: string
    season: string
    archived: boolean
    image_url: string 
    created_at: string
}

export type OutfitGenerateData = {
    garments: Garment[]
    outfit: Outfit
}