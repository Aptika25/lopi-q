<template>
  <InternLayout>
    <div class="space-y-6 select-none font-sans w-full pb-28 sm:pb-8">
      
      <!-- ===== PAGE HEADER ===== -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-[#F8BBD0]/60 pb-4">
        <div>
          <h2 class="font-display font-black text-[#1b1c1c] text-lg sm:text-xl flex items-center gap-2">
            <span class="material-symbols-outlined text-[#f06292] text-2xl">history</span>
            <span>Riwayat &amp; Log Presensi Magang</span>
          </h2>
          <p class="font-sans text-[#8a7176] mt-1 text-xs">Rekapitulasi lengkap pemindaian QR Code presensi Masuk &amp; Pulang tugas siaga peserta magang.</p>
        </div>
        <div class="flex items-center gap-2">
          <span class="px-3 py-1 bg-[#FCE4EC] text-[#ab2c5d] text-xs font-extrabold rounded-full border border-[#F8BBD0]">
            REKAP PRESENSI DIGITAL
          </span>
        </div>
      </div>

      <!-- ===== PROFILE SUMMARY CARD ===== -->
      <div class="bg-white rounded-3xl p-6 sm:p-7 border border-[#F8BBD0] shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-start gap-4">
          <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-[#ab2c5d] via-[#f06292] to-[#fec1d6] text-white font-display font-black text-2xl flex items-center justify-center shrink-0 shadow-md">
            {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'M' }}
          </div>
          <div class="flex-1 min-w-0 space-y-1">
            <h2 class="text-lg font-display font-black text-[#1b1c1c] leading-tight">
              {{ authStore.user?.name || 'Peserta Magang' }}
            </h2>
            
            <div class="text-xs text-[#ab2c5d] font-mono font-bold">
              NISN / NIM: {{ authStore.user?.nip || '-' }}
            </div>

            <div class="text-xs text-[#574146] font-bold uppercase tracking-wide">
              Jurusan: {{ authStore.user?.jabatan || 'Rekayasa Perangkat Lunak' }}
            </div>

            <div class="text-xs text-[#8a7176] font-medium pt-0.5">
              Asal Sekolah / Kampus: {{ authStore.user?.unit_kerja || 'Diskominfo Bulukumba' }}
            </div>
          </div>
        </div>

        <!-- Keluar Sesi Button (Mobile Only) -->
        <div class="sm:hidden w-full pt-3 border-t border-[#F8BBD0]/60">
          <button 
            @click="handleLogout"
            class="w-full py-2.5 px-4 bg-[#FFF5F8] hover:bg-[#FCE4EC] border border-[#F8BBD0] text-[#ab2c5d] font-bold text-xs rounded-full transition-all flex items-center justify-center gap-2 cursor-pointer shadow-2xs group"
          >
            <span>Keluar Sesi</span>
            <span class="material-symbols-outlined text-base group-hover:translate-x-0.5 transition-transform">logout</span>
          </button>
        </div>
      </div>

      <!-- ===== HISTORY CONTAINER ===== -->
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-[#F8BBD0] shadow-sm space-y-6">
        
        <!-- Header Controls & Refresh Button -->
        <div class="flex items-center justify-between pb-4 border-b border-[#F8BBD0]/60 gap-2">
          <div class="min-w-0 flex-1">
            <h3 class="text-sm sm:text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">description</span>
              <span class="truncate">Daftar Log Kehadiran</span>
            </h3>
            <p class="text-xs text-[#8a7176] mt-0.5 hidden sm:block">Log presensi otomatis yang terverifikasi Geofence Posko Siaga 112.</p>
          </div>

          <button 
            @click="fetchHistory"
            class="px-3.5 py-2 bg-[#FFF5F8] hover:bg-[#FCE4EC] text-[#ab2c5d] font-bold text-xs rounded-xl border border-[#F8BBD0] flex items-center gap-1.5 cursor-pointer shrink-0 transition-colors"
          >
            <span class="material-symbols-outlined text-base">refresh</span>
            <span>Refresh</span>
          </button>
        </div>

        <!-- Filter & Sorting Toolbar -->
        <div class="bg-[#FFF5F8] p-4 rounded-2xl border border-[#F8BBD0]/60 text-xs space-y-3 sm:space-y-0 sm:flex sm:items-center sm:justify-between sm:gap-3">
          <!-- Month & Year Filters -->
          <div class="grid grid-cols-2 sm:flex sm:items-center gap-2 flex-1">
            <select v-model="selectedMonth" class="w-full sm:w-auto px-3 py-2 bg-white border border-[#ddbfc5] rounded-xl font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] text-xs">
              <option v-for="m in monthOptions" :key="m.value" :value="m.value">{{ m.label }}</option>
            </select>

            <select v-model="selectedYear" class="w-full sm:w-auto px-3 py-2 bg-white border border-[#ddbfc5] rounded-xl font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] text-xs">
              <option v-for="y in yearOptions" :key="y.value" :value="y.value">{{ y.label }}</option>
            </select>
          </div>

          <!-- Page Size & Reset Filter -->
          <div class="flex items-center justify-between sm:justify-end gap-3 pt-2 sm:pt-0 border-t sm:border-0 border-[#F8BBD0]/60">
            <div class="flex items-center gap-1.5 font-bold text-[#574146]">
              <span class="text-[11px] text-[#8a7176]">Tampilkan:</span>
              <select v-model="pageSize" class="px-2.5 py-1.5 bg-white border border-[#ddbfc5] rounded-xl font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] text-xs">
                <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }} Data</option>
              </select>
            </div>

            <button 
              v-if="selectedMonth !== 'ALL' || selectedYear !== 'ALL'"
              @click="selectedMonth = 'ALL'; selectedYear = 'ALL'"
              class="px-3 py-1.5 bg-[#ab2c5d] hover:bg-[#881b47] text-white font-extrabold rounded-xl border-0 text-xs cursor-pointer whitespace-nowrap shadow-2xs"
            >
              Reset Filter
            </button>
          </div>
        </div>

        <!-- Mobile Card View -->
        <div class="block sm:hidden space-y-3">
          <div 
            v-for="item in paginatedHistoryList" 
            :key="item.id"
            class="bg-white rounded-2xl p-4 border border-[#F8BBD0]/80 shadow-2xs space-y-3"
          >
            <!-- Card Header -->
            <div class="flex items-center justify-between border-b border-[#F8BBD0]/50 pb-2.5">
              <div class="flex items-center gap-1.5 font-mono font-bold text-[#1b1c1c] text-xs">
                <span class="material-symbols-outlined text-base text-[#f06292]">calendar_today</span>
                <span>{{ item.date }}</span>
              </div>
              <span class="px-2.5 py-0.5 bg-[#FCE4EC] border border-[#F8BBD0] text-[#ab2c5d] rounded-full text-[10px] font-extrabold">
                {{ item.shiftName }}
              </span>
            </div>

            <!-- Card Body: Clock In & Clock Out Grid -->
            <div class="grid grid-cols-2 gap-2.5 text-xs">
              <!-- Jam Masuk Box -->
              <div class="bg-emerald-50/70 p-3 rounded-xl border border-emerald-200/80 flex flex-col justify-between space-y-2 min-h-[85px]">
                <div class="text-[10px] font-extrabold text-emerald-800 uppercase tracking-wider flex items-center gap-1">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Jam Masuk
                </div>
                <div class="font-mono text-xs font-black text-[#1b1c1c]">
                  {{ item.clockIn }}
                </div>
                <div>
                  <span v-if="item.clockIn !== '--:--:--'" class="text-[10px] font-bold text-emerald-800 bg-white/90 px-2 py-0.5 rounded border border-emerald-200 inline-block shadow-2xs">
                    ✓ Terverifikasi
                  </span>
                  <span v-else class="text-[10px] font-semibold text-slate-400 italic">--</span>
                </div>
              </div>

              <!-- Jam Pulang Box -->
              <div 
                class="p-3 rounded-xl border flex flex-col justify-between space-y-2 min-h-[85px]"
                :class="item.clockOut !== '--:--:--' ? 'bg-amber-50/70 border-amber-200/80' : 'bg-slate-50 border-slate-200'"
              >
                <div 
                  class="text-[10px] font-extrabold uppercase tracking-wider flex items-center gap-1"
                  :class="item.clockOut !== '--:--:--' ? 'text-amber-800' : 'text-slate-400'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="item.clockOut !== '--:--:--' ? 'bg-amber-500' : 'bg-slate-300'"></span> Jam Pulang
                </div>
                <div class="font-mono text-xs font-black" :class="item.clockOut !== '--:--:--' ? 'text-[#1b1c1c]' : 'text-slate-400'">
                  {{ item.clockOut }}
                </div>
                <div>
                  <span v-if="item.clockOut !== '--:--:--'" class="text-[10px] font-bold text-amber-800 bg-white/90 px-2 py-0.5 rounded border border-amber-200 inline-block shadow-2xs">
                    ✓ Terverifikasi
                  </span>
                  <span v-else class="text-[10px] font-medium text-slate-400 bg-white/70 px-2 py-0.5 rounded border border-slate-200 inline-block">
                    ⏳ Belum Scan
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="filteredHistoryList.length === 0" class="py-8 text-center bg-[#FFF5F8] rounded-2xl border border-dashed border-[#F8BBD0]">
            <div class="inline-flex items-center justify-center gap-1.5 text-xs text-[#8a7176] font-medium">
              <span class="material-symbols-outlined text-[#f06292]">info</span>
              <span>Tidak ada riwayat presensi yang sesuai dengan filter.</span>
            </div>
          </div>
        </div>

        <!-- Desktop History Table -->
        <div class="hidden sm:block overflow-x-auto rounded-2xl border border-[#F8BBD0] shadow-xs">
          <table class="w-full text-left text-xs text-[#1b1c1c]">
            <thead class="bg-[#FFF5F8] text-[#ab2c5d] font-bold uppercase tracking-wider border-b border-[#F8BBD0]">
              <tr>
                <th class="py-3.5 px-4">Tanggal Siaga</th>
                <th class="py-3.5 px-4">Jadwal Shift</th>
                <th class="py-3.5 px-4">Jam Masuk</th>
                <th class="py-3.5 px-4">Jam Pulang</th>
                <th class="py-3.5 px-4">Status Geofence</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#F8BBD0]/40 bg-white">
              <tr v-for="item in paginatedHistoryList" :key="item.id" class="hover:bg-[#FFF5F8]/50 transition-colors">
                <td class="py-3.5 px-4 font-mono font-bold text-[#1b1c1c] whitespace-nowrap">{{ item.date }}</td>
                <td class="py-3.5 px-4">
                  <span class="px-2.5 py-0.5 bg-[#FCE4EC] border border-[#F8BBD0] text-[#ab2c5d] rounded-full text-[10px] font-extrabold">
                    {{ item.shiftName }}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-mono whitespace-nowrap font-extrabold text-emerald-700">
                  {{ item.clockIn }}
                </td>
                <td class="py-3.5 px-4 font-mono whitespace-nowrap font-extrabold text-amber-700">
                  {{ item.clockOut }}
                </td>
                <td class="py-3.5 px-4 whitespace-nowrap">
                  <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border bg-emerald-50 text-emerald-700 border-emerald-200">
                    Posko 112
                  </span>
                </td>
              </tr>
              <tr v-if="filteredHistoryList.length === 0">
                <td colspan="5" class="py-8 text-center">
                  <div class="inline-flex items-center justify-center gap-1.5 text-xs text-[#8a7176] font-medium">
                    <span class="material-symbols-outlined text-[#f06292]">info</span>
                    <span>Tidak ada riwayat presensi yang sesuai dengan filter.</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls Footer -->
        <div v-if="filteredHistoryList.length > 0" class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t border-[#F8BBD0]/60 text-xs">
          <div class="text-[#574146] font-medium text-center sm:text-left">
            Menampilkan <strong class="text-[#1b1c1c] font-mono font-bold">{{ ((currentPage - 1) * pageSize) + 1 }}</strong> - <strong class="text-[#1b1c1c] font-mono font-bold">{{ Math.min(currentPage * pageSize, filteredHistoryList.length) }}</strong> dari <strong class="text-[#1b1c1c] font-mono font-bold">{{ filteredHistoryList.length }}</strong> riwayat presensi
          </div>

          <div class="flex items-center gap-2">
            <button 
              @click="currentPage--" 
              :disabled="currentPage === 1"
              class="px-3 py-1.5 bg-[#FFF5F8] hover:bg-[#FCE4EC] text-[#ab2c5d] font-bold rounded-xl border border-[#F8BBD0] cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs"
            >
              <span class="material-symbols-outlined text-base">chevron_left</span>
              <span>Sebelumnya</span>
            </button>

            <span class="px-3 py-1.5 bg-[#ab2c5d] text-white font-extrabold font-mono rounded-xl text-xs shadow-2xs">
              {{ currentPage }} / {{ totalPages }}
            </span>

            <button 
              @click="currentPage++" 
              :disabled="currentPage >= totalPages"
              class="px-3 py-1.5 bg-[#FFF5F8] hover:bg-[#FCE4EC] text-[#ab2c5d] font-bold rounded-xl border border-[#F8BBD0] cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs"
            >
              <span>Berikutnya</span>
              <span class="material-symbols-outlined text-base">chevron_right</span>
            </button>
          </div>
        </div>

      </div>

    </div>

    <!-- ===== MODAL KONFIRMASI KELUAR SESI ===== -->
    <Teleport to="body">
      <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs select-none">
        <div class="bg-white rounded-3xl border border-[#F8BBD0] shadow-2xl max-w-sm w-full p-6 text-center space-y-4">
          <div class="w-14 h-14 rounded-full bg-[#FCE4EC] text-[#ab2c5d] flex items-center justify-center mx-auto shadow-inner">
            <span class="material-symbols-outlined text-3xl">logout</span>
          </div>
          
          <div class="space-y-1">
            <h3 class="text-base font-extrabold text-[#1b1c1c]">Konfirmasi Keluar Sesi</h3>
            <p class="text-xs text-[#574146] leading-relaxed">
              Apakah Anda yakin ingin keluar dari akun <strong class="text-[#1b1c1c]">{{ authStore.user?.name || 'Peserta Magang' }}</strong>?
            </p>
          </div>

          <div class="grid grid-cols-2 gap-2.5 pt-2">
            <button
              @click="showLogoutModal = false"
              class="w-full py-2.5 px-4 bg-slate-100 hover:bg-slate-200 text-[#574146] font-bold text-xs rounded-xl transition-all cursor-pointer border-0"
            >
              Batal
            </button>
            <button
              @click="executeLogout"
              class="w-full py-2.5 px-4 bg-[#ab2c5d] hover:bg-[#881b47] text-white font-extrabold text-xs rounded-xl shadow-md transition-all cursor-pointer border-0 active:scale-95"
            >
              Ya, Keluar
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </InternLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import InternLayout from '@/layouts/InternLayout.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const schedulesData = ref<any>(null)

const now = new Date()

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

const getShiftNameForDate = (dateStr: string) => {
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
      const isMember = t.members.some((m: any) => {
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
    const dayEntry = schedulesData.value.daysInMonth.find((d: any) => d.date === isoDate || d.dateStr === isoDate)
    if (dayEntry && userTeamCode) {
      if (dayEntry.shift1 === userTeamCode || dayEntry.shift1Team === userTeamCode) return 'Shift Pagi'
      if (dayEntry.shift2 === userTeamCode || dayEntry.shift2Team === userTeamCode) {
        const mode = schedulesData.value?.shiftMode || 2
        return mode === 3 ? 'Shift Sore' : 'Shift Malam'
      }
      if (dayEntry.shift3 === userTeamCode || dayEntry.shift3Team === userTeamCode) return 'Shift Malam'
      if (dayEntry.offTeams && Array.isArray(dayEntry.offTeams) && dayEntry.offTeams.includes(userTeamCode)) return 'OFF (Libur)'
    }
  }

  return 'Shift Siaga 112'
}

const parseTs = (rawTs: string) => {
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

  const sortedHistory = [...authStore.presensiHistory].sort((a: any, b: any) => (a.id || 0) - (b.id || 0))
  const dateMap = new Map()

  sortedHistory.forEach((item: any) => {
    const parsed = parseTs(item.timestamp || '')

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
          clockOut: '--:--:--'
        })
      } else {
        existing.clockIn = parsed.timeFormatted
      }
    } else if (item.type === 'PULANG') {
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
          clockOut: parsed.timeFormatted
        }
        dateMap.set(targetDateKey, targetEntry)
      } else {
        targetEntry.clockOut = parsed.timeFormatted
      }
    }
  })

  const result = Array.from(dateMap.values())
  result.sort((a: any, b: any) => {
    const keyA = a.isoDate || a.date
    const keyB = b.isoDate || b.date
    return keyB.localeCompare(keyA)
  })

  return result
})

const filteredHistoryList = computed(() => {
  return allPairedRows.value.filter((row: any) => {
    if (selectedMonth.value !== 'ALL' && row.month !== selectedMonth.value) return false
    if (selectedYear.value !== 'ALL' && row.year !== selectedYear.value) return false
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
  router.push('/login')
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
