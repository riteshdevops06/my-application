import { useState } from 'react'
import { useNotifications } from '../../context/NotificationContext'

const mockMilestones = [
  {
    id: 1,
    campaignId: 1,
    campaignTitle: 'Blindness Treatment & Eye Surgery',
    ngo: 'Vision Care Foundation',
    title: 'Phase 1: Mobile Eye Clinics Setup',
    description: 'Establish mobile eye clinics in rural areas',
    targetAmount: 20000,
    raised: 20000,
    targetDate: '2024-06-30',
    status: 'completed',
    completedDate: '2024-06-15',
    approvalStatus: 'approved',
    submittedEvidence: ['clinic_photos.jpg', 'equipment_receipts.pdf', 'staff_training_certificates.pdf'],
    adminNotes: 'Excellent execution. All clinics operational as planned.',
  },
  {
    id: 2,
    campaignId: 1,
    campaignTitle: 'Blindness Treatment & Eye Surgery',
    ngo: 'Vision Care Foundation',
    title: 'Phase 2: Cataract Surgeries',
    description: 'Perform 500 cataract surgeries',
    targetAmount: 50000,
    raised: 30000,
    targetDate: '2024-10-31',
    status: 'in_progress',
    approvalStatus: 'not_required',
  },
  {
    id: 3,
    campaignId: 2,
    campaignTitle: 'Animal Vaccination Drive',
    ngo: 'Paws & Claws Rescue',
    title: 'Community Outreach Phase',
    description: 'Vaccinate 1000 animals in rural communities',
    targetAmount: 15000,
    raised: 12000,
    targetDate: '2024-08-15',
    status: 'completed',
    completedDate: '2024-08-10',
    approvalStatus: 'pending_approval',
    submittedEvidence: ['vaccination_records.pdf', 'animal_photos.jpg'],
    submittedDate: '2024-08-12',
  },
]

const approvalStatusColors = {
  pending_approval: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
  not_required: 'bg-gray-100 text-gray-600',
}

const statusColors = {
  pending: 'bg-gray-100 text-gray-800',
  in_progress: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

export default function PlatformAdminMilestones() {
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [milestones, setMilestones] = useState(mockMilestones)
  const [filter, setFilter] = useState('all')
  const [selectedMilestone, setSelectedMilestone] = useState(null)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [showApprovalModal, setShowApprovalModal] = useState(false)
  const [showRejectionModal, setShowRejectionModal] = useState(false)
  const [adminNotes, setAdminNotes] = useState('')
  const [rejectionReason, setRejectionReason] = useState('')
  const [loading, setLoading] = useState(false)
  const filteredMilestones = milestones.filter(milestone => {
    if (filter === 'all') return true
    if (filter === 'needs_approval') return milestone.approvalStatus === 'pending_approval'
    return milestone.approvalStatus === filter
  })

  const handleViewDetails = (milestone) => {
    setSelectedMilestone(milestone)
    setAdminNotes(milestone.adminNotes || '')
    setShowDetailsModal(true)
  }

  const handleApprove = (milestone) => {
    setSelectedMilestone(milestone)
    setAdminNotes(milestone.adminNotes || '')
    setShowApprovalModal(true)
  }

  const handleReject = (milestone) => {
    setSelectedMilestone(milestone)
    setShowRejectionModal(true)
  }

  const confirmApproval = async () => {
    if (!selectedMilestone) return
    setLoading(true)
    try {
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setMilestones(prev => prev.map(milestone => 
        milestone.id === selectedMilestone.id 
          ? { ...milestone, approvalStatus: 'approved', adminNotes, approvedDate: new Date().toISOString().split('T')[0] }
          : milestone
      ))
      
      notifySuccess('Milestone Approved', `${selectedMilestone.title} has been approved`)
      setShowApprovalModal(false)
      setSelectedMilestone(null)
      setAdminNotes('')
    } catch (err) {
      notifyError('Approval Failed', err.message || 'Failed to approve milestone')
    } finally {
      setLoading(false)
    }
  }

  const confirmRejection = async () => {
    if (!selectedMilestone || !rejectionReason.trim()) return
    setLoading(true)
    try {
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setMilestones(prev => prev.map(milestone => 
        milestone.id === selectedMilestone.id 
          ? { ...milestone, approvalStatus: 'rejected', rejectionReason, rejectedDate: new Date().toISOString().split('T')[0] }
          : milestone
      ))
      
      notifySuccess('Milestone Rejected', `${selectedMilestone.title} has been rejected`)
      setShowRejectionModal(false)
      setSelectedMilestone(null)
      setRejectionReason('')
    } catch (err) {
      notifyError('Rejection Failed', err.message || 'Failed to reject milestone')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Milestone Approval</h1>
          <p className="mt-1 text-gray-600">Review and approve completed milestones</p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Pending Approval</p>
          <p className="text-2xl font-bold text-yellow-600">{milestones.filter(m => m.approvalStatus === 'pending_approval').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Approved</p>
          <p className="text-2xl font-bold text-green-600">{milestones.filter(m => m.approvalStatus === 'approved').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Rejected</p>
          <p className="text-2xl font-bold text-red-600">{milestones.filter(m => m.approvalStatus === 'rejected').length}</p>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Value</p>
          <p className="text-xl font-bold text-primary-600">{formatCurrency(milestones.reduce((sum, m) => sum + m.targetAmount, 0))}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex gap-2 flex-wrap">
          {['all', 'needs_approval', 'approved', 'rejected', 'pending_approval'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${
                filter === status
                  ? 'bg-primary-100 text-primary-700'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {status === 'all' ? 'All Milestones' :
               status === 'needs_approval' ? 'Needs Approval' :
               status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ')}
            </button>
          ))}
        </div>
      </div>

      {/* Milestones List */}
      <div className="space-y-4">
        {filteredMilestones.map((milestone) => (
          <div key={milestone.id} className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
              <div className="flex-1">
                <div className="flex items-start justify-between gap-4 mb-2">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{milestone.title}</h3>
                    <p className="text-sm text-gray-600">{milestone.campaignTitle} • {milestone.ngo}</p>
                  </div>
                  <div className="flex flex-col gap-2">
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${approvalStatusColors[milestone.approvalStatus]}`}>
                      {milestone.approvalStatus === 'pending_approval' ? 'Pending Approval' :
                       milestone.approvalStatus === 'not_required' ? 'No Approval Needed' :
                       milestone.approvalStatus.charAt(0).toUpperCase() + milestone.approvalStatus.slice(1)}
                    </span>
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[milestone.status]}`}>
                      {milestone.status === 'in_progress' ? 'In Progress' : milestone.status.charAt(0).toUpperCase() + milestone.status.slice(1)}
                    </span>
                  </div>
                </div>
                
                <p className="text-gray-600 mb-4">{milestone.description}</p>
                
                <div className="grid grid-cols-1 sm:grid-cols-4 gap-4 mb-4">
                  <div>
                    <p className="text-sm text-gray-500">Target Amount</p>
                    <p className="font-medium">{formatCurrency(milestone.targetAmount)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Raised</p>
                    <p className="font-medium text-primary-600">{formatCurrency(milestone.raised)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Target Date</p>
                    <p className="font-medium">{formatDate(milestone.targetDate)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">{milestone.status === 'completed' ? 'Completed' : 'Status'}</p>
                    <p className="font-medium">
                      {milestone.completedDate ? formatDate(milestone.completedDate) : 'In Progress'}
                    </p>
                  </div>
                </div>

                {milestone.submittedEvidence && milestone.submittedEvidence.length > 0 && (
                  <div className="mb-4">
                    <p className="text-sm font-medium text-gray-700 mb-2">Submitted Evidence:</p>
                    <div className="flex flex-wrap gap-2">
                      {milestone.submittedEvidence.map((evidence, index) => (
                        <button
                          key={index}
                          className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-lg hover:bg-blue-200 transition-colors"
                        >
                          📎 {evidence}
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {milestone.adminNotes && (
                  <div className="p-3 bg-green-50 rounded-lg mb-4">
                    <p className="text-sm font-medium text-green-800">Admin Notes:</p>
                    <p className="text-sm text-green-700 mt-1">{milestone.adminNotes}</p>
                  </div>
                )}

                {milestone.rejectionReason && (
                  <div className="p-3 bg-red-50 rounded-lg mb-4">
                    <p className="text-sm font-medium text-red-800">Rejection Reason:</p>
                    <p className="text-sm text-red-700 mt-1">{milestone.rejectionReason}</p>
                  </div>
                )}
              </div>
              
              <div className="flex gap-2">
                <button
                  onClick={() => handleViewDetails(milestone)}
                  className="px-3 py-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 border border-primary-300 rounded-lg transition-colors hover:bg-primary-50"
                >
                  Details
                </button>
                {milestone.approvalStatus === 'pending_approval' && (
                  <>
                    <button
                      onClick={() => handleApprove(milestone)}
                      className="px-3 py-1.5 text-sm font-medium text-green-600 hover:text-green-700 border border-green-300 rounded-lg transition-colors hover:bg-green-50"
                    >
                      Approve
                    </button>
                    <button
                      onClick={() => handleReject(milestone)}
                      className="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 border border-red-300 rounded-lg transition-colors hover:bg-red-50"
                    >
                      Reject
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {filteredMilestones.length === 0 && (
        <div className="bg-white rounded-xl border border-gray-200 text-center py-16">
          <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No milestones found</h3>
          <p className="mt-1 text-sm text-gray-500">
            {filter === 'all' ? 'No milestones to display' : `No milestones with ${filter.replace('_', ' ')} status`}
          </p>
        </div>
      )}
    </div>
  )
}