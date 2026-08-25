<template>
  <InternLayout>
    <div class="space-y-6 select-none font-sans">
      <!-- Page Header -->
      <div class="hidden sm:flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Dashboard Status Siaga Peserta Magang</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs hidden sm:block">Panel operasional siaga 112, status pendaftaran presensi, dan riwayat tugas presensi harian.</p>
        </div>
      </div>
      <div class="bg-white rounded-3xl p-6 sm:p-7 border border-slate-200/90 shadow-sm space-y-5">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
          
          <!-- Profile Info with Avatar Initials -->
          <div class="flex items-start gap-4 min-w-0">
            <!-- Avatar Initials Badge -->
            <div class="w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white font-display font-black text-2xl flex items-center justify-center shrink-0 shadow-md">
              {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'C' }}
            </div>

            <div class="space-y-1 flex-1 min-w-0">
              <!-- Peserta Magang Status Pill -->
              <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 bg-rose-50 text-rose-800 rounded-full text-[10px] sm:text-xs font-extrabold border border-rose-200">
                <span class="w-1.5 h-1.5 rounded-full bg-rose-600 animate-ping"></span>
                <span>PESERTA MAGANG POSKO SIAGA 112</span>
              </div>

              <!-- 1. Selamat Tugas, [Nama]! -->
              <h2 class="text-base sm:text-xl md:text-2xl font-display font-black text-slate-900 leading-tight truncate">
                Selamat Tugas, {{ authStore.user?.name || 'Petugas Peserta Magang' }}!
              </h2>

              <!-- 2. NIP -->
              <div class="text-xs font-mono font-extrabold text-rose-700">
                NIP. {{ authStore.user?.nip || '-' }}
              </div>

              <!-- 3. Jabatan -->
              <div class="text-xs font-extrabold text-slate-800 uppercase tracking-wide">
                {{ authStore.user?.jabatan || 'OPERATOR LAYANAN OPERASIONAL' }}
              </div>

              <!-- 4. Unit Kerja -->
              <div class="text-xs text-slate-500 font-medium">
                {{ authStore.user?.unit_kerja || 'Diskominfo Kabupaten Bulukumba' }}
              </div>
            </div>
          </div>

          <!-- Action Buttons Bar -->
          <div class="grid grid-cols-1 sm:grid-cols-2 md:flex md:items-center gap-2.5 pt-3 md:pt-0 border-t md:border-0 border-slate-100">
            <button 
              @click="showLeaveModal = true"
              class="w-full sm:w-auto px-4 py-3 bg-amber-50/80 hover:bg-amber-100 text-amber-900 border border-amber-200/90 font-bold text-xs rounded-xl shadow-2xs transition-all flex items-center justify-center gap-2 cursor-pointer active:scale-98"
            >
              <span class="material-symbols-outlined text-[18px] text-amber-700">assignment_ind</span>
              <span>Ajukan Sakit / Tukar Shift</span>
            </button>

            <router-link 
              to="/intern/scan" 
              class="w-full sm:w-auto px-5 py-3 bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-2 decoration-none cursor-pointer active:scale-98"
            >
              <span class="material-symbols-outlined text-[18px]">qr_code_scanner</span>
              <span>Buka Kamera Scan QR</span>
            </router-link>
          </div>

        </div>
      </div>

      <!-- Real-Time Status Shift Hari Ini Banner -->
      <div 
        class="rounded-3xl p-5 sm:p-6 border shadow-md transition-all flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 relative overflow-hidden backdrop-blur-md"
        :class="{
          'bg-gradient-to-r from-emerald-950 via-teal-900 to-slate-900 border-emerald-600/50 text-white shadow-emerald-950/20': shiftInfo.badge === 'ACTIVE',
          'bg-gradient-to-r from-rose-950 via-slate-900 to-rose-900 border-rose-600/50 text-white shadow-rose-950/20': shiftInfo.badge === 'OFF',
          'bg-gradient-to-r from-amber-950 via-slate-900 to-amber-900 border-amber-600/50 text-white shadow-amber-950/20': shiftInfo.badge === 'LEAVE',
          'bg-gradient-to-r from-slate-800 via-slate-900 to-slate-800 border-slate-600 text-white animate-pulse': shiftInfo.badge === 'LOADING'
        }"
      >
        <!-- Background Ambient Glow Effect -->
        <div v-if="shiftInfo.badge === 'ACTIVE'" class="absolute -right-8 -bottom-8 w-36 h-36 bg-emerald-500/10 rounded-full blur-2xl pointer-events-none"></div>

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
              <span>STATUS OPERASIONAL SIAGA 112</span>
            </div>
            <h3 class="text-base font-display font-black text-white leading-tight">
              {{ shiftInfo.shiftName }}
            </h3>
            <p class="text-xs font-semibold text-slate-200/90 leading-relaxed">
              {{ shiftInfo.shiftTimeStr }}
            </p>
          </div>
        </div>
      </div>

      <!-- Jadwal Shift Berikutnya (muncul jika hari ini libur/off) -->
      <transition name="fade-next">
        <div
          v-if="nextShiftInfo && shiftInfo.badge === 'OFF'"
          class="flex items-start gap-3 px-5 py-4 bg-amber-50 border border-amber-200 rounded-2xl shadow-xs"
        >
          <div class="w-10 h-10 rounded-xl bg-amber-100 border border-amber-300 flex items-center justify-center shrink-0 text-amber-700">
            <span class="material-symbols-outlined text-[22px]">event_upcoming</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-[9px] font-extrabold text-amber-700 uppercase tracking-widest mb-0.5">Jadwal Shift Berikutnya</div>
            <div class="text-sm font-display font-black text-slate-900 leading-tight">
              {{ nextShiftInfo.label }}
              <span v-if="nextShiftInfo.asReplacer" class="text-amber-600 font-bold text-xs"> · Pengganti {{ nextShiftInfo.replacedName?.split(',')[0] }}</span>
            </div>
            <div class="flex items-center gap-2 mt-1 flex-wrap">
              <span class="text-xs font-bold text-slate-700">
                {{ nextShiftInfo.dayName }}, {{ nextShiftInfo.formattedDate }}
              </span>
              <span class="text-slate-300">·</span>
              <span class="text-xs font-mono font-bold text-emerald-700">
                {{ nextShiftInfo.start }} – {{ nextShiftInfo.end }} WITA
              </span>
            </div>
          </div>
          <div class="text-right shrink-0">
            <div class="text-[9px] font-bold text-amber-600 uppercase tracking-wider">{{ nextShiftInfo.dateStr }}</div>
          </div>
        </div>
        <div v-else-if="shiftInfo.badge === 'OFF' && !nextShiftInfo" class="flex items-center gap-2.5 px-4 py-3 bg-slate-50 border border-slate-200 rounded-2xl text-xs text-slate-500">
          <span class="material-symbols-outlined text-[18px] text-slate-400">calendar_month</span>
          <span>Tidak ada jadwal shift tersisa di bulan ini. Jadwal bulan berikutnya belum tersedia.</span>
        </div>
      </transition>

      <!-- Real-Time Status Siaga Cards (Database PostgreSQL) -->

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        
        <!-- Card 1: Status Masuk Siaga -->
        <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">Presensi Masuk Siaga</span>
            <span 
              class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold"
              :class="authStore.todayStatus?.is_masuk ? 'bg-emerald-100 text-emerald-800 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
            >
              {{ authStore.todayStatus?.is_masuk ? '✓ Sudah Masuk' : 'Belum Presensi' }}
            </span>
          </div>

          <div class="text-2xl font-display font-black text-slate-900 font-mono">
            {{ authStore.todayStatus?.is_masuk ? formatTimeDisplay(authStore.todayStatus?.masuk?.timestamp || authStore.todayStatus?.clock_in_time) : '--:--:--' }}
          </div>

          <div class="text-xs text-slate-500 flex items-center justify-between">
            <span>Jarak ke Posko:</span>
            <strong class="text-slate-900 font-mono font-bold">
              {{ authStore.todayStatus?.is_masuk ? formatDistance(authStore.todayStatus?.masuk) : '--' }}
            </strong>
          </div>
        </div>

        <!-- Card 2: Status Pulang / Selesai Siaga -->
        <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">Presensi Selesai Siaga</span>
            <span 
              class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold"
              :class="authStore.todayStatus?.is_pulang ? 'bg-emerald-100 text-emerald-800 border border-emerald-200' : (authStore.todayStatus?.is_masuk ? 'bg-amber-100 text-amber-800 border border-amber-200' : 'bg-slate-100 text-slate-500 border border-slate-200')"
            >
              {{ authStore.todayStatus?.is_pulang ? '✓ Sudah Selesai' : (authStore.todayStatus?.is_masuk ? '⏳ Belum Scan Pulang' : 'Belum Presensi') }}
            </span>
          </div>

          <div class="text-2xl font-display font-black text-slate-900 font-mono">
            {{ authStore.todayStatus?.is_pulang ? formatTimeDisplay(authStore.todayStatus?.pulang?.timestamp || authStore.todayStatus?.clock_out_time) : '--:--:--' }}
          </div>

          <div class="text-xs text-slate-500 flex items-center justify-between">
            <span>Jarak ke Posko:</span>
            <strong class="text-slate-900 font-mono font-bold">
              {{ authStore.todayStatus?.is_pulang ? formatDistance(authStore.todayStatus?.pulang) : '--' }}
            </strong>
          </div>
        </div>

        <!-- Card 3: Geofence Radar Posko (Database posko_locations) -->
        <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider">Geofence Posko Siaga</span>
            <span class="px-2.5 py-0.5 bg-emerald-100 text-emerald-800 border border-emerald-200 rounded-full text-[10px] font-bold">
              Radius {{ poskoInfo.radius }}m
            </span>
          </div>

          <div class="text-sm font-bold text-slate-900 leading-tight">
            {{ poskoInfo.name }}
          </div>

          <div class="text-xs text-slate-500 font-mono flex items-center justify-between">
            <span>Koordinat:</span>
            <span>{{ poskoInfo.lat.toFixed(4) }}, {{ poskoInfo.lng.toFixed(4) }}</span>
          </div>
        </div>
      </div>

      <!-- Leave Request Modal -->
      <div v-if="showLeaveModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs z-50 flex items-center justify-center p-4 overflow-y-auto select-none">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 sm:p-8 space-y-6 shadow-2xl">
          <div class="flex items-center justify-between border-b border-slate-200 pb-4">
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700 text-2xl">assignment_ind</span>
              <h3 class="text-lg font-display font-black text-slate-900">Form Pengajuan Sakit / Izin / Tukar Shift</h3>
            </div>
            <button @click="showLeaveModal = false" class="text-slate-400 hover:text-slate-700 text-xl font-bold border-0 bg-transparent cursor-pointer">✕</button>
          </div>

          <form @submit.prevent="submitLeaveRequest" class="space-y-4 text-xs">
            <div v-if="formMessage" class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-900 font-bold">
              {{ formMessage }}
            </div>

            <div class="space-y-1">
              <label class="font-bold text-slate-700">Kategori Pengajuan <span class="text-rose-600">*</span></label>
              <select v-model="form.category" required class="w-full p-2.5 bg-slate-50 border border-slate-300 rounded-xl font-bold text-slate-900 focus:ring-2 focus:ring-rose-500">
                <option value="Sakit">Sakit (Wajib Surat Dokter)</option>
                <option value="Izin">Izin Tugas Mendadak</option>
                <option value="Tukar Shift">Tukar Shift (Replacement)</option>
              </select>
            </div>

            <div class="space-y-1">
              <label class="font-bold text-slate-700">Tanggal Shift yang Berhalangan <span class="text-rose-600">*</span></label>
              <input v-model="form.shift_date" type="text" placeholder="Contoh: 01-08-2026" required class="w-full p-2.5 bg-slate-50 border border-slate-300 rounded-xl font-bold text-slate-900" />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-slate-700">Pilih Petugas Pengganti (Standby Backup)</label>
              <select v-model="form.replacement_name" class="w-full p-2.5 bg-slate-50 border border-slate-300 rounded-xl font-bold text-slate-900">
                <option value="">-- Serahkan Pemilihan ke Admin --</option>
                <option v-for="off in replacementOptions" :key="off" :value="off">{{ off }}</option>
              </select>
            </div>

            <div class="space-y-1">
              <label class="font-bold text-slate-700">Alasan Lengkap <span class="text-rose-600">*</span></label>
              <textarea v-model="form.reason" rows="3" placeholder="Tuliskan keterangan sakit atau alasan izin..." required class="w-full p-2.5 bg-slate-50 border border-slate-300 rounded-xl font-medium text-slate-900"></textarea>
            </div>

            <div class="space-y-1">
              <label class="font-bold text-slate-700">Lampirkan Bukti Surat Dokter / Dokumen (URL Foto/PDF)</label>
              <input v-model="form.attachment_url" type="text" placeholder="https://..." class="w-full p-2.5 bg-slate-50 border border-slate-300 rounded-xl font-mono text-[11px] text-slate-800" />
            </div>

            <div class="flex justify-end gap-3 border-t border-slate-200 pt-4">
              <button type="button" @click="showLeaveModal = false" class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-xl border border-slate-300 cursor-pointer">Batal</button>
              <button type="submit" :disabled="submitting" class="px-6 py-2.5 bg-rose-700 hover:bg-rose-800 text-white font-extrabold rounded-xl shadow-md cursor-pointer border-0 disabled:opacity-50">
                {{ submitting ? 'Kirim...' : 'Kirim Pengajuan' }}
              </button>
            </div>
          </form>
        </div>
      </div>

    </div>
  </InternLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import InternLayout from '@/layouts/InternLayout.vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const showLeaveModal = ref(false);
const submitting = ref(false);
const formMessage = ref('');

const poskoInfo = ref({
  name: 'Posko Siaga NTPD 112 Kab. Bulukumba',
  lat: -5.5645,
  lng: 120.1945,
  radius: 2.0
});

const shiftInfo = ref({
  shiftName: 'Memuat jadwal shift...',
  shiftTimeStr: 'Mengambil data dari server...',
  badge: 'LOADING',
  badgeText: '⏳ MEMUAT DATA'
});


const nextShiftInfo = ref(null); // { label, dateStr, dayName, formattedDate, start, end }

const replacementOptions = [
  'A.Mappalua, S.Pd',
  'Suherman, S.Pd',
  'Riswandi Risman',
  'Abil Kizri',
  'Imam Ardiyansah',
  'Abd.Rahim',
  'Munawir Syadzali',
  'Abdullah, S.Kep., Ns',
  'Ismail, S.Sos',
  'Aldi Afdali Saputra'
];

const now = new Date()
const todayDateStr = [
  String(now.getDate()).padStart(2, '0'),
  String(now.getMonth() + 1).padStart(2, '0'),
  now.getFullYear()
].join('-')

const form = ref({
  category: 'Sakit',
  shift_date: todayDateStr,
  replacement_name: '',
  reason: '',
  attachment_url: ''
});

function formatTimeDisplay(val) {
  if (!val) return '--:--:--'
  const str = String(val).trim()
  const match = str.match(/\d{2}:\d{2}:\d{2}/)
  if (match) {
    return `${match[0]} WITA`
  }
  return str.includes('WITA') ? str : `${str} WITA`
}

const formatDistance = (record) => {
  if (!record) return '--'
  const val = record.distance_meters ?? record.distanceMeters ?? record.distance
  if (val !== undefined && val !== null && val !== '') {
    const num = Number(val)
    if (!isNaN(num)) return `${num.toFixed(1)} Meter`
  }
  return '0.8 Meter'
}

const recentHistoryList = computed(() => {
  if (authStore.presensiHistory && Array.isArray(authStore.presensiHistory) && authStore.presensiHistory.length > 0) {
    return authStore.presensiHistory.slice(0, 5).map((item) => {
      const rawTs = item.timestamp || ''
      const parts = rawTs.split(' ')
      let dateKey = parts[0] || '-'
      let timePart = parts[1] || parts[0] || '--:--:--'

      if (dateKey.includes('-') && dateKey.split('-')[0].length === 4) {
        const [y, m, d] = dateKey.split('-')
        dateKey = `${d}-${m}-${y}`
      }

      const timeFormatted = timePart.includes('WITA') ? timePart : `${timePart} WITA`

      return {
        id: item.id,
        date: dateKey,
        type: item.type || 'MASUK',
        clockTime: timeFormatted,
        distance: item.distance_meters ? item.distance_meters.toFixed(1) : '0.8'
      }
    })
  }
  return []
})

const fetchPoskoInfo = async () => {
  try {
    const res = await axios.get('/api/presensi/posko-qr');
    if (res.data) {
      if (res.data.name) poskoInfo.value.name = res.data.name;
      if (res.data.coordinates) {
        if (res.data.coordinates.latitude) poskoInfo.value.lat = res.data.coordinates.latitude;
        if (res.data.coordinates.longitude) poskoInfo.value.lng = res.data.coordinates.longitude;
        if (res.data.coordinates.max_radius_meters) poskoInfo.value.radius = res.data.coordinates.max_radius_meters;
      }
    }
  } catch (e) {}
};

const fetchShiftSchedule = async () => {
  try {
    // Priority 0: Active unclosed shift (e.g. Night Shift started yesterday, awaiting clock-out today)
    if (authStore.todayStatus?.is_masuk && !authStore.todayStatus?.is_pulang) {
      const rawTs = authStore.todayStatus?.masuk?.timestamp || authStore.todayStatus?.clock_in_time || ''
      const cleanTime = formatTimeDisplay(rawTs)
      const isNightShift = rawTs.includes('Shift Kemarin') || rawTs.includes('Kemarin')

      shiftInfo.value = {
        shiftName: isNightShift ? 'Shift 2 (Malam Lintas Hari)' : 'Shift Siaga 112 (Sedang Berjalan)',
        shiftTimeStr: `Scan Masuk: ${cleanTime} ${isNightShift ? '(Piket Malam)' : ''} · Status: Dalam Tugas Siaga`,
        badge: 'ACTIVE',
        badgeText: '🟢 DALAM SIAGA 112'
      }
      return
    }

    const now = new Date()
    const todayStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    const userNip = (authStore.user?.nip || '').replace(/\s/g, '')
    const userName = authStore.user?.name || ''

    const schedRes = await axios.get('/api/admin/schedules')
    const schedData = schedRes.data?.schedules
    if (!schedData || !schedData.daysInMonth || schedData.daysInMonth.length === 0) {
      shiftInfo.value = {
        shiftName: 'Jadwal Shift Belum Dikonfigurasi',
        shiftTimeStr: 'Admin belum menyimpan rincian jadwal shift bulan ini',
        badge: 'OFF',
        badgeText: '⚠ BELUM ADA JADWAL'
      }
      return
    }

    const dayEntry = schedData.daysInMonth.find(d =>
      d.dateStr === todayStr || d.date === todayStr
    ) || schedData.daysInMonth[now.getDate() - 1]

    if (!dayEntry) {
      shiftInfo.value = {
        shiftName: 'Jadwal Hari Ini Tidak Ditemukan',
        shiftTimeStr: 'Hubungi Admin untuk cek tanggal jadwal shift',
        badge: 'OFF',
        badgeText: '⚠ TANGGAL TIDAK SESUAI'
      }
      return
    }

    const teams = schedData.teams || []
    const shiftMode = schedData.shiftMode || 2

    const slotTimes = [
      { key: 'shift1', replKey: 'replacementsShift1', start: '08:00', end: shiftMode === 3 ? '16:00' : '20:00', label: shiftMode === 3 ? 'Shift Pagi' : 'Shift Pagi' },
      { key: 'shift2', replKey: 'replacementsShift2', start: shiftMode === 3 ? '16:00' : '20:00', end: shiftMode === 3 ? '24:00' : '08:00', label: shiftMode === 3 ? 'Shift Sore' : 'Shift Malam' },
      { key: 'shift3', replKey: 'replacementsShift3', start: '00:00', end: '08:00', label: 'Shift Malam' }
    ]

    // ─── 1. Cek apakah user adalah PENGGANTI hari ini ───
    for (let i = 0; i < slotTimes.length; i++) {
      const slot = slotTimes[i]
      const replacements = dayEntry[slot.replKey] || []
      const myReplacement = replacements.find(r =>
        r.replacerNip && r.replacerNip.replace(/\s/g, '') === userNip
      )
      if (myReplacement) {
        shiftInfo.value = {
          shiftName: `${slot.label} — Pengganti ${myReplacement.replacedName.split(',')[0]}`,
          shiftTimeStr: `${slot.start} – ${slot.end} WITA · 🟡 Pengganti Resmi`,
          badge: 'ACTIVE',
          badgeText: '🟢 PENGGANTI AKTIF'
        }
        return
      }
    }

    // ─── 2. Cek tim utama user ───
    let userTeamCode = ''
    for (const t of teams) {
      if (t.members && Array.isArray(t.members)) {
        const found = t.members.some(m => {
          const nipMatch = userNip && m.nip && m.nip.replace(/\s/g, '') === userNip
          const nameMatch = userName && m.name && m.name.toLowerCase().includes(userName.toLowerCase().split(' ')[0])
          return nipMatch || nameMatch
        })
        if (found) { userTeamCode = t.code; break }
      }
    }

    if (!userTeamCode) {
      shiftInfo.value = {
        shiftName: 'Belum Terdaftar di Tim',
        shiftTimeStr: 'Hubungi Admin untuk pendaftaran tim shift',
        badge: 'OFF',
        badgeText: '⚠ BELUM TERDAFTAR'
      }
      return
    }

    // ─── 3. Cek apakah tim libur hari ini ───
    const isOff = dayEntry.offTeams && dayEntry.offTeams.includes(userTeamCode)
    if (isOff) {
      shiftInfo.value = {
        shiftName: `HARI LIBUR — Tim ${userTeamCode}`,
        shiftTimeStr: 'Tidak Ada Shift Hari Ini · Standby Backup',
        badge: 'OFF',
        badgeText: '🔒 LIBUR / OFF'
      }
      findNextShift(schedData, userTeamCode, todayStr, slotTimes, userNip)
      return
    }

    // ─── 4. Tentukan shift tim ───
    for (let i = 0; i < slotTimes.length; i++) {
      const s = slotTimes[i]
      if (dayEntry[s.key] === userTeamCode) {
        shiftInfo.value = {
          shiftName: `${s.label} — Tim ${userTeamCode} (${s.start}–${s.end} WITA)`,
          shiftTimeStr: `${s.start} – ${s.end} WITA`,
          badge: 'ACTIVE',
          badgeText: '🟢 SHIFT AKTIF'
        }
        return
      }
    }

    // ─── 5. Tidak ditemukan di shift hari ini, cari jadwal berikutnya ───
    shiftInfo.value = {
      shiftName: `Standby — Tim ${userTeamCode}`,
      shiftTimeStr: 'Tim tidak terjadwal shift hari ini',
      badge: 'OFF',
      badgeText: '⚡ STANDBY'
    }
    findNextShift(schedData, userTeamCode, todayStr, slotTimes, userNip)

  } catch (e) {
    shiftInfo.value = {
      shiftName: 'Gagal Memuat Jadwal',
      shiftTimeStr: 'Terjadi kesalahan saat mengambil data jadwal dari server',
      badge: 'OFF',
      badgeText: '⚠ TERJADI KESALAHAN'
    }
  }
}

// Cari jadwal shift berikutnya untuk user
const findNextShift = (schedData, userTeamCode, todayStr, slotTimes, userNipClean) => {
  nextShiftInfo.value = null
  const days = schedData.daysInMonth || []
  const dayNamesIndo = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  const monthNamesIndo = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']

  for (const d of days) {
    const dStr = d.dateStr || d.date || ''
    if (dStr <= todayStr) continue // hanya hari mendatang

    // Cek apakah user sebagai pengganti di hari itu
    for (let i = 0; i < slotTimes.length; i++) {
      const slot = slotTimes[i]
      const replacements = d[slot.replKey] || []
      const myEntry = replacements.find(r =>
        r.replacerNip && r.replacerNip.replace(/\s/g, '') === (userNipClean || '')
      )
      if (myEntry) {
        const dateObj = new Date(dStr)
        nextShiftInfo.value = {
          label: slot.label,
          dateStr: dStr,
          dayName: dayNamesIndo[dateObj.getDay()],
          formattedDate: `${String(dateObj.getDate()).padStart(2, '0')} ${monthNamesIndo[dateObj.getMonth()]} ${dateObj.getFullYear()}`,
          start: slot.start,
          end: slot.end,
          asReplacer: true,
          replacedName: myEntry.replacedName
        }
        return
      }
    }

    // Cek apakah tim user dijadwalkan shift di hari itu
    for (const slot of slotTimes) {
      if (d[slot.key] === userTeamCode) {
        const dateObj = new Date(dStr)
        nextShiftInfo.value = {
          label: slot.label,
          dateStr: dStr,
          dayName: dayNamesIndo[dateObj.getDay()],
          formattedDate: `${String(dateObj.getDate()).padStart(2, '0')} ${monthNamesIndo[dateObj.getMonth()]} ${dateObj.getFullYear()}`,
          start: slot.start,
          end: slot.end,
          asReplacer: false,
          replacedName: ''
        }
        return
      }
    }
  }
}



async function submitLeaveRequest() {
  submitting.value = true;
  formMessage.value = '';
  try {

    let currentList = [];
    try {
      const resGet = await axios.get('/api/presensi/leave-requests');
      if (resGet.data && Array.isArray(resGet.data.requests)) {
        currentList = resGet.data.requests;
      }
    } catch (e) {}

    const newReq = {
      id: Date.now(),
      created_at: new Date().toISOString().replace('T', ' ').substring(0, 16),
      user_name: authStore.user?.name || 'Peserta Magang',
      user_nip: authStore.user?.nip || '19940503 202521 1 138',
      category: form.value.category,
      shift_date: form.value.shift_date,
      replacement_name: form.value.replacement_name,
      reason: form.value.reason,
      attachment_url: form.value.attachment_url,
      status: 'PENDING'
    };

    currentList.unshift(newReq);
    await axios.put('/api/presensi/leave-requests', currentList);

    formMessage.value = 'Pengajuan berhasil dikirim! Menunggu konfirmasi Super Admin.';
    setTimeout(() => {
      showLeaveModal.value = false;
      formMessage.value = '';
      form.value.reason = '';
    }, 1500);
  } catch (err) {
    formMessage.value = 'Gagal mengirim pengajuan.';
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await authStore.fetchProfile();
  fetchPoskoInfo();
  await fetchShiftSchedule();
  await authStore.fetchTodayStatus();
  await authStore.fetchHistory();
});
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
