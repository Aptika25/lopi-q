<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    <div class="glass-card w-full max-w-lg rounded-3xl p-6 relative border border-slate-700 shadow-2xl overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-800">
        <div class="flex items-center space-x-3">
          <div class="p-2.5 bg-rose-500/10 rounded-xl text-rose-500">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">Scan QR Code Posko 112</h3>
            <p class="text-xs text-slate-400">Pindai QR Code resmi Posko NTPD 112 Bulukumba</p>
          </div>
        </div>
        
        <button @click="$emit('close')" class="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Mode Selector Tabs -->
      <div class="flex items-center bg-slate-900/80 rounded-xl p-1 mt-4 mb-4 border border-slate-800 text-xs font-semibold">
        <button 
          @click="mode = 'camera'"
          class="flex-1 py-2 rounded-lg transition-all"
          :class="mode === 'camera' ? 'bg-rose-600 text-white shadow' : 'text-slate-400 hover:text-white'"
        >
          📷 Kamera Scanner
        </button>
        <button 
          @click="mode = 'posko_qr'; fetchOfficialQr();"
          class="flex-1 py-2 rounded-lg transition-all"
          :class="mode === 'posko_qr' ? 'bg-rose-600 text-white shadow' : 'text-slate-400 hover:text-white'"
        >
          🖼️ QR Code Posko (Server)
        </button>
      </div>

      <!-- Camera Scanner Mode -->
      <div v-show="mode === 'camera'" class="space-y-4">
        <div class="relative bg-black rounded-2xl overflow-hidden min-h-[260px] flex items-center justify-center border border-slate-800">
          <div id="reader-apps-web" class="w-full"></div>
          <div v-if="scannerError" class="p-4 text-center text-xs text-rose-400">
            {{ scannerError }}
          </div>
        </div>
        <p class="text-[11px] text-slate-400 text-center">Arahkan kamera HP ke QR Code yang terpajang di meja Posko Peserta Magang 112.</p>
      </div>

      <!-- Server Posko QR Mode -->
      <div v-show="mode === 'posko_qr'" class="space-y-4 text-center">
        <div class="bg-white p-4 rounded-2xl inline-block shadow-xl border-4 border-slate-800">
          <img v-if="officialQrImage" :src="officialQrImage" alt="Official Posko 112 QR" class="w-56 h-56 mx-auto" />
          <div v-else class="w-56 h-56 flex items-center justify-center text-slate-500 text-xs font-mono">
            Loading QR Code...
          </div>
        </div>
        <div>
          <div class="text-xs font-bold text-slate-200">QR Code Resmi Posko 112 Bulukumba</div>
          <div class="text-[11px] text-slate-400 mt-0.5">Berlaku otomatis untuk verifikasi presensi call taker di lokasi posko.</div>
        </div>

        <button 
          @click="selectToken(officialQrToken)"
          :disabled="!officialQrToken"
          class="w-full py-3 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white text-xs font-bold rounded-xl shadow-lg transition-all"
        >
          ✓ Gunakan QR Code Posko Ini
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onUnmounted } from 'vue';
import axios from 'axios';
import { Html5QrcodeScanner } from 'html5-qrcode';

const props = defineProps({
  isOpen: Boolean
});

const emit = defineEmits(['close', 'scan-success']);

const mode = ref('posko_qr');
const officialQrImage = ref('');
const officialQrToken = ref('');
const scannerError = ref('');
let html5QrcodeScanner = null;

async function fetchOfficialQr() {
  try {
    const res = await axios.get('/api/presensi/posko-qr');
    if (res.data.success) {
      officialQrImage.value = res.data.qr_image;
      officialQrToken.value = res.data.qr_token;
    }
  } catch (err) {
    console.error('Failed to fetch official posko QR:', err);
  }
}

function startScanner() {
  if (html5QrcodeScanner) return;

  try {
    html5QrcodeScanner = new Html5QrcodeScanner("reader-apps-web", {
      fps: 10,
      qrbox: { width: 220, height: 220 }
    }, false);

    html5QrcodeScanner.render((decodedText) => {
      selectToken(decodedText);
      stopScanner();
    }, () => {});
  } catch (err) {
    scannerError.value = "Tidak dapat mengakses kamera. Pastikan izin kamera telah diberikan.";
  }
}

function stopScanner() {
  if (html5QrcodeScanner) {
    html5QrcodeScanner.clear().catch(err => console.error(err));
    html5QrcodeScanner = null;
  }
}

function selectToken(token) {
  emit('scan-success', token);
  emit('close');
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    fetchOfficialQr();
    if (mode.value === 'camera') {
      setTimeout(startScanner, 300);
    }
  } else {
    stopScanner();
  }
});

watch(mode, (newMode) => {
  if (newMode === 'camera') {
    setTimeout(startScanner, 300);
  } else {
    stopScanner();
  }
});

onUnmounted(() => {
  stopScanner();
});
</script>
