import React from 'react';
import type { AnalysisResponse } from '../types';

interface ResultsSummaryProps {
  summary: AnalysisResponse['summary'];
  agentCount: number;
}

const severityConfig = {
  Critical: {
    color: 'text-red-700',
    bg: 'bg-red-100',
    border: 'border-red-200',
    icon: (
      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  High: {
    color: 'text-orange-700',
    bg: 'bg-orange-100',
    border: 'border-orange-200',
    icon: (
      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  Medium: {
    color: 'text-yellow-700',
    bg: 'bg-yellow-100',
    border: 'border-yellow-200',
    icon: (
      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
      </svg>
    ),
  },
  Low: {
    color: 'text-blue-700',
    bg: 'bg-blue-100',
    border: 'border-blue-200',
    icon: (
      <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
      </svg>
    ),
  },
};

export const ResultsSummary: React.FC<ResultsSummaryProps> = ({ summary, agentCount }) => {
  const { total_findings, findings_by_severity, agents_run } = summary;

  // Get severity breakdown in priority order
  const severityOrder = ['Critical', 'High', 'Medium', 'Low'];
  const severityBreakdown = severityOrder.map(severity => ({
    severity,
    count: findings_by_severity[severity] || 0,
    config: severityConfig[severity as keyof typeof severityConfig],
  })).filter(item => item.count > 0);

  const hasFindings = total_findings > 0;

  return (
    <div className="bg-white rounded-lg shadow-md border border-gray-200 p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">Analysis Summary</h2>
        <div className="flex items-center space-x-2">
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <span className="text-sm text-gray-600">Analysis Complete</span>
        </div>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        {/* Total Findings */}
        <div className="bg-gray-50 rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-gray-900 mb-1">{total_findings}</div>
          <div className="text-sm text-gray-600">Total Findings</div>
        </div>

        {/* Agents Run */}
        <div className="bg-gray-50 rounded-lg p-4 text-center">
          <div className="text-2xl font-bold text-gray-900 mb-1">{agents_run.length}</div>
          <div className="text-sm text-gray-600">Agents Executed</div>
        </div>

        {/* Analysis Status */}
        <div className="bg-gray-50 rounded-lg p-4 text-center">
          <div className={`text-2xl font-bold mb-1 ${hasFindings ? 'text-red-600' : 'text-green-600'}`}>
            {hasFindings ? 'Issues Found' : 'Clean'}
          </div>
          <div className="text-sm text-gray-600">Overall Status</div>
        </div>
      </div>

      {/* Severity Breakdown */}
      {hasFindings && (
        <div className="mb-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Findings by Severity</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {severityBreakdown.map(({ severity, count, config }) => (
              <div
                key={severity}
                className={`rounded-lg border p-4 ${config.bg} ${config.border}`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <div className={config.color}>
                      {config.icon}
                    </div>
                    <span className={`font-medium ${config.color}`}>
                      {severity}
                    </span>
                  </div>
                  <span className={`text-xl font-bold ${config.color}`}>
                    {count}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Agents Information */}
      <div>
        <h3 className="text-lg font-medium text-gray-900 mb-4">Analysis Agents</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {agents_run.map((agent, index) => (
            <div key={index} className="bg-green-50 border border-green-200 rounded-lg p-3">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-sm font-medium text-green-800">{agent}</span>
              </div>
            </div>
          ))}
        </div>
        
        {agentCount > agents_run.length && (
          <div className="mt-2 text-sm text-gray-500">
            {agentCount - agents_run.length} additional agent(s) were disabled for this analysis
          </div>
        )}
      </div>

      {/* No Findings State */}
      {!hasFindings && (
        <div className="text-center py-8">
          <svg className="w-16 h-16 text-green-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Issues Found</h3>
          <p className="text-gray-500">
            Great! Your SBOM analysis completed successfully with no security issues or compliance violations detected.
          </p>
        </div>
      )}
    </div>
  );
};