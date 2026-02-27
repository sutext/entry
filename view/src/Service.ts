export type User = {
  id: number;
  email: string;
  nickname: string;
  phone: string;
  gender: number;
  birthday: number;
  height: number;
  weight: number;
  lastLogin: string;
  username: string;
};
export type RegisterFormData = {
  email: string;
  password: string;
};
const TokenKey = 'token';
export const getToken = () => localStorage.getItem(TokenKey);
export const register = async (formData: RegisterFormData) => {
  const res = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    });
    if (res.status === 200) {
        return res.json().then((data) => {
            localStorage.setItem(TokenKey, data.token);
            return data.user as User;
        });
    } else {
        throw new Error('账号创建失败，请重试。');
    }
};
export type LoginFormData = {
  email: string;
  password: string;
};
export const login = async (formData: LoginFormData) => {

    const res = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    });
    if (res.status === 200) {
        return res.json().then((data) => {
            localStorage.setItem(TokenKey, data.token);
            return data.user as User;
        });
    } else {
        throw new Error('登录失败，请重试。');
    }
}

export const approve = async (search: URLSearchParams)=>{
    const token = localStorage.getItem(TokenKey);
    if (!token) {
        throw new Error('请先登录。');
    }
    const res = await fetch(`/approve?${search.toString()}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
    });
    if (res.status === 200) {
        return res.json().then((data) => {
            return data.user as User;
        });
    } else {
        throw new Error('授权失败，请重试。');
    }
}

export const profile = async () => {
    const token = localStorage.getItem(TokenKey);
    if (!token) {
        throw new Error('请先登录。');
    }
    const res = await fetch('/profile', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
    });
    if (res.status === 200) {
        return res.json().then((data) => {
            return data.user as User;
        });
    } else {
        throw new Error('获取用户信息失败，请重试。');
    }
};