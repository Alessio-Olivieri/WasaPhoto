<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			userId: null,
			username: null,
			loggedIn: true
		}
	},
	methods: {
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

		async getProfile(){
			this.loading = true
			this.errormsg = null
			console.log("getting user profile...")
			const response = await this.$axios.get(`/${this.username}`, { username: this.username });
		}
	},
	mounted() {
		try{
			this.userId = sessionStorage.getItem('authToken').match(/^Bearer (.+)$/)[1]
			this.username = sessionStorage.getItem('username')}
		catch(e){
			this.loggedIn = false
		}
		console.log("authToken: ", this.userId, "username: ", this.username)

		if (this.loggedIn) {
			this.getProfile();
		}
		
	}
}
</script>

<template>
	<div>
		<div v-if="this.loggedIn">
			<h1>{{ this.username }}</h1>
            This is my profile page
        </div>
		<div v-else>
			<h1>Not logged in</h1>
			<p>You are not logged in, please login to see your profile</p>
		</div>
	</div>
</template>

<style>
</style>

