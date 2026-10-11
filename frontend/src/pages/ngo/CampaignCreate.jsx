import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationContext'

const categories = [
  'Healthcare',
  'Education',
  'Animal Welfare',
  'Community Development',
  'Emergency Relief',
  'Environment',
  'Human Rights',
  'Arts & Culture',
]

export default function NGOCampaignCreate() {
  const navigate = useNavigate()
  const { success: notifySuccess, error: notifyError } = useNotifications()
  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState({})
  const [step, setStep] = useState(1) // 1: basic info, 2: details, 3: tiers

  const [formData, setFormData] = useState({
    // Step 1
    title: '',
    category: '',
    shortDescription: '',
    // Step 2
    description: '',
    goal: '',
    startDate: '',
    endDate: '',
    location: '',
    // Step 3
    tiers: [
      { amount: '', title: '', description: '' },
      { amount: '', title: '', description: '' },
      { amount: '', title: '', description: '' },
    ],
  })

  const validateStep = () => {
    const newErrors = {}
    if (step === 1) {
      if (!formData.title.trim()) newErrors.title = 'Title is required'
      if (!formData.category) newErrors.category = 'Category is required'
      if (!formData.shortDescription.trim()) newErrors.shortDescription = 'Short description is required'
      else if (formData.shortDescription.length > 200) newErrors.shortDescription = 'Max 200 characters'
    } else if (step === 2) {
      if (!formData.description.trim()) newErrors.description = 'Full description is required'
      if (!formData.goal || Number(formData.goal) < 1) newErrors.goal = 'Goal amount is required'
      if (!formData.startDate) newErrors.startDate = 'Start date is required'
      if (!formData.endDate) newErrors.endDate = 'End date is required'
      else if (new Date(formData.endDate) <= new Date(formData.startDate)) newErrors.endDate = 'End date must be after start date'
      if (!formData.location.trim()) newErrors.location = 'Location is required'
    } else if (step === 3) {
      const validTiers = formData.tiers.filter(t => t.amount && t.title)
      if (validTiers.length === 0) newErrors.tiers = 'At least one donation tier is required'
      validTiers.forEach((tier, index) => {
        if (!tier.amount || Number(tier.amount) < 1) newErrors[`tier_${index}_amount`] = 'Valid amount required'
        if (!tier.title.trim()) newErrors[`tier_${index}_title`] = 'Title required'
      })
    }
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) setErrors((prev) => ({ ...prev, [field]: null }))
  }

  const handleTierChange = (index, field, value) => {
    setFormData((prev) => {
      const newTiers = [...prev.tiers]
      newTiers[index] = { ...newTiers[index], [field]: value }
      return { ...prev, tiers: newTiers }
    })
    if (errors[`tier_${index}_${field}`]) {
      setErrors((prev) => ({ ...prev, [`tier_${index}_${field}`]: null }))
    }
  }

  const addTier = () => {
    setFormData((prev) => ({ ...prev, tiers: [...prev.tiers, { amount: '', title: '', description: '' }] }))
  }

  const removeTier = (index) => {
    if (formData.tiers.length <= 1) return
    setFormData((prev) => ({ ...prev, tiers: prev.tiers.filter((_, i) => i !== index) }))
  }

  const handleNext = () => {
    if (!validateStep()) return
    setStep((prev) => prev + 1)
  }

  const handleBack = () => {
    setStep((prev) => prev - 1)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!validateStep()) return

    setLoading(true)
    try {
      // TODO: Call API to create campaign
      await new Promise((resolve) => setTimeout(resolve, 1500))
      notifySuccess('Campaign Created', 'Your campaign has been created successfully!')
      navigate('/ngo/campaigns')
    } catch (err) {
      notifyError('Creation Failed', err.message || 'Failed to create campaign')
    } finally {
      setLoading(false)
    }
  }

  const steps = [
    { number: 1, label: 'Basic Info' },
    { number: 2, label: 'Details' },
    { number: 3, label: 'Donation Tiers' },
  ]

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Progress Steps */}
      <div className="flex items-center justify-between">
        {steps.map((s, index) => (
          <div key={s.number} className="flex items-center">
            <div
              className={`w-10 h-10 rounded-full flex items-center justify-center font-medium ${
                step >= s.number
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-200 text-gray-500'
              }`}
            >
              {step > s.number ? (
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg>
              ) : s.number}
            </div>
            <span className={`ml-2 text-sm font-medium ${step >= s.number ? 'text-primary-600' : 'text-gray-500'}`}>
              {s.label}
            </span>
            {index < steps.length - 1 && (
              <div
                className={`ml-2 flex-1 h-0.5 max-w-xs ${step > s.number ? 'bg-primary-600' : 'bg-gray-200'}`}
              />
            )}
          </div>
        ))}
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="bg-white rounded-xl border border-gray-200 p-6">
        <div>
          <h2 className="text-xl font-bold text-gray-900">{steps[step - 1].label}</h2>
          <p className="mt-1 text-gray-600">
            {step === 1 && 'Basic campaign information'}
            {step === 2 && 'Campaign details and timeline'}
            {step === 3 && 'Set up suggested donation amounts'}
          </p>
        </div>

        {/* Step 1: Basic Info */}
        {step === 1 && (
          <div className="mt-6 space-y-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700">Campaign Title *</label>
              <input
                id="title"
                type="text"
                value={formData.title}
                onChange={(e) => handleChange('title', e.target.value)}
                className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.title ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
                placeholder="e.g., Blindness Treatment & Eye Surgery"
                maxLength="100"
              />
              {errors.title && <p className="mt-1 text-sm text-red-600">{errors.title}</p>}
              <p className="mt-1 text-sm text-gray-500">{formData.title.length}/100 characters</p>
            </div>

            <div>
              <label htmlFor="category" className="block text-sm font-medium text-gray-700">Category *</label>
              <select
                id="category"
                value={formData.category}
                onChange={(e) => handleChange('category', e.target.value)}
                className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.category ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
              >
                <option value="">Select a category</option>
                {categories.map((cat) => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>
              {errors.category && <p className="mt-1 text-sm text-red-600">{errors.category}</p>}
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
                placeholder="Brief summary shown in campaign cards (max 200 characters)"
              />
              {errors.shortDescription && <p className="mt-1 text-sm text-red-600">{errors.shortDescription}</p>}
              <p className="mt-1 text-sm text-gray-500">{formData.shortDescription.length}/200 characters</p>
            </div>
          </div>
        )}

        {/* Step 2: Details */}
        {step === 2 && (
          <div className="mt-6 space-y-6">
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">Full Description *</label>
              <textarea
                id="description"
                value={formData.description}
                onChange={(e) => handleChange('description', e.target.value)}
                rows={6}
                className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors.description ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
                placeholder="Detailed description of your campaign, goals, and impact..."
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
                  step="1"
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
                  placeholder="e.g., Kenya, Uganda, Tanzania"
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
          </div>
        )}

        {/* Step 3: Donation Tiers */}
        {step === 3 && (
          <div className="mt-6 space-y-6">
            <div>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Donation Tiers</h3>
                <button type="button" onClick={addTier} className="text-sm text-primary-600 hover:text-primary-700 font-medium">
                  + Add Tier
                </button>
              </div>
              <p className="text-sm text-gray-600">Create suggested donation amounts to encourage giving</p>
            </div>

            <div className="space-y-4">
              {formData.tiers.map((tier, index) => (
                <div key={index} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex items-start justify-between gap-4 mb-4">
                    <h4 className="font-medium text-gray-900">Tier {index + 1}</h4>
                    {formData.tiers.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeTier(index)}
                        className="text-red-600 hover:text-red-700 text-sm"
                      >
                        Remove
                      </button>
                    )}
                  </div>
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Amount ($) *</label>
                      <input
                        type="number"
                        min="1"
                        step="1"
                        value={tier.amount}
                        onChange={(e) => handleTierChange(index, 'amount', e.target.value)}
                        className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors[`tier_${index}_amount`] ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
                      />
                      {errors[`tier_${index}_amount`] && <p className="mt-1 text-sm text-red-600">{errors[`tier_${index}_amount`]}</p>}
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Title *</label>
                      <input
                        type="text"
                        value={tier.title}
                        onChange={(e) => handleTierChange(index, 'title', e.target.value)}
                        className={`mt-1 block w-full px-3 py-2 border rounded-lg ${errors[`tier_${index}_title`] ? 'border-red-500' : 'border-gray-300'} focus:outline-none focus:ring-2 focus:ring-primary-500`}
                        placeholder="e.g., Vision Champion"
                      />
                      {errors[`tier_${index}_title`] && <p className="mt-1 text-sm text-red-600">{errors[`tier_${index}_title`]}</p>}
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Description</label>
                      <input
                        type="text"
                        value={tier.description}
                        onChange={(e) => handleTierChange(index, 'description', e.target.value)}
                        className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                        placeholder="e.g., Funds one cataract surgery"
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {errors.tiers && <p className="text-sm text-red-600">{errors.tiers}</p>}
          </div>
        )}

        {/* Navigation Buttons */}
        <div className="mt-8 flex items-center justify-between pt-6 border-t border-gray-200">
          <button
            type="button"
            onClick={handleBack}
            disabled={step === 1}
            className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Back
          </button>
          <div className="flex gap-3">
            {step < 3 && (
              <button
                type="button"
                onClick={handleNext}
                className="px-6 py-2 bg-primary-600 text-white font-medium rounded-lg hover:bg-primary-700 transition-colors"
              >
                Next
              </button>
            )}
            {step === 3 && (
              <button
                type="submit"
                disabled={loading}
                className="px-6 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50"
              >
                {loading ? 'Creating...' : 'Create Campaign'}
              </button>
            )}
          </div>
        </div>
      </form>
    </div>
  )
}