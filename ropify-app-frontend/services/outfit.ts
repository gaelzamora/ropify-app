import { GeneratedOutfitResponse } from "@/types/outfit";
import { Api } from "./api";

async function generateRandomOutfit(
    save?: boolean
): Promise<GeneratedOutfitResponse> {
    return await Api.get
}