import { useState } from 'react'
import { Link } from 'react-router-dom'

const mockDonations = [
  {
    id: 1,
    campaignId: 1,
    campaignTitle: 'Blindness Treatment & Eye Surgery',
    ngo: 'Vision Care Foundation',
    amount: 100,
    date: '2024-03-15',
    status: 'completed',
    receiptUrl: '/receipts/1.pdf',
    transactionId: 'TXN-2024-001234',
  },
  {
    id: 2,
    campaignId: 2,
    campaignTitle: 'Animal Vaccination Drive',
    ngo: 'Paws & Claws Rescue',
    amount: 50,
    date: '2024-02-28',
    status: 'completed',
    receiptUrl: '/receipts/2.pdf',
    transactionId: 'TXN-2024-001156',
  },
  {
    id: 3,
    campaignId: 3,
    campaignTitle: 'Cancer Treatment Assistance',
    ngo: 'Hope Cancer Support',
    amount: 250,
    date: '2024-01-20',
    status: 'completed',
    receiptUrl: '/receipts/3.pdf',
    transactionId: 'TXN-2024-000987',
  },
]

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

const statusColors = {
  completed: 'bg-green-100 text-green-800',
  pending: 'bg-yellow-100 text-yellow-800',
  failed: 'bg-red-100 text-red-800',
  refunded: 'bg-gray-100 text-gray-800',
}

export default function DonorDonations() {
  const [filter, setFilter] = useState('all')
  const [sortBy, setSortBy] = useState('date_desc')

  const filteredDonations = mockDonations.filter((d) => {
    if (filter === 'all') return true
    return d.status === filter
  }).sort((a, b) => {
    if (sortBy === 'date_desc') return new Date(b.date) - new Date(a.date)
    if (sortBy === 'date_asc') return new Date(a.date) - new Date(b.date)
    if (sortBy === 'amount_desc') return b.amount - a.amount
    if (sortBy === 'amount_asc') return a.amount - b.amount
    return 0
  })

  const totalDonated = mockDonations.reduce((sum, d) => sum + d.amount, 0)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">My Donations</h1>
          <p className="mt-1 text-gray-600">Track your donation history and download receipts</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="bg-primary-50 text-primary-700 px-4 py-2 rounded-lg">
            <p className="text-sm font-medium">Total Donated</p>
            <p className="text-2xl font-bold">{formatCurrency(totalDonated)}</p>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex gap-2 flex-wrap">
            {['all', 'completed', 'pending', 'failed'].map((status) => (
              <button
                key={status}
                onClick={() => setFilter(status)}
                className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${
                  filter === status
                    ? 'bg-primary-100 text-primary-700'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {status.charAt(0).toUpperCase() + status.slice(1)}
              </button>
            ))}
          </div>
          <div className="flex items-center gap-2 ml-auto">
            <label htmlFor="sort" className="text-sm text-gray-600">Sort by:</label>
            <select
              id="sort"
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="date_desc">Newest First</option>
              <option value="date_asc">Oldest First</option>
              <option value="amount_desc">Highest Amount</option>
              <option value="amount_asc">Lowest Amount</option>
            </select>
          </div>
        </div>
      </div>

      {/* Donations List */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        {filteredDonations.length === 0 ? (
          <div className="text-center py-16">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">No donations found</h3>
            <p className="mt-1 text-sm text-gray-500">
              {filter === 'all' ? 'You haven\'t made any donations yet' : `No ${filter} donations found`}
            </p>
            {filter === 'all' && (
              <Link to="/campaigns" className="mt-4 inline-block text-primary-600 hover:text-primary-500 font-medium">
                Browse Campaigns
              </Link>
            )}
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {filteredDonations.map((donation) => (
              <div key={donation.id} className="p-6 hover:bg-gray-50 transition-colors">
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 flex-wrap">
                      <Link to={`/campaigns/${donation.campaignId}`} className="font-semibold text-gray-900 hover:text-primary-600">
                        {donation.campaignTitle}
                      </Link>
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[donation.status]}`}>
                        {donation.status}
                      </span>
                    </div>
                    <p className="mt-1 text-sm text-gray-600">{donation.ngo}</p>
                    <p className="mt-1 text-xs text-gray-500">
                      Transaction: {donation.transactionId} • {formatDate(donation.date)}
                    </p>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <p className="text-2xl font-bold text-gray-900">{formatCurrency(donation.amount)}</p>
                      <p className="text-sm text-gray-500">{formatDate(donation.date)}</p>
                    </div>
                    <div className="flex gap-2">
                      <button className="px-3 py-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 border border-primary-300 rounded-lg transition-colors">
                        View Receipt
                      </button>
                      <Link
                        to={`/campaigns/${donation.campaignId}`}
                        className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors"
                      >
                        Campaign
                      </Link>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Donations</p>
          <p className="mt-2 text-3xl font-bold text-gray-900">{mockDonations.length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Campaigns Supported</p>
          <p className="mt-2 text-3xl font-bold text-gray-900">{new Set(mockDonations.map(d => d.campaignId)).size}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">NGOs Supported</p>
          <p className="mt-2 text-3xl font-bold text-gray-900">{new Set(mockDonations.map(d => d.ngo)).size}</p>
        </div>
      </div>
    </div>
  )
}