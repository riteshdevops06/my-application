import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const mockCampaigns = [
  {
    id: 1,
    title: 'Blindness Treatment & Eye Surgery',
    category: 'Healthcare',
    status: 'active',
    raised: 45000,
    goal: 100000,
    donors: 124,
    startDate: '2024-01-15',
    endDate: '2024-12-31',
    verified: true,
  },
  {
    id: 2,
    title: 'Animal Vaccination Drive',
    category: 'Animal Welfare',
    status: 'active',
    raised: 12000,
    goal: 25000,
    donors: 56,
    startDate: '2024-02-01',
    endDate: '2024-11-15',
    verified: true,
  },
  {
    id: 3,
    title: 'Cancer Treatment Assistance',
    category: 'Healthcare',
    status: 'draft',
    raised: 0,
    goal: 150000,
    donors: 0,
    startDate: '2024-03-01',
    endDate: '2025-01-31',
    verified: false,
  },
]

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const getProgress = (raised, goal) => goal > 0 ? Math.min(Math.round((raised / goal) * 100), 100) : 0

const statusColors = {
  active: 'bg-green-100 text-green-800',
  draft: 'bg-gray-100 text-gray-800',
  completed: 'bg-blue-100 text-blue-800',
  paused: 'bg-yellow-100 text-yellow-800',
}

export default function NGOCampaigns() {
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [campaignToDelete, setCampaignToDelete] = useState(null)

  const handleDelete = async (campaign) => {
    setCampaignToDelete(campaign)
    setShowDeleteModal(true)
  }

  const confirmDelete = async () => {
    if (!campaignToDelete) return
    try {
      // TODO: Call delete API
      await new Promise((resolve) => setTimeout(resolve, 1000))
      notifySuccess('Campaign Deleted', `${campaignToDelete.title} has been deleted`)
      setShowDeleteModal(false)
      setCampaignToDelete(null)
    } catch (err) {
      notifyError('Delete Failed', err.message || 'Failed to delete campaign')
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">My Campaigns</h1>
          <p className="mt-1 text-gray-600">Manage your fundraising campaigns</p>
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

      {/* Campaigns Table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Campaign</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Progress</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donors</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Dates</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {mockCampaigns.map((campaign) => (
                <tr key={campaign.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div>
                      <Link to={`/ngo/campaigns/${campaign.id}/edit`} className="font-medium text-gray-900 hover:text-primary-600">
                        {campaign.title}
                      </Link>
                      {!campaign.verified && (
                        <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
                          Unverified
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{campaign.category}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[campaign.status]}`}>
                      {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="w-48">
                      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full transition-all ${
                            campaign.status === 'active' ? 'bg-primary-600' : 'bg-gray-400'
                          }`}
                          style={{ width: `${getProgress(campaign.raised, campaign.goal)}%` }}
                        />
                      </div>
                      <p className="mt-1 text-xs text-gray-500">
                        {formatCurrency(campaign.raised)} / {formatCurrency(campaign.goal)}
                      </p>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{campaign.donors}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    <p>{new Date(campaign.startDate).toLocaleDateString()}</p>
                    <p className="text-gray-400">to {new Date(campaign.endDate).toLocaleDateString()}</p>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <Link
                        to={`/ngo/campaigns/${campaign.id}/edit`}
                        className="px-3 py-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 border border-primary-300 rounded-lg transition-colors"
                      >
                        Edit
                      </Link>
                      <Link
                        to={`/ngo/campaigns/${campaign.id}/milestones`}
                        className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors"
                      >
                        Milestones
                      </Link>
                      <Link
                        to={`/ngo/campaigns/${campaign.id}/expenses`}
                        className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors"
                      >
                        Expenses
                      </Link>
                      <button
                        onClick={() => handleDelete(campaign)}
                        className="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 border border-red-300 rounded-lg transition-colors"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {mockCampaigns.length === 0 && (
          <div className="text-center py-16">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">No campaigns yet</h3>
            <p className="mt-1 text-sm text-gray-500">Create your first campaign to start fundraising</p>
            <Link to="/ngo/campaigns/create" className="mt-4 inline-block text-primary-600 hover:text-primary-500 font-medium">
              Create Campaign
            </Link>
          </div>
        )}
      </div>

      {/* Delete Modal */}
      {showDeleteModal && campaignToDelete && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowDeleteModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900">Delete Campaign</h2>
              <p className="mt-2 text-gray-600">Are you sure you want to delete <strong>{campaignToDelete.title}</strong>? This action cannot be undone.</p>
              <div className="mt-6 flex gap-3 justify-end">
                <button
                  onClick={() => setShowDeleteModal(false)}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmDelete}
                  className="px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition-colors"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}