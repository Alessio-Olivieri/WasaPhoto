<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			authToken: null,
			loggedIn: true,
			page_number: 0,
			message: "",
			stream_data: null,
		}
	},
	methods: {
        async getStream(page){
			this.loading = true
			console.log(`getting stream of ${sessionStorage.getItem('username')}, at page ${page}...`)
			try{
				this.stream_data = await this.$axios.get(`/stream`, {
					params: {
							page: page,
						},
						headers: {
							'Authorization': this.authToken,
						},
					});
				console.log("page number " + this.stream_data.data.page_number)
				this.page_number = this.stream_data.data.page_number
				for (const post of this.stream_data.data.posts){
					post.showCommentBox = false
					post.showComments=false
					post.CommentsPage=0
					if (post.comments == null){
						post.comments = []
					}
				}
				this.message = ""
				console.log("stream data retrieved succesfully")
			} catch(error){
				if (error.response){
				const statusCode = error.response.status;
                    switch (statusCode) {
						case 400:
							console.error("Bad request")
							this.message = "bad request"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 404:
							console.error("Stream is empty")
							this.message = "Nothing to see here..."
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
				} else {
					console.error("error: ", error)
				}
			}
			this.loading = false
        },

		async addComment(post){
			this.loading = true
			try{
			post.showCommentBox = false
			console.log("adding comment...")
			const response = await this.$axios.post(`/photos/${post.post_id}/comments/`, {
				content: post.commentText
			}, {
				headers: {
						'Authorization': this.authToken,
					},
				});
			post.showComments = true
			post.comments.push(response.data)
			this.message = ""
			} catch(error){
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: post_id not valid or content is empty")
							this.message = "post_id not valid or content is empty"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: you're banned")
							this.message = "Forbidden: you're banned"
							break
						case 404:
							console.error("post not found")
							this.message = "post not found"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
			}
			this.loading = false
		},

		async removeComment(post, i){
			this.loading = true
			try{
				console.log("adding comment...")
				await this.$axios.delete(`/photos/${post.post_id}/comments/${post.comments[i].comment_id}`, {
					headers: {
						'Authorization': this.authToken
					}
				});
				post.comments.splice(i, 1)
				this.message = ""
			} catch(error){
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: post_id or comment_id not valid")
							this.message = "post_id or comment_id not valid"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: you're not the owner of the comment")
							this.message = "Forbidden: you're not the owner of the comment"
							break
						case 404:
							console.error("comment not found")
							this.message = "comment not found"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
			}
			this.loading = false
		},

		async likePost(post){
			this.loading = true
			try{
				console.log("liking post...")
				await this.$axios.put(`/photos/${post.post_id}/likes/${sessionStorage.getItem("username")}`, {}, {
					headers: {
							'Authorization': this.authToken,
						},
					});
				post.likes_count = post.likes_count + 1
				post.is_liked = true
				this.message = ""
				} catch(error){
					const statusCode = error.response.status;
						switch (statusCode) {
							case 400:
								console.error("Bad request: post_id not valid")
								this.message = "post_id not valid"
								break
							case 401:
								console.error("Auth header missing")
								this.message = "Unauthorized"
								break
							case 403:
								console.error("Forbidden: you're banned")
								this.message = "Forbidden: you're banned"
								break
							case 404:
								console.error("post not found")
								this.message = "post not found"
								break
							default:
								console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
								this.message = "Error handling ban request"
						}
				}

			this.loading = false
		},

		async unLikePost(post){
			this.loading = true
			try{
				console.log("unliking post...")
				await this.$axios.delete(`/photos/${post.post_id}/likes/${sessionStorage.getItem("username")}`, {
					headers: {
						'Authorization': this.authToken
					}
				});
				post.likes_count = post.likes_count - 1
				post.is_liked = false
				this.message = ""
			} catch(error){
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: post_id not valid")
							this.message = "post_id not valid"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: you're banned")
							this.message = "Forbidden: you're banned"
							break
						case 404:
							console.error("post not found or like not found")
							this.message = "post not found or like not found"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
			}
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
		<h2 v-if="message!=''">{{ message }}</h2>
		<div v-if="stream_data!=null" class="post-container">
			<button v-if="page_number!=0" @click="getStream(stream_data.data.page_number-1)">previous page</button>
			<button v-if="page_number!=0 && page_number-3>0" @click="getStream(stream_data.data.page_number-3)">{{ page_number-3 }}</button>
			<button v-if="page_number!=0 && page_number-2>0" @click="getStream(stream_data.data.page_number-2)">{{ page_number-2 }}</button>
			<button v-if="page_number!=0 && page_number-1>0" @click="getStream(stream_data.data.page_number-1)">{{ page_number-1 }}</button>
			{{ page_number }} 
			<button v-if="page_number+1<=stream_data.data.number_of_pages" @click="getStream(stream_data.data.page_number+1)">{{ page_number+1 }}</button>
			<button v-if="page_number+2<=stream_data.data.number_of_pages" @click="getStream(stream_data.data.page_number+2)">{{ page_number+2 }}</button>
			<button v-if="page_number+3<=stream_data.data.number_of_pages" @click="getStream(stream_data.data.page_number+3)">{{ page_number+3 }}</button>
			<button v-if="page_number+1<=stream_data.data.number_of_pages" @click="getStream(page_number+1)">next page</button>
			<div v-for="post in stream_data.data.posts" :key="post.post_id" class="post">
				<div class="post-content">
					<router-link :to="'/users/' + post.username">@{{ post.username }} </router-link>
					<p v-if="post.content != 'undefined'">{{ post.content }}</p>
					<p>Likes: {{ post.likes_count }}</p>
					<button v-if="!post.is_liked" @click="likePost(post)">Like</button>
					<button v-if="post.is_liked" @click="unLikePost(post)">un-Like</button>
				</div>
				<img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
				<div v-if="!post.showCommentBox">
					<button @click="(post.showCommentBox = true)">Add Comment</button>
				</div>
				<div v-else>
					<textarea v-model="post.commentText" placeholder="Enter your comment"></textarea>
					<button @click="addComment(post)">Post Comment</button>
				</div>
				<div v-if="post.comments!=null" class="comments">
					<button v-if="!post.showComments" @click="post.showComments=true">Show Comments</button>
					<div v-else>
						<h4 class="comments-heading">Comments:</h4>
						<div v-for="(comment,comment_index) in post.comments" :key="comment.comment_id" class="comments">
							<div class="comment">
								<p>
									<router-link :to="'/users/' + comment.username">@{{ comment.username }} </router-link>
									: {{ comment.content }}
								</p>
								<button v-if="authToken.slice(7) == comment.user_id" @click="removeComment(post, comment_index)">Delete comment</button>
							</div>
						</div>
					</div>
				</div>
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

  
