import { NavLink, useNavigate } from 'react-router-dom'
import { Box, Button, Flex, Text, VStack } from '@chakra-ui/react'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', to: '/' }]

export function AdminSidebar() {
  const token = useAuthStore((s) => s.token)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  return (
    <Flex as="aside" direction="column" w="56" flexShrink={0} bg="white" borderRightWidth="1px">
      <Box px={4} py={4} borderBottomWidth="1px">
        <Text fontSize="sm" fontWeight="semibold">{envConfig.appName}</Text>
      </Box>
      <VStack align="stretch" flex={1} p={3} spacing={1}>
        {navItems.map(({ label, to }) => (
          <Button key={to} as={NavLink} to={to} variant="ghost" justifyContent="flex-start" size="sm">
            {label}
          </Button>
        ))}
      </VStack>
      <Box p={3} borderTopWidth="1px">
        {token ? (
          <Button variant="ghost" size="sm" width="full" onClick={() => { clearAuth(); navigate('/', { replace: true }) }}>
            Sign out
          </Button>
        ) : (
          <Button as={NavLink} to="/login" variant="link" size="sm" width="full">
            Sign in
          </Button>
        )}
      </Box>
    </Flex>
  )
}
