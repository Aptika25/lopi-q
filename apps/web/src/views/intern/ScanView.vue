<template>
  <InternLayout>
    <div class="max-w-4xl mx-auto space-y-6 select-none font-sans relative pb-safe">
      
      <!-- ===== FLOATING TOAST ALERT ===== -->
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

      <!-- ===== PAGE HEADER ===== -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-[#F8BBD0]/60 pb-4">
        <div>
          <h2 class="font-display font-black text-[#1b1c1c] text-lg sm:text-xl flex items-center gap-2">
            <span class="material-symbols-outlined text-[#f06292] text-2xl">qr_code_scanner</span>
            <span>Scan QR Code Presensi Siaga 112</span>
          </h2>
          <p class="font-sans text-[#8a7176] mt-1 text-xs">Pindai QR Code resmi Posko Siaga 112 untuk mencatat Presensi Masuk atau Pulang Tugas.</p>
        </div>
        <div>
          <span class="px-3 py-1 bg-[#FCE4EC] text-[#ab2c5d] text-xs font-extrabold rounded-full border border-[#F8BBD0]">
            PRESENSI DIGITAL MAGANG
          </span>
        </div>
      </div>

      <!-- ===== AUTO-DETECT PRESENSI MODE BANNER ===== -->
      <div 
        class="rounded-3xl p-5 border shadow-sm flex items-center justify-between gap-4 transition-all"
        :class="autoDetectedMode === 'MASUK' ? 'bg-emerald-50/90 border-emerald-200 text-emerald-950' : (autoDetectedMode === 'PULANG' ? 'bg-amber-50/90 border-amber-200 text-amber-950' : 'bg-[#FFF5F8] border-[#F8BBD0] text-[#1b1c1c]')"
      >
        <div class="flex items-center gap-4">
          <div 
            class="w-12 h-12 rounded-2xl flex items-center justify-center font-bold text-white shrink-0 shadow-md"
            :class="autoDetectedMode === 'MASUK' ? 'bg-gradient-to-tr from-emerald-700 to-teal-600' : (autoDetectedMode === 'PULANG' ? 'bg-gradient-to-tr from-[#ab2c5d] to-[#f06292]' : 'bg-slate-600')"
          >
            <span class="material-symbols-outlined text-2xl">
              {{ autoDetectedMode === 'MASUK' ? 'login' : (autoDetectedMode === 'PULANG' ? 'logout' : 'task_alt') }}
            </span>
          </div>
          <div class="space-y-0.5">
            <div class="text-[10px] font-mono font-black uppercase tracking-wider text-[#8a7176]">Deteksi Otomatis Sistem</div>
            <h4 class="text-sm font-black font-display flex items-center gap-2" :class="autoDetectedMode === 'MASUK' ? 'text-emerald-800' : (autoDetectedMode === 'PULANG' ? 'text-[#ab2c5d]' : 'text-[#1b1c1c]')">
              <span>{{ autoDetectedMode === 'MASUK' ? '🟢 MODE AUTO: PRESENSI MASUK' : (autoDetectedMode === 'PULANG' ? '🟠 MODE AUTO: PRESENSI PULANG' : '✅ PRESENSI HARI INI COMPLETED') }}</span>
            </h4>
            <p class="text-xs font-medium text-[#574146] leading-snug">
              <span v-if="autoDetectedMode === 'MASUK'">Sistem mendeteksi Anda belum absen Masuk. Pemindaian QR Code akan otomatis mencatat Waktu Masuk Tugas.</span>
              <span v-else-if="autoDetectedMode === 'PULANG'">Sistem mendeteksi Anda telah absen Masuk. Pemindaian QR Code akan otomatis mencatat Waktu Pulang Tugas.</span>
              <span v-else>Presensi Masuk &amp; Pulang untuk siklus siaga hari ini telah terverifikasi lengkap.</span>
            </p>
          </div>
        </div>

        <div class="hidden sm:flex flex-col items-end shrink-0">
          <span 
            class="px-3 py-1 text-[10px] font-extrabold rounded-full border uppercase tracking-wider shadow-2xs"
            :class="autoDetectedMode === 'MASUK' ? 'bg-emerald-100 text-emerald-800 border-emerald-300' : (autoDetectedMode === 'PULANG' ? 'bg-[#FCE4EC] text-[#ab2c5d] border-[#F8BBD0]' : 'bg-slate-200 text-slate-700 border-slate-300')"
          >
            Sistem Geofence 112
          </span>
        </div>
      </div>

      <!-- ===== MAIN SCANNER & GEOFENCE MAP GRID ===== -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
        
        <!-- ===== LEFT: CAMERA VIEWPORT & MANUAL TOKEN CARD (7 Cols) ===== -->
        <div class="md:col-span-7 bg-white rounded-3xl p-6 border border-[#F8BBD0] shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">photo_camera</span>
              <span>Pemindai Kamera Live</span>
            </h3>
            <span v-if="cameraActive" class="px-2.5 py-1 bg-emerald-50 text-emerald-700 text-[11px] font-extrabold rounded-full border border-emerald-200 flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-emerald-500 animate-ping"></span>
              <span>Kamera Aktif</span>
            </span>
          </div>

          <!-- Camera Stream Box -->
          <div class="relative w-full aspect-square sm:aspect-video bg-slate-900 rounded-2xl overflow-hidden flex items-center justify-center border border-slate-800 shadow-inner">
            <video 
              ref="videoRef" 
              class="w-full h-full object-cover" 
              autoplay 
              playsinline 
              v-show="cameraActive"
            ></video>

            <!-- Scanning Framing Overlay -->
            <div v-if="cameraActive" class="absolute inset-0 pointer-events-none flex items-center justify-center">
              <div class="w-48 h-48 sm:w-56 sm:h-56 border-2 border-[#f06292] rounded-2xl relative shadow-[0_0_30px_rgba(240,98,146,0.5)] animate-pulse">
                <div class="absolute -top-1 -left-1 w-6 h-6 border-t-4 border-l-4 border-[#ab2c5d]"></div>
                <div class="absolute -top-1 -right-1 w-6 h-6 border-t-4 border-r-4 border-[#ab2c5d]"></div>
                <div class="absolute -bottom-1 -left-1 w-6 h-6 border-b-4 border-l-4 border-[#ab2c5d]"></div>
                <div class="absolute -bottom-1 -right-1 w-6 h-6 border-b-4 border-r-4 border-[#ab2c5d]"></div>
              </div>
            </div>

            <!-- Inactive Placeholder -->
            <div v-if="!cameraActive" class="text-center p-6 space-y-3">
              <div class="w-16 h-16 bg-[#FCE4EC] rounded-full flex items-center justify-center mx-auto text-[#ab2c5d]">
                <span class="material-symbols-outlined text-3xl">videocam_off</span>
              </div>
              <div>
                <p class="text-xs font-bold text-slate-300">Kamera Belum Aktif</p>
                <p class="text-[11px] text-slate-400 max-w-xs mx-auto mt-1">
                  Klik tombol "Buka Kamera Scan" untuk pemindaian otomatis.
                </p>
              </div>
            </div>
          </div>

          <!-- Camera Shutter Button -->
          <div class="flex flex-col items-center gap-2 pt-1">
            <button 
              @click="toggleCamera"
              class="w-14 h-14 rounded-full border-4 border-[#ab2c5d] p-1 flex items-center justify-center hover:bg-[#f06292]/10 transition-all active:scale-90 cursor-pointer shadow-md bg-white"
              :title="cameraActive ? 'Matikan Kamera' : 'Aktifkan Kamera'"
            >
              <div class="w-full h-full rounded-full bg-[#ab2c5d] flex items-center justify-center text-white">
                <span class="material-symbols-outlined text-xl">
                  {{ cameraActive ? 'videocam_off' : 'photo_camera' }}
                </span>
              </div>
            </button>
            <span class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">
              {{ cameraActive ? 'Matikan Kamera' : 'Buka Kamera Scan' }}
            </span>
          </div>

          <!-- Manual Token Input Section -->
          <div class="pt-3 border-t border-[#F8BBD0]/60 space-y-2">
            <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider flex items-center gap-1">
              <span class="material-symbols-outlined text-xs text-[#f06292]">key</span>
              <span>Input Token QR Manual</span>
            </label>
            <div class="flex gap-2">
              <input 
                v-model="manualTokenInput" 
                type="text" 
                placeholder="Tempel / ketik token QR posko..." 
                class="flex-1 px-3 py-2 border border-[#ddbfc5] rounded-xl text-xs focus:outline-none focus:border-[#f06292] font-mono bg-[#FFF5F8]"
              />
              <button 
                @click="submitManualToken" 
                :disabled="loading || !manualTokenInput.trim()"
                class="px-4 py-2 bg-[#ab2c5d] hover:bg-[#881b47] disabled:bg-slate-300 text-white font-bold text-xs rounded-xl border-0 cursor-pointer shadow-xs transition-all"
              >
                Kirim
              </button>
            </div>
          </div>
        </div>

        <!-- ===== RIGHT: GEOFENCE GPS & LEAFLET MAP CARD (5 Cols) ===== -->
        <div class="md:col-span-5 bg-white rounded-3xl p-6 border border-[#F8BBD0] shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">location_on</span>
              <span>Status Lokasi Geofence</span>
            </h3>
            <button 
              @click="getGeolocation" 
              class="text-xs text-[#f06292] hover:text-[#ab2c5d] font-bold bg-transparent border-0 cursor-pointer flex items-center gap-1"
            >
              <span class="material-symbols-outlined text-sm" :class="{ 'animate-spin': gpsLoading }">refresh</span>
              <span>Refresh GPS</span>
            </button>
          </div>

          <!-- Leaflet Interactive Map Viewport Container -->
          <div id="scanLeafletMap" class="w-full h-48 rounded-2xl border border-[#ddbfc5] overflow-hidden shadow-inner relative z-0"></div>

          <!-- Location Metadata Badges -->
          <div class="p-4 bg-[#FFF5F8] rounded-2xl border border-[#F8BBD0]/60 space-y-2.5 text-xs">
            <div class="flex items-center justify-between">
              <span class="text-[#574146] font-medium">Jarak ke Posko 112:</span>
              <span class="font-mono font-extrabold text-[#ab2c5d]">
                {{ distanceMeters !== null ? distanceMeters.toFixed(1) + ' meter' : 'Memuat...' }}
              </span>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-[#574146] font-medium">Batas Maksimal Radius:</span>
              <span class="font-mono font-bold text-[#1b1c1c]">{{ maxRadius }} meter</span>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-[#574146] font-medium">Status Geofence:</span>
              <span 
                class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border"
                :class="isWithinRadius ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200'"
              >
                {{ isWithinRadius ? 'Dalam Radius' : 'Di Luar Radius' }}
              </span>
            </div>
          </div>
        </div>

      </div>

    </div>

    <!-- ===== DRAGGABLE FAB CONTAINER FOR IZIN/SAKIT ===== -->
    <div id="draggable-fab-container" class="fixed right-4 z-40 flex flex-col items-end gap-2 select-none bottom-24">
      <!-- Expanded Option Pills -->
      <transition name="fade">
        <div v-if="showCutiOptions" class="flex flex-col gap-2 mb-2">
          <button 
            @click="openLeaveForm('IZIN')"
            class="bg-white border border-[#ddbfc5] text-[#ab2c5d] px-5 py-2 rounded-full shadow-lg font-bold text-xs uppercase hover:bg-[#FCE4EC] transition-all cursor-pointer"
          >
            Izin
          </button>
          <button 
            @click="openLeaveForm('SAKIT')"
            class="bg-white border border-[#ddbfc5] text-[#ab2c5d] px-5 py-2 rounded-full shadow-lg font-bold text-xs uppercase hover:bg-[#FCE4EC] transition-all cursor-pointer"
          >
            Sakit
          </button>
        </div>
      </transition>

      <!-- FAB Button Trigger -->
      <button 
        @click="showCutiOptions = !showCutiOptions"
        class="bg-[#f06292] hover:bg-[#ab2c5d] text-white rounded-full shadow-[0_8px_30px_rgba(240,98,146,0.4)] flex items-center justify-center transition-all active:scale-90 w-12 h-12 border-2 border-white cursor-pointer"
        title="Pengajuan Izin / Sakit"
      >
        <span class="material-symbols-outlined text-2xl">
          {{ showCutiOptions ? 'close' : 'add' }}
        </span>
      </button>
    </div>

    <!-- ===== FORM PENGAJUAN OVERLAY BOTTOM SHEET ===== -->
    <transition name="fade">
      <div v-if="showLeaveOverlay" class="fixed inset-0 bg-slate-900/50 z-50 flex items-end justify-center p-0 backdrop-blur-xs" @click.self="showLeaveOverlay = false">
        <div class="w-full max-w-md bg-white rounded-t-3xl p-6 space-y-4 shadow-2xl border-t border-[#ddbfc5]">
          
          <div class="flex justify-between items-center border-b border-[#ddbfc5]/60 pb-3">
            <h2 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">edit_document</span>
              <span>Form Pengajuan {{ formLeave.category }}</span>
            </h2>
            <button @click="showLeaveOverlay = false" class="p-1 rounded-full hover:bg-slate-100 border-0 bg-transparent cursor-pointer text-slate-400">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <form @submit.prevent="submitLeaveRequest" class="space-y-4 text-xs">
            <!-- Kategori -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Kategori Pengajuan</label>
              <select v-model="formLeave.category" class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none">
                <option value="IZIN">Izin Keperluan Mendesak</option>
                <option value="SAKIT">Sakit (Dengan Surat Keterangan)</option>
              </select>
            </div>

            <!-- Tanggal -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tanggal Shift <span class="text-rose-500">*</span></label>
              <input v-model="formLeave.shift_date" type="date" required class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none" />
            </div>

            <!-- Alasan -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Alasan Detail <span class="text-rose-500">*</span></label>
              <textarea v-model="formLeave.reason" rows="3" required placeholder="Tuliskan alasan detail..." class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none resize-none"></textarea>
            </div>

            <div v-if="formMessage" class="p-3 bg-emerald-50 text-emerald-800 rounded-xl font-bold border border-emerald-200">
              {{ formMessage }}
            </div>

            <!-- Kirim Button -->
            <button type="submit" :disabled="submitting" class="w-full py-3.5 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-full font-bold shadow-lg active:scale-95 transition-all border-0 cursor-pointer">
              Kirim Pengajuan
            </button>
          </form>

        </div>
      </div>
    </transition>

    <!-- ===== SUCCESS VERIFICATION DIALOG ===== -->
    <transition name="fade">
      <div v-if="isSuccessModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
        <div class="bg-white rounded-3xl max-w-sm w-full p-6 text-center border border-[#F8BBD0] shadow-2xl space-y-4">
          <div class="w-16 h-16 rounded-full bg-emerald-100 text-emerald-600 flex items-center justify-center mx-auto border border-emerald-200 shadow-md">
            <span class="material-symbols-outlined text-3xl">task_alt</span>
          </div>
          
          <div class="space-y-1">
            <h3 class="font-bold text-lg text-[#1b1c1c]">Presensi Berhasil!</h3>
            <p class="text-xs text-[#574146] leading-relaxed">{{ successMessage }}</p>
          </div>

          <div class="p-3 bg-[#FFF5F8] rounded-xl text-xs text-[#ab2c5d] font-mono font-bold border border-[#F8BBD0]/60">
            Mengalihkan ke Riwayat dalam {{ redirectCountdown }}s...
          </div>
        </div>
      </div>
    </transition>
  </InternLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import InternLayout from '@/layouts/InternLayout.vue'

const router = useRouter()
const authStore = useAuthStore()

// ===== STATE =====
const isSuccessModalOpen = ref(false)
const successMessage = ref('')
const redirectCountdown = ref(2)
const manualTokenInput = ref('')
const loading = ref(false)

const showCutiOptions = ref(false)
const showLeaveOverlay = ref(false)
const submitting = ref(false)
const formMessage = ref('')

const formLeave = ref({
  category: 'IZIN',
  shift_date: new Date().toISOString().substring(0, 10),
  reason: ''
})

const cameraActive = ref(false)
const videoRef = ref<HTMLVideoElement | null>(null)
const mediaStream = ref<MediaStream | null>(null)
const scannedToken = ref('')

const currentLat = ref<number | null>(-5.5645)
const currentLng = ref<number | null>(120.1945)
const poskoLat = ref(-5.5645)
const poskoLng = ref(120.1945)
const maxRadius = ref(2.0)
const poskoToken = ref('')
const gpsLoading = ref(false)

let scanLoopId: number | null = null
let mapInstance: any = null
let poskoMarker: any = null
let userMarker: any = null
let circleInstance: any = null
let polylineInstance: any = null

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
  }, 5000)
}

// ===== COMPUTEDS =====
const autoDetectedMode = computed(() => {
  if (!authStore.todayStatus?.is_masuk) return 'MASUK'
  if (authStore.todayStatus?.is_masuk && !authStore.todayStatus?.is_pulang) return 'PULANG'
  return 'SELESAI'
})

const distanceMeters = computed(() => {
  if (currentLat.value === null || currentLng.value === null) return null
  return calculateDistance(currentLat.value, currentLng.value, poskoLat.value, poskoLng.value)
})

const isWithinRadius = computed(() => {
  if (distanceMeters.value === null) return false
  return distanceMeters.value <= maxRadius.value
})

// ===== HAVERSINE MATH =====
function calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371000
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLon = (lon2 - lon1) * Math.PI / 180
  const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
            Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
            Math.sin(dLon / 2) * Math.sin(dLon / 2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
  return R * c
}

// ===== LEAFLET MAP =====
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
      html: `<div style="background-color:#ab2c5d;color:white;width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 10px rgba(171,44,93,0.4);border:2px solid white;"><span class="material-symbols-outlined" style="font-size:18px;">local_hospital</span></div>`,
      iconSize: [32, 32],
      iconAnchor: [16, 16]
    })

    const userIcon = L.divIcon({
      className: 'custom-user-pin',
      html: `<div style="background-color:#f06292;color:white;width:32px;height:32px;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 10px rgba(240,98,146,0.4);border:2px solid white;"><span class="material-symbols-outlined" style="font-size:18px;">person_pin_circle</span></div>`,
      iconSize: [32, 32],
      iconAnchor: [16, 16]
    })

    poskoMarker = L.marker([poskoLat.value, poskoLng.value], { icon: poskoIcon }).addTo(mapInstance)
      .bindPopup(`<b>Posko Siaga 112 Bulukumba</b><br>Radius: ${maxRadius.value}m`)

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

  if (userMarker) userMarker.setLatLng([userLat, userLng])
  if (poskoMarker) poskoMarker.setLatLng([poskoLat.value, poskoLng.value])

  if (circleInstance) mapInstance.removeLayer(circleInstance)
  const circleColor = isWithinRadius.value ? '#10b981' : '#f43f5e'
  circleInstance = L.circle([poskoLat.value, poskoLng.value], {
    color: circleColor,
    fillColor: circleColor,
    fillOpacity: 0.25,
    radius: maxRadius.value
  }).addTo(mapInstance)

  if (polylineInstance) mapInstance.removeLayer(polylineInstance)
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

// ===== CAMERA =====
const toggleCamera = () => {
  if (cameraActive.value) {
    stopCamera()
  } else {
    startCamera()
  }
}

const startCamera = async () => {
  try {
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
    showToast('AKSES KAMERA', 'Kamera tidak dapat diakses. Gunakan input token manual.', 'warning')
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
    } catch (e) {}
  }

  if (cameraActive.value && !isSuccessModalOpen.value) {
    scanLoopId = requestAnimationFrame(startScanLoop)
  }
}

const submitManualToken = () => {
  if (!manualTokenInput.value.trim()) return
  scannedToken.value = manualTokenInput.value.trim()
  submitScanPresensi()
}

// ===== GEOLOCATION =====
const getGeolocation = () => {
  if ('geolocation' in navigator) {
    gpsLoading.value = true
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        gpsLoading.value = false
        currentLat.value = parseFloat(pos.coords.latitude.toFixed(7))
        currentLng.value = parseFloat(pos.coords.longitude.toFixed(7))
        updateMapVisuals()
      },
      () => {
        gpsLoading.value = false
        currentLat.value = -5.5645
        currentLng.value = 120.1945
        updateMapVisuals()
      },
      { enableHighAccuracy: true, timeout: 10000 }
    )
  }
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
  } catch (e) {}
}

// ===== SUBMIT PRESENSI =====
const submitScanPresensi = async () => {
  if (loading.value || isSuccessModalOpen.value) return

  if (autoDetectedMode.value === 'SELESAI') {
    showToast('PRESENSI SELESAI', 'Anda telah menyelesaikan presensi Masuk & Pulang untuk siaga hari ini.', 'warning')
    return
  }

  if (!isWithinRadius.value) {
    showToast('PRESENSI DITOLAK', `Jarak Anda (${distanceMeters.value?.toFixed(1)}m) melebihi radius Posko (${maxRadius.value}m).`, 'error')
    return
  }

  const token = scannedToken.value || manualTokenInput.value || poskoToken.value || 'LOPI-Q-POSKO-BULUKUMBA-2026-NTPD112'
  loading.value = true

  try {
    let res
    const mode = autoDetectedMode.value
    if (mode === 'PULANG') {
      res = await authStore.clockOut(currentLat.value || -5.5645, currentLng.value || 120.1945, token)
    } else {
      res = await authStore.clockIn(currentLat.value || -5.5645, currentLng.value || 120.1945, token)
    }

    if (res && res.success) {
      const msg = res.message || `Presensi ${mode} Berhasil!`
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
      showToast('PRESENSI DITOLAK', res.error || 'Gagal verifikasi presensi.', 'error')
    }
  } catch (err: any) {
    showToast('KONEKSI GAGAL', err.response?.data?.error || 'Gagal terhubung ke server presensi.', 'error')
  } finally {
    loading.value = false
  }
}

// ===== LEAVE FORM OVERLAY =====
const openLeaveForm = (category: string) => {
  formLeave.value.category = category
  showCutiOptions.value = false
  showLeaveOverlay.value = true
}

const submitLeaveRequest = async () => {
  submitting.value = true
  formMessage.value = ''
  try {
    let currentList: any[] = []
    try {
      const resGet = await axios.get('/api/presensi/leave-requests')
      if (resGet.data && Array.isArray(resGet.data.requests)) {
        currentList = resGet.data.requests
      }
    } catch (e) {}

    const newReq = {
      id: Date.now(),
      created_at: new Date().toISOString().replace('T', ' ').substring(0, 16),
      user_name: authStore.user?.name || 'Peserta Magang',
      user_nip: authStore.user?.nip || '-',
      category: formLeave.value.category,
      shift_date: formLeave.value.shift_date,
      reason: formLeave.value.reason,
      status: 'PENDING'
    }

    currentList.unshift(newReq)
    await axios.put('/api/presensi/leave-requests', currentList)

    formMessage.value = 'Pengajuan berhasil dikirim! Menunggu konfirmasi Super Admin.'
    setTimeout(() => {
      showLeaveOverlay.value = false
      formMessage.value = ''
      formLeave.value.reason = ''
    }, 1500)
  } catch (err) {
    formMessage.value = 'Gagal mengirim pengajuan.'
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  getGeolocation()
  await fetchPoskoInfo()
  initLeafletMap()
  await authStore.fetchTodayStatus()
  startCamera()
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
.material-symbols-outlined { 
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; 
}
@keyframes scan {
  0% { transform: translateY(-100%); opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { transform: translateY(100%); opacity: 0; }
}
.scan-line {
  animation: scan 2s linear infinite;
}
.pb-safe { 
  padding-bottom: env(safe-area-inset-bottom, 80px); 
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.toast-slide-enter-active, .toast-slide-leave-active { transition: all 0.35s cubic-bezier(0.34, 1.56, 0.64, 1); }
.toast-slide-enter-from, .toast-slide-leave-to { opacity: 0; transform: translateY(-24px) scale(0.95); }
</style>
