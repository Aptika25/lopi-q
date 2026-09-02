<template>
  <InternLayout>
    <div class="select-none font-sans text-[#1b1c1c] pb-safe">
      
      <!-- Main Container -->
      <main class="px-4 sm:px-6 py-6 flex flex-col gap-6 max-w-2xl mx-auto">
        
        <!-- ===== 1. HEADER SECTION ===== -->
        <header class="flex justify-between items-center pt-2 pb-4 border-b border-[#F8BBD0]/60">
          <div>
            <h1 class="text-2xl sm:text-3xl font-bold text-[#1b1c1c] leading-tight">
              Halo, {{ userFirstName }}!
            </h1>
            <div class="flex items-center gap-2 mt-1 text-xs text-[#574146]">
              <span class="font-mono font-bold text-[#ab2c5d]">NISN/NIM: {{ authStore.user?.nip || '0091755987' }}</span>
              <span>•</span>
              <span class="truncate max-w-[180px] sm:max-w-none">{{ authStore.user?.unit_kerja || 'TKJ' }}</span>
            </div>
          </div>

          <!-- User Profile Avatar Circle -->
          <div class="flex flex-col items-center gap-1">
            <div class="h-14 w-14 sm:h-16 sm:w-16 rounded-full border-2 border-[#f06292] p-0.5 bg-[#fbf9f8] shadow-sm relative">
              <div class="w-full h-full object-cover rounded-full bg-gradient-to-tr from-[#ab2c5d] via-[#f06292] to-[#fec1d6] text-white font-black text-xl sm:text-2xl flex items-center justify-center shadow-inner">
                {{ userInitial }}
              </div>
              <span 
                class="absolute bottom-0 right-0 w-4 h-4 rounded-full border-2 border-white"
                :class="authStore.user?.is_active ? 'bg-emerald-500' : 'bg-slate-400'"
                :title="authStore.user?.is_active ? 'Akun Aktif' : 'Non-aktif'"
              ></span>
            </div>
            <span class="text-[10px] font-extrabold uppercase px-2 py-0.5 bg-[#fec1d6] text-[#65394b] rounded-full">
              INTERN
            </span>
          </div>
        </header>

        <!-- ===== 2. STATUS CARD (TODAY) ===== -->
        <section class="bg-white rounded-3xl p-5 sm:p-6 relative overflow-hidden border border-[#F8BBD0] shadow-sm space-y-4">
          <div class="absolute -right-12 -top-12 w-36 h-36 bg-[#f06292]/10 rounded-full blur-2xl pointer-events-none"></div>

          <div class="flex items-center justify-between relative z-10">
            <h2 class="text-lg font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">today</span>
              <span>Status Hari Ini</span>
            </h2>
            <span 
              class="px-3 py-1 rounded-full text-xs font-extrabold uppercase tracking-wider shadow-2xs border"
              :class="authStore.todayStatus?.is_masuk ? 'bg-emerald-50 text-emerald-800 border-emerald-200' : 'bg-[#FFF8E1] text-[#854d0e] border-[#fef08a]'"
            >
              {{ authStore.todayStatus?.is_masuk ? (authStore.todayStatus?.is_pulang ? 'Selesai Tugas' : 'Hadir (Siaga)') : 'Belum Absen' }}
            </span>
          </div>

          <!-- Check In & Check Out Grid -->
          <div class="grid grid-cols-2 gap-4 relative z-10">
            <!-- Check In Box -->
            <div class="bg-[#FFF5F8] p-4 rounded-2xl border border-[#F8BBD0]/60 space-y-1">
              <div class="flex items-center gap-2 text-[#8a7176] text-xs font-bold uppercase tracking-wider">
                <span class="material-symbols-outlined text-lg text-[#f06292]">login</span>
                <span>Check In</span>
              </div>
              <p class="text-2xl sm:text-3xl font-black text-[#ab2c5d] font-mono">
                {{ checkInTimeDisplay }}
              </p>
            </div>

            <!-- Check Out Box -->
            <div class="bg-[#FFF5F8] p-4 rounded-2xl border border-[#F8BBD0]/60 space-y-1">
              <div class="flex items-center gap-2 text-[#8a7176] text-xs font-bold uppercase tracking-wider">
                <span class="material-symbols-outlined text-lg text-[#8a7176]">logout</span>
                <span>Check Out</span>
              </div>
              <p 
                class="text-2xl sm:text-3xl font-black font-mono"
                :class="authStore.todayStatus?.is_pulang ? 'text-[#ab2c5d]' : 'text-[#8a7176]/50'"
              >
                {{ checkOutTimeDisplay }}
              </p>
            </div>
          </div>
        </section>

        <!-- ===== 3. WEEKLY PROGRESS (PROGRES MINGGUAN) ===== -->
        <section class="space-y-2">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">bar_chart</span>
              <span>Progres Mingguan Presensi</span>
            </h3>
            <span class="text-xs text-[#8a7176] font-medium">Minggu Ini</span>
          </div>

          <div class="bg-white rounded-3xl p-5 border border-[#F8BBD0] shadow-sm flex items-end justify-between h-36 gap-3">
            <!-- Monday (S) -->
            <div class="flex-1 bg-[#f06292] rounded-t-xl h-[85%] flex flex-col justify-end items-center pb-2 transition-all hover:opacity-90 shadow-2xs">
              <span class="font-bold text-[11px] text-white">S</span>
            </div>
            <!-- Tuesday (S) -->
            <div class="flex-1 bg-[#f06292] rounded-t-xl h-[100%] flex flex-col justify-end items-center pb-2 transition-all hover:opacity-90 shadow-2xs">
              <span class="font-bold text-[11px] text-white">S</span>
            </div>
            <!-- Wednesday (R) -->
            <div class="flex-1 bg-[#f06292] rounded-t-xl h-[90%] flex flex-col justify-end items-center pb-2 transition-all hover:opacity-90 shadow-2xs">
              <span class="font-bold text-[11px] text-white">R</span>
            </div>
            <!-- Thursday (K) -->
            <div class="flex-1 bg-[#f06292] rounded-t-xl h-[75%] flex flex-col justify-end items-center pb-2 transition-all hover:opacity-90 shadow-2xs">
              <span class="font-bold text-[11px] text-white">K</span>
            </div>
            <!-- Friday (J) -->
            <div class="flex-1 bg-[#FFF5F8] rounded-t-xl h-[40%] flex flex-col justify-end items-center pb-2 transition-all border border-dashed border-[#F8BBD0]">
              <span class="font-bold text-[11px] text-[#8a7176]">J</span>
            </div>
          </div>
        </section>

        <!-- ===== 4. JURNAL KEGIATAN HARIAN ===== -->
        <section class="space-y-2">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">menu_book</span>
              <span>Jurnal Kegiatan Harian</span>
            </h3>
            <button 
              @click="openAddNoteModal"
              class="text-[#ab2c5d] hover:bg-[#FCE4EC] p-1.5 rounded-full transition-colors flex items-center justify-center border-0 bg-transparent cursor-pointer"
              title="Tambah Catatan Jurnal"
            >
              <span class="material-symbols-outlined text-xl">edit_note</span>
            </button>
          </div>

          <div class="bg-white rounded-3xl border border-[#F8BBD0] shadow-sm overflow-hidden">
            <div class="overflow-y-auto max-h-64">
              <table class="w-full text-left border-collapse">
                <thead class="bg-[#FFF5F8] sticky top-0 z-10 border-b border-[#F8BBD0]">
                  <tr>
                    <th class="p-3 text-[11px] font-bold text-[#ab2c5d] uppercase tracking-wider">Nama</th>
                    <th class="p-3 text-[11px] font-bold text-[#ab2c5d] uppercase tracking-wider">Waktu</th>
                    <th class="p-3 text-[11px] font-bold text-[#ab2c5d] uppercase tracking-wider">Kegiatan</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-[#F8BBD0]/40 text-xs">
                  <tr v-for="(item, idx) in activityJournalList" :key="idx" class="hover:bg-[#FFF5F8]/50 transition-colors">
                    <td class="p-3 font-bold text-[#1b1c1c]">{{ item.name }}</td>
                    <td class="p-3 font-mono text-[#574146]">{{ item.time }}</td>
                    <td class="p-3 text-[#1b1c1c] font-medium">{{ item.activity }}</td>
                  </tr>
                  <tr v-if="activityJournalList.length === 0">
                    <td colspan="3" class="p-6 text-center text-[#8a7176] text-xs">
                      Belum ada catatan jurnal kegiatan harian recorded.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </section>

      </main>
    </div>

    <!-- ===== MODAL TAMBAH CATATAN JURNAL ===== -->
    <transition name="fade">
      <div v-if="showNoteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="showNoteModal = false">
        <div class="bg-white rounded-3xl max-w-md w-full p-6 border border-[#F8BBD0] shadow-2xl space-y-4">
          <div class="flex items-center justify-between border-b border-[#F8BBD0]/60 pb-3">
            <h3 class="font-bold text-base text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">edit_note</span>
              <span>Tambah Catatan Jurnal Harian</span>
            </h3>
            <button @click="showNoteModal = false" class="text-slate-400 hover:text-slate-600 border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <form @submit.prevent="addJournalNote" class="space-y-4 text-xs">
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Uraian Kegiatan / Task <span class="text-rose-500">*</span></label>
              <input v-model="newJournalNote" type="text" required placeholder="Contoh: Input laporan harian & kirim rekap presensi" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-xl focus:outline-none focus:border-[#f06292]" />
            </div>

            <div class="flex justify-end gap-3 pt-2">
              <button type="button" @click="showNoteModal = false" class="px-4 py-2 border border-[#ddbfc5] rounded-xl font-bold text-[#574146] bg-white cursor-pointer hover:bg-slate-50">Batal</button>
              <button type="submit" class="px-5 py-2 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-xl font-bold border-0 cursor-pointer shadow-xs">
                Simpan Catatan
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
import { useAuthStore } from '@/stores/auth'
import InternLayout from '@/layouts/InternLayout.vue'

const authStore = useAuthStore()

// State
const showNoteModal = ref(false)
const newJournalNote = ref('')

const defaultSampleJournal = [
  { name: 'Andi', time: '13:15', activity: 'Internal Sync' },
  { name: 'Rina', time: '15:00', activity: 'Quality Check' },
  { name: 'Hikma', time: '16:30', activity: 'Daily Report' }
]

// Dynamic Jurnal Activity List persisted in localStorage
const activityJournalList = ref<any[]>([])

const loadJournalNotes = () => {
  try {
    const key = `lopiq_journal_notes_${authStore.user?.nip || 'default'}`
    const stored = localStorage.getItem(key)
    if (stored) {
      const parsed = JSON.parse(stored)
      if (Array.isArray(parsed) && parsed.length > 0) {
        activityJournalList.value = parsed
        return
      }
    }
  } catch (e) {}
  
  // Fallback to default initial list if empty
  activityJournalList.value = [...defaultSampleJournal]
}

const saveJournalNotes = () => {
  try {
    const key = `lopiq_journal_notes_${authStore.user?.nip || 'default'}`
    localStorage.setItem(key, JSON.stringify(activityJournalList.value))
  } catch (e) {}
}

// Computeds
const userFirstName = computed(() => {
  const name = authStore.user?.name || 'adhe'
  return name.split(' ')[0]
})

const userInitial = computed(() => {
  const name = authStore.user?.name || 'A'
  return name.charAt(0).toUpperCase()
})

const checkInTimeDisplay = computed(() => {
  if (authStore.todayStatus?.masuk?.timestamp) {
    return formatTime(authStore.todayStatus.masuk.timestamp)
  }
  if (authStore.todayStatus?.clock_in_time) {
    return formatTime(authStore.todayStatus.clock_in_time)
  }
  return '--:--'
})

const checkOutTimeDisplay = computed(() => {
  if (authStore.todayStatus?.pulang?.timestamp) {
    return formatTime(authStore.todayStatus.pulang.timestamp)
  }
  if (authStore.todayStatus?.clock_out_time) {
    return formatTime(authStore.todayStatus.clock_out_time)
  }
  return '--:--'
})

// Methods
const formatTime = (ts: string) => {
  if (!ts) return '--:--'
  try {
    const d = new Date(ts)
    if (isNaN(d.getTime())) return ts.substring(0, 5)
    return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  } catch (e) {
    return ts.substring(0, 5)
  }
}

const openAddNoteModal = () => {
  newJournalNote.value = ''
  showNoteModal.value = true
}

const addJournalNote = () => {
  if (!newJournalNote.value.trim()) return
  const now = new Date()
  const timeStr = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
  
  activityJournalList.value.unshift({
    name: userFirstName.value,
    time: timeStr,
    activity: newJournalNote.value.trim()
  })

  saveJournalNotes()
  showNoteModal.value = false
  newJournalNote.value = ''
}

onMounted(async () => {
  await authStore.fetchProfile()
  await authStore.fetchTodayStatus()
  await authStore.fetchHistory()
  loadJournalNotes()
})
</script>

<style scoped>
.pb-safe { 
  padding-bottom: env(safe-area-inset-bottom, 80px); 
}
.material-symbols-outlined { 
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; 
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
