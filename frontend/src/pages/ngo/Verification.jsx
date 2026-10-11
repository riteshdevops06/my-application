import { useState } from 'react'
import { useNotifications } from '../../context/NotificationContext'

const mockVerificationStatus = {
  status: 'pending', // 'pending', 'approved', 'rejected', 'incomplete'
  submittedDate: '2024-02-15',
  completionPercentage: 75,
}

const verificationSteps = [
  {
    id: 1,
    title: 'Organization Documents',
    description: 'Upload official registration documents',
    completed: true,
  },
  {
    id: 2,
    title: 'Financial Information',
    description: 'Provide financial transparency documents',
    completed: true,
  },
  {
    id: 3,
    title: 'Leadership & Contact',
    description: 'Verify leadership and contact information',
    completed: false,
  },
]

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
  incomplete: 'bg-gray-100 text-gray-800',
}

export default function NGOVerification() {
  const { success: notifySuccess } = useNotifications()
  const [verificationData] = useState(mockVerificationStatus)

  const handleSubmitForReview = async () => {
    try {
      // TODO: Call API to submit for review
      await new Promise(resolve => setTimeout(resolve, 1000))
      notifySuccess('Submitted for Review', 'Your verification documents have been submitted for review')
    } catch (err) {
      console.error('Submission failed:', err)
    }
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">NGO Verification</h1>
        <p className="mt-1 text-gray-600">Complete the verification process to unlock all platform features</p>
      </div>

      {/* Status Overview */}
      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <h2 className="text-lg font-semibold text-gray-900">Verification Status</h2>
              <span className={`px-2 py-1 text-sm font-medium rounded-full ${statusColors[verificationData.status]}`}>
                {verificationData.status.charAt(0).toUpperCase() + verificationData.status.slice(1)}
              </span>
            </div>
            
            {verificationData.status === 'pending' && (
              <p className="text-sm text-gray-600">
                Submitted on {new Date(verificationData.submittedDate).toLocaleDateString()} • Under review
              </p>
            )}
            
            {verificationData.status === 'incomplete' && (
              <p className="text-sm text-gray-600">
                Complete all required documents to submit for review
              </p>
            )}
          </div>
          
          <div className="text-center">
            <div className="text-3xl font-bold text-primary-600">{verificationData.completionPercentage}%</div>
            <div className="text-sm text-gray-500">Complete</div>
          </div>
        </div>
        
        {/* Progress Bar */}
        <div className="mt-4">
          <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div 
              className="h-full bg-primary-600 rounded-full transition-all"
              style={{ width: `${verificationData.completionPercentage}%` }}
            />
          </div>
        </div>
      </div>

      {/* Verification Steps */}
      <div className="space-y-6">
        {verificationSteps.map((step) => (
          <div key={step.id} className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center font-medium ${
                step.completed ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'
              }`}>
                {step.completed ? (
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                  </svg>
                ) : step.id}
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-gray-900">{step.title}</h3>
                <p className="text-sm text-gray-600">{step.description}</p>
              </div>
              <div className="text-sm text-gray-500">
                {step.completed ? 'Completed' : 'Pending'}
              </div>
            </div>
            
            {!step.completed && (
              <div className="ml-11">
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
                  <div className="text-center">
                    <svg className="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                    <p className="mt-2 text-sm text-gray-600">
                      Upload your documents here
                    </p>
                    <button className="mt-2 px-3 py-1.5 text-sm font-medium text-primary-600 border border-primary-300 rounded-lg hover:bg-primary-50 transition-colors">
                      Choose Files
                    </button>
                  </div>
                </div>
              </div>
            )}

            {step.completed && (
              <div className="ml-11">
                <div className="p-3 bg-green-50 rounded-lg">
                  <div className="flex items-center gap-2">
                    <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <p className="text-sm text-green-800">Documents uploaded and verified</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Submit for Review */}
      {verificationData.status === 'incomplete' && (
        <div className="bg-white rounded-xl border border-gray-200 p-6 text-center">
          <svg className="mx-auto h-12 w-12 text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">Ready for Review!</h3>
          <p className="text-gray-600 mb-4">
            All required documents have been uploaded. Submit your application for admin review.
          </p>
          <button
            onClick={handleSubmitForReview}
            className="px-6 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 transition-colors"
          >
            Submit for Review
          </button>
        </div>
      )}
    </div>
  )
}