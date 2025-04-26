import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import { faBell as faBellSolid, faMoon, faSun, faExpand, faPlusCircle, faTimes, faMapLocationDot } from '@fortawesome/free-solid-svg-icons'
import { faBell as faBellRegular } from '@fortawesome/free-regular-svg-icons';

library.add(faBellSolid, faBellRegular, faMoon, faSun, faExpand, faPlusCircle, faTimes, faMapLocationDot)

const app = createApp(App)

app.component('font-awesome-icon', FontAwesomeIcon)
app.use(createPinia())

app.mount('#app')


