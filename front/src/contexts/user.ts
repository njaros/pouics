import { createContext } from "react";
import { type User } from "../types/user"

export const UserContext = createContext<User | null>(null)