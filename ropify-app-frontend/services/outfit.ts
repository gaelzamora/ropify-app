import { GeneratedOutfitResponse, OutfitListResponse, OutfitResponse } from "@/types/outfit";
import { Api } from "./api";
import { GarmentOptimized } from "@/types/garment";

async function generateRandomOutfit(
    closetId: string,
    save?: boolean
): Promise<GeneratedOutfitResponse> {
    return await Api.post(`/closet/${closetId}/generate-outfit`, {
        params: { save }
    })

}

async function createOutfit(
    garments: GarmentOptimized[],
    closet_id?: string,
    name?: string,
): Promise<OutfitResponse> {
    return await Api.post(`/outfit`, { closet_id, name, garments })
}

async function getOutfits(
    closet_id: string
): Promise<OutfitListResponse> {
    console.log(closet_id)
    return await Api.get(`/closet/${closet_id}/outfits`)
}

export const outfitService = {
    generateRandomOutfit,
    createOutfit,
    getOutfits
}