<script setup lang="ts">
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import FloatLabel from 'primevue/floatlabel'
import type { User } from '@/types/user'
import { ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import axios from 'axios'
import store from '@/store'

interface LogInResponse {
  user: User
  access_token: string
  refresh_token: string
}

const username = ref<string>('')
const password = ref<string>('')
const errorMessages = ref<string>('')
const toast = useToast()


const handleLogin = async () => {
  try {
    const response = await axios.post<LogInResponse>(
      '/login_user',
      {
        username: username.value,
        password: password.value
      },
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: 'none',
          'Access-Control-Allow-Origin': '*'
        }
      })
    store.setUser(response.data.user, response.data.access_token, response.data.refresh_token)
    toast.add({
      severity: 'success',
      summary: `Welcome ${response.data.user.username}`,
      detail: 'Login Successful',
      life: 3000
    })
  } catch (error: any) {
  if (error.response && error.response.status === 401) {
    errorMessages.value = error.response.data.message
  } else {
    errorMessages.value = 'An error occurred. Please try again later'
  }
  toast.add({
    severity: 'error',
    summary: 'Login Failed',
    detail: errorMessages.value,
    life: 3000
  })

  }
}

</script>

<template>
<div class="flex flex-column row-gap-5">
  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-user"></i>
    </InputGroupAddon>
    <FloatLabel>
      <InputText id="username" v-model="username" />
      <label for="username">Username</label>
    </FloatLabel>
  </InputGroup>
  <InputGroup>
    <InputGroupAddon>
      <i class="pi pi-lock"></i>
    </InputGroupAddon>
    <FloatLabel>
      <InputText id="password" type="password" v-model="password"/>
      <label for="password">Password</label>
    </FloatLabel>
  </InputGroup>

  <Button label="Login" @click="handleLogin"/>
</div>
</template>

<style scoped>

</style>