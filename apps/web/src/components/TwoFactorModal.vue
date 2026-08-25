<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    <div class="glass-card w-full max-w-md rounded-3xl p-6 relative border border-slate-700 shadow-2xl overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-800">
        <div class="flex items-center space-x-3">
          <div class="p-2.5 bg-emerald-500/10 rounded-xl text-emerald-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">Google Authenticator (2FA)</h3>
            <p class="text-xs text-slate-400">Keamanan Otentikasi Dua Langkah</p>
          </div>
        </div>

        <button @click="$emit('close')" class="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- State: 2FA Already Active -->
      <div v-if="authStore.user?.totp_enabled" class="py-6 text-center space-y-4">
        <div class="w-16 h-16 bg-emerald-500/20 text-emerald-400 rounded-full flex items-center justify-center mx-auto border border-emerald-500/30">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <div>
          <h4 class="text-base font-bold text-white">Google Authenticator Aktif</h4>
          <p class="text-xs text-slate-400 mt-1">Akun Anda saat ini dilindungi oleh otentikasi dua langkah (2FA).</p>
        </div>

        <button 
          @click="handleDisable"
          :disabled="loading"
          class="w-full py-3 bg-rose-900/60 hover:bg-rose-900 text-rose-300 font-bold text-xs rounded-xl border border-rose-700/50 transition-all"
        >
          Nonaktifkan 2FA
        </button>
      </div>

      <!-- State: Setup 2FA -->
      <div v-else class="py-4 space-y-4">
        <div v-if="!qrCodeDataUrl" class="text-center py-6">
          <button 
            @click="initSetup"
            :disabled="loading"
            class="px-6 py-3 bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs rounded-xl shadow-lg transition-all"
          >
            Mulai Aktifkan Google Authenticator
          </button>
        </div>

        <div v-else class="space-y-4">
          <div class="text-center">
            <div class="text-xs text-slate-400 mb-2">1. Pindai QR Code menggunakan aplikasi Google Authenticator di HP Anda:</div>
            <div class="bg-white p-3 rounded-2xl inline-block shadow-lg border-2 border-slate-700">
              <img :src="qrCodeDataUrl" alt="2FA QR Code" class="w-48 h-48 mx-auto" />
            </div>
            <div class="text-[11px] font-mono text-slate-400 mt-2 bg-slate-900 p-2 rounded border border-slate-800">
              Secret: <span class="text-amber-400 font-bold">{{ secretKey }}</span>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">2. Masukkan 6 Digit Kode OTP dari Aplikasi:</label>
            <input 
              v-model="totpCode"
              type="text"
              maxlength="6"
              placeholder="Contoh: 123456"
              class="w-full px-4 py-3 bg-slate-900 border border-slate-700 rounded-xl text-center text-xl font-mono tracking-widest text-emerald-400 focus:outline-none focus:border-emerald-500"
            />
          </div>

          <div v-if="errorMsg" class="p-3 bg-rose-950/80 border border-rose-500/40 text-rose-300 text-xs rounded-xl text-center">
            {{ errorMsg }}
          </div>

          <button 
            @click="handleEnable"
            :disabled="loading || totpCode.length !== 6"
            class="w-full py-3 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-bold text-xs rounded-xl shadow-lg transition-all"
          >
            Verifikasi & Aktifkan 2FA
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useAuthStore } from '../stores/auth';

const props = defineProps({
  isOpen: Boolean
});

const emit = defineEmits(['close']);
const authStore = useAuthStore();

const loading = ref(false);
const qrCodeDataUrl = ref('');
const secretKey = ref('');
const totpCode = ref('');
const errorMsg = ref('');

async function initSetup() {
  loading.value = true;
  errorMsg.value = '';
  try {
    const data = await authStore.setup2fa();
    if (data && data.success) {
      qrCodeDataUrl.value = data.qr_code;
      secretKey.value = data.secret;
    }
  } catch (err) {
    errorMsg.value = "Gagal memulai setup 2FA.";
  } finally {
    loading.value = false;
  }
}

async function handleEnable() {
  if (totpCode.value.length !== 6) return;
  loading.value = true;
  errorMsg.value = '';

  try {
    await authStore.enable2fa(totpCode.value);
    emit('close');
  } catch (err) {
    errorMsg.value = authStore.error || "Kode verifikasi salah.";
  } finally {
    loading.value = false;
  }
}

async function handleDisable() {
  if (!confirm("Apakah Anda yakin ingin menonaktifkan 2FA Google Authenticator?")) return;
  loading.value = true;
  try {
    await authStore.disable2fa();
    emit('close');
  } catch (err) {
    errorMsg.value = authStore.error || "Gagal menonaktifkan 2FA.";
  } finally {
    loading.value = false;
  }
}

watch(() => props.isOpen, (newVal) => {
  if (newVal && !authStore.user?.totp_enabled) {
    initSetup();
  } else {
    qrCodeDataUrl.value = '';
    totpCode.value = '';
    errorMsg.value = '';
  }
});
</script>
