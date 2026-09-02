export interface LoginCredentials {
  email: string
  password: string
}

export interface AuthUser {
  id: number
  email: string
  name?: string
}

export interface AuthResponse {
  token: string
  user: AuthUser
}
