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
    <div class="update-username">
      <LoadingSpinner v-if="loading"></LoadingSpinner>
      <div class="update-form">
        <h3 v-if="message" class="update-message">{{ message }}</h3>
        <form @submit.prevent="Update_username" class="update-form__form">
          <label for="username" class="update-form__label">Current username: {{ this.username }}</label>
          <input
            type="text"
            id="new_username_box"
            v-model="new_username"
            required
            minlength="3"
            maxlength="20"
            class="update-form__input"
          />
          <button type="submit" class="update-form__button">
            Update username <svg class="feather">
            </svg>
          </button>
        </form>
        <p v-if="username_update_completed" class="update-message">Update completed</p>
      </div>
    </div>
  </template>
  
  <style scoped>
  .update-username {
    /* Add your styles for the main container here */
    background-color: #f5f5f5; /* Light gray background */
    border: 1px solid #ddd; /* Thin border */
    padding: 20px; /* Padding for inner elements */
    border-radius: 5px; /* Rounded corners */
    box-shadow: 0px 2px 5px rgba(0, 0, 0, 0.1); /* Subtle shadow */
    width: 400px; /* Set a fixed width */
    margin: 0 auto; /* Center the container horizontally */
  }
  
  .update-form {
    display: flex;
    flex-direction: column;
    gap: 10px; /* Spacing between form elements */
  }
  
  .update-message,
  .update-form__label {
    font-weight: bold;
    color: #333; /* Darker text color */
    font-family: monospace; /* Monospaced font for CSV look */
  }
  
  .update-form__input,
  .update-form__button {
    border: none; /* Remove default border */
    background-color: #fff; /* White background for input/button */
    padding: 8px; /* Consistent padding */
    border-radius: 3px; /* Rounded corners */
    box-shadow: inset 0px 1px 2px rgba(0, 0, 0, 0.1); /* Subtle inset shadow */
    transition: all 0.2s ease-in-out; /* Add smooth transitions on hover */
  }
  
  .update-form__input:hover,
  .update-form__button:hover {
    box-shadow: inset 0px 2px 5px rgba(0, 0, 0, 0.2); /* Increase shadow on hover */
  }
  
  .update-form__button {
    cursor: pointer; /* Indicate clickable button */
    background-color: #333; /* Darker background for button */
    color: #fff; /* White text color for button */
  }
  </style>
  