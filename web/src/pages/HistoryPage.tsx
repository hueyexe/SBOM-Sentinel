import React from 'react';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Link } from 'react-router-dom';

export const HistoryPage: React.FC = () => {
  return (
    <div className="max-w-6xl mx-auto p-6 space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Analysis History</h1>
        <p className="text-lg text-gray-600">
          Review your previous SBOM analysis results and track security trends
        </p>
      </div>

      {/* Empty State */}
      <Card>
        <div className="text-center py-16">
          <svg
            className="h-24 w-24 text-gray-300 mx-auto mb-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1}
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
            />
          </svg>
          
          <h3 className="text-xl font-medium text-gray-900 mb-4">No Analysis History Yet</h3>
          
          <p className="text-gray-500 mb-8 max-w-md mx-auto">
            Once you upload and analyze SBOM files, your analysis history will appear here. 
            You'll be able to view detailed results, track findings over time, and compare different analyses.
          </p>

          <div className="space-y-4">
            <Link to="/submit">
              <Button variant="primary" size="lg">
                Upload Your First SBOM
              </Button>
            </Link>
            
            <div className="text-sm text-gray-500">
              Start by uploading an SBOM file to see analysis results
            </div>
          </div>
        </div>
      </Card>

      {/* Features Preview */}
      <Card title="What You'll See Here" subtitle="Preview of upcoming history features">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div className="text-center p-4">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-3">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h4 className="font-semibold text-gray-900 mb-2">Analysis Timeline</h4>
            <p className="text-sm text-gray-600">
              Chronological view of all your SBOM analyses with timestamps and status
            </p>
          </div>

          <div className="text-center p-4">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-3">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <h4 className="font-semibold text-gray-900 mb-2">Findings Summary</h4>
            <p className="text-sm text-gray-600">
              Quick overview of security findings and compliance issues discovered
            </p>
          </div>

          <div className="text-center p-4">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-3">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
            <h4 className="font-semibold text-gray-900 mb-2">Trend Analysis</h4>
            <p className="text-sm text-gray-600">
              Track how your security posture changes over time with visual trends
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
};