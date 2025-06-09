import { Building2, CreditCard, DollarSign, Plus, Settings, Truck, User } from 'lucide-react'
import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useCompanies } from '../context/CompanyContext'
import { loansAPI, vehiclesAPI } from '../services/api'
import type { Loan, Vehicle } from '../types'

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const { companies, selectedCompany, setSelectedCompany, loadCompanies } = useCompanies();
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loans, setLoans] = useState<Loan[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      await loadCompanies();
      const [vehiclesData, loansData] = await Promise.all([
        vehiclesAPI.getAll(),
        loansAPI.getAll(),
      ]);
      
      setVehicles(vehiclesData || []);
      setLoans(loansData || []);
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
      setVehicles([]);
      setLoans([]);
    } finally {
      setIsLoading(false);
    }
  };

  const getCompanyStats = (companyId: string) => {
    const companyVehicles = vehicles?.filter(v => v.company_id === companyId) || [];
    const companyLoans = loans?.filter(l => l.company_id === companyId) || [];
    const totalDebt = companyLoans.reduce((sum, loan) => sum + (loan.remaining_balance || 0), 0);
    const monthlyPayments = companyLoans.reduce((sum, loan) => sum + (loan.monthly_payment || 0), 0);
    
    return {
      vehiclesCount: companyVehicles.length,
      trucksCount: companyVehicles.filter(v => v.type === 'truck').length,
      trailersCount: companyVehicles.filter(v => v.type === 'trailer').length,
      totalDebt,
      monthlyPayments,
      activeLoans: companyLoans.filter(l => l.status === 'active').length,
    };
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Загрузка...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Business Schedule</h1>
              <p className="text-gray-600">Добро пожаловать, {user?.name}</p>
            </div>
            <div className="flex items-center space-x-4">
              {companies && companies.length > 0 && (
                <select
                  value={selectedCompany?.id || ''}
                  onChange={(e) => {
                    const company = companies.find((c: any) => c.id === e.target.value);
                    setSelectedCompany(company || null);
                  }}
                  className="border border-gray-300 rounded-md px-3 py-2"
                >
                  <option value="">Выберите компанию</option>
                  {companies?.map((company: any) => (
                    <option key={company.id} value={company.id}>
                      {company.name}
                    </option>
                  )) || []}
                </select>
              )}
              
              <Link
                to="/companies"
                className="flex items-center px-3 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md"
              >
                <Settings className="h-5 w-5 mr-2" />
                Управление
              </Link>
              
              <Link
                to="/profile"
                className="flex items-center px-3 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md"
              >
                <User className="h-5 w-5 mr-2" />
                Профиль
              </Link>
              <button
                onClick={logout}
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
              >
                Выйти
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {!companies || companies.length === 0 ? (
          <div className="text-center py-12">
            <Building2 className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">Нет компаний</h3>
            <p className="mt-1 text-sm text-gray-500">
              Начните с создания вашей первой компании.
            </p>
            <div className="mt-6">
              <Link
                to="/companies/new"
                className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
              >
                <Plus className="-ml-1 mr-2 h-5 w-5" />
                Добавить компанию
              </Link>
            </div>
          </div>
        ) : (
          <div className="px-4 py-6 sm:px-0">
            {selectedCompany && (
              <>
                {/* Company Info */}
                <div className="mb-8">
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">
                    {selectedCompany.name}
                  </h2>
                  <p className="text-gray-600">EIN: {selectedCompany.ein}</p>
                  <p className="text-gray-600">{selectedCompany.address}</p>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                  {(() => {
                    const stats = getCompanyStats(selectedCompany.id);
                    return (
                      <>
                        <div className="bg-white overflow-hidden shadow rounded-lg">
                          <div className="p-5">
                            <div className="flex items-center">
                              <div className="flex-shrink-0">
                                <Truck className="h-6 w-6 text-gray-400" />
                              </div>
                              <div className="ml-5 w-0 flex-1">
                                <dl>
                                  <dt className="text-sm font-medium text-gray-500 truncate">
                                    Транспорт
                                  </dt>
                                  <dd className="text-lg font-medium text-gray-900">
                                    {stats.vehiclesCount} ({stats.trucksCount} траков, {stats.trailersCount} трейлеров)
                                  </dd>
                                </dl>
                              </div>
                            </div>
                          </div>
                        </div>

                        <div className="bg-white overflow-hidden shadow rounded-lg">
                          <div className="p-5">
                            <div className="flex items-center">
                              <div className="flex-shrink-0">
                                <CreditCard className="h-6 w-6 text-gray-400" />
                              </div>
                              <div className="ml-5 w-0 flex-1">
                                <dl>
                                  <dt className="text-sm font-medium text-gray-500 truncate">
                                    Активные кредиты
                                  </dt>
                                  <dd className="text-lg font-medium text-gray-900">
                                    {stats.activeLoans}
                                  </dd>
                                </dl>
                              </div>
                            </div>
                          </div>
                        </div>

                        <div className="bg-white overflow-hidden shadow rounded-lg">
                          <div className="p-5">
                            <div className="flex items-center">
                              <div className="flex-shrink-0">
                                <DollarSign className="h-6 w-6 text-gray-400" />
                              </div>
                              <div className="ml-5 w-0 flex-1">
                                <dl>
                                  <dt className="text-sm font-medium text-gray-500 truncate">
                                    Общий долг
                                  </dt>
                                  <dd className="text-lg font-medium text-gray-900">
                                    ${stats.totalDebt.toLocaleString()}
                                  </dd>
                                </dl>
                              </div>
                            </div>
                          </div>
                        </div>

                        <div className="bg-white overflow-hidden shadow rounded-lg">
                          <div className="p-5">
                            <div className="flex items-center">
                              <div className="flex-shrink-0">
                                <DollarSign className="h-6 w-6 text-gray-400" />
                              </div>
                              <div className="ml-5 w-0 flex-1">
                                <dl>
                                  <dt className="text-sm font-medium text-gray-500 truncate">
                                    Месячные платежи
                                  </dt>
                                  <dd className="text-lg font-medium text-gray-900">
                                    ${stats.monthlyPayments.toLocaleString()}
                                  </dd>
                                </dl>
                              </div>
                            </div>
                          </div>
                        </div>
                      </>
                    );
                  })()}
                </div>

                {/* Quick Actions */}
                <div className="bg-white shadow rounded-lg p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Быстрые действия</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <Link
                      to={`/vehicles?company_id=${selectedCompany.id}`}
                      className="bg-indigo-50 border border-indigo-200 rounded-lg p-4 hover:bg-indigo-100 transition-colors"
                    >
                      <div className="flex items-center">
                        <Truck className="h-5 w-5 text-indigo-600 mr-2" />
                        <div>
                          <h4 className="font-medium text-indigo-900">Управление транспортом</h4>
                          <p className="text-sm text-indigo-700">Добавить и редактировать транспорт</p>
                        </div>
                      </div>
                    </Link>

                    <Link
                      to={`/vehicles/new?company_id=${selectedCompany.id}`}
                      className="bg-green-50 border border-green-200 rounded-lg p-4 hover:bg-green-100 transition-colors"
                    >
                      <div className="flex items-center">
                        <Plus className="h-5 w-5 text-green-600 mr-2" />
                        <div>
                          <h4 className="font-medium text-green-900">Добавить транспорт</h4>
                          <p className="text-sm text-green-700">Новый трак или трейлер</p>
                        </div>
                      </div>
                    </Link>
                    
                    <Link
                      to="/schedules/debt"
                      className="bg-blue-50 border border-blue-200 rounded-lg p-4 hover:bg-blue-100 transition-colors"
                    >
                      <h4 className="font-medium text-blue-900">График долгов</h4>
                      <p className="text-sm text-blue-700">Просмотр всех долговых обязательств</p>
                    </Link>
                    
                    <Link
                      to="/schedules/amortization"
                      className="bg-purple-50 border border-purple-200 rounded-lg p-4 hover:bg-purple-100 transition-colors"
                    >
                      <h4 className="font-medium text-purple-900">График погашения</h4>
                      <p className="text-sm text-purple-700">Детальный план выплат</p>
                    </Link>
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </main>
    </div>
  );
};

export default Dashboard; 