import React from 'react'
import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom'
import Dashboard from './components/Dashboard'
import Login from './components/Login'
import UserProfile from './components/UserProfile'
import CompanyForm from './components/companies/CompanyForm'
import CompanyList from './components/companies/CompanyList'
import VehicleForm from './components/vehicles/VehicleForm'
import VehicleList from './components/vehicles/VehicleList'
import { AuthProvider, useAuth } from './context/AuthContext'
import { CompanyProvider } from './context/CompanyContext'
import { VehicleProvider } from './context/VehicleContext'

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth();
  
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Загрузка...</div>
      </div>
    );
  }
  
  return user ? <>{children}</> : <Navigate to="/login" />;
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth();
  
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Загрузка...</div>
      </div>
    );
  }
  
  return user ? <Navigate to="/dashboard" /> : <>{children}</>;
};

function App() {
  return (
    <AuthProvider>
      <CompanyProvider>
        <VehicleProvider>
          <Router>
            <div className="App">
              <Routes>
                <Route 
                  path="/login" 
                  element={
                    <PublicRoute>
                      <Login />
                    </PublicRoute>
                  } 
                />
                <Route 
                  path="/dashboard" 
                  element={
                    <ProtectedRoute>
                      <Dashboard />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/profile" 
                  element={
                    <ProtectedRoute>
                      <UserProfile />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/companies" 
                  element={
                    <ProtectedRoute>
                      <CompanyList />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/companies/new" 
                  element={
                    <ProtectedRoute>
                      <CompanyForm />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/companies/:id/edit" 
                  element={
                    <ProtectedRoute>
                      <CompanyForm />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/vehicles" 
                  element={
                    <ProtectedRoute>
                      <VehicleList />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/vehicles/new" 
                  element={
                    <ProtectedRoute>
                      <VehicleForm />
                    </ProtectedRoute>
                  } 
                />
                <Route 
                  path="/vehicles/:id/edit" 
                  element={
                    <ProtectedRoute>
                      <VehicleForm />
                    </ProtectedRoute>
                  } 
                />
                <Route path="/" element={<Navigate to="/dashboard" />} />
              </Routes>
            </div>
          </Router>
        </VehicleProvider>
      </CompanyProvider>
    </AuthProvider>
  );
}

export default App;
