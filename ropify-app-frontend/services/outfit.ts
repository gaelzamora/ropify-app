import { GeneratedOutfitResponse } from "@/types/outfit";
import { Api } from "./api";

async function generateRandomOutfit(
    closetId: string,
    save?: boolean
): Promise<GeneratedOutfitResponse> {
    return await Api.get(`/outfit/generate/${closetId}`, {
        params: { save }
    })

}

export const outfitService = {
    generateRandomOutfit
}