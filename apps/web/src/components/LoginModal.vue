<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    <div class="glass-card w-full max-w-md rounded-3xl p-6 relative border border-slate-700 shadow-2xl overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-800">
        <div class="flex items-center space-x-3">
          <div class="p-2.5 bg-gradient-to-tr from-rose-600 to-amber-500 rounded-xl text-white font-extrabold text-lg shadow-md">
            112
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">Login LOPI-Q</h3>
            <p class="text-xs text-slate-400">Autentikasi NIP Pegawai / Email Dinas</p>
          </div>
        </div>

        <button @click="$emit('close')" class="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form Step 1 -->
      <form v-if="!authStore.otpRequired" @submit.prevent="handleLogin" class="mt-4 space-y-4 text-xs">
        <div>
          <label class="block text-slate-300 font-semibold mb-1">NIP atau Email Dinas</label>
          <div class="relative">
            <input 
              v-model="identifier" 
              type="text" 
              required
              placeholder="Masukkan NIP atau Email Dinas..." 
              class="w-full px-4 py-3 bg-slate-900 border border-slate-700 rounded-xl text-white font-mono focus:outline-none focus:border-rose-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-slate-300 font-semibold mb-1">Kata Sandi (Password)</label>
          <input 
            v-model="password" 
            type="password" 
            required
            placeholder="Masukkan kata sandi..." 
            class="w-full px-4 py-3 bg-slate-900 border border-slate-700 rounded-xl text-white focus:outline-none focus:border-rose-500"
          />
        </div>

        <!-- Quick Selector -->
        <div class="bg-slate-900/60 p-3 rounded-xl border border-slate-800">
          <div class="text-[10px] text-slate-400 uppercase font-bold mb-2">Pilih Akun Demo:</div>
          <div class="grid grid-cols-2 gap-1.5 max-h-36 overflow-y-auto pr-1">
            <button 
              type="button"
              @click="prefill('aswan@bulukumbakab.go.id', 'Asw&a198')"
              class="px-2 py-1.5 bg-rose-950/80 hover:bg-rose-900 text-rose-300 text-[10px] font-bold rounded text-left border border-rose-800/50 truncate"
            >
              👑 Super Admin (Aswan)
            </button>
            <button 
              type="button"
              @click="prefill('19940503202521138', '123456')"
              class="px-2 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] font-medium rounded text-left border border-slate-700 truncate"
            >
              👤 A.Mappalua (NIP)
            </button>
          </div>
        </div>

        <div v-if="authStore.error" class="p-3 bg-rose-950/80 border border-rose-500/40 text-rose-300 text-xs rounded-xl">
          {{ authStore.error }}
        </div>

        <button 
          type="submit" 
          :disabled="authStore.loading"
          class="w-full py-3 bg-gradient-to-r from-rose-600 to-amber-600 hover:from-rose-500 hover:to-amber-500 text-white font-bold rounded-xl shadow-lg transition-all"
        >
          {{ authStore.loading ? 'Memproses Login...' : 'Masuk Sekarang' }}
        </button>
      </form>

      <!-- Form Step 2: 2FA -->
      <form v-else @submit.prevent="handleVerify2fa" class="mt-4 space-y-4 text-xs">
        <div class="text-center py-2">
          <div class="w-12 h-12 bg-amber-500/20 text-amber-400 rounded-full flex items-center justify-center mx-auto mb-2 border border-amber-500/30">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h4 class="text-sm font-bold text-white">Masukkan Kode Google Authenticator</h4>
          <p class="text-xs text-slate-400 mt-0.5">Buka aplikasi Google Authenticator & masukkan 6-digit kode OTP.</p>
        </div>

        <div>
          <input 
            v-model="totpCode" 
            type="text" 
            maxlength="6"
            required
            placeholder="Contoh: 123456" 
            class="w-full px-4 py-3 bg-slate-900 border border-slate-700 rounded-xl text-center text-2xl font-mono tracking-widest text-emerald-400 focus:outline-none focus:border-emerald-500"
          />
        </div>

        <div v-if="authStore.error" class="p-3 bg-rose-950/80 border border-rose-500/40 text-rose-300 text-xs rounded-xl text-center">
          {{ authStore.error }}
        </div>

        <button 
          type="submit" 
          :disabled="authStore.loading || totpCode.length !== 6"
          class="w-full py-3 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-xl shadow-lg transition-all"
        >
          {{ authStore.loading ? 'Verifikasi Kode...' : 'Verifikasi & Masuk' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';

const props = defineProps({
  isOpen: Boolean
});

const emit = defineEmits(['close']);
const authStore = useAuthStore();

const identifier = ref('');
const password = ref('');
const totpCode = ref('');

function prefill(id, p) {
  identifier.value = id;
  password.value = p;
}

async function handleLogin() {
  const res = await authStore.login(identifier.value, password.value);
  if (res.success) {
    emit('close');
  }
}

async function handleVerify2fa() {
  const res = await authStore.verify2fa(totpCode.value);
  if (res.success) {
    emit('close');
  }
}
</script>
