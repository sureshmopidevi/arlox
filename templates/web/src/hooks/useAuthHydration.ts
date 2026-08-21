import { useEffect, useState } from 'react'
import { useAuthStore } from '@/stores/authStore'

export function useAuthHydration() {
  const [hydrated, setHydrated] = useState(
    () => useAuthStore.persist.hasHydrated(),
  )

  useEffect(() => {
    if (hydrated) return
    const unsub = useAuthStore.persist.onFinishHydration(() =>
      setHydrated(true),
    )
    // Re-check in case hydration completed between render and effect
    if (useAuthStore.persist.hasHydrated()) setHydrated(true)
    return unsub
  }, [hydrated])

  return hydrated
}
