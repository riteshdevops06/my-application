import { useState, useRef, forwardRef } from 'react';
import { cn } from '../../utils/cn';
import Button from './Button';
import Alert from './Alert';

const FileUpload = forwardRef(({
  onFileSelect,
  onUploadComplete,
  onUploadError,
  accept = "*/*",
  multiple = false,
  maxSize = 10 * 1024 * 1024, // 10MB
  maxFiles = 5,
  uploadType = "document",
  entityType = "campaign",
  entityId,
  disabled = false,
  className,
  children,
  ...props
}, ref) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState({});
  const [error, setError] = useState(null);
  const [uploadedFiles, setUploadedFiles] = useState([]);
  
  const fileInputRef = useRef(null);
  const dropZoneRef = useRef(null);

  // File validation
  const validateFile = (file) => {
    if (file.size > maxSize) {
      return `File "${file.name}" is too large. Maximum size is ${formatFileSize(maxSize)}.`;
    }
    
    const acceptedTypes = accept.split(',').map(type => type.trim());
    if (accept !== "*/*" && !acceptedTypes.some(type => {
      if (type.startsWith('.')) {
        return file.name.toLowerCase().endsWith(type.toLowerCase());
      }
      return file.type.includes(type.replace('*', ''));
    })) {
      return `File "${file.name}" type is not accepted.`;
    }
    
    return null;
  };

  // Handle file selection
  const handleFiles = async (files) => {
    const fileArray = Array.from(files);
    
    if (multiple && fileArray.length > maxFiles) {
      setError(`Too many files selected. Maximum is ${maxFiles} files.`);
      return;
    }
    
    if (!multiple && fileArray.length > 1) {
      setError('Only one file can be uploaded at a time.');
      return;
    }

    // Validate all files
    const validationErrors = [];
    for (const file of fileArray) {
      const error = validateFile(file);
      if (error) {
        validationErrors.push(error);
      }
    }

    if (validationErrors.length > 0) {
      setError(validationErrors.join(' '));
      return;
    }

    setError(null);
    
    // Call onFileSelect callback
    if (onFileSelect) {
      onFileSelect(fileArray);
    }

    // Start upload
    if (entityId) {
      await uploadFiles(fileArray);
    }
  };

  // Upload files to server
  const uploadFiles = async (files) => {
    setIsUploading(true);
    setUploadProgress({});

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
        setUploadedFiles(prev => [...prev, ...result.files]);
        
        if (onUploadComplete) {
          onUploadComplete(result.files);
        }
      } else {
        // Single file upload
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
        setUploadedFiles(prev => [...prev, result]);
        
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

  // Remove uploaded file
  const removeFile = async (fileId) => {
    try {
      const fileToRemove = uploadedFiles.find(f => f.id === fileId);
      if (!fileToRemove) return;

      const response = await fetch(fileToRemove.url, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (response.ok) {
        setUploadedFiles(prev => prev.filter(f => f.id !== fileId));
      }
    } catch (error) {
      console.error('Error removing file:', error);
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
      handleFiles(files);
    }
  };

  // Click handler for file input
  const handleClick = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  // File input change handler
  const handleInputChange = (e) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFiles(files);
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

  // Get file type icon
  const getFileIcon = (filename) => {
    const ext = filename.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf':
        return '📄';
      case 'doc':
      case 'docx':
        return '📝';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
      case 'webp':
        return '🖼️';
      case 'csv':
      case 'xlsx':
      case 'xls':
        return '📊';
      default:
        return '📎';
    }
  };

  const dropZoneClasses = cn(
    'border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors',
    isDragging
      ? 'border-blue-400 bg-blue-50'
      : 'border-gray-300 hover:border-gray-400',
    disabled && 'opacity-50 cursor-not-allowed',
    className
  );

  return (
    <div ref={ref} {...props}>
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
        
        {children || (
          <div>
            <div className="text-4xl mb-2">📁</div>
            <p className="text-gray-600 mb-2">
              {isUploading ? 'Uploading...' : 'Drag and drop files here, or click to browse'}
            </p>
            <p className="text-sm text-gray-500">
              {accept !== "*/*" && `Accepted types: ${accept}`}
              {maxSize && ` • Max size: ${formatFileSize(maxSize)}`}
              {multiple && ` • Max ${maxFiles} files`}
            </p>
          </div>
        )}
        
        {isUploading && (
          <div className="mt-4">
            <div className="animate-spin rounded-full h-6 w-6 border-2 border-blue-600 border-t-transparent mx-auto"></div>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <Alert type="error" message={error} className="mt-4" />
      )}

      {/* Uploaded Files List */}
      {uploadedFiles.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium text-gray-900 mb-2">Uploaded Files</h4>
          <div className="space-y-2">
            {uploadedFiles.map((file) => (
              <div
                key={file.id}
                className="flex items-center justify-between p-2 bg-gray-50 rounded-lg"
              >
                <div className="flex items-center space-x-2">
                  <span className="text-lg">{getFileIcon(file.filename)}</span>
                  <div>
                    <p className="text-sm font-medium text-gray-900">{file.filename}</p>
                    <p className="text-xs text-gray-500">{formatFileSize(file.size)}</p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <a
                    href={file.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 hover:text-blue-800 text-sm"
                  >
                    View
                  </a>
                  <button
                    onClick={() => removeFile(file.id)}
                    className="text-red-600 hover:text-red-800 text-sm"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
});

FileUpload.displayName = 'FileUpload';

export default FileUpload;