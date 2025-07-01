// SBOM Types
export interface SbomComponent {
  name: string;
  version: string;
  license?: string;
  purl?: string;
  type?: string;
}

export interface SbomMetadata {
  name: string;
  version?: string;
  description?: string;
  author?: string;
  [key: string]: any;
}

export interface Sbom {
  id: string;
  name: string;
  components: SbomComponent[];
  metadata: SbomMetadata;
  created_at?: string;
  updated_at?: string;
}

// API Response Types
export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  message?: string;
}

export interface SubmitSbomResponse {
  id: string;
  message: string;
}

export interface AnalysisResult {
  agent_name: string;
  finding: string;
  severity: 'Low' | 'Medium' | 'High' | 'Critical';
}

export interface AnalysisResponse {
  sbom_id: string;
  results: AnalysisResult[];
  summary: {
    total_findings: number;
    findings_by_severity: Record<string, number>;
    agents_run: string[];
  };
}

export interface AnalysisOptions {
  enableAiHealth: boolean;
  enableProactiveScan: boolean;
  enableVulnScan: boolean;
}

// File Upload Types
export interface FileUploadState {
  file: File | null;
  isUploading: boolean;
  uploadProgress: number;
  error?: string;
}