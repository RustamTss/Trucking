import React, { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react';
import { vehiclesAPI } from '../services/api';
import type { Vehicle } from '../types';

interface VehicleContextType {
  vehicles: Vehicle[];
  isLoading: boolean;
  error: string | null;
  loadVehicles: () => Promise<void>;
  loadVehiclesByCompany: (companyId: string) => Promise<void>;
  createVehicle: (vehicle: Omit<Vehicle, 'id' | 'created_at' | 'updated_at'>) => Promise<Vehicle>;
  updateVehicle: (id: string, vehicle: Partial<Vehicle>) => Promise<Vehicle>;
  deleteVehicle: (id: string) => Promise<void>;
  getVehiclesByCompany: (companyId: string) => Vehicle[];
}

const VehicleContext = createContext<VehicleContextType | undefined>(undefined);

export const useVehicles = (): VehicleContextType => {
  const context = useContext(VehicleContext);
  if (!context) {
    throw new Error('useVehicles must be used within a VehicleProvider');
  }
  return context;
};

interface VehicleProviderProps {
  children: ReactNode;
}

export const VehicleProvider: React.FC<VehicleProviderProps> = ({ children }) => {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadVehicles = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await vehiclesAPI.getAll();
      setVehicles(data || []);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка загрузки транспорта';
      setError(errorMessage);
      setVehicles([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loadVehiclesByCompany = useCallback(async (companyId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      // Предполагаем что API поддерживает фильтрацию по company_id
      const data = await vehiclesAPI.getAll();
      const filteredVehicles = data?.filter(v => v.company_id === companyId) || [];
      setVehicles(filteredVehicles);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка загрузки транспорта';
      setError(errorMessage);
      setVehicles([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createVehicle = useCallback(async (vehicleData: Omit<Vehicle, 'id' | 'created_at' | 'updated_at'>) => {
    setError(null);
    try {
      const newVehicle = await vehiclesAPI.create(vehicleData);
      setVehicles(prev => [...prev, newVehicle]);
      return newVehicle;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка создания транспорта';
      setError(errorMessage);
      throw err;
    }
  }, []);

  const updateVehicle = useCallback(async (id: string, vehicleData: Partial<Vehicle>) => {
    setError(null);
    try {
      const updatedVehicle = await vehiclesAPI.update(id, vehicleData);
      setVehicles(prev => prev.map(vehicle => 
        vehicle.id === id ? updatedVehicle : vehicle
      ));
      return updatedVehicle;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка обновления транспорта';
      setError(errorMessage);
      throw err;
    }
  }, []);

  const deleteVehicle = useCallback(async (id: string) => {
    setError(null);
    try {
      await vehiclesAPI.delete(id);
      setVehicles(prev => prev.filter(vehicle => vehicle.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Ошибка удаления транспорта';
      setError(errorMessage);
      throw err;
    }
  }, []);

  const getVehiclesByCompany = useCallback((companyId: string) => {
    return vehicles.filter(vehicle => vehicle.company_id === companyId);
  }, [vehicles]);

  const value = useMemo(() => ({
    vehicles,
    isLoading,
    error,
    loadVehicles,
    loadVehiclesByCompany,
    createVehicle,
    updateVehicle,
    deleteVehicle,
    getVehiclesByCompany,
  }), [
    vehicles,
    isLoading,
    error,
    loadVehicles,
    loadVehiclesByCompany,
    createVehicle,
    updateVehicle,
    deleteVehicle,
    getVehiclesByCompany,
  ]);

  return (
    <VehicleContext.Provider value={value}>
      {children}
    </VehicleContext.Provider>
  );
}; 