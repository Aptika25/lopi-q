<template>
  <InternLayout>
    <div class="space-y-6 select-none font-sans w-full pb-28 sm:pb-8">
      <!-- Page Header -->
      <div class="hidden sm:flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Laporan Kehadiran Saya</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs hidden sm:block">Riwayat absensi dan log presensi harian petugas Peserta Magang di Posko Siaga NTPD 112 Bulukumba.</p>
        </div>
      </div>
      <div class="bg-white rounded-3xl p-6 sm:p-7 border border-slate-200 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-start gap-4">
          <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white font-display font-black text-2xl flex items-center justify-center shrink-0 shadow-md">
            {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'C' }}
          </div>
          <div class="flex-1 min-w-0 space-y-1">
            <!-- 1. Nama -->
            <h2 class="text-lg font-display font-black text-slate-900 leading-tight">
              {{ authStore.user?.name || 'Peserta Magang 112' }}
            </h2>
            
            <!-- 2. NIP. -->
            <div class="text-xs text-rose-700 font-mono font-bold">
              NIP. {{ authStore.user?.nip || '-' }}
            </div>

            <!-- 3. Jabatan (Tanpa Badge, Plain Text) -->
            <div class="text-xs text-slate-800 font-bold uppercase tracking-wide">
              {{ authStore.user?.jabatan || 'OPERATOR LAYANAN OPERASIONAL' }}
            </div>

            <!-- 4. Unit Kerja -->
            <div class="text-xs text-slate-500 font-medium pt-0.5">
              {{ authStore.user?.unit_kerja || 'Diskominfo Kabupaten Bulukumba' }}
            </div>
          </div>
        </div>

        <!-- Keluar Sesi Button (Mobile Only - Inside Profile Card) -->
        <div class="sm:hidden w-full pt-3 border-t border-slate-100">
          <button 
            @click="handleLogout"
            class="w-full py-2.5 px-4 bg-rose-50/80 hover:bg-rose-100 border border-rose-200/80 text-rose-600 font-bold text-xs rounded-full transition-all flex items-center justify-center gap-2 cursor-pointer shadow-2xs group"
          >
            <span>Keluar Sesi</span>
            <span class="material-symbols-outlined text-[16px] text-rose-600 group-hover:translate-x-0.5 transition-transform">logout</span>
          </button>
        </div>
      </div>

      <!-- History Table Container -->
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-slate-200 shadow-sm space-y-6">
        <!-- Card Header: Title on 1 Single Line + Refresh Button -->
        <div class="flex items-center justify-between pb-4 border-b border-slate-200 gap-2">
          <div class="min-w-0 flex-1">
            <h3 class="text-xs sm:text-base md:text-lg font-display font-black text-slate-900 flex items-center gap-1.5 sm:gap-2">
              <span class="material-symbols-outlined text-rose-700 text-[18px] sm:text-[22px] shrink-0">description</span>
              <span class="truncate">Laporan Riwayat Kehadiran</span>
            </h3>
            <p class="text-xs text-slate-500 mt-1 hidden sm:block">Daftar rekapitulasi presensi siaga 112 yang terverifikasi Geofence.</p>
          </div>

          <button 
            @click="fetchHistory"
            class="px-2.5 py-1.5 sm:px-3.5 sm:py-2 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold text-xs rounded-xl border border-slate-300 flex items-center gap-1 cursor-pointer border-0 shrink-0"
          >
            <span class="material-symbols-outlined text-[14px] sm:text-[16px]">refresh</span>
            <span>Refresh</span>
          </button>
        </div>

        <!-- Filter & Sorting Toolbar (Bulan, Tahun, Tampilkan Per Halaman) -->
        <div class="bg-slate-50 p-3 sm:p-4 rounded-2xl border border-slate-200 text-xs space-y-2.5 sm:space-y-0 sm:flex sm:items-center sm:justify-between sm:gap-3 overflow-hidden">
          <!-- Left: Month & Year Filters (2-col grid on mobile) -->
          <div class="grid grid-cols-2 sm:flex sm:items-center gap-2 flex-1">
            <select v-model="selectedMonth" class="w-full sm:w-auto px-2.5 py-1.5 bg-white border border-slate-300 rounded-xl font-bold text-slate-800 shadow-2xs focus:ring-2 focus:ring-rose-500 text-xs">
              <option v-for="m in monthOptions" :key="m.value" :value="m.value">{{ m.label }}</option>
            </select>

            <select v-model="selectedYear" class="w-full sm:w-auto px-2.5 py-1.5 bg-white border border-slate-300 rounded-xl font-bold text-slate-800 shadow-2xs focus:ring-2 focus:ring-rose-500 text-xs">
              <option v-for="y in yearOptions" :key="y.value" :value="y.value">{{ y.label }}</option>
            </select>
          </div>

          <!-- Right: Page Size & Reset Filter -->
          <div class="flex items-center justify-between sm:justify-end gap-2 pt-2 sm:pt-0 border-t sm:border-0 border-slate-200/80">
            <div class="flex items-center gap-1 text-slate-600 font-bold">
              <span class="text-[11px] text-slate-500 shrink-0">Tampilkan:</span>
              <select v-model="pageSize" class="px-2 py-1.5 bg-white border border-slate-300 rounded-xl font-bold text-slate-800 shadow-2xs focus:ring-2 focus:ring-rose-500 text-xs">
                <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }} Data</option>
              </select>
            </div>

            <button 
              v-if="selectedMonth !== 'ALL' || selectedYear !== 'ALL'"
              @click="selectedMonth = 'ALL'; selectedYear = 'ALL'"
              class="px-2.5 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 font-extrabold rounded-xl border border-rose-200 text-xs cursor-pointer whitespace-nowrap shrink-0"
            >
              Reset Filter
            </button>
          </div>
        </div>

        <!-- Mobile Card View (No Horizontal Scroll Needed) -->
        <div class="block sm:hidden space-y-3">
          <div 
            v-for="item in paginatedHistoryList" 
            :key="item.id"
            class="bg-white rounded-3xl p-4 border border-slate-200/90 shadow-2xs space-y-3 hover:border-slate-300 transition-all"
          >
            <!-- Card Header: Date & Shift -->
            <div class="flex items-center justify-between border-b border-slate-100 pb-2.5">
              <div class="flex items-center gap-1.5 font-mono font-bold text-slate-900 text-xs">
                <span class="material-symbols-outlined text-[16px] text-rose-600">calendar_today</span>
                <span>{{ item.date }}</span>
              </div>
              <span class="px-2.5 py-0.5 bg-rose-50 border border-rose-200/80 text-rose-700 rounded-full text-[10px] font-extrabold">
                {{ item.shiftName }}
              </span>
            </div>

            <!-- Card Body: Clock In & Clock Out Grid with Symmetrical Equal-Height Boxes -->
            <div class="grid grid-cols-2 gap-2.5 text-xs">
              <!-- Jam Masuk Box -->
              <div class="bg-emerald-50/60 p-3 rounded-2xl border border-emerald-100/90 flex flex-col justify-between space-y-2 min-h-[92px]">
                <div class="flex items-center justify-between">
                  <span class="text-[10px] font-extrabold text-emerald-800 uppercase tracking-wider flex items-center gap-1">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Jam Masuk
                  </span>
                </div>

                <div class="font-mono text-xs font-black text-slate-900">
                  {{ item.clockIn }}
                </div>

                <div>
                  <span v-if="item.clockIn !== '--:--:--'" class="text-[10px] font-bold text-emerald-800 bg-white/90 px-2 py-0.5 rounded-md border border-emerald-200/70 inline-flex items-center gap-1 shadow-2xs">
                    ðŸ“ {{ item.masukDistance }}
                  </span>
                  <span v-else class="text-[10px] font-semibold text-slate-400 italic">
                    --
                  </span>
                </div>
              </div>

              <!-- Jam Pulang Box -->
              <div 
                class="p-3 rounded-2xl border flex flex-col justify-between space-y-2 min-h-[92px] transition-colors"
                :class="item.clockOut !== '--:--:--' ? 'bg-amber-50/60 border-amber-100/90' : 'bg-slate-50/70 border-slate-200/70'"
              >
                <div class="flex items-center justify-between">
                  <span 
                    class="text-[10px] font-extrabold uppercase tracking-wider flex items-center gap-1"
                    :class="item.clockOut !== '--:--:--' ? 'text-amber-800' : 'text-slate-400'"
                  >
                    <span class="w-1.5 h-1.5 rounded-full" :class="item.clockOut !== '--:--:--' ? 'bg-amber-500' : 'bg-slate-300'"></span> Jam Pulang
                  </span>
                </div>

                <div class="font-mono text-xs font-black" :class="item.clockOut !== '--:--:--' ? 'text-slate-900' : 'text-slate-400'">
                  {{ item.clockOut }}
                </div>

                <div>
                  <span v-if="item.clockOut !== '--:--:--'" class="text-[10px] font-bold text-amber-800 bg-white/90 px-2 py-0.5 rounded-md border border-amber-200/70 inline-flex items-center gap-1 shadow-2xs">
                    ðŸ“ {{ item.pulangDistance }}
                  </span>
                  <span v-else class="text-[10px] font-medium text-slate-400 bg-white/70 px-2 py-0.5 rounded-md border border-slate-200/60 inline-flex items-center gap-1">
                    â³ Belum Scan
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="filteredHistoryList.length === 0" class="py-6 text-center bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
            <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
              <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
              <span>Tidak ada riwayat presensi yang sesuai dengan filter.</span>
            </div>
          </div>
        </div>

        <!-- Desktop History Table (Hidden on Mobile) -->
        <div class="hidden sm:block overflow-x-auto rounded-2xl border border-slate-200 shadow-xs">
          <table class="w-full text-left text-xs text-slate-700">
            <thead class="bg-slate-100 text-slate-700 font-bold uppercase tracking-wider border-b border-slate-200">
              <tr>
                <th class="py-3.5 px-4">Tanggal Siaga</th>
                <th class="py-3.5 px-4">Jadwal Shift</th>
                <th class="py-3.5 px-4">Jam Masuk (Jarak)</th>
                <th class="py-3.5 px-4">Jam Pulang (Jarak)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 bg-white">
              <tr v-for="item in paginatedHistoryList" :key="item.id" class="hover:bg-slate-50 transition-colors">
                <td class="py-3.5 px-4 font-mono font-bold text-slate-900 whitespace-nowrap">{{ item.date }}</td>
                <td class="py-3.5 px-4 font-bold text-rose-800">
                  <span class="px-2.5 py-0.5 bg-rose-50 border border-rose-200/80 text-rose-700 rounded-full text-[10px] font-extrabold">
                    {{ item.shiftName }}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-mono whitespace-nowrap">
                  <div class="font-extrabold text-emerald-700">{{ item.clockIn }}</div>
                  <div v-if="item.clockIn !== '--:--:--'" class="text-[10px] font-bold text-emerald-800 bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-200 inline-block mt-0.5">
                    ðŸ“ {{ item.masukDistance }}
                  </div>
                </td>
                <td class="py-3.5 px-4 font-mono whitespace-nowrap">
                  <div class="font-extrabold text-amber-700">{{ item.clockOut }}</div>
                  <div v-if="item.clockOut !== '--:--:--'" class="text-[10px] font-bold text-amber-800 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200 inline-block mt-0.5">
                    ðŸ“ {{ item.pulangDistance }}
                  </div>
                </td>
              </tr>
              <tr v-if="filteredHistoryList.length === 0">
                <td colspan="5" class="py-6 text-center">
                  <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
                    <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
                    <span>Tidak ada riwayat presensi yang sesuai dengan filter.</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls Footer -->
        <div v-if="filteredHistoryList.length > 0" class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t border-slate-200 text-xs">
          <div class="text-slate-500 font-semibold text-center sm:text-left">
            Menampilkan <strong class="text-slate-900 font-mono font-bold">{{ ((currentPage - 1) * pageSize) + 1 }}</strong> - <strong class="text-slate-900 font-mono font-bold">{{ Math.min(currentPage * pageSize, filteredHistoryList.length) }}</strong> dari <strong class="text-slate-900 font-mono font-bold">{{ filteredHistoryList.length }}</strong> riwayat presensi
          </div>

          <div class="flex items-center gap-1.5">
            <button 
              @click="currentPage--" 
              :disabled="currentPage === 1"
              class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-xl border border-slate-300 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs"
            >
              <span class="material-symbols-outlined text-[16px]">chevron_left</span>
              <span>Sebelumnya</span>
            </button>

            <span class="px-3 py-1.5 bg-rose-50 border border-rose-200 text-rose-800 font-extrabold font-mono rounded-xl text-xs">
              {{ currentPage }} / {{ totalPages }}
            </span>

            <button 
              @click="currentPage++" 
              :disabled="currentPage >= totalPages"
              class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-xl border border-slate-300 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs"
            >
              <span>Berikutnya</span>
              <span class="material-symbols-outlined text-[16px]">chevron_right</span>
            </button>
          </div>
        </div>

      </div>

    </div>

    <!-- ========== MODAL KONFIRMASI KELUAR SESI ========== -->
    <Teleport to="body">
      <div v-if="showLogoutModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs animate-in fade-in duration-200 select-none">
        <div class="bg-white rounded-3xl border border-slate-200 shadow-2xl max-w-sm w-full p-6 text-center space-y-4 animate-in zoom-in-95 duration-150">
          <div class="w-14 h-14 rounded-full bg-rose-100 text-rose-600 flex items-center justify-center mx-auto shadow-inner">
            <span class="material-symbols-outlined text-3xl">logout</span>
          </div>
          
          <div class="space-y-1">
            <h3 class="text-base font-extrabold text-slate-900">Konfirmasi Keluar Sesi</h3>
            <p class="text-xs text-slate-500 leading-relaxed">
              Apakah Anda yakin ingin keluar dari akun <strong class="text-slate-800">{{ authStore.user?.name || 'Peserta Magang' }}</strong>? Anda harus memasukkan kredensial login kembali untuk masuk.
            </p>
          </div>

          <div class="grid grid-cols-2 gap-2.5 pt-2">
            <button
              @click="showLogoutModal = false"
              class="w-full py-2.5 px-4 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs rounded-xl transition-all cursor-pointer border-0"
            >
              Batal
            </button>
            <button
              @click="executeLogout"
              class="w-full py-2.5 px-4 bg-gradient-to-r from-rose-600 to-rose-700 hover:from-rose-500 hover:to-rose-600 text-white font-extrabold text-xs rounded-xl shadow-md transition-all cursor-pointer border-0 active:scale-95"
            >
              Ya, Keluar
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </InternLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import InternLayout from '@/layouts/InternLayout.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const schedulesData = ref(null)

const now = new Date()
const currentMonthStr = String(now.getMonth() + 1).padStart(2, '0')
const currentYearStr = String(now.getFullYear())

const selectedMonth = ref('ALL')
const selectedYear = ref('ALL')
const pageSize = ref(10)
const currentPage = ref(1)

const monthOptions = [
  { value: 'ALL', label: 'Semua Bulan' },
  { value: '01', label: 'Januari' },
  { value: '02', label: 'Februari' },
  { value: '03', label: 'Maret' },
  { value: '04', label: 'April' },
  { value: '05', label: 'Mei' },
  { value: '06', label: 'Juni' },
  { value: '07', label: 'Juli' },
  { value: '08', label: 'Agustus' },
  { value: '09', label: 'September' },
  { value: '10', label: 'Oktober' },
  { value: '11', label: 'November' },
  { value: '12', label: 'Desember' }
]

const yearOptions = [
  { value: 'ALL', label: 'Semua Tahun' },
  { value: '2024', label: '2024' },
  { value: '2025', label: '2025' },
  { value: '2026', label: '2026' },
  { value: '2027', label: '2027' }
]

const pageSizeOptions = [5, 10, 20, 50, 100]

watch([selectedMonth, selectedYear, pageSize], () => {
  currentPage.value = 1
})

const fetchSchedulesData = async () => {
  try {
    const res = await axios.get('/api/admin/schedules')
    schedulesData.value = res.data?.schedules
  } catch (e) {}
}

const getShiftNameForDate = (dateStr) => {
  let isoDate = dateStr
  if (dateStr.includes('-') && dateStr.split('-')[0].length === 2) {
    const [d, m, y] = dateStr.split('-')
    isoDate = `${y}-${m}-${d}`
  }

  const cleanUserNip = (authStore.user?.nip || '').replace(/\s+/g, '')
  const cleanUserName = (authStore.user?.name || '').toLowerCase()
  const teams = schedulesData.value?.teams || []

  let userTeamCode = ''
  for (const t of teams) {
    if (t.members && Array.isArray(t.members)) {
      const isMember = t.members.some((m) => {
        if (typeof m === 'string') {
          const cleanM = m.replace(/\s+/g, '')
          return (cleanUserNip && cleanM.includes(cleanUserNip)) || (cleanUserName && cleanM.toLowerCase().includes(cleanUserName))
        } else if (typeof m === 'object' && m !== null) {
          const mNip = (m.nip || m.Nip || m.NIP || '').replace(/\s+/g, '')
          const mName = (m.name || m.Name || '').toLowerCase()
          return (cleanUserNip && mNip.includes(cleanUserNip)) || (cleanUserName && mName.includes(cleanUserName))
        }
        return false
      })
      if (isMember) {
        userTeamCode = t.code || t.id || t.name
        break
      }
    }
  }

  if (schedulesData.value && Array.isArray(schedulesData.value.daysInMonth)) {
    const dayEntry = schedulesData.value.daysInMonth.find((d) => d.date === isoDate || d.dateStr === isoDate)
    if (dayEntry && userTeamCode) {
      if (dayEntry.shift1 === userTeamCode || dayEntry.shift1Team === userTeamCode) return 'Shift 1 (Pagi)'
      if (dayEntry.shift2 === userTeamCode || dayEntry.shift2Team === userTeamCode) {
        const mode = schedulesData.value?.shiftMode || 2
        return mode === 3 ? 'Shift 2 (Sore)' : 'Shift 2 (Malam)'
      }
      if (dayEntry.shift3 === userTeamCode || dayEntry.shift3Team === userTeamCode) return 'Shift 3 (Malam)'
      if (dayEntry.offTeams && Array.isArray(dayEntry.offTeams) && dayEntry.offTeams.includes(userTeamCode)) return 'OFF (Libur)'
    }
  }

  return 'Shift Siaga 112'
}

const parseTs = (rawTs) => {
  if (!rawTs) return { dateKey: '--', isoDate: '', month: '', year: '', timeFormatted: '--:--:--', hour: 0 }
  const clean = rawTs.replace('T', ' ').split('.')[0].replace('Z', '')
  const [datePart = '', timePart = ''] = clean.split(' ')

  let year = '', month = '', day = ''
  if (datePart.includes('-')) {
    const p = datePart.split('-')
    if (p[0].length === 4) {
      year = p[0]
      month = p[1]
      day = p[2]
    } else {
      day = p[0]
      month = p[1]
      year = p[2]
    }
  }

  const dateKey = (day && month && year) ? `${day.padStart(2, '0')}-${month.padStart(2, '0')}-${year}` : datePart
  const isoDate = (year && month && day) ? `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}` : datePart
  let cleanTime = timePart ? timePart.substring(0, 8) : '--:--:--'
  let timeFormatted = cleanTime !== '--:--:--' ? (cleanTime.includes('WITA') ? cleanTime : `${cleanTime} WITA`) : '--:--:--'

  let hour = 0
  if (cleanTime && cleanTime.includes(':')) {
    hour = parseInt(cleanTime.split(':')[0], 10)
  }

  return {
    dateKey,
    isoDate,
    month: month ? month.padStart(2, '0') : '',
    year: year || '',
    timeFormatted,
    hour
  }
}

const allPairedRows = computed(() => {
  if (!authStore.presensiHistory || !Array.isArray(authStore.presensiHistory) || authStore.presensiHistory.length === 0) {
    return []
  }

  // Sort history ascending by ID / timestamp so we pair chronologically
  const sortedHistory = [...authStore.presensiHistory].sort((a, b) => (a.id || 0) - (b.id || 0))
  const dateMap = new Map()

  sortedHistory.forEach((item) => {
    const parsed = parseTs(item.timestamp || '')
    const distanceFormatted = (item.distance_meters !== undefined && item.distance_meters !== null) ? `${item.distance_meters.toFixed(1)} Meter` : '0.8 Meter'

    if (item.type === 'MASUK') {
      const existing = dateMap.get(parsed.dateKey)
      if (!existing) {
        dateMap.set(parsed.dateKey, {
          id: item.id,
          date: parsed.dateKey,
          isoDate: parsed.isoDate,
          month: parsed.month,
          year: parsed.year,
          shiftName: getShiftNameForDate(parsed.dateKey),
          clockIn: parsed.timeFormatted,
          masukDistance: distanceFormatted,
          clockOut: '--:--:--',
          pulangDistance: '--'
        })
      } else {
        existing.clockIn = parsed.timeFormatted
        existing.masukDistance = distanceFormatted
      }
    } else if (item.type === 'PULANG') {
      // Check if this PULANG belongs to yesterday's night shift (scan PULANG between 00:00 - 12:00 WITA)
      let targetDateKey = parsed.dateKey
      if (parsed.hour < 12 && parsed.isoDate) {
        const d = new Date(parsed.isoDate)
        d.setDate(d.getDate() - 1)
        const prevDateKey = `${String(d.getDate()).padStart(2, '0')}-${String(d.getMonth() + 1).padStart(2, '0')}-${d.getFullYear()}`
        const prevEntry = dateMap.get(prevDateKey)
        if (prevEntry && prevEntry.clockOut === '--:--:--') {
          targetDateKey = prevDateKey
        }
      }

      let targetEntry = dateMap.get(targetDateKey)
      if (!targetEntry) {
        targetEntry = {
          id: item.id,
          date: targetDateKey,
          isoDate: parsed.isoDate,
          month: parsed.month,
          year: parsed.year,
          shiftName: getShiftNameForDate(targetDateKey),
          clockIn: '--:--:--',
          masukDistance: '--',
          clockOut: parsed.timeFormatted,
          pulangDistance: distanceFormatted
        }
        dateMap.set(targetDateKey, targetEntry)
      } else {
        targetEntry.clockOut = parsed.timeFormatted
        targetEntry.pulangDistance = distanceFormatted
      }
    }
  })

  // Convert map values to array and sort descending by date
  const result = Array.from(dateMap.values())
  result.sort((a, b) => {
    const keyA = a.isoDate || a.date
    const keyB = b.isoDate || b.date
    return keyB.localeCompare(keyA)
  })

  return result
})

const filteredHistoryList = computed(() => {
  return allPairedRows.value.filter((row) => {
    if (selectedMonth.value !== 'ALL' && row.month !== selectedMonth.value) {
      return false
    }
    if (selectedYear.value !== 'ALL' && row.year !== selectedYear.value) {
      return false
    }
    return true
  })
})

const totalPages = computed(() => {
  return Math.ceil(filteredHistoryList.value.length / pageSize.value) || 1
})

const paginatedHistoryList = computed(() => {
  const startIdx = (currentPage.value - 1) * pageSize.value
  const endIdx = startIdx + Number(pageSize.value)
  return filteredHistoryList.value.slice(startIdx, endIdx)
})

const fetchHistory = async () => {
  await authStore.fetchHistory()
}

const showLogoutModal = ref(false)

function handleLogout() {
  showLogoutModal.value = true
}

function executeLogout() {
  showLogoutModal.value = false
  authStore.logout()
  window.location.href = '/login'
}

onMounted(async () => {
  await authStore.fetchProfile()
  await fetchSchedulesData()
  await fetchHistory()
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>





