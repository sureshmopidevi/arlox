import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MantineProvider } from '@mantine/core'
import { LoginForm } from './LoginForm'

const client = new QueryClient()

function wrap(ui: React.ReactNode) {
  return (
    <QueryClientProvider client={client}>
      <MantineProvider>
        <MemoryRouter>{ui}</MemoryRouter>
      </MantineProvider>
    </QueryClientProvider>
  )
}

test('renders sign in heading', () => {
  render(wrap(<LoginForm />))
  expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
})
