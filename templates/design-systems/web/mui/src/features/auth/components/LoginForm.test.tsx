import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import { LoginForm } from './LoginForm'

const client = new QueryClient()
const theme = createTheme()

function wrap(ui: React.ReactNode) {
  return (
    <QueryClientProvider client={client}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <MemoryRouter>{ui}</MemoryRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

test('renders sign in heading', () => {
  render(wrap(<LoginForm />))
  expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
})
