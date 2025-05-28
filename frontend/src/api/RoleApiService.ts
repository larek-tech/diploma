import axiosInstance from './axiosInstance';
import {
  CreateRoleRequest,
  CreateUserRequest,
  ListDomainsResponse,
  ListRolesResponse,
  ListSourcesResponse,
  ListUsersResponse,
  PermittedRoles,
  PermittedUsers,
  RoleForEntity,
  UpdateRoleRequest,
  User,
} from './models/roles';

export class RoleApiService {
    // Управление ролями
    static async createRole(data: CreateRoleRequest): Promise<RoleForEntity> {
        const response = await axiosInstance.post<RoleForEntity>('/api/v1/role/', data);
        return response.data;
    }

    static async listRoles(offset: number, limit: number): Promise<ListRolesResponse> {
        const response = await axiosInstance.get<ListRolesResponse>('/api/v1/role/list', {
            params: { offset, limit },
        });
        return response.data;
    }

    static async getRole(id: number): Promise<RoleForEntity> {
        const response = await axiosInstance.get<RoleForEntity>(`/api/v1/role/${id}`);
        return response.data;
    }

    static async deleteRole(id: number): Promise<void> {
        await axiosInstance.delete(`/api/v1/role/${id}`);
    }

    static async setRole(data: UpdateRoleRequest): Promise<void> {
        await axiosInstance.put('/api/v1/role/set', data);
    }

    static async removeRole(data: UpdateRoleRequest): Promise<void> {
        await axiosInstance.put('/api/v1/role/remove', data);
    }

    // Управление разрешениями доменов
    static async getDomainPermittedRoles(domainId: number): Promise<PermittedRoles> {
        const response = await axiosInstance.get<PermittedRoles>(
            `/api/v1/domain/permissions/roles/${domainId}`
        );
        return response.data;
    }

    static async updateDomainPermittedRoles(
        domainId: number,
        data: PermittedRoles
    ): Promise<PermittedRoles> {
        const response = await axiosInstance.put<PermittedRoles>(
            `/api/v1/domain/permissions/roles/${domainId}`,
            data
        );
        return response.data;
    }

    static async getDomainPermittedUsers(domainId: number): Promise<PermittedUsers> {
        const response = await axiosInstance.get<PermittedUsers>(
            `/api/v1/domain/permissions/users/${domainId}`
        );
        return response.data;
    }

    static async updateDomainPermittedUsers(
        domainId: number,
        data: PermittedUsers
    ): Promise<PermittedUsers> {
        const response = await axiosInstance.put<PermittedUsers>(
            `/api/v1/domain/permissions/users/${domainId}`,
            data
        );
        return response.data;
    }

    // Управление разрешениями источников
    static async getSourcePermittedRoles(sourceId: number): Promise<PermittedRoles> {
        const response = await axiosInstance.get<PermittedRoles>(
            `/api/v1/source/permissions/roles/${sourceId}`
        );
        return response.data;
    }

    static async updateSourcePermittedRoles(
        sourceId: number,
        data: PermittedRoles
    ): Promise<PermittedRoles> {
        const response = await axiosInstance.put<PermittedRoles>(
            `/api/v1/source/permissions/roles/${sourceId}`,
            data
        );
        return response.data;
    }

    static async getSourcePermittedUsers(sourceId: number): Promise<PermittedUsers> {
        const response = await axiosInstance.get<PermittedUsers>(
            `/api/v1/source/permissions/users/${sourceId}`
        );
        return response.data;
    }

    static async updateSourcePermittedUsers(
        sourceId: number,
        data: PermittedUsers
    ): Promise<PermittedUsers> {
        const response = await axiosInstance.put<PermittedUsers>(
            `/api/v1/source/permissions/users/${sourceId}`,
            data
        );
        return response.data;
    }

    // Методы для получения списков пользователей и источников
    static async listUsers(offset: number = 0, limit: number = 100): Promise<ListUsersResponse> {
        const response = await axiosInstance.get<ListUsersResponse>('/api/v1/user/list', {
            params: { offset, limit },
        });
        return response.data;
    }

    static async getSources(offset: number = 0, limit: number = 100): Promise<ListSourcesResponse> {
        const response = await axiosInstance.get<ListSourcesResponse>('/api/v1/source/list', {
            params: { offset, limit },
        });
        return response.data;
    }

    // Назначение роли пользователю
    static async updateUserRole(userId: number, roleId: number): Promise<void> {
        const data: UpdateRoleRequest = {
            userId,
            roleId,
        };
        await axiosInstance.put('/api/v1/role/set', data);
    }

    // Удаление роли у пользователя
    static async removeUserRole(userId: number, roleId: number): Promise<void> {
        const data: UpdateRoleRequest = {
            userId,
            roleId,
        };
        await axiosInstance.put('/api/v1/role/remove', data);
    }

    // Управление пользователями
    static async createUser(data: CreateUserRequest): Promise<User> {
        const response = await axiosInstance.post<User>('/api/v1/user/', data);
        return response.data;
    }

    // Управление доменами
    static async getDomains(offset: number = 0, limit: number = 100): Promise<ListDomainsResponse> {
        const response = await axiosInstance.get<ListDomainsResponse>('/api/v1/domain/list', {
            params: { offset, limit },
        });
        return response.data;
    }
}
