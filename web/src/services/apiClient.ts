import axios, { type AxiosResponse } from 'axios';
import type { Sbom, SubmitSbomResponse, AnalysisResponse, AnalysisOptions } from '../types';

// Configure axios defaults
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.data?.error) {
      throw new Error(error.response.data.message || error.response.data.error);
    }
    throw error;
  }
);

export class ApiService {
  /**
   * Submit an SBOM file for storage
   */
  static async submitSbom(file: File): Promise<SubmitSbomResponse> {
    const formData = new FormData();
    formData.append('sbom', file);

    const response: AxiosResponse<SubmitSbomResponse> = await apiClient.post(
      '/api/v1/sboms',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );

    return response.data;
  }

  /**
   * Retrieve an SBOM by ID
   */
  static async getSbom(id: string): Promise<Sbom> {
    const response: AxiosResponse<Sbom> = await apiClient.get(
      `/api/v1/sboms/get?id=${encodeURIComponent(id)}`
    );

    return response.data;
  }

  /**
   * Analyze an SBOM with optional AI features
   */
  static async analyzeSbom(
    id: string, 
    options: AnalysisOptions = { enableAiHealth: false, enableProactiveScan: false }
  ): Promise<AnalysisResponse> {
    const params = new URLSearchParams();
    
    if (options.enableAiHealth) {
      params.append('enable-ai-health-check', 'true');
    }
    
    if (options.enableProactiveScan) {
      params.append('enable-proactive-scan', 'true');
    }

    const url = `/api/v1/sboms/${encodeURIComponent(id)}/analyze${params.toString() ? '?' + params.toString() : ''}`;
    
    const response: AxiosResponse<AnalysisResponse> = await apiClient.post(url);

    return response.data;
  }

  /**
   * Health check endpoint
   */
  static async healthCheck(): Promise<{ status: string }> {
    const response: AxiosResponse<{ status: string }> = await apiClient.get('/health');
    return response.data;
  }

  /**
   * Submit and analyze SBOM in one step (convenience method)
   */
  static async submitAndAnalyzeSbom(
    file: File,
    options: AnalysisOptions = { enableAiHealth: false, enableProactiveScan: false }
  ): Promise<{ sbomId: string; analysisResult: AnalysisResponse }> {
    // First submit the SBOM
    const submitResponse = await this.submitSbom(file);
    
    // Then analyze it
    const analysisResult = await this.analyzeSbom(submitResponse.id, options);
    
    return {
      sbomId: submitResponse.id,
      analysisResult
    };
  }
}

export default ApiService;