import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import PublicLayout from '../layouts/PublicLayout.vue'
import MainLayout from '../layouts/MainLayout.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import FeedView from '../views/FeedView.vue'
import ProfileView from '../views/ProfileView.vue'
import SettingsView from '../views/SettingsView.vue'
import CreatePostView from '../views/CreatePostView.vue'
import NotificationsView from '../views/NotificationsView.vue'
import GroupsBrowse from '../views/GroupsBrowse.vue'
import GroupDetails from '../views/GroupDetails.vue'

const routes = [
  { path: '/', redirect: '/login' },
  {
    path: '/login',
    component: LoginView,
    meta: { layout: PublicLayout, requiresAuth: false }
  },
  {
    path: '/register',
    component: RegisterView,
    meta: { layout: PublicLayout, requiresAuth: false }
  },
  {
    path: '/feed',
    component: FeedView,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/profile',
    component: ProfileView,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/settings',
    component: SettingsView,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/notifications',
    component: NotificationsView,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/create-post',
    component: CreatePostView,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/groups',
    name: 'Groups',
    component: GroupsBrowse,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/groups/:id',
    name: 'GroupDetails',
    component: GroupDetails,
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/friends',
    name: 'friends',
    component: () => import('@/views/FriendsView.vue'),
    meta: { layout: MainLayout, requiresAuth: true }
  },
  {
    path: '/post/:id',
    name: 'Post',
    component: () => import('@/views/PostView.vue'),
    meta: { layout: MainLayout, requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  // Wait until auth state is checked
  if (!authStore.checked) {
    await authStore.checkAuth()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return '/login'
  }

  if (!to.meta.requiresAuth && authStore.isAuthenticated &&
      (to.path === '/login' || to.path === '/register')) {
    return '/feed'
  }
})

export default router
