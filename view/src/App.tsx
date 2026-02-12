import { 
    HashRouter as Router, 
    Routes, 
    Route,
    Navigate,
    useLocation, 
} from 'react-router-dom';

import { Background, Footer } from './Widgets';
import Login from './Login';
import Approve from './Approve';
import NotFound from './NotFound';
import Register from './Register';
import Profile from './Profile';
import { useState } from 'react';
import Root from './Root';

const Protected = ({ children, isAuthenticated }: { children: React.ReactNode, isAuthenticated: boolean }) => {
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return children;
};
// --- 主应用路由配置 ---
const AppContent = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const { pathname } = useLocation();

  // Root, Profile and NotFound are Fullscreen, others are centered Card
  const isFullscreenLayout = pathname === '/' || pathname === '/profile' || !['/', '/login', '/register', '/approve', '/profile'].includes(pathname);

  return (
    <div className="min-h-screen bg-white sm:bg-slate-50 font-sans text-slate-900 selection:bg-blue-100 flex flex-col">
      <Background />
      
      <div className={`flex-1 flex flex-col ${isFullscreenLayout ? '' : 'items-center justify-center sm:p-4'}`}>
        <Routes>
          <Route path="/" element={<Root />} />
          <Route path="/login" element={<Login onLogin={() => setIsAuthenticated(true)} />} />
          <Route path="/register" element={<Register />} />
          <Route path="/approve" element={<Protected isAuthenticated={isAuthenticated}><Approve /></Protected>} />
          <Route path="/profile" element={<Protected isAuthenticated={isAuthenticated}><Profile onLogout={() => setIsAuthenticated(false)} /></Protected>} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </div>

      {/* Conditional footer: Only on card layouts */}
      {!isFullscreenLayout && <Footer className="hidden sm:block py-8 w-full" />}
    </div>
  );
};


const App = () => (
  <Router>
    <AppContent />
  </Router>
);

export default App;