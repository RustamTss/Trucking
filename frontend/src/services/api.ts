import axios from 'axios'
import type {
	AmortizationScheduleItem,
	Company,
	DashboardStats,
	DebtScheduleItem,
	DepreciationScheduleItem,
	Loan,
	LoginRequest,
	LoginResponse,
	Payment,
	User,
	Vehicle
} from '../types'

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Добавляем токен к каждому запросу
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Обработка ошибок авторизации
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },
  
  register: async (data: { email: string; password: string; name: string }): Promise<LoginResponse> => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },
};

// Companies API
export const companiesAPI = {
  getAll: async (): Promise<Company[]> => {
    const response = await api.get('/companies');
    return response.data;
  },
  
  create: async (data: Omit<Company, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<Company> => {
    const response = await api.post('/companies', data);
    return response.data;
  },
  
  update: async (id: string, data: Partial<Company>): Promise<Company> => {
    const response = await api.put(`/companies/${id}`, data);
    return response.data;
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/companies/${id}`);
  },
};

// Vehicles API
export const vehiclesAPI = {
  getAll: async (companyId?: string): Promise<Vehicle[]> => {
    const params = companyId ? { company_id: companyId } : {};
    const response = await api.get('/vehicles', { params });
    return response.data;
  },
  
  create: async (data: Omit<Vehicle, 'id' | 'created_at' | 'updated_at'>): Promise<Vehicle> => {
    const response = await api.post('/vehicles', data);
    return response.data;
  },
  
  update: async (id: string, data: Partial<Vehicle>): Promise<Vehicle> => {
    const response = await api.put(`/vehicles/${id}`, data);
    return response.data;
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/vehicles/${id}`);
  },
};

// Loans API
export const loansAPI = {
  getAll: async (companyId?: string): Promise<Loan[]> => {
    const params = companyId ? { company_id: companyId } : {};
    const response = await api.get('/loans', { params });
    return response.data;
  },
  
  create: async (data: Omit<Loan, 'id' | 'created_at' | 'updated_at'>): Promise<Loan> => {
    const response = await api.post('/loans', data);
    return response.data;
  },
  
  update: async (id: string, data: Partial<Loan>): Promise<Loan> => {
    const response = await api.put(`/loans/${id}`, data);
    return response.data;
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/loans/${id}`);
  },
};

// Payments API
export const paymentsAPI = {
  getAll: async (): Promise<Payment[]> => {
    const response = await api.get('/payments');
    return response.data;
  },
  
  getByLoan: async (loanId: string): Promise<Payment[]> => {
    const response = await api.get(`/payments/loan/${loanId}`);
    return response.data;
  },
  
  create: async (data: Omit<Payment, 'id' | 'created_at'>): Promise<Payment> => {
    const response = await api.post('/payments', data);
    return response.data;
  },
};

// Schedules API
export const schedulesAPI = {
  getDebtSchedule: async (): Promise<DebtScheduleItem[]> => {
    const response = await api.get('/schedules/debt');
    return response.data;
  },
  
  getAmortizationSchedule: async (loanId?: string): Promise<AmortizationScheduleItem[]> => {
    const params = loanId ? { loan_id: loanId } : {};
    const response = await api.get('/schedules/amortization', { params });
    return response.data;
  },
  
  getDepreciationSchedule: async (companyId?: string): Promise<DepreciationScheduleItem[]> => {
    const params = companyId ? { company_id: companyId } : {};
    const response = await api.get('/schedules/depreciation', { params });
    return response.data;
  },
};

// Stats API
export const statsAPI = {
  getDashboardStats: async (): Promise<DashboardStats> => {
    const response = await api.get('/stats/dashboard');
    return response.data;
  },
};

// Users API
export const usersAPI = {
  getAll: async (): Promise<User[]> => {
    const response = await api.get('/users');
    return response.data;
  },
  
  getProfile: async (): Promise<User> => {
    const response = await api.get('/users/profile');
    return response.data;
  },
  
  getById: async (id: string): Promise<User> => {
    const response = await api.get(`/users/${id}`);
    return response.data;
  },
  
  update: async (id: string, data: Partial<User>): Promise<User> => {
    const response = await api.put(`/users/${id}`, data);
    return response.data;
  },
  
  delete: async (id: string): Promise<void> => {
    await api.delete(`/users/${id}`);
  },
  
  updateProfile: async (data: { name?: string; email?: string; password?: string }): Promise<User> => {
    const userStr = localStorage.getItem('user');
    if (!userStr) throw new Error('User not found');
    const user = JSON.parse(userStr);
    const response = await api.put(`/users/${user.id}`, data);
    return response.data;
  },
}; 