"use client";
import { ReactNode } from 'react';

interface ProtectedRouteProps {
  children: ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  // SEM AUTENTICAÇÃO - sempre renderiza o conteúdo
  return <>{children}</>;
}