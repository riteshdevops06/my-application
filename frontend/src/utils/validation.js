// Validation utilities for frontend forms
import React from 'react';

/**
 * Validation rule functions
 */
export const validationRules = {
  // Required field validation
  required: (value) => {
    if (value === null || value === undefined || value === '') {
      return 'This field is required';
    }
    return null;
  },

  // Email validation
  email: (value) => {
    if (!value) return null;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(value) ? null : 'Please enter a valid email address';
  },

  // Minimum length validation
  minLength: (min) => (value) => {
    if (!value) return null;
    return value.length >= min ? null : `Must be at least ${min} characters`;
  },

  // Maximum length validation
  maxLength: (max) => (value) => {
    if (!value) return null;
    return value.length <= max ? null : `Must be no more than ${max} characters`;
  },

  // Minimum value validation
  min: (min) => (value) => {
    if (value === null || value === undefined || value === '') return null;
    const num = parseFloat(value);
    return num >= min ? null : `Must be at least ${min}`;
  },

  // Maximum value validation
  max: (max) => (value) => {
    if (value === null || value === undefined || value === '') return null;
    const num = parseFloat(value);
    return num <= max ? null : `Must be no more than ${max}`;
  },

  // Pattern matching validation
  pattern: (regex, message = 'Invalid format') => (value) => {
    if (!value) return null;
    return regex.test(value) ? null : message;
  },

  // Numeric validation
  numeric: (value) => {
    if (!value) return null;
    return !isNaN(value) && !isNaN(parseFloat(value)) ? null : 'Must be a number';
  },

  // Integer validation
  integer: (value) => {
    if (!value) return null;
    return Number.isInteger(Number(value)) ? null : 'Must be a whole number';
  },

  // URL validation
  url: (value) => {
    if (!value) return null;
    try {
      new URL(value);
      return null;
    } catch {
      return 'Must be a valid URL';
    }
  },

  // Phone number validation
  phone: (value) => {
    if (!value) return null;
    const phoneRegex = /^\+?[\d\s\-\(\)]{10,}$/;
    return phoneRegex.test(value) ? null : 'Please enter a valid phone number';
  },

  // Date validation
  date: (value) => {
    if (!value) return null;
    const date = new Date(value);
    return !isNaN(date.getTime()) ? null : 'Please enter a valid date';
  },

  // Future date validation
  futureDate: (value) => {
    if (!value) return null;
    const date = new Date(value);
    const now = new Date();
    return date > now ? null : 'Date must be in the future';
  },

  // Past date validation
  pastDate: (value) => {
    if (!value) return null;
    const date = new Date(value);
    const now = new Date();
    return date < now ? null : 'Date must be in the past';
  },

  // Password strength validation
  password: (value) => {
    if (!value) return null;
    
    const errors = [];
    if (value.length < 8) errors.push('at least 8 characters');
    if (!/[A-Z]/.test(value)) errors.push('one uppercase letter');
    if (!/[a-z]/.test(value)) errors.push('one lowercase letter');
    if (!/\d/.test(value)) errors.push('one number');
    if (!/[!@#$%^&*(),.?":{}|<>]/.test(value)) errors.push('one special character');
    
    return errors.length === 0 ? null : `Password must contain ${errors.join(', ')}`;
  },

  // Confirm password validation
  confirmPassword: (originalPassword) => (value) => {
    if (!value) return null;
    return value === originalPassword ? null : 'Passwords do not match';
  },

  // File validation
  file: (options = {}) => (files) => {
    if (!files || files.length === 0) return null;
    
    const { 
      maxSize = 10 * 1024 * 1024, // 10MB
      allowedTypes = [],
      maxFiles = 1 
    } = options;
    
    if (files.length > maxFiles) {
      return `Maximum ${maxFiles} file(s) allowed`;
    }
    
    for (const file of files) {
      if (file.size > maxSize) {
        return `File size must be less than ${formatFileSize(maxSize)}`;
      }
      
      if (allowedTypes.length > 0) {
        const fileType = file.type || '';
        const fileName = file.name || '';
        const extension = fileName.split('.').pop()?.toLowerCase();
        
        const isAllowed = allowedTypes.some(type => {
          if (type.startsWith('.')) {
            return extension === type.substring(1);
          }
          return fileType.includes(type);
        });
        
        if (!isAllowed) {
          return `File type not allowed. Allowed types: ${allowedTypes.join(', ')}`;
        }
      }
    }
    
    return null;
  },

  // Custom validation
  custom: (validator) => validator,
};

/**
 * Validation schema builder
 */
export class ValidationSchema {
  constructor() {
    this.rules = {};
  }

  field(name, ...validators) {
    this.rules[name] = validators;
    return this;
  }

  validate(data) {
    const errors = {};
    
    Object.entries(this.rules).forEach(([fieldName, validators]) => {
      const value = data[fieldName];
      
      for (const validator of validators) {
        const error = typeof validator === 'function' 
          ? validator(value)
          : validationRules[validator]?.(value);
          
        if (error) {
          errors[fieldName] = error;
          break; // Stop at first error for this field
        }
      }
    });
    
    return {
      isValid: Object.keys(errors).length === 0,
      errors
    };
  }
}

/**
 * Form validator class
 */
export class FormValidator {
  constructor(schema) {
    this.schema = schema;
    this.errors = {};
  }

  validate(data) {
    const result = this.schema.validate(data);
    this.errors = result.errors;
    return result.isValid;
  }

  getError(field) {
    return this.errors[field];
  }

  hasError(field) {
    return !!this.errors[field];
  }

  getAllErrors() {
    return this.errors;
  }

  clearErrors() {
    this.errors = {};
  }

  clearError(field) {
    delete this.errors[field];
  }
}

/**
 * Async validation for server-side checks
 */
export class AsyncValidator {
  constructor(baseURL = '/api/v1') {
    this.baseURL = baseURL;
  }

  async validateField(field, value, endpoint) {
    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({ [field]: value }),
      });

      const data = await response.json();
      
      if (!response.ok) {
        return data.error?.message || 'Validation failed';
      }
      
      return data.valid ? null : data.message;
    } catch (error) {
      console.error('Async validation error:', error);
      return 'Unable to validate. Please try again.';
    }
  }

  // Check if email is available
  async checkEmailAvailable(email) {
    return this.validateField('email', email, '/validation/email-available');
  }

  // Check if username is available
  async checkUsernameAvailable(username) {
    return this.validateField('username', username, '/validation/username-available');
  }

  // Validate campaign title uniqueness
  async checkCampaignTitleUnique(title, ngoId) {
    return this.validateField('title', title, `/validation/campaign-title-unique/${ngoId}`);
  }
}

/**
 * Real-time validation hook
 */
export const useValidation = (schema, options = {}) => {
  const {
    validateOnChange = false,
    validateOnBlur = true,
    debounceMs = 300
  } = options;

  const [values, setValues] = React.useState({});
  const [errors, setErrors] = React.useState({});
  const [touched, setTouched] = React.useState({});
  const [isValidating, setIsValidating] = React.useState({});
  
  const debounceTimeouts = React.useRef({});
  const validator = React.useMemo(() => new FormValidator(schema), [schema]);

  const validateField = React.useCallback((name, value) => {
    const tempData = { ...values, [name]: value };
    const fieldSchema = new ValidationSchema().field(name, ...schema.rules[name] || []);
    const result = fieldSchema.validate(tempData);
    
    setErrors(prev => ({
      ...prev,
      [name]: result.errors[name] || null
    }));
    
    return !result.errors[name];
  }, [values, schema]);

  const debouncedValidate = React.useCallback((name, value) => {
    if (debounceTimeouts.current[name]) {
      clearTimeout(debounceTimeouts.current[name]);
    }
    
    debounceTimeouts.current[name] = setTimeout(() => {
      validateField(name, value);
    }, debounceMs);
  }, [validateField, debounceMs]);

  const setValue = React.useCallback((name, value) => {
    setValues(prev => ({ ...prev, [name]: value }));
    
    if (validateOnChange || touched[name]) {
      debouncedValidate(name, value);
    }
  }, [validateOnChange, touched, debouncedValidate]);

  const setFieldTouched = React.useCallback((name, isTouched = true) => {
    setTouched(prev => ({ ...prev, [name]: isTouched }));
    
    if (isTouched && validateOnBlur) {
      validateField(name, values[name]);
    }
  }, [validateOnBlur, validateField, values]);

  const validateForm = React.useCallback(() => {
    const isValid = validator.validate(values);
    setErrors(validator.getAllErrors());
    return isValid;
  }, [validator, values]);

  const resetForm = React.useCallback((newValues = {}) => {
    setValues(newValues);
    setErrors({});
    setTouched({});
    setIsValidating({});
  }, []);

  const getFieldProps = React.useCallback((name) => ({
    value: values[name] || '',
    error: touched[name] ? errors[name] : undefined,
    onChange: (e) => setValue(name, e.target.value),
    onBlur: () => setFieldTouched(name, true),
  }), [values, errors, touched, setValue, setFieldTouched]);

  return {
    values,
    errors,
    touched,
    isValidating,
    setValue,
    setFieldTouched,
    validateField,
    validateForm,
    resetForm,
    getFieldProps,
    isValid: Object.keys(errors).length === 0,
  };
};

/**
 * Sanitization functions
 */
export const sanitize = {
  // Remove HTML tags
  html: (input) => {
    if (!input) return input;
    return input.replace(/<[^>]*>/g, '');
  },

  // Sanitize for URL
  url: (input) => {
    if (!input) return input;
    return encodeURIComponent(input);
  },

  // Remove special characters except specified ones
  text: (input, allowed = '') => {
    if (!input) return input;
    const regex = new RegExp(`[^a-zA-Z0-9\\s${allowed.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}]`, 'g');
    return input.replace(regex, '');
  },

  // Normalize phone number
  phone: (input) => {
    if (!input) return input;
    return input.replace(/[^\d+]/g, '');
  },

  // Normalize email
  email: (input) => {
    if (!input) return input;
    return input.trim().toLowerCase();
  },

  // Remove extra whitespace
  whitespace: (input) => {
    if (!input) return input;
    return input.trim().replace(/\s+/g, ' ');
  },
};

/**
 * Business-specific validation schemas
 */
export const schemas = {
  // User registration schema
  userRegistration: new ValidationSchema()
    .field('name', validationRules.required, validationRules.minLength(2), validationRules.maxLength(100))
    .field('email', validationRules.required, validationRules.email)
    .field('password', validationRules.required, validationRules.password)
    .field('confirmPassword', validationRules.required, validationRules.confirmPassword)
    .field('role', validationRules.required),

  // Campaign creation schema
  campaignCreation: new ValidationSchema()
    .field('title', validationRules.required, validationRules.minLength(3), validationRules.maxLength(200))
    .field('description', validationRules.required, validationRules.minLength(10), validationRules.maxLength(5000))
    .field('goalAmount', validationRules.required, validationRules.numeric, validationRules.min(100))
    .field('category', validationRules.required)
    .field('endDate', validationRules.date, validationRules.futureDate)
    .field('website', validationRules.url),

  // Donation schema
  donation: new ValidationSchema()
    .field('amount', validationRules.required, validationRules.numeric, validationRules.min(1))
    .field('campaignId', validationRules.required)
    .field('donorName', validationRules.maxLength(100))
    .field('donorEmail', validationRules.email)
    .field('message', validationRules.maxLength(500)),

  // Milestone creation schema
  milestone: new ValidationSchema()
    .field('title', validationRules.required, validationRules.minLength(3), validationRules.maxLength(200))
    .field('description', validationRules.required, validationRules.minLength(10))
    .field('targetAmount', validationRules.required, validationRules.numeric, validationRules.min(0))
    .field('targetDate', validationRules.date, validationRules.futureDate),

  // NGO profile schema
  ngoProfile: new ValidationSchema()
    .field('name', validationRules.required, validationRules.minLength(2), validationRules.maxLength(200))
    .field('description', validationRules.required, validationRules.minLength(10))
    .field('website', validationRules.url)
    .field('phone', validationRules.phone)
    .field('email', validationRules.required, validationRules.email)
    .field('taxId', validationRules.required)
    .field('address', validationRules.required),
};

/**
 * Utility functions
 */
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export default validationRules;