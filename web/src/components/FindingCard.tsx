import React from 'react';
import type { AnalysisResult } from '../types';

interface FindingCardProps {
  finding: AnalysisResult;
  index: number;
}

const severityConfig = {
  Critical: {
    color: 'text-red-700',
    bg: 'bg-red-50',
    border: 'border-red-200',
    badgeBg: 'bg-red-100',
    badgeText: 'text-red-800',
    icon: (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  High: {
    color: 'text-orange-700',
    bg: 'bg-orange-50',
    border: 'border-orange-200',
    badgeBg: 'bg-orange-100',
    badgeText: 'text-orange-800',
    icon: (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  Medium: {
    color: 'text-yellow-700',
    bg: 'bg-yellow-50',
    border: 'border-yellow-200',
    badgeBg: 'bg-yellow-100',
    badgeText: 'text-yellow-800',
    icon: (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  Low: {
    color: 'text-blue-700',
    bg: 'bg-blue-50',
    border: 'border-blue-200',
    badgeBg: 'bg-blue-100',
    badgeText: 'text-blue-800',
    icon: (
      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
      </svg>
    ),
  },
};

// Agent icon mapping
const agentIcons = {
  'License Agent': (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
    </svg>
  ),
  'Dependency Health Agent': (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
    </svg>
  ),
  'Proactive Vulnerability Agent': (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
  ),
};

export const FindingCard: React.FC<FindingCardProps> = ({ finding, index }) => {
  const { agent_name, finding: findingText, severity } = finding;
  const config = severityConfig[severity];
  const agentIcon = agentIcons[agent_name as keyof typeof agentIcons] || agentIcons['License Agent'];

  // Extract component name if present in the finding text
  const extractComponentInfo = (text: string) => {
    // Look for patterns like "Component 'name'" or "'name' (version)"
    const componentMatch = text.match(/(?:Component\s+)?'([^']+)'(?:\s+\(([^)]+)\))?/);
    if (componentMatch) {
      return {
        name: componentMatch[1],
        version: componentMatch[2],
        hasComponent: true,
      };
    }
    return { hasComponent: false };
  };

  const componentInfo = extractComponentInfo(findingText);

  return (
    <div className={`rounded-lg border p-6 transition-all duration-200 hover:shadow-md ${config.bg} ${config.border}`}>
      {/* Header */}
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center space-x-3">
          {/* Agent Icon */}
          <div className="text-gray-500">
            {agentIcon}
          </div>
          
          {/* Agent Name */}
          <div>
            <h3 className="font-medium text-gray-900">{agent_name}</h3>
            <div className="text-sm text-gray-500">Finding #{index + 1}</div>
          </div>
        </div>

        {/* Severity Badge */}
        <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.badgeBg} ${config.badgeText}`}>
          <div className={`mr-1 ${config.color}`}>
            {config.icon}
          </div>
          {severity}
        </div>
      </div>

      {/* Component Information */}
      {componentInfo.hasComponent && (
        <div className="mb-4 p-3 bg-white rounded-md border border-gray-200">
          <div className="flex items-center space-x-2 text-sm">
            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <span className="font-medium text-gray-900">Component:</span>
            <span className="text-gray-700">{componentInfo.name}</span>
            {componentInfo.version && (
              <>
                <span className="text-gray-400">•</span>
                <span className="text-gray-600">v{componentInfo.version}</span>
              </>
            )}
          </div>
        </div>
      )}

      {/* Finding Description */}
      <div className="space-y-2">
        <h4 className="font-medium text-gray-900">Finding Details</h4>
        <p className="text-gray-700 leading-relaxed">
          {findingText}
        </p>
      </div>

      {/* Severity Indicator */}
      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center space-x-2">
            <div className={config.color}>
              {config.icon}
            </div>
            <span className={`font-medium ${config.color}`}>
              {severity} Severity
            </span>
          </div>
          
          {/* Action Recommendation */}
          <div className="text-gray-500">
            {severity === 'Critical' && 'Immediate action required'}
            {severity === 'High' && 'Address promptly'}
            {severity === 'Medium' && 'Plan to resolve'}
            {severity === 'Low' && 'Monitor and consider'}
          </div>
        </div>
      </div>
    </div>
  );
};