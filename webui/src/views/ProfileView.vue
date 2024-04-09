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
				await this.$axios.put(`/photos/${post.post_id}/likes/me`, {}, {
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
				await this.$axios.delete(`/photos/${post.post_id}/likes/me`, {
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
				this.profile_data.data.is_following=true
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
				this.profile_data.data.is_following=false
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
				this.profile_data.data.is_banned = true
				this.profile_data.data.is_following = false
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
				this.profile_data.data.is_banned = false
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
	  <div v-if="loggedIn" class="profile-container">
		<h3 v-if="message != ''" class="alert alert-primary">{{ message }}</h3>
		<h1>{{ $route.params.username }}</h1>
		<div v-if="isItMe">
		  <router-link to="/settings" class="profile-settings">
			<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
			Profile settings
		  </router-link>
		</div>
		<p>This is the profile page of {{ $route.params.username }}</p>
		<p v-if="message != ''">Error: {{ message }} </p>
		<div v-if="profile_data">
		  <div v-if="!isItMe">
			<button v-if="!profile_data.data.is_following && !profile_data.data.is_banned" @click="follow()" class="btn btn-dark">Follow</button>
			<button v-if="profile_data.data.is_following" @click="unfollow()" class="btn btn-dark">Unfollow</button>
			<button v-if="!profile_data.data.is_banned" @click="ban" class="btn btn-dark">Ban</button>
			<button v-if="profile_data.data.is_banned" @click="unban" class="btn btn-dark">Unban</button>
		  </div>
		  <p class="followers-count">Number of followers: {{profile_data.data.followers_count}}</p>
		  <button v-if="!show_followers" @click="(show_followers = true)" class="btn btn-dark">Show followers</button>
		  <button v-if="show_followers" @click="(show_followers = false)" class="btn btn-dark">Hide followers</button>
		  <p v-if="show_followers" class="followers-list">Followers:
			<label v-for="username in profile_data.data.followers" :key="username">
			  <router-link :to="'/users/' + username"> @{{ username }} </router-link>
			</label>
		  </p>
		  <div v-if="profile_data.data.posts" class="post-container">
			<h3 class="photos-heading">Photos:</h3>
			<div v-for="(post, post_index) in profile_data.data.posts" :key="post.post_id" class="post">
			  <div class="post-content">
				<p v-if="post.content != 'undefined'">{{ post.content }}</p>
				<p>Likes: {{ post.likes_count }}</p>
				<button v-if="!post.is_liked" @click="likePost(post)" class="btn btn-dark btn-like">Like</button>
				<button v-if="post.is_liked" @click="unLikePost(post)" class="btn btn-dark btn-like">un-Like</button>
			  </div>
			  <img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
			  <button v-if="isItMe" @click="removePost(post_index)" class="btn btn-dark btn-delete">Delete Post</button>
			  <button v-if="!post.showCommentBox" @click="(post.showCommentBox = true)" class="btn btn-dark">Add Comment</button>
			  <div v-else>
				<textarea v-model="post.commentText" placeholder="Enter your comment"></textarea>
				<button @click="addComment(post)" class="btn btn-dark">Post Comment</button>
			  </div>
			  <div v-if="post.comments!=null" class="comments">
				<button v-if="!post.showComments" @click="post.showComments=true" class="btn btn-dark">Show Comments</button>
				<div v-else>
				  <h4 class="comments-heading">Comments:</h4>
				  <div v-for="(comment,comment_index) in post.comments" :key="comment.comment_id" class="comments">
					<div class="comment">
					  <p>
						<router-link :to="'/users/' + comment.username">@{{ comment.username }} </router-link>
						: {{ comment.content }}
					  </p>
					  <button v-if="authToken.slice(7) == comment.user_id" @click="removeComment(post, comment_index)" class="btn btn-dark">Delete comment</button>
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
  
  <style scoped>
  .profile-container {
	background-color: #f5f5f5;
	border: 1px solid #ddd;
   border-radius: 5px;
	padding: 20px;
	width: 80%;
	margin: 0 auto;
  }
  
  .profile-settings {
	color: #337ab7;
	text-decoration: none;
  }
  
  .profile-settings svg {
	height: 1em;
	width: 1em;
	margin-right: 5px;
  }
  
  .alert {
	background-color: #d4edda;
	color: #155724;
	padding: 10px;
	margin-bottom: 20px;
  }
  
  h1 {
	margin-top: 0;
  }
  
  p {
	color: #333;
  }
  
  .btn {
	display: inline-block;
	padding: 5px 10px;
	margin-bottom: 10px;
	color: #fff;
	background-color: #337ab7;
	border: none;
	border-radius: 3px;
	cursor: pointer;
  }
  
  .btn-dark {
	background-color: #333;
  }
  
  .btn-dark:hover {
	background-color: #555;
  }
  
  .btn-like {
	margin-right: 10px;
  }
  
  .btn-delete {
	background-color: red;
  }
  .btn-delete:hover {
	background-color: rgb(251, 94, 94);
  }
  
  .followers-count {
	margin-top: 20px;
  }
  
  .followers-list {
	margin-top: 10px;
  }
  
  .post-container {
	margin-top: 20px;
  }
  
  .photos-heading {
	margin-top: 0;
  }
  
  .post {
	background-color: #f9f9f9;
	border: 1px solid #ddd;
	border-radius: 5px;
	padding: 10px;
	margin-bottom: 10px;
  }
  
  .post-content {
	margin-bottom: 10px;
  }
  
  .post-image {
    max-width: 100%;
    margin-bottom: 10px;
  }
  
  .comments-heading {
	margin-top: 10px;
  }
  
  .comments {
	margin-top: 10px;
  }
  
  .comment {
	background-color: #f5f5f5;
	border: 1px solid #ddd;
	border-radius: 5px;
	padding: 10px;
	margin-bottom: 10px;
  }
  
  textarea {
	width: 100%;
	height: 80px;
	resize: none;
	margin-bottom: 10px;
  }
  
  .router-link-active {
	color: #337ab7;
  }
  </style>