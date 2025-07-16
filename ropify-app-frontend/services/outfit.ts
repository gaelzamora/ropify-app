import { GeneratedOutfitResponse } from "@/types/outfit";
import { Api } from "./api";

async function generateRandomOutfit(
    closetId: string,
    save?: boolean
): Promise<GeneratedOutfitResponse> {
    return await Api.post(`/closet/${closetId}/generate-outfit`, {
        params: { save }
    })

}

export const outfitService = {
    generateRandomOutfit
}