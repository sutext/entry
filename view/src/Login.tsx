
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Lock, 
  ChevronRight, 
  ShieldCheck, 
  User, 
} from 'lucide-react';
import { cardBaseStyles, Footer } from './Widgets';
const Login = ({ onLogin }: { onLogin: () => void }) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState('');

  const handleLogin = (e:React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      onLogin();
      navigate('/profile');
    }, 1200);
  };

  return (
    <div className={`${cardBaseStyles} animate-in fade-in slide-in-from-bottom-4 duration-700`}>
      <div className="flex flex-col items-center mb-8">
        <div className="w-16 h-16 bg-blue-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-blue-200">
          <ShieldCheck className="text-white w-9 h-9" />
        </div>
        <h1 className="text-2xl font-bold text-slate-800">登 录</h1>
        <p className="text-slate-500 text-sm mt-2 text-center">使用您的账号登录以继续访问第三方应用</p>
      </div>

      <form onSubmit={handleLogin} className="space-y-5">
        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-700 ml-1">电子邮箱</label>
          <div className="relative">
            <User className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 w-5 h-5" />
            <input 
              required
              type="email" 
              placeholder="name@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all"
            />
          </div>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium text-slate-700 ml-1">密码</label>
          <div className="relative">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 w-5 h-5" />
            <input 
              required
              type="password" 
              placeholder="••••••••"
              className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all"
            />
          </div>
        </div>

        <div className="flex items-center justify-between text-sm py-2">
          <label className="flex items-center space-x-2 cursor-pointer group">
            <input type="checkbox" className="w-4 h-4 rounded border-slate-300 text-blue-500 focus:ring-blue-400" />
            <span className="text-slate-500 group-hover:text-slate-700 transition-colors">记住我</span>
          </label>
          <a href="#" className="text-blue-500 hover:text-blue-600 font-medium transition-colors">忘记密码？</a>
        </div>

        <button 
          disabled={isLoading}
          type="submit" 
          className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3.5 rounded-xl shadow-md shadow-blue-100 transition-all flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70"
        >
          {isLoading ? (
            <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
          ) : (
            <>
              <span>立即登录</span>
              <ChevronRight className="w-5 h-5" />
            </>
          )}
        </button>
      </form>

      <div className="mt-8 pt-8 border-t border-slate-100 flex flex-col items-center">
        <p className="text-sm text-slate-500">还没有账号？</p>
        <button onClick={() => navigate('/register')} className="mt-2 text-blue-500 font-semibold hover:underline decoration-2 underline-offset-4">创建一个新账号</button>
      </div>
      <Footer className="mt-8 sm:hidden opacity-60" />
    </div>
  );
};

export default Login;