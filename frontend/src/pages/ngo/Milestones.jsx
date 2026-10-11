import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const mockMilestones = [
  {
    id: 1,
    title: 'Phase 1: Mobile Eye Clinics Setup',
    description: 'Establish mobile eye clinics in rural areas',
    targetAmount: 20000,
    raised: 20000,
    targetDate: '2024-06-30',
    status: 'completed',
    completedDate: '2024-06-15',
  },
  {
    id: 2,
    title: 'Phase 2: Cataract Surgeries',
    description: 'Perform 500 cataract surgeries across the mobile clinics',
    targetAmount: 50000,
    raised: 30000,
    targetDate: '2024-10-31',
    status: 'in_progress',
  },
]

const mockCampaign = {
  id: 1,
  title: 'Blindness Treatment & Eye Surgery',
  goal: 100000,
  raised: 50000,
}

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  overdue: 'bg-red-100 text-red-800',
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

export default function NGOMilestones() {
  const { id } = useParams()
  const { success: notifySuccess } = useNotifications()
  const [milestones, setMilestones] = useState(mockMilestones)
  const [showCreateModal, setShowCreateModal] = useState(false)

  const handleComplete = async (milestone) => {
    setMilestones(prev => prev.map(m => 
      m.id === milestone.id 
        ? { ...m, status: 'completed', completedDate: new Date().toISOString().split('T')[0] }
        : m
    ))
    notifySuccess('Milestone Completed', `${milestone.title} has been marked as completed!`)
  }

  const getProgress = (raised, target) => target > 0 ? Math.min(Math.round((raised / target) * 100), 100) : 0

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Campaign Milestones</h1>
          <p className="mt-1 text-gray-600">{mockCampaign.title}</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center px-4 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 transition-colors"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add Milestone
        </button>
      </div>

      {/* Campaign Overview */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div>
            <p className="text-sm text-gray-600">Total Goal</p>
            <p className="text-2xl font-bold text-gray-900">{formatCurrency(mockCampaign.goal)}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Total Raised</p>
            <p className="text-2xl font-bold text-primary-600">{formatCurrency(mockCampaign.raised)}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Milestones</p>
            <p className="text-2xl font-bold text-gray-900">
              {milestones.filter(m => m.status === 'completed').length} / {milestones.length}
            </p>
          </div>
        </div>
        <div className="mt-4">
          <div className="h-3 bg-gray-200 rounded-full overflow-hidden">
            <div 
              className="h-full bg-primary-600 rounded-full transition-all"
              style={{ width: `${getProgress(mockCampaign.raised, mockCampaign.goal)}%` }}
            />
          </div>
          <p className="mt-2 text-sm text-gray-600">
            {getProgress(mockCampaign.raised, mockCampaign.goal)}% of total goal achieved
          </p>
        </div>
      </div>

      {/* Milestones List */}
      <div className="space-y-4">
        {milestones.map((milestone, index) => {
          const progress = getProgress(milestone.raised, milestone.targetAmount)
          
          return (
            <div key={milestone.id} className="bg-white rounded-xl border border-gray-200 p-6">
              <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="text-sm font-medium text-gray-500">Phase {index + 1}</span>
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[milestone.status]}`}>
                      {milestone.status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                    </span>
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-1">{milestone.title}</h3>
                  <p className="text-gray-600 mb-4">{milestone.description}</p>
                  
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
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
                  </div>

                  {milestone.status !== 'completed' && (
                    <div className="mb-4">
                      <div className="flex justify-between text-sm mb-1">
                        <span>Progress</span>
                        <span>{progress}%</span>
                      </div>
                      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div 
                          className="h-full bg-primary-600 rounded-full transition-all"
                          style={{ width: `${progress}%` }}
                        />
                      </div>
                    </div>
                  )}

                  {milestone.completedDate && (
                    <p className="text-sm text-green-600">
                      ✓ Completed on {formatDate(milestone.completedDate)}
                    </p>
                  )}
                </div>
                
                <div className="flex gap-2">
                  {milestone.status !== 'completed' && (
                    <button
                      onClick={() => handleComplete(milestone)}
                      className="px-3 py-1.5 text-sm font-medium text-green-600 hover:text-green-700 border border-green-300 rounded-lg transition-colors hover:bg-green-50"
                    >
                      Mark Complete
                    </button>
                  )}
                  <button className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors hover:bg-gray-50">
                    Edit
                  </button>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Create Modal (placeholder) */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowCreateModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Create New Milestone</h2>
              <p className="text-gray-600 mb-4">Milestone creation form would go here.</p>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => setShowCreateModal(false)}
                  className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={() => setShowCreateModal(false)}
                  className="px-4 py-2 bg-primary-600 text-white rounded-lg font-medium hover:bg-primary-700 transition-colors"
                >
                  Create
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}