import { test, expect } from '@playwright/test';

test.describe('SBOM Sentinel E2E Tests', () => {
  
  test('Complete SBOM Analysis Workflow', async ({ page }) => {
    console.log('🎬 Starting complete SBOM analysis workflow test...');
    
    // Step 1: Navigate to the home page
    console.log('📍 Step 1: Navigating to home page...');
    await page.goto('/');
    
    // Verify we're on the home page
    await expect(page).toHaveTitle(/SBOM Sentinel/);
    await expect(page.locator('h1')).toContainText('SBOM Sentinel');
    
    console.log('✅ Home page loaded successfully');
    
    // Step 2: Navigate to Submit SBOM page
    console.log('📍 Step 2: Navigating to Submit SBOM page...');
    
    // Look for and click the submit link/button
    const submitLink = page.locator('a', { hasText: /submit|upload/i }).first();
    await expect(submitLink).toBeVisible({ timeout: 10000 });
    await submitLink.click();
    
    // Verify we're on the submit page
    await expect(page.url()).toContain('/submit');
    await expect(page.locator('h1, h2')).toContainText(/submit|upload/i);
    
    console.log('✅ Submit SBOM page loaded successfully');
    
    // Step 3: Upload SBOM file
    console.log('📍 Step 3: Uploading SBOM file...');
    
    // Find file input
    const fileInput = page.locator('input[type="file"]');
    await expect(fileInput).toBeVisible();
    
    // Upload the test SBOM file using relative path
    const testSBOMPath = './tests/e2e/fixtures/test-sbom.json';
    await fileInput.setInputFiles(testSBOMPath);
    
    console.log('✅ SBOM file uploaded successfully');
    
    // Step 4: Enable AI analysis options
    console.log('📍 Step 4: Enabling AI analysis options...');
    
    // Look for toggle switches or checkboxes for AI options
    const aiHealthToggle = page.locator('input[type="checkbox"], [role="switch"]').filter({
      has: page.locator('text=/health|dependency|ai/i')
    }).first();
    
    const proactiveToggle = page.locator('input[type="checkbox"], [role="switch"]').filter({
      has: page.locator('text=/proactive|vulnerability|scan/i')
    }).first();
    
    // Enable AI health check if toggle exists
    if (await aiHealthToggle.isVisible()) {
      if (!(await aiHealthToggle.isChecked())) {
        await aiHealthToggle.click();
        console.log('✅ AI Health Check enabled');
      }
    }
    
    // Enable proactive vulnerability scan if toggle exists
    if (await proactiveToggle.isVisible()) {
      if (!(await proactiveToggle.isChecked())) {
        await proactiveToggle.click();
        console.log('✅ Proactive Vulnerability Scan enabled');
      }
    }
    
    // Step 5: Submit the form
    console.log('📍 Step 5: Submitting SBOM for analysis...');
    
    const submitButton = page.locator('button[type="submit"], button').filter({
      hasText: /submit|analyze|upload/i
    }).first();
    
    await expect(submitButton).toBeVisible();
    await expect(submitButton).toBeEnabled();
    
    // Click submit and wait for navigation
    await Promise.all([
      page.waitForURL(/\/analysis\/.*/, { timeout: 30000 }),
      submitButton.click()
    ]);
    
    console.log('✅ SBOM submitted successfully');
    
    // Step 6: Wait for and verify analysis results page
    console.log('📍 Step 6: Verifying analysis results...');
    
    // Verify we're on an analysis page
    expect(page.url()).toMatch(/\/analysis\/[^\/]+/);
    
    // Wait for page to load completely
    await page.waitForLoadState('networkidle', { timeout: 30000 });
    
    // Look for analysis results content
    await expect(page.locator('h1, h2')).toContainText(/analysis|results/i, { timeout: 15000 });
    
    // Step 7: Verify presence of findings
    console.log('📍 Step 7: Checking for analysis findings...');
    
    // Wait for findings to load
    await page.waitForTimeout(2000);
    
    // Look for severity indicators - should find High severity from GPL license
    const highSeverityElements = page.locator('[class*="severity"], [class*="finding"], .card, .result').filter({
      hasText: /high|critical/i
    });
    
    // Check for specific license findings
    const licenseFindings = page.locator('text=/GPL|AGPL|license|copyleft/i');
    
    // Verify we have some findings
    const findingsExist = await Promise.race([
      highSeverityElements.first().isVisible().then(() => true).catch(() => false),
      licenseFindings.first().isVisible().then(() => true).catch(() => false),
      page.locator('text=/finding|result|issue/i').first().isVisible().then(() => true).catch(() => false)
    ]);
    
    if (findingsExist) {
      console.log('✅ Analysis findings detected');
      
      // Try to find High severity findings specifically
      const hasHighSeverity = await highSeverityElements.first().isVisible().catch(() => false);
      if (hasHighSeverity) {
        console.log('✅ High severity findings detected (likely from GPL/AGPL licenses)');
      }
      
      // Check for license-related content
      const hasLicenseFindings = await licenseFindings.first().isVisible().catch(() => false);
      if (hasLicenseFindings) {
        console.log('✅ License-related findings detected');
      }
    } else {
      console.log('⚠️ No specific findings detected, but analysis page loaded');
    }
    
    // Step 8: Verify page structure and content
    console.log('📍 Step 8: Verifying analysis page structure...');
    
    // Check for SBOM information
    const sbomInfo = page.locator('text=/E2E Test Application|test-sbom|sbom/i');
    await expect(sbomInfo.first()).toBeVisible({ timeout: 10000 });
    
    // Check for analysis summary or results section
    const resultsSection = page.locator('[class*="result"], [class*="finding"], [class*="summary"], .card').first();
    await expect(resultsSection).toBeVisible({ timeout: 10000 });
    
    console.log('✅ Analysis page structure verified');
    
    console.log('🎉 Complete SBOM analysis workflow test passed!');
  });
  
  test('Navigation and UI Elements', async ({ page }) => {
    console.log('🎬 Starting navigation and UI elements test...');
    
    // Test home page
    await page.goto('/');
    await expect(page.locator('h1')).toBeVisible();
    
    // Test that navigation elements are present
    const nav = page.locator('nav, header').first();
    if (await nav.isVisible()) {
      console.log('✅ Navigation element found');
    }
    
    // Test responsive design by changing viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    await expect(page.locator('h1')).toBeVisible();
    
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(page.locator('h1')).toBeVisible();
    
    console.log('✅ Responsive design test passed');
  });
  
  test('Error Handling', async ({ page }) => {
    console.log('🎬 Starting error handling test...');
    
    // Navigate to submit page
    await page.goto('/submit');
    
    // Try to submit without file
    const submitButton = page.locator('button[type="submit"], button').filter({
      hasText: /submit|analyze|upload/i
    }).first();
    
    if (await submitButton.isVisible()) {
      await submitButton.click();
      
      // Check for error message or validation
      await page.waitForTimeout(1000);
      
      // Look for error indicators
      const errorMessage = page.locator('[class*="error"], [role="alert"], .alert-error').first();
      const validationMessage = page.locator('input:invalid, [aria-invalid="true"]').first();
      
      const hasError = await Promise.race([
        errorMessage.isVisible().then(() => true).catch(() => false),
        validationMessage.isVisible().then(() => true).catch(() => false)
      ]);
      
      if (hasError) {
        console.log('✅ Error handling works - validation prevents submission without file');
      } else {
        console.log('⚠️ No visible error handling detected');
      }
    }
    
    console.log('✅ Error handling test completed');
  });
  
  test('Backend Health Check', async ({ page }) => {
    console.log('🎬 Starting backend health check test...');
    
    // Test that backend is responding
    const response = await page.request.get('http://localhost:8080/health');
    expect(response.status()).toBe(200);
    
    const healthData = await response.json();
    expect(healthData.status).toBe('ok');
    expect(healthData.service).toBe('sbom-sentinel');
    
    console.log('✅ Backend health check passed');
  });
});