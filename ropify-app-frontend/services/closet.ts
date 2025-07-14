import { ClosetListResponse, ClosetResponse } from "@/types/closet";
import { Api } from "./api";

// Definir un tipo más apropiado para la entrada de archivos
type ClosetImageInput = {
  uri?: string;
  name?: string;
  type?: string;
} | Blob | File | any;

async function createOne(
    name: string,
    closetImage: ClosetImageInput // Usar el tipo personalizado
): Promise<ClosetResponse>{
    const formData = new FormData();
    formData.append('name', name);
    
    // Si es un objeto File/Blob del navegador
    if (closetImage && typeof closetImage === 'object') {
        if (closetImage.uri) {
            // Caso React Native
            formData.append('closet_image', {
                uri: closetImage.uri,
                name: closetImage.name || 'closet_image.jpg',
                type: closetImage.type || 'image/jpeg'
            } as any);
        } else {
            // Caso Web (File/Blob)
            formData.append('closet_image', closetImage);
        }
    }
    
    return Api.post("/closet", formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}


async function getMany(): Promise<ClosetListResponse> {
    return Api.get("/closet")
}

async function updateOne(
    id: string
): Promise<ClosetResponse> {
    return Api.put(`/closet/${id}`)
}

async function deleteOne(
    id: string
): Promise<Response> {
    return Api.delete(`/closet/${id}`)
}

export const closetService = {
    createOne,
    getMany,
    updateOne,
    deleteOne
}