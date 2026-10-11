import { useState, useRef } from 'react';
import { cn } from '../../utils/cn';
import Button from './Button';
import Alert from './Alert';

const ImageUpload = ({
  onImageSelect,
  onUploadComplete,
  onUploadError,
  multiple = false,
  maxSize = 5 * 1024 * 1024, // 5MB
  maxImages = 5,
  aspectRatio = null, // e.g., "16:9", "1:1"
  uploadType = "image",
  entityType = "campaign",
  entityId,
  disabled = false,
  className,
  placeholder,
  ...props
}) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState(null);
  const [images, setImages] = useState([]);
  const [previews, setPreviews] = useState([]);
  
  const fileInputRef = useRef(null);
  const dropZoneRef = useRef(null);

  const acceptedTypes = "image/jpeg,image/jpg,image/png,image/gif,image/webp";

  // Create image preview
  const createPreview = (file) => {
    return new Promise((resolve) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        resolve({
          id: Date.now() + Math.random(),
          file,
          url: e.target.result,
          name: file.name,
          size: file.size,
        });
      };
      reader.readAsDataURL(file);
    });
  };

  // Validate image file
  const validateImage = (file) => {
    if (file.size > maxSize) {
      return `Image "${file.name}" is too large. Maximum size is ${formatFileSize(maxSize)}.`;
    }
    
    if (!file.type.startsWith('image/')) {
      return `File "${file.name}" is not a valid image.`;
    }
    
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      return `Image type "${file.type}" is not supported.`;
    }
    
    return null;
  };

  // Handle image selection
  const handleImages = async (files) => {
    const fileArray = Array.from(files);
    
    if (multiple && fileArray.length > maxImages) {
      setError(`Too many images selected. Maximum is ${maxImages} images.`);
      return;
    }
    
    if (!multiple && fileArray.length > 1) {
      setError('Only one image can be uploaded at a time.');
      return;
    }

    // Validate all images
    const validationErrors = [];
    for (const file of fileArray) {
      const error = validateImage(file);
      if (error) {
        validationErrors.push(error);
      }
    }

    if (validationErrors.length > 0) {
      setError(validationErrors.join(' '));
      return;
    }

    setError(null);
    
    // Create previews
    const newPreviews = await Promise.all(fileArray.map(createPreview));
    
    if (multiple) {
      setPreviews(prev => [...prev, ...newPreviews]);
    } else {
      setPreviews(newPreviews);
    }
    
    // Call callback
    if (onImageSelect) {
      onImageSelect(fileArray, newPreviews);
    }

    // Upload if entity ID provided
    if (entityId) {
      await uploadImages(fileArray);
    }
  };

  // Upload images to server
  const uploadImages = async (files) => {
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
        setImages(prev => [...prev, ...result.files]);
        
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
        setImages([result]);
        
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

  // Remove preview
  const removePreview = (previewId) => {
    setPreviews(prev => prev.filter(p => p.id !== previewId));
  };

  // Remove uploaded image
  const removeImage = async (imageId) => {
    try {
      const imageToRemove = images.find(img => img.id === imageId);
      if (!imageToRemove) return;

      const response = await fetch(imageToRemove.url, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (response.ok) {
        setImages(prev => prev.filter(img => img.id !== imageId));
      }
    } catch (error) {
      console.error('Error removing image:', error);
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
    
    const files = Array.from(e.dataTransfer.files).filter(file => 
      file.type.startsWith('image/')
    );
    
    if (files.length > 0) {
      handleImages(files);
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
      handleImages(files);
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
    'border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition-colors',
    isDragging
      ? 'border-blue-400 bg-blue-50'
      : 'border-gray-300 hover:border-gray-400',
    disabled && 'opacity-50 cursor-not-allowed',
    aspectRatio && 'relative',
    className
  );

  const getAspectRatioClasses = () => {
    if (!aspectRatio) return '';
    
    const ratios = {
      '1:1': 'aspect-square',
      '16:9': 'aspect-video',
      '4:3': 'aspect-[4/3]',
      '3:2': 'aspect-[3/2]',
    };
    
    return ratios[aspectRatio] || '';
  };

  return (
    <div {...props}>
      {/* Upload Area */}
      <div
        ref={dropZoneRef}
        className={cn(dropZoneClasses, getAspectRatioClasses())}
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        onClick={handleClick}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept={acceptedTypes}
          multiple={multiple}
          onChange={handleInputChange}
          className="hidden"
          disabled={disabled}
        />
        
        <div className={cn(aspectRatio && 'absolute inset-0 flex flex-col justify-center')}>
          <div className="text-4xl mb-2">🖼️</div>
          <p className="text-gray-600 mb-2">
            {isUploading ? 'Uploading images...' : 
             placeholder || 'Drag and drop images here, or click to browse'}
          </p>
          <p className="text-sm text-gray-500">
            PNG, JPG, GIF, WebP up to {formatFileSize(maxSize)}
            {multiple && ` • Max ${maxImages} images`}
          </p>
        </div>
        
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

      {/* Image Previews */}
      {previews.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium text-gray-900 mb-2">Selected Images</h4>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {previews.map((preview) => (
              <div key={preview.id} className="relative group">
                <div className="aspect-square rounded-lg overflow-hidden bg-gray-100">
                  <img
                    src={preview.url}
                    alt={preview.name}
                    className="w-full h-full object-cover"
                  />
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    removePreview(preview.id);
                  }}
                  className="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  ×
                </button>
                <p className="text-xs text-gray-500 mt-1 truncate">{preview.name}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Uploaded Images */}
      {images.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium text-gray-900 mb-2">Uploaded Images</h4>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {images.map((image) => (
              <div key={image.id} className="relative group">
                <div className="aspect-square rounded-lg overflow-hidden bg-gray-100">
                  <img
                    src={image.url}
                    alt={image.filename}
                    className="w-full h-full object-cover"
                  />
                </div>
                <div className="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center space-x-2">
                  <a
                    href={image.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="bg-white text-gray-900 px-2 py-1 rounded text-xs hover:bg-gray-100"
                  >
                    View
                  </a>
                  <button
                    onClick={() => removeImage(image.id)}
                    className="bg-red-500 text-white px-2 py-1 rounded text-xs hover:bg-red-600"
                  >
                    Delete
                  </button>
                </div>
                <p className="text-xs text-gray-500 mt-1 truncate">{image.filename}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default ImageUpload;