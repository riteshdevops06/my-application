import { useState } from 'react'
import { Link } from 'react-router-dom'

const mockCampaigns = [
  {
    id: 1,
    title: 'Blindness Treatment & Eye Surgery',
    ngo: 'Vision Care Foundation',
    category: 'Healthcare',
    image: 'https://images.unsplash.com/photo-1576091160550-2173dba999ef?w=400&h=250&fit=crop',
    raised: 45000,
    goal: 100000,
    status: 'active',
    endDate: '2024-12-31',
  },
  {
    id: 2,
    title: 'Animal Vaccination Drive',
    ngo: 'Paws & Claws Rescue',
    category: 'Animal Welfare',
    image: 'https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?w=400&h=250&fit=crop',
    raised: 12000,
    goal: 25000,
    status: 'active',
    endDate: '2024-11-15',
  },
  {
    id: 3,
    title: 'Cancer Treatment Assistance',
    ngo: 'Hope Cancer Support',
    category: 'Healthcare',
    image: 'https://images.unsplash.com/photo-1576091160399-112ba8d25d1d?w=400&h=250&fit=crop',
    raised: 78000,
    goal: 150000,
    status: 'active',
    endDate: '2025-01-31',
  },
  {
    id: 4,
    title: 'Clean Water for Rural Villages',
    ngo: 'Water Life Initiative',
    category: 'Community Development',
    image: 'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=400&h=250&fit=crop',
    raised: 25000,
    goal: 50000,
    status: 'active',
    endDate: '2024-10-30',
  },
  {
    id: 5,
    title: 'Education for Underprivileged Children',
    ngo: 'Bright Future Academy',
    category: 'Education',
    image: 'https://images.unsplash.com/photo-1503676260728-1c00da094a0b?w=400&h=250&fit=crop',
    raised: 60000,
    goal: 120000,
    status: 'completed',
    endDate: '2024-08-15',
  },
  {
    id: 6,
    title: 'Disaster Relief Fund',
    ngo: 'Global Aid Network',
    category: 'Emergency Relief',
    image: 'https://images.unsplash.com/photo-1590080877112-5b9b0e6b8b7b?w=400&h=250&fit=crop',
    raised: 200000,
    goal: 500000,
    status: 'active',
    endDate: '2024-12-01',
  },
]

const categories = ['All', 'Healthcare', 'Education', 'Animal Welfare', 'Community Development', 'Emergency Relief']
const statuses = ['All', 'Active', 'Completed', 'Upcoming']

export default function DonorCampaigns() {
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('All')
  const [status, setStatus] = useState('All')
  const [view, setView] = useState('grid') // grid or list

  const filteredCampaigns = mockCampaigns.filter((campaign) => {
    const matchesSearch = campaign.title.toLowerCase().includes(search.toLowerCase()) ||
      campaign.ngo.toLowerCase().includes(search.toLowerCase())
    const matchesCategory = category === 'All' || campaign.category === category
    const matchesStatus = status === 'All' || campaign.status === status.toLowerCase()
    return matchesSearch && matchesCategory && matchesStatus
  })

  const getProgress = (raised, goal) => Math.min(Math.round((raised / goal) * 100), 100)

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
  }

  const getStatusBadge = (status) => {
    const styles = {
      active: 'bg-green-100 text-green-800',
      completed: 'bg-blue-100 text-blue-800',
      upcoming: 'bg-yellow-100 text-yellow-800',
    }
    return styles[status] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Browse Campaigns</h1>
        <p className="mt-1 text-gray-600">Discover and support causes that matter to you</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4 sm:p-6">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="relative flex-1">
            <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              placeholder="Search campaigns..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
          <select
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          >
            {categories.map((cat) => (
              <option key={cat} value={cat}>{cat}</option>
            ))}
          </select>
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          >
            {statuses.map((s) => (
              <option key={s} value={s}>{s}</option>
            ))}
          </select>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => setView('grid')}
              className={`p-2 rounded-lg ${view === 'grid' ? 'bg-primary-100 text-primary-600' : 'text-gray-400 hover:text-gray-600'}`}
              aria-label="Grid view"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
              </svg>
            </button>
            <button
              onClick={() => setView('list')}
              className={`p-2 rounded-lg ${view === 'list' ? 'bg-primary-100 text-primary-600' : 'text-gray-400 hover:text-gray-600'}`}
              aria-label="List view"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Campaigns Grid/List */}
      <div className={`grid gap-6 ${view === 'grid' ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3' : 'grid-cols-1'}`}>
        {filteredCampaigns.map((campaign) => {
          const progress = getProgress(campaign.raised, campaign.goal)

          if (view === 'grid') {
            return (
              <Link
                key={campaign.id}
                to={`/campaigns/${campaign.id}`}
                className="bg-white rounded-xl border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow"
              >
                <div className="relative h-48 bg-gray-200">
                  <img
                    src={campaign.image}
                    alt={campaign.title}
                    className="w-full h-full object-cover"
                  />
                  <span className={`absolute top-3 right-3 px-2 py-1 text-xs font-medium rounded-full ${getStatusBadge(campaign.status)}`}>
                    {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                  </span>
                </div>
                <div className="p-5">
                  <p className="text-sm text-gray-500">{campaign.category}</p>
                  <h3 className="mt-1 text-lg font-semibold text-gray-900 line-clamp-2">{campaign.title}</h3>
                  <p className="mt-1 text-sm text-gray-600">{campaign.ngo}</p>
                  <div className="mt-4">
                    <div className="flex items-center justify-between text-sm mb-1">
                      <span className="font-medium">{formatCurrency(campaign.raised)}</span>
                      <span className="text-gray-500">of {formatCurrency(campaign.goal)}</span>
                    </div>
                    <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-primary-600 rounded-full transition-all"
                        style={{ width: `${progress}%` }}
                      />
                    </div>
                    <p className="mt-1 text-xs text-gray-500">{progress}% funded • Ends {new Date(campaign.endDate).toLocaleDateString()}</p>
                  </div>
                </div>
              </Link>
            )
          }

          return (
            <div key={campaign.id} className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <div className="flex flex-col sm:flex-row">
                <div className="relative w-full sm:w-64 h-48 sm:h-auto bg-gray-200 flex-shrink-0">
                  <img
                    src={campaign.image}
                    alt={campaign.title}
                    className="w-full h-full object-cover"
                  />
                  <span className={`absolute top-3 right-3 px-2 py-1 text-xs font-medium rounded-full ${getStatusBadge(campaign.status)}`}>
                    {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                  </span>
                </div>
                <div className="p-5 flex flex-col justify-between">
                  <div>
                    <p className="text-sm text-gray-500">{campaign.category}</p>
                    <h3 className="mt-1 text-lg font-semibold text-gray-900">{campaign.title}</h3>
                    <p className="mt-1 text-sm text-gray-600">{campaign.ngo}</p>
                  </div>
                  <div className="mt-4 sm:mt-0 sm:ml-6 flex-1">
                    <div className="flex items-center justify-between text-sm mb-1">
                      <span className="font-medium">{formatCurrency(campaign.raised)}</span>
                      <span className="text-gray-500">of {formatCurrency(campaign.goal)}</span>
                    </div>
                    <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-primary-600 rounded-full transition-all"
                        style={{ width: `${progress}%` }}
                      />
                    </div>
                    <p className="mt-1 text-xs text-gray-500">{progress}% funded • Ends {new Date(campaign.endDate).toLocaleDateString()}</p>
                    <Link
                      to={`/campaigns/${campaign.id}`}
                      className="mt-4 inline-block text-primary-600 hover:text-primary-500 font-medium text-sm"
                    >
                      View Details →
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {filteredCampaigns.length === 0 && (
        <div className="text-center py-16">
          <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No campaigns found</h3>
          <p className="mt-1 text-sm text-gray-500">Try adjusting your filters or search terms</p>
        </div>
      )}
    </div>
  )
}