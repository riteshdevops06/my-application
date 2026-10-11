import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const mockCampaigns = [
  {
    id: 1,
    title: 'Blindness Treatment & Eye Surgery',
    ngo: 'Vision Care Foundation',
    ngoId: 1,
    category: 'Healthcare',
    status: 'active',
    raised: 45000,
    goal: 100000,
    donors: 124,
    startDate: '2024-01-15',
    endDate: '2024-12-31',
    approved: true,
    approvedDate: '2024-01-20',
    flagged: false,
    milestones: 3,
    completedMilestones: 1,
    pendingExpenses: 2,
    location: 'Kenya, Uganda, Tanzania',
  },
  {
    id: 2,
    title: 'Animal Vaccination Drive',
    ngo: 'Paws & Claws Rescue',
    ngoId: 2,
    category: 'Animal Welfare',
    status: 'active',
    raised: 12000,
    goal: 25000,
    donors: 56,
    startDate: '2024-02-01',
    endDate: '2024-11-15',
    approved: true,
    approvedDate: '2024-02-03',
    flagged: false,
    milestones: 2,
    completedMilestones: 0,
    pendingExpenses: 1,
    location: 'Lagos, Nigeria',
  },
  {
    id: 3,
    title: 'Cancer Treatment Assistance',
    ngo: 'Hope Cancer Support',
    ngoId: 3,
    category: 'Healthcare',
    status: 'pending_approval',
    raised: 0,
    goal: 150000,
    donors: 0,
    startDate: '2024-03-01',
    endDate: '2025-01-31',
    approved: false,
    flagged: false,
    milestones: 4,
    completedMilestones: 0,
    pendingExpenses: 0,
    location: 'Multiple locations',
  },
  {
    id: 4,
    title: 'Emergency Flood Relief',
    ngo: 'Global Aid Network',
    ngoId: 4,
    category: 'Emergency Relief',
    status: 'active',
    raised: 200000,
    goal: 500000,
    donors: 892,
    startDate: '2024-02-10',
    endDate: '2024-06-30',
    approved: true,
    approvedDate: '2024-02-10',
    flagged: true,
    flagReason: 'Unusually high donation velocity requires review',
    milestones: 1,
    completedMilestones: 1,
    pendingExpenses: 5,
    location: 'East Africa',
  },
]

const statusColors = {
  active: 'bg-green-100 text-green-800',
  pending_approval: 'bg-yellow-100 text-yellow-800',
  paused: 'bg-orange-100 text-orange-800',
  completed: 'bg-blue-100 text-blue-800',
  rejected: 'bg-red-100 text-red-800',
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
const getProgress = (raised, goal) => goal > 0 ? Math.min(Math.round((raised / goal) * 100), 100) : 0

export default function PlatformAdminCampaigns() {
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [campaigns, setCampaigns] = useState(mockCampaigns)
  const [filter, setFilter] = useState('all')
  const [selectedCampaign, setSelectedCampaign] = useState(null)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [showApprovalModal, setShowApprovalModal] = useState(false)
  const [showSuspendModal, setShowSuspendModal] = useState(false)
  const [suspendReason, setSuspendReason] = useState('')
  const [loading, setLoading] = useState(false)

  const filteredCampaigns = campaigns.filter(campaign => {
    if (filter === 'all') return true
    if (filter === 'flagged') return campaign.flagged
    if (filter === 'high_activity') return campaign.donors > 500
    return campaign.status === filter
  })

  const handleViewDetails = (campaign) => {
    setSelectedCampaign(campaign)
    setShowDetailsModal(true)
  }

  const handleApprove = (campaign) => {
    setSelectedCampaign(campaign)
    setShowApprovalModal(true)
  }

  const handleSuspend = (campaign) => {
    setSelectedCampaign(campaign)
    setShowSuspendModal(true)
  }

  const confirmApproval = async () => {
    if (!selectedCampaign) return
    setLoading(true)
    try {
      // TODO: Call API to approve campaign
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setCampaigns(prev => prev.map(campaign => 
        campaign.id === selectedCampaign.id 
          ? { ...campaign, status: 'active', approved: true, approvedDate: new Date().toISOString().split('T')[0] }
          : campaign
      ))
      
      notifySuccess('Campaign Approved', `${selectedCampaign.title} has been approved and is now active`)
      setShowApprovalModal(false)
      setSelectedCampaign(null)
    } catch (err) {
      notifyError('Approval Failed', err.message || 'Failed to approve campaign')
    } finally {
      setLoading(false)
    }
  }

  const confirmSuspend = async () => {
    if (!selectedCampaign || !suspendReason.trim()) return
    setLoading(true)
    try {
      // TODO: Call API to suspend campaign
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setCampaigns(prev => prev.map(campaign => 
        campaign.id === selectedCampaign.id 
          ? { ...campaign, status: 'paused', flagged: true, flagReason: suspendReason }
          : campaign
      ))
      
      notifySuccess('Campaign Suspended', `${selectedCampaign.title} has been suspended`)
      setShowSuspendModal(false)
      setSelectedCampaign(null)
      setSuspendReason('')
    } catch (err) {
      notifyError('Suspension Failed', err.message || 'Failed to suspend campaign')
    } finally {
      setLoading(false)
    }
  }

  const getStatusIcon = (campaign) => {
    if (campaign.flagged) {
      return <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" /></svg>
    }
    
    switch (campaign.status) {
      case 'active':
        return <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
      case 'pending_approval':
        return <svg className="w-4 h-4 text-yellow-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" /></svg>
      case 'paused':
        return <svg className="w-4 h-4 text-orange-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zM7 8a1 1 0 012 0v4a1 1 0 01-2 0V8zm5-1a1 1 0 00-1 1v4a1 1 0 002 0V8a1 1 0 00-1-1z" clipRule="evenodd" /></svg>
      default:
        return <svg className="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" /></svg>
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Campaign Oversight</h1>
          <p className="mt-1 text-gray-600">Monitor and manage all platform campaigns</p>
        </div>
        <div className="flex items-center gap-4 text-sm text-gray-500">
          <span className="font-medium">Total Campaigns:</span>
          <span>{campaigns.length}</span>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-5 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Active</p>
          <p className="text-2xl font-bold text-green-600">{campaigns.filter(c => c.status === 'active').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Pending</p>
          <p className="text-2xl font-bold text-yellow-600">{campaigns.filter(c => c.status === 'pending_approval').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Flagged</p>
          <p className="text-2xl font-bold text-red-600">{campaigns.filter(c => c.flagged).length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Raised</p>
          <p className="text-xl font-bold text-primary-600">{formatCurrency(campaigns.reduce((sum, c) => sum + c.raised, 0))}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Donors</p>
          <p className="text-2xl font-bold text-gray-900">{campaigns.reduce((sum, c) => sum + c.donors, 0).toLocaleString()}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex gap-2 flex-wrap">
          {['all', 'active', 'pending_approval', 'paused', 'flagged', 'high_activity'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${
                filter === status
                  ? 'bg-primary-100 text-primary-700'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {status === 'all' ? 'All Campaigns' : 
               status === 'pending_approval' ? 'Pending Approval' :
               status === 'high_activity' ? 'High Activity' :
               status.charAt(0).toUpperCase() + status.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Campaigns Table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Campaign</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">NGO</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Progress</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donors</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Timeline</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {filteredCampaigns.map((campaign) => (
                <tr key={campaign.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="font-medium text-gray-900">{campaign.title}</p>
                        {campaign.flagged && (
                          <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                          </svg>
                        )}
                      </div>
                      <p className="text-sm text-gray-500">{campaign.category}</p>
                      <p className="text-xs text-gray-400">{campaign.location}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <Link
                      to={`/admin/ngos`}
                      className="text-sm font-medium text-primary-600 hover:text-primary-700"
                    >
                      {campaign.ngo}
                    </Link>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      {getStatusIcon(campaign)}
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[campaign.status]}`}>
                        {campaign.status === 'pending_approval' ? 'Pending' : 
                         campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                      </span>
                    </div>
                    {campaign.flagged && campaign.flagReason && (
                      <p className="text-xs text-red-600 mt-1">⚠ {campaign.flagReason}</p>
                    )}
                  </td>
                  <td className="px-6 py-4">
                    <div className="w-32">
                      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-primary-600 rounded-full transition-all"
                          style={{ width: `${getProgress(campaign.raised, campaign.goal)}%` }}
                        />
                      </div>
                      <p className="text-xs text-gray-500 mt-1">
                        {formatCurrency(campaign.raised)} / {formatCurrency(campaign.goal)}
                      </p>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{campaign.donors.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    <p>{formatDate(campaign.startDate)}</p>
                    <p className="text-gray-400">to {formatDate(campaign.endDate)}</p>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleViewDetails(campaign)}
                        className="px-3 py-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 border border-primary-300 rounded-lg transition-colors hover:bg-primary-50"
                      >
                        Details
                      </button>
                      {campaign.status === 'pending_approval' && (
                        <button
                          onClick={() => handleApprove(campaign)}
                          className="px-3 py-1.5 text-sm font-medium text-green-600 hover:text-green-700 border border-green-300 rounded-lg transition-colors hover:bg-green-50"
                        >
                          Approve
                        </button>
                      )}
                      {(campaign.status === 'active' || campaign.flagged) && (
                        <button
                          onClick={() => handleSuspend(campaign)}
                          className="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 border border-red-300 rounded-lg transition-colors hover:bg-red-50"
                        >
                          Suspend
                        </button>
                      )}
                      <Link
                        to={`/admin/milestones?campaign=${campaign.id}`}
                        className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors hover:bg-gray-50"
                      >
                        Milestones
                      </Link>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Details Modal */}
      {showDetailsModal && selectedCampaign && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowDetailsModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-2xl shadow-xl max-h-[90vh] overflow-y-auto">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold text-gray-900">Campaign Details</h2>
                <button
                  onClick={() => setShowDetailsModal(false)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <div className="space-y-6">
                <div>
                  <h3 className="font-semibold text-lg text-gray-900">{selectedCampaign.title}</h3>
                  <p className="text-gray-600">by {selectedCampaign.ngo}</p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Category</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedCampaign.category}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Status</label>
                    <div className="mt-1 flex items-center gap-2">
                      {getStatusIcon(selectedCampaign)}
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[selectedCampaign.status]}`}>
                        {selectedCampaign.status === 'pending_approval' ? 'Pending' : 
                         selectedCampaign.status.charAt(0).toUpperCase() + selectedCampaign.status.slice(1)}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Goal</label>
                    <p className="mt-1 text-lg font-bold text-gray-900">{formatCurrency(selectedCampaign.goal)}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Raised</label>
                    <p className="mt-1 text-lg font-bold text-primary-600">{formatCurrency(selectedCampaign.raised)}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Progress</label>
                    <p className="mt-1 text-lg font-bold text-gray-900">{getProgress(selectedCampaign.raised, selectedCampaign.goal)}%</p>
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Milestones</label>
                    <p className="mt-1 text-sm text-gray-900">
                      {selectedCampaign.completedMilestones} of {selectedCampaign.milestones} completed
                    </p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Pending Expenses</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedCampaign.pendingExpenses} awaiting review</p>
                  </div>
                </div>

                {selectedCampaign.flagged && selectedCampaign.flagReason && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Flag Reason</label>
                    <div className="p-3 bg-red-50 rounded-lg">
                      <p className="text-sm text-red-800">{selectedCampaign.flagReason}</p>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Approval Modal */}
      {showApprovalModal && selectedCampaign && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowApprovalModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Approve Campaign</h2>
              <p className="text-gray-600 mb-6">
                Are you sure you want to approve <strong>{selectedCampaign.title}</strong>? This will make it visible to donors and allow fundraising.
              </p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => setShowApprovalModal(false)}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmApproval}
                  disabled={loading}
                  className="px-4 py-2 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Approving...' : 'Approve Campaign'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Suspend Modal */}
      {showSuspendModal && selectedCampaign && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowSuspendModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Suspend Campaign</h2>
              <p className="text-gray-600 mb-4">
                Please provide a reason for suspending <strong>{selectedCampaign.title}</strong>:
              </p>
              <textarea
                value={suspendReason}
                onChange={(e) => setSuspendReason(e.target.value)}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                placeholder="Explain why this campaign is being suspended..."
              />
              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => setShowSuspendModal(false)}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmSuspend}
                  disabled={loading || !suspendReason.trim()}
                  className="px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Suspending...' : 'Suspend Campaign'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}