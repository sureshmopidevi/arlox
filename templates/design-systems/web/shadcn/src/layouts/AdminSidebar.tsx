import { NavLink, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', to: '/' }]

export function AdminSidebar() {
  const token = useAuthStore((s) => s.token)
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  function handleSignOut() {
    clearAuth()
    navigate('/', { replace: true })
  }

  return (
    <aside className="flex w-56 shrink-0 flex-col border-r bg-card">
      <div className="border-b px-4 py-4">
        <span className="text-sm font-semibold tracking-tight">{envConfig.appName}</span>
      </div>
      <nav className="flex flex-1 flex-col gap-0.5 p-3">
        {navItems.map(({ label, to }) => (
          <NavLink
            key={to}
            to={to}
            end
            className={({ isActive }) =>
              cn(
                'rounded-md px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
              )
            }
          >
            {label}
          </NavLink>
        ))}
      </nav>
      <div className="border-t p-3">
        {token ? (
          <Button variant="ghost" className="w-full justify-start" onClick={handleSignOut}>
            Sign out
          </Button>
        ) : (
          <Button variant="link" className="w-full justify-start px-3" asChild>
            <NavLink to="/login">Sign in</NavLink>
          </Button>
        )}
      </div>
    </aside>
  )
}
