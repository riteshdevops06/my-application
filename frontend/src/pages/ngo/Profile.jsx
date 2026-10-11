import { useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { useNotifications } from '../../context/NotificationContext'

export default function NGOProfile() {
  const { user, updateUser } = useAuth()
  const { success: notifySuccess, error: notifyError } = useNotifications()

  const [activeTab, setActiveTab] = useState('organization')
  const [loading, setLoading] = useState(false)

  // Organization profile data
  const [orgData, setOrgData] = useState({
    name: user?.organization?.name || 'Vision Care Foundation',
    description: user?.organization?.description || 'Dedicated to preventing blindness and providing eye care to underserved communities.',
    category: user?.organization?.category || 'Healthcare',
    website: user?.organization?.website || 'https://visioncare.org',
    phone: user?.organization?.phone || '+254-700-123456',
    email: user?.organization?.email || 'contact@visioncare.org',
    address: user?.organization?.address || '123 Nairobi Road, Nairobi, Kenya',
    registrationNumber: user?.organization?.registrationNumber || 'NGO-2020-001234',
  })

  const tabs = [
    { id: 'organization', label: 'Organization' },
    { id: 'personal', label: 'Personal Info' },
    { id: 'security', label: 'Security' },
    { id: 'notifications', label: 'Notifications' },
  ]

  const handleOrgSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    try {
      // TODO: Call API to update organization profile
      await new Promise((resolve) => setTimeout(resolve, 1000))
      notifySuccess('Profile Updated', 'Organization profile has been updated successfully!')
    } catch (err) {
      notifyError('Update Failed', err.message || 'Failed to update profile')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">NGO Profile Settings</h1>
        <p className="mt-1 text-gray-600">Manage your organization and account settings</p>
      </div>

      {/* Verification Status */}
      <div className="bg-white rounded-xl border border-gray-200 p-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-yellow-100 rounded-full flex items-center justify-center">
            <svg className="w-5 h-5 text-yellow-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="flex-1">
            <p className="text-sm font-medium text-yellow-800">NGO Verification Pending</p>
            <p className="text-sm text-yellow-700">Complete your verification to unlock all features</p>
          </div>
          <button className="px-3 py-1.5 text-sm font-medium text-yellow-700 border border-yellow-300 rounded-lg hover:bg-yellow-50 transition-colors">
            Complete Verification
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="flex gap-8" aria-label="Profile tabs">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div>
        {activeTab === 'organization' && (
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="mb-6">
              <h2 className="text-lg font-semibold text-gray-900">Organization Information</h2>
              <p className="text-sm text-gray-600">Update your NGO's public profile and details</p>
            </div>
            
            <form onSubmit={handleOrgSubmit} className="space-y-6">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                    Organization Name *
                  </label>
                  <input
                    id="name"
                    type="text"
                    value={orgData.name}
                    onChange={(e) => setOrgData(prev => ({ ...prev, name: e.target.value }))}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>

                <div>
                  <label htmlFor="category" className="block text-sm font-medium text-gray-700 mb-1">
                    Category *
                  </label>
                  <select
                    id="category"
                    value={orgData.category}
                    onChange={(e) => setOrgData(prev => ({ ...prev, category: e.target.value }))}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="Healthcare">Healthcare</option>
                    <option value="Education">Education</option>
                    <option value="Animal Welfare">Animal Welfare</option>
                    <option value="Community Development">Community Development</option>
                    <option value="Emergency Relief">Emergency Relief</option>
                    <option value="Environment">Environment</option>
                  </select>
                </div>
              </div>

              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                  Description *
                </label>
                <textarea
                  id="description"
                  value={orgData.description}
                  onChange={(e) => setOrgData(prev => ({ ...prev, description: e.target.value }))}
                  rows={4}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  placeholder="Describe your organization's mission and activities..."
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                    Organization Email *
                  </label>
                  <input
                    id="email"
                    type="email"
                    value={orgData.email}
                    onChange={(e) => setOrgData(prev => ({ ...prev, email: e.target.value }))}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>

                <div>
                  <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-1">
                    Phone Number
                  </label>
                  <input
                    id="phone"
                    type="tel"
                    value={orgData.phone}
                    onChange={(e) => setOrgData(prev => ({ ...prev, phone: e.target.value }))}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="address" className="block text-sm font-medium text-gray-700 mb-1">
                  Address
                </label>
                <textarea
                  id="address"
                  value={orgData.address}
                  onChange={(e) => setOrgData(prev => ({ ...prev, address: e.target.value }))}
                  rows={2}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>

              <div className="pt-4 border-t border-gray-200">
                <button
                  type="submit"
                  disabled={loading}
                  className="px-6 py-2 bg-primary-600 text-white font-medium rounded-lg hover:bg-primary-700 transition-colors disabled:opacity-50"
                >
                  {loading ? 'Saving...' : 'Save Organization Info'}
                </button>
              </div>
            </form>
          </div>
        )}

        {activeTab === 'personal' && (
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="text-center py-16">
              <p className="text-gray-500">Personal information settings would go here.</p>
            </div>
          </div>
        )}

        {activeTab === 'security' && (
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="text-center py-16">
              <p className="text-gray-500">Security and password settings would go here.</p>
            </div>
          </div>
        )}

        {activeTab === 'notifications' && (
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="text-center py-16">
              <p className="text-gray-500">Notification preferences would go here.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}