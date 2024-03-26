<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			authToken: null,
			loggedIn: true,
			stream_data: null,
		}
	},
	methods: {
        async getStream(page){
			this.loading = true
			console.log(`getting stream of ${sessionStorage.getItem('username')}`)
			this.stream_data = await this.$axios.get(`/stream`, {
				params: {
                        page: page,
                    },
                    headers: {
                        'Authorization': this.authToken,
                    },
                });
			console.log("stream data retrieved succesfully")
			this.loading = false
        },

		async likePost(post_id){
			this.loading = true
			console.log("liking post...")
			const response = await this.$axios.put(`/photos/${post_id}/likes/${sessionStorage.getItem("username")}`, {}, {
				headers: {
                        'Authorization': this.authToken,
                    },
                });
			this.loading = false
		},

		async unLikePost(post_id){
			this.loading = true
			console.log("unliking post...")
			await this.$axios.delete(`/photos/${post_id}/likes/${sessionStorage.getItem("username")}`, {
				headers: {
					'Authorization': this.authToken
				}
			});
			this.loading = false
		},
	},
	mounted() {
		try {
			this.authToken = sessionStorage.getItem('authToken');
		} catch (e) {
			this.loggedIn = false;
		}

		console.log("authToken: '", this.authToken, "'",
					"authUsername: ", sessionStorage.getItem('username'),
					"loggedIn: ", this.loggedIn);
		this.getStream(0)
	},
	computed: {
    	username() {
      	return sessionStorage.getItem("username");
    	}
	}

}
</script>

<template>
	<div>
		<h1>This is the stream of {{ username }}</h1>
		<div v-if="this.stream_data">
			<div v-for="post in this.stream_data.data" :key="post" class="post">
				{{ post.username }}
				<div class="post-content">
					<p v-if="post.content">{{ post.content }}</p>
					<button v-if="!post.is_liked" @click="likePost(post.post_id)">Like</button>
					<button v-if="post.is_liked" @click="unLikePost(post.post_id)">un-Like</button>
				</div>
					<img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
			</div>
		</div>
	</div>
  </template>
  
  <style>
  /* Basic styles for the post container */
  .post-container {
	margin: 1rem 0;
	border: 1px solid #ddd;
	padding: 1rem;
	border-radius: 5px;
  }
  
  /* Styles for the photos heading */
  .photos-heading {
	font-weight: bold;
	margin-bottom: 0.5rem;
  }
  
  /* Styles for each post */
  .post {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 1rem;
  border: 1px solid #ddd; /* Add border */
  padding: 1rem; /* Add padding */
  border-radius: 5px; /* Add border radius */
}
  
  /* Styles for post image with fixed size and centering */
  .post-image {
	width: 150px;
	height: 150px;
	object-fit: cover;
	margin-bottom: 1rem;
	border: 1px solid #3f8317;
  }
  
  /* Styles for post content area */
  .post-content {
	text-align: center;
  }
  
  /* Styles for comments heading */
  .comments-heading {
	font-weight: bold;
	margin-top: 0.5rem;
  }
  
  /* Styles for individual comments */
  .comments {
	margin-bottom: 0.25rem;
  }

  .comment{
	  border: 1px solid #ddd;
	  padding: 0.5rem;
	  border-radius: 5px;
  }
  </style>
  
