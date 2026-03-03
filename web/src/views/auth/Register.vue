<template>
  <div class="register-container">
    <div class="register-card">
      <h1 class="register-title">注册 Alike</h1>
      <p class="register-subtitle">开启你的社交之旅</p>
      
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="phone">手机号</label>
          <input
            id="phone"
            v-model="formData.phone"
            type="tel"
            placeholder="请输入手机号"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            required
            minlength="6"
          />
        </div>
        
        <div class="form-group">
          <label for="nickname">昵称</label>
          <input
            id="nickname"
            v-model="formData.nickname"
            type="text"
            placeholder="请输入昵称"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="bio">个人简介</label>
          <textarea
            id="bio"
            v-model="formData.bio"
            placeholder="介绍一下自己..."
            rows="3"
          ></textarea>
        </div>
        
        <button type="submit" class="register-btn" :disabled="isLoading">
          {{ isLoading ? '注册中...' : '注册' }}
        </button>
      </form>
      
      <div class="register-footer">
        <p>已有账号？<router-link to="/login">立即登录</router-link></p>
      </div>
      
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const formData = ref({
  phone: '',
  password: '',
  nickname: '',
  bio: ''
})

const isLoading = ref(false)
const errorMessage = ref('')

const handleRegister = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    
    const result = await userStore.register({
      phone: formData.value.phone,
      password: formData.value.password,
      nickname: formData.value.nickname,
      bio: formData.value.bio
    })
    
    if (result.success) {
      router.push('/')
    } else {
      errorMessage.value = result.message || '注册失败'
    }
  } catch (error) {
    errorMessage.value = '注册失败，请重试'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.register-card {
  background: white;
  border-radius: 12px;
  padding: 2.5rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

.register-title {
  font-size: 2rem;
  font-weight: bold;
  color: #667eea;
  text-align: center;
  margin: 0 0 0.5rem 0;
}

.register-subtitle {
  text-align: center;
  color: #666;
  margin: 0 0 2rem 0;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #333;
}

.form-group input,
.form-group textarea {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  font-family: inherit;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.register-btn {
  padding: 0.875rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.3s;
}

.register-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.register-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-footer {
  margin-top: 1.5rem;
  text-align: center;
  color: #666;
}

.register-footer a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.register-footer a:hover {
  text-decoration: underline;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #fee;
  color: #c33;
  border-radius: 6px;
  text-align: center;
}
</style>