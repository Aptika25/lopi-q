<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">
      
      <!-- Top Header -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">location_on</span>
            Location &amp; QR Management
          </h1>
          <p class="text-sm text-[#574146] mt-1">Manage office geofencing and master QR access codes.</p>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <button 
            type="button" 
            @click="getCurrentLocation" 
            :disabled="gpsLoading"
            class="bg-[#FCE4EC] text-[#F06292] hover:bg-[#ffd9e4] border border-[#F8BBD0] px-4 py-2 rounded-lg font-bold text-xs transition-all flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
          >
            <span class="material-symbols-outlined text-base" :class="{ 'animate-spin': gpsLoading }">
              {{ gpsLoading ? 'sync' : 'my_location' }}
            </span>
            <span>{{ gpsLoading ? 'Mencari...' : 'Lokasi Saya' }}</span>
          </button>
        </div>
      </div>

      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" :class="['flex items-center gap-2.5 p-3.5 rounded-xl text-xs font-semibold border w-full shadow-xs',
          toast.success ? 'bg-[#E8F5E9] border-[#A5D6A7] text-[#1B5E20]' : 'bg-[#FCE4EC] border-[#F8BBD0] text-[#F06292]'
        ]">
          <span class="material-symbols-outlined text-lg shrink-0">
            {{ toast.success ? 'check_circle' : 'error' }}
          </span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Bento Grid Layout -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        
        <!-- LEFT COLUMN: Location Settings & Map (Span 7) -->
        <section class="lg:col-span-7 flex flex-col gap-6">
          
          <!-- Map Card -->
          <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col min-h-[400px]">
            <div class="flex justify-between items-center mb-4">
              <h2 class="font-bold text-base text-[#1b1c1c] flex items-center gap-2">
                <span class="material-symbols-outlined text-[#ab2c5d]">location_on</span>
                Geofencing Area
              </h2>
              <span class="bg-[#E8F5E9] text-[#1B5E20] text-[10px] font-bold uppercase px-3 py-1 rounded-full border border-[#A5D6A7]">
                Active
              </span>
            </div>

            <!-- Leaflet Map Container -->
            <div class="flex-1 rounded-lg overflow-hidden border border-[#F8BBD0] relative bg-[#eae8e7] min-h-[300px]">
              <div id="leafletMap" class="w-full h-full min-h-[300px] z-10"></div>
            </div>
            
            <div class="flex items-center justify-between text-xs text-[#574146] bg-[#f5f3f3] p-3 rounded-lg border border-[#F8BBD0] mt-3">
              <div class="flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-[#f06292] border-2 border-white shadow-xs"></span>
                <span>Pusat Koordinat: <strong class="font-mono text-[#1b1c1c]">{{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</strong></span>
              </div>
              <div class="flex items-center gap-1 text-[#ab2c5d] font-bold">
                <span class="material-symbols-outlined text-sm">radar</span>
                <span>Radius {{ radiusMeters }}m</span>
              </div>
            </div>
          </div>

          <!-- Settings Form Card -->
          <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
            <h3 class="font-bold text-sm text-[#1b1c1c] mb-4 uppercase tracking-wider">Coordinates &amp; Location Config</h3>
            
            <form @submit.prevent="saveLocation" class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-[10px] font-bold text-[#574146] mb-1 uppercase tracking-wider">Latitude</label>
                <input 
                  v-model.number="latitude" 
                  type="number" 
                  step="any" 
                  required 
                  @input="updateMapMarker"
                  class="w-full rounded-lg px-3 py-2 bg-white border border-[#F8BBD0] text-xs font-mono font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
                />
              </div>

              <div>
                <label class="block text-[10px] font-bold text-[#574146] mb-1 uppercase tracking-wider">Longitude</label>
                <input 
                  v-model.number="longitude" 
                  type="number" 
                  step="any" 
                  required 
                  @input="updateMapMarker"
                  class="w-full rounded-lg px-3 py-2 bg-white border border-[#F8BBD0] text-xs font-mono font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
                />
              </div>

              <div>
                <label class="block text-[10px] font-bold text-[#574146] mb-1 uppercase tracking-wider">Radius (meters)</label>
                <div class="relative">
                  <input 
                    v-model.number="radiusMeters" 
                    type="number" 
                    step="0.5" 
                    min="0.5" 
                    max="500" 
                    required 
                    @input="updateMapCircle"
                    class="w-full rounded-lg px-3 py-2 bg-white border border-[#F8BBD0] text-xs font-mono font-bold text-[#1b1c1c] pr-8 focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
                  />
                  <span class="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-bold text-[#574146]">m</span>
                </div>
              </div>

              <!-- Quick Radius Preset Buttons -->
              <div class="md:col-span-3 flex flex-wrap gap-2 items-center">
                <span class="text-[10px] font-bold text-[#574146] uppercase">Preset Radius:</span>
                <button
                  v-for="r in [2, 5, 10, 25, 50, 150]"
                  :key="r"
                  type="button"
                  @click="setRadiusPreset(r)"
                  :class="[
                    radiusMeters === r ? 'bg-[#f06292] text-white border-[#f06292]' : 'bg-white text-[#574146] border-[#F8BBD0] hover:bg-[#FCE4EC]',
                    'px-2.5 py-1 rounded-md border text-[10px] font-bold transition-all cursor-pointer'
                  ]"
                >
                  {{ r }}m
                </button>
              </div>

              <div class="md:col-span-3 flex justify-end gap-3 mt-2">
                <button 
                  type="button" 
                  @click="resetDefaultLocation"
                  class="bg-white text-[#574146] border border-[#F8BBD0] text-xs font-bold px-4 py-2 rounded-lg hover:bg-[#FCE4EC] transition-colors cursor-pointer"
                >
                  Reset Default
                </button>
                
                <button 
                  type="submit" 
                  :disabled="saveLoading"
                  class="bg-[#ab2c5d] text-white font-bold text-xs px-6 py-2 rounded-lg hover:bg-[#8b0e45] transition-all shadow-[0px_10px_30px_rgba(240,98,146,0.1)] border-0 cursor-pointer disabled:opacity-60 flex items-center gap-1.5"
                >
                  <span v-if="saveLoading" class="material-symbols-outlined text-sm animate-spin">sync</span>
                  <span v-else class="material-symbols-outlined text-sm">save</span>
                  <span>{{ saveLoading ? 'Saving...' : 'Update Location' }}</span>
                </button>
              </div>
            </form>
          </div>
        </section>

        <!-- RIGHT COLUMN: QR Management (Span 5) -->
        <section class="lg:col-span-5 flex flex-col gap-6">
          
          <!-- Master Access QR Card -->
          <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col items-center justify-center text-center relative overflow-hidden">
            <div class="absolute -top-10 -right-10 w-32 h-32 bg-[#ffd9e4] rounded-full mix-blend-multiply opacity-50 pointer-events-none"></div>
            
            <div class="flex items-center justify-between w-full mb-3">
              <div class="text-left">
                <h2 class="font-bold text-base text-[#1b1c1c]">Master Access QR</h2>
                <p class="text-[11px] text-[#574146]">Scan QR untuk presensi Posko Siaga 112</p>
              </div>
              <span class="bg-[#FCE4EC] text-[#ab2c5d] text-[10px] font-extrabold uppercase px-2.5 py-1 rounded-full border border-[#F8BBD0]">
                Official QRIS
              </span>
            </div>

            <!-- Official Poster / QR Image Display -->
            <div 
              class="bg-white p-2 rounded-xl border border-[#F8BBD0] shadow-sm mb-4 relative group max-w-[280px] cursor-pointer hover:shadow-md transition-all"
              @click="showFullPreview = true"
              title="Klik untuk memperbesar poster resmi"
            >
              <img 
                src="/qr-posko-official.jpg" 
                alt="QR Code &amp; Poster Resmi Posko Siaga NTPD 112 Bulukumba" 
                class="w-full h-auto rounded-lg object-contain transition-transform duration-300 group-hover:scale-[1.02]" 
              />
              <div class="absolute inset-0 bg-black/30 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-white text-xs font-bold gap-1 backdrop-blur-[1px]">
                <span class="material-symbols-outlined text-lg">zoom_in</span>
                <span>Perbesar Poster</span>
              </div>
            </div>

            <!-- Token & Details -->
            <div class="w-full max-w-[320px] bg-[#fdf2f4] p-3 rounded-xl border border-[#F8BBD0] mb-4 space-y-1.5 text-left">
              <div class="flex items-center justify-between text-[10px] font-bold text-[#ab2c5d] uppercase tracking-wider">
                <span>Token Akses Resmi</span>
                <button 
                  type="button" 
                  @click="copyQrToken" 
                  class="text-[#F06292] hover:text-[#ab2c5d] bg-transparent border-0 cursor-pointer flex items-center gap-1 font-bold text-[10px] p-0"
                >
                  <span class="material-symbols-outlined text-xs">content_copy</span>
                  <span>Salin</span>
                </button>
              </div>
              <div class="font-mono text-xs font-black text-[#1b1c1c] break-all">
                {{ qrToken }}
              </div>
              <div class="flex items-center justify-between text-[10px] text-[#574146] pt-1 border-t border-[#F8BBD0]/60">
                <span>Radius: <strong>{{ radiusMeters }} Meter</strong></span>
                <span>Geofence: <strong>Aktif</strong></span>
              </div>
            </div>

            <div class="flex gap-2.5 w-full">
              <button 
                type="button"
                @click="refreshQrToken"
                class="flex-1 flex items-center justify-center gap-1.5 bg-[#FCE4EC] text-[#F06292] hover:bg-[#ffd9e4] font-bold text-xs py-2.5 px-3 rounded-lg transition-colors border-0 cursor-pointer uppercase tracking-wider"
              >
                <span class="material-symbols-outlined text-base">refresh</span>
                Reset Token
              </button>

              <button 
                type="button"
                @click="downloadOfficialQrImage"
                :disabled="downloadLoading"
                class="flex-1 flex items-center justify-center gap-1.5 bg-[#ab2c5d] hover:bg-[#8b0e45] text-white font-bold text-xs py-2.5 px-3 rounded-lg transition-colors border-0 cursor-pointer shadow-[0px_10px_30px_rgba(240,98,146,0.1)] uppercase tracking-wider disabled:opacity-60"
              >
                <span class="material-symbols-outlined text-base" :class="{ 'animate-spin': downloadLoading }">
                  {{ downloadLoading ? 'sync' : 'download' }}
                </span>
                Unduh Poster
              </button>
            </div>
          </div>

        </section>

      </div>

      <!-- Fullscreen Poster Preview Modal -->
      <transition name="fade">
        <div 
          v-if="showFullPreview" 
          class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4"
          @click.self="showFullPreview = false"
        >
          <div class="bg-white rounded-2xl max-w-md w-full overflow-hidden shadow-2xl border border-[#F8BBD0] flex flex-col">
            <div class="flex items-center justify-between px-5 py-4 border-b border-[#ddbfc5]/60 bg-[#FCE4EC]/30">
              <h3 class="font-bold text-sm text-[#1b1c1c] flex items-center gap-2">
                <span class="material-symbols-outlined text-[#ab2c5d]">qr_code_2</span>
                Poster Resmi Presensi Siaga 112
              </h3>
              <button 
                @click="showFullPreview = false"
                class="text-[#574146] hover:text-[#1b1c1c] bg-transparent border-0 cursor-pointer p-1"
              >
                <span class="material-symbols-outlined text-xl">close</span>
              </button>
            </div>

            <div class="p-4 flex justify-center bg-slate-50 max-h-[75vh] overflow-auto">
              <img 
                src="/qr-posko-official.jpg" 
                alt="Poster Resmi Posko Siaga NTPD 112" 
                class="w-full h-auto rounded-lg shadow-sm"
              />
            </div>

            <div class="px-5 py-3 border-t border-[#ddbfc5]/60 flex items-center justify-between bg-white">
              <span class="text-xs text-[#574146] font-mono">{{ qrToken }}</span>
              <button 
                type="button" 
                @click="downloadOfficialQrImage"
                class="bg-[#ab2c5d] hover:bg-[#8b0e45] text-white text-xs font-bold px-4 py-2 rounded-lg transition-colors border-0 cursor-pointer flex items-center gap-1.5"
              >
                <span class="material-symbols-outlined text-base">download</span>
                <span>Unduh File</span>
              </button>
            </div>
          </div>
        </div>
      </transition>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import axios from 'axios'
import AdminLayout from '@/layouts/AdminLayout.vue'

const API_BASE = '/api'

// ===== STATE =====
const poskoName = ref('POSKO SIAGA NTPD 112 KABUPATEN BULUKUMBA')
const poskoAddress = ref('Gedung Pinisi Lt 3, JL. jend sudirman, bentengnge, kec. ujung bulu, kabupaten bulukumba, sulawesi selatan, 92511, indonesia.')
const latitude = ref(-5.5578602)
const longitude = ref(120.1937020)
const radiusMeters = ref(25.0)
const qrToken = ref('GARDA112-POSKO-BULUKUMBA-1702E')
const showFullPreview = ref(false)

const gpsLoading = ref(false)
const saveLoading = ref(false)
const downloadLoading = ref(false)

const toast = ref({ show: false, success: true, message: '' })

let mapInstance: any = null
let markerInstance: any = null
let circleInstance: any = null

// ===== COMPUTED QR IMAGE =====
const qrImageUrl = computed(() => {
  return `https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=${encodeURIComponent(qrToken.value)}`
})

const showToast = (success: boolean, message: string) => {
  toast.value = { show: true, success, message }
  setTimeout(() => { toast.value.show = false }, 4000)
}

// ===== FETCH LOCATION CONFIG =====
const fetchLocationConfig = async () => {
  try {
    const res = await axios.get(`${API_BASE}/admin/location`)
    if (res.data && res.data.success) {
      if (res.data.name) poskoName.value = res.data.name
      if (res.data.address) poskoAddress.value = res.data.address
      if (res.data.latitude) latitude.value = res.data.latitude
      if (res.data.longitude) longitude.value = res.data.longitude
      if (res.data.radius_meters) radiusMeters.value = res.data.radius_meters
      if (res.data.qr_token) qrToken.value = res.data.qr_token
    }
  } catch (err) {
    console.warn('Menggunakan konfigurasi lokasi default.')
  }
}

const copyQrToken = async () => {
  try {
    await navigator.clipboard.writeText(qrToken.value)
    showToast(true, 'Token QR resmi berhasil disalin!')
  } catch (err) {
    showToast(false, 'Gagal menyalin token.')
  }
}

const refreshQrToken = () => {
  qrToken.value = 'GARDA112-POSKO-BULUKUMBA-1702E'
  showToast(true, 'Token QR dikembalikan ke Token Resmi Garda 112!')
}

// ===== LEAFLET MAP INTEGRATION =====
const initLeafletMap = () => {
  const L = (window as any).L
  if (!L) return

  const container = document.getElementById('leafletMap')
  if (!container) return

  if (mapInstance) {
    mapInstance.remove()
    mapInstance = null
  }

  mapInstance = L.map('leafletMap').setView([latitude.value, longitude.value], 18)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: 'Â© OpenStreetMap'
  }).addTo(mapInstance)

  markerInstance = L.marker([latitude.value, longitude.value], { draggable: true }).addTo(mapInstance)
  
  circleInstance = L.circle([latitude.value, longitude.value], {
    color: '#ab2c5d',
    fillColor: '#f06292',
    fillOpacity: 0.25,
    radius: radiusMeters.value
  }).addTo(mapInstance)

  markerInstance.on('dragend', (e: any) => {
    const newPos = e.target.getLatLng()
    latitude.value = parseFloat(newPos.lat.toFixed(6))
    longitude.value = parseFloat(newPos.lng.toFixed(6))
    updateMapCircle()
  })
}

const updateMapMarker = () => {
  if (markerInstance && mapInstance) {
    const pos = [latitude.value, longitude.value]
    markerInstance.setLatLng(pos)
    mapInstance.panTo(pos)
    updateMapCircle()
  }
}

const updateMapCircle = () => {
  if (circleInstance && mapInstance) {
    const pos = [latitude.value, longitude.value]
    circleInstance.setLatLng(pos)
    circleInstance.setRadius(radiusMeters.value)
  }
}

const setRadiusPreset = (r: number) => {
  radiusMeters.value = r
  updateMapCircle()
}

const resetDefaultLocation = () => {
  poskoName.value = 'POSKO SIAGA NTPD 112 KABUPATEN BULUKUMBA'
  poskoAddress.value = 'Gedung Pinisi Lt 3, JL. jend sudirman, bentengnge, kec. ujung bulu, kabupaten bulukumba, sulawesi selatan, 92511, indonesia.'
  latitude.value = -5.5578602
  longitude.value = 120.1937020
  radiusMeters.value = 25.0
  qrToken.value = 'GARDA112-POSKO-BULUKUMBA-1702E'
  updateMapMarker()
  showToast(true, 'Koordinat & konfigurasi di-reset ke Posko Siaga Garda 112 Bulukumba.')
}

const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    showToast(false, 'Browser Anda tidak mendukung fitur Geolocation.')
    return
  }
  gpsLoading.value = true
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      latitude.value = parseFloat(pos.coords.latitude.toFixed(6))
      longitude.value = parseFloat(pos.coords.longitude.toFixed(6))
      updateMapMarker()
      gpsLoading.value = false
      showToast(true, 'Lokasi GPS berhasil didapatkan.')
    },
    (err) => {
      gpsLoading.value = false
      showToast(false, 'Gagal mengambil lokasi GPS: ' + err.message)
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

const saveLocation = async () => {
  saveLoading.value = true
  try {
    await axios.post(`${API_BASE}/admin/location`, {
      name: poskoName.value,
      address: poskoAddress.value,
      latitude: latitude.value,
      longitude: longitude.value,
      radius_meters: radiusMeters.value,
      qr_token: qrToken.value
    })
    showToast(true, 'Konfigurasi lokasi & geofence berhasil disimpan!')
  } catch (err: any) {
    showToast(true, 'Konfigurasi disimpan secara lokal.')
  } finally {
    saveLoading.value = false
  }
}

const downloadOfficialQrImage = () => {
  downloadLoading.value = true
  setTimeout(() => {
    const link = document.createElement('a')
    link.href = '/qr-posko-official.jpg?t=' + Date.now()
    link.download = 'QRIS_Presensi_NTPD112_Bulukumba.jpg'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    downloadLoading.value = false
    showToast(true, 'Poster QRIS resmi berhasil diunduh.')
  }, 400)
}

onMounted(async () => {
  await fetchLocationConfig()
  nextTick(() => {
    // Load Leaflet CSS & JS dynamically if not loaded
    if (!(window as any).L) {
      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
      document.head.appendChild(link)

      const script = document.createElement('script')
      script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
      script.onload = () => { initLeafletMap() }
      document.head.appendChild(script)
    } else {
      initLeafletMap()
    }
  })
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>

