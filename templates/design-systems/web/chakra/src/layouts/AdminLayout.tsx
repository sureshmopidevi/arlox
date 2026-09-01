import { Flex } from '@chakra-ui/react'
import { Outlet } from 'react-router-dom'
import { AdminSidebar } from './AdminSidebar'

export function AdminLayout() {
  return (
    <Flex h="100vh" overflow="hidden" bg="gray.100">
      <AdminSidebar />
      <Flex as="main" flex={1} overflow="auto" p={6}>
        <Outlet />
      </Flex>
    </Flex>
  )
}
