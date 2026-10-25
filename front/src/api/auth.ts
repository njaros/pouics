/**
 * File containing the API call functions to login or register
 */

import type { AuthResponse, LoginRequest, RegisterRequest } from "../types/api";
import { request } from "./client";

export async function register(body: RegisterRequest): Promise<AuthResponse> {
    return request("register", {
        method: "POST",
        body: JSON.stringify(body)
    })
}

export async function login(body: LoginRequest): Promise<AuthResponse> {
    return request("login", {
        method: "POST",
        body: JSON.stringify(body)
    })
}
