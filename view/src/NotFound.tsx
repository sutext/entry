import { 
  FileQuestion,
  Home
} from 'lucide-react';
import { cardBaseStyles, Footer } from './Widgets';
import { useNavigate } from 'react-router-dom';
const NotFound = () => {
  const navigate = useNavigate();

  return (
    <div className={`${cardBaseStyles} animate-in fade-in zoom-in-95 duration-500 flex flex-col items-center text-center`}>
      <div className="w-20 h-20 bg-slate-50 rounded-3xl flex items-center justify-center mb-6 relative">
        <FileQuestion className="text-slate-400 w-10 h-10" />
        <div className="absolute -top-1 -right-1 w-6 h-6 bg-red-400 rounded-full border-4 border-white flex items-center justify-center text-[10px] text-white font-bold">!</div>
      </div>
      
      <h1 className="text-4xl font-black text-slate-200 mb-2">404</h1>
      <h2 className="text-xl font-bold text-slate-800 mb-3">页面未找到</h2>
      <p className="text-slate-500 text-sm mb-8 max-w-[240px]">
        抱歉，您访问的路径似乎消失在数字黑洞中了。
      </p>

      <button 
        onClick={() => navigate('/login')}
        className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3.5 rounded-xl shadow-md shadow-blue-100 transition-all flex items-center justify-center space-x-2 active:scale-[0.98]"
      >
        <Home className="w-4 h-4" />
        <span>回到首页</span>
      </button>

      <Footer className="mt-10 sm:hidden opacity-60" />
    </div>
  );
};
export default NotFound