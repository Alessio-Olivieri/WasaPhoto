<script>
const authToken = sessionStorage.getItem('authToken');
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			loggedIn: true,
			profile_data: null,
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
			this.profile_data = await this.$axios.get(`/users/${this.$route.params.username}`, {
                    headers: {
                        'Authorization': authToken,
                    },
                });
			console.log("user profile data retrieved succesfully")
			this.loading = false
		},

		async follow() {
			this.loading = true
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
			this.loading = false
		},

		async unfollow(){
			this.loading = true
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
			this.loading = false
		},
		
		async ban(){
			this.loading = true
			console.log("following user...")
			const response = await this.$axios.put(`/banned/${this.$route.params.username}`, {}, {
				headers: {
                        'Authorization': this.authToken,
                    },
                });
			this.refresh()
			this.loading = false
		},

		async unban(){
			this.loading = true
			console.log("unfollowing user...")
			await this.$axios.delete(`/banned/${this.$route.params.username}`, {
				headers: {
					'Authorization': this.authToken
				}
			});
			this.refresh()
			this.loading = false
		},
		async likePost(post){
			this.loading = true
			console.log("liking post...")
			await this.$axios.put(`/photos/${post.post_id}/likes/${sessionStorage.getItem("username")}`, {}, {
				headers: {
                        'Authorization': this.authToken,
                    },
                });
			post.likes_count = post.likes_count + 1
			post.is_liked = true
			// this.refresh()
			this.loading = false
		},

		async unLikePost(post){
			this.loading = true
			console.log("unliking post...")
			await this.$axios.delete(`/photos/${post.post_id}/likes/${sessionStorage.getItem("username")}`, {
				headers: {
					'Authorization': this.authToken
				}
			});
			post.likes_count = post.likes_count - 1
			post.is_liked = false
			this.loading = false
		},

		async showTextBox(post_id){
			console.log("showing text box for post_id: ", post_id);
			this.commentBox_to_show[post_id] = true;
		},
		async removeCommentBox(post_id) {
			delete this.commentBox_to_show[post_id];
		},

		async addComment(post_id, commentText){
			this.loading = true
			this.removeCommentBox(post_id)
			console.log("adding comment...")
			const response = await this.$axios.post(`/photos/${post_id}/comments/`, {
				content: commentText
			}, {
				headers: {
						'Authorization': this.authToken,
					},
				});
			this.refresh()
			this.loading = false
		},

		async removeComment(post_id, comment_id){
			this.loading = true
			console.log("adding comment...")
			await this.$axios.delete(`/photos/${post_id}/comments/${comment_id}`, {
				headers: {
					'Authorization': this.authToken
				}
			});
			this.refresh()
			this.loading = false
		}
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
		}

	},	

}
</script>

<template>
	<div>
		<div v-if="loggedIn">
			<h1>{{ $route.params.username }}</h1>
			<div v-if="isItMe">
				Change username Button
			</div>
			<p>This is the profile page of {{ $route.params.username }}</p>
			<div v-if="profile_data">
				<div v-if="!isItMe">
					<button v-if="!profile_data.data.isFollowing && !profile_data.data.isBanned" @click="follow()">Follow</button>
					<button v-if="profile_data.data.isFollowing" @click="unfollow()">Unfollow</button>
					<button v-if="!profile_data.data.isBanned" @click="ban">Ban</button>
					<button v-if="profile_data.data.isBanned" @click="unban">Unban</button>
				</div>
				Number of followers: {{profile_data.data.followers_count}}
				<p>Followers:
					<label v-for="username in profile_data.data.followers" :key="username">{{ username }} </label>
				</p>
				<div v-if="profile_data.data.posts" class="post-container">
					<h3 class="photos-heading">Photos:</h3>
					<div v-for="post in profile_data.data.posts" :key="post" class="post">
						<div class="post-content">
							<p v-if="post.content">{{ post.content }}</p>
							<p>Likes: {{ post.likes_count }}</p>
							<button v-if="!post.is_liked" @click="likePost(post)">Like</button>
							<button v-if="post.is_liked" @click="unLikePost(post)">un-Like</button>
						</div>
						<div v-if="!post.showTextBox">
							<img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
							<button v-if="!this.commentBox_to_show.hasOwnProperty(post.post_id)" @click="showTextBox(post.post_id)">Add Comment</button>
							{{commentBox_to_show}}
						</div>
						<div v-if="this.commentBox_to_show.hasOwnProperty(post.post_id)">
							<textarea v-model="post.commentText" placeholder="Enter your comment"></textarea>
							<button @click="addComment(post.post_id, post.commentText)">Post Comment</button>
						</div>
						<div v-if="post.comments" class="comments">
							<h4 class="comments-heading">Comments:</h4>
							<div v-for="comment in post.comments" :key="comment" class="comments">
								<div class="comment">
									<p>From {{ comment.username }} : {{ comment.content }}</p>
									<button @click="removeComment(post.post_id, comment.comment_id)">Delete comment</button>
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
  
