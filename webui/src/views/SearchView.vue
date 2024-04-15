<script>
const authToken = sessionStorage.getItem('authToken');
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            userId: null,
            search_result : null,
            username: "",
            message: ""
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
	<h2 class="alert alert-primary" v-if="message">{{ message }}</h2>
	<div class="search-container">
		<h2>Search.</h2>
		<form @submit.prevent="search" class="search-form">
			<label for="username" class="search-form__label">User to search:</label>
			<input
				type="text"
				id="username"
				v-model="username"
				required
				minlength="3"
				maxlength="20"
				class="search-form__input"
				/>
			<button type="submit" class="search-form__button">
				Search 
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#search" />
				</svg>
			</button>
		</form>
		<div v-if="search_result != null" class="search-results">
			<p>Search completed</p>
			<ul class="search-results__list">
				<li v-for="username in search_result" :key="username" class="search-results__item">
					<router-link :to="'/users/' + username">{{ username }}</router-link>
				</li>
			</ul>
		</div>
	</div>
</template>

  
  <style scoped>
  .search-container {
    /* Add your styles for the main container here */
    background-color: #f5f5f5; /* Light gray background */
    border: 1px solid #ddd; /* Thin border */
    padding: 20px; /* Padding for inner elements */
    border-radius: 5px; /* Rounded corners */
    width: 400px; /* Set a fixed width */
    margin: 0 auto; /* Center the container horizontally */
  }
  
  .search-form {
    display: flex;
    gap: 10px; /* Spacing between form elements */
    align-items: center; /* Align label and input vertically */
  }
  
  .search-form__label {
    font-weight: bold;
    color: #333; /* Darker text color */
  }
  
  .search-form__input,
  .search-form__button {
    border: none; /* Remove default border */
    background-color: #fff; /* White background for input/button */
    padding: 8px; /* Consistent padding */
    border-radius: 3px; /* Rounded corners */
    box-shadow: inset 0px 1px 2px rgba(0, 0, 0, 0.1); /* Subtle inset shadow */
    transition: all 0.2s ease-in-out; /* Add smooth transitions on hover */
  }
  
  .search-form__input:hover,
  .search-form__button:hover {
    box-shadow: inset 0px 2px 5px rgba(0, 0, 0, 0.2); /* Increase shadow on hover */
  }
  
  .search-form__button {
    cursor: pointer; /* Indicate clickable button */
    background-color: #333; /* Darker background for button */
    color: #fff; /* White text color for button */
    padding: 10px 15px; /* Adjust padding for better button size */
  }
  
  .search-results {
    /* Add styles for the search results container here */
    margin-top: 20px; /* Spacing after the form */
  }
  
  .search-results__list {
    list-style: none; /* Remove default bullet points */
    padding: 0; /* Reset padding for better spacing */
  }
  
  .search-results__item {
    margin-bottom: 5px; /* Spacing between results */
  }
  
  .search-results__item a { /* Style the router-link */
    color: #333; /* Darker text color for links */
    text-decoration: none; /* Remove underline */
    transition: color 0.2s ease-in-out; /* Smooth color transition on hover */
  }
  
  .search-results__item a:hover {
    color: #007bff; /* Blue color on hover */
  }
  </style>
  
