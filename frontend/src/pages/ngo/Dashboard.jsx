import { Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

const statCards = [
  { title: 'Total Raised', value: '$45,000', icon: 'money', color: 'green', trend: '+12%' },
  { title: 'Active Campaigns', value: '3', icon: 'campaigns', color: 'blue', trend: '+1' },
  { title: 'Total Donors', value: '234', icon: 'users', color: 'purple', trend: '+23' },
  { title: 'Pending Verification', value: '1', icon: 'pending', color: 'orange', trend: 'Review needed' },
]

const statIcons = {
  money: <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
  campaigns: <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>,
  users: <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>,
  pending: <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>,
}

const colorClasses = {
  green: 'bg-green-50 text-green-600 border-green-200',
  blue: 'bg-blue-50 text-blue-600 border-blue-200',
  purple: 'bg-purple-50 text-purple-600 border-purple-200',
  orange: 'bg-orange-50 text-orange-600 border-orange-200',
}

const recentDonations = [
  { id: 1, donor: 'John Smith', amount: 100, campaign: 'Blindness Treatment', date: '2024-03-15' },
  { id: 2, donor: 'Sarah Johnson', amount: 50, campaign: 'Animal Vaccination', date: '2024-03-14' },
  { id: 3, donor: 'Michael Brown', amount: 250, campaign: 'Cancer Treatment', date: '2024-03-13' },
  { id: 4, donor: 'Emily Davis', amount: 75, campaign: 'Blindness Treatment', date: '2024-03-12' },
]

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)

export default function NGODashboard() {
  const { user } = useAuth()

  return (
    <div className="space-y-8">
      {/* Welcome Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Welcome, {user?.name || 'NGO Admin'}!</h1>
          <p className="mt-1 text-gray-600">Manage your campaigns and track donations</p>
        </div>
        <Link
          to="/ngo/campaigns/create"
          className="inline-flex items-center px-4 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 transition-colors"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Create Campaign
        </Link>
      </div>

      {/* Verification Status Alert */}
      <div className="bg-yellow-50 border border-yellow-200 rounded-xl p-4">
        <div className="flex items-start gap-3">
          <svg className="w-5 h-5 text-yellow-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          <div className="flex-1">
            <h3 className="text-sm font-medium text-yellow-800">NGO Verification Pending</h3>
            <p className="mt-1 text-sm text-yellow-700">Your NGO verification is under review. You can create campaigns but cannot publish them until verified.</p>
            <Link to="/ngo/verification" className="mt-2 inline-block text-sm font-medium text-yellow-800 underline hover:text-yellow-900">
              View verification status →
            </Link>
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat) => (
          <div key={stat.title} className={`p-6 rounded-xl border transition-all ${colorClasses[stat.color]}`}>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium">{stat.title}</p>
                <p className="mt-2 text-3xl font-bold">{stat.value}</p>
                <p className="mt-1 text-sm font-medium">{stat.trend}</p>
              </div>
              <div className="p-3 rounded-lg bg-white/50">
                {statIcons[stat.icon]}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Link to="/ngo/campaigns" className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow">
          <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Manage Campaigns</h3>
          <p className="mt-2 text-gray-600">View and edit your campaigns</p>
        </Link>

        <Link to="/ngo/campaigns/create" className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow">
          <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Create Campaign</h3>
          <p className="mt-2 text-gray-600">Start a new fundraising campaign</p>
        </Link>

        <Link to="/ngo/verification" className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow">
          <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">NGO Verification</h3>
          <p className="mt-2 text-gray-600">Submit verification documents</p>
        </Link>

        <Link to="/ngo/profile" className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow">
          <div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Profile Settings</h3>
          <p className="mt-2 text-gray-600">Update your NGO profile</p>
        </Link>
      </div>

      {/* Recent Donations */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="p-6 border-b border-gray-200 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">Recent Donations</h2>
          <Link to="#" className="text-sm text-primary-600 hover:text-primary-500 font-medium">View all</Link>
        </div>
        <div className="divide-y divide-gray-200">
          {recentDonations.map((donation) => (
            <div key={donation.id} className="p-6 hover:bg-gray-50 flex items-center justify-between">
              <div>
                <p className="font-medium text-gray-900">{donation.donor}</p>
                <p className="text-sm text-gray-500">{donation.campaign} • {new Date(donation.date).toLocaleDateString()}</p>
              </div>
              <p className="text-lg font-bold text-green-600">{formatCurrency(donation.amount)}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Campaign Performance Chart Placeholder */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Fundraising Progress (Last 30 Days)</h2>
        <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
          <p className="text-gray-500">Chart.js integration coming in Phase 12</p>
        </div>
      </div>
    </div>
  )
}