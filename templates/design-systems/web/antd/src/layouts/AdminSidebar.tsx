import { NavLink, useNavigate } from 'react-router-dom'
import { Button, Menu } from 'antd'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', key: '/', to: '/' }]

export function AdminSidebar() {
  const token = useAuthStore((s) => s.token)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  return (
    <aside className="flex w-56 shrink-0 flex-col border-r bg-white">
      <div className="border-b px-4 py-4">
        <span className="text-sm font-semibold">{envConfig.appName}</span>
      </div>
      <Menu
        mode="inline"
        className="flex-1 border-0"
        items={navItems.map(({ label, key, to }) => ({
          key,
          label: <NavLink to={to}>{label}</NavLink>,
        }))}
      />
      <div className="border-t p-3">
        {token ? (
          <Button type="text" block onClick={() => { clearAuth(); navigate('/', { replace: true }) }}>
            Sign out
          </Button>
        ) : (
          <NavLink to="/login">
            <Button type="link" block>Sign in</Button>
          </NavLink>
        )}
      </div>
    </aside>
  )
}
