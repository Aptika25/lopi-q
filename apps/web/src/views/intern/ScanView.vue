<template>
  <InternLayout>
    <div class="max-w-4xl mx-auto space-y-6 select-none font-sans relative">
      
      <!-- Floating Top Toast Notification Overlay (Instant Alert without Scrolling) -->
      <transition name="toast-slide">
        <div 
          v-if="toast.show" 
          class="fixed top-4 left-4 right-4 sm:left-auto sm:right-6 sm:max-w-md z-50 p-4 rounded-2xl shadow-2xl border backdrop-blur-md flex items-start gap-3 select-none"
          :class="{
            'bg-rose-950/95 border-rose-500 text-white shadow-rose-950/60': toast.type === 'error',
            'bg-amber-950/95 border-amber-500 text-white shadow-amber-950/60': toast.type === 'warning',
            'bg-sky-950/95 border-sky-500 text-white shadow-sky-950/60': toast.type === 'info',
            'bg-emerald-950/95 border-emerald-500 text-white shadow-emerald-950/60': toast.type === 'success'
          }"
        >
          <div 
            class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-inner"
            :class="{
              'bg-rose-500/30 text-rose-300': toast.type === 'error',
              'bg-amber-500/30 text-amber-300': toast.type === 'warning',
              'bg-sky-500/30 text-sky-300': toast.type === 'info',
              'bg-emerald-500/30 text-emerald-300': toast.type === 'success'
            }"
          >
            <span class="material-symbols-outlined text-2xl">
              {{ toast.type === 'error' ? 'report_problem' : (toast.type === 'warning' ? 'lock_clock' : 'info') }}
            </span>
          </div>

          <div class="flex-1 min-w-0 space-y-0.5">
            <h4 class="text-xs font-mono font-black uppercase tracking-wider text-amber-300 flex items-center gap-1.5">
              <span>{{ toast.title }}</span>
            </h4>
            <p class="text-xs font-semibold leading-relaxed opacity-95 text-slate-100">
              {{ toast.message }}
            </p>
          </div>

          <button 
            @click="toast.show = false" 
            class="text-slate-300 hover:text-white text-base font-bold p-1 bg-transparent border-0 cursor-pointer shrink-0"
          >
            ✕
          </button>
        </div>
      </transition>

      <!-- Page Header -->
      <div class="hidden sm:flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Scan QR Code Presensi Siaga 112</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs hidden sm:block">Lakukan pemindaian QR Code resmi pada Posko Siaga 112 untuk melakukan presensi Masuk atau Pulang Tugas.</p>
        </div>
      </div>

      <!-- Shift Schedule & Leave Status Card -->
      <div 
        class="rounded-3xl p-5 border shadow-sm transition-all flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
        :class="{
          'bg-gradient-to-r from-emerald-900 via-teal-900 to-slate-900 border-emerald-700/80 text-white': shiftState.badge === 'ACTIVE',
          'bg-gradient-to-r from-rose-950 via-slate-900 to-rose-900 border-rose-700/80 text-white': shiftState.badge === 'OFF' || shiftState.badge === 'EARLY_LOCKED',
          'bg-gradient-to-r from-amber-950 via-slate-900 to-amber-900 border-amber-700/80 text-white': shiftState.badge === 'LEAVE'
        }"
      >
        <div class="flex items-center gap-3.5 min-w-0">
          <div 
            class="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0 shadow-md border"
            :class="{
              'bg-emerald-500/20 border-emerald-400/40 text-emerald-300': shiftState.badge === 'ACTIVE',
              'bg-rose-500/20 border-rose-400/40 text-rose-300': shiftState.badge === 'OFF' || shiftState.badge === 'EARLY_LOCKED',
              'bg-amber-500/20 border-amber-400/40 text-amber-300': shiftState.badge === 'LEAVE'
            }"
          >
            <span class="material-symbols-outlined text-2xl">
              {{ shiftState.badge === 'ACTIVE' ? 'event_available' : (shiftState.badge === 'LEAVE' ? 'event_busy' : 'lock_clock') }}
            </span>
          </div>

          <div class="space-y-0.5 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-xs font-mono font-extrabold uppercase tracking-wider text-amber-300">Jadwal Shift Hari Ini</span>
              <span 
                class="px-2 py-0.5 text-[10px] font-extrabold rounded-full border uppercase tracking-wider"
                :class="{
                  'bg-emerald-500/20 border-emerald-400/50 text-emerald-200': shiftState.badge === 'ACTIVE',
                  'bg-rose-500/20 border-rose-400/50 text-rose-200': shiftState.badge === 'OFF' || shiftState.badge === 'EARLY_LOCKED',
                  'bg-amber-500/20 border-amber-400/50 text-amber-200': shiftState.badge === 'LEAVE'
                }"
              >
                {{ shiftState.badge === 'ACTIVE' ? '🟢 SHIFT TERBUKA' : (shiftState.badge === 'LEAVE' ? 'ℹ️ CUTI / IZIN' : (shiftState.badge === 'OFF' ? '🔒 LIBUR / OFF' : '⏳ BELUM WAKTU SHIFT')) }}
              </span>
            </div>
            <h3 class="text-sm font-display font-black truncate text-white">
              {{ shiftState.shiftName }}
            </h3>
            <p class="text-[11px] font-semibold text-slate-300">
              {{ shiftState.shiftTimeStr }}
            </p>
          </div>
        </div>

        <div v-if="shiftState.isLocked" class="w-full sm:w-auto text-right">
          <div class="px-3 py-1.5 bg-rose-500/20 border border-rose-400/40 rounded-xl text-[11px] font-bold text-rose-200 leading-snug">
            {{ shiftState.lockReason }}
          </div>
        </div>
      </div>

      <!-- Auto-Detect Scan Status Banner (Tombol Manual Dihapus, Presensi Otomatis oleh Sistem) -->
      <div 
        class="rounded-3xl p-4 border shadow-xs flex items-center justify-between gap-3 transition-all"
        :class="autoDetectedMode === 'MASUK' ? 'bg-emerald-50/90 border-emerald-200 text-emerald-950' : (autoDetectedMode === 'PULANG' ? 'bg-amber-50/90 border-amber-200 text-amber-950' : 'bg-slate-100 border-slate-200 text-slate-800')"
      >
        <div class="flex items-center gap-3.5">
          <div 
            class="w-11 h-11 rounded-2xl flex items-center justify-center font-bold text-white shrink-0 shadow-md"
            :class="autoDetectedMode === 'MASUK' ? 'bg-gradient-to-tr from-emerald-700 to-teal-600' : (autoDetectedMode === 'PULANG' ? 'bg-gradient-to-tr from-amber-700 to-rose-600' : 'bg-slate-600')"
          >
            <span class="material-symbols-outlined text-[24px]">
              {{ autoDetectedMode === 'MASUK' ? 'login' : (autoDetectedMode === 'PULANG' ? 'logout' : 'task_alt') }}
            </span>
          </div>
          <div>
            <div class="text-[10px] font-mono font-black uppercase tracking-wider text-slate-400">Mode Presensi Otomatis Sistem</div>
            <h4 class="text-xs font-black font-display flex items-center gap-2" :class="autoDetectedMode === 'MASUK' ? 'text-emerald-800' : (autoDetectedMode === 'PULANG' ? 'text-amber-800' : 'text-slate-800')">
              <span>{{ autoDetectedMode === 'MASUK' ? '🟢 MODE AUTO: PRESENSI MASUK DUTY' : (autoDetectedMode === 'PULANG' ? '🟠 MODE AUTO: PRESENSI PULANG DUTY' : '✅ PRESENSI HARI INI LENGKAP') }}</span>
            </h4>
            <p class="text-[11px] font-semibold text-slate-600 mt-0.5 leading-snug">
              <span v-if="autoDetectedMode === 'MASUK'">Sistem mendeteksi Anda belum scan Masuk. Hasil scan kamera akan otomatis mencatat Jam Masuk Tugas.</span>
              <span v-else-if="autoDetectedMode === 'PULANG'">Sistem mendeteksi Anda telah scan Masuk. Hasil scan kamera akan otomatis mencatat Jam Pulang Tugas.</span>
              <span v-else>Anda telah menyelesaikan seluruh siklus presensi Masuk &amp; Pulang untuk siaga shift ini.</span>
            </p>
          </div>
        </div>

        <div class="hidden sm:flex flex-col items-end shrink-0">
          <span 
            class="px-3 py-1 text-[10px] font-extrabold rounded-full border uppercase tracking-wider shadow-2xs"
            :class="autoDetectedMode === 'MASUK' ? 'bg-emerald-100 text-emerald-800 border-emerald-300' : (autoDetectedMode === 'PULANG' ? 'bg-amber-100 text-amber-800 border-amber-300' : 'bg-slate-200 text-slate-700 border-slate-300')"
          >
            Otomatis Deteksi Shift
          </span>
          <span class="text-[9px] font-mono font-semibold text-slate-400 mt-1">Toleransi Keterlambatan: 30 Menit</span>
        </div>
      </div>

      <!-- Main Scanner & GPS Geofence Grid -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
        
        <!-- Left Column: Camera Scanner Card (7 cols) -->
        <div class="md:col-span-7 bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-base font-display font-bold text-slate-800 flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">photo_camera</span>
              <span>Pemindai Kamera Live</span>
            </h2>
            <span v-if="cameraActive" class="px-2.5 py-1 bg-emerald-100 text-emerald-800 text-[11px] font-extrabold rounded-full flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping"></span>
              <span>Kamera Aktif</span>
            </span>
          </div>

          <!-- Camera Viewport Box -->
          <div class="relative w-full aspect-square sm:aspect-video bg-slate-900 rounded-2xl overflow-hidden flex items-center justify-center border border-slate-800 shadow-inner">
            <!-- Video Stream -->
            <video 
              ref="videoRef" 
              class="w-full h-full object-cover" 
              autoplay 
              playsinline 
              v-show="cameraActive"
            ></video>

            <!-- Camera Scanner Reticle Overlay -->
            <div v-if="cameraActive" class="absolute inset-0 pointer-events-none flex items-center justify-center">
              <div class="w-48 h-48 sm:w-56 sm:h-56 border-2 border-rose-500 rounded-2xl relative shadow-[0_0_30px_rgba(225,29,72,0.4)] animate-pulse">
                <div class="absolute -top-1 -left-1 w-6 h-6 border-t-4 border-l-4 border-amber-400"></div>
                <div class="absolute -top-1 -right-1 w-6 h-6 border-t-4 border-r-4 border-amber-400"></div>
                <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-4 border-l-4 border-amber-400"></div>
                <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-4 border-r-4 border-amber-400"></div>
              </div>
            </div>

            <!-- Camera Inactive Placeholder -->
            <div v-if="!cameraActive" class="text-center p-6 space-y-3">
              <div class="w-16 h-16 bg-rose-950/60 rounded-full flex items-center justify-center mx-auto text-rose-400">
                <span class="material-symbols-outlined text-3xl">videocam_off</span>
              </div>
              <div>
                <p class="text-xs font-bold text-slate-300">Kamera Belum Aktif</p>
                <p class="text-[11px] text-slate-500 max-w-xs mx-auto mt-1">
                  Klik tombol "Buka Kamera Scan" di bawah untuk mengaktifkan pemindai kamera live.
                </p>
              </div>
            </div>
          </div>

          <!-- Camera Control Actions -->
          <div class="flex gap-3">
            <button 
              v-if="!cameraActive"
              @click="startCamera"
              :disabled="shiftState.isLocked && !authStore.todayStatus?.is_masuk"
              class="w-full py-3 px-4 bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-extrabold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer border-0 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              <span class="material-symbols-outlined text-[18px]">videocam</span>
              <span>Buka Kamera Scan</span>
            </button>
            <button 
              v-else
              @click="stopCamera"
              class="w-full py-3 px-4 bg-slate-100 hover:bg-slate-200 text-slate-700 font-extrabold text-xs rounded-xl border border-slate-200 transition-all flex items-center justify-center gap-2 cursor-pointer border-0"
            >
              <span class="material-symbols-outlined text-[18px]">videocam_off</span>
              <span>Tutup Kamera Pemindai</span>
            </button>
          </div>

          <!-- Scanned QR Token Result & Action Alerts -->
          <div v-if="actionError" class="p-3 bg-rose-50 border border-rose-200 rounded-xl text-rose-800 text-xs flex items-center gap-2">
            <span class="material-symbols-outlined text-[18px] shrink-0 text-rose-600">error</span>
            <span>{{ actionError }}</span>
          </div>

          <div v-if="scanMode === 'PULANG' && !clockOutWindowInfo.isAllowed" class="p-3.5 bg-amber-50 border border-amber-300 rounded-2xl text-amber-900 text-xs flex items-center gap-2 font-bold shadow-xs">
            <span class="material-symbols-outlined text-amber-600 text-[20px] shrink-0">lock_clock</span>
            <span>{{ clockOutWindowInfo.lockReason }}</span>
          </div>

          <div v-if="scannedToken" class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-900 text-xs flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined text-emerald-600 text-[18px]">check_circle</span>
              <span>Token QR Terdeteksi: <strong class="font-mono text-emerald-800">{{ scannedToken }}</strong></span>
            </div>
          </div>
        </div>

        <!-- Right Column: Live GPS & Geofence Map (5 cols) -->
        <div class="md:col-span-5 space-y-6">
          <!-- Geofence Monitor & Interactive Map Card -->
          <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-4">
            <div class="flex items-center justify-between">
              <h2 class="text-base font-display font-bold text-slate-800 flex items-center gap-2">
                <span class="material-symbols-outlined text-rose-700">map</span>
                <span>Peta Geofence &amp; Radar GPS</span>
              </h2>
              <span class="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider">Live Maps</span>
            </div>

            <!-- Distance Metric Box -->
            <div 
              class="p-4 rounded-2xl border text-center transition-all"
              :class="isWithinRadius ? 'bg-emerald-50/80 border-emerald-200 text-emerald-950' : 'bg-rose-50/80 border-rose-200 text-rose-950'"
            >
              <div class="text-[10px] font-bold uppercase tracking-wider text-slate-500">Jarak Anda ke Posko Siaga 112</div>
              <div class="text-3xl font-black font-display my-1" :class="isWithinRadius ? 'text-emerald-700' : 'text-rose-700'">
                {{ distanceMeters !== null ? `${distanceMeters.toFixed(2)} m` : 'Memuat GPS...' }}
              </div>
              <div class="inline-flex items-center gap-1 text-[11px] font-bold px-2.5 py-0.5 rounded-full" :class="isWithinRadius ? 'bg-emerald-200 text-emerald-900' : 'bg-rose-200 text-rose-900'">
                <span class="material-symbols-outlined text-[14px]">{{ isWithinRadius ? 'check' : 'warning' }}</span>
                <span>{{ isWithinRadius ? `DI DALAM RADIUS (<= ${maxRadius}m)` : `DILUAR RADIUS (> ${maxRadius}m)` }}</span>
              </div>
            </div>

            <!-- Interactive Leaflet Map for Mobile & Desktop -->
            <div class="relative">
              <div id="scanLeafletMap" class="w-full h-60 rounded-2xl border border-slate-200 shadow-inner overflow-hidden z-10"></div>
              
              <!-- Map Legend Overlay -->
              <div class="absolute bottom-2 left-2 right-2 bg-white/90 backdrop-blur-xs p-2 rounded-xl border border-slate-200/80 text-[10px] flex items-center justify-around font-semibold text-slate-700 z-20 shadow-xs">
                <div class="flex items-center gap-1">
                  <span class="w-2.5 h-2.5 rounded-full bg-rose-600 inline-block"></span>
                  <span>Posko 112</span>
                </div>
                <div class="flex items-center gap-1">
                  <span class="w-2.5 h-2.5 rounded-full bg-sky-600 inline-block"></span>
                  <span>Posisi Anda</span>
                </div>
                <div class="flex items-center gap-1">
                  <span class="w-2.5 h-2.5 rounded-full" :class="isWithinRadius ? 'bg-emerald-500' : 'bg-rose-500'"></span>
                  <span>Radius {{ maxRadius }}m</span>
                </div>
              </div>
            </div>

            <!-- Anti-Fake GPS Warning Alert -->
            <div v-if="isMockDetected" class="p-4 bg-rose-100 border-2 border-rose-300 rounded-2xl text-rose-900 text-xs space-y-1.5 shadow-sm">
              <div class="flex items-center gap-2 font-extrabold text-rose-800">
                <span class="material-symbols-outlined text-rose-600 text-[20px]">security_update_warning</span>
                <span>TERDETEKSI MANIPULASI GPS / FAKE GPS</span>
              </div>
              <p class="text-[11px] text-rose-700 font-semibold leading-relaxed">
                {{ mockReason || 'Sistem mendeteksi penggunaan aplikasi Lokasi Palsu / Mock Location. Presensi ditolak!' }}
              </p>
              <div class="text-[10px] font-bold text-rose-950 bg-white/80 p-2.5 rounded-xl border border-rose-200 mt-1">
                🔒 Peringatan Keamanan: Harap matikan aplikasi Fake GPS / Opsi Pengembang Mock Location di HP Anda dan gunakan sinyal GPS hardware asli perangkat.
              </div>
            </div>

            <!-- GPS Coordinates Details -->
            <div class="bg-slate-50 p-3 rounded-xl border border-slate-200 text-xs space-y-1">
              <div class="flex justify-between text-slate-600">
                <span>Latitude Hardware:</span>
                <span class="font-mono font-bold text-slate-900">{{ currentLat?.toFixed(6) || '...' }}</span>
              </div>
              <div class="flex justify-between text-slate-600">
                <span>Longitude Hardware:</span>
                <span class="font-mono font-bold text-slate-900">{{ currentLng?.toFixed(6) || '...' }}</span>
              </div>
              <div class="flex justify-between text-slate-600 pt-1 border-t border-slate-200">
                <span>Akurasi Sensor Sinyal:</span>
                <span class="font-mono font-bold text-emerald-700">± {{ gpsAccuracy }} meter</span>
              </div>
            </div>

            <!-- Refresh GPS Button -->
            <button 
              @click="getGeolocation"
              :disabled="gpsLoading"
              class="w-full py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold text-xs rounded-xl transition-all flex items-center justify-center gap-1.5 cursor-pointer border-0 disabled:opacity-50"
            >
              <span class="material-symbols-outlined text-[16px]" :class="{ 'animate-spin': gpsLoading }">
                {{ gpsLoading ? 'sync' : 'refresh' }}
              </span>
              <span>{{ gpsLoading ? 'Memindai Satelit GPS...' : 'Perbarui Lokasi Hardware GPS' }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Full-Screen Success Overlay Modal -->
      <transition name="fade">
        <div v-if="isSuccessModalOpen" class="fixed inset-0 bg-slate-900/80 backdrop-blur-md z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-3xl max-w-sm w-full p-8 text-center space-y-5 shadow-2xl border border-slate-200">
            <div class="w-20 h-20 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto border-4 border-emerald-200 animate-pulse">
              <span class="material-symbols-outlined text-5xl">check_circle</span>
            </div>

            <div class="space-y-1">
              <span class="px-3 py-1 bg-emerald-100 text-emerald-800 text-[10px] font-extrabold rounded-full uppercase tracking-wider">
                PRESENSI {{ scanMode }} BERHASIL!
              </span>
              <h3 class="text-lg font-display font-black text-slate-900 mt-2">
                {{ authStore.user?.name || 'Petugas Peserta Magang' }}
              </h3>
              <p class="text-xs text-slate-500 font-medium leading-relaxed">
                {{ successMessage }}
              </p>
            </div>

            <div class="bg-slate-50 p-3.5 rounded-2xl border border-slate-200 text-xs font-mono space-y-1 text-slate-700 text-left">
              <div class="flex justify-between">
                <span>Jarak ke Posko:</span>
                <strong class="text-slate-900">{{ distanceMeters?.toFixed(2) }} meter</strong>
              </div>
              <div class="flex justify-between pt-1 border-t border-slate-200">
                <span>Status Geofence:</span>
                <strong class="text-emerald-700">✓ Valid (&lt;= {{ maxRadius }}m)</strong>
              </div>
            </div>

            <div class="text-[11px] font-bold text-rose-700 animate-pulse flex items-center justify-center gap-1.5">
              <span class="material-symbols-outlined text-[16px] animate-spin">sync</span>
              <span>Menutup kamera &amp; membuka Laporan ({{ redirectCountdown }}s)...</span>
            </div>
          </div>
        </div>
      </transition>

    </div>
  </InternLayout>
</template>

<script setup lang="ts">
import InternLayout from '@/layouts/InternLayout.vue'
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { calculateDistance } from '@/utils/haversine'

const authStore = useAuthStore()
const router = useRouter()

const autoDetectedMode = computed<'MASUK' | 'PULANG' | 'SELESAI'>(() => {
  if (!authStore.todayStatus?.is_masuk) {
    return 'MASUK'
  }
  if (!authStore.todayStatus?.is_pulang) {
    return 'PULANG'
  }
  return 'SELESAI'
})

const scanMode = computed<'MASUK' | 'PULANG' | 'SELESAI'>(() => {
  return autoDetectedMode.value
})

const isSuccessModalOpen = ref(false)
const successMessage = ref('')
const redirectCountdown = ref(2)

const toast = ref<{
  show: boolean
  type: 'error' | 'warning' | 'info' | 'success'
  title: string
  message: string
}>({
  show: false,
  type: 'error',
  title: '',
  message: ''
})

let toastTimer: any = null

const showToast = (title: string, message: string, type: 'error' | 'warning' | 'info' | 'success' = 'error') => {
  if (toastTimer) clearTimeout(toastTimer)
  toast.value = { show: true, type, title, message }
  toastTimer = setTimeout(() => {
    toast.value.show = false
  }, 6000)
}

const shiftState = ref<{
  isLocked: boolean
  badge: 'ACTIVE' | 'EARLY_LOCKED' | 'OFF' | 'LEAVE' | 'FINISHED'
  shiftName: string
  shiftTimeStr: string
  lockReason: string
}>({
  isLocked: false,
  badge: 'ACTIVE',
  shiftName: 'Shift Siaga 112',
  shiftTimeStr: 'Jendela Presensi Terbuka',
  lockReason: ''
})

const currentShiftStartHour = ref<number>(8)
const currentShiftEndHour = ref<number>(16)

const clockOutWindowInfo = computed(() => {
  // If user has NOT scanned MASUK yet, clock-out is not applicable
  if (!authStore.todayStatus?.is_masuk) {
    return {
      isAllowed: false,
      startStr: '15:30',
      endStr: '23:59',
      lockReason: 'Anda belum melakukan Presensi Masuk.'
    }
  }

  // If user ALREADY completed both MASUK & PULANG
  if (authStore.todayStatus?.is_masuk && authStore.todayStatus?.is_pulang) {
    return {
      isAllowed: false,
      startStr: '00:00',
      endStr: '23:59',
      lockReason: 'Tugas Siaga Hari ini telah selesai (Absen Masuk & Pulang Terverifikasi).'
    }
  }

  // User HAS scanned MASUK and NOT yet scanned PULANG
  const now = new Date()
  const nowHour = now.getHours()
  const nowMins = nowHour * 60 + now.getMinutes()

  const startH = currentShiftStartHour.value
  const endH = currentShiftEndHour.value

  // Check based on Shift Type
  // 1. Shift Pagi / Shift 1 (start ~8:00, end ~16:00 or 20:00)
  if (startH >= 6 && startH <= 10) {
    const minClockOutMins = (endH === 16) ? (15 * 60 + 30) : (19 * 60 + 30)
    const isAllowed = nowMins >= minClockOutMins || nowHour >= endH || nowHour < 4
    const openTimeStr = (endH === 16) ? '15:30' : '19:30'

    return {
      isAllowed,
      startStr: openTimeStr,
      endStr: '23:59',
      lockReason: isAllowed 
        ? '' 
        : `Anda telah melakukan Presensi Masuk. Presensi Pulang untuk Shift Pagi baru dibuka pukul ${openTimeStr} WITA (30 menit sebelum piket selesai).`
    }
  }

  // 2. Shift Sore / Shift 2 (3-shift mode: 16:00 - 24:00)
  if (startH >= 14 && startH <= 17) {
    const minClockOutMins = 23 * 60 + 30
    const isAllowed = nowMins >= minClockOutMins || nowHour < 4
    return {
      isAllowed,
      startStr: '23:30',
      endStr: '03:00',
      lockReason: isAllowed 
        ? '' 
        : `Anda telah melakukan Presensi Masuk. Presensi Pulang untuk Shift Sore baru dibuka pukul 23:30 WITA.`
    }
  }

  // 3. Shift Malam / Night Shift (20:00 - 08:00 or 00:00 - 08:00)
  if (startH >= 18 || startH === 0) {
    const isAllowed = (nowMins >= 7 * 60 + 30 && nowHour < 15) || nowHour >= 20 || nowHour < 4
    return {
      isAllowed,
      startStr: '07:30',
      endStr: '15:00',
      lockReason: isAllowed 
        ? '' 
        : `Anda telah melakukan Presensi Masuk. Presensi Pulang Shift Malam baru dibuka pukul 07:30 WITA (30 menit sebelum piket selesai).`
    }
  }

  return {
    isAllowed: true,
    startStr: '00:00',
    endStr: '23:59',
    lockReason: ''
  }
})

const cameraActive = ref(false)
const videoRef = ref<HTMLVideoElement | null>(null)
const mediaStream = ref<MediaStream | null>(null)
const scannedToken = ref('')

const currentLat = ref<number | null>(-5.5645)
const currentLng = ref<number | null>(120.1945)
const gpsAccuracy = ref<number>(5.0)

const poskoLat = ref(-5.5645)
const poskoLng = ref(120.1945)
const maxRadius = ref(2.0)
const poskoToken = ref('')

const gpsLoading = ref(false)
const isMockDetected = ref(false)
const mockReason = ref('')
const lastGpsTime = ref<number>(0)

const loading = ref(false)
const actionError = ref('')

let scanLoopId: number | null = null
let mapInstance: any = null
let poskoMarker: any = null
let userMarker: any = null
let circleInstance: any = null
let polylineInstance: any = null

const distanceMeters = computed(() => {
  if (currentLat.value === null || currentLng.value === null) return null
  return calculateDistance(currentLat.value, currentLng.value, poskoLat.value, poskoLng.value)
})

const isWithinRadius = computed(() => {
  if (distanceMeters.value === null) return false
  return distanceMeters.value <= maxRadius.value
})

// ===== SHIFT SCHEDULE & LEAVE VALIDATION =====
const evaluateShiftSchedule = async () => {
  try {
    // Priority 0: If user has an active, unclosed MASUK record, determine shift type from clock-in timestamp
    if (authStore.todayStatus?.is_masuk && !authStore.todayStatus?.is_pulang) {
      const masukTimeStr = authStore.todayStatus?.masuk?.timestamp || authStore.todayStatus?.clock_in_time || ''
      const isYesterdayShift = masukTimeStr.includes('Shift Kemarin')

      let clockInHour = -1
      const timeMatch = masukTimeStr.match(/(\d{1,2}):(\d{2})/)
      if (timeMatch) {
        clockInHour = parseInt(timeMatch[1], 10)
      }

      // If clock-in was in evening (>= 18:00 or <= 04:00) or yesterday's shift -> Night Shift (20:00 - 08:00)
      if (isYesterdayShift || clockInHour >= 18 || (clockInHour >= 0 && clockInHour <= 4)) {
        currentShiftStartHour.value = 20
        currentShiftEndHour.value = 8
        shiftState.value = {
          isLocked: false,
          badge: 'ACTIVE',
          shiftName: 'Shift Malam (Sedang Berjalan)',
          shiftTimeStr: 'Presensi Selesai Siaga (Pulang) Terbuka',
          lockReason: ''
        }
        return
      }

      // If clock-in was in afternoon (>= 12:00 && <= 17:00) -> Shift Sore (16:00 - 24:00)
      if (clockInHour >= 12 && clockInHour <= 17) {
        currentShiftStartHour.value = 16
        currentShiftEndHour.value = 24
        shiftState.value = {
          isLocked: false,
          badge: 'ACTIVE',
          shiftName: 'Shift Sore (Sedang Berjalan)',
          shiftTimeStr: 'Presensi Selesai Siaga (Pulang) Terbuka',
          lockReason: ''
        }
        return
      }

      // Default for morning clock-in -> Shift Pagi (08:00 - 20:00)
      currentShiftStartHour.value = 8
      currentShiftEndHour.value = 20
      shiftState.value = {
        isLocked: false,
        badge: 'ACTIVE',
        shiftName: 'Shift Pagi (Sedang Berjalan)',
        shiftTimeStr: 'Presensi Selesai Siaga (Pulang) Terbuka',
        lockReason: ''
      }
      return
    }
    const now = new Date()
    const todayStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    const userNip = authStore.user?.nip || ''
    const userName = authStore.user?.name || ''

    // 1. Fetch Leave Requests
    const leaveRes = await axios.get('/api/presensi/leave-requests')
    const leaves = leaveRes.data?.requests || []

    const activeLeave = leaves.find((l: any) => {
      if (l.status !== 'Approved') return false
      const matchesUser = (userNip && l.nip === userNip) || (userName && l.userName === userName)
      if (!matchesUser) return false
      return todayStr >= l.startDate && todayStr <= l.endDate
    })

    if (activeLeave) {
      shiftState.value = {
        isLocked: true,
        badge: 'LEAVE',
        shiftName: 'Izin / Cuti Resmi (Disetujui Admin)',
        shiftTimeStr: `${activeLeave.startDate} s/d ${activeLeave.endDate}`,
        lockReason: `Anda dalam masa Cuti / Izin: "${activeLeave.reason || 'Dinas Luar'}". Presensi ditolak.`
      }
      return
    }

    // 2. Fetch Schedules
    const schedRes = await axios.get('/api/admin/schedules')
    const schedData = schedRes.data?.schedules
    if (schedData && schedData.daysInMonth) {
      const dayEntry = schedData.daysInMonth.find((d: any) => d.dateStr === todayStr || d.date === todayStr) || schedData.daysInMonth[now.getDate() - 1]
      const teams = schedData.teams || []
      const shiftMode = schedData.shiftMode || 2

      const slotTimes = [
        { shiftKey: 'shift1', replKey: 'replacementsShift1', startHour: 8, endHour: shiftMode === 3 ? 16 : 20, label: shiftMode === 3 ? 'Shift Pagi' : 'Shift Pagi' },
        { shiftKey: 'shift2', replKey: 'replacementsShift2', startHour: shiftMode === 3 ? 16 : 20, endHour: shiftMode === 3 ? 24 : 8, label: shiftMode === 3 ? 'Shift Sore' : 'Shift Malam' },
        { shiftKey: 'shift3', replKey: 'replacementsShift3', startHour: 0, endHour: 8, label: 'Shift Malam' }
      ]

      const userNipClean = userNip.replace(/\s/g, '')

      // ─── 1. Cek apakah user terdaftar sebagai PENGGANTI hari ini ───
      let foundAsReplacer = false
      for (const slot of slotTimes) {
        const replacements: any[] = dayEntry?.[slot.replKey] || []
        const myEntry = replacements.find((r: any) =>
          r.replacerNip && r.replacerNip.replace(/\s/g, '') === userNipClean
        )
        if (myEntry) {
          foundAsReplacer = true
          const startHour = slot.startHour
          const endHour = slot.endHour
          currentShiftStartHour.value = startHour
          currentShiftEndHour.value = endHour

          const currentHour = now.getHours()
          const currentMin = now.getMinutes()
          const currentTotalMins = currentHour * 60 + currentMin
          const shiftStartMins = startHour * 60
          const earlyWindowMins = shiftStartMins - 30

          if (!authStore.todayStatus?.is_masuk && currentTotalMins < earlyWindowMins && currentTotalMins > 0) {
            const wH = String(Math.floor(earlyWindowMins / 60)).padStart(2, '0')
            const wM = String(earlyWindowMins % 60).padStart(2, '0')
            shiftState.value = {
              isLocked: true,
              badge: 'EARLY_LOCKED',
              shiftName: `${slot.label} — Pengganti ${myEntry.replacedName.split(',')[0]}`,
              shiftTimeStr: `Jadwal mulai ${String(startHour).padStart(2, '0')}:00 WITA`,
              lockReason: `Belum waktu shift. Presensi Masuk dibuka pukul ${wH}:${wM} WITA.`
            }
          } else {
            shiftState.value = {
              isLocked: false,
              badge: 'ACTIVE',
              shiftName: `${slot.label} — Pengganti ${myEntry.replacedName.split(',')[0]} 🟡`,
              shiftTimeStr: `${String(startHour).padStart(2,'0')}:00 – ${String(endHour === 24 ? 0 : endHour).padStart(2,'0')}:00 WITA`,
              lockReason: ''
            }
          }
          break
        }
      }
      if (foundAsReplacer) return

      // ─── 2. Cek tim utama user ───
      let userTeamCode = ''
      for (const t of teams) {
        if (t.members && Array.isArray(t.members)) {
          if (t.members.some((m: any) =>
            (userNipClean && m.nip && m.nip.replace(/\s/g, '') === userNipClean) ||
            (userName && m.name && m.name.toLowerCase().includes(userName.toLowerCase().split(' ')[0]))
          )) {
            userTeamCode = t.code
            break
          }
        }
      }

      if (!userTeamCode && teams.length > 0) userTeamCode = teams[0].code

      if (dayEntry) {
        if (dayEntry.offTeams && dayEntry.offTeams.includes(userTeamCode)) {
          // If user already scanned MASUK, do NOT lock them even if team shows OFF today
          if (!authStore.todayStatus?.is_masuk) {
            shiftState.value = {
              isLocked: true,
              badge: 'OFF',
              shiftName: `HARI LIBUR (OFF) — Tim ${userTeamCode}`,
              shiftTimeStr: 'Tidak Ada Jadwal Shift Hari Ini',
              lockReason: `Tim ${userTeamCode} dijadwalkan LIBUR hari ini. Presensi terkunci.`
            }
            return
          }
        }

        for (const slot of slotTimes) {
          if (dayEntry[slot.shiftKey] === userTeamCode) {
            const startHour = slot.startHour
            const endHour = slot.endHour
            currentShiftStartHour.value = startHour
            currentShiftEndHour.value = endHour

            const currentHour = now.getHours()
            const currentMin = now.getMinutes()
            const currentTotalMins = currentHour * 60 + currentMin
            const shiftStartMins = startHour * 60
            const earlyWindowMins = shiftStartMins - 30

            if (!authStore.todayStatus?.is_masuk && currentTotalMins < earlyWindowMins && currentTotalMins > 0) {
              const wH = String(Math.floor(earlyWindowMins / 60)).padStart(2, '0')
              const wM = String(earlyWindowMins % 60).padStart(2, '0')
              shiftState.value = {
                isLocked: true,
                badge: 'EARLY_LOCKED',
                shiftName: `${slot.label} — Tim ${userTeamCode}`,
                shiftTimeStr: `Jadwal Shift Dimulai ${String(startHour).padStart(2, '0')}:00 WITA`,
                lockReason: `Belum waktu shift. Presensi Masuk dibuka pukul ${wH}:${wM} WITA.`
              }
            } else {
              shiftState.value = {
                isLocked: false,
                badge: 'ACTIVE',
                shiftName: `${slot.label} — Tim ${userTeamCode}`,
                shiftTimeStr: `Presensi Pulang dibuka 30 min sblm shift berakhir`,
                lockReason: ''
              }
            }
            return
          }
        }
      }
    }

    // Priority 0 fallback if already scanned MASUK
    if (authStore.todayStatus?.is_masuk && !authStore.todayStatus?.is_pulang) {
      shiftState.value = {
        isLocked: false,
        badge: 'ACTIVE',
        shiftName: 'Shift Siaga 112 (Sedang Berjalan)',
        shiftTimeStr: 'Presensi Selesai Siaga (Pulang) Terbuka',
        lockReason: ''
      }
    }
  } catch (e) {
    console.warn('Gagal memverifikasi jadwal shift:', e)
  }
}


// ===== LEAFLET MAP INTEGRATION FOR PESERTA MAGANG =====
const initLeafletMap = () => {
  if (typeof (window as any).L === 'undefined') {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
    document.head.appendChild(link)

    const script = document.createElement('script')
    script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
    script.onload = () => setupLeaflet()
    document.body.appendChild(script)
  } else {
    setupLeaflet()
  }
}

const setupLeaflet = () => {
  nextTick(() => {
    const L = (window as any).L
    if (!L) return

    const container = document.getElementById('scanLeafletMap')
    if (!container) return

    if (mapInstance) {
      mapInstance.remove()
      mapInstance = null
    }

    mapInstance = L.map('scanLeafletMap').setView([poskoLat.value, poskoLng.value], 17)

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; OpenStreetMap &bull; LOPI-Q'
    }).addTo(mapInstance)

    const poskoIcon = L.divIcon({
      className: 'custom-posko-pin',
      html: `<div style="background-color:#e11d48;color:white;width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 10px rgba(225,29,72,0.4);border:2px solid white;"><span class="material-symbols-outlined" style="font-size:18px;">local_hospital</span></div>`,
      iconSize: [32, 32],
      iconAnchor: [16, 16]
    })

    const userIcon = L.divIcon({
      className: 'custom-user-pin',
      html: `<div style="background-color:#0284c7;color:white;width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 10px rgba(2,132,199,0.4);border:2px solid white;"><span class="material-symbols-outlined" style="font-size:18px;">person_pin_circle</span></div>`,
      iconSize: [32, 32],
      iconAnchor: [16, 16]
    })

    poskoMarker = L.marker([poskoLat.value, poskoLng.value], { icon: poskoIcon }).addTo(mapInstance)
      .bindPopup(`<b>Posko Siaga 112 Bulukumba</b><br>Radius Geofence: ${maxRadius.value}m`)

    const userLat = currentLat.value || poskoLat.value
    const userLng = currentLng.value || poskoLng.value
    userMarker = L.marker([userLat, userLng], { icon: userIcon }).addTo(mapInstance)
      .bindPopup(`<b>Posisi Peserta Magang</b>`)

    updateMapVisuals()
  })
}

const updateMapVisuals = () => {
  const L = (window as any).L
  if (!L || !mapInstance) return

  const userLat = currentLat.value || poskoLat.value
  const userLng = currentLng.value || poskoLng.value

  if (userMarker) {
    userMarker.setLatLng([userLat, userLng])
  }
  if (poskoMarker) {
    poskoMarker.setLatLng([poskoLat.value, poskoLng.value])
  }

  // Draw Geofence Circle
  if (circleInstance) {
    mapInstance.removeLayer(circleInstance)
  }
  const circleColor = isWithinRadius.value && !isMockDetected.value ? '#10b981' : '#f43f5e'
  circleInstance = L.circle([poskoLat.value, poskoLng.value], {
    color: circleColor,
    fillColor: circleColor,
    fillOpacity: 0.25,
    radius: maxRadius.value
  }).addTo(mapInstance)

  // Connector Line
  if (polylineInstance) {
    mapInstance.removeLayer(polylineInstance)
  }
  polylineInstance = L.polyline([
    [poskoLat.value, poskoLng.value],
    [userLat, userLng]
  ], {
    color: circleColor,
    weight: 3,
    dashArray: '6, 6'
  }).addTo(mapInstance)

  try {
    const bounds = L.latLngBounds([
      [poskoLat.value, poskoLng.value],
      [userLat, userLng]
    ])
    mapInstance.fitBounds(bounds, { padding: [35, 35], maxZoom: 18 })
  } catch (e) {}
}

const fetchPoskoInfo = async () => {
  try {
    const res = await axios.get('/api/presensi/posko-qr')
    if (res.data && res.data.coordinates) {
      if (res.data.coordinates.latitude) poskoLat.value = res.data.coordinates.latitude
      if (res.data.coordinates.longitude) poskoLng.value = res.data.coordinates.longitude
      if (res.data.coordinates.max_radius_meters) maxRadius.value = res.data.coordinates.max_radius_meters
    }
    if (res.data && res.data.qr_token) {
      poskoToken.value = res.data.qr_token
    }
    updateMapVisuals()
  } catch (e) {
    console.warn('Menggunakan data posko default.')
  }
}

const startCamera = async () => {
  try {
    actionError.value = ''
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' }
    })
    mediaStream.value = stream
    if (videoRef.value) {
      videoRef.value.srcObject = stream
    }
    cameraActive.value = true
    startScanLoop()
  } catch (err: any) {
    console.warn('Gagal membuka kamera:', err.message)
    actionError.value = 'Akses kamera tidak aktif atau tidak diizinkan di browser.'
  }
}

const stopCamera = () => {
  if (scanLoopId !== null) {
    cancelAnimationFrame(scanLoopId)
    scanLoopId = null
  }
  if (mediaStream.value) {
    mediaStream.value.getTracks().forEach(track => track.stop())
    mediaStream.value = null
  }
  cameraActive.value = false
}

const startScanLoop = async () => {
  if (!cameraActive.value || !videoRef.value || isSuccessModalOpen.value || loading.value) return

  if ('BarcodeDetector' in window) {
    try {
      const barcodeDetector = new (window as any).BarcodeDetector({ formats: ['qr_code'] })
      const barcodes = await barcodeDetector.detect(videoRef.value)
      if (barcodes && barcodes.length > 0) {
        scannedToken.value = barcodes[0].rawValue
        submitScanPresensi()
        return
      }
    } catch (e) {
      // Loop Frame
    }
  }

  if (cameraActive.value && !isSuccessModalOpen.value) {
    scanLoopId = requestAnimationFrame(startScanLoop)
  }
}

// ===== ANTI-FAKE GPS & HARDWARE SENSOR VERIFICATION =====
const getGeolocation = () => {
  if ('geolocation' in navigator) {
    gpsLoading.value = true
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        gpsLoading.value = false
        const coords = pos.coords as any

        if (coords.isMock === true || coords.mocked === true || (pos as any).isMock === true) {
          isMockDetected.value = true
          mockReason.value = 'Aplikasi Fake GPS / Mock Location terdeteksi aktif pada perangkat HP Anda.'
          updateMapVisuals()
          return
        }

        if (coords.accuracy !== undefined && coords.accuracy === 0) {
          isMockDetected.value = true
          mockReason.value = 'Presisi sinyal GPS tidak valid (Accuracy = 0m). Terdeteksi manipulasi sinyal GPS.'
          updateMapVisuals()
          return
        }

        if (currentLat.value !== null && currentLng.value !== null && lastGpsTime.value > 0) {
          const timeDiff = (Date.now() - lastGpsTime.value) / 1000 // sec
          const distMoved = calculateDistance(currentLat.value, currentLng.value, pos.coords.latitude, pos.coords.longitude)
          if (timeDiff > 0 && timeDiff < 3) {
            const speedKmH = (distMoved / timeDiff) * 3.6
            if (speedKmH > 150) {
              isMockDetected.value = true
              mockReason.value = `Lompatan lokasi terdeteksi (${distMoved.toFixed(0)}m dlm ${timeDiff.toFixed(1)}s). Terindikasi Fake GPS.`
              updateMapVisuals()
              return
            }
          }
        }

        isMockDetected.value = false
        mockReason.value = ''
        lastGpsTime.value = Date.now()
        currentLat.value = parseFloat(pos.coords.latitude.toFixed(7))
        currentLng.value = parseFloat(pos.coords.longitude.toFixed(7))
        gpsAccuracy.value = coords.accuracy ? parseFloat(coords.accuracy.toFixed(1)) : 5.0
        updateMapVisuals()
      },
      (err) => {
        gpsLoading.value = false
        console.warn('GPS error, using default:', err)
        currentLat.value = -5.5645
        currentLng.value = 120.1945
        updateMapVisuals()
      },
      { 
        enableHighAccuracy: true,
        maximumAge: 0,
        timeout: 10000 
      }
    )
  }
}

// ===== UNIFIED SMART PRESENSI SUBMISSION =====
const submitScanPresensi = async () => {
  if (loading.value || isSuccessModalOpen.value) return

  // Prevent submit if already completed both MASUK & PULANG
  if (authStore.todayStatus?.is_masuk && authStore.todayStatus?.is_pulang) {
    actionError.value = `Presensi Ditolak! Anda telah menyelesaikan presensi Masuk & Pulang untuk siaga shift hari ini.`
    showToast('PRESENSI SELESAI', actionError.value, 'warning')
    return
  }

  // Prevent submit if shift is locked
  if (shiftState.value.isLocked && scanMode.value === 'MASUK' && !authStore.todayStatus?.is_masuk) {
    actionError.value = `Presensi Ditolak! ${shiftState.value.lockReason}`
    showToast('PRESENSI DITOLAK', actionError.value, 'error')
    return
  }

  // Prevent early clock out before window opens
  if (scanMode.value === 'PULANG' && !clockOutWindowInfo.value.isAllowed) {
    actionError.value = `Presensi Pulang Ditolak! ${clockOutWindowInfo.value.lockReason}`
    showToast('PRESENSI PULANG DITOLAK', actionError.value, 'error')
    return
  }

  if (isMockDetected.value) {
    actionError.value = 'Presensi Ditolak! Harap matikan aplikasi Fake GPS / Mock Location terlebih dahulu.'
    return
  }

  if (!isWithinRadius.value) {
    actionError.value = `Presensi Ditolak! Jarak Anda (${distanceMeters.value?.toFixed(1)}m) melebihi batas radius Posko (${maxRadius.value}m).`
    return
  }

  const token = scannedToken.value || poskoToken.value || 'LOPI-Q-POSKO-BULUKUMBA-2026-NTPD112'
  loading.value = true
  actionError.value = ''

  try {
    let res
    if (scanMode.value === 'PULANG') {
      res = await authStore.clockOut(currentLat.value || -5.5645, currentLng.value || 120.1945, token)
    } else {
      res = await authStore.clockIn(currentLat.value || -5.5645, currentLng.value || 120.1945, token)
    }

    if (res && res.success) {
      if ('vibrate' in navigator) {
        try {
          navigator.vibrate([120, 80, 120])
        } catch (e) {}
      }

      const msg = res.message || `Presensi ${scanMode.value} Siaga 112 Berhasil!`
      successMessage.value = msg
      showToast('PRESENSI BERHASIL', msg, 'success')
      isSuccessModalOpen.value = true
      stopCamera()

      let seconds = 2
      redirectCountdown.value = seconds
      const interval = setInterval(() => {
        seconds--
        redirectCountdown.value = seconds
        if (seconds <= 0) {
          clearInterval(interval)
          router.push('/intern/history')
        }
      }, 750)
    } else {
      actionError.value = res.error || 'Gagal mengirim presensi.'
      showToast('PRESENSI DITOLAK', actionError.value, 'error')
    }
  } catch (err: any) {
    actionError.value = err.response?.data?.error || 'Gagal terhubung ke server presensi.'
    showToast('KONEKSI GAGAL', actionError.value, 'error')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  getGeolocation()
  await fetchPoskoInfo()
  initLeafletMap()

  await authStore.fetchTodayStatus()
  await evaluateShiftSchedule()

  if (authStore.todayStatus?.is_masuk && authStore.todayStatus?.is_pulang) {
    showToast('PRESENSI HARI INI SELESAI', 'Tugas Siaga Hari ini telah selesai (Absen Masuk & Pulang Terverifikasi).', 'success')
  }

  if (!shiftState.value.isLocked || authStore.todayStatus?.is_masuk) {
    startCamera()
  }
})

onUnmounted(() => {
  stopCamera()
  if (mapInstance) {
    mapInstance.remove()
    mapInstance = null
  }
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.toast-slide-enter-active, .toast-slide-leave-active { transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1); }
.toast-slide-enter-from, .toast-slide-leave-to { opacity: 0; transform: translateY(-24px) scale(0.95); }
</style>
