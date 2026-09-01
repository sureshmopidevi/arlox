import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Alert, Box, Button, Paper, TextField, Typography } from '@mui/material'
import { useLoginMutation } from '../hooks/useAuthMutations'

export function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()
  const { mutate: login, isPending, error } = useLoginMutation()

  function handleSubmit(e: FormEvent) {
    e.preventDefault()
    login({ email, password }, { onSuccess: () => void navigate('/') })
  }

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h5" component="h1" gutterBottom>
        Sign in
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Enter your credentials to continue.
      </Typography>
      <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <TextField
          label="Email"
          type="email"
          autoComplete="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <TextField
          label="Password"
          type="password"
          autoComplete="current-password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {error && <Alert severity="error">{error.message}</Alert>}
        <Button type="submit" variant="contained" disabled={isPending} fullWidth>
          {isPending ? 'Signing in…' : 'Sign in'}
        </Button>
      </Box>
    </Paper>
  )
}
