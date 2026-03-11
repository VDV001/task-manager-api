import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
// TODO: uncomment when router is created (Task 3)
// import router from './router'
import './assets/main.css'

const app = createApp(App)
app.use(createPinia())
// TODO: uncomment when router is created (Task 3)
// app.use(router)
app.mount('#app')
