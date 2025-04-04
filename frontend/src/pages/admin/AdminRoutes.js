import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import AdminPage from './AdminPage';
import CoursesAdmin from './CoursesAdmin';
import StatsDashboard from './StatsDashboard';
import MessagesAdmin from './MessagesAdmin';
import ProfileAdmin from './ProfileAdmin';
import ProtectedRoute from '../../components/ProtectedRoute';

const AdminRoutes = () => {
  return (
    <ProtectedRoute>
      <Routes>
        <Route path="/" element={<AdminPage />}>
          <Route index element={<Navigate to="stats" replace />} />
          <Route path="courses" element={<CoursesAdmin />} />
          <Route path="stats" element={<StatsDashboard />} />
          <Route path="messages" element={<MessagesAdmin />} />
          <Route path="profile" element={<ProfileAdmin />} />
        </Route>
      </Routes>
    </ProtectedRoute>
  );
};

export default AdminRoutes;