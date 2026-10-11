import { useState, useEffect, useMemo, useRef } from 'react';
import { cn } from '../../utils/cn';
import Input from './Input';
import Select from './Select';
import Button from './Button';

// Advanced Search Component
const AdvancedSearch = ({
  onSearch,
  searchFields = [],
  filterOptions = [],
  sortOptions = [],
  initialFilters = {},
  placeholder = "Search...",
  className,
  showFilters = true,
  showSort = true,
  debounceMs = 300,
  ...props
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [activeFilters, setActiveFilters] = useState(initialFilters);
  const [sortBy, setSortBy] = useState('');
  const [sortDirection, setSortDirection] = useState('asc');
  const [isExpanded, setIsExpanded] = useState(false);
  const debounceRef = useRef(null);

  // Debounced search effect
  useEffect(() => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      if (onSearch) {
        onSearch({
          searchTerm,
          filters: activeFilters,
          sortBy,
          sortDirection,
        });
      }
    }, debounceMs);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [searchTerm, activeFilters, sortBy, sortDirection, onSearch, debounceMs]);

  // Update filter
  const updateFilter = (key, value) => {
    setActiveFilters(prev => {
      if (value === '' || value === null || value === undefined) {
        const { [key]: removed, ...rest } = prev;
        return rest;
      }
      return { ...prev, [key]: value };
    });
  };

  // Clear all filters
  const clearFilters = () => {
    setSearchTerm('');
    setActiveFilters({});
    setSortBy('');
    setSortDirection('asc');
  };

  // Get active filter count
  const activeFilterCount = Object.keys(activeFilters).filter(
    key => activeFilters[key] !== '' && activeFilters[key] !== null
  ).length;

  return (
    <div className={cn('bg-white border border-gray-200 rounded-lg p-4', className)} {...props}>
      {/* Search Bar */}
      <div className="flex items-center space-x-3">
        <div className="flex-1">
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <Input
              type="text"
              placeholder={placeholder}
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Filter Toggle */}
        {(showFilters || showSort) && (
          <Button
            variant="outline"
            onClick={() => setIsExpanded(!isExpanded)}
            className="flex items-center space-x-1"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.707A1 1 0 013 7V4z" />
            </svg>
            <span>Filters</span>
            {activeFilterCount > 0 && (
              <span className="bg-blue-100 text-blue-800 text-xs rounded-full px-2 py-0.5">
                {activeFilterCount}
              </span>
            )}
            <svg 
              className={cn('w-4 h-4 transform transition-transform', isExpanded && 'rotate-180')} 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </Button>
        )}

        {/* Clear Filters */}
        {(activeFilterCount > 0 || searchTerm) && (
          <Button variant="outline" onClick={clearFilters} size="sm">
            Clear
          </Button>
        )}
      </div>

      {/* Expanded Filters */}
      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {/* Filter Options */}
            {showFilters && filterOptions.map((filter) => (
              <div key={filter.key}>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {filter.label}
                </label>
                
                {filter.type === 'select' ? (
                  <Select
                    value={activeFilters[filter.key] || ''}
                    onChange={(e) => updateFilter(filter.key, e.target.value)}
                    size="sm"
                  >
                    <option value="">All {filter.label}</option>
                    {filter.options.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </Select>
                ) : filter.type === 'date' ? (
                  <Input
                    type="date"
                    value={activeFilters[filter.key] || ''}
                    onChange={(e) => updateFilter(filter.key, e.target.value)}
                    size="sm"
                  />
                ) : filter.type === 'range' ? (
                  <div className="flex space-x-2">
                    <Input
                      type="number"
                      placeholder="Min"
                      value={activeFilters[`${filter.key}_min`] || ''}
                      onChange={(e) => updateFilter(`${filter.key}_min`, e.target.value)}
                      size="sm"
                    />
                    <Input
                      type="number"
                      placeholder="Max"
                      value={activeFilters[`${filter.key}_max`] || ''}
                      onChange={(e) => updateFilter(`${filter.key}_max`, e.target.value)}
                      size="sm"
                    />
                  </div>
                ) : (
                  <Input
                    type="text"
                    placeholder={`Filter by ${filter.label.toLowerCase()}`}
                    value={activeFilters[filter.key] || ''}
                    onChange={(e) => updateFilter(filter.key, e.target.value)}
                    size="sm"
                  />
                )}
              </div>
            ))}

            {/* Sort Options */}
            {showSort && sortOptions.length > 0 && (
              <div className="lg:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Sort By
                </label>
                <div className="flex space-x-2">
                  <Select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value)}
                    size="sm"
                    className="flex-1"
                  >
                    <option value="">Default Order</option>
                    {sortOptions.map((option) => (
                      <option key={option.key} value={option.key}>
                        {option.label}
                      </option>
                    ))}
                  </Select>
                  <Select
                    value={sortDirection}
                    onChange={(e) => setSortDirection(e.target.value)}
                    size="sm"
                  >
                    <option value="asc">Ascending</option>
                    <option value="desc">Descending</option>
                  </Select>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Active Filters Display */}
      {activeFilterCount > 0 && (
        <div className="mt-3 flex flex-wrap gap-2">
          {Object.entries(activeFilters).map(([key, value]) => {
            if (!value) return null;
            
            const filter = filterOptions.find(f => f.key === key || key.startsWith(f.key));
            const filterLabel = filter?.label || key;
            
            return (
              <span
                key={key}
                className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                {filterLabel}: {value}
                <button
                  onClick={() => updateFilter(key, '')}
                  className="ml-1 inline-flex items-center justify-center w-4 h-4 rounded-full bg-blue-200 text-blue-600 hover:bg-blue-300"
                >
                  ×
                </button>
              </span>
            );
          })}
        </div>
      )}
    </div>
  );
};

// Quick Search Component
const QuickSearch = ({
  value,
  onChange,
  onSearch,
  placeholder = "Quick search...",
  suggestions = [],
  showSuggestions = true,
  loading = false,
  className,
  ...props
}) => {
  const [isFocused, setIsFocused] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const inputRef = useRef(null);
  const suggestionRefs = useRef([]);

  // Handle keyboard navigation
  const handleKeyDown = (e) => {
    if (!showSuggestions || !suggestions.length) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex(prev => 
          prev < suggestions.length - 1 ? prev + 1 : 0
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex(prev => 
          prev > 0 ? prev - 1 : suggestions.length - 1
        );
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0) {
          const selected = suggestions[selectedIndex];
          onChange(selected.value || selected);
          if (onSearch) onSearch(selected.value || selected);
        } else if (onSearch) {
          onSearch(value);
        }
        setIsFocused(false);
        break;
      case 'Escape':
        setIsFocused(false);
        inputRef.current?.blur();
        break;
    }
  };

  // Scroll selected suggestion into view
  useEffect(() => {
    if (selectedIndex >= 0 && suggestionRefs.current[selectedIndex]) {
      suggestionRefs.current[selectedIndex].scrollIntoView({
        block: 'nearest',
        behavior: 'smooth'
      });
    }
  }, [selectedIndex]);

  return (
    <div className={cn('relative', className)} {...props}>
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          {loading ? (
            <div className="animate-spin rounded-full h-4 w-4 border-2 border-gray-300 border-t-blue-600" />
          ) : (
            <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          )}
        </div>
        
        <Input
          ref={inputRef}
          type="text"
          placeholder={placeholder}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setTimeout(() => setIsFocused(false), 200)}
          className="pl-10 pr-10"
        />
        
        {value && (
          <button
            onClick={() => {
              onChange('');
              inputRef.current?.focus();
            }}
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>

      {/* Suggestions Dropdown */}
      {isFocused && showSuggestions && suggestions.length > 0 && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-auto">
          {suggestions.map((suggestion, index) => {
            const isObject = typeof suggestion === 'object';
            const displayText = isObject ? suggestion.label : suggestion;
            const value = isObject ? suggestion.value : suggestion;
            
            return (
              <button
                key={index}
                ref={el => suggestionRefs.current[index] = el}
                onClick={() => {
                  onChange(value);
                  if (onSearch) onSearch(value);
                  setIsFocused(false);
                }}
                className={cn(
                  'w-full px-4 py-2 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none',
                  selectedIndex === index && 'bg-blue-50 text-blue-700'
                )}
              >
                <div>
                  <div className="text-sm font-medium">{displayText}</div>
                  {isObject && suggestion.description && (
                    <div className="text-xs text-gray-500">{suggestion.description}</div>
                  )}
                </div>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};

// Search Results Component
const SearchResults = ({
  results = [],
  loading = false,
  error = null,
  emptyMessage = "No results found",
  renderItem,
  onLoadMore,
  hasMore = false,
  className,
  ...props
}) => {
  if (loading && results.length === 0) {
    return (
      <div className={cn('flex items-center justify-center py-12', className)}>
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-2 border-gray-300 border-t-blue-600 mx-auto mb-2" />
          <p className="text-gray-500">Searching...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={cn('flex items-center justify-center py-12', className)}>
        <div className="text-center">
          <div className="text-red-500 mb-2">
            <svg className="w-8 h-8 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-red-600 font-medium">Search Error</p>
          <p className="text-gray-500 text-sm">{error.message || 'An error occurred while searching'}</p>
        </div>
      </div>
    );
  }

  if (results.length === 0) {
    return (
      <div className={cn('flex items-center justify-center py-12', className)}>
        <div className="text-center">
          <div className="text-gray-400 mb-2">
            <svg className="w-8 h-8 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <p className="text-gray-500">{emptyMessage}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={cn('space-y-4', className)} {...props}>
      {results.map((item, index) => (
        <div key={item.id || index}>
          {renderItem ? renderItem(item, index) : (
            <div className="p-4 border border-gray-200 rounded-lg">
              {JSON.stringify(item)}
            </div>
          )}
        </div>
      ))}

      {/* Load More Button */}
      {hasMore && (
        <div className="flex justify-center pt-4">
          <Button
            onClick={onLoadMore}
            variant="outline"
            loading={loading}
          >
            Load More Results
          </Button>
        </div>
      )}

      {/* Loading indicator for additional results */}
      {loading && results.length > 0 && (
        <div className="flex justify-center py-4">
          <div className="animate-spin rounded-full h-6 w-6 border-2 border-gray-300 border-t-blue-600" />
        </div>
      )}
    </div>
  );
};

// Filter Tags Component
const FilterTags = ({ filters, onRemoveFilter, onClearAll, className }) => {
  const filterEntries = Object.entries(filters).filter(([_, value]) => 
    value !== null && value !== undefined && value !== ''
  );

  if (filterEntries.length === 0) return null;

  return (
    <div className={cn('flex flex-wrap gap-2 items-center', className)}>
      <span className="text-sm text-gray-600">Active filters:</span>
      
      {filterEntries.map(([key, value]) => (
        <span
          key={key}
          className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
        >
          {key}: {value}
          <button
            onClick={() => onRemoveFilter(key)}
            className="ml-2 inline-flex items-center justify-center w-4 h-4 rounded-full bg-blue-200 text-blue-600 hover:bg-blue-300"
          >
            ×
          </button>
        </span>
      ))}
      
      <Button
        variant="ghost"
        size="sm"
        onClick={onClearAll}
        className="text-gray-500 hover:text-gray-700"
      >
        Clear all
      </Button>
    </div>
  );
};

export { AdvancedSearch, QuickSearch, SearchResults, FilterTags };
export default AdvancedSearch;