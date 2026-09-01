import { Flex } from '@chakra-ui/react'
import { LoginForm } from '@/features/auth/components/LoginForm'

export function LoginPage() {
  return (
    <Flex minH="100vh" align="center" justify="center" bg="gray.50" px={4}>
      <Flex w="full" maxW="md">
        <LoginForm />
      </Flex>
    </Flex>
  )
}
