import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

// Import specific icons
import { faBell as faBellSolid } from '@fortawesome/free-solid-svg-icons'
import { faBell as faBellRegular } from '@fortawesome/free-regular-svg-icons'; // Regular (outline) bell icon

// Add icons to the library
library.add(faBellSolid, faBellRegular)

const app = createApp(App)

// Register the component globally
app.component('font-awesome-icon', FontAwesomeIcon)

app.mount('#app')

