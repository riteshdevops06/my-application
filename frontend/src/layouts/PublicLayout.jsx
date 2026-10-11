import { Outlet, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function PublicLayout() {
  const { isAuthenticated } = useAuth()

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation */}
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo */}
            <Link to="/campaigns" className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.154a4 4 0 110 5.692M15.211 4.154a4 4 0 110 5.692m6.656 10.879a4 4 0 11-5.692 0m5.692 0a4 4 0 110-5.692m-5.692 0a4 4 0 110 5.692m0-11.384a4 4 0 115.692 0" />
                </svg>
              </div>
              <span className="text-xl font-bold text-gray-900">NGO Tracker</span>
            </Link>

            {/* Navigation Links */}
            <div className="hidden md:flex items-center space-x-8">
              <Link to="/campaigns" className="text-gray-600 hover:text-primary-600 font-medium transition-colors">
                Browse Campaigns
              </Link>
              {isAuthenticated ? (
                <>
                  <Link to="/donor/dashboard" className="text-gray-600 hover:text-primary-600 font-medium transition-colors">
                    Dashboard
                  </Link>
                </>
              ) : (
                <div className="flex items-center space-x-4">
                  <Link to="/login" className="text-gray-600 hover:text-primary-600 font-medium transition-colors">
                    Login
                  </Link>
                  <Link
                    to="/register"
                    className="bg-primary-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-primary-700 transition-colors"
                  >
                    Register
                  </Link>
                </div>
              )}
            </div>

            {/* Mobile menu button */}
            <div className="md:hidden">
              <button className="text-gray-600 hover:text-gray-900">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="min-h-[calc(100vh-64px)]">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-white text-lg font-semibold mb-4">NGO Tracker</h3>
              <p className="text-sm">
                Transparent donation tracking for a better world.
              </p>
            </div>
            <div>
              <h4 className="text-white text-sm font-semibold mb-4">Quick Links</h4>
              <ul className="space-y-2 text-sm">
                <li><Link to="/campaigns" className="hover:text-white transition-colors">Browse Campaigns</Link></li>
                <li><Link to="/login" className="hover:text-white transition-colors">Login</Link></li>
                <li><Link to="/register" className="hover:text-white transition-colors">Register</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="text-white text-sm font-semibold mb-4">For NGOs</h4>
              <ul className="space-y-2 text-sm">
                <li><Link to="/register?role=ngo" className="hover:text-white transition-colors">Register NGO</Link></li>
                <li><Link to="/ngo/verification" className="hover:text-white transition-colors">Verification</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="text-white text-sm font-semibold mb-4">Legal</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">Privacy Policy</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Terms of Service</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-sm">
            &copy; 2024 NGO Donation Transparency Tracker. All rights reserved.
          </div>
        </div>
      </footer>
    </div>
  )
}