import { createApp } from 'vue'
import './assets/styles/base.css'
import App from './App.vue'

import pinia from './store'

const app = createApp(App)
app.use(pinia)
app.mount('#app')
