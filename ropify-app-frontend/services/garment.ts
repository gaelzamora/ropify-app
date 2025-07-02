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

async function analyzeGarmentImage(imageUri: string) {
  const formData = new FormData();
  
  const filename = imageUri.split('/').pop();
  const match = /\.(\w+)$/.exec(filename || '');
  const type = match ? `image/${match[1]}` : 'image';
  
  formData.append('image', {
    uri: imageUri,
    name: filename,
    type,
  } as any);
  
  return Api.post('/garment/analyze', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

export const garmentService = {
    filterGarments,
    analyzeGarmentImage
}