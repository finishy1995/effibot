import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        component: () => import('@/views/chat/Mind.vue')
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router