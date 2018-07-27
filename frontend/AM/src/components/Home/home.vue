<template>
<div>
  <router-link ></router-link>
</div>
</template>

<script>
import { getplaylistrecom } from 'mock/api'

export default {
  data () {
    return {
      data: ''
    }
  },
  created () {
    getplaylistrecom()
      .then(response => {
        this.data = response.data
      })
      .catch(function (error) {
        console.log(error)
      })
  }
}
</script>

<style scoped>
</style>
