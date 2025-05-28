// API types for roles management

import {Domain, Source} from './domain';

export interface RoleForEntity {
    id: number;
    name: string;
    createdAt: {
        seconds: number;
        nanos: number;
    };
}

export interface CreateRoleRequest {
    name: string;
}

export interface UpdateRoleRequest {
    roleId: number;
    userId: number;
}

export interface ListRolesResponse {
    roles: RoleForEntity[];
}

export interface PermittedRoles {
    resourceId: number;
    roleIds: number[];
}

export interface PermittedUsers {
    resourceId: number;
    userIds: number[];
}

export interface GetResourcePermissionsRequest {
    resourceId: number;
}

// User interfaces
export interface User {
    id: number;
    login: string;
    firstName: string;
    lastName: string;
    email: string;
    roles: number[] | null;
    createdAt: {
        seconds: number;
        nanos: number;
    };
}

export interface ListUsersResponse {
    users: User[];
}

// User creation interface
export interface CreateUserRequest {
    email: string;
    password: string;
}

export interface ListSourcesResponse {
    sources: Source[];
}

export interface ListDomainsResponse {
    domains: Domain[];
}

// Domain permissions interfaces
export interface DomainPermissions {
    domainId: number;
    userIds: number[];
    roleIds: number[];
}

export interface UpdateDomainPermissionsRequest {
    domainId: number;
    userIds?: number[];
    roleIds?: number[];
}
