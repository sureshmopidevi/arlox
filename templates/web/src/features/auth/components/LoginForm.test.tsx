import { expect, test } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { Providers } from '@/app/providers'
import { LoginForm } from './LoginForm'

test('renders sign in heading', () => {
  render(
    <Providers>
      <MemoryRouter>
        <LoginForm />
      </MemoryRouter>
    </Providers>,
  )
  expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
})
