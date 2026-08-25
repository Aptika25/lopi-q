<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import OtpInput from '@/components/OtpInput.vue'

const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const isBackupMode = ref(false)
const otpCode = ref('')
const isResetSelfMode = ref(false)
const successMessage = ref('')

const toggleResetSelfMode = () => {
  isResetSelfMode.value = !isResetSelfMode.value
  isBackupMode.value = isResetSelfMode.value
  otpCode.value = ''
  errorMessage.value = ''
}

// Setup 2FA States
const showSetup2FA = ref(false)
const qrCodeUrl = ref('')
const totpSecret = ref('')
const setupStep = ref(1) // 1 = scan QR, 2 = show backup codes
const backupCodes = ref<string[]>([])
const setupCode = ref('')
const setupError = ref('')
const setupLoading = ref(false)

const init2faSetup = async () => {
  setupLoading.value = true
  setupError.value = ''
  showSetup2FA.value = true
  setupStep.value = 1
  setupCode.value = ''
  
  try {
    const data = await authStore.setup2fa()
    if (data && data.secret) {
      totpSecret.value = data.secret
      qrCodeUrl.value = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(data.otpauth_url || '')}`
    } else {
      setupError.value = 'Gagal menyiapkan data QR Code.'
    }
  } catch (err) {
    console.error(err)
    setupError.value = 'Gagal menghubungi server untuk konfigurasi 2FA.'
  } finally {
    setupLoading.value = false
  }
}

const handleConfirmSetup = async () => {
  if (!setupCode.value || setupCode.value.length !== 6 || isNaN(Number(setupCode.value))) {
    setupError.value = 'Kode OTP harus berupa 6 digit angka.'
    return
  }
  
  setupLoading.value = true
  setupError.value = ''
  
  try {
    const data = await authStore.enable2fa(setupCode.value, totpSecret.value)
    if (data && data.success) {
      backupCodes.value = data.backup_codes || []
      setupStep.value = 2
    } else {
      setupError.value = authStore.error || 'Kode OTP salah. Silakan coba lagi.'
    }
  } catch (err) {
    console.error(err)
    setupError.value = 'Gagal memverifikasi kode OTP.'
  } finally {
    setupLoading.value = false
  }
}

const handleSetupComplete = async () => {
  loading.value = true
  try {
    const res = await authStore.login(email.value, password.value)
    if (res.otpRequired) {
      showSetup2FA.value = false
      otpCode.value = ''
      isBackupMode.value = false
      isResetSelfMode.value = false
    } else if (res.success) {
      router.push('/admin')
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const copyBackupCodes = () => {
  const text = backupCodes.value.join('\n')
  navigator.clipboard.writeText(text)
  alert('Kode backup pemulihan berhasil disalin!')
}

const downloadBackupCodes = () => {
  const text = `KODE BACKUP PEMULIHAN LOPI-Q BULUKUMBA\nIdentifier: ${email.value}\nTanggal: ${new Date().toLocaleDateString()}\n\nSimpan kode ini baik-baik. Setiap kode hanya bisa digunakan sekali untuk menonaktifkan atau melewati 2FA:\n\n` + backupCodes.value.join('\n')
  const blob = new Blob([text], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `backup-codes-lopi-q-${email.value.replace(/[^a-zA-Z0-9]/g, '_')}.txt`
  a.click()
  URL.revokeObjectURL(url)
}

const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = 'Email atau NIP wajib diisi.'
    return
  }

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const res = await authStore.login(email.value, password.value)
    if (res.otpSetupRequired) {
      otpCode.value = ''
      isBackupMode.value = false
      isResetSelfMode.value = false
      await init2faSetup()
    } else if (res.otpRequired) {
      otpCode.value = ''
      isBackupMode.value = false
      isResetSelfMode.value = false
    } else if (res.success) {
      if (authStore.isAdmin) {
        router.push('/admin')
      } else {
        router.push('/calltaker/scan')
      }
    } else {
      errorMessage.value = authStore.error || 'Email/NIP atau kata sandi tidak valid.'
    }
  } catch (err: any) {
    console.error(err)
    errorMessage.value = 'Gagal terhubung ke server autentikasi. Silakan coba beberapa saat lagi.'
  } finally {
    loading.value = false
  }
}

const handleSelfReset2fa = async () => {
  if (!otpCode.value || otpCode.value.length !== 8) {
    errorMessage.value = 'Kode backup harus 8 karakter alfanumerik.'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const success = await authStore.selfReset2fa(otpCode.value)
    if (success) {
      authStore.otpRequired = false
      isResetSelfMode.value = false
      isBackupMode.value = false
      otpCode.value = ''
      successMessage.value = '2FA berhasil dinonaktifkan secara mandiri. Silakan login kembali menggunakan password.'
    } else {
      errorMessage.value = authStore.error || 'Gagal melakukan reset mandiri. Pastikan kode backup Anda benar.'
    }
  } catch (err: any) {
    console.error(err)
    errorMessage.value = 'Gagal melakukan reset 2FA secara mandiri.'
  } finally {
    loading.value = false
  }
}

const handleVerify2fa = async () => {
  if (isResetSelfMode.value) {
    await handleSelfReset2fa()
    return
  }

  if (!otpCode.value) {
    errorMessage.value = isBackupMode.value ? 'Kode backup wajib diisi.' : 'Kode OTP wajib diisi.'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const success = await authStore.verify2fa(otpCode.value)
    if (success) {
      if (authStore.isAdmin) {
        router.push('/admin')
      } else {
        router.push('/calltaker/scan')
      }
    } else {
      errorMessage.value = authStore.error || 'Kode verifikasi salah atau kedaluwarsa.'
    }
  } catch (err: any) {
    console.error(err)
    errorMessage.value = 'Gagal verifikasi kode OTP. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div class="min-h-screen w-full bg-background flex items-center justify-center relative overflow-hidden font-body-lg text-on-background select-none p-gutter">
    <!-- Subtle animated gradient background -->
    <div class="absolute inset-0 bg-gradient-to-br from-surface-container-lowest via-secondary-fixed to-primary-fixed opacity-70 z-0"></div>
    <!-- Abstract blurred shapes for "Ethereal Presence" -->
    <div class="absolute top-0 left-0 w-[40rem] h-[40rem] bg-secondary-container rounded-full mix-blend-multiply filter blur-[100px] opacity-30 animate-pulse z-0 -translate-x-1/2 -translate-y-1/2"></div>
    <div class="absolute bottom-0 right-0 w-[30rem] h-[30rem] bg-primary-container rounded-full mix-blend-multiply filter blur-[80px] opacity-20 z-0 translate-x-1/3 translate-y-1/3"></div>

    <main class="relative z-10 w-full max-w-md my-auto">
      <div class="glass-card rounded-xl p-container-padding flex flex-col gap-stack-lg shadow-[0px_10px_30px_rgba(240,98,146,0.1)]">
        <!-- Branding -->
        <div class="text-center flex flex-col items-center gap-stack-sm">
          <div class="w-16 h-16 bg-surface-container-lowest rounded-full flex items-center justify-center border border-outline-variant shadow-sm mb-2">
            <span class="material-symbols-outlined text-[32px] text-primary" style="font-variation-settings: 'FILL' 1;">schedule</span>
          </div>
          <h1 class="font-headline-lg text-headline-lg md:font-headline-lg-mobile md:text-headline-lg-mobile text-primary font-bold tracking-tight">LOPI-Q</h1>
          <p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-widest">
            Logbook, Online Presence, and Internship Quality Management System
          </p>
        </div>

        <!-- Forced 2FA Setup Flow -->
        <div v-if="showSetup2FA" class="flex flex-col gap-stack-md">
          <div class="flex flex-col gap-stack-sm text-center">
            <h2 class="font-title-md text-title-md text-on-surface">
              {{ setupStep === 1 ? 'Aktivasi Keamanan 2FA' : 'Simpan Kode Backup Pemulihan' }}
            </h2>
            <p class="font-body-sm text-body-sm text-on-surface-variant">
              {{ setupStep === 1 ? 'Wajib mengaktifkan 2FA Google Authenticator.' : 'Simpan kode pemulihan ini di tempat aman.' }}
            </p>
          </div>

          <!-- Step 1: Scan QR -->
          <div v-if="setupStep === 1" class="flex flex-col gap-stack-md">
            <div v-if="setupError" class="p-3 rounded-lg bg-error-container border border-outline-variant text-on-error-container text-body-sm flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">error</span>
              <span>{{ setupError }}</span>
            </div>

            <div class="flex items-center gap-4 bg-surface-container-low p-3 rounded-xl border border-outline-variant">
              <div class="relative w-24 h-24 flex items-center justify-center border border-outline-variant p-1 rounded-lg bg-white shrink-0">
                <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="Scan QR Code" class="w-full h-full object-contain" />
                <span v-else class="material-symbols-outlined text-[24px] animate-spin text-primary">sync</span>
              </div>
              <div class="flex flex-col text-left">
                <p class="font-body-sm text-body-sm text-on-surface-variant">
                  Pindai QR ini via <strong>Google Authenticator</strong>.
                </p>
                <div class="bg-surface-container-lowest p-2 rounded border border-outline-variant mt-2">
                  <p class="font-label-md text-label-md text-on-surface-variant uppercase">Kunci Secret</p>
                  <code class="font-body-sm font-bold text-primary break-all block mt-1">{{ totpSecret || 'Memuat...' }}</code>
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-stack-sm">
              <label class="font-label-md text-label-md text-on-surface-variant text-center" for="setup-otp">Kode OTP 6-Digit</label>
              <OtpInput v-model="setupCode" :disabled="setupLoading" />
            </div>

            <button 
              @click="handleConfirmSetup"
              class="w-full flex justify-center items-center gap-2 py-3 px-4 rounded-lg text-on-primary bg-primary-container hover:bg-on-primary-fixed-variant transition-all font-title-md text-title-md border-0 cursor-pointer disabled:opacity-50"
              type="button"
              :disabled="setupLoading || !setupCode"
            >
              <span>{{ setupLoading ? 'Memverifikasi...' : 'Aktifkan & Verifikasi 2FA' }}</span>
            </button>
            <button @click="showSetup2FA = false" type="button" class="text-body-sm text-on-surface-variant hover:text-primary bg-transparent border-0 cursor-pointer">
              Batal
            </button>
          </div>

          <!-- Step 2: Backup codes -->
          <div v-else class="flex flex-col gap-stack-md">
            <div class="grid grid-cols-2 gap-2 bg-surface-container-lowest p-3 rounded-lg border border-outline-variant text-center font-mono text-body-sm font-bold">
              <div v-for="c in backupCodes" :key="c" class="py-1 border-b border-outline-variant/30 text-primary">
                {{ c }}
              </div>
            </div>
            <div class="flex gap-2">
              <button @click="copyBackupCodes" class="flex-1 py-2 border border-outline-variant rounded-lg bg-surface-container-lowest text-on-surface font-label-md text-label-md cursor-pointer" type="button">
                Salin Semua
              </button>
              <button @click="downloadBackupCodes" class="flex-1 py-2 border border-outline-variant rounded-lg bg-surface-container-lowest text-on-surface font-label-md text-label-md cursor-pointer" type="button">
                Unduh .txt
              </button>
            </div>
            <button @click="handleSetupComplete" class="w-full py-3 rounded-lg text-on-primary bg-primary-container font-title-md text-title-md border-0 cursor-pointer" type="button">
              Lanjutkan ke Dashboard
            </button>
          </div>
        </div>

        <!-- Standard Login & OTP Forms -->
        <div v-else class="flex flex-col gap-stack-md">
          <div class="text-center flex flex-col gap-stack-sm">
            <h2 class="font-title-md text-title-md text-on-surface">
              {{ isResetSelfMode ? 'Reset 2FA' : (authStore.otpRequired ? 'Verifikasi Dua Langkah' : 'Masuk ke Sistem') }}
            </h2>
            <p class="font-body-sm text-body-sm text-on-surface-variant">
              {{ isResetSelfMode ? 'Masukkan Kode Backup 8-karakter.' : (authStore.otpRequired ? 'Masukkan kode 6-digit Google Authenticator.' : 'Silakan masukkan kredensial untuk mengakses sistem LOPI-Q.') }}
            </p>
          </div>

          <div v-if="successMessage" class="p-3 rounded-lg bg-surface-container-lowest border border-outline-variant text-primary text-body-sm">
            {{ successMessage }}
          </div>
          <div v-if="errorMessage" class="p-3 rounded-lg bg-error-container border border-outline-variant/30 text-on-error-container text-body-sm">
            {{ errorMessage }}
          </div>

          <!-- Password Login Form -->
          <form v-if="!authStore.otpRequired" @submit.prevent="handleLogin" class="flex flex-col gap-stack-md">
            <div class="flex flex-col gap-stack-sm">
              <label class="font-label-md text-label-md text-on-surface-variant" for="email">Email Address / NIP</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <span class="material-symbols-outlined text-outline text-[20px]">mail</span>
                </div>
                <input 
                  v-model="email"
                  autocomplete="email" 
                  class="block w-full pl-10 pr-3 py-2 border border-outline-variant rounded-lg bg-surface-container-lowest text-on-surface focus:outline-none focus:ring-2 focus:ring-primary-container focus:border-primary-container transition-all shadow-sm font-body-sm text-body-sm placeholder-on-surface-variant/50" 
                  id="email" 
                  name="email" 
                  placeholder="admin@lopi-q.bulukumbakab.go.id" 
                  required 
                  type="text"
                  :disabled="loading"
                />
              </div>
            </div>

            <div class="flex flex-col gap-stack-sm">
              <label class="font-label-md text-label-md text-on-surface-variant" for="password">Password</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <span class="material-symbols-outlined text-outline text-[20px]">lock</span>
                </div>
                <input 
                  v-model="password"
                  autocomplete="current-password" 
                  class="block w-full pl-10 pr-10 py-2 border border-outline-variant rounded-lg bg-surface-container-lowest text-on-surface focus:outline-none focus:ring-2 focus:ring-primary-container focus:border-primary-container transition-all shadow-sm font-body-sm text-body-sm placeholder-on-surface-variant/50" 
                  id="password" 
                  name="password" 
                  placeholder="••••••••" 
                  required 
                  :type="showPassword ? 'text' : 'password'"
                  :disabled="loading"
                />
                <button 
                  type="button"
                  @click="togglePassword" 
                  class="absolute inset-y-0 right-0 pr-3 flex items-center text-outline hover:text-primary border-0 bg-transparent cursor-pointer"
                >
                  <span class="material-symbols-outlined text-[20px]">{{ showPassword ? 'visibility' : 'visibility_off' }}</span>
                </button>
              </div>
            </div>

            <div class="flex items-center justify-between mt-1">
              <div class="flex items-center">
                <input class="h-4 w-4 text-primary-container focus:ring-primary-container border-outline-variant rounded cursor-pointer" id="remember-me" name="remember-me" type="checkbox">
                <label class="ml-2 block font-body-sm text-body-sm text-on-surface-variant cursor-pointer" for="remember-me">
                  Remember me
                </label>
              </div>
              <a @click.prevent="errorMessage = 'Silakan hubungi Super Admin untuk mereset kata sandi.'" href="#" class="font-label-md text-label-md text-primary hover:text-on-primary-fixed-variant transition-colors font-semibold decoration-none">
                Forgot password?
              </a>
            </div>

            <button 
              class="mt-2 w-full flex justify-center items-center gap-2 py-3 px-4 border border-transparent rounded-lg shadow-[0px_4px_14px_rgba(240,98,146,0.2)] text-on-primary bg-primary-container hover:bg-on-primary-fixed-variant hover:shadow-[0px_6px_20px_rgba(240,98,146,0.3)] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-container font-title-md text-title-md text-[16px] transition-all transform active:scale-[0.98] cursor-pointer" 
              type="submit"
              :disabled="loading"
            >
              <span>{{ loading ? 'Mengautentikasi...' : 'Sign in' }}</span>
              <span class="material-symbols-outlined text-[20px]">arrow_forward</span>
            </button>
          </form>

          <!-- OTP Form -->
          <form v-else @submit.prevent="handleVerify2fa" class="flex flex-col gap-stack-md">
            <div class="flex flex-col gap-stack-sm">
              <label class="font-label-md text-label-md text-on-surface-variant" for="otp">
                {{ isResetSelfMode ? 'Kode Backup (8 Karakter)' : (isBackupMode ? 'Kode Backup (8 Karakter)' : 'Kode OTP 6-Digit') }}
              </label>
              <div v-if="!isResetSelfMode && !isBackupMode" class="mt-1">
                <OtpInput v-model="otpCode" :disabled="loading" />
              </div>
              <input 
                v-else
                v-model="otpCode"
                class="block w-full px-3 py-2 border border-outline-variant rounded-lg bg-surface-container-lowest text-on-surface font-body-sm text-body-sm"
                placeholder="a1b2c3d4"
                required
                type="text"
                maxlength="8"
                :disabled="loading"
              />
            </div>

            <div class="flex items-center justify-between text-xs">
              <button type="button" @click="toggleResetSelfMode" class="text-error bg-transparent border-0 cursor-pointer font-label-md text-label-md">
                {{ isResetSelfMode ? 'Batal Reset' : 'Lupa Authenticator?' }}
              </button>
              <button v-if="!isResetSelfMode" type="button" @click="isBackupMode = !isBackupMode; otpCode = '';" class="text-primary bg-transparent border-0 cursor-pointer font-label-md text-label-md">
                {{ isBackupMode ? 'Gunakan App Authenticator' : 'Gunakan Kode Backup' }}
              </button>
            </div>

            <button class="mt-2 w-full py-3 rounded-lg text-on-primary bg-primary-container font-title-md text-title-md border-0 cursor-pointer" type="submit" :disabled="loading">
              {{ loading ? 'Memproses...' : 'Verifikasi & Masuk' }}
            </button>
            <button type="button" @click="authStore.otpRequired = false; isResetSelfMode = false; isBackupMode = false;" class="text-body-sm text-on-surface-variant hover:text-primary bg-transparent border-0 cursor-pointer">
              Kembali ke Login
            </button>
          </form>
        </div>

        <div class="mt-4 text-center border-t border-outline-variant/30 pt-4">
          <p class="font-body-sm text-body-sm text-on-surface-variant">
            Secure access for authorized personnel only.
          </p>
          <router-link to="/" class="mt-2 inline-flex items-center gap-1 font-label-md text-label-md text-primary hover:text-on-primary-fixed-variant decoration-none">
            <span class="material-symbols-outlined text-[16px]">arrow_back</span>
            Kembali ke Beranda
          </router-link>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.material-symbols-outlined {
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}
.material-symbols-outlined.fill {
  font-variation-settings: 'FILL' 1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
