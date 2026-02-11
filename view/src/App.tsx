import { 
    BrowserRouter as Router, 
    Routes, 
    Route, 
} from 'react-router-dom';

import { Background, Footer } from './Widgets';
import Login from './Login';
import Authorize from './Authorize';
import NotFound from './NotFound';

const AppContent = () => {
  return (
    <div className="min-h-screen bg-white sm:bg-slate-50 font-sans text-slate-900 selection:bg-blue-100 flex items-center justify-center sm:p-4">
      <Background />
      <div className="w-full flex justify-center">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/authorize" element={<Authorize />} />
          <Route path="*" element={<NotFound />} /> {/* Default Fallback */}
        </Routes>
      </div>
      <Footer className="hidden sm:block fixed bottom-8 w-full text-center" />
    </div>
  );
};

const App = () => (
  <Router>
    <AppContent />
  </Router>
);

export default App;