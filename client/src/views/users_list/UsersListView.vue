<script setup lang="ts">
import $api from '@/packages/api/client';
import { shallowRef } from 'vue';

const shit = shallowRef({
  id: 0,
  date: '',
});
const count = shallowRef<number>(1);

async function findUsers() {
  try {
    const response = await $api('/bot_users', {
      params: {
        limit: 100,
        page: count.value,
        query: 'Elena',
        next_cursor_date: shit.value.date,
        next_cursor_id: shit.value.id,
        age_from: 18,
        age_till: 25,
      },
    });
    shit.value.date = response.data[response.data.length - 1].join_date;
    shit.value.id = response.data[response.data.length - 1].u_id;
    count.value += 1;

    console.log(response);
  } catch (e) {
    console.error(e);
  }
}
</script>

<template>
  <div>
    USERS
    <button @click="findUsers">
      FIND
    </button>
  </div>
</template>

<style scoped>

</style>
