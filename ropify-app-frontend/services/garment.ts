import { GarmentListResponse } from "@/types/garment";
import { Api } from "./api";

async function filterGarments(
    page: number, 
    limit: number, 
    user_id?: string,
    category?: string,
    color?: string, 
    brand?: string,
): Promise<GarmentListResponse> {
    if (category === "all") {
        category = ""
    }

    return Api.get("/garment", {
        params: { page, limit, color, brand, category, user_id }
    })
}

export const garmentService = {
    filterGarments
}