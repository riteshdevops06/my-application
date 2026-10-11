import { useState, useEffect } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useNotifications } from '../../context/NotificationContext'

const ROLE_OPTIONS = [
  { value: 'donor', label: 'Donor', description: 'I want to donate to campaigns and track impact' },
  { value: 'ngo_admin', label: 'NGO Administrator', description: 'I represent an NGO and want to manage campaigns' },
]

export default function Register() {
  const [searchParams] = useSearchParams()
  const [step, setStep] = useState(1) // 1: role selection, 2: form
  const [role, setRole] = useState(searchParams.get('role') || 'donor')
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    // NGO-specific fields
    ngoName: '',
    ngoRegistrationNumber: '',
    ngoAddress: '',
    ngoPhone: '',
    ngoWebsite: '',
    ngoDescription: '',
  })
  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState(false)

  const { register } = useAuth()
  const { error: notifyError, success: notifySuccess } = useNotifications()
  const navigate = useNavigate()

  // Auto-select role from URL param
  useEffect(() => {
    if (searchParams.get('role')) {
      setRole(searchParams.get('role'))
      setStep(2)
    }
  }, [searchParams])

  const validateStep1 = () => {
    if (!role) {
      setErrors({ role: 'Please select an account type' })
      return false
    }
    setErrors({})
    return true
  }

  const validateStep2 = () => {
    const newErrors = {}
    if (!formData.name.trim()) newErrors.name = 'Full name is required'
    if (!formData.email.trim()) newErrors.email = 'Email is required'
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = 'Invalid email format'
    if (!formData.password) newErrors.password = 'Password is required'
    else if (formData.password.length < 8) newErrors.password = 'Password must be at least 8 characters'
    if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = 'Passwords do not match'

    if (role === 'ngo_admin') {
      if (!formData.ngoName.trim()) newErrors.ngoName = 'NGO name is required'
      if (!formData.ngoRegistrationNumber.trim()) newErrors.ngoRegistrationNumber = 'Registration number is required'
      if (!formData.ngoAddress.trim()) newErrors.ngoAddress = 'Address is required'
      if (!formData.ngoPhone.trim()) newErrors.ngoPhone = 'Phone number is required'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
    // Clear error when user types
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: null }))
    }
  }

  const handleRoleSelect = (selectedRole) => {
    setRole(selectedRole)
    setStep(2)
  }

  const handleBack = () => {
    setStep(1)
    setErrors({})
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!validateStep2()) return

    setLoading(true)
    try {
      const userData = {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        role,
        // NGO-specific data
        ...(role === 'ngo_admin' && {
          ngoName: formData.ngoName,
          ngoRegistrationNumber: formData.ngoRegistrationNumber,
          ngoAddress: formData.ngoAddress,
          ngoPhone: formData.ngoPhone,
          ngoWebsite: formData.ngoWebsite,
          ngoDescription: formData.ngoDescription,
        }),
      }
      await register(userData)
      notifySuccess('Registration Successful', 'Welcome to NGO Donation Transparency Tracker!')
    } catch (err) {
      notifyError('Registration Failed', err.message || 'An error occurred during registration')
    } finally {
      setLoading(false)
    }
  }

  // Step 1: Role Selection
  if (step === 1) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-2xl w-full space-y-8">
          <div className="text-center">
            <Link to="/campaigns" className="inline-flex">
              <div className="w-12 h-12 bg-primary-600 rounded-xl flex items-center justify-center mx-auto">
                <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.154a4 4 0 110 5.692M15.211 4.154a4 4 0 110 5.692m6.656 10.879a4 4 0 11-5.692 0m5.692 0a4 4 0 110-5.692m-5.692 0a4 4 0 110 5.692m0-11.384a4 4 0 115.692 0" />
                </svg>
              </div>
            </Link>
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">Create your account</h2>
            <p className="mt-2 text-sm text-gray-600">Choose the account type that best describes you</p>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {ROLE_OPTIONS.map((option) => (
              <button
                key={option.value}
                onClick={() => handleRoleSelect(option.value)}
                className={`relative p-6 border-2 rounded-xl text-left transition-all ${
                  role === option.value
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                }`}
              >
                <div className="flex items-start">
                  <div
                    className={`h-5 w-5 rounded border-2 flex items-center justify-center mr-4 flex-shrink-0 ${
                      role === option.value
                        ? 'border-primary-500 bg-primary-500'
                        : 'border-gray-300'
                    }`}
                  >
                    {role === option.value && (
                      <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    )}
                  </div>
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">{option.label}</h3>
                    <p className="mt-1 text-sm text-gray-600">{option.description}</p>
                  </div>
                </div>
              </button>
            ))}
          </div>

          <p className="text-center text-sm text-gray-600">
            Already have an account?{' '}
            <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    )
  }

  // Step 2: Registration Form
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <button
            onClick={handleBack}
            className="text-gray-500 hover:text-gray-700 text-sm font-medium flex items-center"
          >
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Back
          </button>
          <div className="mt-4 text-center">
            <h2 className="text-3xl font-extrabold text-gray-900">
              {role === 'donor' ? 'Donor Registration' : 'NGO Registration'}
            </h2>
            <p className="mt-2 text-sm text-gray-600">
              {role === 'donor'
                ? 'Create your donor account to start supporting causes'
                : 'Register your NGO to manage campaigns and track donations'}
            </p>
          </div>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          {/* Common fields */}
          <div>
            <label htmlFor="name" className="sr-only">Full name</label>
            <input
              id="name"
              name="name"
              type="text"
              autoComplete="name"
              required
              value={formData.name}
              onChange={handleChange}
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.name ? 'border-red-500' : 'border-gray-300'
              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
              placeholder="Full name"
            />
            {errors.name && <p className="mt-1 text-sm text-red-600">{errors.name}</p>}
          </div>

          <div>
            <label htmlFor="email" className="sr-only">Email address</label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              required
              value={formData.email}
              onChange={handleChange}
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.email ? 'border-red-500' : 'border-gray-300'
              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
              placeholder="Email address"
            />
            {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email}</p>}
          </div>

          <div>
            <label htmlFor="password" className="sr-only">Password</label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="new-password"
              required
              value={formData.password}
              onChange={handleChange}
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.password ? 'border-red-500' : 'border-gray-300'
              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
              placeholder="Password (min 8 characters)"
            />
            {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password}</p>}
          </div>

          <div>
            <label htmlFor="confirmPassword" className="sr-only">Confirm password</label>
            <input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              autoComplete="new-password"
              required
              value={formData.confirmPassword}
              onChange={handleChange}
              className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                errors.confirmPassword ? 'border-red-500' : 'border-gray-300'
              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
              placeholder="Confirm password"
            />
            {errors.confirmPassword && <p className="mt-1 text-sm text-red-600">{errors.confirmPassword}</p>}
          </div>

          {/* NGO-specific fields */}
          {role === 'ngo_admin' && (
            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">NGO Information</h3>
              <div>
                <label htmlFor="ngoName" className="sr-only">NGO name</label>
                <input
                  id="ngoName"
                  name="ngoName"
                  type="text"
                  required
                  value={formData.ngoName}
                  onChange={handleChange}
                  className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                    errors.ngoName ? 'border-red-500' : 'border-gray-300'
                  } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
                  placeholder="NGO legal name"
                />
                {errors.ngoName && <p className="mt-1 text-sm text-red-600">{errors.ngoName}</p>}
              </div>

              <div>
                <label htmlFor="ngoRegistrationNumber" className="sr-only">Registration number</label>
                <input
                  id="ngoRegistrationNumber"
                  name="ngoRegistrationNumber"
                  type="text"
                  required
                  value={formData.ngoRegistrationNumber}
                  onChange={handleChange}
                  className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                    errors.ngoRegistrationNumber ? 'border-red-500' : 'border-gray-300'
                  } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
                  placeholder="Registration number"
                />
                {errors.ngoRegistrationNumber && <p className="mt-1 text-sm text-red-600">{errors.ngoRegistrationNumber}</p>}
              </div>

              <div>
                <label htmlFor="ngoAddress" className="sr-only">Address</label>
                <textarea
                  id="ngoAddress"
                  name="ngoAddress"
                  required
                  value={formData.ngoAddress}
                  onChange={handleChange}
                  rows={3}
                  className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                    errors.ngoAddress ? 'border-red-500' : 'border-gray-300'
                  } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
                  placeholder="Registered address"
                />
                {errors.ngoAddress && <p className="mt-1 text-sm text-red-600">{errors.ngoAddress}</p>}
              </div>

              <div>
                <label htmlFor="ngoPhone" className="sr-only">Phone number</label>
                <input
                  id="ngoPhone"
                  name="ngoPhone"
                  type="tel"
                  required
                  value={formData.ngoPhone}
                  onChange={handleChange}
                  className={`appearance-none rounded-md relative block w-full px-3 py-2 border ${
                    errors.ngoPhone ? 'border-red-500' : 'border-gray-300'
                  } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm`}
                  placeholder="Phone number"
                />
                {errors.ngoPhone && <p className="mt-1 text-sm text-red-600">{errors.ngoPhone}</p>}
              </div>

              <div>
                <label htmlFor="ngoWebsite" className="sr-only">Website (optional)</label>
                <input
                  id="ngoWebsite"
                  name="ngoWebsite"
                  type="url"
                  value={formData.ngoWebsite}
                  onChange={handleChange}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                  placeholder="Website URL (optional)"
                />
              </div>

              <div>
                <label htmlFor="ngoDescription" className="sr-only">Description (optional)</label>
                <textarea
                  id="ngoDescription"
                  name="ngoDescription"
                  value={formData.ngoDescription}
                  onChange={handleChange}
                  rows={3}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
                  placeholder="Brief description of your NGO (optional)"
                />
              </div>
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={loading}
              className="w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? (
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
              ) : (
                'Create account'
              )}
            </button>
          </div>
        </form>

        <p className="text-center text-sm text-gray-600">
          Already have an account?{' '}
          <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}