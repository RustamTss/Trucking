import { Save, Truck } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useCompanies } from '../../context/CompanyContext';
import { useVehicles } from '../../context/VehicleContext';
import BackButton from '../common/BackButton';

interface FormData {
  company_id: string;
  type: 'truck' | 'trailer';
  vin: string;
  make: string;
  model: string;
  year: string;
  purchase_price: string;
  purchase_date: string;
  status: 'active' | 'inactive' | 'sold';
}

const VehicleForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const { companies } = useCompanies();
  const { vehicles, createVehicle, updateVehicle, isLoading } = useVehicles();
  
  const isEditing = Boolean(id);
  const existingVehicle = isEditing ? vehicles.find(v => v.id === id) : null;
  const preselectedCompanyId = searchParams.get('company_id') || '';

  const [formData, setFormData] = useState<FormData>({
    company_id: preselectedCompanyId,
    type: 'truck',
    vin: '',
    make: '',
    model: '',
    year: new Date().getFullYear().toString(),
    purchase_price: '0',
    purchase_date: new Date().toISOString().split('T')[0],
    status: 'active',
  });

  const [errors, setErrors] = useState<Partial<FormData>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (isEditing && existingVehicle) {
      setFormData({
        company_id: existingVehicle.company_id,
        type: existingVehicle.type,
        vin: existingVehicle.vin,
        make: existingVehicle.make,
        model: existingVehicle.model,
        year: existingVehicle.year.toString(),
        purchase_price: existingVehicle.purchase_price.toString(),
        purchase_date: existingVehicle.purchase_date.split('T')[0],
        status: existingVehicle.status,
      });
    }
  }, [isEditing, existingVehicle]);

  const validateForm = (): boolean => {
    const newErrors: Partial<FormData> = {};

    if (!formData.company_id) {
      newErrors.company_id = 'Выберите компанию';
    }

    if (!formData.vin.trim()) {
      newErrors.vin = 'VIN обязателен';
    } else if (!/^[A-HJ-NPR-Z0-9]{17}$/i.test(formData.vin)) {
      newErrors.vin = 'VIN должен содержать 17 символов (без I, O, Q)';
    }

    if (!formData.make.trim()) {
      newErrors.make = 'Марка обязательна';
    }

    if (!formData.model.trim()) {
      newErrors.model = 'Модель обязательна';
    }

    if (formData.year < '1900' || formData.year > (new Date().getFullYear() + 1).toString()) {
      newErrors.year = 'Некорректный год выпуска';
    }

    if (formData.purchase_price <= '0') {
      newErrors.purchase_price = 'Цена покупки должна быть больше 0';
    }

    if (!formData.purchase_date) {
      newErrors.purchase_date = 'Дата покупки обязательна';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      const vehicleData = {
        ...formData,
        year: parseInt(formData.year),
        purchase_price: parseFloat(formData.purchase_price),
        purchase_date: new Date(formData.purchase_date).toISOString(),
      };

      if (isEditing && id) {
        await updateVehicle(id, vehicleData);
      } else {
        await createVehicle(vehicleData);
      }

      // Возвращаемся к списку транспорта для выбранной компании
      navigate(`/vehicles?company_id=${formData.company_id}`);
    } catch (error) {
      console.error('Ошибка сохранения транспорта:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    let processedValue: any = value;
    if (name === 'year' || name === 'purchase_price') {
      processedValue = parseFloat(value) || 0;
    }
    
    setFormData(prev => ({ ...prev, [name]: processedValue }));
    
    // Очищаем ошибку при изменении поля
    if (errors[name as keyof FormData]) {
      setErrors(prev => ({ ...prev, [name]: undefined }));
    }
  };

  const getBackUrl = () => {
    if (formData.company_id) {
      return `/vehicles?company_id=${formData.company_id}`;
    }
    return '/dashboard';
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <BackButton to={getBackUrl()} className="mb-4" />
          <div className="flex items-center">
            <Truck className="h-8 w-8 text-indigo-600 mr-3" />
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                {isEditing ? 'Редактировать транспорт' : 'Добавить транспорт'}
              </h1>
              <p className="text-gray-600">
                {isEditing ? 'Обновите информацию о транспорте' : 'Заполните данные нового транспорта'}
              </p>
            </div>
          </div>
        </div>

        {/* Form */}
        <div className="bg-white shadow rounded-lg">
          <form onSubmit={handleSubmit} className="p-6 space-y-6">
            {/* Компания */}
            <div>
              <label htmlFor="company_id" className="block text-sm font-medium text-gray-700 mb-2">
                Компания *
              </label>
              <select
                id="company_id"
                name="company_id"
                value={formData.company_id}
                onChange={handleChange}
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                  errors.company_id ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                }`}
                disabled={isSubmitting}
                aria-describedby={errors.company_id ? 'company_id-error' : undefined}
              >
                <option value="">Выберите компанию</option>
                {companies?.map(company => (
                  <option key={company.id} value={company.id}>
                    {company.name}
                  </option>
                )) || []}
              </select>
              {errors.company_id && (
                <p id="company_id-error" className="mt-1 text-sm text-red-600">
                  {errors.company_id}
                </p>
              )}
            </div>

            {/* Тип транспорта */}
            <div>
              <label htmlFor="type" className="block text-sm font-medium text-gray-700 mb-2">
                Тип транспорта *
              </label>
              <select
                id="type"
                name="type"
                value={formData.type}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                disabled={isSubmitting}
              >
                <option value="truck">Трак</option>
                <option value="trailer">Трейлер</option>
              </select>
            </div>

            {/* VIN */}
            <div>
              <label htmlFor="vin" className="block text-sm font-medium text-gray-700 mb-2">
                VIN номер *
              </label>
              <input
                type="text"
                id="vin"
                name="vin"
                value={formData.vin}
                onChange={handleChange}
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 uppercase ${
                  errors.vin ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                }`}
                placeholder="1HGBH41JXMN109186"
                maxLength={17}
                disabled={isSubmitting}
                aria-describedby={errors.vin ? 'vin-error' : undefined}
              />
              {errors.vin && (
                <p id="vin-error" className="mt-1 text-sm text-red-600">
                  {errors.vin}
                </p>
              )}
            </div>

            {/* Марка и модель */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="make" className="block text-sm font-medium text-gray-700 mb-2">
                  Марка *
                </label>
                <input
                  type="text"
                  id="make"
                  name="make"
                  value={formData.make}
                  onChange={handleChange}
                  className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                    errors.make ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                  }`}
                  placeholder="Freightliner"
                  disabled={isSubmitting}
                  aria-describedby={errors.make ? 'make-error' : undefined}
                />
                {errors.make && (
                  <p id="make-error" className="mt-1 text-sm text-red-600">
                    {errors.make}
                  </p>
                )}
              </div>

              <div>
                <label htmlFor="model" className="block text-sm font-medium text-gray-700 mb-2">
                  Модель *
                </label>
                <input
                  type="text"
                  id="model"
                  name="model"
                  value={formData.model}
                  onChange={handleChange}
                  className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                    errors.model ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                  }`}
                  placeholder="Cascadia"
                  disabled={isSubmitting}
                  aria-describedby={errors.model ? 'model-error' : undefined}
                />
                {errors.model && (
                  <p id="model-error" className="mt-1 text-sm text-red-600">
                    {errors.model}
                  </p>
                )}
              </div>
            </div>

            {/* Год выпуска */}
            <div>
              <label htmlFor="year" className="block text-sm font-medium text-gray-700 mb-2">
                Год выпуска *
              </label>
              <input
                type="number"
                id="year"
                name="year"
                value={formData.year}
                onChange={handleChange}
                min="1900"
                max={new Date().getFullYear() + 1}
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                  errors.year ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                }`}
                disabled={isSubmitting}
                aria-describedby={errors.year ? 'year-error' : undefined}
              />
              {errors.year && (
                <p id="year-error" className="mt-1 text-sm text-red-600">
                  {errors.year}
                </p>
              )}
            </div>

            {/* Цена покупки */}
            <div>
              <label htmlFor="purchase_price" className="block text-sm font-medium text-gray-700 mb-2">
                Цена покупки ($) *
              </label>
              <input
                type="number"
                id="purchase_price"
                name="purchase_price"
                value={formData.purchase_price}
                onChange={handleChange}
                min="0"
                step="0.01"
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                  errors.purchase_price ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                }`}
                placeholder="85000.00"
                disabled={isSubmitting}
                aria-describedby={errors.purchase_price ? 'purchase_price-error' : undefined}
              />
              {errors.purchase_price && (
                <p id="purchase_price-error" className="mt-1 text-sm text-red-600">
                  {errors.purchase_price}
                </p>
              )}
            </div>

            {/* Дата покупки */}
            <div>
              <label htmlFor="purchase_date" className="block text-sm font-medium text-gray-700 mb-2">
                Дата покупки *
              </label>
              <input
                type="date"
                id="purchase_date"
                name="purchase_date"
                value={formData.purchase_date}
                onChange={handleChange}
                className={`w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 ${
                  errors.purchase_date ? 'border-red-300 focus:border-red-500' : 'border-gray-300 focus:border-indigo-500'
                }`}
                disabled={isSubmitting}
                aria-describedby={errors.purchase_date ? 'purchase_date-error' : undefined}
              />
              {errors.purchase_date && (
                <p id="purchase_date-error" className="mt-1 text-sm text-red-600">
                  {errors.purchase_date}
                </p>
              )}
            </div>

            {/* Статус */}
            <div>
              <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-2">
                Статус
              </label>
              <select
                id="status"
                name="status"
                value={formData.status}
                onChange={handleChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                disabled={isSubmitting}
              >
                <option value="active">Активный</option>
                <option value="inactive">Неактивный</option>
                <option value="sold">Продан</option>
              </select>
            </div>

            {/* Actions */}
            <div className="flex items-center justify-between pt-6 border-t border-gray-200">
              <BackButton to={getBackUrl()} />
              
              <button
                type="submit"
                disabled={isSubmitting || isLoading}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                aria-label={isEditing ? 'Сохранить изменения' : 'Создать транспорт'}
              >
                <Save className="w-4 h-4 mr-2" />
                {isSubmitting ? 'Сохранение...' : (isEditing ? 'Сохранить' : 'Создать')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default VehicleForm; 