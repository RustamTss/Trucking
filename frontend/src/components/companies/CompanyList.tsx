import { Building2, Edit, Plus, Trash2 } from 'lucide-react';
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useCompanies } from '../../context/CompanyContext';
import BackButton from '../common/BackButton';

const CompanyList: React.FC = () => {
  const { companies, deleteCompany, isLoading } = useCompanies();
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const handleDelete = async (id: string, name: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить компанию "${name}"? Это действие нельзя отменить.`)) {
      return;
    }

    setDeletingId(id);
    try {
      await deleteCompany(id);
    } catch (error) {
      console.error('Ошибка удаления компании:', error);
      alert('Ошибка при удалении компании');
    } finally {
      setDeletingId(null);
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent, action: () => void) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      action();
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-center py-12">
            <div className="text-lg text-gray-500">Загрузка компаний...</div>
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
              <Building2 className="h-8 w-8 text-indigo-600 mr-3" />
              <div>
                <h1 className="text-3xl font-bold text-gray-900">Управление компаниями</h1>
                <p className="text-gray-600">
                  Всего компаний: {companies.length}
                </p>
              </div>
            </div>
            <Link
              to="/companies/new"
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              aria-label="Добавить новую компанию"
            >
              <Plus className="w-4 h-4 mr-2" />
              Добавить компанию
            </Link>
          </div>
        </div>

        {/* Company List */}
        {companies.length === 0 ? (
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
                aria-label="Создать первую компанию"
              >
                <Building2 className="-ml-1 mr-2 h-5 w-5" />
                Добавить компанию
              </Link>
            </div>
          </div>
        ) : (
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {companies.map((company) => (
                <li key={company.id} className="px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center">
                        <Building2 className="h-5 w-5 text-gray-400 mr-3" />
                        <div>
                          <h3 className="text-lg font-medium text-gray-900 truncate">
                            {company.name}
                          </h3>
                          <div className="mt-1 text-sm text-gray-500 space-y-1">
                            <p>EIN: {company.ein}</p>
                            <p>{company.address}</p>
                            {company.phone && <p>Телефон: {company.phone}</p>}
                            {company.email && <p>Email: {company.email}</p>}
                          </div>
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center space-x-2">
                      {/* Edit Button */}
                      <Link
                        to={`/companies/${company.id}/edit`}
                        className="inline-flex items-center p-2 border border-gray-300 rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                        aria-label={`Редактировать компанию ${company.name}`}
                      >
                        <Edit className="h-4 w-4" />
                      </Link>

                      {/* Delete Button */}
                      <button
                        onClick={() => handleDelete(company.id, company.name)}
                        onKeyDown={(e) => handleKeyDown(e, () => handleDelete(company.id, company.name))}
                        disabled={deletingId === company.id}
                        className="inline-flex items-center p-2 border border-red-300 rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                        tabIndex={0}
                        aria-label={`Удалить компанию ${company.name}`}
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};

export default CompanyList; 