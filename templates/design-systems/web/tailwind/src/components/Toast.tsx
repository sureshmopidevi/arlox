import { useEffect, useState } from 'react'

export type ToastType = 'success' | 'error' | 'info'

interface ToastProps {
  message: string
  type?: ToastType
  duration?: number
  onClose?: () => void
}

const styles: Record<ToastType, string> = {
  success: 'bg-green-600 text-white',
  error: 'bg-red-600 text-white',
  info: 'bg-gray-800 text-white',
}

export function Toast({
  message,
  type = 'info',
  duration = 3500,
  onClose,
}: ToastProps) {
  const [visible, setVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false)
      onClose?.()
    }, duration)
    return () => clearTimeout(timer)
  }, [duration, onClose])

  if (!visible) return null

  return (
    <div
      role="status"
      className={[
        'fixed bottom-4 right-4 z-50 flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg text-sm max-w-sm',
        styles[type],
      ].join(' ')}
    >
      <span className="flex-1">{message}</span>
      <button
        onClick={() => {
          setVisible(false)
          onClose?.()
        }}
        className="opacity-70 hover:opacity-100 transition-opacity shrink-0"
        aria-label="Dismiss"
      >
        ✕
      </button>
    </div>
  )
}
