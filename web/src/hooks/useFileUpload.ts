import { useState, useCallback } from 'react';
import type { FileUploadState } from '../types';

export const useFileUpload = () => {
  const [uploadState, setUploadState] = useState<FileUploadState>({
    file: null,
    isUploading: false,
    uploadProgress: 0,
    error: undefined,
  });

  const setFile = useCallback((file: File | null) => {
    setUploadState(prev => ({
      ...prev,
      file,
      error: undefined,
    }));
  }, []);

  const setUploading = useCallback((isUploading: boolean) => {
    setUploadState(prev => ({
      ...prev,
      isUploading,
      uploadProgress: isUploading ? 0 : prev.uploadProgress,
    }));
  }, []);

  const setProgress = useCallback((uploadProgress: number) => {
    setUploadState(prev => ({
      ...prev,
      uploadProgress: Math.min(100, Math.max(0, uploadProgress)),
    }));
  }, []);

  const setError = useCallback((error: string | undefined) => {
    setUploadState(prev => ({
      ...prev,
      error,
      isUploading: false,
    }));
  }, []);

  const reset = useCallback(() => {
    setUploadState({
      file: null,
      isUploading: false,
      uploadProgress: 0,
      error: undefined,
    });
  }, []);

  const validateFile = useCallback((file: File): string | null => {
    // Check file size (max 32MB to match backend)
    const maxSize = 32 * 1024 * 1024; // 32MB
    if (file.size > maxSize) {
      return 'File size exceeds 32MB limit';
    }

    // Check file extension
    const allowedExtensions = ['.json', '.xml', '.spdx'];
    const fileName = file.name.toLowerCase();
    const hasValidExtension = allowedExtensions.some(ext => fileName.endsWith(ext));
    
    if (!hasValidExtension) {
      return 'Invalid file type. Please upload a JSON, XML, or SPDX file';
    }

    return null;
  }, []);

  return {
    uploadState,
    setFile,
    setUploading,
    setProgress,
    setError,
    reset,
    validateFile,
  };
};