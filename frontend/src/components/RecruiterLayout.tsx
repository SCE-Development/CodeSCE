import { NavLink, Outlet } from 'react-router-dom'

const navItems = [
  { to: '/recruiter', label: 'Dashboard', end: true },
  { to: '/recruiter/candidates', label: 'Candidates' },
  { to: '/recruiter/assessments', label: 'Assessments' },
]

function RecruiterLayout() {
  return (
    <div className="min-h-screen flex bg-gray-50">
      <aside className="w-60 bg-white border-r border-gray-200 flex flex-col">
        <div className="px-6 py-5 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">CodeSCE</h2>
          <p className="text-xs text-gray-500">Recruiter</p>
        </div>
        <nav className="flex-1 px-3 py-4 space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                `block px-3 py-2 rounded-md text-sm font-medium ${
                  isActive
                    ? 'bg-gray-900 text-white'
                    : 'text-gray-700 hover:bg-gray-100'
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}

export default RecruiterLayout
