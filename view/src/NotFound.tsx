import { 
  FileQuestion,
  Home
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const NotFound = () => {
  const navigate = useNavigate();
  return (
    <div className="w-full min-h-screen flex flex-col items-center justify-center p-6 text-center animate-in fade-in duration-700">
      <div className="relative mb-12">
        <h1 className="text-[12rem] md:text-[18rem] font-black text-slate-100 leading-none select-none">404</h1>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <FileQuestion className="text-blue-500 w-24 h-24 mb-6 animate-bounce duration-[3000ms]" />
          <h2 className="text-3xl font-bold text-slate-800">页面不存在</h2>
          <button onClick={() => navigate('/')} className="mt-8 px-10 py-4 bg-slate-900 text-white font-bold rounded-2xl shadow-xl flex items-center space-x-2 active:scale-95 transition-all">
            <Home className="w-5 h-5" /> <span>回到首页</span>
          </button>
        </div>
      </div>
    </div>
  );
};
export default NotFound