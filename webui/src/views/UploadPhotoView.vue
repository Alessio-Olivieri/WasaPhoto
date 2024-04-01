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
							this.message = "Error handling ban request"
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
	<div>
		<div>
			<h1>So... {{ this.authUser }}, you want to upload a picture?</h1>
            <div class="row" style="margin: 10px;">
                <textarea class="form-control" id="picture-text-box" placeholder="Enter a description for your picture" v-model="caption"></textarea>
            </div>
            <div>
                <h3 v-if="message!=''">{{ message }}</h3>
                <input type="file" id="picture" accept=".png,.jpg,image/jpeg" @change="onFileChange"/>
                <button @click="uploadPicture">Upload Picture</button>
            </div>  
        </div>
        <div v-if="response_post != null" >
            <h2>Image uploaded:</h2>
            <p>User: {{ response_post.data.username }}&#10;</p>
            <p v-if="response_post.data.content">Text: {{ response_post.data.content }}&#10;</p>
            <p v-if="response_post.data.image!=null"><img v-bind:src="dataURI" alt="Post Image"></p>
        </div>
	</div>
</template>

<style>
img {
  max-width: 1000px;
  max-height: 700px;
}

</style>

