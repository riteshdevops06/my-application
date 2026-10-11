import { useState, useRef } from 'react';
import { cn } from '../../utils/cn';
import Alert from './Alert';
import Button from './Button';

const DocumentUpload = ({
  onDocumentSelect,
  onUploadComplete,
  onUploadError,
  accept = ".pdf,.doc,.docx,.txt,.csv,.xlsx,.xls",
  multiple = true,
  maxSize = 10 * 1024 * 1024, // 10MB
  maxFiles = 10,
  uploadType = "document",
  entityType = "campaign",
  entityId,
  disabled = false,
  className,
  title = "Upload Documents",
  description = "Upload receipts, reports, or other supporting documents",
  ...props
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState(null);
  const [documents, setDocuments] = useState([]);
  
  const fileInputRef = useRef(null);
  const dropZoneRef = useRef(null);

  // Get document type icon
  const getDocumentIcon = (filename) => {
    const ext = filename.split('.').pop()?.toLowerCase();
    const iconMap = {
      pdf: '📄',
      doc: '📝',
      docx: '📝',
      txt: '📄',
      csv: '📊',
      xlsx: '📊',
      xls: '📊',
      zip: '📦',
      rar: '📦',
      default: '📎'
    };
    return iconMap[ext] || iconMap.default;
  };

  // Get document type color
  const getDocumentColor = (filename) => {
    const ext = filename.split('.').pop()?.toLowerCase();
    const colorMap = {
      pdf: 'text-red-600 bg-red-50',
      doc: 'text-blue-600 bg-blue-50',
      docx: 'text-blue-600 bg-blue-50',
      txt: 'text-gray-600 bg-gray-50',
      csv: 'text-green-600 bg-green-50',
      xlsx: 'text-green-600 bg-green-50',
      xls: 'text-green-600 bg-green-50',
      zip: 'text-purple-600 bg-purple-50',
      rar: 'text-purple-600 bg-purple-50',
      default: 'text-gray-600 bg-gray-50'
    };
    return colorMap[ext] || colorMap.default;
  };

  // Validate document
  const validateDocument = (file) => {
    if (file.size > maxSize) {
      return `Document "${file.name}" is too large. Maximum size is ${formatFileSize(maxSize)}.`;
    }
    
    const acceptedExtensions = accept.split(',').map(ext => ext.trim().toLowerCase());
    const fileExt = '.' + file.name.split('.').pop()?.toLowerCase();
    
    if (!acceptedExtensions.includes(fileExt)) {
      return `Document type "${fileExt}" is not accepted. Allowed types: ${accept}`;
    }
    
    return null;
  };

  // Handle document selection
  const handleDocuments = async (files) => {
    const fileArray = Array.from(files);
    
    if (multiple && fileArray.length > maxFiles) {
      setError(`Too many documents selected. Maximum is ${maxFiles} files.`);
      return;
    }
    
    if (!multiple && fileArray.length > 1) {
      setError('Only one document can be uploaded at a time.');
      return;
    }

    // Validate all documents
    const validationErrors = [];
    for (const file of fileArray) {
      const error = validateDocument(file);
      if (error) {
        validationErrors.push(error);
      }
    }

    if (validationErrors.length > 0) {
      setError(validationErrors.join(' '));
      return;
    }

    setError(null);
    
    // Call callback
    if (onDocumentSelect) {
      onDocumentSelect(fileArray);
    }

    // Upload if entity ID provided
    if (entityId) {
      await uploadDocuments(fileArray);
    }
  };

  // Upload documents to server
  const uploadDocuments = async (files) => {
    setIsUploading(true);

    try {
      const formData = new FormData();
      
      if (multiple) {
        files.forEach(file => {
          formData.append('files', file);
        });
        formData.append('type', uploadType);
        formData.append('entity_type', entityType);
        formData.append('entity_id', entityId);

        const response = await fetch('/api/v1/files/upload/multiple', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
          body: formData,
        });

        if (!response.ok) {
          throw new Error(`Upload failed: ${response.statusText}`);
        }

        const result = await response.json();
        setDocuments(prev => [...prev, ...result.files]);
        
        if (onUploadComplete) {
          onUploadComplete(result.files);
        }
      } else {
        formData.append('file', files[0]);
        formData.append('type', uploadType);
        formData.append('entity_type', entityType);
        formData.append('entity_id', entityId);

        const response = await fetch('/api/v1/files/upload', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
          body: formData,
        });

        if (!response.ok) {
          throw new Error(`Upload failed: ${response.statusText}`);
        }

        const result = await response.json();
        setDocuments([result]);
        
        if (onUploadComplete) {
          onUploadComplete([result]);
        }
      }
    } catch (error) {
      setError(error.message);
      if (onUploadError) {
        onUploadError(error);
      }
    } finally {
      setIsUploading(false);
    }
  };

  // Remove document
  const removeDocument = async (documentId) => {
    try {
      const docToRemove = documents.find(doc => doc.id === documentId);
      if (!docToRemove) return;

      const response = await fetch(docToRemove.url, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (response.ok) {
        setDocuments(prev => prev.filter(doc => doc.id !== documentId));
      }
    } catch (error) {
      console.error('Error removing document:', error);
    }
  };

  // Drag and drop handlers
  const handleDragEnter = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.currentTarget === dropZoneRef.current) {
      setIsDragging(false);
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    if (disabled) return;
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleDocuments(files);
    }
  };

  // Click handler
  const handleClick = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  // Input change handler
  const handleInputChange = (e) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleDocuments(files);
    }
  };

  // Format file size
  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const dropZoneClasses = cn(
    'border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-all duration-200',
    isDragging
      ? 'border-blue-400 bg-blue-50 scale-105'
      : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50',
    disabled && 'opacity-50 cursor-not-allowed',
    className
  );

  return (
    <div {...props}>
      {/* Header */}
      <div className="mb-4">
        <h3 className="text-lg font-medium text-gray-900">{title}</h3>
        <p className="text-sm text-gray-600">{description}</p>
      </div>

      {/* Drop Zone */}
      <div
        ref={dropZoneRef}
        className={dropZoneClasses}
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        onClick={handleClick}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept={accept}
          multiple={multiple}
          onChange={handleInputChange}
          className="hidden"
          disabled={disabled}
        />
        
        <div className="space-y-4">
          <div className="text-5xl">📎</div>
          <div>
            <p className="text-lg text-gray-700 font-medium mb-2">
              {isUploading ? 'Uploading documents...' : 'Drop documents here or click to browse'}
            </p>
            <p className="text-sm text-gray-500">
              Supported formats: {accept.replace(/\./g, '').toUpperCase()}
            </p>
            <p className="text-sm text-gray-500">
              Max file size: {formatFileSize(maxSize)}
              {multiple && ` • Max ${maxFiles} files`}
            </p>
          </div>
          
          {!isUploading && (
            <Button variant="outline" size="sm">
              Choose Files
            </Button>
          )}
        </div>
        
        {isUploading && (
          <div className="mt-4">
            <div className="animate-spin rounded-full h-8 w-8 border-2 border-blue-600 border-t-transparent mx-auto"></div>
            <p className="text-sm text-gray-600 mt-2">Processing upload...</p>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <Alert type="error" message={error} className="mt-4" />
      )}

      {/* Uploaded Documents List */}
      {documents.length > 0 && (
        <div className="mt-6">
          <h4 className="text-md font-medium text-gray-900 mb-4">Uploaded Documents</h4>
          <div className="space-y-3">
            {documents.map((doc) => (
              <div
                key={doc.id}
                className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
              >
                <div className="flex items-center space-x-4">
                  <div className={cn(
                    'flex-shrink-0 w-12 h-12 rounded-lg flex items-center justify-center text-xl',
                    getDocumentColor(doc.filename)
                  )}>
                    {getDocumentIcon(doc.filename)}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h5 className="text-sm font-medium text-gray-900 truncate">
                      {doc.filename}
                    </h5>
                    <p className="text-sm text-gray-500">
                      {formatFileSize(doc.size)}
                    </p>
                    <p className="text-xs text-gray-400">
                      Uploaded {new Date().toLocaleDateString()}
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  <a
                    href={doc.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center px-3 py-1.5 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 transition-colors"
                  >
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    View
                  </a>
                  <button
                    onClick={() => removeDocument(doc.id)}
                    className="inline-flex items-center px-3 py-1.5 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50 transition-colors"
                  >
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default DocumentUpload;