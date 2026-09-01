import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Alert,
  AlertIcon,
  Box,
  Button,
  FormControl,
  FormLabel,
  Heading,
  Input,
  Stack,
  Text,
} from '@chakra-ui/react'
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
    <Box bg="white" p={8} borderRadius="xl" shadow="sm" borderWidth="1px">
      <Heading size="lg" mb={1}>
        Sign in
      </Heading>
      <Text color="gray.500" fontSize="sm" mb={6}>
        Enter your credentials to continue.
      </Text>
      <Stack as="form" spacing={4} onSubmit={handleSubmit}>
        <FormControl isRequired>
          <FormLabel>Email</FormLabel>
          <Input type="email" autoComplete="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Password</FormLabel>
          <Input type="password" autoComplete="current-password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </FormControl>
        {error && (
          <Alert status="error">
            <AlertIcon />
            {error.message}
          </Alert>
        )}
        <Button type="submit" colorScheme="blue" isLoading={isPending} width="full">
          Sign in
        </Button>
      </Stack>
    </Box>
  )
}
