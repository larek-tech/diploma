export interface LoginParams {
    email: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    meta: {
        roles: Role[];
        userId: number;
    };
    type: string;
}

export enum Role {
    Admin = 'ADMIN',
    Analyst = 'ANALYST',
    Buyer = 'BUYER',
}
