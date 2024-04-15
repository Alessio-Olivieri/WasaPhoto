<script>
export default {
    data: function() {
        return {
            username: "",
            errormsg: null,
            loading: false,
            userId: null,
            message: ""
        }
    },
    methods: {
        async login() {
            this.loading = true;
            if (this.username.length < 3) {
                console.error("Username must be at least 3 characters long");
                return;
            }
            try {
                console.log("Logging in with username: " + this.username);
                const response = await this.$axios.post('/login', {}, {
                    params: {
                        username: this.username,
                    }
                },);
                console.log("Response: ", response.data);
                this.userId = response.data.userId; 
                this.saveToSessionStorage();
                console.log("Token saved correctly")
            } catch (error) {
                console.error("Error while logging in:", error);
                console.error("Error message:", error.message);
            }
            this.loading = false;
            //console.log("pushing userId in $router")
            //this.$router.push(toString(this.username));
            location.reload();
        },

        saveToSessionStorage() {
            const bearerToken = `Bearer ${this.userId}`;
            sessionStorage.setItem('authToken', bearerToken);
            sessionStorage.setItem('username', this.username)
        },
    },
}
</script>

<template>
	<div class="login-container">
		<h3 v-if="message != ''" class="alert alert-primary">{{ message }}</h3>
		<LoadingSpinner v-if="loading"></LoadingSpinner>
		<div class="login-form">
			<h2>Autenticati.</h2>
			<form @submit.prevent="login">
				<label class="login-label" for="username">Username:</label>
				<input type="text" id="username" v-model="username" required minlength="3" maxlength="20"
					style="padding: 6px;" />
				<button type="submit" class="btn btn-sm btn-outline-primary"
					style="padding: 8px; float: right; font-size: 20px;">
					Login 
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#key" />
					</svg>
				</button>
			</form>
			<div v-if="userId != null">
				<p>Logged in as {{ username }}, Welcome!</p>
			</div>
		</div>
	</div>
</template>


<style>
</style>
