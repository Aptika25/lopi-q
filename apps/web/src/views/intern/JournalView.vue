<template>
  <InternLayout>
    <div class="min-h-screen font-sans text-[#1b1c1c] pb-28 px-4 py-6 max-w-3xl mx-auto w-full space-y-5">

      <!-- Header -->
      <div class="flex items-start justify-between gap-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d]">edit_note</span>
            Jurnal Kegiatan Harian
          </h1>
          <p class="text-xs text-[#574146] mt-1">Catat setiap kegiatan beserta bukti dukungnya. Pembimbing akan menyetujui atau memberi catatan per kegiatan.</p>
        </div>
        <button
          @click="openForm"
          class="shrink-0 px-4 py-2 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-full text-xs font-bold shadow-md flex items-center gap-1.5 border-0 cursor-pointer"
        >
          <span class="material-symbols-outlined text-base">add</span>
          <span>Tambah</span>
        </button>
      </div>

      <!-- Date picker + daily summary -->
      <div class="bg-white rounded-2xl p-4 border border-[#F8BBD0] shadow-sm flex flex-col sm:flex-row gap-3 sm:items-center sm:justify-between">
        <div class="flex items-center gap-2">
          <button @click="shiftDay(-1)" class="p-1.5 rounded-full hover:bg-[#FCE4EC] border-0 bg-transparent cursor-pointer text-[#ab2c5d]" title="Hari sebelumnya">
            <span class="material-symbols-outlined">chevron_left</span>
          </button>
          <input v-model="selectedDate" :max="today" type="date" class="px-3 py-1.5 border border-[#ddbfc5] rounded-lg text-xs font-bold focus:outline-none focus:border-[#f06292]" />
          <button @click="shiftDay(1)" :disabled="selectedDate >= today" class="p-1.5 rounded-full hover:bg-[#FCE4EC] border-0 bg-transparent cursor-pointer text-[#ab2c5d] disabled:opacity-30" title="Hari berikutnya">
            <span class="material-symbols-outlined">chevron_right</span>
          </button>
        </div>
        <div class="flex gap-2 text-[10px] font-bold">
          <span class="px-2 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">{{ countBy('APPROVED') }} disetujui</span>
          <span class="px-2 py-1 rounded-full bg-amber-50 text-amber-700 border border-amber-200">{{ countBy('PENDING') }} menunggu</span>
          <span class="px-2 py-1 rounded-full bg-rose-50 text-rose-700 border border-rose-200">{{ countBy('REJECTED') }} ditolak</span>
        </div>
      </div>

      <div class="text-xs font-bold text-[#574146]">{{ formatDateId(selectedDate) }}</div>

      <div v-if="errorMessage" class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs font-bold">{{ errorMessage }}</div>

      <!-- Timeline -->
      <div v-if="loading" class="text-center text-xs text-slate-500 py-8">Memuat jurnal...</div>
      <div v-else-if="entries.length === 0" class="bg-white rounded-2xl border border-dashed border-[#ddbfc5] p-8 text-center text-xs text-slate-500">
        Belum ada kegiatan pada tanggal ini. Tekan <strong>Tambah</strong> untuk mencatat kegiatan.
      </div>
      <ol v-else class="relative border-l-2 border-[#F8BBD0] ml-3 space-y-4">
        <li v-for="e in entries" :key="e.id" class="ml-5">
          <span class="absolute -left-[9px] mt-4 w-4 h-4 rounded-full border-2 border-white" :class="dotClass(e.status)"></span>
          <div class="bg-white rounded-2xl p-4 border border-[#F8BBD0] shadow-sm space-y-2">
            <div class="flex items-start justify-between gap-2">
              <div>
                <div class="font-mono text-xs font-extrabold text-[#ab2c5d]">{{ e.activity_time }} WITA</div>
                <h3 class="font-bold text-sm leading-snug">{{ e.title }}</h3>
              </div>
              <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border whitespace-nowrap flex items-center gap-1" :class="journalStatusMeta[e.status]?.badge">
                <span class="material-symbols-outlined text-xs">{{ journalStatusMeta[e.status]?.icon }}</span>
                {{ journalStatusMeta[e.status]?.label }}
              </span>
            </div>
            <p v-if="e.description" class="text-xs text-[#574146] leading-relaxed whitespace-pre-line">{{ e.description }}</p>

            <div class="flex items-center gap-3">
              <img v-if="evidenceUrls[e.id]" :src="evidenceUrls[e.id]" alt="Bukti dukung" class="w-20 h-20 object-cover rounded-xl border border-[#ddbfc5] cursor-pointer" @click="showEvidence(e)" />
              <button v-else-if="e.has_evidence" @click="showEvidence(e)" class="flex items-center gap-1 text-[11px] font-bold text-[#ab2c5d] bg-[#FFF5F8] px-3 py-1.5 rounded-lg border border-[#F8BBD0] cursor-pointer">
                <span class="material-symbols-outlined text-sm">{{ e.evidence_mime === 'application/pdf' ? 'picture_as_pdf' : 'image' }}</span>
                <span>Lihat Bukti Dukung</span>
              </button>
            </div>

            <div v-if="e.status !== 'PENDING' || e.reviewer_note" class="p-3 rounded-xl text-[11px] leading-relaxed" :class="e.status === 'REJECTED' ? 'bg-rose-50 text-rose-800' : 'bg-[#FFF5F8] text-[#574146]'">
              <div class="font-bold">Catatan Pembimbing{{ e.reviewed_by ? ` (${e.reviewed_by})` : '' }}:</div>
              <div class="whitespace-pre-line">{{ e.reviewer_note || '-' }}</div>
            </div>

            <div v-if="e.status === 'PENDING'" class="flex justify-end">
              <button @click="removeEntry(e)" class="text-[11px] font-bold text-rose-600 bg-transparent border-0 cursor-pointer flex items-center gap-1">
                <span class="material-symbols-outlined text-sm">delete</span>
                <span>Hapus</span>
              </button>
            </div>
          </div>
        </li>
      </ol>

      <!-- Add entry bottom sheet -->
      <div v-if="showForm" class="fixed inset-0 bg-slate-900/50 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4" @click.self="showForm = false">
        <form @submit.prevent="submitEntry" class="w-full max-w-md bg-white rounded-t-3xl sm:rounded-3xl p-6 space-y-4 shadow-2xl border-t border-[#ddbfc5] text-xs max-h-[92vh] overflow-y-auto">
          <div class="flex justify-between items-center border-b border-[#ddbfc5]/60 pb-3">
            <h2 class="text-base font-bold flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">edit_document</span>
              Tambah Kegiatan
            </h2>
            <button type="button" @click="showForm = false" class="p-1 rounded-full hover:bg-slate-100 border-0 bg-transparent cursor-pointer text-slate-400">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tanggal <span class="text-rose-500">*</span></label>
              <input v-model="form.activity_date" :max="today" type="date" required class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none" />
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Jam <span class="text-rose-500">*</span></label>
              <input v-model="form.activity_time" type="time" required class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none" />
            </div>
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Kegiatan <span class="text-rose-500">*</span></label>
            <input v-model="form.title" type="text" required maxlength="255" placeholder="Contoh: Mengunggah berita ke website" class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none" />
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Uraian (opsional)</label>
            <textarea v-model="form.description" rows="3" placeholder="Jelaskan singkat apa yang dikerjakan..." class="w-full p-3 rounded-xl border border-[#ddbfc5] bg-[#f5f3f3] focus:border-[#f06292] focus:outline-none resize-none"></textarea>
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Bukti Dukung <span class="text-rose-500">*</span></label>
            <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp,application/pdf" required @change="onFileChange" class="w-full text-[11px] file:mr-3 file:px-3 file:py-2 file:rounded-lg file:border-0 file:bg-[#FCE4EC] file:text-[#ab2c5d] file:font-bold" />
            <p class="text-[10px] text-[#8a7176]">Foto / screenshot (JPG, PNG, WEBP) atau PDF, maksimal 3 MB.</p>
            <img v-if="previewUrl" :src="previewUrl" alt="Pratinjau" class="mt-2 max-h-40 rounded-xl border border-[#ddbfc5]" />
          </div>

          <div v-if="formError" class="p-3 bg-rose-50 text-rose-700 rounded-xl font-bold border border-rose-200">{{ formError }}</div>

          <button type="submit" :disabled="submitting" class="w-full py-3.5 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-full font-bold shadow-lg border-0 cursor-pointer disabled:opacity-60">
            {{ submitting ? 'Mengirim...' : 'Kirim Kegiatan' }}
          </button>
        </form>
      </div>
    </div>
  </InternLayout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import InternLayout from '@/layouts/InternLayout.vue'
import { type JournalEntry, journalStatusMeta, loadEvidenceUrl, openEvidence, todayWita, formatDateId } from '@/utils/journal'

const MAX_BYTES = 3 * 1024 * 1024

const today = todayWita()
const selectedDate = ref(today)
const entries = ref<JournalEntry[]>([])
const evidenceUrls = ref<Record<number, string>>({})
const loading = ref(false)
const errorMessage = ref('')

const showForm = ref(false)
const submitting = ref(false)
const formError = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const previewUrl = ref('')
const form = ref({ activity_date: today, activity_time: '', title: '', description: '' })

function countBy(status: string) {
  return entries.value.filter(e => e.status === status).length
}

function dotClass(status: string) {
  if (status === 'APPROVED') return 'bg-emerald-500'
  if (status === 'REJECTED') return 'bg-rose-500'
  return 'bg-amber-400'
}

function shiftDay(delta: number) {
  const d = new Date(`${selectedDate.value}T00:00:00`)
  d.setDate(d.getDate() + delta)
  const next = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  if (next <= today) selectedDate.value = next
}

async function loadEntries() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.get('/api/journal', { params: { date: selectedDate.value } })
    entries.value = res.data.entries || []
    for (const e of entries.value) {
      if (e.has_evidence && e.evidence_mime.startsWith('image/') && !evidenceUrls.value[e.id]) {
        loadEvidenceUrl(e.id).then(url => { evidenceUrls.value[e.id] = url }).catch(() => {})
      }
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal memuat jurnal kegiatan.'
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

function nowTime() {
  return new Intl.DateTimeFormat('en-GB', { timeZone: 'Asia/Makassar', hour: '2-digit', minute: '2-digit', hour12: false }).format(new Date())
}

function resetFile() {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = ''
  selectedFile.value = null
  if (fileInput.value) fileInput.value.value = ''
}

function openForm() {
  formError.value = ''
  resetFile()
  form.value = { activity_date: selectedDate.value, activity_time: selectedDate.value === today ? nowTime() : '', title: '', description: '' }
  showForm.value = true
}

function onFileChange(ev: Event) {
  formError.value = ''
  const file = (ev.target as HTMLInputElement).files?.[0] || null
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = ''
  selectedFile.value = null
  if (!file) return
  if (file.size > MAX_BYTES) {
    formError.value = 'Ukuran bukti dukung maksimal 3 MB.'
    if (fileInput.value) fileInput.value.value = ''
    return
  }
  selectedFile.value = file
  if (file.type.startsWith('image/')) previewUrl.value = URL.createObjectURL(file)
}

async function submitEntry() {
  if (!selectedFile.value) {
    formError.value = 'Bukti dukung wajib diunggah.'
    return
  }
  submitting.value = true
  formError.value = ''
  const data = new FormData()
  data.append('activity_date', form.value.activity_date)
  data.append('activity_time', form.value.activity_time)
  data.append('title', form.value.title)
  data.append('description', form.value.description)
  data.append('evidence', selectedFile.value)
  try {
    await axios.post('/api/journal', data)
    showForm.value = false
    resetFile()
    if (selectedDate.value !== form.value.activity_date) selectedDate.value = form.value.activity_date
    else await loadEntries()
  } catch (err: any) {
    formError.value = err.response?.data?.error || 'Gagal mengirim kegiatan.'
  } finally {
    submitting.value = false
  }
}

async function removeEntry(e: JournalEntry) {
  if (!confirm(`Hapus kegiatan "${e.title}"?`)) return
  try {
    await axios.delete(`/api/journal/${e.id}`)
    await loadEntries()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal menghapus kegiatan.'
  }
}

watch(selectedDate, loadEntries)
onMounted(loadEntries)
onUnmounted(resetFile)
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
