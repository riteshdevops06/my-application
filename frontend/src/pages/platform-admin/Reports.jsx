import { useState } from 'react'

const mockReportData = {
  totalDonations: 2847563,
  totalDonors: 15432,
  totalNGOs: 156,
  totalCampaigns: 423,
  averageDonation: 184.52,
  donationGrowth: 15.2,
  topCategories: [
    { name: 'Healthcare', amount: 1250000, campaigns: 89, percentage: 43.9 },
    { name: 'Education', amount: 856000, campaigns: 67, percentage: 30.1 },
    { name: 'Emergency Relief', amount: 445000, campaigns: 34, percentage: 15.6 },
    { name: 'Community Development', amount: 296563, campaigns: 45, percentage: 10.4 },
  ],
  topNGOs: [
    { name: 'Vision Care Foundation', raised: 125000, campaigns: 3, donors: 2150 },
    { name: 'Hope Cancer Support', raised: 89000, campaigns: 2, donors: 1834 },
    { name: 'Global Aid Network', raised: 67000, campaigns: 1, donors: 892 },
    { name: 'Education for All', raised: 45000, campaigns: 2, donors: 1245 },
  ],
  monthlyTrends: [
    { month: 'Jan', donations: 185000, donors: 1240 },
    { month: 'Feb', donations: 220000, donors: 1450 },
    { month: 'Mar', donations: 195000, donors: 1380 },
    { month: 'Apr', donations: 275000, donors: 1890 },
    { month: 'May', donations: 310000, donors: 2100 },
    { month: 'Jun', donations: 285000, donors: 1950 },
  ],
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatNumber = (num) => new Intl.NumberFormat('en-US').format(num)

export default function PlatformAdminReports() {
  const [selectedTimeRange, setSelectedTimeRange] = useState('6months')
  const [reportType, setReportType] = useState('overview')
  const [data] = useState(mockReportData)

  const exportReport = (format) => {
    // TODO: Implement export functionality
    console.log(`Exporting ${reportType} report in ${format} format`)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Analytics & Reports</h1>
          <p className="mt-1 text-gray-600">Platform performance insights and detailed reports</p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => exportReport('csv')}
            className="px-3 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Export CSV
          </button>
          <button
            onClick={() => exportReport('pdf')}
            className="px-3 py-2 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors"
          >
            Export PDF
          </button>
        </div>
      </div>

      {/* Controls */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <div>
            <label htmlFor="reportType" className="block text-sm font-medium text-gray-700 mb-1">Report Type</label>
            <select
              id="reportType"
              value={reportType}
              onChange={(e) => setReportType(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="overview">Platform Overview</option>
              <option value="donations">Donation Analytics</option>
              <option value="ngos">NGO Performance</option>
              <option value="campaigns">Campaign Analytics</option>
            </select>
          </div>
          <div>
            <label htmlFor="timeRange" className="block text-sm font-medium text-gray-700 mb-1">Time Range</label>
            <select
              id="timeRange"
              value={selectedTimeRange}
              onChange={(e) => setSelectedTimeRange(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="1month">Last Month</option>
              <option value="3months">Last 3 Months</option>
              <option value="6months">Last 6 Months</option>
              <option value="1year">Last Year</option>
              <option value="all">All Time</option>
            </select>
          </div>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Donations</p>
              <p className="text-2xl font-bold text-gray-900">{formatCurrency(data.totalDonations)}</p>
              <p className="text-sm text-green-600 mt-1">+{data.donationGrowth}% vs last period</p>
            </div>
            <div className="p-3 bg-green-100 rounded-lg">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Donors</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(data.totalDonors)}</p>
              <p className="text-sm text-gray-500 mt-1">Avg: {formatCurrency(data.averageDonation)}</p>
            </div>
            <div className="p-3 bg-blue-100 rounded-lg">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active NGOs</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(data.totalNGOs)}</p>
              <p className="text-sm text-gray-500 mt-1">{Math.round(data.totalCampaigns / data.totalNGOs)} campaigns avg</p>
            </div>
            <div className="p-3 bg-purple-100 rounded-lg">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m0 0H7m-2 0h2M9 7h6m-6 4h6m-6 4h6" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Campaigns</p>
              <p className="text-2xl font-bold text-gray-900">{formatNumber(data.totalCampaigns)}</p>
              <p className="text-sm text-gray-500 mt-1">{Math.round(data.totalCampaigns * 0.68)} active</p>
            </div>
            <div className="p-3 bg-orange-100 rounded-lg">
              <svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Top Categories */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Top Categories by Donations</h2>
          <div className="space-y-4">
            {data.topCategories.map((category, index) => (
              <div key={category.name} className="flex items-center gap-4">
                <div className="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center text-sm font-medium text-primary-600">
                  {index + 1}
                </div>
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <p className="font-medium text-gray-900">{category.name}</p>
                    <div className="text-right">
                      <p className="font-semibold text-gray-900">{formatCurrency(category.amount)}</p>
                      <p className="text-xs text-gray-500">{category.percentage}%</p>
                    </div>
                  </div>
                  <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary-600 rounded-full transition-all"
                      style={{ width: `${category.percentage}%` }}
                    />
                  </div>
                  <p className="text-xs text-gray-500 mt-1">{category.campaigns} campaigns</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Top NGOs */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Top Performing NGOs</h2>
          <div className="space-y-4">
            {data.topNGOs.map((ngo, index) => (
              <div key={ngo.name} className="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 transition-colors">
                <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center text-sm font-medium text-green-600">
                  {index + 1}
                </div>
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{ngo.name}</p>
                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <span>{formatCurrency(ngo.raised)} raised</span>
                    <span>{ngo.campaigns} campaigns</span>
                    <span>{ngo.donors.toLocaleString()} donors</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Monthly Trends Chart */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Monthly Donation Trends</h2>
        <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
          <div className="text-center">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <p className="mt-2 text-sm text-gray-500">Interactive charts coming soon</p>
            <p className="text-xs text-gray-400">Chart.js integration in development</p>
          </div>
        </div>
        
        {/* Simple Data Display */}
        <div className="mt-6 grid grid-cols-2 sm:grid-cols-6 gap-4">
          {data.monthlyTrends.map((month) => (
            <div key={month.month} className="text-center">
              <p className="text-sm font-medium text-gray-900">{month.month}</p>
              <p className="text-lg font-bold text-primary-600">{formatCurrency(month.donations)}</p>
              <p className="text-xs text-gray-500">{month.donors.toLocaleString()} donors</p>
            </div>
          ))}
        </div>
      </div>

      {/* System Health */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">System Health</h2>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="w-16 h-16 mx-auto mb-3 rounded-full bg-green-100 flex items-center justify-center">
              <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p className="text-2xl font-bold text-green-600">98.5%</p>
            <p className="text-sm text-gray-600">System Uptime</p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 mx-auto mb-3 rounded-full bg-blue-100 flex items-center justify-center">
              <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <p className="text-2xl font-bold text-blue-600">0.24s</p>
            <p className="text-sm text-gray-600">Avg Response Time</p>
          </div>
          <div className="text-center">
            <div className="w-16 h-16 mx-auto mb-3 rounded-full bg-purple-100 flex items-center justify-center">
              <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
              </svg>
            </div>
            <p className="text-2xl font-bold text-purple-600">99.9%</p>
            <p className="text-sm text-gray-600">Data Integrity</p>
          </div>
        </div>
      </div>
    </div>
  )
}