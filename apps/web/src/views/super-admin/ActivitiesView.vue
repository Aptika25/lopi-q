<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">

      <!-- Header Section -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill">history_edu</span>
            Jurnal Aktivitas
          </h1>
          <p class="text-sm text-[#574146] mt-1">Aktivitas harian setiap peserta magang beserta bukti dukung. Setujui atau tolak tiap kegiatan dengan catatan.</p>
        </div>

        <div class="flex items-center gap-2 shrink-0">
          <span class="text-xs font-bold text-[#ab2c5d] bg-[#ffd9e4] px-4 py-2 rounded-lg border border-[#F8BBD0] flex items-center gap-1.5 shadow-xs">
            <span class="material-symbols-outlined text-base">calendar_today</span>
            <span>{{ filterDate ? formatDateId(filterDate) : 'Semua Tanggal' }}</span>
          </span>
        </div>
      </div>

      <!-- Summary -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div class="bg-white/90 rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-[#8a7176]">Total Kegiatan</div>
          <div class="text-2xl font-extrabold text-[#1b1c1c] mt-1">{{ entries.length }}</div>
        </div>
        <div v-for="(meta, key) in journalStatusMeta" :key="key" class="bg-white/90 rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-[#8a7176]">{{ meta.label }}</div>
          <div class="text-2xl font-extrabold text-[#1b1c1c] mt-1">{{ entries.filter(e => e.status === key).length }}</div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white/85 backdrop-blur-md rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col md:flex-row gap-3 items-stretch md:items-center">
        <div class="flex items-center gap-1">
          <input v-model="filterDate" type="date" class="px-3 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292]" />
          <button v-if="filterDate" @click="filterDate = ''" class="px-2 py-2 text-[11px] font-bold text-[#ab2c5d] bg-transparent border-0 cursor-pointer" title="Tampilkan semua tanggal">Semua</button>
          <button v-else @click="filterDate = todayWita()" class="px-2 py-2 text-[11px] font-bold text-[#ab2c5d] bg-transparent border-0 cursor-pointer">Hari ini</button>
        </div>
        <select v-model="filterUser" class="w-full md:w-60 px-3 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] cursor-pointer">
          <option :value="0">Semua Peserta Magang</option>
          <option v-for="intern in interns" :key="intern.id" :value="intern.id">{{ intern.name }}</option>
        </select>
        <select v-model="filterStatus" class="w-full md:w-48 px-3 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] cursor-pointer">
          <option value="">Semua Status</option>
          <option v-for="(meta, key) in journalStatusMeta" :key="key" :value="key">{{ meta.label }}</option>
        </select>
        <div class="relative w-full md:w-64 md:ml-auto">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">search</span>
          <input v-model="searchQuery" type="text" placeholder="Cari kegiatan..." class="w-full pl-10 pr-4 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-medium focus:outline-none focus:border-[#f06292]" />
        </div>
      </div>

      <div v-if="errorMessage" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs font-bold">{{ errorMessage }}</div>

      <div v-if="loading" class="text-center text-xs text-slate-500 py-8">Memuat jurnal aktivitas...</div>

      <div v-else-if="groups.length === 0" class="bg-white/90 rounded-xl border border-[#F8BBD0] py-8 text-center">
        <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
          <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
          <span>Belum ada log aktivitas yang diunggah.</span>
        </div>
      </div>

      <!-- One card per intern per day -->
      <div v-for="g in groups" v-else :key="g.key" class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        <div class="bg-[#FCE4EC] border-b border-[#F8BBD0] px-5 py-3 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-full bg-white text-[#ab2c5d] flex items-center justify-center font-extrabold text-sm border border-[#F8BBD0]">{{ g.userName.charAt(0) }}</div>
            <div>
              <div class="font-bold text-sm text-[#1b1c1c]">{{ g.userName }}</div>
              <div class="text-[10px] text-[#574146] font-medium">{{ g.institution || (g.userNip ? `NISN: ${g.userNip}` : '') }}</div>
            </div>
          </div>
          <div class="text-[11px] font-bold text-[#574146] flex items-center gap-2">
            <span>{{ formatDateId(g.date) }}</span>
            <span class="px-2 py-0.5 rounded-full bg-white border border-[#F8BBD0]">{{ g.items.length }} kegiatan</span>
          </div>
        </div>

        <div class="divide-y divide-[#F8BBD0]">
          <div v-for="e in g.items" :key="e.id" class="p-5 flex flex-col lg:flex-row gap-4">
            <div class="flex gap-4 flex-1 min-w-0">
              <div class="font-mono text-xs font-extrabold text-[#ab2c5d] w-14 shrink-0 pt-0.5">{{ e.activity_time }}</div>
              <div class="flex-1 min-w-0 space-y-1.5">
                <div class="flex items-start gap-2 flex-wrap">
                  <h4 class="font-bold text-xs text-[#1b1c1c]">{{ e.title }}</h4>
                  <span class="px-2 py-0.5 rounded-full text-[9px] font-extrabold uppercase border whitespace-nowrap" :class="journalStatusMeta[e.status]?.badge">{{ journalStatusMeta[e.status]?.label }}</span>
                </div>
                <p v-if="e.description" class="text-[11px] text-[#574146] leading-relaxed whitespace-pre-line">{{ e.description }}</p>
                <div v-if="e.reviewed_by" class="text-[10px] text-[#8a7176]">Direview oleh {{ e.reviewed_by }} &middot; {{ e.reviewed_at }}</div>
              </div>
              <div class="shrink-0">
                <img v-if="evidenceUrls[e.id]" :src="evidenceUrls[e.id]" alt="Bukti dukung" class="w-24 h-24 object-cover rounded-lg border border-[#ddbfc5] cursor-zoom-in" @click="showEvidence(e)" />
                <button v-else-if="e.has_evidence" @click="showEvidence(e)" class="flex flex-col items-center justify-center w-24 h-24 rounded-lg border border-[#ddbfc5] bg-[#FFF5F8] text-[#ab2c5d] text-[10px] font-bold cursor-pointer gap-1">
                  <span class="material-symbols-outlined">{{ e.evidence_mime === 'application/pdf' ? 'picture_as_pdf' : 'image' }}</span>
                  Bukti Dukung
                </button>
              </div>
            </div>

            <!-- Review panel -->
            <div class="lg:w-80 shrink-0 space-y-2">
              <textarea
                v-model="notes[e.id]"
                rows="2"
                placeholder="Catatan untuk peserta (wajib bila tidak disetujui)"
                class="w-full px-3 py-2 border border-[#ddbfc5] rounded-lg text-[11px] focus:outline-none focus:border-[#f06292] resize-none"
              ></textarea>
              <div class="flex gap-2">
                <button
                  @click="review(e, 'APPROVED')"
                  :disabled="savingId === e.id"
                  class="flex-1 px-3 py-1.5 rounded-lg text-[11px] font-bold border cursor-pointer flex items-center justify-center gap-1 disabled:opacity-60"
                  :class="e.status === 'APPROVED' ? 'bg-emerald-600 text-white border-emerald-600' : 'bg-white text-emerald-700 border-emerald-300 hover:bg-emerald-50'"
                >
                  <span class="material-symbols-outlined text-sm">check</span> Setujui
                </button>
                <button
                  @click="review(e, 'REJECTED')"
                  :disabled="savingId === e.id"
                  class="flex-1 px-3 py-1.5 rounded-lg text-[11px] font-bold border cursor-pointer flex items-center justify-center gap-1 disabled:opacity-60"
                  :class="e.status === 'REJECTED' ? 'bg-rose-600 text-white border-rose-600' : 'bg-white text-rose-700 border-rose-300 hover:bg-rose-50'"
                >
                  <span class="material-symbols-outlined text-sm">close</span> Tidak Disetujui
                </button>
              </div>
              <div v-if="rowErrors[e.id]" class="text-[10px] font-bold text-rose-600">{{ rowErrors[e.id] }}</div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import axios from 'axios'
import AdminLayout from '@/layouts/AdminLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { type JournalEntry, journalStatusMeta, loadEvidenceUrl, openEvidence, todayWita, formatDateId } from '@/utils/journal'

const authStore = useAuthStore()

const filterDate = ref(todayWita())
const filterUser = ref(0)
const filterStatus = ref('')
const searchQuery = ref('')

const entries = ref<JournalEntry[]>([])
const evidenceUrls = ref<Record<number, string>>({})
const notes = ref<Record<number, string>>({})
const rowErrors = ref<Record<number, string>>({})
const savingId = ref<number | null>(null)
const loading = ref(false)
const errorMessage = ref('')

const interns = computed(() =>
  (authStore.usersList || []).filter((u: any) => {
    const role = String(u.role || '').toLowerCase()
    return role === 'intern' || role === 'peserta'
  })
)

const internById = computed(() => new Map(interns.value.map((u: any) => [u.id, u])))

const groups = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  const map = new Map<string, { key: string; date: string; userName: string; userNip: string; institution: string; items: JournalEntry[] }>()
  for (const e of entries.value) {
    if (q && !e.title.toLowerCase().includes(q) && !e.description.toLowerCase().includes(q) && !e.user_name.toLowerCase().includes(q)) continue
    const key = `${e.activity_date}|${e.user_id}`
    if (!map.has(key)) {
      const profile: any = internById.value.get(e.user_id)
      map.set(key, { key, date: e.activity_date, userName: e.user_name || profile?.name || `Peserta #${e.user_id}`, userNip: e.user_nip, institution: profile?.jabatan || '', items: [] })
    }
    map.get(key)!.items.push(e)
  }
  return Array.from(map.values())
})

async function loadEntries() {
  loading.value = true
  errorMessage.value = ''
  try {
    const params: Record<string, string | number> = {}
    if (filterDate.value) params.date = filterDate.value
    if (filterUser.value) params.user_id = filterUser.value
    if (filterStatus.value) params.status = filterStatus.value
    const res = await axios.get('/api/journal', { params })
    entries.value = res.data.entries || []
    for (const e of entries.value) {
      if (notes.value[e.id] === undefined) notes.value[e.id] = e.reviewer_note || ''
      if (e.has_evidence && e.evidence_mime.startsWith('image/') && !evidenceUrls.value[e.id]) {
        loadEvidenceUrl(e.id).then(url => { evidenceUrls.value[e.id] = url }).catch(() => {})
      }
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal memuat jurnal aktivitas.'
  } finally {
    loading.value = false
  }
}

async function showEvidence(e: JournalEntry) {
  try {
    await openEvidence(e.id)
  } catch {
    errorMessage.value = 'Gagal membuka bukti dukung.'
  }
}

async function review(e: JournalEntry, status: 'APPROVED' | 'REJECTED') {
  rowErrors.value[e.id] = ''
  const note = (notes.value[e.id] || '').trim()
  if (status === 'REJECTED' && !note) {
    rowErrors.value[e.id] = 'Isi catatan alasan sebelum menolak kegiatan.'
    return
  }
  savingId.value = e.id
  try {
    const res = await axios.put(`/api/admin/journal/${e.id}/review`, { status, note })
    const idx = entries.value.findIndex(x => x.id === e.id)
    if (idx >= 0 && res.data.entry) entries.value[idx] = res.data.entry
  } catch (err: any) {
    rowErrors.value[e.id] = err.response?.data?.error || 'Gagal menyimpan review.'
  } finally {
    savingId.value = null
  }
}

watch([filterDate, filterUser, filterStatus], loadEntries)

onMounted(async () => {
  await Promise.all([loadEntries(), authStore.fetchUsers()])
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
