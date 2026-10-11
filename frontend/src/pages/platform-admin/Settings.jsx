import { useState } from 'react';
import Card from '../../components/ui/Card';
import Button from '../../components/ui/Button';
import Input from '../../components/ui/Input';
import Select from '../../components/ui/Select';
import Textarea from '../../components/ui/Textarea';
import Alert from '../../components/ui/Alert';

const Settings = () => {
  const [activeTab, setActiveTab] = useState('general');
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(false);
  
  // General Settings
  const [generalSettings, setGeneralSettings] = useState({
    platformName: 'NGO Donation Tracker',
    platformDescription: 'A transparent platform for tracking NGO donations and milestones',
    contactEmail: 'admin@ngotracker.com',
    supportEmail: 'support@ngotracker.com',
    maxFileSize: '10',
    allowedFileTypes: 'pdf,doc,docx,jpg,png',
  });

  // Campaign Settings  
  const [campaignSettings, setCampaignSettings] = useState({
    minGoalAmount: '1000',
    maxGoalAmount: '1000000',
    defaultCampaignDuration: '90',
    requireApproval: true,
    autoPublish: false,
    maxImagesPerCampaign: '5',
  });

  // NGO Verification Settings
  const [verificationSettings, setVerificationSettings] = useState({
    requireTaxExempt: true,
    requireBoard: true,
    requireFinancials: true,
    autoReview: false,
    verificationExpiry: '365',
    requiredDocuments: 'registration,tax-exempt,financial-report',
  });

  // Payment Settings
  const [paymentSettings, setPaymentSettings] = useState({
    processingFee: '2.9',
    minimumDonation: '5',
    maximumDonation: '50000',
    currency: 'USD',
    payoutSchedule: 'weekly',
    requireSSL: true,
  });

  // Notification Settings
  const [notificationSettings, setNotificationSettings] = useState({
    emailNotifications: true,
    smsNotifications: false,
    pushNotifications: true,
    adminAlerts: true,
    donorReceipts: true,
    campaignUpdates: true,
  });

  const handleSave = async (settingsType) => {
    setLoading(true);
    try {
      // Here would be API call to save settings
      // await api.put(`/api/admin/settings/${settingsType}`, settingsData);
      
      // Mock API delay
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (error) {
      console.error('Error saving settings:', error);
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
    { id: 'general', label: 'General', icon: '⚙️' },
    { id: 'campaigns', label: 'Campaigns', icon: '🎯' },
    { id: 'verification', label: 'NGO Verification', icon: '✅' },
    { id: 'payments', label: 'Payments', icon: '💳' },
    { id: 'notifications', label: 'Notifications', icon: '🔔' },
  ];

  const renderGeneralSettings = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Platform Name"
          value={generalSettings.platformName}
          onChange={(e) => setGeneralSettings(prev => ({ ...prev, platformName: e.target.value }))}
        />
        <Input
          label="Contact Email"
          type="email"
          value={generalSettings.contactEmail}
          onChange={(e) => setGeneralSettings(prev => ({ ...prev, contactEmail: e.target.value }))}
        />
      </div>
      
      <Textarea
        label="Platform Description"
        value={generalSettings.platformDescription}
        onChange={(e) => setGeneralSettings(prev => ({ ...prev, platformDescription: e.target.value }))}
        rows={3}
      />
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Support Email"
          type="email"
          value={generalSettings.supportEmail}
          onChange={(e) => setGeneralSettings(prev => ({ ...prev, supportEmail: e.target.value }))}
        />
        <Input
          label="Max File Size (MB)"
          type="number"
          value={generalSettings.maxFileSize}
          onChange={(e) => setGeneralSettings(prev => ({ ...prev, maxFileSize: e.target.value }))}
        />
      </div>
      
      <Input
        label="Allowed File Types (comma separated)"
        value={generalSettings.allowedFileTypes}
        onChange={(e) => setGeneralSettings(prev => ({ ...prev, allowedFileTypes: e.target.value }))}
      />
      
      <Button onClick={() => handleSave('general')} loading={loading}>
        Save General Settings
      </Button>
    </div>
  );

  const renderCampaignSettings = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Minimum Goal Amount ($)"
          type="number"
          value={campaignSettings.minGoalAmount}
          onChange={(e) => setCampaignSettings(prev => ({ ...prev, minGoalAmount: e.target.value }))}
        />
        <Input
          label="Maximum Goal Amount ($)"
          type="number"
          value={campaignSettings.maxGoalAmount}
          onChange={(e) => setCampaignSettings(prev => ({ ...prev, maxGoalAmount: e.target.value }))}
        />
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Default Campaign Duration (days)"
          type="number"
          value={campaignSettings.defaultCampaignDuration}
          onChange={(e) => setCampaignSettings(prev => ({ ...prev, defaultCampaignDuration: e.target.value }))}
        />
        <Input
          label="Max Images per Campaign"
          type="number"
          value={campaignSettings.maxImagesPerCampaign}
          onChange={(e) => setCampaignSettings(prev => ({ ...prev, maxImagesPerCampaign: e.target.value }))}
        />
      </div>
      
      <div className="space-y-4">
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={campaignSettings.requireApproval}
            onChange={(e) => setCampaignSettings(prev => ({ ...prev, requireApproval: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Require admin approval for new campaigns</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={campaignSettings.autoPublish}
            onChange={(e) => setCampaignSettings(prev => ({ ...prev, autoPublish: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Auto-publish approved campaigns</span>
        </label>
      </div>
      
      <Button onClick={() => handleSave('campaigns')} loading={loading}>
        Save Campaign Settings
      </Button>
    </div>
  );

  const renderVerificationSettings = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Verification Expiry (days)"
          type="number"
          value={verificationSettings.verificationExpiry}
          onChange={(e) => setVerificationSettings(prev => ({ ...prev, verificationExpiry: e.target.value }))}
        />
        <Input
          label="Required Documents (comma separated)"
          value={verificationSettings.requiredDocuments}
          onChange={(e) => setVerificationSettings(prev => ({ ...prev, requiredDocuments: e.target.value }))}
        />
      </div>
      
      <div className="space-y-4">
        <h3 className="text-lg font-medium text-gray-900">Verification Requirements</h3>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={verificationSettings.requireTaxExempt}
            onChange={(e) => setVerificationSettings(prev => ({ ...prev, requireTaxExempt: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Require tax-exempt status documentation</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={verificationSettings.requireBoard}
            onChange={(e) => setVerificationSettings(prev => ({ ...prev, requireBoard: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Require board of directors information</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={verificationSettings.requireFinancials}
            onChange={(e) => setVerificationSettings(prev => ({ ...prev, requireFinancials: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Require financial statements</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={verificationSettings.autoReview}
            onChange={(e) => setVerificationSettings(prev => ({ ...prev, autoReview: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Enable automated review (experimental)</span>
        </label>
      </div>
      
      <Button onClick={() => handleSave('verification')} loading={loading}>
        Save Verification Settings
      </Button>
    </div>
  );

  const renderPaymentSettings = () => (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Processing Fee (%)"
          type="number"
          step="0.1"
          value={paymentSettings.processingFee}
          onChange={(e) => setPaymentSettings(prev => ({ ...prev, processingFee: e.target.value }))}
        />
        <Select
          label="Currency"
          value={paymentSettings.currency}
          onChange={(e) => setPaymentSettings(prev => ({ ...prev, currency: e.target.value }))}
        >
          <option value="USD">USD - US Dollar</option>
          <option value="EUR">EUR - Euro</option>
          <option value="GBP">GBP - British Pound</option>
          <option value="CAD">CAD - Canadian Dollar</option>
        </Select>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Input
          label="Minimum Donation ($)"
          type="number"
          value={paymentSettings.minimumDonation}
          onChange={(e) => setPaymentSettings(prev => ({ ...prev, minimumDonation: e.target.value }))}
        />
        <Input
          label="Maximum Donation ($)"
          type="number"
          value={paymentSettings.maximumDonation}
          onChange={(e) => setPaymentSettings(prev => ({ ...prev, maximumDonation: e.target.value }))}
        />
      </div>
      
      <Select
        label="Payout Schedule"
        value={paymentSettings.payoutSchedule}
        onChange={(e) => setPaymentSettings(prev => ({ ...prev, payoutSchedule: e.target.value }))}
      >
        <option value="daily">Daily</option>
        <option value="weekly">Weekly</option>
        <option value="bi-weekly">Bi-weekly</option>
        <option value="monthly">Monthly</option>
      </Select>
      
      <label className="flex items-center space-x-2">
        <input
          type="checkbox"
          checked={paymentSettings.requireSSL}
          onChange={(e) => setPaymentSettings(prev => ({ ...prev, requireSSL: e.target.checked }))}
          className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
        />
        <span>Require SSL for all payment transactions</span>
      </label>
      
      <Button onClick={() => handleSave('payments')} loading={loading}>
        Save Payment Settings
      </Button>
    </div>
  );

  const renderNotificationSettings = () => (
    <div className="space-y-6">
      <h3 className="text-lg font-medium text-gray-900">Notification Preferences</h3>
      
      <div className="space-y-4">
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.emailNotifications}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, emailNotifications: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Enable email notifications</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.smsNotifications}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, smsNotifications: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Enable SMS notifications</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.pushNotifications}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, pushNotifications: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Enable push notifications</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.adminAlerts}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, adminAlerts: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Send alerts to administrators</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.donorReceipts}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, donorReceipts: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Automatically send donation receipts</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={notificationSettings.campaignUpdates}
            onChange={(e) => setNotificationSettings(prev => ({ ...prev, campaignUpdates: e.target.checked }))}
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
          <span>Send campaign update notifications</span>
        </label>
      </div>
      
      <Button onClick={() => handleSave('notifications')} loading={loading}>
        Save Notification Settings
      </Button>
    </div>
  );

  const renderTabContent = () => {
    switch (activeTab) {
      case 'general':
        return renderGeneralSettings();
      case 'campaigns':
        return renderCampaignSettings();
      case 'verification':
        return renderVerificationSettings();
      case 'payments':
        return renderPaymentSettings();
      case 'notifications':
        return renderNotificationSettings();
      default:
        return renderGeneralSettings();
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Platform Settings</h1>
          <p className="mt-2 text-gray-600">
            Configure platform-wide settings and preferences.
          </p>
        </div>

        {saved && (
          <Alert
            type="success"
            message="Settings saved successfully!"
            className="mb-6"
          />
        )}

        <div className="bg-white rounded-lg shadow">
          {/* Tab Navigation */}
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-4 px-6 text-sm font-medium border-b-2 transition-colors ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <span className="mr-2">{tab.icon}</span>
                  {tab.label}
                </button>
              ))}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {renderTabContent()}
          </div>
        </div>

        {/* System Information */}
        <Card className="mt-8">
          <div className="p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">System Information</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-sm">
              <div>
                <span className="font-medium text-gray-700">Platform Version:</span>
                <p className="text-gray-600">v1.0.0</p>
              </div>
              <div>
                <span className="font-medium text-gray-700">Last Updated:</span>
                <p className="text-gray-600">October 11, 2026</p>
              </div>
              <div>
                <span className="font-medium text-gray-700">Environment:</span>
                <p className="text-gray-600">Production</p>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Settings;