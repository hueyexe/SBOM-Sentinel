import { FullConfig } from '@playwright/test';

export default async function globalSetup(config: FullConfig) {
  console.log('🔧 Starting global E2E test setup...');

  // Wait for backend to be ready (assume it's started externally)
  await waitForServer('http://localhost:8080/health', 30000);

  console.log('✅ Global E2E test setup complete!');
}

async function waitForServer(url: string, timeout: number): Promise<void> {
  console.log(`🕐 Waiting for server at ${url}...`);
  
  const startTime = Date.now();
  
  while (Date.now() - startTime < timeout) {
    try {
      const response = await fetch(url);
      if (response.ok) {
        console.log(`✅ Server is ready at ${url}`);
        return;
      }
    } catch (error) {
      // Server not ready yet, continue waiting
    }
    
    // Wait 1 second before next attempt
    await new Promise(resolve => setTimeout(resolve, 1000));
  }
  
  throw new Error(`Server at ${url} did not become ready within ${timeout}ms`);
}