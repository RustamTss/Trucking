import React, { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react';
import { companiesAPI } from '../services/api';
import type { Company } from '../types';

interface CompanyContextType {
  companies: Company[];
  selectedCompany: Company | null;
  isLoading: boolean;
  error: string | null;
  loadCompanies: () => Promise<void>;
  createCompany: (company: Omit<Company, 'id' | 'created_at' | 'updated_at'>) => Promise<Company>;
  updateCompany: (id: string, company: Partial<Company>) => Promise<Company>;
  deleteCompany: (id: string) => Promise<void>;
  setSelectedCompany: (company: Company | null) => void;
}

const CompanyContext = createContext<CompanyContextType | undefined>(undefined);

export const useCompanies = (): CompanyContextType => {
  const context = useContext(CompanyContext);
  if (!context) {
    throw new Error('useCompanies must be used within a CompanyProvider');
  }
  return context;
};

interface CompanyProviderProps {
  children: ReactNode;
}

export const CompanyProvider: React.FC<CompanyProviderProps> = ({ children }) => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadCompanies = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await companiesAPI.getAll();
      setCompanies(data || []);
      
      // Автоматически выбираем первую компанию если нет выбранной
      if (!selectedCompany && data && data.length > 0) {
        setSelectedCompany(data[0]);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка загрузки компаний';
      setError(errorMessage);
      setCompanies([]);
    } finally {
      setIsLoading(false);
    }
  }, [selectedCompany]);

  const createCompany = useCallback(async (companyData: Omit<Company, 'id' | 'created_at' | 'updated_at'>) => {
    setError(null);
    try {
      const newCompany = await companiesAPI.create(companyData);
      setCompanies(prev => [...prev, newCompany]);
      
      // Автоматически выбираем новую компанию если это первая
      if (companies.length === 0) {
        setSelectedCompany(newCompany);
      }
      
      return newCompany;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка создания компании';
      setError(errorMessage);
      throw err;
    }
  }, [companies.length]);

  const updateCompany = useCallback(async (id: string, companyData: Partial<Company>) => {
    setError(null);
    try {
      const updatedCompany = await companiesAPI.update(id, companyData);
      setCompanies(prev => prev.map(company => 
        company.id === id ? updatedCompany : company
      ));
      
      // Обновляем выбранную компанию если это она
      if (selectedCompany?.id === id) {
        setSelectedCompany(updatedCompany);
      }
      
      return updatedCompany;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка обновления компании';
      setError(errorMessage);
      throw err;
    }
  }, [selectedCompany?.id]);

  const deleteCompany = useCallback(async (id: string) => {
    setError(null);
    try {
      await companiesAPI.delete(id);
      setCompanies(prev => prev.filter(company => company.id !== id));
      
      // Сбрасываем выбранную компанию если удаляем её
      if (selectedCompany?.id === id) {
        const remainingCompanies = companies.filter(company => company.id !== id);
        setSelectedCompany(remainingCompanies.length > 0 ? remainingCompanies[0] : null);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка удаления компании';
      setError(errorMessage);
      throw err;
    }
  }, [companies, selectedCompany?.id]);

  const value = useMemo(() => ({
    companies,
    selectedCompany,
    isLoading,
    error,
    loadCompanies,
    createCompany,
    updateCompany,
    deleteCompany,
    setSelectedCompany,
  }), [
    companies,
    selectedCompany,
    isLoading,
    error,
    loadCompanies,
    createCompany,
    updateCompany,
    deleteCompany,
  ]);

  return (
    <CompanyContext.Provider value={value}>
      {children}
    </CompanyContext.Provider>
  );
}; 