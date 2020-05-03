<template>
  <b-container class="mt-3">
    <h1>Create Post</h1>
    <b-form>
      <b-form-group>
        <label>Title</label>
        <b-input v-model="title"></b-input>
      </b-form-group>
      <b-button variant="primary" @click="submitPost">Submit</b-button>
    </b-form>

    <hr>

    <h1>Posts</h1>
    <b-row>
      <b-col cols="3" v-for="post in posts" :key="post.id">
        <PostCard :post="post"></PostCard>
      </b-col>
    </b-row>

  </b-container>
</template>

<script>
  import axios from 'axios'
  import PostCard from "@/components/PostCard";
  import {POST_SERVICE, QUERY_SERVICE} from "@/consts";

  export default {
    components: {PostCard},

    data() {
      return {
        posts: [],
        title: '',
      }
    },

    async mounted() {
      const res = await axios.get(`${QUERY_SERVICE}/posts`)
      this.posts = res.data
    },

    methods: {
      submitPost() {
        axios.post(`${POST_SERVICE}/posts/create`, {title: this.title})
      },

    }

  }
</script>