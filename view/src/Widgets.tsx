
export const Background = () => (
  <div className="fixed inset-0 -z-10 overflow-hidden pointer-events-none">
    <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-blue-50 rounded-full blur-3xl opacity-60"></div>
    <div className="absolute bottom-[-5%] right-[-5%] w-[30%] h-[30%] bg-sky-100 rounded-full blur-3xl opacity-50"></div>
  </div>
);

export const Footer = ({ className }:{className:string}) => (
  <footer className={`${className} text-slate-400 text-xs text-center`}>
    <p>© 2024 AzureAuth Security. 保留所有权利。</p>
    <div className="mt-2 space-x-3">
      <a href="#" className="hover:text-blue-400 transition-colors">隐私政策</a>
      <span>•</span>
      <a href="#" className="hover:text-blue-400 transition-colors">服务条款</a>
    </div>
  </footer>
);

export const cardBaseStyles = "w-full max-w-md bg-white sm:rounded-3xl sm:shadow-sm sm:border sm:border-slate-100 p-8 md:p-10 transition-all duration-500 ease-in-out";
