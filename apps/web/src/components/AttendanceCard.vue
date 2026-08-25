<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
    <!-- Left Column: Live Geofence Radar & Clock In/Out -->
    <div class="lg:col-span-7 space-y-6">
      <div class="glass-card rounded-3xl p-6 sm:p-8 relative overflow-hidden">
        <!-- Ambient background gradient -->
        <div class="absolute top-0 right-0 w-96 h-96 bg-rose-600/10 rounded-full blur-3xl pointer-events-none"></div>
        
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-bold text-white flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Radar Geofence Posko 112
            </h2>
            <p class="text-xs text-slate-400 mt-0.5">Batas Radius Presensi Maksimal: <span class="text-rose-400 font-bold">2.0 Meter</span></p>
          </div>

          <button 
            @click="refreshLocation" 
            :disabled="locating"
            class="p-2.5 bg-slate-800/80 hover:bg-slate-700/80 text-slate-300 rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all border border-slate-700/60"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" :class="{ 'animate-spin': locating }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>{{ locating ? 'Mencari GPS...' : 'Update Lokasi' }}</span>
          </button>
        </div>

        <!-- Radar Visualizer Box -->
        <div class="relative bg-slate-950/80 border border-slate-800 rounded-2xl p-6 mb-6 flex flex-col items-center justify-center text-center">
          <div class="relative w-48 h-48 flex items-center justify-center my-4">
            <div class="absolute inset-0 rounded-full border border-slate-800"></div>
            <div class="absolute inset-4 rounded-full border border-slate-800/60"></div>
            <div 
              class="absolute inset-8 rounded-full border transition-all duration-500"
              :class="isWithinRadius ? 'border-emerald-500/40 bg-emerald-500/5 animate-ping-slow' : 'border-rose-500/40 bg-rose-500/5'"
            ></div>

            <div class="absolute inset-0 rounded-full overflow-hidden pointer-events-none">
              <div class="w-1/2 h-1/2 bg-gradient-to-tr from-transparent to-rose-500/20 origin-bottom-right animate-radar"></div>
            </div>

            <div class="relative z-10 flex flex-col items-center">
              <div class="w-10 h-10 rounded-full bg-rose-600 border-2 border-white flex items-center justify-center shadow-lg shadow-rose-900/60">
                <span class="text-xs font-black text-white">112</span>
              </div>
              <span class="text-[10px] font-bold text-slate-300 mt-1 bg-slate-900/90 px-2 py-0.5 rounded border border-slate-800">POSKO</span>
            </div>
          </div>

          <div class="mt-2 space-y-2">
            <div class="text-3xl font-black tracking-tight" :class="isWithinRadius ? 'text-emerald-400' : 'text-rose-400'">
              {{ currentDistance !== null ? `${currentDistance} Meter` : 'Menghitung GPS...' }}
            </div>
            
            <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-bold border"
                 :class="isWithinRadius ? 'bg-emerald-500/10 text-emerald-300 border-emerald-500/30 glow-green' : 'bg-rose-500/10 text-rose-300 border-rose-500/30 glow-red'">
              <span class="w-2 h-2 rounded-full" :class="isWithinRadius ? 'bg-emerald-400 animate-pulse' : 'bg-rose-400'"></span>
              <span>{{ isWithinRadius ? 'LOKASI SESUAI (DALAM RADIUS <= 2M)' : 'DI LUAR RADIUS (LEBIH DARI 2M)' }}</span>
            </div>

            <div class="text-xs text-slate-400 font-mono">
              Lat: {{ userLat ? userLat.toFixed(6) : '-' }}, Lng: {{ userLng ? userLng.toFixed(6) : '-' }}
            </div>
          </div>
        </div>

        <!-- Simulation Controls -->
        <div class="bg-slate-900/50 border border-slate-800 rounded-xl p-4 mb-6">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Simulasi Koordinat Posko (Untuk Uji Coba)</span>
            <span class="text-[10px] text-amber-400 font-mono">GPS Test Tool</span>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <button 
              @click="simulateCoords(-5.5645, 120.1945)"
              class="px-2.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-xs text-emerald-400 font-semibold rounded-lg border border-emerald-500/30 transition-all text-center"
            >
              🎯 Tepat di Posko (0m)
            </button>
            <button 
              @click="simulateCoords(-5.56453, 120.19452)"
              class="px-2.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-xs text-emerald-400 font-semibold rounded-lg border border-emerald-500/30 transition-all text-center"
            >
              📍 Dekat Posko (4.2m)
            </button>
            <button 
              @click="simulateCoords(-5.56480, 120.19480)"
              class="px-2.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-xs text-rose-400 font-semibold rounded-lg border border-rose-500/30 transition-all text-center"
            >
              ❌ Luar Posko (45.8m)
            </button>
          </div>
        </div>

        <!-- QR Code Status & Trigger -->
        <div class="bg-slate-900/80 border border-slate-800 rounded-xl p-4 flex items-center justify-between mb-6">
          <div class="flex items-center space-x-3">
            <div class="p-2.5 bg-rose-500/10 rounded-xl border border-rose-500/20 text-rose-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
              </svg>
            </div>
            <div>
              <div class="text-sm font-bold text-white">QR Code Presensi Posko</div>
              <div class="text-xs" :class="qrVerified ? 'text-emerald-400 font-medium' : 'text-slate-400'">
                {{ qrVerified ? '✓ QR Code Posko Terverifikasi' : 'Belum Scan QR Code Posko' }}
              </div>
            </div>
          </div>

          <button 
            @click="$emit('open-qr-scanner')"
            class="px-4 py-2 bg-gradient-to-r from-rose-600 to-amber-600 hover:from-rose-500 hover:to-amber-500 text-white text-xs font-bold rounded-xl shadow-md transition-all flex items-center gap-1.5"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
            </svg>
            <span>Scan QR Code</span>
          </button>
        </div>

        <!-- Action Buttons -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <button 
            @click="handleClockIn"
            :disabled="!canClockIn || submitting"
            class="py-4 px-6 rounded-2xl font-black text-sm tracking-wider uppercase transition-all duration-300 flex items-center justify-center gap-2 shadow-lg"
            :class="canClockIn ? 'bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white shadow-emerald-900/40 glow-green' : 'bg-slate-800/80 text-slate-500 border border-slate-800 cursor-not-allowed'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
            </svg>
            <span>{{ authStore.todayStatus.is_masuk ? 'Sudah Presensi Masuk' : 'Presensi Masuk' }}</span>
          </button>

          <button 
            @click="handleClockOut"
            :disabled="!canClockOut || submitting"
            class="py-4 px-6 rounded-2xl font-black text-sm tracking-wider uppercase transition-all duration-300 flex items-center justify-center gap-2 shadow-lg"
            :class="canClockOut ? 'bg-gradient-to-r from-indigo-600 to-rose-600 hover:from-indigo-500 hover:to-rose-500 text-white shadow-indigo-900/40 glow-red' : 'bg-slate-800/80 text-slate-500 border border-slate-800 cursor-not-allowed'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            <span>{{ authStore.todayStatus.is_pulang ? 'Sudah Presensi Pulang' : 'Presensi Pulang' }}</span>
          </button>
        </div>

        <div v-if="alertMessage" class="mt-4 p-4 rounded-xl text-xs font-semibold border flex items-center justify-between"
             :class="alertSuccess ? 'bg-emerald-950/80 border-emerald-500/40 text-emerald-300' : 'bg-rose-950/80 border-rose-500/40 text-rose-300'">
          <span>{{ alertMessage }}</span>
          <button @click="alertMessage = ''" class="text-slate-400 hover:text-white">&times;</button>
        </div>
      </div>
    </div>

    <!-- Right Column: User Info & Today Status -->
    <div class="lg:col-span-5 space-y-6">
      <div class="glass-card rounded-3xl p-6 relative overflow-hidden">
        <h3 class="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Informasi Pengguna
        </h3>

        <div v-if="authStore.user" class="space-y-4">
          <div class="bg-slate-900/80 rounded-2xl p-4 border border-slate-800">
            <div class="text-xs text-slate-400">Nama Lengkap</div>
            <div class="text-lg font-bold text-white">{{ authStore.user.name }}</div>
            <div class="text-xs font-mono text-rose-400 mt-1">NIP. {{ authStore.user.nip || '-' }}</div>
            <div class="text-xs font-mono text-slate-400 mt-0.5">Email: {{ authStore.user.email || '-' }}</div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="bg-slate-900/60 rounded-xl p-3 border border-slate-800">
              <div class="text-[10px] text-slate-400 uppercase font-semibold">Jabatan</div>
              <div class="text-xs font-bold text-slate-200 mt-1">{{ authStore.user.jabatan || 'Penata Layanan Operasional' }}</div>
            </div>
            <div class="bg-slate-900/60 rounded-xl p-3 border border-slate-800">
              <div class="text-[10px] text-slate-400 uppercase font-semibold">Unit Kerja</div>
              <div class="text-xs font-bold text-slate-200 mt-1">{{ authStore.user.unit_kerja || 'Diskominfo' }}</div>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-8 space-y-3">
          <p class="text-xs text-slate-400">Silakan masuk menggunakan NIP atau Email Dinas Anda untuk presensi.</p>
          <button 
            @click="$emit('open-login')"
            class="px-5 py-2.5 bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs rounded-xl shadow-md transition-all"
          >
            Login Akun
          </button>
        </div>
      </div>

      <!-- Today Status -->
      <div class="glass-card rounded-3xl p-6">
        <h3 class="text-sm font-bold uppercase tracking-wider text-slate-400 mb-4 flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Status Presensi Hari Ini
        </h3>

        <div class="space-y-3">
          <div class="p-4 rounded-2xl border flex items-center justify-between"
               :class="authStore.todayStatus.is_masuk ? 'bg-emerald-950/40 border-emerald-500/30' : 'bg-slate-900/40 border-slate-800'">
            <div>
              <div class="text-xs text-slate-400">Presensi Masuk</div>
              <div class="text-sm font-bold" :class="authStore.todayStatus.is_masuk ? 'text-emerald-400' : 'text-slate-500'">
                {{ authStore.todayStatus.is_masuk ? formatTime(authStore.todayStatus.masuk.timestamp) : 'Belum Melakukan' }}
              </div>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 rounded-lg"
                  :class="authStore.todayStatus.is_masuk ? 'bg-emerald-500/20 text-emerald-300' : 'bg-slate-800 text-slate-500'">
              {{ authStore.todayStatus.is_masuk ? 'Tercatat' : 'Belum' }}
            </span>
          </div>

          <div class="p-4 rounded-2xl border flex items-center justify-between"
               :class="authStore.todayStatus.is_pulang ? 'bg-indigo-950/40 border-indigo-500/30' : 'bg-slate-900/40 border-slate-800'">
            <div>
              <div class="text-xs text-slate-400">Presensi Pulang</div>
              <div class="text-sm font-bold" :class="authStore.todayStatus.is_pulang ? 'text-indigo-400' : 'text-slate-500'">
                {{ authStore.todayStatus.is_pulang ? formatTime(authStore.todayStatus.pulang.timestamp) : 'Belum Melakukan' }}
              </div>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 rounded-lg"
                  :class="authStore.todayStatus.is_pulang ? 'bg-indigo-500/20 text-indigo-300' : 'bg-slate-800 text-slate-500'">
              {{ authStore.todayStatus.is_pulang ? 'Tercatat' : 'Belum' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { calculateDistance, POSKO_COORDS } from '../utils/haversine';

const props = defineProps({
  scannedQr: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['open-qr-scanner', 'open-login']);
const authStore = useAuthStore();

const userLat = ref(-5.5645);
const userLng = ref(120.1945);
const locating = ref(false);
const submitting = ref(false);
const alertMessage = ref('');
const alertSuccess = ref(false);

const currentDistance = computed(() => {
  if (!userLat.value || !userLng.value) return null;
  return calculateDistance(userLat.value, userLng.value);
});

const isWithinRadius = computed(() => {
  return currentDistance.value !== null && currentDistance.value <= POSKO_COORDS.maxRadiusMeters;
});

const qrVerified = computed(() => {
  return props.scannedQr && props.scannedQr.includes('POSKO_112_BULUKUMBA');
});

const canClockIn = computed(() => {
  return authStore.isAuthenticated && isWithinRadius.value && qrVerified.value && !authStore.todayStatus.is_masuk;
});

const canClockOut = computed(() => {
  return authStore.isAuthenticated && isWithinRadius.value && qrVerified.value && authStore.todayStatus.is_masuk && !authStore.todayStatus.is_pulang;
});

function refreshLocation() {
  if (!navigator.geolocation) {
    alertMessage.value = "Geolocation tidak didukung oleh peramban Anda.";
    alertSuccess.value = false;
    return;
  }

  locating.value = true;
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      userLat.value = pos.coords.latitude;
      userLng.value = pos.coords.longitude;
      locating.value = false;
    },
    () => {
      userLat.value = -5.5645;
      userLng.value = 120.1945;
      locating.value = false;
    },
    { enableHighAccuracy: true, timeout: 5000 }
  );
}

function simulateCoords(lat, lng) {
  userLat.value = lat;
  userLng.value = lng;
}

async function handleClockIn() {
  if (!authStore.isAuthenticated) {
    emit('open-login');
    return;
  }

  submitting.value = true;
  alertMessage.value = '';

  try {
    const res = await authStore.clockIn(userLat.value, userLng.value, props.scannedQr);
    alertMessage.value = res.message;
    alertSuccess.value = true;
  } catch (err) {
    alertMessage.value = authStore.error || "Gagal melakukan presensi masuk.";
    alertSuccess.value = false;
  } finally {
    submitting.value = false;
  }
}

async function handleClockOut() {
  if (!authStore.isAuthenticated) {
    emit('open-login');
    return;
  }

  submitting.value = true;
  alertMessage.value = '';

  try {
    const res = await authStore.clockOut(userLat.value, userLng.value, props.scannedQr);
    alertMessage.value = res.message;
    alertSuccess.value = true;
  } catch (err) {
    alertMessage.value = authStore.error || "Gagal melakukan presensi pulang.";
    alertSuccess.value = false;
  } finally {
    submitting.value = false;
  }
}

function formatTime(isoStr) {
  if (!isoStr) return '-';
  const d = new Date(isoStr);
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WITA';
}

onMounted(() => {
  refreshLocation();
});
</script>
