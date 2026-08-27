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
    errorMessage.value = 'NIP/Email dan Password wajib diisi.'
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
        router.push('/intern/scan')
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
        router.push('/intern/scan')
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

<<template>
  <div class="min-h-screen w-full flex bg-gradient-to-tr from-surface-container-low via-surface-bright to-primary-fixed/20 text-on-background antialiased font-body-md text-body-md select-none">
    <!-- Split Screen Layout -->
    <div class="flex w-full min-h-screen">
      <!-- Left Side: Image/Brand (Hidden on Mobile) -->
      <div class="hidden lg:flex w-1/2 flex-col justify-between relative overflow-hidden" style="background-color: #5e002b;">
        <!-- Decorative Background Image -->
        <div 
          class="absolute inset-0 bg-cover bg-center opacity-40" 
          style="background-image: url('https://images.unsplash.com/photo-1582139329536-e7284fece509?q=80&w=1000')"
        ></div>
        <!-- Overlaid Premium Chronos Primary Gradient -->
        <div class="absolute inset-0" style="background: linear-gradient(to bottom right, rgba(171, 44, 93, 0.85), rgba(94, 0, 43, 0.90));"></div>
        <div class="absolute inset-0" style="background: linear-gradient(to top, rgba(37, 24, 30, 0.5), transparent);"></div>
        
        <div class="relative z-10 p-xl flex flex-col h-full justify-between">
          <!-- Branding -->
          <div class="flex flex-col gap-0.5">
            <h1 class="font-headline-md text-headline-md text-white tracking-tight">LOPI-Q</h1>
            <p class="font-label-md text-label-md text-primary-fixed tracking-wider uppercase mt-1">DISKOMINFO KABUPATEN BULUKUMBA</p>
            <p class="font-body-sm text-xs text-primary-fixed-dim opacity-95 mt-1 leading-snug">Layanan Otomasi Presensi, Informasi Aktivitas, dan Quality Monitoring Magang</p>
          </div>
          <!-- Mission Statement -->
          <div class="max-w-md">
            <h2 class="font-display-lg text-display-lg text-white mb-md leading-tight whitespace-pre-line">Presensi Digital Siaga,
Otomasi Logbook Magang.</h2>
            <p class="font-body-md text-body-md text-primary-fixed-dim opacity-90 leading-relaxed">Sistem informasi manajemen presensi dan laporan aktivitas harian terintegrasi dengan pengawasan lokasi geofence presisi radius 2 meter untuk seluruh Siswa dan Mahasiswa Magang Diskominfo Kabupaten Bulukumba.</p>
          </div>
          <!-- Footer Branding -->
          <div class="font-body-sm text-body-sm text-primary-fixed">
            &copy; 2026 Diskominfo dan Persandian Kab. Bulukumba &bull; LOPI-Q
          </div>
        </div>
      </div>

      <!-- Right Side: Login Form -->
      <div class="w-full lg:w-1/2 flex items-center justify-center p-md lg:p-xl bg-gradient-to-tr from-surface-container-low/40 to-surface-bright">
        <div class="w-full max-w-md glass-card rounded-3xl p-xl flex flex-col gap-lg shadow-[0px_10px_30px_rgba(240,98,146,0.15)] border border-outline-variant/30">
          <!-- Mobile Branding (Visible only on small screens) -->
          <div class="flex lg:hidden flex-col items-center gap-1 mb-sm justify-center text-center">
            <h1 class="font-headline-md text-headline-md text-primary">LOPI-Q</h1>
            <p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">NTPD KABUPATEN BULUKUMBA</p>
          </div>

          <!-- Forced 2FA Setup Flow -->
          <div v-if="showSetup2FA" class="flex flex-col gap-3">
            <!-- Setup Header -->
            <div class="flex flex-col gap-0.5">
              <h2 class="text-base font-bold text-on-surface tracking-tight">
                {{ setupStep === 1 ? 'Aktivasi Keamanan 2FA' : 'Simpan Kode Backup Pemulihan' }}
              </h2>
              <p class="text-[11px] text-on-surface-variant leading-relaxed">
                {{ setupStep === 1 ? 'Anda wajib mengaktifkan Autentikasi Dua Faktor (2FA) Google Authenticator.' : 'Simpan kode di bawah jika Anda kehilangan akses Authenticator.' }}
              </p>
            </div>

            <!-- Setup Step 1: Scan QR Code & Enter OTP (Compact Layout) -->
            <div v-if="setupStep === 1" class="flex flex-col gap-3">
              <div v-if="setupError" class="p-2.5 rounded-xl bg-error-container border border-outline-variant/30 text-on-error-container text-[11px] flex items-start gap-2 shadow-sm">
                <span class="material-symbols-outlined text-[16px] shrink-0 mt-0.5">error</span>
                <span>{{ setupError }}</span>
              </div>

              <!-- Side-by-side QR Code and Secret Key -->
              <div class="flex items-center gap-4 bg-surface-container-low border border-outline-variant/30 p-3 rounded-2xl">
                <!-- QR Code Container -->
                <div class="relative w-28 h-28 flex items-center justify-center border border-outline-variant/50 p-1.5 rounded-xl shadow-inner bg-surface-container-lowest shrink-0">
                  <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="Scan QR Code" class="w-full h-full object-contain rounded-lg" />
                  <div v-else class="animate-pulse flex flex-col items-center justify-center gap-1.5 text-on-surface-variant/40">
                    <span class="material-symbols-outlined text-[28px] animate-spin text-primary">sync</span>
                    <span class="text-[9px]">Memuat...</span>
                  </div>
                </div>
                <!-- Manual Instructions -->
                <div class="flex flex-col text-left justify-center min-w-0">
                  <p class="text-[10px] text-on-surface-variant leading-relaxed">
                    Pindai QR ini via aplikasi <strong>Google Authenticator</strong>.
                  </p>
                  <div class="bg-surface-container-lowest p-2 rounded-lg border border-outline-variant/30 mt-2">
                    <p class="text-[9px] text-on-surface-variant font-bold uppercase tracking-wider leading-none">Kunci Manual (Secret Key)</p>
                    <code class="text-[10px] font-bold text-primary break-all select-all font-mono block mt-1">{{ totpSecret || 'Memuat...' }}</code>
                  </div>
                </div>
              </div>

              <!-- Input Code Field -->
              <div class="flex flex-col gap-1">
                <label class="text-[11px] font-bold text-on-surface-variant uppercase tracking-wider text-center" for="setup-otp">Masukkan Kode OTP 6-Digit</label>
                <div class="mt-1">
                  <OtpInput 
                    v-model="setupCode"
                    :disabled="setupLoading"
                  />
                </div>
              </div>

              <!-- Action buttons -->
              <div class="flex flex-col gap-2 mt-1">
                <button 
                  @click="handleConfirmSetup"
                  class="w-full py-2.5 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs rounded-xl hover:shadow-lg transition-all active:scale-[0.98] flex items-center justify-center gap-1.5 shadow-[0px_4px_14px_rgba(240,98,146,0.2)] disabled:opacity-50 duration-200 cursor-pointer border-0" 
                  type="button"
                  :disabled="setupLoading || !setupCode"
                >
                  <span>{{ setupLoading ? 'Memverifikasi...' : 'Aktifkan & Verifikasi 2FA' }}</span>
                  <span v-if="!setupLoading" class="material-symbols-outlined text-[16px]">security</span>
                </button>

                <button 
                  type="button" 
                  @click="showSetup2FA = false"
                  class="w-full py-2 text-center text-xs text-on-surface-variant hover:text-primary transition-colors bg-transparent border-0 cursor-pointer"
                >
                  Batal & Kembali
                </button>
              </div>
            </div>

            <!-- Setup Step 2: Show Backup Codes -->
            <div v-else class="flex flex-col gap-3">
              <div class="p-3 rounded-xl bg-tertiary-fixed border border-outline-variant text-on-tertiary-container text-[11px] flex items-start gap-2 shadow-sm">
                <span class="material-symbols-outlined text-[18px] shrink-0 mt-0.5 text-primary">warning</span>
                <div>
                  <span class="font-bold">PERINGATAN:</span> Harap simpan kode pemulihan ini sekarang sebelum masuk!
                </div>
              </div>

              <!-- Backup codes grid -->
              <div class="grid grid-cols-2 gap-2 bg-surface-container-lowest text-on-surface p-3 rounded-xl font-mono text-xs font-bold tracking-wider text-center border border-outline-variant shadow-inner">
                <div v-for="c in backupCodes" :key="c" class="py-1 px-2 border-b border-outline-variant/30 text-primary hover:text-on-primary-fixed-variant transition-colors">
                  {{ c }}
                </div>
              </div>

              <!-- Backup tools -->
              <div class="flex gap-2">
                <button 
                  @click="copyBackupCodes"
                  class="flex-1 py-2 px-3 border border-outline-variant hover:bg-surface-container-low transition-colors text-xs font-bold rounded-xl flex items-center justify-center gap-1.5 cursor-pointer bg-surface-container-lowest text-on-surface"
                  type="button"
                >
                  <span class="material-symbols-outlined text-[14px]">content_copy</span>
                  <span>Salin Semua</span>
                </button>
                <button 
                  @click="downloadBackupCodes"
                  class="flex-1 py-2 px-3 border border-outline-variant hover:bg-surface-container-low transition-colors text-xs font-bold rounded-xl flex items-center justify-center gap-1.5 cursor-pointer bg-surface-container-lowest text-on-surface"
                  type="button"
                >
                  <span class="material-symbols-outlined text-[14px]">download</span>
                  <span>Unduh .txt</span>
                </button>
              </div>

              <!-- Done button -->
              <button 
                @click="handleSetupComplete"
                class="w-full py-2.5 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs rounded-xl hover:shadow-lg transition-all active:scale-[0.98] flex items-center justify-center gap-1.5 shadow-[0px_4px_14px_rgba(240,98,146,0.2)] duration-200 cursor-pointer border-0" 
                type="button"
              >
                <span>Saya Sudah Simpan, Lanjutkan</span>
                <span class="material-symbols-outlined text-[16px]">done_all</span>
              </button>
            </div>
          </div>

          <!-- Standard Login & OTP Forms -->
          <div v-else class="flex flex-col gap-lg">
            <!-- Form Header -->
            <div class="flex flex-col gap-xs">
              <h2 class="font-headline-sm text-headline-sm text-on-surface tracking-tight">
                {{ isResetSelfMode ? 'Reset Keamanan 2FA' : (authStore.otpRequired ? 'Verifikasi Dua Langkah' : 'Masuk ke Sistem') }}
              </h2>
              <p class="font-body-sm text-body-sm text-on-surface-variant">
                {{ isResetSelfMode ? 'Masukkan salah satu Kode Backup Pemulihan Anda (8 karakter alfanumerik) untuk menonaktifkan 2FA secara permanen.' : (authStore.otpRequired ? 'Masukkan kode autentikasi dari aplikasi Google Authenticator Anda atau kode backup.' : 'Silakan masukkan kredensial Anda untuk mengakses dashboard admin.') }}
              </p>
            </div>

            <!-- Success Alert -->
            <transition name="fade">
              <div v-if="successMessage" class="p-4 rounded-xl bg-surface-container-lowest border border-outline-variant text-primary text-body-sm flex items-start gap-3 shadow-sm">
                <span class="material-symbols-outlined text-[20px] shrink-0 mt-0.5 text-primary">check_circle</span>
                <div>
                  <span class="font-bold">Sukses:</span> {{ successMessage }}
                </div>
              </div>
            </transition>

            <!-- Error Alert -->
            <transition name="fade">
              <div v-if="errorMessage" class="p-4 rounded-xl bg-error-container border border-outline-variant/30 text-on-error-container text-body-sm flex items-start gap-3 shadow-sm">
                <span class="material-symbols-outlined text-[20px] shrink-0 mt-0.5">error</span>
                <div>
                  <span class="font-bold">Gagal masuk:</span> {{ errorMessage }}
                </div>
              </div>
            </transition>

            <!-- Form -->
            <form v-if="!authStore.otpRequired" @submit.prevent="handleLogin" class="flex flex-col gap-md">
              <!-- Email Field -->
              <div class="flex flex-col gap-xs">
                <label class="font-label-md text-label-md text-on-surface-variant" for="email">Alamat Email</label>
                <div class="relative">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-[20px]">mail</span>
                  <input 
                    v-model="email"
                    class="w-full pl-11 pr-4 py-3 bg-surface-container-lowest border border-outline-variant focus:border-primary focus:ring-4 focus:ring-primary-container/20 rounded-xl focus:outline-none transition-all font-body-sm text-body-sm text-on-surface placeholder:text-on-surface-variant/50" 
                    id="email" 
                    name="email" 
                    placeholder="admin@lopi-q.bulukumbakab.go.id" 
                    required 
                    type="email"
                    :disabled="loading"
                  />
                </div>
              </div>

              <!-- Password Field -->
              <div class="flex flex-col gap-xs">
                <label class="font-label-md text-label-md text-on-surface-variant" for="password">Kata Sandi</label>
                <div class="relative">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-[20px]">lock</span>
                  <input 
                    v-model="password"
                    class="w-full pl-11 pr-11 py-3 bg-surface-container-lowest border border-outline-variant focus:border-primary focus:ring-4 focus:ring-primary-container/20 rounded-xl focus:outline-none transition-all font-body-sm text-body-sm text-on-surface placeholder:text-on-surface-variant/50" 
                    id="password" 
                    name="password" 
                    placeholder="••••••••" 
                    required 
                    :type="showPassword ? 'text' : 'password'"
                    :disabled="loading"
                  />
                  <button 
                    class="absolute right-4 top-1/2 -translate-y-1/2 text-outline hover:text-primary transition-colors focus:outline-none flex items-center justify-center border-0 bg-transparent cursor-pointer" 
                    @click="togglePassword" 
                    type="button"
                  >
                    <span class="material-symbols-outlined text-[20px]">
                      {{ showPassword ? 'visibility' : 'visibility_off' }}
                    </span>
                  </button>
                </div>
              </div>

              <!-- Options Row -->
              <div class="flex items-center justify-between mt-xs">
                <div class="flex items-center gap-2">
                  <input 
                    class="w-4 h-4 text-primary-container border-outline-variant rounded focus:ring-primary-container bg-surface cursor-pointer" 
                    id="remember" 
                    name="remember" 
                    type="checkbox"
                  />
                  <label class="font-body-sm text-body-sm text-on-surface-variant select-none cursor-pointer" for="remember">Ingat Saya</label>
                </div>
                <button 
                  type="button"
                  @click="errorMessage = 'Silakan hubungi Super Admin Diskominfo untuk mereset kata sandi Anda.'"
                  class="font-label-md text-label-md text-primary hover:text-on-primary-fixed-variant transition-colors bg-transparent border-0 cursor-pointer font-semibold"
                >
                  Lupa Kata Sandi?
                </button>
              </div>

              <!-- Submit Button -->
              <button 
                class="mt-sm w-full py-3 px-4 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-label-md text-label-md rounded-xl hover:shadow-[0px_6px_20px_rgba(240,98,146,0.3)] transition-all active:scale-[0.98] flex items-center justify-center gap-2 shadow-[0px_4px_14px_rgba(240,98,146,0.2)] disabled:opacity-50 disabled:scale-100 duration-200 cursor-pointer border-0" 
                type="submit"
                :disabled="loading"
              >
                <svg v-if="loading" class="animate-spin h-4 w-4 text-on-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ loading ? 'Mengautentikasi...' : 'Masuk' }}</span>
                <span v-if="!loading" class="material-symbols-outlined text-[18px]">login</span>
              </button>
            </form>

            <!-- 2FA Form -->
            <form v-else @submit.prevent="handleVerify2fa" class="flex flex-col gap-md">
              <div class="flex flex-col gap-xs">
                <label class="font-label-md text-label-md text-on-surface-variant" for="otp">
                  {{ isResetSelfMode ? 'Kode Backup Pemulihan (8 Karakter)' : (isBackupMode ? 'Kode Backup Pemulihan (8 Karakter)' : 'Kode Google Authenticator (6 Digit)') }}
                </label>
                <div v-if="!isResetSelfMode && !isBackupMode" class="mt-1">
                  <OtpInput 
                    v-model="otpCode"
                    :disabled="loading"
                  />
                </div>
                <div v-else class="relative">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-[20px]">
                    key
                  </span>
                  <input 
                    v-model="otpCode"
                    class="w-full pl-11 pr-4 py-3 bg-surface-container-lowest border border-outline-variant focus:border-primary rounded-xl focus:outline-none transition-all font-body-sm text-body-sm text-on-surface placeholder:text-on-surface-variant/50" 
                    id="otp" 
                    name="otp" 
                    placeholder="a1b2c3d4" 
                    required 
                    type="text"
                    maxlength="8"
                    :disabled="loading"
                    autocomplete="one-time-code"
                  />
                </div>
              </div>

              <!-- Options Row for 2FA/Reset -->
              <div class="flex items-center justify-between mt-xs">
                <button 
                  type="button" 
                  @click="toggleResetSelfMode"
                  class="font-label-md text-label-md text-error hover:text-red-700 transition-colors bg-transparent border-0 cursor-pointer"
                >
                  {{ isResetSelfMode ? 'Batal Reset 2FA' : 'Lupa Authenticator?' }}
                </button>
                <button 
                  v-if="!isResetSelfMode"
                  type="button" 
                  @click="isBackupMode = !isBackupMode; otpCode = '';"
                  class="font-label-md text-label-md text-primary hover:text-on-primary-fixed-variant transition-colors bg-transparent border-0 cursor-pointer"
                >
                  {{ isBackupMode ? 'Gunakan Google Authenticator' : 'Gunakan Kode Backup' }}
                </button>
              </div>

              <!-- Verify/Reset Button -->
              <button 
                class="mt-sm w-full py-3 px-4 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-label-md text-label-md rounded-xl hover:shadow-[0px_6px_20px_rgba(240,98,146,0.3)] transition-all active:scale-[0.98] flex items-center justify-center gap-2 shadow-[0px_4px_14px_rgba(240,98,146,0.2)] disabled:opacity-50 disabled:scale-100 duration-200 cursor-pointer border-0" 
                type="submit"
                :disabled="loading"
              >
                <svg v-if="loading" class="animate-spin h-4 w-4 text-on-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ loading ? (isResetSelfMode ? 'Memproses Reset...' : 'Memverifikasi...') : (isResetSelfMode ? 'Reset 2FA Sekarang' : 'Verifikasi & Masuk') }}</span>
                <span v-if="!loading" class="material-symbols-outlined text-[18px]">
                  {{ isResetSelfMode ? 'lock_reset' : 'verified_user' }}
                </span>
              </button>

              <!-- Back to Password Login -->
              <button 
                type="button" 
                @click="authStore.otpRequired = false; isResetSelfMode = false; isBackupMode = false; otpCode = '';"
                class="w-full py-2.5 text-center text-xs text-on-surface-variant hover:text-primary transition-colors bg-transparent border-0 cursor-pointer"
              >
                Kembali ke Login Password
              </button>
            </form>
          </div>

          <!-- Footer Links -->
          <div class="mt-auto pt-lg border-t border-outline-variant/30 flex flex-col gap-3 text-center text-[11px]">
            <p class="text-on-surface-variant/70 leading-relaxed font-medium">
              © 2026 <span class="font-semibold text-on-surface">LOPI-Q Kab. Bulukumba</span>.<br/>
              Developed by <a href="https://diskominfo.bulukumbakab.go.id" target="_blank" rel="noopener noreferrer" class="font-bold text-primary hover:underline">Diskominfo dan Persandian Kab. Bulukumba</a>.
            </p>
            <div class="flex justify-center mt-1">
              <router-link class="px-4 py-2 rounded-full border border-outline-variant/40 hover:border-primary hover:bg-surface-container-lowest text-on-surface-variant hover:text-primary transition-all duration-200 inline-flex items-center justify-center gap-1.5 font-label-md text-xs shadow-sm bg-surface-container-lowest/50 cursor-pointer decoration-none" to="/">
                <span class="material-symbols-outlined text-[16px]">arrow_back</span>
                Kembali ke Halaman Utama
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
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
