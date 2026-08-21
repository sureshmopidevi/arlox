import axios from 'axios'

export function handleApiError(error: unknown): never {
  if (axios.isAxiosError(error)) {
    const message =
      (error.response?.data as { message?: string })?.message ?? error.message
    throw new Error(message)
  }
  throw error
}
