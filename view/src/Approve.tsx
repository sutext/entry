import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  CheckCircle2, 
  ShieldCheck, 
  Layout, 
  Globe, 
  Smartphone 
} from 'lucide-react';
import { cardBaseStyles, Footer } from './Widgets';

const Approve = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);

  const handleApprove = () => {
    setIsLoading(true);
    setTimeout(() => {
      alert('授权成功！正在重定向回第三方应用...');
      setIsLoading(false);
    }, 1000);
  };

  return (
    <div className={`${cardBaseStyles} animate-in fade-in zoom-in-95 duration-500`}>
      <div className="flex justify-between items-center mb-10 relative">
        <div className="w-14 h-14 bg-blue-100 rounded-2xl flex items-center justify-center z-10">
          <ShieldCheck className="text-blue-600 w-8 h-8" />
        </div>
        <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-24 border-t-2 border-dashed border-slate-200"></div>
        <div className="w-14 h-14 bg-indigo-50 rounded-2xl flex items-center justify-center z-10 border border-indigo-100">
          <Layout className="text-indigo-500 w-8 h-8" />
        </div>
      </div>

      <div className="text-center mb-8">
        <h2 className="text-xl font-bold text-slate-800">应用授权请求</h2>
        <p className="text-slate-500 text-sm mt-2 leading-relaxed">
          <span className="font-semibold text-slate-700">"Creative Studio"</span> 正在请求访问您的账号权限。
        </p>
      </div>

      <div className="bg-slate-50 rounded-2xl p-6 mb-8 border border-slate-100">
        <p className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-4">该应用将能够：</p>
        <ul className="space-y-4">
          <li className="flex items-start space-x-3">
            <div className="mt-1 bg-green-100 rounded-full p-0.5">
              <CheckCircle2 className="w-4 h-4 text-green-600" />
            </div>
            <div className="text-sm">
              <p className="font-medium text-slate-700">访问您的公开信息</p>
              <p className="text-slate-500 text-xs">包括姓名、头像和性别</p>
            </div>
          </li>
          <li className="flex items-start space-x-3">
            <div className="mt-1 bg-green-100 rounded-full p-0.5">
              <CheckCircle2 className="w-4 h-4 text-green-600" />
            </div>
            <div className="text-sm">
              <p className="font-medium text-slate-700">查看您的电子邮箱地址</p>
              <p className="text-slate-500 text-xs">仅用于身份验证和通知</p>
            </div>
          </li>
        </ul>
      </div>

      <div className="flex flex-col space-y-3">
        <button 
          onClick={handleApprove}
          disabled={isLoading}
          className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3.5 rounded-xl shadow-lg shadow-blue-100 transition-all active:scale-[0.98] flex items-center justify-center"
        >
          {isLoading ? (
            <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
          ) : "授权并继续"}
        </button>
        <button 
          onClick={() => navigate('/login')}
          className="w-full bg-white hover:bg-slate-50 text-slate-500 font-medium py-3.5 rounded-xl border border-slate-200 transition-all active:scale-[0.98]"
        >
          返回登录
        </button>
      </div>

      <div className="mt-8 flex items-center justify-center space-x-4 text-slate-400">
        <div className="flex items-center space-x-1">
          <Globe className="w-3.5 h-3.5" />
          <span className="text-xs">官方验证</span>
        </div>
        <div className="w-1 h-1 bg-slate-300 rounded-full"></div>
        <div className="flex items-center space-x-1">
          <Smartphone className="w-3.5 h-3.5" />
          <span className="text-xs">加密连接</span>
        </div>
      </div>
      <Footer className="mt-10 sm:hidden opacity-60" />
    </div>
  );
};

export default Approve