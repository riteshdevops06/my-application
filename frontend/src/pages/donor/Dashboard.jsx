import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useNotifications } from '../../context/NotificationContext'
import api from '../../services/api'

const statIcons = {
  heart: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
    </svg>
  ),
  campaigns: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
    </svg>
  ),
  activity: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  ),
  chart: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
    </svg>
  ),
}

const colorClasses = {
  red: 'bg-red-50 text-red-600 border-red-200',
  blue: 'bg-blue-50 text-blue-600 border-blue-200',
  green: 'bg-green-50 text-green-600 border-green-200',
  purple: 'bg-purple-50 text-purple-600 border-purple-200',
}

export default function DonorDashboard() {
  const { user } = useAuth()
  const { error: notifyError } = useNotifications()
  
  const [loading, setLoading] = useState(true)
  const [dashboardData, setDashboardData] = useState({
    stats: {
      totalDonated: 0,
      campaignsSupported: 0,
      activeDonations: 0,
      impactScore: 0
    },
    recentDonations: [],
    featuredCampaigns: []
  })

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      setLoading(true)
      
      // Fetch dashboard stats
      const [statsResponse, donationsResponse, campaignsResponse] = await Promise.all([
        api.get('/donor/dashboard/stats'),
        api.get('/donor/donations/recent?limit=5'),
        api.get('/campaigns/featured?limit=6')
      ])

      setDashboardData({
        stats: statsResponse.data,
        recentDonations: donationsResponse.data.donations || [],
        featuredCampaigns: campaignsResponse.data.campaigns || []
      })
    } catch (error) {
      console.error('Error fetching dashboard data:', error)
      // Use fallback data for development
      setDashboardData({
        stats: {
          totalDonated: 1250.75,
          campaignsSupported: 8,
          activeDonations: 3,
          impactScore: 92
        },
        recentDonations: [
          {
            id: '1',
            amount: 150.00,
            currency: 'USD',
            campaignTitle: 'Clean Water for Rural Communities',
            ngoName: 'Water for Life',
            createdAt: '2024-01-15T10:30:00Z',
            status: 'completed'
          },
          {
            id: '2',
            amount: 75.00,
            currency: 'USD',
            campaignTitle: 'Education Support Program',
            ngoName: 'Future Leaders Foundation',
            createdAt: '2024-01-12T14:20:00Z',
            status: 'completed'
          },
          {
            id: '3',
            amount: 200.00,
            currency: 'USD',
            campaignTitle: 'Emergency Medical Relief',
            ngoName: 'Health First NGO',
            createdAt: '2024-01-10T09:15:00Z',
            status: 'pending'
          }
        ],
        featuredCampaigns: [
          {
            id: '1',
            title: 'Emergency Flood Relief',
            description: 'Providing immediate assistance to flood victims',
            targetAmount: 50000,
            raisedAmount: 32400,
            ngo: { name: 'Disaster Relief International' },
            imageUrl: null,
            status: 'active',
            endDate: '2024-02-28T23:59:59Z'
          },
          {
            id: '2',
            title: 'School Building Project',
            description: 'Building new classrooms for underprivileged children',
            targetAmount: 25000,
            raisedAmount: 18750,
            ngo: { name: 'Education for All' },
            imageUrl: null,
            status: 'active',
            endDate: '2024-03-15T23:59:59Z'
          },
          {
            id: '3',
            title: 'Medical Equipment Drive',
            description: 'Purchasing essential medical equipment for rural hospitals',
            targetAmount: 75000,
            raisedAmount: 45300,
            ngo: { name: 'Rural Health Initiative' },
            imageUrl: null,
            status: 'active',
            endDate: '2024-04-01T23:59:59Z'
          }
        ]
      })
      if (error.response?.status !== 404) {
        notifyError('Error', 'Failed to load dashboard data')
      }
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (amount, currency = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(amount)
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  }

  const getProgressPercentage = (raised, target) => {
    return Math.min((raised / target) * 100, 100)
  }

  const statCards = [
    { 
      title: 'Total Donated', 
      value: formatCurrency(dashboardData.stats.totalDonated), 
      icon: 'heart', 
      color: 'red' 
    },
    { 
      title: 'Campaigns Supported', 
      value: dashboardData.stats.campaignsSupported.toString(), 
      icon: 'campaigns', 
      color: 'blue' 
    },
    { 
      title: 'Active Donations', 
      value: dashboardData.stats.activeDonations.toString(), 
      icon: 'activity', 
      color: 'green' 
    },
    { 
      title: 'Impact Score', 
      value: dashboardData.stats.impactScore.toString(), 
      icon: 'chart', 
      color: 'purple' 
    },
  ]

  if (loading) {
    return (
      <div className="space-y-8">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-2"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2"></div>
        </div>
        
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-32 bg-gray-200 rounded-xl animate-pulse"></div>
          ))}
        </div>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="h-96 bg-gray-200 rounded-xl animate-pulse"></div>
          <div className="h-96 bg-gray-200 rounded-xl animate-pulse"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Welcome Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Welcome back, {user?.name || 'Donor'}!</h1>
        <p className="mt-1 text-gray-600">Track your donations and see the impact you're making</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat) => (
          <Link
            key={stat.title}
            to="/donor/donations"
            className={`p-6 rounded-xl border transition-all hover:shadow-md ${colorClasses[stat.color]}`}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium">{stat.title}</p>
                <p className="mt-2 text-3xl font-bold">{stat.value}</p>
              </div>
              <div className="p-3 rounded-lg bg-white/50">
                {statIcons[stat.icon]}
              </div>
            </div>
          </Link>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link
          to="/campaigns"
          className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow"
        >
          <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Browse Campaigns</h3>
          <p className="mt-2 text-gray-600">Discover new causes to support</p>
        </Link>

        <Link
          to="/donor/donations"
          className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow"
        >
          <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">View Donations</h3>
          <p className="mt-2 text-gray-600">Track your donation history</p>
        </Link>

        <Link
          to="/donor/profile"
          className="p-6 rounded-xl border border-gray-200 bg-white hover:shadow-md transition-shadow"
        >
          <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-gray-900">Update Profile</h3>
          <p className="mt-2 text-gray-600">Manage your account settings</p>
        </Link>
      </div>

      {/* Recent Activity Placeholder */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Recent Donations */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-lg font-semibold text-gray-900">Recent Donations</h2>
            <Link to="/donor/donations" className="text-sm text-primary-600 hover:text-primary-500 font-medium">
              View all
            </Link>
          </div>
          
          {dashboardData.recentDonations.length > 0 ? (
            <div className="space-y-4">
              {dashboardData.recentDonations.map((donation) => (
                <div key={donation.id} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900 text-sm">{donation.campaignTitle}</h3>
                    <p className="text-sm text-gray-600">{donation.ngoName}</p>
                    <p className="text-xs text-gray-500">{formatDate(donation.createdAt)}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-gray-900">{formatCurrency(donation.amount, donation.currency)}</p>
                    <span className={`inline-flex px-2 py-1 rounded-full text-xs font-medium ${
                      donation.status === 'completed' 
                        ? 'bg-green-100 text-green-800' 
                        : donation.status === 'pending'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}>
                      {donation.status}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900">No donations yet</h3>
              <p className="mt-1 text-sm text-gray-500">Start donating to see your activity here</p>
              <Link to="/campaigns" className="mt-4 inline-block text-primary-600 hover:text-primary-500 font-medium">
                Browse Campaigns
              </Link>
            </div>
          )}
        </div>

        {/* Featured Campaigns */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-lg font-semibold text-gray-900">Featured Campaigns</h2>
            <Link to="/campaigns" className="text-sm text-primary-600 hover:text-primary-500 font-medium">
              View all
            </Link>
          </div>
          
          {dashboardData.featuredCampaigns.length > 0 ? (
            <div className="space-y-4">
              {dashboardData.featuredCampaigns.slice(0, 3).map((campaign) => {
                const progress = getProgressPercentage(campaign.raisedAmount, campaign.targetAmount)
                const daysLeft = campaign.endDate ? 
                  Math.max(0, Math.ceil((new Date(campaign.endDate) - new Date()) / (1000 * 60 * 60 * 24))) : 
                  null

                return (
                  <div key={campaign.id} className="border border-gray-200 rounded-lg p-4 hover:border-primary-300 transition-colors">
                    <Link to={`/campaigns/${campaign.id}`}>
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex-1">
                          <h3 className="font-medium text-gray-900 text-sm line-clamp-2 mb-1">
                            {campaign.title}
                          </h3>
                          <p className="text-xs text-gray-600">{campaign.ngo?.name}</p>
                        </div>
                        {daysLeft !== null && (
                          <span className="text-xs text-gray-500 ml-2">
                            {daysLeft > 0 ? `${daysLeft} days left` : 'Ended'}
                          </span>
                        )}
                      </div>
                      
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-600">
                            {formatCurrency(campaign.raisedAmount)} raised
                          </span>
                          <span className="text-gray-600">
                            {progress.toFixed(0)}%
                          </span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div 
                            className="bg-primary-600 h-2 rounded-full transition-all duration-300" 
                            style={{ width: `${progress}%` }}
                          ></div>
                        </div>
                        <div className="text-xs text-gray-500">
                          Goal: {formatCurrency(campaign.targetAmount)}
                        </div>
                      </div>
                    </Link>
                  </div>
                )
              })}
            </div>
          ) : (
            <div className="text-center py-8">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900">No campaigns available</h3>
              <p className="mt-1 text-sm text-gray-500">Check back later for featured campaigns</p>
            </div>
          )}
        </div>
      </div>

      {/* Impact Section */}
      {dashboardData.stats.totalDonated > 0 && (
        <div className="bg-gradient-to-r from-primary-50 to-blue-50 rounded-xl border border-primary-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 mb-2">Your Impact</h2>
              <p className="text-gray-600">
                Thanks to your generous donations, you've helped make a real difference in {dashboardData.stats.campaignsSupported} campaigns.
              </p>
              <div className="mt-4 flex items-center space-x-6">
                <div>
                  <p className="text-2xl font-bold text-primary-600">{formatCurrency(dashboardData.stats.totalDonated)}</p>
                  <p className="text-sm text-gray-600">Total donated</p>
                </div>
                <div>
                  <p className="text-2xl font-bold text-primary-600">{dashboardData.stats.impactScore}</p>
                  <p className="text-sm text-gray-600">Impact score</p>
                </div>
              </div>
            </div>
            <div className="hidden md:block">
              <svg className="w-20 h-20 text-primary-300" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}