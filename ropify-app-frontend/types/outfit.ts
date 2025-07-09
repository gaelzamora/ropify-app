import { ApiResponse } from "./api";
import { Garment } from "./garment"; // Asegúrate de importar el tipo Garment si ya existe

// Tipos para outfits individuales y listas
export type OutfitResponse = ApiResponse<{ data: Outfit }>
export type OutfitListResponse = ApiResponse<Outfit[]>

// Tipo específico para la respuesta del generador de outfits
export type GeneratedOutfitResponse = ApiResponse<{
    outfit: Outfit;
    garments: Garment[];
}>

// Actualiza el tipo Outfit para incluir todos los campos
export type Outfit = {
    id: string
    user_id: string
    name: string
    garment_ids: string[]
    tags: string[] | null
    occasion: string
    season: string
    archived: boolean
    image_url: string // Campo nuevo
    created_at: string // Campo nuevo
}