<template>
  <InternLayout>
    <div class="min-h-screen flex flex-col font-sans text-[#1b1c1c] pb-safe relative select-none">

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

      <!-- ===== MAIN CONTENT AREA ===== -->
      <main class="flex-grow flex flex-col items-center justify-center px-4 py-6 sm:py-8 max-w-2xl mx-auto w-full space-y-6">

        <!-- Scanner Container (Soft Minimalism) -->
        <div class="w-full max-w-sm flex flex-col items-center space-y-4">
          
          <!-- Header Title Section -->
          <div class="text-center space-y-1">
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-extrabold uppercase tracking-wider bg-[#FCE4EC] text-[#ab2c5d] border border-[#F8BBD0] mb-1">
              <span class="w-1.5 h-1.5 rounded-full bg-[#f06292] animate-ping"></span>
              <span>MODE AUTO: {{ autoDetectedMode === 'MASUK' ? 'CHECK-IN MASUK' : (autoDetectedMode === 'PULANG' ? 'CHECK-OUT PULANG' : 'SELESAI') }}</span>
            </div>
            <h1 class="text-2xl sm:text-3xl font-bold text-[#1b1c1c]">
              {{ autoDetectedMode === 'MASUK' ? 'Scan to Check-in' : (autoDetectedMode === 'PULANG' ? 'Scan to Check-out' : 'Presensi Selesai') }}
            </h1>
            <p class="text-xs sm:text-sm text-[#574146]">
              Position the QR code within the frame.
            </p>
          </div>

          <!-- The Scanner 'Window' -->
          <div class="relative w-64 h-64 rounded-3xl overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.1)] border border-[#ddbfc5] bg-white flex items-center justify-center p-4">
            
            <!-- Corner Indicators (Simulated) -->
            <div class="absolute top-4 left-4 w-8 h-8 border-t-2 border-l-2 border-[#ab2c5d] rounded-tl-lg z-20"></div>
            <div class="absolute top-4 right-4 w-8 h-8 border-t-2 border-r-2 border-[#ab2c5d] rounded-tr-lg z-20"></div>
            <div class="absolute bottom-4 left-4 w-8 h-8 border-b-2 border-l-2 border-[#ab2c5d] rounded-bl-lg z-20"></div>
            <div class="absolute bottom-4 right-4 w-8 h-8 border-b-2 border-r-2 border-[#ab2c5d] rounded-br-lg z-20"></div>

            <!-- Video Stream Feed -->
            <video 
              ref="videoRef" 
              class="absolute inset-0 w-full h-full object-cover z-10" 
              autoplay 
              playsinline 
              v-show="cameraActive"
            ></video>

            <!-- Animated Scanning Line -->
            <div v-if="cameraActive" class="absolute w-[80%] h-0.5 bg-[#f06292] rounded-full shadow-[0_0_8px_#f06292] scan-line z-20"></div>

            <!-- Fallback Icon Watermark when camera inactive -->
            <div v-if="!cameraActive" class="flex flex-col items-center justify-center text-center space-y-2 relative z-10 p-4">
              <span class="material-symbols-outlined text-[#8a7176] opacity-20 text-6xl">qr_code</span>
              <p class="text-xs text-[#574146] font-medium">Klik tombol kamera di bawah untuk mulai scan.</p>
            </div>
          </div>

          <!-- Camera Shutter / Trigger Button -->
          <div class="flex flex-col items-center gap-2 pt-1">
            <button 
              @click="toggleCamera"
              class="w-16 h-16 rounded-full border-4 border-[#ab2c5d] p-1 flex items-center justify-center hover:bg-[#f06292]/10 transition-all active:scale-90 cursor-pointer shadow-md bg-white"
              :title="cameraActive ? 'Matikan Kamera' : 'Aktifkan Kamera'"
            >
              <div class="w-full h-full rounded-full bg-[#ab2c5d] flex items-center justify-center text-white">
                <span class="material-symbols-outlined text-2xl">
                  {{ cameraActive ? 'videocam_off' : 'photo_camera' }}
                </span>
              </div>
            </button>
            <span class="text-[11px] font-bold text-[#574146] uppercase tracking-wider">
              {{ cameraActive ? 'Matikan Kamera' : 'Buka Kamera Scan' }}
            </span>
          </div>

        </div>

        <!-- ===== GEOFENCE MAP CARD SECTION ===== -->
        <div class="w-full max-w-md bg-white rounded-3xl p-5 border border-[#F8BBD0] shadow-sm space-y-3">
          <div class="flex items-center justify-between border-b border-[#F8BBD0]/60 pb-2">
            <h3 class="text-sm font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">location_on</span>
              <span>Geofence Posko Siaga 112</span>
            </h3>
            <button 
              @click="getGeolocation" 
              class="text-xs text-[#f06292] hover:text-[#ab2c5d] font-bold bg-transparent border-0 cursor-pointer flex items-center gap-1"
            >
              <span class="material-symbols-outlined text-sm" :class="{ 'animate-spin': gpsLoading }">refresh</span>
              <span>Refresh</span>
            </button>
          </div>

          <!-- Leaflet Map Viewport Container -->
          <div id="scanLeafletMap" class="w-full h-44 rounded-2xl border border-[#ddbfc5] overflow-hidden shadow-inner relative z-0"></div>

          <!-- Geofence Info Badges -->
          <div class="flex items-center justify-between text-xs bg-[#FFF5F8] p-3 rounded-xl border border-[#F8BBD0]/60">
            <div class="flex items-center gap-1.5 text-[#574146]">
              <span class="font-medium">Jarak Posko:</span>
              <span class="font-mono font-extrabold text-[#ab2c5d]">
                {{ distanceMeters !== null ? distanceMeters.toFixed(1) + ' m' : 'Memuat...' }}
              </span>
            </div>
            <span 
              class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border"
              :class="isWithinRadius ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-rose-50 text-rose-700 border-rose-200'"
            >
              {{ isWithinRadius ? 'Dalam Radius' : 'Di Luar Radius' }}
            </span>
          </div>
        </div>

      </main>

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

    </div>
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
    showToast('AKSES KAMERA', 'Kamera tidak dapat diakses.', 'warning')
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
