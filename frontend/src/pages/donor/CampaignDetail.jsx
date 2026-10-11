import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useNotifications } from '../../context/NotificationContext'

const mockCampaign = {
  id: 1,
  title: 'Blindness Treatment & Eye Surgery',
  ngo: 'Vision Care Foundation',
  ngoDescription: 'Dedicated to preventing blindness and providing eye care to underserved communities.',
  category: 'Healthcare',
  image: 'https://images.unsplash.com/photo-1576091160550-2173dba999ef?w=800&h=400&fit=crop',
  description: 'Millions of people in developing countries suffer from preventable blindness. This campaign funds cataract surgeries, glaucoma treatments, and provides corrective lenses to those in need. Every $50 can restore sight to one person.',
  raised: 45000,
  goal: 100000,
  status: 'active',
  endDate: '2024-12-31',
  startDate: '2024-01-15',
  location: 'Kenya, Uganda, Tanzania',
  milestones: [
    { id: 1, title: 'Phase 1: Mobile Eye Clinics Setup', targetAmount: 20000, raised: 20000, status: 'completed', description: 'Established 3 mobile eye clinics in rural areas' },
    { id: 2, title: 'Phase 2: Cataract Surgeries', targetAmount: 50000, raised: 25000, status: 'in_progress', description: 'Perform 500 cataract surgeries' },
    { id: 3, title: 'Phase 3: Follow-up Care & Glasses', targetAmount: 30000, raised: 0, status: 'pending', description: 'Post-operative care and corrective lenses' },
  ],
  expenses: [
    { id: 1, title: 'Medical Equipment', amount: 15000, category: 'Equipment', date: '2024-02-01', status: 'verified' },
    { id: 2, title: 'Medical Supplies', amount: 8000, category: 'Supplies', date: '2024-02-15', status: 'verified' },
    { id: 3, title: 'Staff Salaries', amount: 12000, category: 'Personnel', date: '2024-03-01', status: 'verified' },
    { id: 4, title: 'Transportation', amount: 5000, category: 'Operations', date: '2024-03-10', status: 'pending' },
  ],
  updates: [
    { id: 1, title: 'First Mobile Clinic Operational', date: '2024-02-20', content: 'Our first mobile eye clinic has started operations in rural Kenya, serving 50+ patients per week.' },
    { id: 2, title: '100 Cataract Surgeries Completed', date: '2024-04-15', content: 'We have successfully completed 100 cataract surgeries with a 98% success rate.' },
  ],
}

const donationTiers = [
  { amount: 25, title: 'Vision Supporter', description: 'Provides basic eye screening for 2 people' },
  { amount: 50, title: 'Sight Restorer', description: 'Funds one cataract surgery', recommended: true },
  { amount: 100, title: 'Vision Champion', description: 'Funds surgery + post-op care for 1 person' },
  { amount: 250, title: 'Sight Ambassador', description: 'Funds a full day of mobile clinic operations' },
  { amount: 500, title: 'Vision Hero', description: 'Funds a week of mobile clinic operations' },
]

export default function DonorCampaignDetail() {
  const { id } = useParams()
  const { user, isAuthenticated } = useAuth()
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [campaign, setCampaign] = useState(mockCampaign)
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedTier, setSelectedTier] = useState(50)
  const [customAmount, setCustomAmount] = useState('')
  const [donationLoading, setDonationLoading] = useState(false)
  const [showDonationModal, setShowDonationModal] = useState(false)

  const getProgress = (raised, goal) => Math.min(Math.round((raised / goal) * 100), 100)
  const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)

  const handleDonate = async () => {
    if (!isAuthenticated) {
      notifyError('Login Required', 'Please log in to make a donation')
      return
    }
    const amount = customAmount || selectedTier
    if (!amount || amount < 1) {
      notifyError('Invalid Amount', 'Please enter a valid donation amount')
      return
    }
    setDonationLoading(true)
    try {
      // TODO: Call donation API
      await new Promise((resolve) => setTimeout(resolve, 1500))
      notifySuccess('Donation Successful!', `Thank you for donating ${formatCurrency(amount)}!`)
      setShowDonationModal(false)
      setCustomAmount('')
    } catch (err) {
      notifyError('Donation Failed', err.message || 'Please try again')
    } finally {
      setDonationLoading(false)
    }
  }

  const tabs = [
    { id: 'overview', label: 'Overview', icon: 'info' },
    { id: 'milestones', label: 'Milestones', icon: 'flag' },
    { id: 'expenses', label: 'Expenses', icon: 'receipt' },
    { id: 'updates', label: 'Updates', icon: 'news' },
  ]

  const tabIcons = {
    info: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
    flag: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
    receipt: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>,
    news: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" /></svg>,
  }

  const statusColors = {
    completed: 'bg-green-100 text-green-800',
    in_progress: 'bg-blue-100 text-blue-800',
    pending: 'bg-yellow-100 text-yellow-800',
  }

  const expenseStatusColors = {
    verified: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    rejected: 'bg-red-100 text-red-800',
  }

  return (
    <div className="space-y-6">
      {/* Campaign Header */}
      <div className="relative">
        <img src={campaign.image} alt={campaign.title} className="w-full h-64 object-cover" />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
        <div className="absolute bottom-0 left-0 right-0 p-6 text-white">
          <div className="max-w-7xl mx-auto flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
            <div>
              <span className="inline-block px-3 py-1 text-sm font-medium bg-primary-600 rounded-full mb-2">{campaign.category}</span>
              <h1 className="text-3xl font-bold">{campaign.title}</h1>
              <p className="mt-1 text-gray-300">by {campaign.ngo}</p>
            </div>
            <div className="w-full sm:w-auto">
              <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 text-center">
                <p className="text-sm text-gray-300">Raised</p>
                <p className="text-4xl font-bold">{formatCurrency(campaign.raised)}</p>
                <p className="text-sm text-gray-300 mt-1">of {formatCurrency(campaign.goal)} goal</p>
                <div className="mt-4 h-3 bg-white/20 rounded-full overflow-hidden">
                  <div className="h-full bg-primary-400 rounded-full" style={{ width: `${getProgress(campaign.raised, campaign.goal)}%` }} />
                </div>
                <p className="mt-2 text-sm text-gray-300">{getProgress(campaign.raised, campaign.goal)}% funded</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Donate Button */}
        <div className="sticky top-20 z-10 sm:static">
          <button
            onClick={() => setShowDonationModal(true)}
            className="w-full sm:w-auto bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold text-lg hover:bg-primary-700 transition-colors flex items-center justify-center gap-2"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.154a4 4 0 110 5.692M15.211 4.154a4 4 0 110 5.692m6.656 10.879a4 4 0 11-5.692 0m5.692 0a4 4 0 110-5.692m-5.692 0a4 4 0 110 5.692m0-11.384a4 4 0 115.692 0" />
            </svg>
            Donate Now
          </button>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-200">
          <nav className="flex gap-8" aria-label="Campaign tabs">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === tab.id
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                {tabIcons[tab.icon]}
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        <div className="py-6">
          {activeTab === 'overview' && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              <div className="lg:col-span-2 space-y-6">
                <div>
                  <h2 className="text-xl font-bold text-gray-900">About this Campaign</h2>
                  <p className="mt-4 text-gray-600 leading-relaxed">{campaign.description}</p>
                </div>
                <div>
                  <h2 className="text-xl font-bold text-gray-900">About {campaign.ngo}</h2>
                  <p className="mt-4 text-gray-600">{campaign.ngoDescription}</p>
                </div>
                <div>
                  <h2 className="text-xl font-bold text-gray-900">Campaign Details</h2>
                  <dl className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                    <div><dt className="text-gray-500">Status</dt><dd className="font-medium capitalize">{campaign.status.replace('_', ' ')}</dd></div>
                    <div><dt className="text-gray-500">Location</dt><dd className="font-medium">{campaign.location}</dd></div>
                    <div><dt className="text-gray-500">Start Date</dt><dd className="font-medium">{new Date(campaign.startDate).toLocaleDateString()}</dd></div>
                    <div><dt className="text-gray-500">End Date</dt><dd className="font-medium">{new Date(campaign.endDate).toLocaleDateString()}</dd></div>
                    <div><dt className="text-gray-500">Category</dt><dd className="font-medium">{campaign.category}</dd></div>
                  </dl>
                </div>
              </div>
              <div className="space-y-6">
                <div className="bg-white rounded-xl border border-gray-200 p-6 sticky top-24">
                  <h3 className="text-lg font-semibold text-gray-900">Donation Tiers</h3>
                  <p className="mt-2 text-sm text-gray-600">Choose a tier or enter a custom amount</p>
                  <div className="mt-4 space-y-2">
                    {donationTiers.map((tier) => (
                      <button
                        key={tier.amount}
                        onClick={() => { setSelectedTier(tier.amount); setCustomAmount('') }}
                        className={`w-full p-4 rounded-lg border-2 text-left transition-all ${
                          selectedTier === tier.amount && !customAmount
                            ? 'border-primary-500 bg-primary-50'
                            : 'border-gray-200 hover:border-gray-300'
                        }`}
                      >
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="font-medium text-gray-900">{tier.title}</p>
                            <p className="text-sm text-gray-600">{tier.description}</p>
                          </div>
                          <span className="text-xl font-bold text-primary-600">${tier.amount}</span>
                        </div>
                        {tier.recommended && <p className="mt-1 text-xs text-primary-600 font-medium">Recommended</p>}
                      </button>
                    ))}
                    <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
                      <label className="block text-sm font-medium text-gray-700">Custom Amount</label>
                      <div className="mt-2 flex items-center">
                        <span className="text-gray-500">$</span>
                        <input
                          type="number"
                          min="1"
                          value={customAmount}
                          onChange={(e) => { setCustomAmount(e.target.value); setSelectedTier(null) }}
                          className="flex-1 ml-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                          placeholder="Enter amount"
                        />
                      </div>
                    </div>
                  </div>
                  <button
                    onClick={handleDonate}
                    disabled={donationLoading}
                    className="mt-6 w-full bg-primary-600 text-white py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors disabled:opacity-50"
                  >
                    {donationLoading ? 'Processing...' : `Donate ${formatCurrency(customAmount || selectedTier)}`}
                  </button>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'milestones' && (
            <div className="space-y-4">
              <h2 className="text-xl font-bold text-gray-900">Campaign Milestones</h2>
              <p className="text-gray-600">Track progress through each phase of the campaign</p>
              <div className="mt-4 space-y-4">
                {campaign.milestones.map((milestone) => (
                  <div key={milestone.id} className="bg-white rounded-xl border border-gray-200 p-6">
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3">
                          <span className={`px-3 py-1 text-sm font-medium rounded-full ${statusColors[milestone.status]}`}>
                            {milestone.status.replace('_', ' ')}
                          </span>
                          <h3 className="text-lg font-semibold text-gray-900">{milestone.title}</h3>
                        </div>
                        <p className="mt-2 text-gray-600">{milestone.description}</p>
                        <div className="mt-3">
                          <div className="flex justify-between text-sm mb-1">
                            <span>Raised: {formatCurrency(milestone.raised)}</span>
                            <span>Target: {formatCurrency(milestone.targetAmount)}</span>
                          </div>
                          <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                            <div className={`h-full rounded-full transition-all ${
                              milestone.status === 'completed' ? 'bg-green-500' : 'bg-primary-500'
                            }`} style={{ width: `${Math.min(Math.round((milestone.raised / milestone.targetAmount) * 100), 100)}%` }} />
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'expenses' && (
            <div className="space-y-4">
              <h2 className="text-xl font-bold text-gray-900">Campaign Expenses</h2>
              <p className="text-gray-600">Transparent breakdown of how funds are spent</p>
              <div className="mt-4 bg-white rounded-xl border border-gray-200 overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Expense</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {campaign.expenses.map((expense) => (
                      <tr key={expense.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 text-sm font-medium text-gray-900">{expense.title}</td>
                        <td className="px-6 py-4 text-sm text-gray-600">{expense.category}</td>
                        <td className="px-6 py-4 text-sm font-medium text-gray-900">{formatCurrency(expense.amount)}</td>
                        <td className="px-6 py-4 text-sm text-gray-600">{new Date(expense.date).toLocaleDateString()}</td>
                        <td className="px-6 py-4">
                          <span className={`px-2 py-1 text-xs font-medium rounded-full ${expenseStatusColors[expense.status]}`}>
                            {expense.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {activeTab === 'updates' && (
            <div className="space-y-4">
              <h2 className="text-xl font-bold text-gray-900">Campaign Updates</h2>
              <p className="text-gray-600">Latest progress reports from the NGO</p>
              <div className="mt-4 space-y-4">
                {campaign.updates.map((update) => (
                  <div key={update.id} className="bg-white rounded-xl border border-gray-200 p-6">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <h3 className="text-lg font-semibold text-gray-900">{update.title}</h3>
                        <p className="mt-2 text-gray-600">{update.content}</p>
                      </div>
                      <span className="text-sm text-gray-500 whitespace-nowrap">{new Date(update.date).toLocaleDateString()}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Donation Modal */}
      {showDonationModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowDonationModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900">Confirm Donation</h2>
              <p className="mt-2 text-gray-600">You are donating to <strong>{campaign.title}</strong></p>
              <div className="mt-4 p-4 bg-gray-50 rounded-lg">
                <div className="flex justify-between text-lg">
                  <span className="font-medium">Amount</span>
                  <span className="font-bold text-primary-600">{formatCurrency(customAmount || selectedTier)}</span>
                </div>
              </div>
              <div className="mt-6 flex gap-3">
                <button
                  onClick={() => setShowDonationModal(false)}
                  className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDonate}
                  disabled={donationLoading}
                  className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg font-medium hover:bg-primary-700 transition-colors disabled:opacity-50"
                >
                  {donationLoading ? 'Processing...' : 'Confirm Donation'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}