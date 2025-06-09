import { Edit, Filter, Plus, Trash2, Truck } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useCompanies } from '../../context/CompanyContext';
import { useVehicles } from '../../context/VehicleContext';
import BackButton from '../common/BackButton';

const VehicleList: React.FC = () => {
  const { companies } = useCompanies();
  const { vehicles, deleteVehicle, isLoading, loadVehicles } = useVehicles();
  const [searchParams, setSearchParams] = useSearchParams();
  const [deletingId, setDeletingId] = useState<string | null>(null);
  
  const selectedCompanyId = searchParams.get('company_id') || '';
  const selectedCompany = companies.find(c => c.id === selectedCompanyId);

  useEffect(() => {
    loadVehicles();
  }, [loadVehicles]);

  const filteredVehicles = selectedCompanyId 
    ? vehicles.filter(v => v.company_id === selectedCompanyId)
    : vehicles;

  const handleDelete = async (id: string, vehicleInfo: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить транспорт "${vehicleInfo}"? Это действие нельзя отменить.`)) {
      return;
    }

    setDeletingId(id);
    try {
      await deleteVehicle(id);
    } catch (error) {
      console.error('Ошибка удаления транспорта:', error);
      alert('Ошибка при удалении транспорта');
    } finally {
      setDeletingId(null);
    }
  };

  const handleCompanyChange = (companyId: string) => {
    if (companyId) {
      setSearchParams({ company_id: companyId });
    } else {
      setSearchParams({});
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent, action: () => void) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      action();
    }
  };

  const getVehicleTypeLabel = (type: string) => {
    return type === 'truck' ? 'Трак' : 'Трейлер';
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'active': return 'Активный';
      case 'inactive': return 'Неактивный';
      case 'sold': return 'Продан';
      default: return status;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-800';
      case 'inactive': return 'bg-yellow-100 text-yellow-800';
      case 'sold': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-center py-12">
            <div className="text-lg text-gray-500">Загрузка транспорта...</div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <BackButton to="/dashboard" className="mb-4" />
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <Truck className="h-8 w-8 text-indigo-600 mr-3" />
              <div>
                <h1 className="text-3xl font-bold text-gray-900">Управление транспортом</h1>
                <p className="text-gray-600">
                  {selectedCompany 
                    ? `${selectedCompany.name}: ${filteredVehicles.length} единиц транспорта`
                    : `Всего транспорта: ${filteredVehicles.length}`
                  }
                </p>
              </div>
            </div>
            <Link
              to={`/vehicles/new${selectedCompanyId ? `?company_id=${selectedCompanyId}` : ''}`}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              aria-label="Добавить новый транспорт"
            >
              <Plus className="w-4 h-4 mr-2" />
              Добавить транспорт
            </Link>
          </div>
        </div>

        {/* Company Filter */}
        <div className="mb-6">
          <div className="flex items-center space-x-4">
            <Filter className="h-5 w-5 text-gray-400" />
            <select
              value={selectedCompanyId}
              onChange={(e) => handleCompanyChange(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-2 bg-white shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Все компании</option>
              {companies?.map(company => (
                <option key={company.id} value={company.id}>
                  {company.name}
                </option>
              )) || []}
            </select>
          </div>
        </div>

        {/* Vehicle List */}
        {filteredVehicles.length === 0 ? (
          <div className="text-center py-12">
            <Truck className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">Нет транспорта</h3>
            <p className="mt-1 text-sm text-gray-500">
              {selectedCompany 
                ? `В компании "${selectedCompany.name}" пока нет транспорта.`
                : 'Начните с добавления первого транспорта.'
              }
            </p>
            <div className="mt-6">
              <Link
                to={`/vehicles/new${selectedCompanyId ? `?company_id=${selectedCompanyId}` : ''}`}
                className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
                aria-label="Добавить первый транспорт"
              >
                <Truck className="-ml-1 mr-2 h-5 w-5" />
                Добавить транспорт
              </Link>
            </div>
          </div>
        ) : (
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {filteredVehicles.map((vehicle) => {
                const company = companies.find(c => c.id === vehicle.company_id);
                const vehicleInfo = `${vehicle.year} ${vehicle.make} ${vehicle.model}`;
                
                return (
                  <li key={vehicle.id} className="px-6 py-4">
                    <div className="flex items-center justify-between">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center">
                          <Truck className="h-5 w-5 text-gray-400 mr-3" />
                          <div>
                            <div className="flex items-center space-x-2">
                              <h3 className="text-lg font-medium text-gray-900">
                                {vehicleInfo}
                              </h3>
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                {getVehicleTypeLabel(vehicle.type)}
                              </span>
                              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(vehicle.status)}`}>
                                {getStatusLabel(vehicle.status)}
                              </span>
                            </div>
                            <div className="mt-1 text-sm text-gray-500 space-y-1">
                              <p>VIN: {vehicle.vin}</p>
                              <p>Цена покупки: ${vehicle.purchase_price.toLocaleString()}</p>
                              <p>Дата покупки: {new Date(vehicle.purchase_date).toLocaleDateString('ru-RU')}</p>
                              {company && !selectedCompany && (
                                <p>Компания: {company.name}</p>
                              )}
                            </div>
                          </div>
                        </div>
                      </div>
                      
                      <div className="flex items-center space-x-2">
                        {/* Edit Button */}
                        <Link
                          to={`/vehicles/${vehicle.id}/edit`}
                          className="inline-flex items-center p-2 border border-gray-300 rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                          aria-label={`Редактировать транспорт ${vehicleInfo}`}
                        >
                          <Edit className="h-4 w-4" />
                        </Link>

                        {/* Delete Button */}
                        <button
                          onClick={() => handleDelete(vehicle.id, vehicleInfo)}
                          onKeyDown={(e) => handleKeyDown(e, () => handleDelete(vehicle.id, vehicleInfo))}
                          disabled={deletingId === vehicle.id}
                          className="inline-flex items-center p-2 border border-red-300 rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                          tabIndex={0}
                          aria-label={`Удалить транспорт ${vehicleInfo}`}
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  </li>
                );
              })}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};

export default VehicleList; 