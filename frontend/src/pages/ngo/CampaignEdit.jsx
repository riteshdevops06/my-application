import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const categories = [
  'Healthcare', 'Education', 'Animal Welfare', 'Community Development',
  'Emergency Relief', 'Environment', 'Human Rights', 'Arts & Culture',
]

export default function NGOCampaignEdit() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [errors, setErrors] = useState({})

  const [formData, setFormData] = useState({
    title: '',
    category: '',
    shortDescription: '',
    description: '',
    goal: '',
    startDate: '',
    endDate: '',
    location: '',
    status: 'draft',
  })

  // Mock data - replace with API call
  useEffect(() => {
    const mockData = {
      title: 'Blindness Treatment & Eye Surgery',
      category: 'Healthcare',
      shortDescription: 'Funding cataract surgeries and eye care in East Africa',
      description: 'Millions of people in developing countries suffer from preventable blindness. This campaign funds cataract surgeries, glaucoma treatments, and provides corrective lenses to those in need.',
      goal: '100000',
      startDate: '2024-01-15',
      endDate: '2024-12-31',
      location: 'Kenya, Uganda, Tanzania',
      status: 'active',
    }
    setFormData(mockData)
    setLoading(false)
  }, [id])

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) setErrors((prev) => ({ ...prev, [field]: null }))
  }

  const validate = () => {
    const newErrors = {}
    if (!formData.title.trim()) newErrors.title = 'Title is required'
    if (!formData.category) newErrors.category = 'Category is required'
    if (!formData.shortDescription.trim()) newErrors.shortDescription = 'Short description is required'
    if (!formData.description.trim()) newErrors.description = 'Full description is required'
    if (!formData.goal || Number(formData.goal) < 1) newErrors.goal = 'Goal amount is required'
    if (!formData.startDate) newErrors.startDate = 'Start date is required'
    if (!formData.endDate) newErrors.endDate = 'End date is required'
    else if (new Date(formData.endDate) <= new Date(formData.startDate)) newErrors.endDate = 'End date must be after start date'
    if (!formData.location.trim()) newErrors.location = 'Location is required'
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!validate()) return

    setSaving(true)
    try {
      // TODO: Call API to update campaign
      await new Promise((resolve) => setTimeout(resolve, 1000))
      notifySuccess('Campaign Updated', 'Changes have been saved successfully!')
    } catch (err) {
      notifyError('Update Failed', err.message || 'Failed to update campaign')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Edit Campaign</h1>
          <p className="mt-1 text-gray-600">Update your campaign details</p>
        </div>
        <div className="flex gap-2">
          <span className={`px-3 py-1 text-sm font-medium rounded-full ${
            formData.status === 'active' ? 'bg-green-100 text-green-800' :
            formData.status === 'draft' ? 'bg-gray-100 text-gray-800' :
            'bg-blue-100 text-blue-800'
          }`}>
            {formData.status}
          </span>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="bg-white rounded-xl border border-gray-200 p-6 space-y-6">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-gray-700">Campaign Title *</label>
          <input
            id="title"
            type="text"
            value={formData.title}
            onChange={(e) => handleChange('title', e.target.value)}
            className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.title ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
          />
          {errors.title && <p className="mt-1 text-sm text-red-600">{errors.title}</p>}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <label htmlFor="category" className="block text-sm font-medium text-gray-700">Category *</label>
            <select
              id="category"
              value={formData.category}
              onChange={(e) => handleChange('category', e.target.value)}
              className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.category ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
            >
              {categories.map((cat) => <option key={cat} value={cat}>{cat}</option>)}
            </select>
            {errors.category && <p className="mt-1 text-sm text-red-600">{errors.category}</p>}
          </div>

          <div>
            <label htmlFor="status" className="block text-sm font-medium text-gray-700">Status</label>
            <select
              id="status"
              value={formData.status}
              onChange={(e) => handleChange('status', e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="completed">Completed</option>
            </select>
          </div>
        </div>

        <div>
          <label htmlFor="shortDescription" className="block text-sm font-medium text-gray-700">Short Description *</label>
          <textarea
            id="shortDescription"
            value={formData.shortDescription}
            onChange={(e) => handleChange('shortDescription', e.target.value)}
            rows={3}
            maxLength={200}
            className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.shortDescription ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
          />
          {errors.shortDescription && <p className="mt-1 text-sm text-red-600">{errors.shortDescription}</p>}
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-gray-700">Full Description *</label>
          <textarea
            id="description"
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
            rows={6}
            className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.description ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
          />
          {errors.description && <p className="mt-1 text-sm text-red-600">{errors.description}</p>}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <label htmlFor="goal" className="block text-sm font-medium text-gray-700">Fundraising Goal ($) *</label>
            <input
              id="goal"
              type="number"
              min="1"
              value={formData.goal}
              onChange={(e) => handleChange('goal', e.target.value)}
              className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.goal ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
            />
            {errors.goal && <p className="mt-1 text-sm text-red-600">{errors.goal}</p>}
          </div>

          <div>
            <label htmlFor="location" className="block text-sm font-medium text-gray-700">Location *</label>
            <input
              id="location"
              type="text"
              value={formData.location}
              onChange={(e) => handleChange('location', e.target.value)}
              className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.location ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
            />
            {errors.location && <p className="mt-1 text-sm text-red-600">{errors.location}</p>}
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div>
            <label htmlFor="startDate" className="block text-sm font-medium text-gray-700">Start Date *</label>
            <input
              id="startDate"
              type="date"
              value={formData.startDate}
              onChange={(e) => handleChange('startDate', e.target.value)}
              className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.startDate ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
            />
            {errors.startDate && <p className="mt-1 text-sm text-red-600">{errors.startDate}</p>}
          </div>

          <div>
            <label htmlFor="endDate" className="block text-sm font-medium text-gray-700">End Date *</label>
            <input
              id="endDate"
              type="date"
              value={formData.endDate}
              onChange={(e) => handleChange('endDate', e.target.value)}
              className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.endDate ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
            />
            {errors.endDate && <p className="mt-1 text-sm text-red-600">{errors.endDate}</p>}
          </div>
        </div>

        <div className="pt-4 border-t border-gray-200 flex justify-end gap-3">
          <button
            type="button"
            onClick={() => navigate('/ngo/campaigns')}
            className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={saving}
            className="px-6 py-2 bg-primary-600 text-white font-medium rounded-lg hover:bg-primary-700 transition-colors disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>
    </div>
  )
}