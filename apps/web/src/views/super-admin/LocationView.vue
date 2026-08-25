<template>
  <AdminLayout>
    <div class="flex flex-col gap-6 w-full select-none font-sans">

      <!-- Top Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Manajemen Lokasi Posko &amp; Geofence QR Code</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Konfigurasi koordinat geofence Posko Siaga 112 dan kelola QR Code resmi untuk presensi Call Taker.</p>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span class="px-3.5 py-1.5 bg-emerald-50 text-emerald-800 border border-emerald-200/80 rounded-full text-xs font-extrabold flex items-center gap-1.5 shadow-2xs">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-ping"></span>
            <span>Batas Radius: {{ radiusMeters }} Meter</span>
          </span>
        </div>
      </div>

      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" :class="['flex items-center gap-2.5 p-3.5 rounded-2xl text-xs font-semibold border w-full shadow-sm',
          toast.success ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-red-50 border-red-200 text-red-800'
        ]">
          <span class="material-symbols-outlined text-[18px] shrink-0">
            {{ toast.success ? 'check_circle' : 'error' }}
          </span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Main Grid: Left Config Form + Right Map & QR Code -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">

        <!-- LEFT COLUMN: Form Config (5 cols) -->
        <div class="lg:col-span-5 bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col gap-5">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <h3 class="font-display font-black text-slate-900 text-sm flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">edit_location_alt</span>
              <span>Edit Koordinat &amp; Geofence</span>
            </h3>
            <button
              type="button"
              @click="getCurrentLocation"
              :disabled="gpsLoading"
              class="py-1.5 px-3 bg-rose-50 hover:bg-rose-100 border border-rose-200 text-rose-800 text-[11px] font-bold rounded-xl transition-all flex items-center gap-1 cursor-pointer disabled:opacity-50"
              title="Gunakan lokasi GPS browser saat ini"
            >
              <span class="material-symbols-outlined text-[15px]" :class="{ 'animate-spin': gpsLoading }">
                {{ gpsLoading ? 'sync' : 'my_location' }}
              </span>
              <span>{{ gpsLoading ? 'Mencari...' : 'Lokasi Saya' }}</span>
            </button>
          </div>

          <form @submit.prevent="saveLocation" class="flex flex-col gap-4">

            <!-- Nama Lokasi -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Nama Lokasi Posko <span class="text-rose-500">*</span></label>
              <input
                v-model="poskoName"
                type="text"
                required
                placeholder="Contoh: Posko Siaga NTPD 112 Kabupaten Bulukumba"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs font-semibold"
              />
            </div>

            <!-- Alamat Lengkap -->
            <div class="flex flex-col gap-1">
              <div class="flex items-center justify-between">
                <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Alamat Lengkap Posko</label>
                <span v-if="isReverseGeocoding" class="text-[10px] font-bold text-rose-700 flex items-center gap-1 animate-pulse">
                  <span class="material-symbols-outlined text-[13px] animate-spin">sync</span>
                  <span>Mengisi Alamat Otomatis...</span>
                </span>
              </div>
              <textarea
                v-model="poskoAddress"
                rows="2"
                placeholder="Contoh: Jl. Jend. Sudirman No. 1, Caile, Kec. Ujung Bulu..."
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs font-semibold resize-none"
              ></textarea>
            </div>

            <!-- Latitude & Longitude Inputs -->
            <div class="grid grid-cols-2 gap-3">
              <div class="flex flex-col gap-1">
                <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Latitude <span class="text-rose-500">*</span></label>
                <input
                  v-model.number="latitude"
                  type="number"
                  step="any"
                  required
                  @input="updateMapMarker"
                  class="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none font-mono text-xs font-bold text-rose-700"
                />
              </div>
              <div class="flex flex-col gap-1">
                <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Longitude <span class="text-rose-500">*</span></label>
                <input
                  v-model.number="longitude"
                  type="number"
                  step="any"
                  required
                  @input="updateMapMarker"
                  class="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none font-mono text-xs font-bold text-rose-700"
                />
              </div>
            </div>

            <!-- Batas Radius Presensi (Geofence) -->
            <div class="flex flex-col gap-1.5 bg-slate-50 p-4 border border-slate-200/80 rounded-2xl">
              <div class="flex items-center justify-between">
                <label class="text-[10px] font-extrabold text-slate-600 uppercase tracking-widest">Batas Radius Presensi (Meter) <span class="text-rose-500">*</span></label>
                <span class="font-mono text-xs font-bold text-rose-800 bg-rose-100 px-2 py-0.5 rounded-md">{{ radiusMeters }} Meter</span>
              </div>
              <input
                v-model.number="radiusMeters"
                type="number"
                step="0.5"
                min="0.5"
                max="500"
                required
                @input="updateMapCircle"
                class="w-full px-3.5 py-2 bg-white border border-slate-200 focus:border-rose-500 rounded-xl font-mono text-xs font-bold text-slate-800"
              />
              
              <!-- Quick Radius Preset Buttons -->
              <div class="flex flex-wrap gap-1.5 mt-1">
                <span class="text-[10px] font-bold text-slate-400 self-center mr-1">Preset:</span>
                <button
                  v-for="r in [2, 5, 10, 25, 50]"
                  :key="r"
                  type="button"
                  @click="setRadiusPreset(r)"
                  :class="[
                    radiusMeters === r ? 'bg-rose-700 text-white border-rose-700' : 'bg-white text-slate-700 border-slate-200 hover:bg-slate-100',
                    'px-2.5 py-1 rounded-lg border text-[10px] font-bold transition-all cursor-pointer'
                  ]"
                >
                  {{ r }}m
                </button>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-3 pt-2">
              <button
                type="button"
                @click="resetDefaultLocation"
                class="py-2.5 px-4 bg-white hover:bg-slate-100 text-slate-700 font-bold text-xs rounded-xl border border-slate-200 transition-all cursor-pointer"
              >
                Reset Default
              </button>
              <button
                type="submit"
                :disabled="saveLoading"
                class="flex-1 py-2.5 px-4 bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-1.5 cursor-pointer border-0 disabled:opacity-60"
              >
                <span v-if="saveLoading" class="material-symbols-outlined text-[16px] animate-spin">sync</span>
                <span class="material-symbols-outlined text-[16px]" v-else>save</span>
                <span>{{ saveLoading ? 'Menyimpan...' : 'Simpan Konfigurasi' }}</span>
              </button>
            </div>
          </form>
        </div>

        <!-- RIGHT COLUMN: Interactive Leaflet Map + High-Res QR Generator (7 cols) -->
        <div class="lg:col-span-7 flex flex-col gap-6">

          <!-- Map Box -->
          <div class="bg-white border border-slate-200 rounded-3xl p-5 shadow-sm flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <h3 class="font-display font-black text-slate-900 text-sm flex items-center gap-2">
                <span class="material-symbols-outlined text-rose-700">map</span>
                <span>Peta Interaktif Geofence Posko Siaga</span>
              </h3>
              <span class="text-[10px] font-bold text-slate-400">Klik / Seret Penanda untuk Atur Koordinat</span>
            </div>

            <!-- Leaflet Container -->
            <div class="relative w-full h-[320px] rounded-2xl overflow-hidden border border-slate-200 shadow-inner">
              <div id="leafletMap" class="w-full h-full z-10"></div>
            </div>

            <div class="flex items-center justify-between text-[11px] text-slate-500 bg-slate-50 p-3 rounded-xl border border-slate-200/60">
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full bg-rose-500 border-2 border-white shadow-xs"></span>
                <span>Pusat Koordinat: <strong class="font-mono text-slate-900">{{ latitude.toFixed(6) }}, {{ longitude.toFixed(6) }}</strong></span>
              </div>
              <div class="flex items-center gap-1.5 text-rose-700 font-bold">
                <span class="material-symbols-outlined text-[14px]">radar</span>
                <span>Radius {{ radiusMeters }}m</span>
              </div>
            </div>
          </div>

          <!-- Official QR Code Generator Card -->
          <div class="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col items-center justify-center text-center gap-4">
            <h3 class="font-display font-black text-slate-900 text-sm flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">qr_code_2</span>
              <span>QR Code Resmi Presensi Posko 112</span>
            </h3>

            <!-- QR Code Frame Preview -->
            <div class="p-5 bg-gradient-to-b from-slate-50 to-rose-50/40 rounded-3xl border border-slate-200 shadow-sm flex flex-col items-center gap-2 max-w-sm w-full">
              <div class="bg-white p-4 rounded-2xl border border-slate-200 shadow-md">
                <img :src="qrImageUrl" alt="QR Code Posko 112" class="w-44 h-44 object-contain rounded-xl" />
              </div>
              <div class="text-[11px] font-mono font-extrabold text-slate-800 tracking-wider mt-1">{{ qrToken }}</div>
              <div class="text-[10px] text-slate-500 font-semibold">{{ poskoName }}</div>
            </div>

            <!-- Download Button Only (Print Button Removed) -->
            <div class="w-full max-w-sm">
              <button
                @click="downloadOfficialQrImage"
                :disabled="downloadLoading"
                class="w-full py-3 px-4 bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-extrabold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer border-0 disabled:opacity-60"
              >
                <span class="material-symbols-outlined text-[18px]" :class="{ 'animate-spin': downloadLoading }">
                  {{ downloadLoading ? 'sync' : 'download' }}
                </span>
                <span>{{ downloadLoading ? 'Mengekspor Gambar HD...' : 'Unduh Image QR Resmi (PNG)' }}</span>
              </button>
            </div>
          </div>

        </div>

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

// ===== LOAD LOCATION CONFIG FROM BACKEND =====
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

const isReverseGeocoding = ref(false)

const reverseGeocodeAddress = async (lat: number, lng: number) => {
  isReverseGeocoding.value = true
  try {
    const res = await axios.get(`https://nominatim.openstreetmap.org/reverse`, {
      params: {
        format: 'json',
        lat,
        lon: lng,
        zoom: 18,
        addressdetails: 1
      }
    })
    if (res.data && res.data.display_name) {
      poskoAddress.value = res.data.display_name
    }
  } catch (err) {
    console.warn('Gagal mengambil alamat reverse geocoding:', err)
  } finally {
    isReverseGeocoding.value = false
  }
}

// ===== LEAFLET MAP INITIALIZATION =====
const initLeafletMap = () => {
  if (typeof (window as any).L === 'undefined') {
    // Inject Leaflet CSS & JS dynamically
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

    const container = document.getElementById('leafletMap')
    if (!container) return

    if (mapInstance) {
      mapInstance.remove()
      mapInstance = null
    }

    mapInstance = L.map('leafletMap').setView([latitude.value, longitude.value], 17)

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(mapInstance)

    // Red Posko Marker
    const icon = L.divIcon({
      className: 'custom-leaflet-marker',
      html: `<div style="background:#be123c;width:24px;height:24px;border-radius:50%;border:3px solid white;box-shadow:0 4px 6px -1px rgba(0,0,0,0.3);"></div>`,
      iconSize: [24, 24],
      iconAnchor: [12, 12]
    })

    markerInstance = L.marker([latitude.value, longitude.value], {
      draggable: true,
      icon
    }).addTo(mapInstance)

    // Geofence Circle
    circleInstance = L.circle([latitude.value, longitude.value], {
      color: '#be123c',
      fillColor: '#f43f5e',
      fillOpacity: 0.25,
      radius: radiusMeters.value
    }).addTo(mapInstance)

    // On Marker Drag
    markerInstance.on('dragend', (event: any) => {
      const position = event.target.getLatLng()
      latitude.value = parseFloat(position.lat.toFixed(7))
      longitude.value = parseFloat(position.lng.toFixed(7))
      circleInstance.setLatLng(position)
      reverseGeocodeAddress(latitude.value, longitude.value)
    })

    // On Map Click
    mapInstance.on('click', (e: any) => {
      const { lat, lng } = e.latlng
      latitude.value = parseFloat(lat.toFixed(7))
      longitude.value = parseFloat(lng.toFixed(7))
      markerInstance.setLatLng([lat, lng])
      circleInstance.setLatLng([lat, lng])
      reverseGeocodeAddress(latitude.value, longitude.value)
    })
  })
}

const updateMapMarker = () => {
  if (mapInstance && markerInstance && circleInstance && (window as any).L) {
    const L = (window as any).L
    const pos = new L.LatLng(latitude.value, longitude.value)
    markerInstance.setLatLng(pos)
    circleInstance.setLatLng(pos)
    mapInstance.panTo(pos)
  }
}

const updateMapCircle = () => {
  if (circleInstance) {
    circleInstance.setRadius(radiusMeters.value)
  }
}

const setRadiusPreset = (r: number) => {
  radiusMeters.value = r
  updateMapCircle()
}

// ===== GET CURRENT GPS LOCATION =====
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    showToast(false, 'Browser Anda tidak mendukung Geolocation GPS.')
    return
  }
  gpsLoading.value = true
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      latitude.value = parseFloat(pos.coords.latitude.toFixed(7))
      longitude.value = parseFloat(pos.coords.longitude.toFixed(7))
      updateMapMarker()
      reverseGeocodeAddress(latitude.value, longitude.value)
      gpsLoading.value = false
      showToast(true, 'Lokasi GPS berhasil didapatkan secara otomatis!')
    },
    (err) => {
      gpsLoading.value = false
      showToast(false, 'Gagal mendapatkan lokasi GPS: ' + err.message)
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

const resetDefaultLocation = () => {
  poskoName.value = 'Posko Siaga NTPD 112 Kabupaten Bulukumba'
  poskoAddress.value = 'Jl. Jend. Sudirman No. 1, Caile, Kec. Ujung Bulu, Kabupaten Bulukumba, Sulawesi Selatan'
  latitude.value = -5.5645
  longitude.value = 120.1945
  radiusMeters.value = 2.0
  updateMapMarker()
  updateMapCircle()
}

// ===== SAVE LOCATION =====
const saveLocation = async () => {
  saveLoading.value = true
  try {
    const res = await axios.put(`${API_BASE}/admin/location`, {
      name: poskoName.value,
      address: poskoAddress.value,
      latitude: latitude.value,
      longitude: longitude.value,
      radius_meters: radiusMeters.value
    })
    if (res.data && res.data.qr_token) {
      qrToken.value = res.data.qr_token
    }
    showToast(true, 'Konfigurasi lokasi, geofence & QR Code baru berhasil disimpan!')
  } catch (err: any) {
    showToast(false, err.response?.data?.error || 'Gagal menyimpan lokasi.')
  } finally {
    saveLoading.value = false
  }
}

// Helper to draw centered wrapped text line by line on canvas
const drawWrappedText = (
  ctx: CanvasRenderingContext2D,
  text: string,
  x: number,
  startY: number,
  maxWidth: number,
  lineHeight: number
): number => {
  const words = text.split(' ')
  let line = ''
  const lines: string[] = []

  for (let n = 0; n < words.length; n++) {
    const testLine = line + (line ? ' ' : '') + words[n]
    const metrics = ctx.measureText(testLine)
    if (metrics.width > maxWidth && n > 0) {
      lines.push(line)
      line = words[n]
    } else {
      line = testLine
    }
  }
  lines.push(line)

  lines.forEach((l, i) => {
    ctx.fillText(l, x, startY + i * lineHeight)
  })

  return lines.length
}

// ===== HIGH-RES OFFICIAL QR IMAGE CANVAS GENERATOR (MULTILINE RAPI) =====
const downloadOfficialQrImage = () => {
  downloadLoading.value = true

  const canvas = document.createElement('canvas')
  canvas.width = 900
  canvas.height = 1200
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    downloadLoading.value = false
    return
  }

  // 1. Background White
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, 900, 1200)

  // 2. Top Header Gradient
  const grad = ctx.createLinearGradient(0, 0, 900, 0)
  grad.addColorStop(0, '#be123c') // rose-700
  grad.addColorStop(0.5, '#9f1239') // rose-800
  grad.addColorStop(1, '#d97706') // amber-600
  ctx.fillStyle = grad
  ctx.fillRect(0, 0, 900, 150)

  // Top Header Text
  ctx.fillStyle = '#ffffff'
  ctx.font = 'bold 24px Montserrat, sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('PEMERINTAH KABUPATEN BULUKUMBA', 450, 52)

  ctx.font = '18px Montserrat, sans-serif'
  ctx.fillText('DINAS KOMUNIKASI INFORMATIKA DAN PERSANDIAN', 450, 85)

  ctx.fillStyle = '#fef08a' // yellow-200
  ctx.font = 'bold 16px sans-serif'
  ctx.fillText('NTPD 112 - POSKO PRESENSI SIAGA CALL TAKER', 450, 122)

  // Outer Border Frame Box
  ctx.strokeStyle = '#cbd5e1'
  ctx.lineWidth = 3
  ctx.strokeRect(35, 175, 830, 980)

  // Inner Title Location (Wrapped cleanly if long)
  ctx.fillStyle = '#0f172a'
  ctx.font = 'bold 24px Montserrat, sans-serif'
  ctx.textAlign = 'center'
  const titleLines = drawWrappedText(ctx, poskoName.value.toUpperCase(), 450, 225, 760, 32)

  // Address (Wrapped cleanly with max width 760px)
  const addressStartY = 225 + titleLines * 32 + 10
  ctx.fillStyle = '#475569'
  ctx.font = '15px sans-serif'
  const addressLines = drawWrappedText(ctx, poskoAddress.value, 450, addressStartY, 760, 22)

  // QR Code Border Box Position
  const qrY = addressStartY + addressLines * 22 + 20
  ctx.fillStyle = '#f8fafc'
  ctx.fillRect(210, qrY, 480, 480)
  ctx.strokeStyle = '#cbd5e1'
  ctx.lineWidth = 2
  ctx.strokeRect(210, qrY, 480, 480)

  // Load Image QR Code into Canvas
  const qrImg = new Image()
  qrImg.crossOrigin = 'Anonymous'
  qrImg.onload = () => {
    ctx.drawImage(qrImg, 240, qrY + 30, 420, 420)

    // Token Code
    ctx.fillStyle = '#be123c'
    ctx.font = 'bold 20px monospace'
    ctx.fillText(qrToken.value, 450, qrY + 480 + 35)

    // Metadata Coordinates & Radius
    ctx.fillStyle = '#0f172a'
    ctx.font = 'bold 17px monospace'
    ctx.fillText(`KOORDINAT: ${latitude.value.toFixed(7)}, ${longitude.value.toFixed(7)}`, 450, qrY + 480 + 75)

    ctx.fillStyle = '#15803d' // emerald-700
    ctx.font = 'bold 16px sans-serif'
    ctx.fillText(`BATAS RADIUS GEOFENCE: ${radiusMeters.value} METER`, 450, qrY + 480 + 110)

    // Usage Instruction
    ctx.fillStyle = '#475569'
    ctx.font = 'italic 15px sans-serif'
    ctx.fillText('Scan QR Code ini menggunakan aplikasi LOPI-Q untuk melakukan presensi siaga.', 450, qrY + 480 + 155)

    // Timestamp & Footer
    ctx.fillStyle = '#94a3b8'
    ctx.font = '13px monospace'
    ctx.fillText(`Official Generated by LOPI-Q • ${new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}`, 450, qrY + 480 + 205)

    // Trigger Download
    const link = document.createElement('a')
    link.download = `QR-Presensi-LOPI-Q-${poskoName.value.replace(/[^a-zA-Z0-9]/g, '-')}.png`
    link.href = canvas.toDataURL('image/png')
    link.click()

    downloadLoading.value = false
    showToast(true, 'Image QR Code Resmi berhasil diunduh dalam kualitas HD!')
  }
  qrImg.onerror = () => {
    downloadLoading.value = false
    showToast(false, 'Gagal memuat QR Code image.')
  }
  qrImg.src = qrImageUrl.value
}

onMounted(async () => {
  await fetchLocationConfig()
  initLeafletMap()
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
