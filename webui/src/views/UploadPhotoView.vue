<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
            authUser: null,
			authToken: null,
            picture: null,
            caption: null,
            imagePresent: false,
            captionPresent: false,
            response_post: null,
            message: "",
		}
	},
	methods: {
		async refresh() {
			console.log("Refreshing data...");
		},

        onFileChange(event) {
            this.picture = event.target.files[0];
        },

		async uploadPicture() {
			console.log("Uploading picture...");
			// Add picture using file reference
			const formData = new FormData();
            try{
  			    formData.append('picture', this.picture);  
                this.imagePresent = true; 
			}   
            catch (error) {
                console.warn("No picture selected")
                this.imagePresent = false;
            }
            
            //add caption using textarea
                formData.append('caption', this.caption);
            if (this.caption == null) {
                console.warn("No caption entered")
                this.captionPresent = false;
            } else { this.captionPresent = true; }
            
            // Check if image or caption are present
            if (this.imagePresent || this.captionPresent) {
                console.log("Uploading...")
            }
            else {
                console.warn("Cannot upload picture without both image and caption")
                return;
            }

            // Send the picture to the server
            try {
            this.response_post = await this.$axios.post('/photos/', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',  // Set content type for file upload
                    'Authorization': this.authToken,
                },
            });
            this.message = ""
			console.log("Picture uploaded successfully!")
			} catch (error) {
				if (error.response){
				const statusCode = error.response.status;
                    switch (statusCode) {
						case 400:
							console.error("Bad request: No content and no photo")
							this.message = "Insert content or photo"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling  request"
					}
				} else {
					console.error("error: ", error)
				}
			}
        }
	},
	mounted() {
        this.authToken = sessionStorage.getItem('authToken');
        this.authUser = sessionStorage.getItem('username');
        console.log("authToken: '", this.authToken, "'",
                    "authUser: ", this.authUser,);
	},
    computed: {
        dataURI() {
        return `data:image/png;base64,${this.response_post.data.image}`; // Replace "png" with the actual image format
        }
  }

}
</script>

<template>
	<div class="upload-container" style="background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin: 0 auto; max-width: 400px;">
		<h2 class="alert alert-primary" v-if="message">{{ message }}</h2>
		<div class="upload-form">
			<h1>So... {{ this.authUser }}, you want to upload a picture?</h1>
			<div class="form-group">
				<label for="picture-text-box" class="form-label">Picture description:</label>
				<textarea
					class="form-control"
					id="picture-text-box"
					placeholder="Enter a description for your picture"
					v-model="caption"
					></textarea>
			</div>
			<div class="form-group">
				<input type="file" id="picture" accept=".png,.jpg,image/jpeg" @change="onFileChange" required/>
				<button type="button" class="btn btn-primary" @click="uploadPicture">Upload Picture</button>
			</div>
		</div>
		<div v-if="response_post != null" class="upload-results">
			<h2>Image uploaded:</h2>
			<p>User: {{ response_post.data.username }}</p>
			<p v-if="response_post.data.content!='null'">Text: {{ response_post.data.content }}</p>
			<p v-if="response_post.data.image!=null">
				<img v-bind:src="dataURI" alt="Post Image" class="post-image"/>
			</p>
		</div>
	</div>
</template>

  
  <style>
  
    .upload-form {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
  
    .upload-form h1 {
      font-size: 1.5rem; /* Adjust heading size */
    }
  
    .form-label {
      font-weight: bold;
      margin-bottom: 5px; /* Add spacing after label */
    }
  
    .form-control,
    .btn-primary {
      border-radius: 3px;
      padding: 8px 10px; /* Adjust padding for better fit */
    }
  
    .btn-primary {
      background-color: #333;
      color: #fff;
    }
  
    .upload-results {
      margin-top: 20px;
      border-top: 1px solid #ddd; /* Add top border */
      padding-top: 10px; /* Add padding after border */
    }
    .post-image {
    max-width: 100%;
    margin-bottom: 10px;
  }
  </style>
  
