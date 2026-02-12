import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {  Footer } from "./Widgets";
import { ArrowLeft, Calendar, Camera, LogOut, Mail, Phone, Ruler, Shield, User, Weight, type LucideProps } from "lucide-react";
import React from "react";

/* eslint-disable @typescript-eslint/no-unused-vars */

const Profile = ({ onLogout }: { onLogout: () => void }) => {
  const navigate = useNavigate();

  const [user] = useState({
    nickname: '张小帅',
    email: 'zhang@example.com',
    phone: '138 0013 8000',
    gender: '男',
    birthday: '1998-06-15',
    height: '180 cm',
    weight: '75 kg',
    lastLogin: '2024-05-20 14:30'
  });

  return (
    <div className="w-full min-h-screen bg-slate-50/50 flex flex-col animate-in fade-in duration-500">
      {/* Top Header */}
      <header className="bg-white border-b border-slate-100 px-6 py-4 md:px-12 flex justify-between items-center sticky top-0 z-20">
        <div className="flex items-center space-x-4">
          <button onClick={() => navigate('/authorize')} className="p-2 hover:bg-slate-50 rounded-full text-slate-400">
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h1 className="text-lg font-bold text-slate-800">个人账号设置</h1>
        </div>
        <button onClick={onLogout} className="flex items-center space-x-2 text-sm font-semibold text-red-500 hover:bg-red-50 px-4 py-2 rounded-full transition-all">
          <LogOut className="w-4 h-4" />
          <span>安全退出</span>
        </button>
      </header>

      <div className="flex-1 flex flex-col md:flex-row max-w-6xl mx-auto w-full p-6 md:p-12 gap-8">
        {/* Sidebar info */}
        <div className="md:w-1/3 space-y-6">
          <div className="bg-white p-8 rounded-[2rem] border border-slate-100 shadow-sm flex flex-col items-center text-center">
            <div className="relative group mb-4">
              <div className="w-32 h-32 bg-gradient-to-tr from-blue-500 to-indigo-500 rounded-[2.5rem] flex items-center justify-center text-white text-5xl font-bold shadow-2xl shadow-blue-100 ring-8 ring-white">
                {user.nickname[0]}
              </div>
              <button className="absolute bottom-0 right-0 w-10 h-10 bg-white border border-slate-100 rounded-2xl shadow-lg flex items-center justify-center text-blue-500 hover:scale-110 transition-transform">
                <Camera className="w-5 h-5" />
              </button>
            </div>
            <h2 className="text-2xl font-bold text-slate-800">{user.nickname}</h2>
            <p className="text-slate-400 text-sm mt-1">AzureAuth 高级认证用户</p>
            <div className="mt-6 pt-6 border-t border-slate-50 w-full">
              <p className="text-xs text-slate-400 uppercase tracking-widest font-bold">最后登录</p>
              <p className="text-sm text-slate-600 mt-1 font-medium">{user.lastLogin}</p>
            </div>
          </div>

          <div className="bg-blue-600 p-6 rounded-[2rem] text-white shadow-xl shadow-blue-200">
            <h4 className="font-bold flex items-center">
              <Shield className="w-4 h-4 mr-2" />
              账户安全等级
            </h4>
            <div className="mt-4 flex items-center justify-between">
              <div className="h-2 flex-1 bg-white/20 rounded-full overflow-hidden">
                <div className="h-full bg-white w-4/5"></div>
              </div>
              <span className="ml-4 font-bold text-sm">较高</span>
            </div>
            <p className="text-xs text-blue-100 mt-4 leading-relaxed">您的账户已启用双重验证。建议定期更换密码以确保持续安全。</p>
          </div>
        </div>

        {/* Details Grid */}
        <div className="md:w-2/3 space-y-6">
          <div className="bg-white rounded-[2rem] border border-slate-100 shadow-sm overflow-hidden">
            <div className="px-8 py-6 border-b border-slate-50 flex justify-between items-center">
              <h3 className="font-bold text-slate-800">基本资料</h3>
              <button className="text-blue-500 text-sm font-bold hover:underline">编辑信息</button>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2">
              <InfoRow label="电子邮箱" value={user.email} icon={<Mail />} />
              <InfoRow label="联系电话" value={user.phone} icon={<Phone />} />
              <InfoRow label="性别" value={user.gender} icon={<User />} />
              <InfoRow label="出生日期" value={user.birthday} icon={<Calendar />} />
            </div>
          </div>

          <div className="bg-white rounded-[2rem] border border-slate-100 shadow-sm overflow-hidden">
            <div className="px-8 py-6 border-b border-slate-50">
              <h3 className="font-bold text-slate-800">身体健康数据</h3>
            </div>
            <div className="p-8 grid grid-cols-2 gap-6">
              <div className="bg-slate-50 p-6 rounded-3xl border border-slate-100/50">
                <div className="flex items-center text-slate-400 mb-2">
                  <Ruler className="w-4 h-4 mr-2" />
                  <span className="text-sm font-bold uppercase tracking-wider">身高</span>
                </div>
                <p className="text-3xl font-black text-slate-800">{user.height}</p>
              </div>
              <div className="bg-slate-50 p-6 rounded-3xl border border-slate-100/50">
                <div className="flex items-center text-slate-400 mb-2">
                  <Weight className="w-4 h-4 mr-2" />
                  <span className="text-sm font-bold uppercase tracking-wider">体重</span>
                </div>
                <p className="text-3xl font-black text-slate-800">{user.weight}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <Footer className="py-8 bg-white" />
    </div>
  );
};
export default Profile;

const InfoRow = ({ label, value, icon }: { label: string, value: string, icon: React.ReactElement<LucideProps> }) => (
  <div className="p-8 border-b md:border-r border-slate-50 last:border-0 hover:bg-slate-50 transition-colors group cursor-default">
    <div className="flex items-center space-x-4">
      <div className="w-10 h-10 bg-slate-50 group-hover:bg-white rounded-xl flex items-center justify-center text-slate-400 group-hover:text-blue-500 border border-transparent group-hover:border-slate-100 transition-all">
        {React.cloneElement(icon, { size: 18 })}
      </div>
      <div>
        <p className="text-[10px] uppercase tracking-widest text-slate-400 font-bold">{label}</p>
        <p className="text-base font-bold text-slate-700 mt-0.5">{value}</p>
      </div>
    </div>
  </div>
);