import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { FileUpload } from '../components/FileUpload';
import { Button } from '../components/Button';
import { useFileUpload } from '../hooks/useFileUpload';
import ApiService from '../services/apiClient';
import type { AnalysisOptions } from '../types';

export const SubmitPage: React.FC = () => {
  const navigate = useNavigate();
  const { uploadState, setFile, setUploading, setError, validateFile } = useFileUpload();
  const [analysisOptions, setAnalysisOptions] = useState<AnalysisOptions>({
    enableAiHealth: false,
    enableProactiveScan: false,
    enableVulnScan: false,
  });
  const [submitSuccess, setSubmitSuccess] = useState<string | null>(null);

  const handleFileSelect = (file: File) => {
    const error = validateFile(file);
    if (error) {
      setError(error);
      return;
    }

    setFile(file);
    setSubmitSuccess(null);
  };

  const handleSubmit = async () => {
    if (!uploadState.file) {
      setError('Please select a file to upload');
      return;
    }

    setUploading(true);
    setError(undefined);

    try {
      // Submit the SBOM
      const response = await ApiService.submitSbom(uploadState.file);
      setSubmitSuccess(`SBOM uploaded successfully with ID: ${response.id}`);
      
      // Redirect to analysis page after a short delay
      setTimeout(() => {
        navigate(`/analysis/${encodeURIComponent(response.id)}`, {
          state: { analysisOptions }
        });
      }, 2000);

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Upload failed';
      setError(errorMessage);
    } finally {
      setUploading(false);
    }
  };

  const handleAnalysisOptionChange = (option: keyof AnalysisOptions, value: boolean) => {
    setAnalysisOptions(prev => ({
      ...prev,
      [option]: value,
    }));
  };

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Upload SBOM for Analysis</h1>
        <p className="text-lg text-gray-600">
          Submit your Software Bill of Materials for comprehensive security analysis
        </p>
      </div>

      {/* Upload Section */}
      <Card title="SBOM File Upload" subtitle="Select your SBOM file to begin analysis">
        <div className="space-y-6">
          <FileUpload
            onFileSelect={handleFileSelect}
            error={uploadState.error}
            loading={uploadState.isUploading}
            progress={uploadState.uploadProgress}
          />

          {uploadState.file && !uploadState.error && (
            <div className="bg-green-50 border border-green-200 rounded-md p-4">
              <div className="flex items-center">
                <svg className="w-5 h-5 text-green-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div>
                  <p className="text-sm font-medium text-green-800">
                    {uploadState.file.name}
                  </p>
                  <p className="text-sm text-green-600">
                    {(uploadState.file.size / (1024 * 1024)).toFixed(2)} MB
                  </p>
                </div>
              </div>
            </div>
          )}

          {submitSuccess && (
            <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
              <div className="flex items-center">
                <svg className="w-5 h-5 text-blue-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p className="text-sm text-blue-800">{submitSuccess}</p>
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Analysis Options */}
      <Card title="Analysis Options" subtitle="Configure which analysis agents to run">
        <div className="space-y-4">
          <div className="bg-gray-50 rounded-lg p-4">
            <h4 className="text-sm font-medium text-gray-900 mb-3">Available Analysis Agents</h4>
            <div className="space-y-3">
              {/* License Analysis - Always enabled */}
              <div className="flex items-center justify-between p-3 bg-white rounded-md border">
                <div className="flex items-center">
                  <div className="w-2 h-2 bg-green-500 rounded-full mr-3"></div>
                  <div>
                    <p className="text-sm font-medium text-gray-900">License Compliance Agent</p>
                    <p className="text-xs text-gray-500">Detects high-risk copyleft licenses</p>
                  </div>
                </div>
                <span className="text-xs text-green-600 font-medium">Always Enabled</span>
              </div>

              {/* AI Health Check */}
              <div className="flex items-center justify-between p-3 bg-white rounded-md border">
                <div className="flex items-center">
                  <div className={`w-2 h-2 rounded-full mr-3 ${analysisOptions.enableAiHealth ? 'bg-blue-500' : 'bg-gray-300'}`}></div>
                  <div>
                    <p className="text-sm font-medium text-gray-900">AI Dependency Health Agent</p>
                    <p className="text-xs text-gray-500">Uses local LLM to assess project health</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={analysisOptions.enableAiHealth}
                    onChange={(e) => handleAnalysisOptionChange('enableAiHealth', e.target.checked)}
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>

              {/* Proactive Scan */}
              <div className="flex items-center justify-between p-3 bg-white rounded-md border">
                <div className="flex items-center">
                  <div className={`w-2 h-2 rounded-full mr-3 ${analysisOptions.enableProactiveScan ? 'bg-blue-500' : 'bg-gray-300'}`}></div>
                  <div>
                    <p className="text-sm font-medium text-gray-900">Proactive Vulnerability Discovery</p>
                    <p className="text-xs text-gray-500">RAG-powered early threat detection</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={analysisOptions.enableProactiveScan}
                    onChange={(e) => handleAnalysisOptionChange('enableProactiveScan', e.target.checked)}
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>

              {/* Known Vulnerability Scanning */}
              <div className="flex items-center justify-between p-3 bg-white rounded-md border">
                <div className="flex items-center">
                  <div className={`w-2 h-2 rounded-full mr-3 ${analysisOptions.enableVulnScan ? 'bg-blue-500' : 'bg-gray-300'}`}></div>
                  <div>
                    <p className="text-sm font-medium text-gray-900">Known Vulnerability Scanning (CVE)</p>
                    <p className="text-xs text-gray-500">Checks against OSV.dev vulnerability database</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={analysisOptions.enableVulnScan}
                    onChange={(e) => handleAnalysisOptionChange('enableVulnScan', e.target.checked)}
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
          </div>

          {(analysisOptions.enableAiHealth || analysisOptions.enableProactiveScan) && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4">
              <div className="flex items-start">
                <svg className="w-5 h-5 text-yellow-400 mr-2 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.728-.833-2.498 0L3.316 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
                <div>
                  <p className="text-sm font-medium text-yellow-800">AI Features Enabled</p>
                  <p className="text-sm text-yellow-700">
                    Make sure Ollama is running locally for AI-powered analysis features to work properly.
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Submit Button */}
      <div className="flex justify-center">
        <Button
          onClick={handleSubmit}
          loading={uploadState.isUploading}
          disabled={!uploadState.file || !!uploadState.error}
          size="lg"
          className="px-12"
        >
          {uploadState.isUploading ? 'Uploading...' : 'Submit SBOM for Analysis'}
        </Button>
      </div>
    </div>
  );
};