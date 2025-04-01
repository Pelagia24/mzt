import User from "../models/User.ts";

export interface AuthResponse {
    access_token: 'string';
}

export interface LogoutResponse {
    acknowledged: boolean;
    deletedCount: number;
}

export interface UserInfoResponse {
    user: User;
}