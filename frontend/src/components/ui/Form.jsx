import { createContext, useContext, useState, useEffect } from 'react';
import { cn } from '../../utils/cn';
import Input from './Input';
import Select from './Select';
import Textarea from './Textarea';
import Button from './Button';

// Form Context
const FormContext = createContext({});

// Validation rules
const validationRules = {
  required: (value) => {
    if (value === null || value === undefined || value === '') {
      return 'This field is required';
    }
    return null;
  },
  
  email: (value) => {
    if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      return 'Please enter a valid email address';
    }
    return null;
  },
  
  minLength: (min) => (value) => {
    if (value && value.length < min) {
      return `Must be at least ${min} characters long`;
    }
    return null;
  },
  
  maxLength: (max) => (value) => {
    if (value && value.length > max) {
      return `Must be no more than ${max} characters long`;
    }
    return null;
  },
  
  min: (min) => (value) => {
    if (value !== null && value !== undefined && parseFloat(value) < min) {
      return `Must be at least ${min}`;
    }
    return null;
  },
  
  max: (max) => (value) => {
    if (value !== null && value !== undefined && parseFloat(value) > max) {
      return `Must be no more than ${max}`;
    }
    return null;
  },
  
  pattern: (regex, message = 'Invalid format') => (value) => {
    if (value && !regex.test(value)) {
      return message;
    }
    return null;
  },
  
  custom: (validator) => validator,
};

// Main Form component
const Form = ({ 
  children, 
  onSubmit, 
  initialValues = {}, 
  validationSchema = {},
  validateOnChange = false,
  validateOnBlur = true,
  className,
  ...props 
}) => {
  const [values, setValues] = useState(initialValues);
  const [errors, setErrors] = useState({});
  const [touched, setTouched] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitCount, setSubmitCount] = useState(0);

  // Validate a single field
  const validateField = (name, value) => {
    const fieldRules = validationSchema[name];
    if (!fieldRules) return null;

    for (const rule of fieldRules) {
      let validator;
      let ruleParams;

      if (typeof rule === 'string') {
        validator = validationRules[rule];
      } else if (typeof rule === 'object') {
        const [ruleName, ...params] = Object.keys(rule);
        validator = validationRules[ruleName];
        ruleParams = params.length > 0 ? rule[ruleName] : undefined;
      } else if (typeof rule === 'function') {
        validator = rule;
      }

      if (validator) {
        const error = ruleParams !== undefined 
          ? validator(ruleParams)(value)
          : validator(value);
        
        if (error) {
          return error;
        }
      }
    }

    return null;
  };

  // Validate all fields
  const validateForm = () => {
    const newErrors = {};
    
    Object.keys(validationSchema).forEach(fieldName => {
      const error = validateField(fieldName, values[fieldName]);
      if (error) {
        newErrors[fieldName] = error;
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle field change
  const handleFieldChange = (name, value) => {
    setValues(prev => ({ ...prev, [name]: value }));

    if (validateOnChange || (submitCount > 0)) {
      const error = validateField(name, value);
      setErrors(prev => ({
        ...prev,
        [name]: error
      }));
    }
  };

  // Handle field blur
  const handleFieldBlur = (name) => {
    setTouched(prev => ({ ...prev, [name]: true }));

    if (validateOnBlur || (submitCount > 0)) {
      const error = validateField(name, values[name]);
      setErrors(prev => ({
        ...prev,
        [name]: error
      }));
    }
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitCount(prev => prev + 1);
    
    const isValid = validateForm();
    if (!isValid) {
      // Mark all fields as touched to show errors
      const allTouched = {};
      Object.keys(validationSchema).forEach(field => {
        allTouched[field] = true;
      });
      setTouched(allTouched);
      return;
    }

    if (onSubmit) {
      setIsSubmitting(true);
      try {
        await onSubmit(values, { setErrors, setValues, resetForm });
      } catch (error) {
        console.error('Form submission error:', error);
      } finally {
        setIsSubmitting(false);
      }
    }
  };

  // Reset form
  const resetForm = () => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
    setSubmitCount(0);
  };

  // Get field props
  const getFieldProps = (name) => ({
    name,
    value: values[name] || '',
    error: (touched[name] || submitCount > 0) ? errors[name] : undefined,
    onChange: (e) => {
      const value = e.target.type === 'checkbox' ? e.target.checked : e.target.value;
      handleFieldChange(name, value);
    },
    onBlur: () => handleFieldBlur(name),
  });

  const contextValue = {
    values,
    errors,
    touched,
    isSubmitting,
    submitCount,
    handleFieldChange,
    handleFieldBlur,
    getFieldProps,
    validateField,
    resetForm,
  };

  return (
    <FormContext.Provider value={contextValue}>
      <form onSubmit={handleSubmit} className={cn('space-y-6', className)} {...props}>
        {children}
      </form>
    </FormContext.Provider>
  );
};

// Field component
const Field = ({ 
  name, 
  label, 
  type = 'text', 
  component = 'input',
  placeholder,
  helpText,
  required,
  className,
  ...props 
}) => {
  const { getFieldProps } = useContext(FormContext);
  const fieldProps = getFieldProps(name);

  const renderField = () => {
    switch (component) {
      case 'textarea':
        return (
          <Textarea
            {...fieldProps}
            {...props}
            placeholder={placeholder}
            className={className}
          />
        );
      
      case 'select':
        return (
          <Select
            {...fieldProps}
            {...props}
            className={className}
          >
            {props.children}
          </Select>
        );
      
      default:
        return (
          <Input
            {...fieldProps}
            {...props}
            type={type}
            placeholder={placeholder}
            className={className}
          />
        );
    }
  };

  return (
    <div className="space-y-1">
      {label && (
        <label className="block text-sm font-medium text-gray-700">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      {renderField()}
      
      {helpText && (
        <p className="text-sm text-gray-500">{helpText}</p>
      )}
      
      {fieldProps.error && (
        <p className="text-sm text-red-600">{fieldProps.error}</p>
      )}
    </div>
  );
};

// FieldArray component for dynamic fields
const FieldArray = ({ name, children }) => {
  const { values, handleFieldChange } = useContext(FormContext);
  const arrayValue = values[name] || [];

  const add = (item = {}) => {
    const newArray = [...arrayValue, item];
    handleFieldChange(name, newArray);
  };

  const remove = (index) => {
    const newArray = arrayValue.filter((_, i) => i !== index);
    handleFieldChange(name, newArray);
  };

  const move = (from, to) => {
    const newArray = [...arrayValue];
    const [removed] = newArray.splice(from, 1);
    newArray.splice(to, 0, removed);
    handleFieldChange(name, newArray);
  };

  const update = (index, item) => {
    const newArray = [...arrayValue];
    newArray[index] = item;
    handleFieldChange(name, newArray);
  };

  return children({
    fields: arrayValue,
    add,
    remove,
    move,
    update,
  });
};

// FormSection component for grouping fields
const FormSection = ({ title, description, children, className, collapsible = false }) => {
  const [isOpen, setIsOpen] = useState(true);

  return (
    <div className={cn('border border-gray-200 rounded-lg p-6', className)}>
      {title && (
        <div 
          className={cn(
            'mb-4 pb-4 border-b border-gray-200',
            collapsible && 'cursor-pointer'
          )}
          onClick={() => collapsible && setIsOpen(!isOpen)}
        >
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">{title}</h3>
            {collapsible && (
              <svg
                className={cn('w-5 h-5 transform transition-transform', isOpen && 'rotate-180')}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            )}
          </div>
          {description && (
            <p className="text-sm text-gray-600 mt-1">{description}</p>
          )}
        </div>
      )}
      
      {(!collapsible || isOpen) && (
        <div className="space-y-4">
          {children}
        </div>
      )}
    </div>
  );
};

// FormActions component for submit/cancel buttons
const FormActions = ({ 
  submitText = 'Submit', 
  cancelText = 'Cancel', 
  onCancel,
  submitDisabled = false,
  className,
  ...props 
}) => {
  const { isSubmitting } = useContext(FormContext);

  return (
    <div className={cn('flex justify-end space-x-3 pt-6 border-t border-gray-200', className)} {...props}>
      {onCancel && (
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting}
        >
          {cancelText}
        </Button>
      )}
      
      <Button
        type="submit"
        loading={isSubmitting}
        disabled={submitDisabled || isSubmitting}
      >
        {submitText}
      </Button>
    </div>
  );
};

// Hook to use form context
const useFormContext = () => {
  const context = useContext(FormContext);
  if (!context) {
    throw new Error('useFormContext must be used within a Form component');
  }
  return context;
};

// Validation schema builder
const createValidationSchema = (schema) => {
  const processedSchema = {};
  
  Object.entries(schema).forEach(([field, rules]) => {
    if (Array.isArray(rules)) {
      processedSchema[field] = rules;
    } else if (typeof rules === 'object') {
      processedSchema[field] = Object.entries(rules).map(([rule, params]) => ({
        [rule]: params
      }));
    }
  });
  
  return processedSchema;
};

// Attach sub-components
Form.Field = Field;
Form.FieldArray = FieldArray;
Form.Section = FormSection;
Form.Actions = FormActions;

export default Form;
export { useFormContext, createValidationSchema, validationRules };