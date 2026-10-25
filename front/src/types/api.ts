import { type User } from "./user"

export interface RegisterRequest {
    name: string,
    password: string
}

export interface LoginRequest extends RegisterRequest {}

export interface AuthResponse {
    user: User,
    token: string
}