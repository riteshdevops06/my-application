import { useMemo } from 'react';
import { cn } from '../../utils/cn';

// Base Chart component
const Chart = ({ children, title, description, className, ...props }) => (
  <div className={cn('bg-white p-6 rounded-lg shadow', className)} {...props}>
    {(title || description) && (
      <div className="mb-4">
        {title && <h3 className="text-lg font-semibold text-gray-900">{title}</h3>}
        {description && <p className="text-sm text-gray-600 mt-1">{description}</p>}
      </div>
    )}
    {children}
  </div>
);

// Line Chart Component
const LineChart = ({ 
  data = [], 
  xKey, 
  yKey, 
  height = 300, 
  color = '#3B82F6',
  showGrid = true,
  showDots = true,
  animate = true,
  className 
}) => {
  const { points, maxY, minY } = useMemo(() => {
    if (!data.length) return { points: [], maxY: 0, minY: 0 };
    
    const values = data.map(item => parseFloat(item[yKey]) || 0);
    const maxY = Math.max(...values);
    const minY = Math.min(...values);
    const range = maxY - minY;
    
    const points = data.map((item, index) => {
      const x = (index / (data.length - 1)) * 100;
      const y = 100 - (((parseFloat(item[yKey]) || 0) - minY) / range) * 100;
      return { x, y, value: item[yKey], label: item[xKey] };
    });
    
    return { points, maxY, minY };
  }, [data, xKey, yKey]);

  const pathD = points.length > 1 
    ? `M ${points[0].x} ${points[0].y} ` + 
      points.slice(1).map(p => `L ${p.x} ${p.y}`).join(' ')
    : '';

  return (
    <div className={cn('relative', className)} style={{ height }}>
      <svg viewBox="0 0 100 100" className="w-full h-full">
        {/* Grid lines */}
        {showGrid && (
          <g className="opacity-20">
            {[0, 25, 50, 75, 100].map(y => (
              <line
                key={y}
                x1="0"
                y1={y}
                x2="100"
                y2={y}
                stroke="#6B7280"
                strokeWidth="0.2"
              />
            ))}
            {points.map((point, i) => (
              <line
                key={i}
                x1={point.x}
                y1="0"
                x2={point.x}
                y2="100"
                stroke="#6B7280"
                strokeWidth="0.1"
              />
            ))}
          </g>
        )}
        
        {/* Line */}
        {pathD && (
          <path
            d={pathD}
            fill="none"
            stroke={color}
            strokeWidth="2"
            className={animate ? 'animate-pulse' : ''}
          />
        )}
        
        {/* Dots */}
        {showDots && points.map((point, index) => (
          <circle
            key={index}
            cx={point.x}
            cy={point.y}
            r="1.5"
            fill={color}
            className="hover:r-2 transition-all cursor-pointer"
          />
        ))}
      </svg>
      
      {/* Y-axis labels */}
      <div className="absolute left-0 top-0 h-full flex flex-col justify-between text-xs text-gray-500 -ml-12">
        <span>{maxY}</span>
        <span>{((maxY + minY) / 2).toFixed(1)}</span>
        <span>{minY}</span>
      </div>
      
      {/* X-axis labels */}
      <div className="absolute bottom-0 left-0 w-full flex justify-between text-xs text-gray-500 -mb-6">
        {points.map((point, index) => (
          <span key={index}>{point.label}</span>
        ))}
      </div>
    </div>
  );
};

// Bar Chart Component
const BarChart = ({ 
  data = [], 
  xKey, 
  yKey, 
  height = 300, 
  color = '#3B82F6',
  showValues = true,
  className 
}) => {
  const { bars, maxValue } = useMemo(() => {
    if (!data.length) return { bars: [], maxValue: 0 };
    
    const values = data.map(item => parseFloat(item[yKey]) || 0);
    const maxValue = Math.max(...values);
    
    const bars = data.map((item, index) => {
      const value = parseFloat(item[yKey]) || 0;
      const height = (value / maxValue) * 100;
      const x = (index / data.length) * 100;
      const width = 100 / data.length * 0.8;
      
      return {
        x: x + (100 / data.length - width) / 2,
        y: 100 - height,
        width,
        height,
        value,
        label: item[xKey]
      };
    });
    
    return { bars, maxValue };
  }, [data, xKey, yKey]);

  return (
    <div className={cn('relative', className)} style={{ height }}>
      <svg viewBox="0 0 100 100" className="w-full h-full">
        {bars.map((bar, index) => (
          <g key={index}>
            <rect
              x={bar.x}
              y={bar.y}
              width={bar.width}
              height={bar.height}
              fill={color}
              className="hover:opacity-80 transition-opacity cursor-pointer"
            />
            {showValues && (
              <text
                x={bar.x + bar.width / 2}
                y={bar.y - 1}
                textAnchor="middle"
                className="text-xs fill-gray-600"
                fontSize="3"
              >
                {bar.value}
              </text>
            )}
          </g>
        ))}
      </svg>
      
      {/* Y-axis labels */}
      <div className="absolute left-0 top-0 h-full flex flex-col justify-between text-xs text-gray-500 -ml-12">
        <span>{maxValue}</span>
        <span>{(maxValue / 2).toFixed(1)}</span>
        <span>0</span>
      </div>
      
      {/* X-axis labels */}
      <div className="absolute bottom-0 left-0 w-full flex justify-between text-xs text-gray-500 -mb-6">
        {bars.map((bar, index) => (
          <span key={index} className="text-center" style={{ width: `${100 / bars.length}%` }}>
            {bar.label}
          </span>
        ))}
      </div>
    </div>
  );
};

// Pie Chart Component
const PieChart = ({ 
  data = [], 
  labelKey, 
  valueKey, 
  size = 300, 
  colors = ['#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6'],
  showLegend = true,
  className 
}) => {
  const { segments, total } = useMemo(() => {
    const total = data.reduce((sum, item) => sum + (parseFloat(item[valueKey]) || 0), 0);
    let currentAngle = 0;
    
    const segments = data.map((item, index) => {
      const value = parseFloat(item[valueKey]) || 0;
      const percentage = (value / total) * 100;
      const angle = (value / total) * 360;
      
      const startAngle = currentAngle;
      const endAngle = currentAngle + angle;
      currentAngle += angle;
      
      // Calculate path for pie slice
      const startAngleRad = (startAngle * Math.PI) / 180;
      const endAngleRad = (endAngle * Math.PI) / 180;
      
      const largeArcFlag = angle > 180 ? 1 : 0;
      const x1 = 50 + 40 * Math.cos(startAngleRad);
      const y1 = 50 + 40 * Math.sin(startAngleRad);
      const x2 = 50 + 40 * Math.cos(endAngleRad);
      const y2 = 50 + 40 * Math.sin(endAngleRad);
      
      const pathData = [
        `M 50 50`,
        `L ${x1} ${y1}`,
        `A 40 40 0 ${largeArcFlag} 1 ${x2} ${y2}`,
        'Z'
      ].join(' ');
      
      return {
        path: pathData,
        color: colors[index % colors.length],
        label: item[labelKey],
        value,
        percentage: percentage.toFixed(1)
      };
    });
    
    return { segments, total };
  }, [data, labelKey, valueKey, colors]);

  return (
    <div className={cn('flex items-center justify-center', className)}>
      <div className="flex items-center space-x-8">
        <div style={{ width: size, height: size }}>
          <svg viewBox="0 0 100 100" className="w-full h-full">
            {segments.map((segment, index) => (
              <path
                key={index}
                d={segment.path}
                fill={segment.color}
                className="hover:opacity-80 transition-opacity cursor-pointer"
                stroke="#ffffff"
                strokeWidth="0.5"
              />
            ))}
          </svg>
        </div>
        
        {showLegend && (
          <div className="space-y-2">
            {segments.map((segment, index) => (
              <div key={index} className="flex items-center space-x-2">
                <div
                  className="w-3 h-3 rounded-full"
                  style={{ backgroundColor: segment.color }}
                ></div>
                <span className="text-sm text-gray-700">
                  {segment.label}: {segment.percentage}%
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

// Donut Chart Component
const DonutChart = ({ 
  data = [], 
  labelKey, 
  valueKey, 
  size = 300, 
  colors = ['#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6'],
  showCenter = true,
  centerContent,
  className 
}) => {
  const { segments, total } = useMemo(() => {
    const total = data.reduce((sum, item) => sum + (parseFloat(item[valueKey]) || 0), 0);
    let currentAngle = 0;
    
    const segments = data.map((item, index) => {
      const value = parseFloat(item[valueKey]) || 0;
      const percentage = (value / total) * 100;
      const angle = (value / total) * 360;
      
      const startAngle = currentAngle;
      const endAngle = currentAngle + angle;
      currentAngle += angle;
      
      const startAngleRad = (startAngle * Math.PI) / 180;
      const endAngleRad = (endAngle * Math.PI) / 180;
      
      const largeArcFlag = angle > 180 ? 1 : 0;
      
      // Outer arc
      const x1 = 50 + 40 * Math.cos(startAngleRad);
      const y1 = 50 + 40 * Math.sin(startAngleRad);
      const x2 = 50 + 40 * Math.cos(endAngleRad);
      const y2 = 50 + 40 * Math.sin(endAngleRad);
      
      // Inner arc  
      const x3 = 50 + 25 * Math.cos(endAngleRad);
      const y3 = 50 + 25 * Math.sin(endAngleRad);
      const x4 = 50 + 25 * Math.cos(startAngleRad);
      const y4 = 50 + 25 * Math.sin(startAngleRad);
      
      const pathData = [
        `M ${x1} ${y1}`,
        `A 40 40 0 ${largeArcFlag} 1 ${x2} ${y2}`,
        `L ${x3} ${y3}`,
        `A 25 25 0 ${largeArcFlag} 0 ${x4} ${y4}`,
        'Z'
      ].join(' ');
      
      return {
        path: pathData,
        color: colors[index % colors.length],
        label: item[labelKey],
        value,
        percentage: percentage.toFixed(1)
      };
    });
    
    return { segments, total };
  }, [data, labelKey, valueKey, colors]);

  return (
    <div className={cn('relative flex items-center justify-center', className)}>
      <div style={{ width: size, height: size }}>
        <svg viewBox="0 0 100 100" className="w-full h-full">
          {segments.map((segment, index) => (
            <path
              key={index}
              d={segment.path}
              fill={segment.color}
              className="hover:opacity-80 transition-opacity cursor-pointer"
              stroke="#ffffff"
              strokeWidth="0.5"
            />
          ))}
        </svg>
        
        {showCenter && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="text-center">
              {centerContent || (
                <>
                  <div className="text-2xl font-bold text-gray-900">{total}</div>
                  <div className="text-sm text-gray-600">Total</div>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// Area Chart Component  
const AreaChart = ({ 
  data = [], 
  xKey, 
  yKey, 
  height = 300, 
  color = '#3B82F6',
  opacity = 0.3,
  className 
}) => {
  const { points, areaPath, maxY, minY } = useMemo(() => {
    if (!data.length) return { points: [], areaPath: '', maxY: 0, minY: 0 };
    
    const values = data.map(item => parseFloat(item[yKey]) || 0);
    const maxY = Math.max(...values);
    const minY = Math.min(...values);
    const range = maxY - minY;
    
    const points = data.map((item, index) => {
      const x = (index / (data.length - 1)) * 100;
      const y = 100 - (((parseFloat(item[yKey]) || 0) - minY) / range) * 100;
      return { x, y, value: item[yKey], label: item[xKey] };
    });
    
    const linePath = points.length > 1 
      ? `M ${points[0].x} ${points[0].y} ` + 
        points.slice(1).map(p => `L ${p.x} ${p.y}`).join(' ')
      : '';
      
    const areaPath = linePath 
      ? `${linePath} L ${points[points.length - 1].x} 100 L ${points[0].x} 100 Z`
      : '';
    
    return { points, areaPath, maxY, minY };
  }, [data, xKey, yKey]);

  return (
    <div className={cn('relative', className)} style={{ height }}>
      <svg viewBox="0 0 100 100" className="w-full h-full">
        {/* Area fill */}
        {areaPath && (
          <path
            d={areaPath}
            fill={color}
            fillOpacity={opacity}
          />
        )}
        
        {/* Line */}
        {points.length > 1 && (
          <path
            d={`M ${points[0].x} ${points[0].y} ` + 
              points.slice(1).map(p => `L ${p.x} ${p.y}`).join(' ')}
            fill="none"
            stroke={color}
            strokeWidth="2"
          />
        )}
      </svg>
    </div>
  );
};

// Attach sub-components
Chart.Line = LineChart;
Chart.Bar = BarChart;
Chart.Pie = PieChart;
Chart.Donut = DonutChart;
Chart.Area = AreaChart;

export default Chart;