import { Navigate, Route, Routes } from 'react-router-dom';
import { Layout } from '@/components/Layout';
import { HomePage } from '@/pages/HomePage';
import { LoginPage } from '@/pages/LoginPage';
import { RegisterPage } from '@/pages/RegisterPage';
import { ProductDetailPage } from '@/pages/ProductDetailPage';
import { SeckillDetailPage } from '@/pages/SeckillDetailPage';
import { OrdersPage } from '@/pages/OrdersPage';
import { SellerPage } from '@/pages/SellerPage';

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="product/:id" element={<ProductDetailPage />} />
        <Route path="seckill/:id" element={<SeckillDetailPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="seller" element={<SellerPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
