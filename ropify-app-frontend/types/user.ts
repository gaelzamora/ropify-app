import { ApiResponse } from "./api"

export type AuthResponse = ApiResponse<{ token: string, user: User }>


export type User = {
    id: number;
    username: string;
    firstName: string;
    lastName: string;
    email: string;
    avatar_url?: string;
    bio?: string;
    google_id?: string
    created_at: string;
};