import { useState } from 'react'
import { useNotifications } from '../../context/NotificationContext'

const mockNGOs = [
  {
    id: 1,
    name: 'Vision Care Foundation',
    email: 'contact@visioncare.org',
    category: 'Healthcare',
    registrationNumber: 'NGO-2020-001234',
    status: 'verified',
    verifiedDate: '2024-02-15',
    submittedDate: '2024-01-15',
    activeCampaigns: 3,
    totalRaised: 125000,
    location: 'Nairobi, Kenya',
    website: 'https://visioncare.org',
    phone: '+254-700-123456',
    documents: ['registration_cert', 'tax_cert', 'financial_report', 'board_list'],
  },
  {
    id: 2,
    name: 'Hope Cancer Support',
    email: 'info@hopecancer.org',
    category: 'Healthcare',
    registrationNumber: 'NGO-2023-005678',
    status: 'pending',
    submittedDate: '2024-03-10',
    activeCampaigns: 1,
    totalRaised: 45000,
    location: 'Lagos, Nigeria',
    website: 'https://hopecancer.org',
    phone: '+234-800-987654',
    documents: ['registration_cert', 'tax_cert', 'financial_report'],
    pendingDocuments: ['board_list', 'audit_report'],
  },
  {
    id: 3,
    name: 'Clean Water Initiative',
    email: 'hello@cleanwater.org',
    category: 'Community Development',
    registrationNumber: 'NGO-2022-003456',
    status: 'rejected',
    rejectedDate: '2024-03-05',
    submittedDate: '2024-02-20',
    activeCampaigns: 0,
    totalRaised: 0,
    location: 'Accra, Ghana',
    rejectionReason: 'Financial documents do not meet transparency requirements. Please resubmit with certified audit reports.',
    documents: ['registration_cert'],
  },
  {
    id: 4,
    name: 'Education for All Foundation',
    email: 'support@educationforall.org',
    category: 'Education',
    registrationNumber: 'NGO-2023-007890',
    status: 'incomplete',
    submittedDate: '2024-03-12',
    activeCampaigns: 0,
    totalRaised: 0,
    location: 'Kampala, Uganda',
    documents: ['registration_cert'],
    missingDocuments: ['tax_cert', 'financial_report', 'board_list'],
  },
]

const statusColors = {
  verified: 'bg-green-100 text-green-800',
  pending: 'bg-yellow-100 text-yellow-800',
  rejected: 'bg-red-100 text-red-800',
  incomplete: 'bg-gray-100 text-gray-800',
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

export default function PlatformAdminNGOs() {
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [ngos, setNGOs] = useState(mockNGOs)
  const [filter, setFilter] = useState('all')
  const [selectedNGO, setSelectedNGO] = useState(null)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [showApprovalModal, setShowApprovalModal] = useState(false)
  const [showRejectionModal, setShowRejectionModal] = useState(false)
  const [rejectionReason, setRejectionReason] = useState('')
  const [loading, setLoading] = useState(false)

  const filteredNGOs = ngos.filter(ngo => {
    if (filter === 'all') return true
    return ngo.status === filter
  })

  const handleViewDetails = (ngo) => {
    setSelectedNGO(ngo)
    setShowDetailsModal(true)
  }

  const handleApprove = (ngo) => {
    setSelectedNGO(ngo)
    setShowApprovalModal(true)
  }

  const handleReject = (ngo) => {
    setSelectedNGO(ngo)
    setShowRejectionModal(true)
  }

  const confirmApproval = async () => {
    if (!selectedNGO) return
    setLoading(true)
    try {
      // TODO: Call API to approve NGO
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setNGOs(prev => prev.map(ngo => 
        ngo.id === selectedNGO.id 
          ? { ...ngo, status: 'verified', verifiedDate: new Date().toISOString().split('T')[0] }
          : ngo
      ))
      
      notifySuccess('NGO Approved', `${selectedNGO.name} has been verified and approved`)
      setShowApprovalModal(false)
      setSelectedNGO(null)
    } catch (err) {
      notifyError('Approval Failed', err.message || 'Failed to approve NGO')
    } finally {
      setLoading(false)
    }
  }

  const confirmRejection = async () => {
    if (!selectedNGO || !rejectionReason.trim()) return
    setLoading(true)
    try {
      // TODO: Call API to reject NGO
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setNGOs(prev => prev.map(ngo => 
        ngo.id === selectedNGO.id 
          ? { ...ngo, status: 'rejected', rejectedDate: new Date().toISOString().split('T')[0], rejectionReason }
          : ngo
      ))
      
      notifySuccess('NGO Rejected', `${selectedNGO.name} verification has been rejected`)
      setShowRejectionModal(false)
      setSelectedNGO(null)
      setRejectionReason('')
    } catch (err) {
      notifyError('Rejection Failed', err.message || 'Failed to reject NGO')
    } finally {
      setLoading(false)
    }
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'verified':
        return <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" /></svg>
      case 'pending':
        return <svg className="w-4 h-4 text-yellow-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" /></svg>
      case 'rejected':
        return <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" /></svg>
      default:
        return <svg className="w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" /></svg>
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">NGO Management</h1>
          <p className="mt-1 text-gray-600">Verify, approve, and manage registered NGOs</p>
        </div>
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <span className="font-medium">Total NGOs:</span>
          <span>{ngos.length}</span>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Verified NGOs</p>
          <p className="text-2xl font-bold text-green-600">{ngos.filter(n => n.status === 'verified').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Pending Review</p>
          <p className="text-2xl font-bold text-yellow-600">{ngos.filter(n => n.status === 'pending').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Rejected</p>
          <p className="text-2xl font-bold text-red-600">{ngos.filter(n => n.status === 'rejected').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Incomplete</p>
          <p className="text-2xl font-bold text-gray-600">{ngos.filter(n => n.status === 'incomplete').length}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex gap-2 flex-wrap">
          {['all', 'pending', 'verified', 'rejected', 'incomplete'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${
                filter === status
                  ? 'bg-primary-100 text-primary-700'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {status === 'all' ? 'All NGOs' : status.charAt(0).toUpperCase() + status.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* NGOs List */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">NGO</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Campaigns</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Raised</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {filteredNGOs.map((ngo) => (
                <tr key={ngo.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div>
                      <p className="font-medium text-gray-900">{ngo.name}</p>
                      <p className="text-sm text-gray-500">{ngo.email}</p>
                      <p className="text-xs text-gray-400">{ngo.location}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{ngo.category}</td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      {getStatusIcon(ngo.status)}
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[ngo.status]}`}>
                        {ngo.status.charAt(0).toUpperCase() + ngo.status.slice(1)}
                      </span>
                    </div>
                    {ngo.status === 'verified' && ngo.verifiedDate && (
                      <p className="text-xs text-gray-500 mt-1">Verified {formatDate(ngo.verifiedDate)}</p>
                    )}
                    {ngo.status === 'rejected' && ngo.rejectedDate && (
                      <p className="text-xs text-red-600 mt-1">Rejected {formatDate(ngo.rejectedDate)}</p>
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{ngo.activeCampaigns}</td>
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{formatCurrency(ngo.totalRaised)}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">{formatDate(ngo.submittedDate)}</td>
                  <td className="px-6 py-4">
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleViewDetails(ngo)}
                        className="px-3 py-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 border border-primary-300 rounded-lg transition-colors hover:bg-primary-50"
                      >
                        Details
                      </button>
                      {ngo.status === 'pending' && (
                        <>
                          <button
                            onClick={() => handleApprove(ngo)}
                            className="px-3 py-1.5 text-sm font-medium text-green-600 hover:text-green-700 border border-green-300 rounded-lg transition-colors hover:bg-green-50"
                          >
                            Approve
                          </button>
                          <button
                            onClick={() => handleReject(ngo)}
                            className="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 border border-red-300 rounded-lg transition-colors hover:bg-red-50"
                          >
                            Reject
                          </button>
                        </>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Details Modal */}
      {showDetailsModal && selectedNGO && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowDetailsModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-2xl shadow-xl max-h-[90vh] overflow-y-auto">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold text-gray-900">NGO Details</h2>
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
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Organization Name</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.name}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Category</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.category}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Registration Number</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.registrationNumber}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Status</label>
                    <div className="mt-1 flex items-center gap-2">
                      {getStatusIcon(selectedNGO.status)}
                      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[selectedNGO.status]}`}>
                        {selectedNGO.status.charAt(0).toUpperCase() + selectedNGO.status.slice(1)}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Email</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.email}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Phone</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.phone || 'Not provided'}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Website</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.website || 'Not provided'}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Location</label>
                    <p className="mt-1 text-sm text-gray-900">{selectedNGO.location}</p>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Submitted Documents</label>
                  <div className="space-y-2">
                    {selectedNGO.documents?.map((doc, index) => (
                      <div key={index} className="flex items-center gap-2 p-2 bg-green-50 rounded-lg">
                        <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                        </svg>
                        <span className="text-sm text-green-800">{doc.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}</span>
                      </div>
                    ))}
                    {selectedNGO.missingDocuments?.map((doc, index) => (
                      <div key={index} className="flex items-center gap-2 p-2 bg-red-50 rounded-lg">
                        <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                        </svg>
                        <span className="text-sm text-red-800">Missing: {doc.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {selectedNGO.rejectionReason && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">Rejection Reason</label>
                    <div className="p-3 bg-red-50 rounded-lg">
                      <p className="text-sm text-red-800">{selectedNGO.rejectionReason}</p>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Approval Modal */}
      {showApprovalModal && selectedNGO && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowApprovalModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Approve NGO</h2>
              <p className="text-gray-600 mb-6">
                Are you sure you want to approve <strong>{selectedNGO.name}</strong>? This will grant them full access to create and manage campaigns.
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
                  {loading ? 'Approving...' : 'Approve NGO'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Rejection Modal */}
      {showRejectionModal && selectedNGO && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowRejectionModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Reject NGO</h2>
              <p className="text-gray-600 mb-4">
                Please provide a reason for rejecting <strong>{selectedNGO.name}</strong>:
              </p>
              <textarea
                value={rejectionReason}
                onChange={(e) => setRejectionReason(e.target.value)}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                placeholder="Explain why this NGO verification is being rejected..."
              />
              <div className="flex justify-end gap-3 mt-6">
                <button
                  onClick={() => setShowRejectionModal(false)}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={confirmRejection}
                  disabled={loading || !rejectionReason.trim()}
                  className="px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Rejecting...' : 'Reject NGO'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}