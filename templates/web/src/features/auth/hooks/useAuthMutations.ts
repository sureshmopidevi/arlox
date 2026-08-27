import { useMutation } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/authStore'
import { authApi } from '../api/authApi'

export function useLoginMutation() {
  const setAuth = useAuthStore((s) => s.setAuth)

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.token, data.user)
    },
  })
}
