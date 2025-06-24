import { AuthResponse, User } from "@/types/user"
import { Api } from "./api"

async function login(email: string, password: string): Promise<AuthResponse> {
    return Api.post("/auth/login", { email, password })
}

async function register(firstName: string, lastName: string, username: string, email: string, password: string): Promise<AuthResponse> {
    console.log("Mandando info..,")
    console.log("first name: ", firstName, "lastname: ", lastName, " username: ", username, "email: ", email, " password: ", password)
    return Api.post("/auth/register", { firstName, lastName, username, email, password })
}

async function getUser(id: number): Promise<User> {
    const response = await Api.get(`/users/user/${id}`)
    return response.data
}

const userService = {
    login,
    register,
    getUser,
}

export { userService }