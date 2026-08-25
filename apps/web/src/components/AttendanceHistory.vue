<template>
  <div class="glass-card rounded-3xl p-6 sm:p-8 space-y-6">
    <div class="flex items-center justify-between pb-4 border-b border-slate-800">
      <div>
        <h2 class="text-xl font-bold text-white flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Riwayat Presensi Peserta Magang 112
        </h2>
        <p class="text-xs text-slate-400 mt-1">Log presensi masuk & pulang beserta lokasi GPS & verifikasi radius 10 meter.</p>
      </div>

      <button 
        @click="authStore.fetchHistory" 
        class="p-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-semibold flex items-center gap-1.5 transition-all border border-slate-700/60"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        <span>Refresh Log</span>
      </button>
    </div>

    <div class="overflow-x-auto rounded-2xl border border-slate-800">
      <table class="w-full text-left text-xs text-slate-300">
        <thead class="bg-slate-900/90 text-slate-400 font-semibold uppercase tracking-wider border-b border-slate-800">
          <tr>
            <th class="py-3.5 px-4">Waktu (WITA)</th>
            <th class="py-3.5 px-4">Peserta Magang / NIP</th>
            <th class="py-3.5 px-4">Tipe Presensi</th>
            <th class="py-3.5 px-4">Jarak Geofence</th>
            <th class="py-3.5 px-4">Status QR Code</th>
            <th class="py-3.5 px-4 text-right">Koordinat GPS</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800/60 bg-slate-950/40">
          <tr v-for="h in authStore.presensiHistory" :key="h.id" class="hover:bg-slate-900/60 transition-colors">
            <td class="py-3.5 px-4 font-mono text-slate-200">
              {{ formatDate(h.timestamp) }}
            </td>
            <td class="py-3.5 px-4">
              <div class="font-bold text-white">{{ h.user_name }}</div>
              <div class="font-mono text-[11px] text-rose-400">NIP. {{ h.user_nip || '-' }}</div>
            </td>
            <td class="py-3.5 px-4">
              <span 
                class="px-2.5 py-1 rounded-lg text-[10px] font-extrabold uppercase tracking-wider"
                :class="h.type === 'MASUK' ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30' : 'bg-indigo-500/20 text-indigo-300 border border-indigo-500/30'"
              >
                {{ h.type }}
              </span>
            </td>
            <td class="py-3.5 px-4">
              <span class="font-bold font-mono text-emerald-400">
                {{ h.distance_meters }}m
              </span>
              <span class="text-[10px] text-slate-400 ml-1">(<= 10.0m)</span>
            </td>
            <td class="py-3.5 px-4">
              <span class="text-xs text-emerald-400 font-semibold flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                QR Valid
              </span>
            </td>
            <td class="py-3.5 px-4 text-right font-mono text-[11px] text-slate-400">
              {{ h.latitude ? h.latitude.toFixed(5) : '-' }}, {{ h.longitude ? h.longitude.toFixed(5) : '-' }}
            </td>
          </tr>

          <tr v-if="authStore.presensiHistory.length === 0">
            <td colspan="6" class="py-8 text-center text-slate-500 text-xs font-medium">
              Belum ada riwayat presensi hari ini.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();

function formatDate(isoStr) {
  if (!isoStr) return '-';
  const d = new Date(isoStr);
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) + ' ' + 
         d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WITA';
}

onMounted(() => {
  authStore.fetchHistory();
});
</script>
