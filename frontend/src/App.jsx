import { Routes, Route, Navigate, Outlet } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import { NotificationProvider } from './context/NotificationContext'

// Layout components
import PublicLayout from './layouts/PublicLayout'
import DonorLayout from './layouts/DonorLayout'
import NGOLayout from './layouts/NGOLayout'
import PlatformAdminLayout from './layouts/PlatformAdminLayout'

// Auth pages
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import ForgotPassword from './pages/auth/ForgotPassword'

// Donor pages
import DonorDashboard from './pages/donor/Dashboard'
import DonorCampaigns from './pages/donor/Campaigns'
import DonorCampaignDetail from './pages/donor/CampaignDetail'
import DonorDonations from './pages/donor/Donations'
import DonorProfile from './pages/donor/Profile'

// NGO pages
import NGODashboard from './pages/ngo/Dashboard'
import NGOCampaigns from './pages/ngo/Campaigns'
import NGOCampaignCreate from './pages/ngo/CampaignCreate'
import NGOCampaignEdit from './pages/ngo/CampaignEdit'
import NGOMilestones from './pages/ngo/Milestones'
import NGOExpenses from './pages/ngo/Expenses'
import NGOProfile from './pages/ngo/Profile'
import NGOVerification from './pages/ngo/Verification'

// Platform Admin pages
import PlatformAdminDashboard from './pages/platform-admin/Dashboard'
import PlatformAdminNGOs from './pages/platform-admin/NGOs'
import PlatformAdminCampaigns from './pages/platform-admin/Campaigns'
import PlatformAdminMilestones from './pages/platform-admin/Milestones'
import PlatformAdminReports from './pages/platform-admin/Reports'
import PlatformAdminSettings from './pages/platform-admin/Settings'

// Protected route components
function ProtectedRoute({ children, allowedRoles }) {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (allowedRoles && !allowedRoles.includes(user.role)) {
    // Redirect to appropriate dashboard based on role
    switch (user.role) {
      case 'donor':
        return <Navigate to="/donor/dashboard" replace />
      case 'ngo_admin':
        return <Navigate to="/ngo/dashboard" replace />
      case 'platform_admin':
        return <Navigate to="/admin/dashboard" replace />
      default:
        return <Navigate to="/login" replace />
    }
  }

  return children
}

// Public route (redirects authenticated users to their dashboard)
function PublicRoute({ children }) {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
      </div>
    )
  }

  if (user) {
    switch (user.role) {
      case 'donor':
        return <Navigate to="/donor/dashboard" replace />
      case 'ngo_admin':
        return <Navigate to="/ngo/dashboard" replace />
      case 'platform_admin':
        return <Navigate to="/admin/dashboard" replace />
      default:
        return <Navigate to="/login" replace />
    }
  }

  return children
}

function AppRoutes() {
  return (
    <Routes>
      {/* Public routes */}
      <Route element={<PublicLayout />}>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <Login />
            </PublicRoute>
          }
        />
        <Route
          path="/register"
          element={
            <PublicRoute>
              <Register />
            </PublicRoute>
          }
        />
        <Route
          path="/forgot-password"
          element={
            <PublicRoute>
              <ForgotPassword />
            </PublicRoute>
          }
        />
        {/* Public campaign browsing for donors (no login required) */}
        <Route path="/campaigns" element={<DonorCampaigns />} />
        <Route path="/campaigns/:id" element={<DonorCampaignDetail />} />
      </Route>

      {/* Donor protected routes */}
      <Route
        path="/donor"
        element={
          <ProtectedRoute allowedRoles={['donor']}>
            <DonorLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<DonorDashboard />} />
        <Route path="donations" element={<DonorDonations />} />
        <Route path="profile" element={<DonorProfile />} />
      </Route>

      {/* NGO Admin protected routes */}
      <Route
        path="/ngo"
        element={
          <ProtectedRoute allowedRoles={['ngo_admin']}>
            <NGOLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<NGODashboard />} />
        <Route path="campaigns" element={<NGOCampaigns />} />
        <Route path="campaigns/create" element={<NGOCampaignCreate />} />
        <Route path="campaigns/:id/edit" element={<NGOCampaignEdit />} />
        <Route path="campaigns/:id/milestones" element={<NGOMilestones />} />
        <Route path="campaigns/:id/expenses" element={<NGOExpenses />} />
        <Route path="profile" element={<NGOProfile />} />
        <Route path="verification" element={<NGOVerification />} />
      </Route>

      {/* Platform Admin protected routes */}
      <Route
        path="/admin"
        element={
          <ProtectedRoute allowedRoles={['platform_admin']}>
            <PlatformAdminLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<PlatformAdminDashboard />} />
        <Route path="ngos" element={<PlatformAdminNGOs />} />
        <Route path="campaigns" element={<PlatformAdminCampaigns />} />
        <Route path="milestones" element={<PlatformAdminMilestones />} />
        <Route path="reports" element={<PlatformAdminReports />} />
        <Route path="settings" element={<PlatformAdminSettings />} />
      </Route>

      {/* Redirect root to login */}
      <Route path="/" element={<Navigate to="/login" replace />} />
      {/* Catch-all 404 */}
      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <AuthProvider>
      <NotificationProvider>
        <AppRoutes />
      </NotificationProvider>
    </AuthProvider>
  )
}

export default App