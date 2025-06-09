import { AlertTriangle, Edit2, Save, Trash2, User, X } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { usersAPI } from '../services/api';
import type { User as UserType } from '../types';

const UserProfile: React.FC = () => {
  const [user, setUser] = useState<UserType | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  const [editData, setEditData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    try {
      setIsLoading(true);
      const profileData = await usersAPI.getProfile();
      setUser(profileData);
      setEditData({
        name: profileData.name,
        email: profileData.email,
        password: '',
        confirmPassword: ''
      });
    } catch (err: any) {
      setError('Ошибка загрузки профиля');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async () => {
    if (!user) return;
    
    if (editData.password && editData.password !== editData.confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    try {
      setIsSaving(true);
      setError('');
      
      const updateData: { name?: string; email?: string; password?: string } = {};
      
      if (editData.name !== user.name) {
        updateData.name = editData.name;
      }
      
      if (editData.email !== user.email) {
        updateData.email = editData.email;
      }
      
      if (editData.password) {
        updateData.password = editData.password;
      }

      const updatedUser = await usersAPI.updateProfile(updateData);
      setUser(updatedUser);
      
      // Обновляем данные в localStorage
      localStorage.setItem('user', JSON.stringify(updatedUser));
      
      setIsEditing(false);
      setSuccess('Профиль успешно обновлен');
      setEditData(prev => ({ ...prev, password: '', confirmPassword: '' }));
      
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка обновления профиля');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!user) return;
    
    try {
      setIsDeleting(true);
      await usersAPI.delete(user.id);
      
      // Очищаем localStorage и перенаправляем на вход
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка удаления аккаунта');
      setIsDeleting(false);
    }
  };

  const handleCancel = () => {
    if (!user) return;
    
    setEditData({
      name: user.name,
      email: user.email,
      password: '',
      confirmPassword: ''
    });
    setIsEditing(false);
    setError('');
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">Профиль не найден</p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 flex items-center">
            <User className="mr-3 h-8 w-8 text-indigo-600" />
            Профиль пользователя
          </h2>
          
          {!isEditing && (
            <button
              onClick={() => setIsEditing(true)}
              className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
            >
              <Edit2 className="w-4 h-4 mr-2" />
              Редактировать
            </button>
          )}
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {success && (
          <div className="mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
            {success}
          </div>
        )}

        <div className="space-y-6">
          {/* Имя */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Имя
            </label>
            {isEditing ? (
              <input
                type="text"
                value={editData.name}
                onChange={(e) => setEditData(prev => ({ ...prev, name: e.target.value }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              />
            ) : (
              <p className="text-gray-900 py-2">{user.name}</p>
            )}
          </div>

          {/* Email */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            {isEditing ? (
              <input
                type="email"
                value={editData.email}
                onChange={(e) => setEditData(prev => ({ ...prev, email: e.target.value }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              />
            ) : (
              <p className="text-gray-900 py-2">{user.email}</p>
            )}
          </div>

          {/* Пароль (только при редактировании) */}
          {isEditing && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Новый пароль (оставьте пустым, если не хотите менять)
                </label>
                <input
                  type="password"
                  value={editData.password}
                  onChange={(e) => setEditData(prev => ({ ...prev, password: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Подтвердите новый пароль
                </label>
                <input
                  type="password"
                  value={editData.confirmPassword}
                  onChange={(e) => setEditData(prev => ({ ...prev, confirmPassword: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </div>
            </>
          )}

          {/* Дата создания */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Дата регистрации
            </label>
            <p className="text-gray-900 py-2">
              {new Date(user.created_at).toLocaleDateString('ru-RU')}
            </p>
          </div>

          {/* Кнопки */}
          {isEditing && (
            <div className="flex justify-between pt-4">
              <div className="space-x-3">
                <button
                  onClick={handleSave}
                  disabled={isSaving}
                  className="flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50"
                >
                  <Save className="w-4 h-4 mr-2" />
                  {isSaving ? 'Сохранение...' : 'Сохранить'}
                </button>
                
                <button
                  onClick={handleCancel}
                  className="flex items-center px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
                >
                  <X className="w-4 h-4 mr-2" />
                  Отмена
                </button>
              </div>

              <button
                onClick={() => setShowDeleteModal(true)}
                className="flex items-center px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
              >
                <Trash2 className="w-4 h-4 mr-2" />
                Удалить аккаунт
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Модальное окно удаления */}
      {showDeleteModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <div className="flex items-center mb-4">
              <AlertTriangle className="h-6 w-6 text-red-600 mr-3" />
              <h3 className="text-lg font-medium text-gray-900">
                Удаление аккаунта
              </h3>
            </div>
            
            <p className="text-gray-600 mb-6">
              Вы действительно хотите удалить свой аккаунт? Это действие нельзя отменить.
              Будут удалены все ваши компании, транспорт, кредиты и платежи.
            </p>
            
            <div className="flex justify-end space-x-3">
              <button
                onClick={() => setShowDeleteModal(false)}
                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
              >
                Отмена
              </button>
              
              <button
                onClick={handleDelete}
                disabled={isDeleting}
                className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50"
              >
                {isDeleting ? 'Удаление...' : 'Удалить'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserProfile; 