import { ArrowLeft, Mail, User,Lock, UserPlus } from "lucide-react";
import { cardBaseStyles, Footer } from "./Widgets";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const Register = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);

  const handleRegister = (e:React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      // 注册成功后跳转登录
      alert('账号创建成功！请登录。');
      navigate('/login');
    }, 1500);
  };

  return (
    <div className={`${cardBaseStyles} animate-in fade-in slide-in-from-bottom-4 duration-700`}>
      <button 
        onClick={() => navigate('/login')}
        className="mb-6 flex items-center text-slate-400 hover:text-blue-500 transition-colors text-sm font-medium group"
      >
        <ArrowLeft className="w-4 h-4 mr-1 group-hover:-translate-x-1 transition-transform" />
        返回登录
      </button>

      <div className="flex flex-col items-center mb-8">
        <div className="w-16 h-16 bg-blue-50 rounded-2xl flex items-center justify-center mb-4 border border-blue-100">
          <UserPlus className="text-blue-500 w-8 h-8" />
        </div>
        <h1 className="text-2xl font-bold text-slate-800">创建新账号</h1>
        <p className="text-slate-500 text-sm mt-2 text-center">加入我们，开启安全便捷的授权体验</p>
      </div>

      <form onSubmit={handleRegister} className="space-y-4">
        <div className="space-y-1.5">
          <label className="text-sm font-medium text-slate-700 ml-1">用户名</label>
          <div className="relative">
            <User className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 w-5 h-5" />
            <input 
              required
              type="text" 
              placeholder="您的姓名"
              className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all"
            />
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-slate-700 ml-1">电子邮箱</label>
          <div className="relative">
            <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 w-5 h-5" />
            <input 
              required
              type="email" 
              placeholder="name@example.com"
              className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all"
            />
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-slate-700 ml-1">设置密码</label>
          <div className="relative">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 w-5 h-5" />
            <input 
              required
              type="password" 
              placeholder="至少 8 位字符"
              className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-400 focus:bg-white transition-all"
            />
          </div>
        </div>

        <div className="pt-2">
          <p className="text-[11px] text-slate-400 text-center mb-4">
            点击注册即表示您同意我们的 <a href="#" className="text-blue-500 hover:underline">服务条款</a> 和 <a href="#" className="text-blue-500 hover:underline">隐私政策</a>
          </p>
          <button 
            disabled={isLoading}
            type="submit" 
            className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3.5 rounded-xl shadow-md shadow-blue-100 transition-all flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70"
          >
            {isLoading ? (
              <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            ) : (
              <span>注册账号</span>
            )}
          </button>
        </div>
      </form>

      <Footer className="mt-8 sm:hidden opacity-60" />
    </div>
  );
};

export default Register;