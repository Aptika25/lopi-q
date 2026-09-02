<template>
  <InternLayout>
    <div class="space-y-6 select-none font-sans">
      
      <!-- ===== 1. PAGE HEADER ===== -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-[#F8BBD0]/60 pb-4">
        <div>
          <h2 class="font-display font-black text-[#1b1c1c] text-lg sm:text-xl flex items-center gap-2">
            <span class="material-symbols-outlined text-[#f06292] text-2xl">dashboard</span>
            <span>Dashboard Siaga Peserta Magang</span>
          </h2>
          <p class="font-sans text-[#8a7176] text-xs mt-1">Panel utama presensi siaga 112, status verifikasi 2FA, dan riwayat tugas harian.</p>
        </div>
        <div class="flex items-center gap-2">
          <span class="px-3 py-1 bg-[#FCE4EC] text-[#ab2c5d] text-xs font-extrabold rounded-full border border-[#F8BBD0] flex items-center gap-1.5 shadow-2xs">
            <span class="w-2 h-2 rounded-full bg-[#f06292] animate-ping"></span>
            <span>SIAGA 112 ONLINE</span>
          </span>
        </div>
      </div>

      <!-- ===== 2. HERO PROFILE CARD ===== -->
      <div class="bg-gradient-to-br from-white via-[#FFF5F8] to-[#FCE4EC]/50 rounded-3xl p-6 sm:p-7 border border-[#F8BBD0] shadow-sm space-y-6 relative overflow-hidden">
        <!-- Decorative Ambient Glow -->
        <div class="absolute -right-10 -bottom-10 w-48 h-48 bg-[#fec1d6]/30 rounded-full blur-3xl pointer-events-none"></div>

        <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6 relative z-10">
          
          <!-- Avatar & Detailed Profile Metadata -->
          <div class="flex items-start gap-4 sm:gap-5 min-w-0">
            <!-- Avatar Circle with Initial -->
            <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-2xl bg-gradient-to-tr from-[#ab2c5d] via-[#f06292] to-[#fec1d6] text-white font-display font-black text-2xl sm:text-3xl flex items-center justify-center shrink-0 shadow-md border-2 border-white">
              {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'P' }}
            </div>

            <div class="space-y-1.5 flex-1 min-w-0">
              <!-- Category Pill -->
              <div class="inline-flex items-center gap-1.5 px-3 py-0.5 bg-[#FCE4EC] text-[#ab2c5d] rounded-full text-[10px] sm:text-xs font-extrabold border border-[#F8BBD0]">
                <span class="w-1.5 h-1.5 rounded-full bg-[#f06292] animate-pulse"></span>
                <span>PESERTA MAGANG POSKO SIAGA 112</span>
              </div>

              <!-- Name & Greeting -->
              <h1 class="text-lg sm:text-2xl font-display font-black text-[#1b1c1c] leading-tight truncate">
                Selamat Tugas, {{ authStore.user?.name || 'Peserta Magang' }}!
              </h1>

              <!-- Identitas NISN / NIM -->
              <div class="text-xs font-mono font-extrabold text-[#ab2c5d] flex items-center gap-1">
                <span class="material-symbols-outlined text-[16px]">badge</span>
                <span>NISN / NIM: {{ authStore.user?.nip || '-' }}</span>
              </div>

              <!-- Grid Detail: Jurusan & Asal Sekolah -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-1 text-xs">
                <div class="flex items-center gap-1.5 text-[#574146]">
                  <span class="material-symbols-outlined text-[16px] text-[#f06292]">school</span>
                  <span class="font-bold text-[#1b1c1c]">Jurusan:</span>
                  <span class="truncate">{{ authStore.user?.unit_kerja || 'Rekayasa Perangkat Lunak' }}</span>
                </div>
                <div class="flex items-center gap-1.5 text-[#574146]">
                  <span class="material-symbols-outlined text-[16px] text-[#f06292]">location_city</span>
                  <span class="font-bold text-[#1b1c1c]">Asal Sekolah / Univ:</span>
                  <span class="truncate">{{ authStore.user?.jabatan || 'SMK Negeri 1 Bulukumba' }}</span>
                </div>
              </div>

              <!-- Status Badges: 2FA & Akun -->
              <div class="flex flex-wrap items-center gap-2 pt-2">
                <!-- Status 2FA Badge -->
                <span 
                  v-if="authStore.user?.totp_enabled" 
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-emerald-50 text-emerald-700 border border-emerald-200"
                >
                  <span class="material-symbols-outlined text-[14px] text-emerald-600">verified_user</span>
                  <span>2FA Terverifikasi</span>
                </span>
                <span 
                  v-else 
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-slate-100 text-slate-500 border border-slate-200"
                >
                  <span class="material-symbols-outlined text-[14px] text-slate-400">gpp_maybe</span>
                  <span>Belum Setup 2FA</span>
                </span>

                <!-- Status Akun Badge -->
                <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[10px] font-extrabold bg-emerald-50 text-emerald-800 border border-emerald-200">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping"></span>
                  <span>Akun Aktif (Siap Login)</span>
                </span>
              </div>
            </div>
          </div>

          <!-- Action Buttons Bar -->
          <div class="flex flex-col sm:flex-row lg:flex-col gap-2.5 shrink-0 pt-4 lg:pt-0 border-t lg:border-0 border-[#F8BBD0]/60">
            <router-link 
              to="/intern/scan" 
              class="px-5 py-3 bg-gradient-to-r from-[#f06292] via-[#ab2c5d] to-[#881b47] hover:from-[#ab2c5d] hover:to-[#701339] text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-2 decoration-none cursor-pointer active:scale-98"
            >
              <span class="material-symbols-outlined text-lg">qr_code_scanner</span>
              <span>Buka Kamera Scan QR</span>
            </router-link>

            <button 
              @click="showLeaveModal = true"
              class="px-4 py-3 bg-white hover:bg-[#FCE4EC] text-[#574146] border border-[#F8BBD0] font-bold text-xs rounded-xl shadow-2xs transition-all flex items-center justify-center gap-2 cursor-pointer active:scale-98"
            >
              <span class="material-symbols-outlined text-lg text-[#ab2c5d]">edit_document</span>
              <span>Ajukan Sakit / Izin / Shift</span>
            </button>
          </div>

        </div>
      </div>

      <!-- ===== 3. STATUS OPERASIONAL SHIFT BANNER ===== -->
      <div 
        class="rounded-3xl p-5 sm:p-6 border shadow-md transition-all flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 relative overflow-hidden"
        :class="{
          'bg-gradient-to-r from-emerald-950 via-teal-900 to-slate-900 border-emerald-600/50 text-white': shiftInfo.badge === 'ACTIVE',
          'bg-gradient-to-r from-slate-900 via-rose-950 to-slate-900 border-rose-600/50 text-white': shiftInfo.badge === 'OFF',
          'bg-gradient-to-r from-amber-950 via-slate-900 to-amber-900 border-amber-600/50 text-white': shiftInfo.badge === 'LEAVE',
          'bg-gradient-to-r from-slate-800 via-slate-900 to-slate-800 border-slate-600 text-white animate-pulse': shiftInfo.badge === 'LOADING'
        }"
      >
        <div class="flex items-center gap-4 relative z-10">
          <div 
            class="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0 border shadow-inner transition-all"
            :class="shiftInfo.badge === 'ACTIVE' ? 'bg-emerald-500/20 border-emerald-400/30 text-emerald-300' : 'bg-white/10 border-white/20 text-amber-300'"
          >
            <span class="material-symbols-outlined text-2xl">
              {{ shiftInfo.badge === 'ACTIVE' ? 'verified_user' : 'event_upcoming' }}
            </span>
          </div>
          <div class="space-y-1">
            <div 
              class="text-[10px] font-mono font-black uppercase tracking-widest flex items-center gap-1.5"
              :class="shiftInfo.badge === 'ACTIVE' ? 'text-emerald-400' : 'text-amber-300'"
            >
              <span class="w-2 h-2 rounded-full" :class="shiftInfo.badge === 'ACTIVE' ? 'bg-emerald-400 animate-ping' : 'bg-amber-400'"></span>
              <span>STATUS OPERASIONAL SHIFT MAGANG HARI INI</span>
            </div>
            <h3 class="font-display font-black text-base sm:text-lg leading-snug">
              {{ shiftInfo.shiftName }}
            </h3>
            <p class="text-xs font-mono text-slate-300">
              {{ shiftInfo.shiftTimeStr }}
            </p>
          </div>
        </div>

        <div class="shrink-0 relative z-10">
          <span 
            class="px-4 py-2 rounded-xl text-xs font-black uppercase tracking-wider border shadow-sm inline-block"
            :class="{
              'bg-emerald-500 text-slate-950 border-emerald-400': shiftInfo.badge === 'ACTIVE',
              'bg-rose-900/80 text-rose-200 border-rose-700': shiftInfo.badge === 'OFF',
              'bg-amber-500 text-slate-950 border-amber-400': shiftInfo.badge === 'LEAVE',
              'bg-slate-700 text-slate-300 border-slate-600': shiftInfo.badge === 'LOADING'
            }"
          >
            {{ shiftInfo.badgeText }}
          </span>
        </div>
      </div>

      <!-- ===== 4. PRESENSI MASUK & PULANG STATUS CARDS ===== -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        
        <!-- Card Presensi Masuk -->
        <div class="bg-white rounded-3xl p-6 border border-[#F8BBD0] shadow-sm space-y-4 relative overflow-hidden">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-emerald-100 text-emerald-800 flex items-center justify-center font-bold">
                <span class="material-symbols-outlined text-xl">login</span>
              </div>
              <div>
                <h3 class="font-bold text-[#1b1c1c] text-sm">Presensi Masuk</h3>
                <p class="text-[11px] text-[#8a7176]">Absensi awal tugas harian</p>
              </div>
            </div>
            <span 
              class="px-2.5 py-1 rounded-full text-[10px] font-extrabold uppercase border"
              :class="authStore.todayStatus?.is_masuk ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-amber-50 text-amber-700 border-amber-200'"
            >
              {{ authStore.todayStatus?.is_masuk ? 'Sudah Masuk' : 'Belum Absen' }}
            </span>
          </div>

          <div class="p-4 bg-[#FFF5F8] rounded-2xl border border-[#F8BBD0]/60 space-y-2">
            <div class="flex items-center justify-between text-xs">
              <span class="text-[#574146] font-medium">Waktu Scan Masuk:</span>
              <span class="font-mono font-bold text-[#1b1c1c]">
                {{ authStore.todayStatus?.masuk?.timestamp ? formatTime(authStore.todayStatus.masuk.timestamp) : '-' }}
              </span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-[#574146] font-medium">Jarak Lokasi Posko:</span>
              <span class="font-mono font-bold text-[#ab2c5d]">
                {{ authStore.todayStatus?.masuk?.distance ? authStore.todayStatus.masuk.distance.toFixed(1) + ' meter' : '-' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Card Presensi Pulang -->
        <div class="bg-white rounded-3xl p-6 border border-[#F8BBD0] shadow-sm space-y-4 relative overflow-hidden">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-sky-100 text-sky-800 flex items-center justify-center font-bold">
                <span class="material-symbols-outlined text-xl">logout</span>
              </div>
              <div>
                <h3 class="font-bold text-[#1b1c1c] text-sm">Presensi Pulang</h3>
                <p class="text-[11px] text-[#8a7176]">Absensi selesai tugas harian</p>
              </div>
            </div>
            <span 
              class="px-2.5 py-1 rounded-full text-[10px] font-extrabold uppercase border"
              :class="authStore.todayStatus?.is_pulang ? 'bg-sky-50 text-sky-700 border-sky-200' : 'bg-slate-100 text-slate-500 border-slate-200'"
            >
              {{ authStore.todayStatus?.is_pulang ? 'Sudah Pulang' : 'Belum Absen' }}
            </span>
          </div>

          <div class="p-4 bg-[#FFF5F8] rounded-2xl border border-[#F8BBD0]/60 space-y-2">
            <div class="flex items-center justify-between text-xs">
              <span class="text-[#574146] font-medium">Waktu Scan Pulang:</span>
              <span class="font-mono font-bold text-[#1b1c1c]">
                {{ authStore.todayStatus?.pulang?.timestamp ? formatTime(authStore.todayStatus.pulang.timestamp) : '-' }}
              </span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-[#574146] font-medium">Jarak Lokasi Posko:</span>
              <span class="font-mono font-bold text-[#ab2c5d]">
                {{ authStore.todayStatus?.pulang?.distance ? authStore.todayStatus.pulang.distance.toFixed(1) + ' meter' : '-' }}
              </span>
            </div>
          </div>
        </div>

      </div>

      <!-- ===== 5. QUICK STATISTIK CARDS ===== -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        
        <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#F8BBD0] shadow-xs space-y-2">
          <div class="flex items-center justify-between text-[#ab2c5d]">
            <span class="text-xs font-bold uppercase tracking-wider text-[#8a7176]">Total Presensi</span>
            <span class="material-symbols-outlined text-xl">event_available</span>
          </div>
          <div class="text-xl sm:text-2xl font-black text-[#1b1c1c] font-mono">
            {{ totalPresensiCount }} <span class="text-xs font-sans font-normal text-[#8a7176]">hari</span>
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#F8BBD0] shadow-xs space-y-2">
          <div class="flex items-center justify-between text-emerald-600">
            <span class="text-xs font-bold uppercase tracking-wider text-[#8a7176]">Tepat Waktu</span>
            <span class="material-symbols-outlined text-xl">check_circle</span>
          </div>
          <div class="text-xl sm:text-2xl font-black text-[#1b1c1c] font-mono">
            100%
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#F8BBD0] shadow-xs space-y-2">
          <div class="flex items-center justify-between text-sky-600">
            <span class="text-xs font-bold uppercase tracking-wider text-[#8a7176]">Kehadiran</span>
            <span class="material-symbols-outlined text-xl">verified</span>
          </div>
          <div class="text-xl sm:text-2xl font-black text-[#1b1c1c] font-mono">
            Aktif
          </div>
        </div>

        <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#F8BBD0] shadow-xs space-y-2">
          <div class="flex items-center justify-between text-[#f06292]">
            <span class="text-xs font-bold uppercase tracking-wider text-[#8a7176]">Verifikasi 2FA</span>
            <span class="material-symbols-outlined text-xl">security</span>
          </div>
          <div class="text-sm font-black text-[#1b1c1c] truncate">
            {{ authStore.user?.totp_enabled ? 'Terverifikasi' : 'Belum Setup' }}
          </div>
        </div>

      </div>

      <!-- ===== 6. RIWAYAT PRESENSI TERBARU TABLE ===== -->
      <div class="bg-white rounded-3xl border border-[#F8BBD0] shadow-sm overflow-hidden space-y-4 p-5 sm:p-6">
        <div class="flex items-center justify-between border-b border-[#F8BBD0]/60 pb-3">
          <div class="flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d]">history</span>
            <h3 class="font-bold text-[#1b1c1c] text-sm sm:text-base">Riwayat Presensi Terbaru Anda</h3>
          </div>
          <router-link to="/intern/history" class="text-xs font-bold text-[#f06292] hover:text-[#ab2c5d] decoration-none flex items-center gap-1">
            <span>Lihat Semua</span>
            <span class="material-symbols-outlined text-sm">arrow_forward</span>
          </router-link>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                <th class="py-3 px-4">Tanggal</th>
                <th class="py-3 px-4">Presensi Masuk</th>
                <th class="py-3 px-4">Presensi Pulang</th>
                <th class="py-3 px-4">Status Geofence</th>
              </tr>
            </thead>
            <tbody class="text-xs divide-y divide-[#F8BBD0]">
              <tr v-for="(item, idx) in recentHistory" :key="idx" class="hover:bg-[#FFF5F8] transition-colors">
                <td class="py-3.5 px-4 font-bold text-[#1b1c1c]">
                  {{ item.date || item.tanggal || '-' }}
                </td>
                <td class="py-3.5 px-4 font-mono text-[#574146]">
                  {{ item.masuk || item.clock_in || '-' }}
                </td>
                <td class="py-3.5 px-4 font-mono text-[#574146]">
                  {{ item.pulang || item.clock_out || '-' }}
                </td>
                <td class="py-3.5 px-4">
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                    <span>Dalam Radius Posko</span>
                  </span>
                </td>
              </tr>
              <tr v-if="recentHistory.length === 0">
                <td colspan="4" class="py-8 text-center text-[#8a7176] text-xs">
                  Belum ada data riwayat presensi recorded.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>

    <!-- ===== MODAL PENGAJUAN SAKIT / IZIN / SHIFT ===== -->
    <transition name="fade">
      <div v-if="showLeaveModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="showLeaveModal = false">
        <div class="bg-white rounded-3xl max-w-md w-full p-6 border border-[#F8BBD0] shadow-2xl space-y-4">
          <div class="flex items-center justify-between border-b border-[#F8BBD0]/60 pb-3">
            <h3 class="font-bold text-base text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#f06292]">edit_document</span>
              <span>Form Pengajuan Izin / Sakit / Shift</span>
            </h3>
            <button @click="showLeaveModal = false" class="text-slate-400 hover:text-slate-600 border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <form @submit.prevent="submitLeaveRequest" class="space-y-4 text-xs">
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Kategori Pengajuan <span class="text-rose-500">*</span></label>
              <select v-model="formLeave.category" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-xl focus:outline-none focus:border-[#f06292]">
                <option value="SAKIT">Sakit (Dengan Surat Keterangan)</option>
                <option value="IZIN">Izin Keperluan Mendesak</option>
                <option value="TUKAR_SHIFT">Tukar Shift Magang</option>
              </select>
            </div>

            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tanggal Shift <span class="text-rose-500">*</span></label>
              <input v-model="formLeave.shift_date" type="date" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-xl focus:outline-none focus:border-[#f06292]" />
            </div>

            <div v-if="formLeave.category === 'TUKAR_SHIFT'" class="space-y-1">
              <label class="font-bold text-[#574146]">Nama Peserta Magang Pengganti</label>
              <input v-model="formLeave.replacement_name" type="text" placeholder="Contoh: Sarah Jenkins" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-xl focus:outline-none focus:border-[#f06292]" />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Alasan Keterangan <span class="text-rose-500">*</span></label>
              <textarea v-model="formLeave.reason" rows="3" required placeholder="Tuliskan alasan lengkap pengajuan..." class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-xl focus:outline-none focus:border-[#f06292]"></textarea>
            </div>

            <div v-if="formMessage" class="p-3 bg-emerald-50 text-emerald-800 rounded-xl font-bold border border-emerald-200">
              {{ formMessage }}
            </div>

            <div class="flex justify-end gap-3 pt-2">
              <button type="button" @click="showLeaveModal = false" class="px-4 py-2 border border-[#ddbfc5] rounded-xl font-bold text-[#574146] bg-white cursor-pointer">Batal</button>
              <button type="submit" :disabled="submitting" class="px-5 py-2 bg-[#F06292] hover:bg-[#ab2c5d] text-white rounded-xl font-bold border-0 cursor-pointer shadow-xs">
                Kirim Pengajuan
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </InternLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import InternLayout from '@/layouts/InternLayout.vue'

const authStore = useAuthStore()

// State
const showLeaveModal = ref(false)
const submitting = ref(false)
const formMessage = ref('')

const formLeave = ref({
  category: 'SAKIT',
  shift_date: new Date().toISOString().substring(0, 10),
  replacement_name: '',
  reason: ''
})

const shiftInfo = ref({
  shiftName: 'Shift Siaga 112 Active',
  shiftTimeStr: '08:00 – 16:00 WITA · Operasional Normal',
  badge: 'ACTIVE',
  badgeText: '🟢 SHIFT AKTIF'
})

// Computeds
const totalPresensiCount = computed(() => {
  return (authStore.presensiHistory || []).length
})

const recentHistory = computed(() => {
  return (authStore.presensiHistory || []).slice(0, 5)
})

// Functions
const formatTime = (ts: string) => {
  if (!ts) return '-'
  try {
    const d = new Date(ts)
    if (isNaN(d.getTime())) return ts
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')} WITA`
  } catch (e) {
    return ts
  }
}

const fetchShiftSchedule = async () => {
  try {
    if (authStore.todayStatus?.is_masuk && !authStore.todayStatus?.is_pulang) {
      shiftInfo.value = {
        shiftName: 'Shift Siaga Peserta Magang 112 (Sedang Berjalan)',
        shiftTimeStr: 'Status: Dalam Tugas Siaga',
        badge: 'ACTIVE',
        badgeText: '🟢 DALAM SIAGA 112'
      }
      return
    }
  } catch (err) {
    // Fallback
  }
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
      replacement_name: formLeave.value.replacement_name,
      reason: formLeave.value.reason,
      status: 'PENDING'
    }

    currentList.unshift(newReq)
    await axios.put('/api/presensi/leave-requests', currentList)

    formMessage.value = 'Pengajuan berhasil dikirim! Menunggu konfirmasi Super Admin.'
    setTimeout(() => {
      showLeaveModal.value = false
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
  await authStore.fetchProfile()
  await fetchShiftSchedule()
  await authStore.fetchTodayStatus()
  await authStore.fetchHistory()
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
