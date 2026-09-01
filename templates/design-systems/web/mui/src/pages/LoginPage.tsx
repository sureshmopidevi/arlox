import { Box } from '@mui/material'
import { LoginForm } from '@/features/auth/components/LoginForm'

export function LoginPage() {
  return (
    <Box sx={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: 'grey.50', px: 2 }}>
      <Box sx={{ width: '100%', maxWidth: 420 }}>
        <LoginForm />
      </Box>
    </Box>
  )
}
