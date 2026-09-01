import { Center, Container } from '@mantine/core'
import { LoginForm } from '@/features/auth/components/LoginForm'

export function LoginPage() {
  return (
    <Center mih="100vh" bg="gray.0">
      <Container size="xs" w="100%">
        <LoginForm />
      </Container>
    </Center>
  )
}
