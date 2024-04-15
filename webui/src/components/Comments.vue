<script>
export default {
    props: ['post'],
    data: function() {
		return {
			errormsg: null,
			loading: false,
			profile_data: null,
			message:"",
			commentBox_to_show: {},
            authToken:""
		}
	},
    methods: {
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
                    post.comments.unshift(response.data)
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
                                    this.message = "Error handling  request"
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
                                    this.message = "Error handling  request"
                            }
                    }
                    this.loading = false
                },
            },
            mounted() {
                this.authToken = sessionStorage.getItem('authToken')
            }
        }
</script>

<template>
    <div v-if="!post.showCommentBox">
        <button @click="(post.showCommentBox = true)" class="btn btn-dark">Add Comment</button>
    </div>
    <div v-else>
        <textarea v-model="post.commentText" placeholder="Enter your comment"></textarea>
        <button @click="addComment(post)" class="btn btn-dark">Post Comment</button>
    </div>
    <div v-if="post.comments!=null" class="comments">
        <button v-if="post.showComments" @click="(post.showComments = false)" class="btn btn-dark">Hide Comments</button>
        <button v-if="!post.showComments" @click="post.showComments=true" class="btn btn-dark">Show Comments</button>
        <div v-else>
            <h4 class="comments-heading">Comments:</h4>
            <div v-if="post.showComments" v-for="(comment,comment_index) in post.comments" :key="comment.comment_id" class="comments">
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
</template>

<style>
  .comments-heading {
    margin-top: 10px;
  }

  .comments .comment {
    background-color: #f5f5f5;
    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 10px;
    margin-bottom: 10px;
  }
</style>