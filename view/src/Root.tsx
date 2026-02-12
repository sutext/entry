import { ChevronRight, Key, Shield, ShieldCheck, Sparkles, User, UserPlus, Zap } from "lucide-react";
import React from "react";
import { useNavigate } from "react-router-dom";
import { cardBaseStyles, Footer } from "./Widgets";

const Root = () => {
  const navigate = useNavigate();

  const features = [
    { icon: <Key className="text-blue-500" />, title: "安全登录", desc: "采用行业标准的身份验证流程" },
    { icon: <UserPlus className="text-indigo-500" />, title: "快速注册", desc: "极简流程，即刻开启安全旅程" },
    { icon: <ShieldCheck className="text-sky-500" />, title: "OAuth2 授权", desc: "精准控制应用对隐私数据的访问" },
    { icon: <User className="text-emerald-500" />, title: "个人资料", desc: "完善的个人中心，随时管理信息" }
  ];

  return (
    <div className="w-full flex flex-col items-center">
      {/* PC 端导航栏 - 仅在桌面端显示 */}
      <nav className="hidden md:flex w-full justify-between items-center px-12 py-8 max-w-7xl animate-in fade-in duration-500">
        <div className="flex items-center space-x-2">
          <Shield className="text-blue-600 w-8 h-8" />
          <span className="font-bold text-2xl text-slate-800">AzureAuth</span>
        </div>
        <div className="flex items-center space-x-6">
          <button onClick={() => navigate('/login')} className="text-slate-600 font-medium hover:text-blue-600 transition-colors">登录</button>
          <button onClick={() => navigate('/register')} className="bg-slate-900 text-white px-6 py-2.5 rounded-full font-bold hover:bg-slate-800 transition-all shadow-lg shadow-slate-200">注册账号</button>
        </div>
      </nav>

      {/* PC 端内容布局 - 沉浸式两栏 */}
      <div className="hidden md:flex flex-row items-center justify-between w-full max-w-7xl px-12 py-16 gap-12">
        <div className="w-1/2 space-y-8 animate-in slide-in-from-left-8 duration-700">
          <div className="inline-flex items-center space-x-2 px-4 py-1.5 bg-blue-50 text-blue-600 rounded-full text-sm font-bold border border-blue-100">
            <Sparkles className="w-4 h-4" />
            <span>全新 OAuth2 协议支持</span>
          </div>
          <h1 className="text-7xl font-black text-slate-800 leading-[1.1]">
            统一安全<br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-500">身份授权中心</span>
          </h1>
          <p className="text-xl text-slate-500 max-w-lg leading-relaxed">
            下一代身份验证解决方案。为开发者提供极速集成体验，为用户提供固若金汤的隐私保护。
          </p>
          <div className="flex space-x-4">
            <button onClick={() => navigate('/login')} className="px-10 py-5 bg-blue-600 text-white font-bold rounded-2xl shadow-xl shadow-blue-200 hover:bg-blue-700 transition-all flex items-center group active:scale-95">
              立即开始 <ChevronRight className="ml-2 group-hover:translate-x-1 transition-transform" />
            </button>
          </div>
        </div>
        <div className="w-1/2 grid grid-cols-2 gap-6 animate-in fade-in zoom-in-95 duration-1000">
          {features.map((f, i) => (
            <div key={i} className="p-8 bg-white/60 backdrop-blur-md rounded-[2.5rem] border border-white hover:border-blue-200 shadow-sm transition-all group">
              <div className="w-14 h-14 bg-white rounded-2xl flex items-center justify-center shadow-sm border border-slate-50 mb-6 group-hover:scale-110 transition-transform">
                {f.icon}
              </div>
              <h3 className="text-lg font-bold text-slate-800 mb-2">{f.title}</h3>
              <p className="text-sm text-slate-500 leading-relaxed">{f.desc}</p>
            </div>
          ))}
        </div>
      </div>

      {/* 移动端内容布局 - 精致卡片设计 */}
      <div className="md:hidden flex flex-col items-center w-full px-4 pt-6 pb-12">
        <div className={`${cardBaseStyles} animate-in fade-in slide-in-from-bottom-4 duration-700`}>
          <div className="flex flex-col items-center mb-6">
            <div className="w-14 h-14 bg-gradient-to-tr from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-blue-200">
              <Shield className="text-white w-8 h-8" />
            </div>
            <h1 className="text-xl font-bold text-slate-800 tracking-tight text-center">AzureAuth 安全中心</h1>
            <p className="text-slate-500 text-xs mt-2 text-center max-w-[240px]">
              受信任的 OAuth2 授权服务
            </p>
          </div>

          <div className="space-y-2 mb-8">
            {features.map((f, i) => (
              <div key={i} className="flex items-center p-3 bg-slate-50/50 rounded-2xl border border-slate-100">
                <div className="w-9 h-9 bg-white rounded-xl flex items-center justify-center shadow-sm mr-3 shrink-0">
                  {React.cloneElement(f.icon, { size: 18 })}
                </div>
                <div>
                  <h3 className="text-xs font-bold text-slate-700">{f.title}</h3>
                  <p className="text-[10px] text-slate-400 line-clamp-1">{f.desc}</p>
                </div>
              </div>
            ))}
          </div>

          <div className="flex flex-col space-y-3">
            <button onClick={() => navigate('/login')} className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-4 rounded-2xl shadow-md shadow-blue-100 flex items-center justify-center space-x-2 active:scale-95 transition-all">
              <Zap className="w-4 h-4" />
              <span>登录账户</span>
            </button>
            <button onClick={() => navigate('/register')} className="w-full bg-white hover:bg-slate-50 text-slate-600 font-bold py-4 rounded-2xl border border-slate-200 active:scale-95 transition-all">
              注册新账号
            </button>
          </div>
        </div>
      </div>
      <Footer className="py-8 bg-white w-full" />
    </div>
  );
};
export default Root;