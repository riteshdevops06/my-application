import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'

const mockSystemMetrics = {
  totalNGOs: 156,
  verifiedNGOs: 134,
  pendingVerifications: 22,
  totalCampaigns: 423,
  activeCampaigns: 287,
  totalDonations: 2847563,
  totalDonors: 15432,
  platformGrowth: 15.2, // percentage
  systemHealth: 98.5, // percentage
  averageResponseTime: 0.24, // seconds
}

const mockRecentActivities = [
  { id: 1, type: 'ngo_verification', message: 'Vision Care Foundation submitted verification documents', timestamp: '2024-03-15T10:30:00Z', status: 'pending' },
  { id: 2, type: 'campaign_approval', message: 'Cancer Treatment Assistance campaign approved', timestamp: '2024-03-15T09:45:00Z', status: 'approved' },
  { id: 3, type: 'expense_review', message: 'Medical Equipment expense requires review', timestamp: '2024-03-15T09:20:00Z', status: 'attention' },
  { id: 4, type: 'milestone_completed', message: 'Water Project Phase 1 milestone completed', timestamp: '2024-03-15T08:15:00Z', status: 'completed' },
  { id: 5, type: 'user_registration', message: 'Hope Cancer Support registered as new NGO', timestamp: '2024-03-14T16:30:00Z', status: 'info' },
]

const mockPendingActions = [
  { id: 1, type: 'ngo_verification', title: 'NGO Verification Pending', count: 22, urgency: 'high', link: '/admin/ngos' },
  { id: 2, type: 'expense_approval', title: 'Expense Approvals Needed', count: 15, urgency: 'medium', link: '/admin/campaigns' },
  { id: 3, type: 'milestone_review', title: 'Milestone Reviews', count: 8, urgency: 'low', link: '/admin/milestones' },
  { id: 4, type: 'campaign_flags', title: 'Flagged Campaigns', count: 3, urgency: 'high', link: '/admin/campaigns' },
]

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatNumber = (num) => new Intl.NumberFormat('en-US').format(num)
const formatRelativeTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diffInHours = Math.floor((now - time) / (1000 * 60 * 60))
  
  if (diffInHours < 1) return 'Just now'
  if (diffInHours < 24) return `${diffInHours}h ago`
  const diffInDays = Math.floor(diffInHours / 24)
  return `${diffInDays}d ago`
}

const urgencyColors = {
  high: 'bg-red-100 text-red-800 border-red-200',
  medium: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  low: 'bg-green-100 text-green-800 border-green-200',
}

const activityStatusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  attention: 'bg-red-100 text-red-800',
  completed: 'bg-blue-100 text-blue-800',
  info: 'bg-gray-100 text-gray-800',
}

const activityIcons = {
  ngo_verification: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg>,
  campaign_approval: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>,
  expense_review: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
  milestone_completed: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
  user_registration: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" /></svg>,
}

export default function PlatformAdminDashboard() {
  const [metrics] = useState(mockSystemMetrics)
  const [recentActivities] = useState(mockRecentActivities)
  const [pendingActions] = useState(mockPendingActions)

  return (
    <div className="space-y-8">
      {/* Welcome Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Platform Administration</h1>
          <p className="mt-1 text-gray-600">Monitor and manage the NGO donation platform</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1">
            <div className={`w-2 h-2 rounded-full ${metrics.systemHealth > 95 ? 'bg-green-500' : metrics.systemHealth > 80 ? 'bg-yellow-500' : 'bg-red-500'}`}></div>
            <span className="text-sm font-medium text-gray-600">
              System Health: {metrics.systemHealth}%
            </span>
          </div>
        </div>
      </div>

      {/* Key Metrics Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total NGOs</p>
              <p className="text-3xl font-bold text-gray-900">{formatNumber(metrics.totalNGOs)}</p>
              <p className="text-sm text-green-600 mt-1">
                {formatNumber(metrics.verifiedNGOs)} verified
              </p>
            </div>
            <div className="p-3 bg-blue-100 rounded-lg">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m0 0H7m-2 0h2M9 7h6m-6 4h6m-6 4h6" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Campaigns</p>
              <p className="text-3xl font-bold text-gray-900">{formatNumber(metrics.activeCampaigns)}</p>
              <p className="text-sm text-gray-500 mt-1">
                of {formatNumber(metrics.totalCampaigns)} total
              </p>
            </div>
            <div className="p-3 bg-green-100 rounded-lg">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Donations</p>
              <p className="text-3xl font-bold text-gray-900">{formatCurrency(metrics.totalDonations)}</p>
              <p className="text-sm text-green-600 mt-1">
                +{metrics.platformGrowth}% this month
              </p>
            </div>
            <div className="p-3 bg-purple-100 rounded-lg">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Donors</p>
              <p className="text-3xl font-bold text-gray-900">{formatNumber(metrics.totalDonors)}</p>
              <p className="text-sm text-gray-500 mt-1">
                Avg: {metrics.averageResponseTime}s response
              </p>
            </div>
            <div className="p-3 bg-orange-100 rounded-lg">
              <svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Pending Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-6">
        {pendingActions.map((action) => (
          <Link
            key={action.id}
            to={action.link}
            className={`block p-4 rounded-xl border-2 hover:shadow-lg transition-all ${urgencyColors[action.urgency]}`}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="font-semibold">{action.title}</p>
                <p className="text-2xl font-bold mt-1">{action.count}</p>
              </div>
              <div className="text-right">
                <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                  action.urgency === 'high' ? 'bg-red-200 text-red-900' :
                  action.urgency === 'medium' ? 'bg-yellow-200 text-yellow-900' :
                  'bg-green-200 text-green-900'
                }`}>
                  {action.urgency.toUpperCase()}
                </span>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Recent Activities */}
        <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">Recent Activities</h2>
              <Link to="/admin/reports" className="text-sm text-primary-600 hover:text-primary-500 font-medium">
                View all →
              </Link>
            </div>
          </div>
          <div className="divide-y divide-gray-200">
            {recentActivities.map((activity) => (
              <div key={activity.id} className="p-6 hover:bg-gray-50 transition-colors">
                <div className="flex items-start gap-3">
                  <div className={`p-2 rounded-lg flex-shrink-0 ${
                    activity.status === 'pending' ? 'bg-yellow-100 text-yellow-600' :
                    activity.status === 'approved' ? 'bg-green-100 text-green-600' :
                    activity.status === 'attention' ? 'bg-red-100 text-red-600' :
                    activity.status === 'completed' ? 'bg-blue-100 text-blue-600' :
                    'bg-gray-100 text-gray-600'
                  }`}>
                    {activityIcons[activity.type]}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900">{activity.message}</p>
                    <div className="mt-1 flex items-center gap-2">
                      <span className={`px-2 py-0.5 text-xs font-medium rounded-full ${activityStatusColors[activity.status]}`}>
                        {activity.status}
                      </span>
                      <span className="text-xs text-gray-500">{formatRelativeTime(activity.timestamp)}</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-6">Quick Actions</h2>
          <div className="space-y-4">
            <Link
              to="/admin/ngos"
              className="flex items-center p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors group"
            >
              <div className="p-2 bg-blue-100 rounded-lg group-hover:bg-blue-200 transition-colors">
                <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <div className="ml-4 flex-1">
                <p className="font-medium text-gray-900">Review NGO Verifications</p>
                <p className="text-sm text-gray-500">{metrics.pendingVerifications} pending reviews</p>
              </div>
              <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>

            <Link
              to="/admin/campaigns"
              className="flex items-center p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors group"
            >
              <div className="p-2 bg-green-100 rounded-lg group-hover:bg-green-200 transition-colors">
                <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div className="ml-4 flex-1">
                <p className="font-medium text-gray-900">Monitor Campaigns</p>
                <p className="text-sm text-gray-500">{metrics.activeCampaigns} active campaigns</p>
              </div>
              <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>

            <Link
              to="/admin/milestones"
              className="flex items-center p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors group"
            >
              <div className="p-2 bg-purple-100 rounded-lg group-hover:bg-purple-200 transition-colors">
                <svg className="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div className="ml-4 flex-1">
                <p className="font-medium text-gray-900">Review Milestones</p>
                <p className="text-sm text-gray-500">8 milestones need review</p>
              </div>
              <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>

            <Link
              to="/admin/reports"
              className="flex items-center p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors group"
            >
              <div className="p-2 bg-orange-100 rounded-lg group-hover:bg-orange-200 transition-colors">
                <svg className="w-5 h-5 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <div className="ml-4 flex-1">
                <p className="font-medium text-gray-900">Generate Reports</p>
                <p className="text-sm text-gray-500">Analytics and insights</p>
              </div>
              <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>

            <Link
              to="/admin/settings"
              className="flex items-center p-4 rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors group"
            >
              <div className="p-2 bg-gray-100 rounded-lg group-hover:bg-gray-200 transition-colors">
                <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <div className="ml-4 flex-1">
                <p className="font-medium text-gray-900">Platform Settings</p>
                <p className="text-sm text-gray-500">Configure system parameters</p>
              </div>
              <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Link>
          </div>
        </div>
      </div>

      {/* Platform Performance Chart Placeholder */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Platform Performance (Last 30 Days)</h2>
        <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
          <div className="text-center">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <p className="mt-2 text-sm text-gray-500">Performance analytics chart coming soon</p>
            <p className="text-xs text-gray-400">Chart.js integration in development</p>
          </div>
        </div>
      </div>
    </div>
  )
}