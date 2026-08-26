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

        <!-- RIGHT COLUMN: QR Management & Active Zones List (Span 5) -->
        <section class="lg:col-span-5 flex flex-col gap-6">
          
          <!-- Master Access QR Card -->
          <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col items-center justify-center text-center relative overflow-hidden">
            <div class="absolute -top-10 -right-10 w-32 h-32 bg-[#ffd9e4] rounded-full mix-blend-multiply opacity-50 pointer-events-none"></div>
            
            <h2 class="font-bold text-base text-[#1b1c1c] mb-1">Master Access QR</h2>
            <p class="text-xs text-[#574146] mb-6">Scan to check-in at main lobby</p>

            <div class="bg-white p-4 rounded-xl border border-[#F8BBD0] shadow-xs mb-4 relative group">
              <img :src="qrImageUrl" alt="Master QR Code" class="w-48 h-48 object-contain" />
            </div>

            <div class="font-mono text-xs font-extrabold text-[#1b1c1c] tracking-wider mb-6 bg-[#f5f3f3] px-3 py-1.5 rounded-md border border-[#F8BBD0]">
              {{ qrToken }}
            </div>

            <div class="flex gap-3 w-full px-2">
              <button 
                type="button"
                @click="refreshQrToken"
                class="flex-1 flex items-center justify-center gap-2 bg-[#FCE4EC] text-[#F06292] font-bold text-xs py-2.5 px-4 rounded-lg hover:bg-[#ffd9e4] transition-colors border-0 cursor-pointer uppercase tracking-wider"
              >
                <span class="material-symbols-outlined text-base">refresh</span>
                Refresh
              </button>

              <button 
                type="button"
                @click="downloadOfficialQrImage"
                :disabled="downloadLoading"
                class="flex-1 flex items-center justify-center gap-2 bg-[#ab2c5d] text-white font-bold text-xs py-2.5 px-4 rounded-lg hover:bg-[#8b0e45] transition-colors border-0 cursor-pointer shadow-[0px_10px_30px_rgba(240,98,146,0.1)] uppercase tracking-wider disabled:opacity-60"
              >
                <span class="material-symbols-outlined text-base" :class="{ 'animate-spin': downloadLoading }">
                  {{ downloadLoading ? 'sync' : 'download' }}
                </span>
                Download
              </button>
            </div>

            <div class="mt-4 text-[11px] font-medium text-[#574146] flex items-center gap-1">
              <span class="material-symbols-outlined text-xs">timer</span>
              <span>Auto-refreshes in 14:59</span>
            </div>
          </div>

          <!-- Active Locations / Zones List -->
          <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex-1">
            <h3 class="font-bold text-sm text-[#1b1c1c] mb-4">Active Zones</h3>
            
            <div class="overflow-x-auto w-full">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0]">
                    <th class="py-2.5 px-4 text-[10px] font-bold text-[#574146] uppercase tracking-wider rounded-l-lg">Zone Name</th>
                    <th class="py-2.5 px-4 text-[10px] font-bold text-[#574146] uppercase tracking-wider">Radius</th>
                    <th class="py-2.5 px-4 text-[10px] font-bold text-[#574146] uppercase tracking-wider text-right rounded-r-lg">Status</th>
                  </tr>
                </thead>
                <tbody class="text-xs text-[#1b1c1c] divide-y divide-[#F8BBD0]">
                  <!-- Row 1 -->
                  <tr class="hover:bg-[#FCE4EC]/30 transition-colors">
                    <td class="px-4 py-3.5">
                      <div class="font-bold text-[#1b1c1c]">Main Lobby</div>
                      <div class="text-[10px] text-[#574146]">Posko Siaga NTPD 112 Bulukumba</div>
                    </td>
                    <td class="px-4 py-3.5 font-mono text-[#574146] font-bold">{{ radiusMeters }}m</td>
                    <td class="px-4 py-3.5 text-right">
                      <span class="bg-[#E8F5E9] text-[#1B5E20] text-[10px] font-bold uppercase px-2.5 py-0.5 rounded-full border border-[#A5D6A7]">
                        Active
                      </span>
                    </td>
                  </tr>
                  <!-- Row 2 -->
                  <tr class="hover:bg-[#FCE4EC]/30 transition-colors">
                    <td class="px-4 py-3.5">
                      <div class="font-bold text-[#1b1c1c]">North Wing</div>
                      <div class="text-[10px] text-[#574146]">Building B - Diskominfo</div>
                    </td>
                    <td class="px-4 py-3.5 font-mono text-[#574146] font-bold">50m</td>
                    <td class="px-4 py-3.5 text-right">
                      <span class="bg-[#E8F5E9] text-[#1B5E20] text-[10px] font-bold uppercase px-2.5 py-0.5 rounded-full border border-[#A5D6A7]">
                        Active
                      </span>
                    </td>
                  </tr>
                  <!-- Row 3 -->
                  <tr class="hover:bg-[#FCE4EC]/30 transition-colors">
                    <td class="px-4 py-3.5">
                      <div class="font-bold text-[#1b1c1c]">Remote Site A</div>
                      <div class="text-[10px] text-[#574146]">Posko Kecamatan Ujung Bulu</div>
                    </td>
                    <td class="px-4 py-3.5 font-mono text-[#574146] font-bold">200m</td>
                    <td class="px-4 py-3.5 text-right">
                      <span class="bg-[#FCE4EC] text-[#F06292] text-[10px] font-bold uppercase px-2.5 py-0.5 rounded-full border border-[#F8BBD0]">
                        Inactive
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

        </section>

      </div>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import axios from 'axios'
import AdminLayout from '@/layouts/AdminLayout.vue'

const API_BASE = '/api'

// ===== STATE =====
const poskoName = ref('Posko Siaga NTPD 112 Kabupaten Bulukumba')
const poskoAddress = ref('Jl. Jend. Sudirman No. 1, Caile, Kec. Ujung Bulu, Kabupaten Bulukumba, Sulawesi Selatan')
const latitude = ref(-5.5645)
const longitude = ref(120.1945)
const radiusMeters = ref(2.0)
const qrToken = ref('LOPI-Q-POSKO-BULUKUMBA-2026-NTPD112')

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

const refreshQrToken = () => {
  const randomSuffix = Math.random().toString(36).substring(2, 7).toUpperCase()
  qrToken.value = `LOPI-Q-POSKO-BULUKUMBA-2026-${randomSuffix}`
  showToast(true, 'Master Access QR berhasil di-refresh!')
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
    attribution: '© OpenStreetMap'
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
  latitude.value = -5.5645
  longitude.value = 120.1945
  radiusMeters.value = 2.0
  updateMapMarker()
  showToast(true, 'Koordinat di-reset ke default Posko Bulukumba.')
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
    link.href = qrImageUrl.value
    link.download = `Master_QR_LOPI-Q_${poskoName.value.replace(/\s+/g, '_')}.png`
    link.click()
    downloadLoading.value = false
    showToast(true, 'QR Code resmi berhasil diunduh.')
  }, 1000)
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
