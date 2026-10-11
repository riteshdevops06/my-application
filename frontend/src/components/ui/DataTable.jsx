import { useState, useMemo } from 'react';
import { cn } from '../../utils/cn';
import Button from './Button';
import Input from './Input';
import Select from './Select';

const DataTable = ({
  data = [],
  columns = [],
  searchable = true,
  sortable = true,
  filterable = false,
  pagination = true,
  pageSize = 10,
  loading = false,
  emptyMessage = "No data available",
  className,
  onRowClick,
  selectable = false,
  onSelectionChange,
  ...props
}) => {
  const [search, setSearch] = useState('');
  const [sortField, setSortField] = useState('');
  const [sortDirection, setSortDirection] = useState('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [filters, setFilters] = useState({});
  const [selectedRows, setSelectedRows] = useState(new Set());

  // Filter data based on search and filters
  const filteredData = useMemo(() => {
    let filtered = [...data];

    // Apply search filter
    if (search && searchable) {
      const searchLower = search.toLowerCase();
      filtered = filtered.filter(item =>
        columns.some(column => {
          const value = getNestedValue(item, column.key);
          return String(value).toLowerCase().includes(searchLower);
        })
      );
    }

    // Apply column filters
    if (filterable) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value && value !== '') {
          filtered = filtered.filter(item => {
            const itemValue = getNestedValue(item, key);
            if (typeof value === 'string') {
              return String(itemValue).toLowerCase().includes(value.toLowerCase());
            }
            return itemValue === value;
          });
        }
      });
    }

    return filtered;
  }, [data, search, filters, columns, searchable, filterable]);

  // Sort data
  const sortedData = useMemo(() => {
    if (!sortField || !sortable) return filteredData;

    return [...filteredData].sort((a, b) => {
      const aValue = getNestedValue(a, sortField);
      const bValue = getNestedValue(b, sortField);

      let comparison = 0;
      if (aValue > bValue) comparison = 1;
      if (aValue < bValue) comparison = -1;

      return sortDirection === 'desc' ? comparison * -1 : comparison;
    });
  }, [filteredData, sortField, sortDirection, sortable]);

  // Paginate data
  const paginatedData = useMemo(() => {
    if (!pagination) return sortedData;

    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return sortedData.slice(startIndex, endIndex);
  }, [sortedData, currentPage, pageSize, pagination]);

  // Calculate pagination info
  const totalPages = Math.ceil(sortedData.length / pageSize);
  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, sortedData.length);

  // Helper function to get nested object values
  const getNestedValue = (obj, path) => {
    return path.split('.').reduce((current, key) => current?.[key], obj);
  };

  // Handle sorting
  const handleSort = (field) => {
    if (!sortable) return;
    
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('asc');
    }
  };

  // Handle row selection
  const handleRowSelection = (rowId, checked) => {
    const newSelection = new Set(selectedRows);
    if (checked) {
      newSelection.add(rowId);
    } else {
      newSelection.delete(rowId);
    }
    setSelectedRows(newSelection);
    
    if (onSelectionChange) {
      onSelectionChange(Array.from(newSelection));
    }
  };

  // Handle select all
  const handleSelectAll = (checked) => {
    if (checked) {
      const allIds = new Set(paginatedData.map(item => item.id || item._id));
      setSelectedRows(allIds);
      if (onSelectionChange) {
        onSelectionChange(Array.from(allIds));
      }
    } else {
      setSelectedRows(new Set());
      if (onSelectionChange) {
        onSelectionChange([]);
      }
    }
  };

  // Render filter input for a column
  const renderFilter = (column) => {
    if (!column.filterable) return null;

    if (column.filterOptions) {
      return (
        <Select
          value={filters[column.key] || ''}
          onChange={(e) => setFilters(prev => ({ ...prev, [column.key]: e.target.value }))}
          className="text-xs"
          size="sm"
        >
          <option value="">All</option>
          {column.filterOptions.map(option => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </Select>
      );
    }

    return (
      <Input
        type="text"
        placeholder={`Filter ${column.title}`}
        value={filters[column.key] || ''}
        onChange={(e) => setFilters(prev => ({ ...prev, [column.key]: e.target.value }))}
        className="text-xs"
        size="sm"
      />
    );
  };

  // Render cell content
  const renderCell = (item, column) => {
    const value = getNestedValue(item, column.key);
    
    if (column.render) {
      return column.render(value, item);
    }
    
    if (column.type === 'date' && value) {
      return new Date(value).toLocaleDateString();
    }
    
    if (column.type === 'currency' && value) {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(value);
    }
    
    return value;
  };

  // Render sort icon
  const renderSortIcon = (field) => {
    if (!sortable || sortField !== field) {
      return (
        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
        </svg>
      );
    }

    return sortDirection === 'asc' ? (
      <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
      </svg>
    ) : (
      <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
      </svg>
    );
  };

  const allSelected = paginatedData.length > 0 && paginatedData.every(item => selectedRows.has(item.id || item._id));
  const someSelected = paginatedData.some(item => selectedRows.has(item.id || item._id));

  return (
    <div className={cn('bg-white rounded-lg shadow', className)} {...props}>
      {/* Header with search and actions */}
      {(searchable || filterable) && (
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            {searchable && (
              <div className="flex-1 max-w-md">
                <Input
                  type="text"
                  placeholder="Search..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="w-full"
                />
              </div>
            )}
            
            {selectedRows.size > 0 && (
              <div className="flex items-center space-x-2">
                <span className="text-sm text-gray-600">
                  {selectedRows.size} selected
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setSelectedRows(new Set());
                    if (onSelectionChange) onSelectionChange([]);
                  }}
                >
                  Clear
                </Button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {selectable && (
                <th className="w-12 px-4 py-3 text-left">
                  <input
                    type="checkbox"
                    checked={allSelected}
                    ref={input => {
                      if (input) input.indeterminate = someSelected && !allSelected;
                    }}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
              )}
              
              {columns.map((column) => (
                <th key={column.key} className="px-4 py-3 text-left">
                  <div className="flex flex-col space-y-2">
                    <div
                      className={cn(
                        'flex items-center space-x-1',
                        sortable && column.sortable !== false ? 'cursor-pointer hover:text-gray-900' : '',
                        'text-xs font-medium text-gray-500 uppercase tracking-wider'
                      )}
                      onClick={() => column.sortable !== false && handleSort(column.key)}
                    >
                      <span>{column.title}</span>
                      {sortable && column.sortable !== false && renderSortIcon(column.key)}
                    </div>
                    
                    {filterable && column.filterable && (
                      <div className="min-w-0">
                        {renderFilter(column)}
                      </div>
                    )}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          
          <tbody className="divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={columns.length + (selectable ? 1 : 0)} className="px-4 py-8 text-center">
                  <div className="flex items-center justify-center space-x-2">
                    <div className="animate-spin rounded-full h-5 w-5 border-2 border-blue-600 border-t-transparent"></div>
                    <span className="text-gray-500">Loading...</span>
                  </div>
                </td>
              </tr>
            ) : paginatedData.length === 0 ? (
              <tr>
                <td colSpan={columns.length + (selectable ? 1 : 0)} className="px-4 py-8 text-center text-gray-500">
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              paginatedData.map((item, index) => {
                const rowId = item.id || item._id || index;
                const isSelected = selectedRows.has(rowId);
                
                return (
                  <tr
                    key={rowId}
                    className={cn(
                      'hover:bg-gray-50 transition-colors',
                      onRowClick && 'cursor-pointer',
                      isSelected && 'bg-blue-50'
                    )}
                    onClick={() => onRowClick?.(item)}
                  >
                    {selectable && (
                      <td className="px-4 py-4">
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={(e) => {
                            e.stopPropagation();
                            handleRowSelection(rowId, e.target.checked);
                          }}
                          className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                      </td>
                    )}
                    
                    {columns.map((column) => (
                      <td
                        key={column.key}
                        className={cn(
                          'px-4 py-4 text-sm',
                          column.align === 'center' && 'text-center',
                          column.align === 'right' && 'text-right',
                          column.className
                        )}
                      >
                        {renderCell(item, column)}
                      </td>
                    ))}
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {pagination && !loading && sortedData.length > 0 && (
        <div className="px-4 py-3 border-t border-gray-200 flex items-center justify-between">
          <div className="text-sm text-gray-700">
            Showing {startItem} to {endItem} of {sortedData.length} results
          </div>
          
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === 1}
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
            >
              Previous
            </Button>
            
            <div className="flex items-center space-x-1">
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                let pageNum;
                if (totalPages <= 5) {
                  pageNum = i + 1;
                } else if (currentPage <= 3) {
                  pageNum = i + 1;
                } else if (currentPage >= totalPages - 2) {
                  pageNum = totalPages - 4 + i;
                } else {
                  pageNum = currentPage - 2 + i;
                }
                
                return (
                  <Button
                    key={pageNum}
                    variant={currentPage === pageNum ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => setCurrentPage(pageNum)}
                  >
                    {pageNum}
                  </Button>
                );
              })}
            </div>
            
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === totalPages}
              onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default DataTable;