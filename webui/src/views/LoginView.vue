<script>
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            userID: null,
            username: "",
        }
    },
    methods: {
        async login() {
            this.loading = true;
            try {
                console.log("Logging in with username: " + this.username);
                const response = await this.$axios.post('/login', { username: this.username }); // Adjusted to match the Go backend
                console.log("Response: ", response.data);
                this.userID = response.data.userId; // Adjusted to correctly access the userId in the response
                this.saveTokenToSessionStorage(this.userID);
            } catch (error) {
                console.error("Error while logging in:", error);
                console.error("Error message:", error.message);
            }
            this.loading = false;
        },

        saveTokenToSessionStorage(userID) {
            const bearerToken = `Bearer ${userID}`;
            sessionStorage.setItem('authToken', bearerToken);
        },
    },
}
</script>

<template>
    <div class="login-container">
        <LoadingSpinner v-if="loading"></LoadingSpinner>
        <div class="login-form">
            <h2>Autenticati bastardo</h2>
            <form @submit.prevent="login">
                <label class="login-label" for="username">Username:</label>
                <input type="text" id="username" v-model="username" required minlength="3" maxlength="20"
                    style="padding: 6px;" />
                <button type="submit" class="btn btn-sm btn-outline-primary"
                    style="padding: 8px; float: right; font-size: 20px;">Login <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#key" />
                    </svg></button>
            </form>
            <div v-if="userID !== null">
                <p>Login successful! User identifier: {{ userID }}</p>
            </div>
        </div>
    </div>
</template>

<style>
</style>
