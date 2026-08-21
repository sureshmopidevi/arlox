import { apiClient } from '@/api/client'
import type { AuthResponse, LoginCredentials } from '../types'

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    const data = await apiClient.post('/auth/login', credentials)
    return data as unknown as AuthResponse
  },
}
