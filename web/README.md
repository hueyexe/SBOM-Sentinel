# SBOM Sentinel Web Dashboard

A modern React TypeScript web interface for the SBOM Sentinel platform, providing an intuitive way to upload, analyze, and manage Software Bill of Materials (SBOM) files.

## 🚀 Features

- **Drag & Drop File Upload**: Modern file upload interface with drag-and-drop support
- **Real-time Analysis**: Live progress tracking for SBOM analysis operations
- **AI-Powered Options**: Configure AI health checks and proactive vulnerability discovery
- **Responsive Design**: Professional UI that works on desktop and mobile devices
- **Type-Safe**: Full TypeScript implementation with comprehensive type definitions
- **Modern Architecture**: Clean component structure following React best practices

## 🏗️ Architecture

The web dashboard follows a professional React application structure:

```
src/
├── components/          # Reusable UI components
│   ├── Button.tsx       # Custom button component with variants
│   ├── Card.tsx         # Container component for content sections
│   └── FileUpload.tsx   # Advanced file upload with drag-and-drop
├── hooks/               # Custom React hooks
│   └── useFileUpload.ts # File upload state management
├── layouts/             # Application layout components
│   └── MainLayout.tsx   # Main app layout with sidebar navigation
├── pages/               # Top-level page components
│   ├── DashboardPage.tsx    # Welcome/overview page
│   ├── SubmitPage.tsx       # SBOM upload and submission
│   ├── AnalysisPage.tsx     # Analysis results display
│   └── HistoryPage.tsx      # Analysis history (placeholder)
├── services/            # External API communication
│   └── apiClient.ts     # Backend API integration
└── types/               # TypeScript type definitions
    └── index.ts         # Shared interface definitions
```

## 🛠️ Technology Stack

- **React 18** - Modern React with hooks and functional components
- **TypeScript** - Full type safety and excellent developer experience
- **Vite** - Fast build tool and development server
- **Tailwind CSS** - Utility-first CSS framework for rapid styling
- **React Router** - Client-side routing for single-page application
- **Axios** - HTTP client for API communication

## 🚦 Getting Started

### Prerequisites

- Node.js 18+ and npm
- SBOM Sentinel backend server running on `http://localhost:8080`

### Installation

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Start the development server:**
   ```bash
   npm run dev
   ```

3. **Open your browser:**
   Navigate to `http://localhost:5173` to view the dashboard.

### Build for Production

```bash
# Build the application
npm run build

# Preview the production build
npm run preview
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the web directory to configure the API endpoint:

```env
VITE_API_URL=http://localhost:8080
```

### Backend Integration

The dashboard connects to the SBOM Sentinel backend via REST API:

- **Submit SBOM**: `POST /api/v1/sboms`
- **Get SBOM**: `GET /api/v1/sboms/get?id={id}`
- **Analyze SBOM**: `POST /api/v1/sboms/{id}/analyze`
- **Health Check**: `GET /health`

## 📱 User Interface

### Dashboard Page (`/`)
- Welcome screen with platform overview
- Feature highlights and quick actions
- Navigation to upload and history sections

### Submit Page (`/submit`)
- Drag-and-drop file upload interface
- Analysis configuration options:
  - License Compliance (always enabled)
  - AI Dependency Health Check (optional)
  - Proactive Vulnerability Discovery (optional)
- Real-time upload progress and feedback

### Analysis Page (`/analysis/{id}`)
- Analysis progress tracking
- SBOM information display
- Results placeholder (ready for future expansion)

### History Page (`/history`)
- Placeholder for analysis history
- Empty state with guided next steps

## 🎨 Design System

### Components

#### Button
```tsx
<Button variant="primary" size="lg" loading={isSubmitting}>
  Submit SBOM
</Button>
```

Variants: `primary`, `secondary`, `danger`
Sizes: `sm`, `md`, `lg`

#### Card
```tsx
<Card title="Upload SBOM" subtitle="Select your file">
  <FileUpload onFileSelect={handleFile} />
</Card>
```

#### FileUpload
```tsx
<FileUpload
  onFileSelect={handleFile}
  accept=".json,.xml,.spdx"
  maxSize={32 * 1024 * 1024}
  loading={uploading}
  error={error}
/>
```

### Styling

The application uses Tailwind CSS for consistent, responsive design:

- **Color Scheme**: Professional blue and gray palette
- **Typography**: Inter font family for excellent readability
- **Spacing**: Consistent spacing scale using Tailwind classes
- **Responsive**: Mobile-first responsive design approach

## 🔌 API Integration

### Type Safety

All API interactions are fully typed:

```typescript
interface SubmitSbomResponse {
  id: string;
  message: string;
}

interface AnalysisOptions {
  enableAiHealth: boolean;
  enableProactiveScan: boolean;
}
```

### Error Handling

Comprehensive error handling with user-friendly messages:

```typescript
try {
  const response = await ApiService.submitSbom(file);
  // Handle success
} catch (error) {
  const message = error instanceof Error ? error.message : 'Upload failed';
  setError(message);
}
```

## 🚀 Development

### File Upload Flow

1. User drags file or clicks to browse
2. File validation (size, type)
3. Upload with progress tracking
4. Success feedback with SBOM ID
5. Automatic redirect to analysis page

### State Management

Uses React hooks for local state management:

- `useFileUpload` - File upload state and validation
- `useState` - Component-level state
- `useNavigate` - Programmatic navigation

### Custom Hooks

#### useFileUpload
```typescript
const {
  uploadState,
  setFile,
  setUploading,
  setError,
  validateFile,
  reset
} = useFileUpload();
```

## 🔄 Future Enhancements

### Planned Features

1. **Real Analysis Integration**
   - Live analysis results display
   - Progress tracking for each analysis agent
   - Detailed findings visualization

2. **Analysis History**
   - SBOM analysis history management
   - Search and filtering capabilities
   - Comparison between analyses

3. **Enhanced Visualizations**
   - Component dependency graphs
   - Security findings charts
   - License compliance summaries

4. **User Experience**
   - Dark mode support
   - Export functionality
   - Advanced filtering and search

### Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the existing code style and patterns
4. Add TypeScript types for all new features
5. Test thoroughly before submitting PR

## 📝 Notes

- The dashboard is designed to work seamlessly with the SBOM Sentinel Go backend
- All API calls include proper error handling and loading states
- The UI is optimized for modern browsers with ES6+ support
- Responsive design ensures functionality across device sizes

---

**Built with ❤️ for supply chain security**
