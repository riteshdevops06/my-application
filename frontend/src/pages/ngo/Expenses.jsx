import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const mockExpenses = [
  {
    id: 1,
    title: 'Medical Equipment Purchase',
    description: 'Surgical instruments and examination tools',
    amount: 15000,
    category: 'Equipment',
    date: '2024-02-01',
    status: 'approved',
    receipt: 'receipt-001.pdf',
    approvedBy: 'Platform Admin',
    approvedDate: '2024-02-05',
  },
  {
    id: 2,
    title: 'Medical Supplies',
    description: 'Medications, bandages, and consumables',
    amount: 8000,
    category: 'Supplies',
    date: '2024-02-15',
    status: 'approved',
    receipt: 'receipt-002.pdf',
  },
  {
    id: 3,
    title: 'Staff Salaries - February',
    description: 'Monthly salaries for medical staff',
    amount: 12000,
    category: 'Personnel',
    date: '2024-03-01',
    status: 'pending',
    receipt: 'payroll-002.pdf',
  },
]

const mockCampaign = {
  id: 1,
  title: 'Blindness Treatment & Eye Surgery',
  raised: 50000,
}

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
}

const formatCurrency = (amount) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 }).format(amount)
const formatDate = (date) => new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })

export default function NGOExpenses() {
  const { id } = useParams()
  const { success: notifySuccess } = useNotifications()
  const [expenses, setExpenses] = useState(mockExpenses)
  const [filter, setFilter] = useState('all')
  const [showCreateModal, setShowCreateModal] = useState(false)

  const filteredExpenses = expenses.filter(expense => {
    if (filter === 'all') return true
    return expense.status === filter
  })

  const totalExpenses = expenses.reduce((sum, exp) => sum + exp.amount, 0)
  const approvedExpenses = expenses.filter(exp => exp.status === 'approved').reduce((sum, exp) => sum + exp.amount, 0)
  const pendingExpenses = expenses.filter(exp => exp.status === 'pending').reduce((sum, exp) => sum + exp.amount, 0)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Campaign Expenses</h1>
          <p className="mt-1 text-gray-600">{mockCampaign.title}</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center px-4 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 transition-colors"
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Add Expense
        </button>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-4 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Raised</p>
          <p className="text-2xl font-bold text-gray-900">{formatCurrency(mockCampaign.raised)}</p>
        </div>
        
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Total Expenses</p>
          <p className="text-2xl font-bold text-red-600">{formatCurrency(totalExpenses)}</p>
        </div>
        
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Approved</p>
          <p className="text-2xl font-bold text-green-600">{formatCurrency(approvedExpenses)}</p>
        </div>
        
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <p className="text-sm text-gray-600">Pending Approval</p>
          <p className="text-2xl font-bold text-yellow-600">{formatCurrency(pendingExpenses)}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex gap-2 flex-wrap">
          {['all', 'pending', 'approved', 'rejected'].map((status) => (
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
      </div>

      {/* Expenses List */}
      <div className="space-y-4">
        {filteredExpenses.length === 0 ? (
          <div className="bg-white rounded-xl border border-gray-200 text-center py-16">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">
              {filter === 'all' ? 'No expenses yet' : `No ${filter} expenses`}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {filter === 'all' ? 'Add your first expense to track campaign spending' : `No expenses with ${filter} status found`}
            </p>
          </div>
        ) : (
          filteredExpenses.map((expense) => (
            <div key={expense.id} className="bg-white rounded-xl border border-gray-200 p-6">
              <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-lg font-semibold text-gray-900">{expense.title}</h3>
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[expense.status]}`}>
                      {expense.status.charAt(0).toUpperCase() + expense.status.slice(1)}
                    </span>
                  </div>
                  
                  <p className="text-gray-600 mb-4">{expense.description}</p>
                  
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-4">
                    <div>
                      <p className="text-sm text-gray-500">Amount</p>
                      <p className="font-semibold text-lg">{formatCurrency(expense.amount)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-500">Category</p>
                      <p className="font-medium">{expense.category}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-500">Date</p>
                      <p className="font-medium">{formatDate(expense.date)}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-500">Receipt</p>
                      <button className="text-primary-600 hover:text-primary-700 font-medium text-sm">
                        {expense.receipt}
                      </button>
                    </div>
                  </div>

                  {expense.status === 'approved' && expense.approvedBy && (
                    <div className="p-3 bg-green-50 rounded-lg">
                      <p className="text-sm text-green-800">
                        ✓ Approved by {expense.approvedBy} on {formatDate(expense.approvedDate)}
                      </p>
                    </div>
                  )}

                  {expense.status === 'pending' && (
                    <div className="p-3 bg-yellow-50 rounded-lg">
                      <p className="text-sm text-yellow-800">
                        ⏳ This expense is pending approval from the platform administrator.
                      </p>
                    </div>
                  )}
                </div>
                
                <div className="flex gap-2">
                  {expense.status === 'pending' && (
                    <button className="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-900 border border-gray-300 rounded-lg transition-colors hover:bg-gray-50">
                      Edit
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Create Modal (placeholder) */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 overflow-y-auto" role="dialog" aria-modal="true">
          <div className="flex min-h-full items-center justify-center p-4">
            <div className="fixed inset-0 bg-black/50" onClick={() => setShowCreateModal(false)} />
            <div className="relative bg-white rounded-xl p-6 w-full max-w-md shadow-xl">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Add New Expense</h2>
              <p className="text-gray-600 mb-4">Expense creation form would go here.</p>
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
                  Submit
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}