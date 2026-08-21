import { NavLink, useNavigate } from 'react-router-dom'
import { envConfig } from '@/config/envConfig'
import { useAuthStore } from '@/stores/authStore'

const navItems = [{ label: 'Home', to: '/' }]

export function AdminSidebar() {
  const clearAuth = useAuthStore((s) => s.clearAuth)
  const navigate = useNavigate()

  function handleSignOut() {
    clearAuth()
    navigate('/login', { replace: true })
  }

  return (
    <aside className="w-56 shrink-0 bg-white border-r border-gray-200 flex flex-col">
      <div className="px-4 py-4 border-b border-gray-200">
        <span className="text-sm font-semibold text-gray-800 tracking-tight">
          {envConfig.appName}
        </span>
      </div>

      <nav className="flex-1 p-3 space-y-0.5">
        {navItems.map(({ label, to }) => (
          <NavLink
            key={to}
            to={to}
            end
            className={({ isActive }) =>
              [
                'flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors',
                isActive
                  ? 'bg-indigo-50 text-indigo-700'
                  : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900',
              ].join(' ')
            }
          >
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="p-3 border-t border-gray-200">
        <button
          onClick={handleSignOut}
          className="w-full px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-700 rounded-md text-left transition-colors"
        >
          Sign out
        </button>
      </div>
    </aside>
  )
}
