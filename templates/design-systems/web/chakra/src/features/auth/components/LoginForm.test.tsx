import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ChakraProvider } from '@chakra-ui/react'
import { LoginForm } from './LoginForm'

const client = new QueryClient()

function wrap(ui: React.ReactNode) {
  return (
    <QueryClientProvider client={client}>
      <ChakraProvider>
        <MemoryRouter>{ui}</MemoryRouter>
      </ChakraProvider>
    </QueryClientProvider>
  )
}

test('renders sign in heading', () => {
  render(wrap(<LoginForm />))
  expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
})
