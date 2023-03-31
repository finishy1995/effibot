import { createApp } from 'vue'
// @ts-ignore
import router from "./router"
import App from './App.vue'
// element UI
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './main.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { jsMind } from 'jsmind'

// pinia
import { createPinia } from "pinia";

const pinia = createPinia();

const app = createApp(App)
app.use(router).use(pinia).use(jsMind).use(ElementPlus).mount('#app')


// 全局图标组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}