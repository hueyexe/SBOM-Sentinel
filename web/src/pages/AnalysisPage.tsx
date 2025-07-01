import React, { useState, useEffect } from 'react';
import { useParams, useLocation } from 'react-router-dom';
import { Card } from '../components/Card';
import { ResultsSummary } from '../components/ResultsSummary';
import { FindingCard } from '../components/FindingCard';
import ApiService from '../services/apiClient';
import type { AnalysisOptions, AnalysisResponse, Sbom } from '../types';

interface LocationState {
  analysisOptions?: AnalysisOptions;
}

export const AnalysisPage: React.FC = () => {
  const { sbomId } = useParams<{ sbomId: string }>();
  const location = useLocation();
  const state = location.state as LocationState;
  const analysisOptions = state?.analysisOptions || {
    enableAiHealth: false,
    enableProactiveScan: false,
  };

  const [analysisData, setAnalysisData] = useState<AnalysisResponse | null>(null);
  const [sbomData, setSbomData] = useState<Sbom | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAnalysisData = async () => {
      if (!sbomId) {
        setError('No SBOM ID provided');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);

        // Fetch both SBOM data and analysis results in parallel
        const [sbomResponse, analysisResponse] = await Promise.all([
          ApiService.getSbom(sbomId),
          ApiService.analyzeSbom(sbomId, analysisOptions)
        ]);

        setSbomData(sbomResponse);
        setAnalysisData(analysisResponse);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to fetch analysis data';
        setError(errorMessage);
        console.error('Analysis fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchAnalysisData();
  }, [sbomId, analysisOptions.enableAiHealth, analysisOptions.enableProactiveScan]);

  // Calculate total number of available agents
  const totalAgents = 1 + (analysisOptions.enableAiHealth ? 1 : 0) + (analysisOptions.enableProactiveScan ? 1 : 0);

  if (loading) {
    return (
      <div className="max-w-6xl mx-auto p-6 space-y-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">SBOM Analysis Results</h1>
          <p className="text-lg text-gray-600">
            Analysis for SBOM: {sbomId}
          </p>
        </div>

        {/* Loading State */}
        <Card title="Analysis Status" subtitle="Running analysis on your SBOM">
          <div className="space-y-4">
            <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
              <div className="flex items-center">
                <svg className="animate-spin h-5 w-5 text-blue-600 mr-3" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <div>
                  <p className="text-sm font-medium text-blue-800">Analysis in Progress</p>
                  <p className="text-sm text-blue-600">
                    Running {totalAgents} analysis agent{totalAgents > 1 ? 's' : ''} on your SBOM. This may take a few moments.
                  </p>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* License Analysis */}
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-2">
                  <h4 className="text-sm font-medium text-gray-900">License Agent</h4>
                  <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
                </div>
                <p className="text-xs text-gray-500">Checking license compliance</p>
              </div>

              {/* AI Health Analysis */}
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-2">
                  <h4 className="text-sm font-medium text-gray-900">AI Health Agent</h4>
                  <div className={`w-2 h-2 rounded-full ${analysisOptions.enableAiHealth ? 'bg-blue-500 animate-pulse' : 'bg-gray-300'}`}></div>
                </div>
                <p className="text-xs text-gray-500">
                  {analysisOptions.enableAiHealth ? 'Analyzing dependency health' : 'Disabled'}
                </p>
              </div>

              {/* Proactive Scan */}
              <div className="bg-white border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-2">
                  <h4 className="text-sm font-medium text-gray-900">Proactive Scan</h4>
                  <div className={`w-2 h-2 rounded-full ${analysisOptions.enableProactiveScan ? 'bg-blue-500 animate-pulse' : 'bg-gray-300'}`}></div>
                </div>
                <p className="text-xs text-gray-500">
                  {analysisOptions.enableProactiveScan ? 'Discovering emerging threats' : 'Disabled'}
                </p>
              </div>
            </div>
          </div>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto p-6 space-y-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">SBOM Analysis Results</h1>
          <p className="text-lg text-gray-600">
            Analysis for SBOM: {sbomId}
          </p>
        </div>

        {/* Error State */}
        <Card title="Analysis Error" subtitle="There was a problem analyzing your SBOM">
          <div className="text-center py-12">
            <svg className="h-16 w-16 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.728-.833-2.498 0L3.316 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Analysis Failed</h3>
            <p className="text-gray-500 mb-4">{error}</p>
            <button
              onClick={() => window.location.reload()}
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Retry Analysis
            </button>
          </div>
        </Card>
      </div>
    );
  }

  if (!analysisData || !sbomData) {
    return (
      <div className="max-w-6xl mx-auto p-6 space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">SBOM Analysis Results</h1>
          <p className="text-lg text-gray-600">No analysis data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">SBOM Analysis Results</h1>
        <p className="text-lg text-gray-600">
          Analysis for <span className="font-semibold">{sbomData.name}</span>
        </p>
        <p className="text-sm text-gray-500 mt-1">
          SBOM ID: {sbomId}
        </p>
      </div>

      {/* Results Summary */}
      <ResultsSummary summary={analysisData.summary} agentCount={totalAgents} />

      {/* SBOM Information */}
      <Card title="SBOM Information" subtitle="Details about the analyzed SBOM">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-3">
            <div>
              <span className="text-sm font-medium text-gray-900">Name:</span>
              <span className="ml-2 text-sm text-gray-700">{sbomData.name}</span>
            </div>
            <div>
              <span className="text-sm font-medium text-gray-900">Components:</span>
              <span className="ml-2 text-sm text-gray-700">{sbomData.components?.length || 0}</span>
            </div>
            <div>
              <span className="text-sm font-medium text-gray-900">Created:</span>
              <span className="ml-2 text-sm text-gray-700">
                {sbomData.created_at ? new Date(sbomData.created_at).toLocaleString() : 'Unknown'}
              </span>
            </div>
          </div>
          <div className="space-y-3">
            <div>
              <span className="text-sm font-medium text-gray-900">Analysis Options:</span>
            </div>
            <ul className="list-disc list-inside ml-4 text-sm text-gray-600 space-y-1">
              <li>License Compliance Analysis: <span className="text-green-600 font-medium">Enabled</span></li>
              <li>AI Dependency Health: <span className={`font-medium ${analysisOptions.enableAiHealth ? 'text-green-600' : 'text-gray-500'}`}>
                {analysisOptions.enableAiHealth ? 'Enabled' : 'Disabled'}
              </span></li>
              <li>Proactive Vulnerability Discovery: <span className={`font-medium ${analysisOptions.enableProactiveScan ? 'text-green-600' : 'text-gray-500'}`}>
                {analysisOptions.enableProactiveScan ? 'Enabled' : 'Disabled'}
              </span></li>
            </ul>
          </div>
        </div>
      </Card>

      {/* Analysis Findings */}
      {analysisData.results.length > 0 ? (
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold text-gray-900">Analysis Findings</h2>
            <div className="text-sm text-gray-500">
              {analysisData.results.length} finding{analysisData.results.length !== 1 ? 's' : ''}
            </div>
          </div>
          
          <div className="space-y-4">
            {analysisData.results.map((result, index) => (
              <FindingCard
                key={`${result.agent_name}-${index}`}
                finding={result}
                index={index}
              />
            ))}
          </div>
        </div>
      ) : (
        <Card title="Analysis Complete" subtitle="No issues found in your SBOM">
          <div className="text-center py-12">
            <svg className="w-16 h-16 text-green-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No Issues Found</h3>
            <p className="text-gray-500">
              Excellent! Your SBOM analysis completed successfully with no security issues or compliance violations detected.
            </p>
          </div>
        </Card>
      )}
    </div>
  );
};