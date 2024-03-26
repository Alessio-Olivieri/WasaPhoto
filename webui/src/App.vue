
<script>
const authToken = sessionStorage.getItem('authToken');
import { RouterLink, RouterView} from 'vue-router'
export default {
	methods: {
		logout() {
			localStorage.clear();
			sessionStorage.clear();
			location.reload();
		},

		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

	},

	data() {
		return {
			username: sessionStorage.getItem('username'),
			mypath: "/users/" + authToken,
			streampath: "/users/" + authToken + "/stream/",
		}
	},

	mounted() {
		console.log("username: ", this.username)
		if (this.username) {
			this.$router.push('/users/' + this.username)
		}
		else {this.$router.push('/login')}
		console.log("authToken: ", sessionStorage.getItem('authToken'))
		console.log("username: ", this.username)
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Example App</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item" v-if=this.username>
							<RouterLink to="/stream" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Stream
							</RouterLink>
						</li>
						<li class="nav-item" v-if=this.username>
							<RouterLink :to="'/users/' + this.username" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								profile of {{ username }}
							</RouterLink>
						</li>
						<li class="nav-item" v-if=this.username>
							<RouterLink :to="'/photos/'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Upload Picture
							</RouterLink>
						</li>
						<li class="nav-item" v-if="!this.username">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Login
							</RouterLink>
						</li>
						<li class="nav-item" v-if=this.username>
							<RouterLink to="/users/" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#search" />
								</svg>
								Search user
							</RouterLink>
						</li>
						<li class="nav-item" v-if=this.username>
							<button class="nav-link" @click="logout">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Logout
							</button>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style></style>
