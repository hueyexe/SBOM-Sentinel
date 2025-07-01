import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MainLayout } from './layouts/MainLayout';
import { DashboardPage } from './pages/DashboardPage';
import { SubmitPage } from './pages/SubmitPage';
import { AnalysisPage } from './pages/AnalysisPage';
import { HistoryPage } from './pages/HistoryPage';

function App() {
  return (
    <Router>
      <MainLayout>
        <Routes>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/submit" element={<SubmitPage />} />
          <Route path="/analysis/:sbomId" element={<AnalysisPage />} />
          <Route path="/history" element={<HistoryPage />} />
        </Routes>
      </MainLayout>
    </Router>
  );
}

export default App;
