<script>
export default{
    data: function() {
            return {
                errormsg: null,
                loading: false,
                username: sessionStorage.getItem('username'),
                message: "",
                new_username: null,
                username_update_completed: false
            }
        },
    methods: {
        async Update_username(){
            this.loading=true
            console.log("changing username button pressed\n changing username")
            try {
            await this.$axios.put('/settings/username', { 'username' : this.new_username }, {
                headers: {
                    'Content-Type': 'application/json',  // Set content type for file upload
                    'Authorization': sessionStorage.getItem('authToken'),
                },
            });
            console.log("Username updated successfully!")
            this.username_update_completed = true
            sessionStorage.setItem('username', this.new_username)
            this.message = ""
            location.reload();
			} catch (error) {
				if (error.response){
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 401:
                            console.error("Unauthorized: user not logged in")
                            this.message = "user not logged in"
                            break
                        case 400:
                            console.error("Bad request: New username wrong format")
                            this.message = "insert new username"
                            break
                        case 409:
                            console.error("Conflict: username already exists")
                            this.message = "username already exists"
                            break
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                            this.message = "error updating username"
                    }
                } else {
                    console.error("error: ", error)
                }
            }
            this.loading = false
        }
    }
}

</script>


<template>
    <div class="login-container">
        <LoadingSpinner v-if="loading"></LoadingSpinner>
        <div class="login-form">
            <h3 v-if="message!=null">{{ message }}</h3>
            <form @submit.prevent="Update_username">
                <label class="login-label" for="username">Current username: {{ this.username }}</label>
                <input type="text" id="new_username_box" v-model="new_username" required minlength="3" maxlength="20"
                    style="padding: 6px;" />
                <button type="submit" class="btn btn-sm btn-outline-primary"
                    style="padding: 8px; font-size: 15px;">Update username <svg class="feather">
                    </svg></button>
            </form>
            <p v-if="username_update_completed">Search completed</p>
        </div>
    </div>
</template>