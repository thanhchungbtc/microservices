<template>
  <b-card>
    <b-card-title>{{ post.title }}</b-card-title>
    <ul>
      <li v-for="comment in post.comments" :key="comment.id">{{ filteredContent(comment) }}</li>
    </ul>

    <div class="text-muted font-weight-bold">Comment</div>
    <b-form-group>
      <b-input placeholder="What is your though?" v-model="content"></b-input>
    </b-form-group>
    <b-button variant="primary" @click="submit">Submit</b-button>
  </b-card>
</template>
<script>
  import axios from 'axios'
  import {COMMENT_SERVICE} from "@/consts";

  export default {
    name: 'post-card',
    props: ['post'],

    data() {
      return {
        content: '',
      }
    },
    methods: {
      submit() {
        axios.post(`${COMMENT_SERVICE}/posts/${this.post.id}/comments`, {
          content: this.content,
        })
      },

      filteredContent(comment) {
        const status = comment.status
        if (status === 'pending') {
          return "pending"
        }
        if (status === 'reject') {
          return 'This comment is being rejected'
        }
        return comment.content
      }
    }
  }
</script>