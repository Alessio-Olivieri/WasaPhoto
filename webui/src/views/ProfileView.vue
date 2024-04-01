<script>
const authToken = sessionStorage.getItem('authToken');
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			loggedIn: true,
			profile_data: null,
			show_followers: false,
			message:"",
			commentBox_to_show: {},
		}
	},
	watch: {
		'$route.params.username': function() {
			if (this.$route.params.username != null) {
				this.refresh();
			}
  		}
	},

	methods: {
		async refresh() {
			console.log("Refreshing data...");
			this.profile_data = await this.$axios.get(`/users/${this.$route.params.username}`, {
                    headers: {
                        'Authorization': authToken,
                    },
                });
		},

		async getProfile(){
			this.loading = true
			console.log("getting user profile...")
			try {
				this.profile_data = await this.$axios.get(`/users/${this.$route.params.username}`, {
                    headers: {
                        'Authorization': authToken,
                    },
                })
			}
			catch (error) {
				if (error.response){
				const statusCode = error.response.status;
                    switch (statusCode) {
						case 400:
							console.error("Bad request: username is not valid or empty")
							this.message = "username is not valid or empty"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: user is banned")
							this.message = "User is banned"
							break
						case 404:
							console.error("user with requested username not exists")
							this.message = "user with requested username not exists"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
				} else {
					console.error("error: ", error)
				}
			}
			console.log("user profile data retrieved succesfully, parsing data...")
			if (this.profile_data.data.posts == null){
				this.profile_data.data.posts = []
			}
			for (const post of this.profile_data.data.posts){
				post.showCommentBox = false
				post.showComments=false
				post.CommentsPage=0
				if (post.comments == null){
					post.comments = []
				}
			}
			console.log("user profile data parsed succesfully")
			this.loading = false
		},

		async removePost(i){
			this.loading = true
			console.log("adding comment...")
			try{
			await this.$axios.delete(`/photos/${this.profile_data.data.posts[i].post_id}`, {
				headers: {
					'Authorization': this.authToken
				}
			});
			this.message = ""
			} catch(error){
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: post_id not valid")
							this.message = "Error handling ban request"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: you're not the owner of the post")
							this.message = "Forbidden: you're not the owner of the post"
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
			this.profile_data.data.posts.splice(i, 1)
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

		async follow() {
			this.loading = true
			try{
				console.log("following user...")
				const response = await this.$axios.put(`/followed/${this.$route.params.username}`, {}, {
					headers: {
							'Authorization': this.authToken,
						},
					});
				this.profile_data.data.followers_count = this.profile_data.data.followers_count + 1
				this.profile_data.data.isFollowing=true
				if (this.profile_data.data.followers == null){
					this.profile_data.data.followers = []
				}
				this.profile_data.data.followers.push(sessionStorage.getItem('username'))
				this.message = ""
			} catch(error) {
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: profile username not present or bad format or you banned the user")
							this.message = "Error handling ban request"
							break
						case
							401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: user has banned you")
							this.message = "Forbidden: user has banned you"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
							this.message = "Error handling ban request"
					}
			}
			this.loading = false
		},

		async unfollow(){
			this.loading = true
			try{
				console.log("unfollowing user...")
				await this.$axios.delete(`/followed/${this.$route.params.username}`, {
					headers: {
						'Authorization': this.authToken
					}
				});
				this.profile_data.data.followers_count = this.profile_data.data.followers_count - 1
				this.profile_data.data.isFollowing=false
				this.profile_data.data.followers.splice(
					this.profile_data.data.followers.findIndex(
						follower => follower === sessionStorage.getItem('username')
						), 1);
				this.message = ""
			} catch (error) {
				const statusCode = error.response.status;
					switch (statusCode) {
						case 400:
							console.error("Bad request: profile username not present or user doesn't exists")
							this.message = "Error handling ban request"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 403:
							console.error("Forbidden: user has banned you")
							this.message = "Forbidden: user has banned you"
							break
						case 404:
							console.error("Already unfollowed")
							this.message = "Already unfollowed"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error);
							this.message = "Error handling ban request"
					}
			}
			this.loading = false
		},
		
		async ban(){
			this.loading = true
			try{
				const response = await this.$axios.put(`/banned/${this.$route.params.username}`, {}, {
					headers: {
							'Authorization': this.authToken,
						},
					});
				this.profile_data.data.isBanned = true
				this.profile_data.data.isFollowing = false
				this.message = ""
			} catch (error) {
				const statusCode = error.response.status;
                    switch (statusCode) {
						case 400:
							console.error("Bad request: profile username not present or user doesn't exists")
							this.message = "Error handling ban request"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						default:
							console.error(`Unhandled HTTP Error (${statusCode}):`, error);
							this.message = "Error handling ban request"
					}
			}
			this.loading = false
		},

		async unban(){
			this.loading = true
			try{
				await this.$axios.delete(`/banned/${this.$route.params.username}`, {
					headers: {
						'Authorization': this.authToken
					}
				});
				this.profile_data.data.isBanned = false
				this.message = ""
			} catch (error) {
				if (error.response){
				const statusCode = error.response.status;
                    switch (statusCode) {
						case 400:
							console.error("Bad request: profile username not present or user doesn't exists")
							this.message = "profile username not present or user doesn't exists"
							break
						case 401:
							console.error("Auth header missing")
							this.message = "Unauthorized"
							break
						case 404:
							console.error("user is not banned")
							this.message = "user is not banned"
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
	},
	mounted() {
		try {
			this.authToken = sessionStorage.getItem('authToken');
		} catch (e) {
			this.loggedIn = false;
		}

		console.log("authToken: '", this.authToken, "'",
					"authUser: ", sessionStorage.getItem('username'),
					"loggedIn: ", this.loggedIn);

		if (this.loggedIn) {
		this.getProfile();
		}
	},
	computed:{
		isItMe: function(){
			return sessionStorage.getItem('username') == this.$route.params.username;
		},

	},	

}
</script>

<template>
	<div>
		<h3 v-if="message != ''">{{ message }}</h3>
		<div v-if="loggedIn">
			<h1>{{ $route.params.username }}</h1>
			<div v-if="isItMe">
				<RouterLink to="/settings">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
					Profile settings
				</RouterLink>			
			</div>
			<p>This is the profile page of {{ $route.params.username }}</p>
			<p v-if="message != ''">Error: {{ message }} </p>
			<div v-if="profile_data">
				<div v-if="!isItMe">
					<button v-if="!profile_data.data.isFollowing && !profile_data.data.isBanned" @click="follow()">Follow</button>
					<button v-if="profile_data.data.isFollowing" @click="unfollow()">Unfollow</button>
					<button v-if="!profile_data.data.isBanned" @click="ban">Ban</button>
					<button v-if="profile_data.data.isBanned" @click="unban">Unban</button>
				</div>
				Number of followers: {{profile_data.data.followers_count}}
				<button v-if="!show_followers" @click="(show_followers = true)">show followers</button>
				<button v-if="show_followers" @click="(show_followers = false)">hide followers</button>
				<p v-if="show_followers">Followers:
					<label v-for="username in profile_data.data.followers" :key="username">
						<router-link :to="'/users/' + username"> @{{ username }} </router-link>
					</label>
				</p>
				<div v-if="profile_data.data.posts" class="post-container">
					<h3 class="photos-heading">Photos:</h3>
					<div v-for="(post, post_index) in profile_data.data.posts" :key="post.post_id" class="post">						<div class="post-content">
							<p v-if="post.content != 'undefined'">{{ post.content }}</p>
							<p>Likes: {{ post.likes_count }}</p>
							<button v-if="!post.is_liked" @click="likePost(post)">Like</button>
							<button v-if="post.is_liked" @click="unLikePost(post)">un-Like</button>
						</div>
						<img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
						<button v-if="isItMe" @click="removePost(post_index)"> Delete Post</button>
						<button v-if="!post.showCommentBox" @click="(post.showCommentBox = true)">Add Comment</button>
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
  
