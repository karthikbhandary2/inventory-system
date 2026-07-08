import { Routes, Route, NavLink } from 'react-router-dom'
import ProductsPage from './pages/ProductsPage'
import ReportPage from './pages/ReportPage'
import AuditPage from './pages/AuditPage'
import LoginGate from './components/LoginGate'

export default function App() {
  return (
    <LoginGate>
      <div className="min-h-screen bg-gray-50">
        <nav className="bg-white border-b px-6 py-3 flex items-center gap-6">
          <span className="font-semibold text-sm">Inventory</span>
          <NavTab to="/">Products</NavTab>
          <NavTab to="/report">Report</NavTab>
          <NavTab to="/audit">Audit log</NavTab>
        </nav>
        <Routes>
          <Route path="/" element={<ProductsPage />} />
          <Route path="/report" element={<ReportPage />} />
          <Route path="/audit" element={<AuditPage />} />
        </Routes>
      </div>
    </LoginGate>
  )
}

function NavTab({ to, children }) {
  return (
    <NavLink
      to={to}
      end={to === '/'}
      className={({ isActive }) =>
        `text-sm pb-1 border-b-2 ${
          isActive ? 'border-blue-600 text-blue-600' : 'border-transparent text-gray-500'
        }`
      }
    >
      {children}
    </NavLink>
  )
}