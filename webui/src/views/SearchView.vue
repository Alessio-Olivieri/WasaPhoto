<script>
const authToken = sessionStorage.getItem('authToken');
const username = sessionStorage.getItem('username');
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            userId: null,
            search_result : null,
        }
    },
    methods: {
        async search() {
            this.loading = true;
            try {
                console.log("Searching for: " + this.username);
                const response = await this.$axios.get('/users/', {
                    params: {
                        username: this.username,
                    },
                    headers: {
                      'Authorization': authToken,
                    },
                });
                console.log("Response: ", response.data);
                this.search_result =  response.data.users;
            } catch (error) {
                if (error.response){
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 401:
                            console.error("Unauthorized: user not logged in")
                            this.message = "user not logged in"
                            break
                        case 400:
                            console.error("Bad request: Search username wrong format")
                            this.message = "insert Search username"
                            break
                        case 404:
                            console.error("No username matches the search")
                            this.message = "No username matches the search"
                            break
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                            this.message = "error updating username"
                    }
                } else {
                    console.error("error: ", error)
                }
            }

            this.loading = false;
        },
    },
}
</script>

<template>
    <div class="login-container">
        <LoadingSpinner v-if="loading"></LoadingSpinner>
        <div class="login-form">
            <h2>Search.</h2>
            <form @submit.prevent="search">
                <label class="login-label" for="username">User to search:</label>
                <input type="text" id="username" v-model="username" required minlength="3" maxlength="20"
                    style="padding: 6px;" />
                <button type="submit" class="btn btn-sm btn-outline-primary"
                    style="padding: 8px; font-size: 15px;">Search <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#search" />
                    </svg></button>
            </form>
            <div v-if="search_result != null">
                <p>Search completed</p>
                <ul>
                    <li v-for="username in search_result" :key="username">
                        <router-link :to="'/users/' + username">{{ username }}</router-link>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<style>
</style>
