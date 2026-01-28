import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
// if the file inside router folder is not named index.js, the import would be:
// import router from './router/someOtherFileName.js'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

/* notes:

this file does 3 things:

    1. creates the Vue application instance with "createApp(App)"
    2. integrates Pinia for state management with "app.use(createPinia())"
    3. sets up Vue Router for navigation with "app.use(router)"
    4. mounts the app to the DOM element with id 'app' using "app.mount('#app')" (it links the app to the HTML)

*/