import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import SearchView from '../views/SearchView.vue'
import StreamView from '../views/StreamView.vue'
import UploadPhotoView from '../views/UploadPhotoView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/stream', component: StreamView},
		{path: '/login', component: LoginView},
		{path: '/users/:username', component: ProfileView},
		{path: '/users/', component: SearchView},
		{path: '/photos/', component: UploadPhotoView},

	]
})

export default router
