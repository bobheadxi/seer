<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <NewTeam />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import NewTeam from '@/components/NewTeam.vue'; // @ is an alias to /src

@Component({
  components: {
    NewTeam,
  },
})
export default class Home extends Vue {}
</script>
