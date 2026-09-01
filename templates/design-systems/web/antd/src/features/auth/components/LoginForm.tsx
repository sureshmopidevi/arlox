import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Alert, Button, Card, Form, Input, Typography } from 'antd'
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
    <Card>
      <Typography.Title level={3}>Sign in</Typography.Title>
      <Typography.Paragraph type="secondary">Enter your credentials to continue.</Typography.Paragraph>
      <Form layout="vertical" onSubmitCapture={handleSubmit}>
        <Form.Item label="Email" required>
          <Input
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
          />
        </Form.Item>
        <Form.Item label="Password" required>
          <Input.Password
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Item>
        {error && <Alert type="error" message={error.message} className="mb-4" showIcon />}
        <Button type="primary" htmlType="submit" block loading={isPending}>
          Sign in
        </Button>
      </Form>
    </Card>
  )
}
