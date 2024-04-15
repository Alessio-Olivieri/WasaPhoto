<script>
import pageCounter from '../components/PageCounter.vue'
export default {
	components: {
		pageCounter
	},
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
							this.message = "Error handling  request"
					}
				} else {
					console.error("error: ", error)
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
								this.message = "Error handling  request"
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
							this.message = "Error handling  request"
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
	<div class="general">
		<h1>This is the stream of {{ username }}</h1>
		<h3 v-if="message != ''" class="alert alert-primary">{{ message }}</h3>
		<div v-if="stream_data!=null" class="stream-container">
			<!-- <pageCounter :page_number="page_number" :number_of_pages="stream_data.data.number_of_pages"></pageCounter>  -->
			<div class="btn-group btn-group-toggle" data-toggle="buttons">
				<label v-if="page_number!=0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-1)">
					previous page
				</label>
				<label v-if="page_number!=0 && page_number-3>=0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-3)">
					{{ page_number-3 }}
				</label>
				<label v-if="page_number!=0 && page_number-2>=0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-2)">
					{{ page_number-2 }}
				</label>
				<label v-if="page_number!=0 && page_number-1>=0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-1)">
					{{ page_number-1 }}
				</label>
				<label class="btn btn-dark active">{{ page_number }} </label>
				<label v-if="page_number+1<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+1)">
					{{ page_number+1 }}
				</label>
				<label v-if="page_number+2<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+2)">
					{{ page_number+2 }}
				</label>
				<label v-if="page_number+3<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+3)">
					{{ page_number+3 }}
				</label>
				<label v-if="page_number+1<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(page_number+1)">
					next page
				</label>
			</div>
			<div v-for="post in stream_data.data.posts" :key="post.post_id" class="post-container">
				<div class="post-data ">
					<router-link :to="'/users/' + post.username">@{{ post.username }} </router-link>
					<p v-if="post.content != 'null'">{{ post.content }}</p>
					<p>Likes: {{ post.likes_count }}</p>
					<button v-if="!post.is_liked" @click="likePost(post)" class="btn btn-dark">Like</button>
					<button v-if="post.is_liked" @click="unLikePost(post)" class="btn btn-dark">un-Like</button>
				</div>
				<img v-if="post.image" :src="`data:image/png;base64,${post.image}`" alt="Post Image" class="post-image">
				<Comments :post=post></Comments>
			</div>
			<div class="btn-group btn-group-toggle" data-toggle="buttons">
				<label v-if="page_number!=0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-1)">
					previous page
				</label>
				<label v-if="page_number!=0 && page_number-3>0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-3)">
					{{ page_number-3 }}
				</label>
				<label v-if="page_number!=0 && page_number-2>0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-2)">
					{{ page_number-2 }}
				</label>
				<label v-if="page_number!=0 && page_number-1>0" class="btn btn-dark" @click="getStream(stream_data.data.page_number-1)">
					{{ page_number-1 }}
				</label>
				<label class="btn btn-dark active">{{ page_number }} </label>
				<label v-if="page_number+1<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+1)">
					{{ page_number+1 }}
				</label>
				<label v-if="page_number+2<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+2)">
					{{ page_number+2 }}
				</label>
				<label v-if="page_number+3<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(stream_data.data.page_number+3)">
					{{ page_number+3 }}
				</label>
				<label v-if="page_number+1<=stream_data.data.number_of_pages" class="btn btn-dark" @click="getStream(page_number+1)">
					next page
				</label>
			</div>
		</div>		
	</div>
  </template>
  
  <style scoped>
  .alert-primary {
    background-color: #007bff;
    color: #fff;
    padding: 10px;
    margin-bottom: 20px;
  }

  .general {
	text-align: center;
  }

  .btn-group {
    margin-bottom: 10px;
  }

  .btn-group label {
    margin-right: 5px;
  }

  .btn-group .btn-dark {
    background-color: #333;
  }
  .btn-group .btn-dark:hover {
    background-color: #555;
  }

  .stream-container {
    background-color: #f9f9f9;	text-align: center;

    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 10px;
    margin-bottom: 10px;
  }

  .post-container {
    background-color: #f9f9f9;
    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 10px;
    margin-bottom: 10px;
  }

  .post-data {
	margin-bottom: 10px;
  }

  .post-image {
    max-width: 100%;
    margin-bottom: 10px;
  }

  textarea {
	width: 100%;
	height: 80px;
	resize: none;
	margin-bottom: 10px;
  }
  
</style>

  
