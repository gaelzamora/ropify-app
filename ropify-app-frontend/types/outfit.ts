import { ApiResponse } from "./api";
import { Garment, GarmentOptimized } from "./garment"; 

export type OutfitResponse = ApiResponse<{ data: Outfit }>
export type OutfitListResponse = ApiResponse<{ data: Outfits }>

export type GeneratedOutfitResponse = ApiResponse<{
    outfit: Outfit;
    garments: Garment[];
}>

export type Outfit = {
    id: string
    user_id: string
    name: string
    garments: GarmentOptimized[]
    occasion: string
    archived: boolean
    image_url: string 
    created_at: string
}

export type OutfitGenerateData = {
    garments: Garment[]
    outfit: Outfit
}

export type Outfits = {
    closet_name: string
    outfits: Outfit[]
}