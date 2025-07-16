import { GarmentListResponse } from "@/types/garment";
import { Api } from "./api";

async function filterGarmentsFromCloset(
    closet_id: string,
    page: number, 
    limit: number, 
    category?: string,
): Promise<GarmentListResponse> {
    if (category === "all") {
        category = ""
    }

    return Api.get(`/closet/${closet_id}/filter-garments`, {
        params: { page, limit, category }
    })
}

async function analyzeGarmentImage(
  closet_id: string,
  imageUri: string
) {
  const formData = new FormData();
  
  const filename = imageUri.split('/').pop();
  const match = /\.(\w+)$/.exec(filename || '');
  const type = match ? `image/${match[1]}` : 'image';
  
  formData.append('image', {
    uri: imageUri,
    name: filename,
    type,
  } as any);
  
  return Api.post(`/closet/${closet_id}/analyze-and-add`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

async function deleteMultipleGarmentsFromCloset(
  closet_id: string, 
  garment_ids: string[]
) {
  return Api.post(`/closet/${closet_id}/garments/remove`, { 
    garment_ids: garment_ids
  });
}

export const garmentService = {
    filterGarmentsFromCloset,
    analyzeGarmentImage,
    deleteMultipleGarmentsFromCloset,
}