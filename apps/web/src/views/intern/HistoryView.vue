<template>
  <InternLayout>
    <div class="space-y-6 select-none font-sans w-full max-w-[1200px] mx-auto pb-28 sm:pb-8">
      
      <!-- ===== 1. PROFILE CARD SECTION ===== -->
      <section class="mt-2">
        <div class="bg-white border border-[#ddbfc5] rounded-2xl p-5 shadow-sm flex flex-col gap-4">
          <div class="flex items-center gap-4">
            <!-- Avatar Box -->
            <div class="w-16 h-16 rounded-xl bg-[#f06292] flex items-center justify-center shadow-sm shrink-0">
              <span class="text-white font-bold text-2xl">
                {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'A' }}
              </span>
            </div>

            <!-- Profile Info -->
            <div class="flex flex-col justify-center min-w-0">
              <h3 class="text-lg font-bold text-[#1b1c1c] leading-tight truncate">
                {{ authStore.user?.name || 'adhe anisa' }}
              </h3>
              <p class="text-[10px] font-bold text-[#ab2c5d] uppercase tracking-wider mt-0.5">
                NISN. {{ authStore.user?.nip || '0091755987' }}
              </p>
              <p class="text-xs text-[#574146] font-semibold mt-0.5 truncate">
                {{ authStore.user?.unit_kerja || 'SMKS TI BULUKUMBA' }}
              </p>
            </div>
          </div>

          <div class="h-px w-full bg-[#ddbfc5]/50"></div>

          <!-- Keluar Sesi Button -->
          <button 
            @click="handleLogout"
            class="w-full py-2.5 px-4 border border-[#f06292] rounded-full flex items-center justify-center gap-2 text-[#ab2c5d] font-bold text-xs hover:bg-[#f06292]/10 transition-colors active:scale-[0.98] cursor-pointer bg-transparent"
          >
            <span>Keluar Sesi</span>
            <span class="material-symbols-outlined text-lg">logout</span>
          </button>
        </div>
      </section>

      <!-- ===== 2. MONTH FILTER SECTION ===== -->
      <section class="py-2 flex items-center justify-between gap-3 bg-white p-4 rounded-2xl border border-[#ddbfc5] shadow-xs">
        <button 
          @click="prevMonth"
          class="w-10 h-10 flex items-center justify-center rounded-full bg-[#f5f3f3] hover:bg-[#eae8e7] transition-colors active:scale-95 text-[#574146] border-0 cursor-pointer"
        >
          <span class="material-symbols-outlined">chevron_left</span>
        </button>

        <div class="flex items-center gap-2">
          <select 
            v-model="selectedMonth" 
            class="px-3 py-1.5 bg-[#f5f3f3] border border-[#ddbfc5] rounded-xl font-bold text-[#1b1c1c] text-xs focus:outline-none focus:border-[#f06292]"
          >
            <option v-for="m in monthOptions" :key="m.value" :value="m.value">{{ m.label }}</option>
          </select>

          <select 
            v-model="selectedYear" 
            class="px-3 py-1.5 bg-[#f5f3f3] border border-[#ddbfc5] rounded-xl font-bold text-[#1b1c1c] text-xs focus:outline-none focus:border-[#f06292]"
          >
            <option v-for="y in yearOptions" :key="y.value" :value="y.value">{{ y.label }}</option>
          </select>
        </div>

        <button 
          @click="nextMonth"
          class="w-10 h-10 flex items-center justify-center rounded-full bg-[#f5f3f3] hover:bg-[#eae8e7] transition-colors active:scale-95 text-[#574146] border-0 cursor-pointer"
        >
          <span class="material-symbols-outlined">chevron_right</span>
        </button>
      </section>

      <!-- ===== 3. HISTORY LIST SECTION ===== -->
      <section class="flex flex-col gap-3">
        <article 
          v-for="item in paginatedHistoryList" 
          :key="item.id"
          class="bg-white border border-[#ddbfc5] rounded-2xl p-4 flex flex-col gap-3 hover:bg-[#ffd9e4]/10 transition-colors shadow-sm"
        >
          <!-- Item Header: Date & Badge -->
          <div class="flex justify-between items-start">
            <h3 class="text-xs font-bold text-[#574146] flex items-center gap-1.5">
              <span class="material-symbols-outlined text-sm text-[#f06292]">calendar_today</span>
              <span>{{ item.date }}</span>
            </h3>

            <!-- Status Badge (Hadir, Terlambat, Izin/Cuti) -->
            <span 
              class="text-[10px] font-extrabold uppercase px-3 py-1 rounded-full tracking-wider border"
              :class="{
                'bg-[#E8F5E9] text-[#1B5E20] border-emerald-200': item.statusBadge === 'HADIR',
                'bg-[#FFF8E1] text-[#FF8F00] border-amber-200': item.statusBadge === 'TERLAMBAT',
                'bg-[#FCE4EC] text-[#F06292] border-rose-200': item.statusBadge === 'IZIN' || item.statusBadge === 'SAKIT'
              }"
            >
              {{ item.statusBadge }}
            </span>
          </div>

          <!-- Item Body: Clock In & Clock Out Grid -->
          <div class="flex items-center gap-6 mt-1">
            <!-- Masuk Box -->
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d] text-xl">login</span>
              <div class="flex flex-col">
                <span class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Masuk</span>
                <span class="text-base font-black font-mono text-[#1b1c1c] leading-tight">{{ item.clockIn }}</span>
              </div>
            </div>

            <!-- Divider Line -->
            <div class="h-8 w-px bg-[#ddbfc5]/60"></div>

            <!-- Keluar Box -->
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined text-[#8a7176] text-xl">logout</span>
              <div class="flex flex-col">
                <span class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Keluar</span>
                <span class="text-base font-black font-mono leading-tight" :class="item.clockOut !== '--:--' ? 'text-[#1b1c1c]' : 'text-slate-400'">
                  {{ item.clockOut }}
                </span>
              </div>
            </div>
          </div>
        </article>

        <!-- Empty State -->
        <div v-if="filteredHistoryList.length === 0" class="py-10 text-center bg-white rounded-2xl border border-dashed border-[#ddbfc5]">
          <div class="inline-flex items-center justify-center gap-2 text-xs text-[#8a7176] font-medium">
            <span class="material-symbols-outlined text-[#f06292]">info</span>
            <span>Tidak ada riwayat presensi yang ditemukan untuk periode ini.</span>
          </div>
        </div>

        <!-- Pagination Footer -->
        <div v-if="filteredHistoryList.length > 0" class="flex items-center justify-between pt-3 text-xs">
          <span class="text-[#574146] font-medium">
            Data {{ ((currentPage - 1) * pageSize) + 1 }} - {{ Math.min(currentPage * pageSize, filteredHistoryList.length) }} dari {{ filteredHistoryList.length }}
          </span>

          <div class="flex items-center gap-1.5">
            <button 
              @click="currentPage--" 
              :disabled="currentPage === 1"
              class="px-3 py-1.5 bg-[#f5f3f3] hover:bg-[#eae8e7] text-[#ab2c5d] font-bold rounded-xl border border-[#ddbfc5] cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed text-xs"
            >
              Prev
            </button>

            <span class="px-3 py-1.5 bg-[#ab2c5d] text-white font-mono font-bold rounded-xl text-xs">
              {{ currentPage }} / {{ totalPages }}
            </span>

            <button 
              @click="currentPage++" 
              :disabled="currentPage >= totalPages"
              class="px-3 py-1.5 bg-[#f5f3f3] hover:bg-[#eae8e7] text-[#ab2c5d] font-bold rounded-xl border border-[#ddbfc5] cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed text-xs"
            >
              Next
            </button>
          </div>
        </div>

      </section>

    </div>

    <!-- ===== MODAL KONFIRMASI KELUAR SESI ===== -->
    <Teleport to="body">
      <div v-if="showLogoutModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs select-none">
        <div class="bg-white rounded-3xl border border-[#ddbfc5] shadow-2xl max-w-sm w-full p-6 text-center space-y-4">
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
              class="w-full py-2.5 px-4 bg-[#f5f3f3] hover:bg-slate-200 text-[#574146] font-bold text-xs rounded-xl transition-all cursor-pointer border-0"
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
import InternLayout from '@/layouts/InternLayout.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const now = new Date()
const currentMonthVal = String(now.getMonth() + 1).padStart(2, '0')
const currentYearVal = String(now.getFullYear())

const selectedMonth = ref(currentMonthVal)
const selectedYear = ref(currentYearVal)
const pageSize = ref(10)
const currentPage = ref(1)
const showLogoutModal = ref(false)

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

watch([selectedMonth, selectedYear], () => {
  currentPage.value = 1
})

const prevMonth = () => {
  if (selectedMonth.value === 'ALL') {
    selectedMonth.value = '01'
    return
  }
  let m = parseInt(selectedMonth.value, 10) - 1
  if (m < 1) {
    m = 12
    let y = parseInt(selectedYear.value, 10) - 1
    selectedYear.value = String(y)
  }
  selectedMonth.value = String(m).padStart(2, '0')
}

const nextMonth = () => {
  if (selectedMonth.value === 'ALL') {
    selectedMonth.value = '12'
    return
  }
  let m = parseInt(selectedMonth.value, 10) + 1
  if (m > 12) {
    m = 1
    let y = parseInt(selectedYear.value, 10) + 1
    selectedYear.value = String(y)
  }
  selectedMonth.value = String(m).padStart(2, '0')
}

const parseTs = (rawTs: string) => {
  if (!rawTs) return { dateKey: '--', isoDate: '', month: '', year: '', timeFormatted: '--:--', hour: 0, min: 0 }
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

  const dateKey = (day && month && year) ? `${day.padStart(2, '0')}/${month.padStart(2, '0')}/${year}` : datePart
  const isoDate = (year && month && day) ? `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}` : datePart
  let cleanTime = timePart ? timePart.substring(0, 5) : '--:--'

  let hour = 0
  let min = 0
  if (cleanTime && cleanTime.includes(':')) {
    const parts = cleanTime.split(':')
    hour = parseInt(parts[0], 10)
    min = parseInt(parts[1], 10)
  }

  return {
    dateKey,
    isoDate,
    month: month ? month.padStart(2, '0') : '',
    year: year || '',
    timeFormatted: cleanTime,
    hour,
    min
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
      const isLate = parsed.hour > 8 || (parsed.hour === 8 && parsed.min > 15)
      const status = isLate ? 'TERLAMBAT' : 'HADIR'

      if (!existing) {
        dateMap.set(parsed.dateKey, {
          id: item.id,
          date: parsed.dateKey,
          isoDate: parsed.isoDate,
          month: parsed.month,
          year: parsed.year,
          statusBadge: status,
          clockIn: parsed.timeFormatted,
          clockOut: '--:--'
        })
      } else {
        existing.clockIn = parsed.timeFormatted
        existing.statusBadge = status
      }
    } else if (item.type === 'PULANG') {
      let targetDateKey = parsed.dateKey
      let targetEntry = dateMap.get(targetDateKey)
      if (!targetEntry) {
        targetEntry = {
          id: item.id,
          date: targetDateKey,
          isoDate: parsed.isoDate,
          month: parsed.month,
          year: parsed.year,
          statusBadge: 'HADIR',
          clockIn: '--:--',
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
  await authStore.fetchHistory()
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
