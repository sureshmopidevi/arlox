import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Alert, Button, Paper, PasswordInput, Stack, Text, TextInput, Title } from '@mantine/core'
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
    <Paper p="xl" radius="md" withBorder>
      <Title order={3}>Sign in</Title>
      <Text c="dimmed" size="sm" mb="md">
        Enter your credentials to continue.
      </Text>
      <form onSubmit={handleSubmit}>
        <Stack gap="md">
          <TextInput
            label="Email"
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <PasswordInput
            label="Password"
            autoComplete="current-password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {error && <Alert color="red">{error.message}</Alert>}
          <Button type="submit" loading={isPending} fullWidth>
            Sign in
          </Button>
        </Stack>
      </form>
    </Paper>
  )
}
