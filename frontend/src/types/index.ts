export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Company {
  id: string;
  user_id: string;
  name: string;
  ein: string;
  address: string;
  phone?: string;
  email?: string;
  created_at: string;
  updated_at: string;
}

export interface Vehicle {
  id: string;
  company_id: string;
  type: 'truck' | 'trailer';
  vin: string;
  make: string;
  model: string;
  year: number;
  purchase_price: number;
  purchase_date: string;
  status: 'active' | 'inactive' | 'sold';
  created_at: string;
  updated_at: string;
}

export interface Loan {
  id: string;
  vehicle_id: string;
  company_id: string;
  lender: string;
  principal_amount: number;
  interest_rate: number;
  term_months: number;
  start_date: string;
  monthly_payment: number;
  remaining_balance: number;
  status: 'active' | 'paid_off';
  created_at: string;
  updated_at: string;
}

export interface Payment {
  id: string;
  loan_id: string;
  payment_date: string;
  principal_paid: number;
  interest_paid: number;
  total_paid: number;
  remaining_balance: number;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface DebtScheduleItem {
  company_name: string;
  total_debt: number;
  monthly_payment: number;
  vehicles_count: number;
}

export interface AmortizationScheduleItem {
  payment_number: number;
  payment_date: string;
  principal_payment: number;
  interest_payment: number;
  total_payment: number;
  remaining_balance: number;
}

export interface DepreciationScheduleItem {
  vehicle_id: string;
  vehicle_name: string;
  purchase_price: number;
  current_value: number;
  depreciation_amount: number;
  age_years: number;
}

export interface DashboardStats {
  total_companies: number;
  total_vehicles: number;
  total_active_loans: number;
  total_debt: number;
  monthly_payments: number;
  total_asset_value: number;
  total_payments_year: number;
} 