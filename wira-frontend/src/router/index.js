import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '@/views/LoginView.vue';
import DashboardView from '@/views/DashboardView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true },
    },
    {
      path: '/logout',
      name: 'logout',
      beforeEnter(to, from, next) {
        // Clear token and redirect to login
        localStorage.removeItem('wira_token');
        next('/login');
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('wira_token');

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Redirect unauthenticated users to login
    next('/login');
  } else if (to.name === 'login' && isAuthenticated) {
    // Redirect authenticated users away from login
    next('/dashboard');
  } else {
    next();
  }
});

export default router;
