import axios from 'axios'

export function handleApiError(error: unknown): never {
  if (axios.isAxiosError(error)) {
    const body = error.response?.data as { error?: string; message?: string } | undefined
    const message = body?.error ?? body?.message ?? error.message
    throw new Error(message)
  }
  throw error
}
