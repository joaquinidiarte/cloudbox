import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Layout from './components/Layout'
import { Toaster } from '@/components/ui/toaster'

function PrivateRoute({ children }) {
  const { token } = useAuthStore()
  return token ? children : <Navigate to="/login" />
}

function PublicRoute({ children }) {
  const { token } = useAuthStore()
  return !token ? children : <Navigate to="/dashboard" />
}

function App() {
  return (
    <Router>
      <Toaster />
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" />} />
        
        <Route path="/login" element={
          <PublicRoute>
            <Login />
          </PublicRoute>
        } />
        
        <Route path="/register" element={
          <PublicRoute>
            <Register />
          </PublicRoute>
        } />
        
        <Route path="/dashboard" element={
          <PrivateRoute>
            <Layout>
              <Dashboard />
            </Layout>
          </PrivateRoute>
        } />
      </Routes>
    </Router>
  )
}

export default App
