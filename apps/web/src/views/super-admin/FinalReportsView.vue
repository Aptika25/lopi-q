<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">

      <!-- Header Section -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill">menu_book</span>
            Laporan Akhir
          </h1>
          <p class="text-sm text-[#574146] mt-1">Pantau pengumpulan, review, dan penilaian laporan akhir peserta magang.</p>
        </div>

        <button
          @click="openForm()"
          class="shrink-0 px-4 py-2 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-lg text-xs font-bold shadow-sm flex items-center gap-1.5 border-0 cursor-pointer transition-colors"
        >
          <span class="material-symbols-outlined text-base">add</span>
          <span>Tambah Laporan</span>
        </button>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div v-for="(meta, key) in statusMeta" :key="key" class="bg-white/90 rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-[#8a7176]">{{ meta.label }}</div>
          <div class="text-2xl font-extrabold text-[#1b1c1c] mt-1">{{ reports.filter(r => r.status === key).length }}</div>
        </div>
      </div>

      <!-- Filters & Search Bar -->
      <div class="bg-white/85 backdrop-blur-md rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="relative w-full md:w-60">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">filter_list</span>
          <select
            v-model="filterStatus"
            class="w-full pl-10 pr-8 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] appearance-none cursor-pointer"
          >
            <option value="">Semua Status</option>
            <option v-for="(meta, key) in statusMeta" :key="key" :value="key">{{ meta.label }}</option>
          </select>
        </div>
        <div class="relative w-full md:w-72">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari peserta, instansi, atau judul..."
            class="w-full pl-10 pr-4 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-medium text-[#1b1c1c] focus:outline-none focus:border-[#f06292]"
          />
        </div>
      </div>

      <div v-if="errorMessage" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs font-bold">{{ errorMessage }}</div>

      <!-- Reports Table -->
      <div class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                <th class="py-3.5 px-6 min-w-[200px]">Peserta</th>
                <th class="py-3.5 px-6 min-w-[240px]">Judul Laporan</th>
                <th class="py-3.5 px-6 w-32">Dikumpulkan</th>
                <th class="py-3.5 px-6 w-36">Status</th>
                <th class="py-3.5 px-6 w-20">Nilai</th>
                <th class="py-3.5 px-6 w-24 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="text-xs text-[#1b1c1c] bg-white divide-y divide-[#F8BBD0]">
              <tr v-if="loading">
                <td colspan="6" class="py-6 text-center text-slate-500 font-medium">Memuat laporan akhir...</td>
              </tr>
              <tr v-for="r in filteredReports" v-else :key="r.id" class="hover:bg-[#FCE4EC]/30 transition-colors">
                <td class="py-4 px-6 align-top">
                  <div class="font-bold text-xs leading-tight">{{ r.intern_name }}</div>
                  <div v-if="r.institution" class="text-[10px] text-[#574146] font-medium mt-0.5">{{ r.institution }}</div>
                </td>
                <td class="py-4 px-6 align-top">
                  <div class="font-bold text-xs">{{ r.title }}</div>
                  <a v-if="r.document_url" :href="r.document_url" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1 text-[11px] font-semibold text-[#ab2c5d] hover:underline mt-0.5">
                    <span class="material-symbols-outlined text-sm">open_in_new</span>
                    <span>Buka Dokumen</span>
                  </a>
                  <div v-if="r.reviewer_notes" class="text-[11px] text-[#574146] leading-relaxed mt-1 line-clamp-2">Catatan: {{ r.reviewer_notes }}</div>
                </td>
                <td class="py-4 px-6 align-top font-mono font-bold text-[#574146] whitespace-nowrap">{{ formatDate(r.submitted_at) }}</td>
                <td class="py-4 px-6 align-top">
                  <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border whitespace-nowrap" :class="statusMeta[r.status]?.badge">
                    {{ statusMeta[r.status]?.label || r.status }}
                  </span>
                </td>
                <td class="py-4 px-6 align-top font-mono font-extrabold text-sm">{{ r.score ?? '-' }}</td>
                <td class="py-4 px-6 align-top text-right whitespace-nowrap">
                  <button @click="openForm(r)" class="p-1.5 rounded-lg hover:bg-[#FCE4EC] text-[#ab2c5d] border-0 bg-transparent cursor-pointer" title="Review / Ubah">
                    <span class="material-symbols-outlined text-base">rate_review</span>
                  </button>
                  <button @click="removeReport(r)" class="p-1.5 rounded-lg hover:bg-rose-50 text-rose-600 border-0 bg-transparent cursor-pointer" title="Hapus">
                    <span class="material-symbols-outlined text-base">delete</span>
                  </button>
                </td>
              </tr>
              <tr v-if="!loading && filteredReports.length === 0">
                <td colspan="6" class="py-6 text-center">
                  <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
                    <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
                    <span>Belum ada laporan akhir.</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Form Modal -->
      <div v-if="showForm" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="showForm = false">
        <form @submit.prevent="saveReport" class="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto p-6 space-y-4 border border-[#F8BBD0] shadow-2xl text-xs">
          <div class="flex justify-between items-center border-b border-[#F8BBD0]/60 pb-3">
            <h2 class="text-base font-bold text-[#1b1c1c]">{{ form.id ? 'Review Laporan Akhir' : 'Tambah Laporan Akhir' }}</h2>
            <button type="button" @click="showForm = false" class="p-1 rounded-full hover:bg-slate-100 border-0 bg-transparent cursor-pointer text-slate-400">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Peserta Magang <span class="text-rose-500">*</span></label>
            <select v-model.number="form.intern_id" @change="onInternChange" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] bg-white">
              <option :value="0" disabled>Pilih peserta...</option>
              <option v-if="form.intern_id && !interns.some((u: any) => u.id === form.intern_id)" :value="form.intern_id">{{ form.intern_name }}</option>
              <option v-for="intern in interns" :key="intern.id" :value="intern.id">{{ intern.name }}</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Asal Sekolah / Kampus</label>
            <input v-model="form.institution" type="text" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Judul Laporan <span class="text-rose-500">*</span></label>
            <input v-model="form.title" required type="text" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Link Dokumen (Google Drive / PDF)</label>
            <input v-model="form.document_url" type="url" placeholder="https://..." class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tanggal Dikumpulkan</label>
              <input v-model="form.submitted_at" type="date" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Status</label>
              <select v-model="form.status" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] bg-white">
                <option v-for="(meta, key) in statusMeta" :key="key" :value="key">{{ meta.label }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Nilai (0 - 100)</label>
              <input v-model="form.score" type="number" min="0" max="100" step="0.5" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Catatan Reviewer</label>
            <textarea v-model="form.reviewer_notes" rows="3" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] resize-none"></textarea>
          </div>

          <div v-if="formError" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 font-bold">{{ formError }}</div>

          <div class="flex justify-end gap-2 pt-2">
            <button type="button" @click="showForm = false" class="px-4 py-2 rounded-lg border border-[#ddbfc5] bg-white font-bold text-[#574146] cursor-pointer">Batal</button>
            <button type="submit" :disabled="saving" class="px-4 py-2 rounded-lg bg-[#ab2c5d] hover:bg-[#881b47] text-white font-bold border-0 cursor-pointer disabled:opacity-60">
              {{ saving ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import AdminLayout from '@/layouts/AdminLayout.vue'
import { useAuthStore } from '@/stores/auth'

interface FinalReport {
  id: number
  intern_id: number
  intern_name: string
  institution: string
  title: string
  document_url: string
  submitted_at: string
  status: string
  score: number | null
  reviewer_notes: string
}

const authStore = useAuthStore()

const statusMeta: Record<string, { label: string; badge: string }> = {
  BELUM_DIKUMPULKAN: { label: 'Belum Dikumpulkan', badge: 'bg-slate-100 text-slate-600 border-slate-200' },
  MENUNGGU_REVIEW: { label: 'Menunggu Review', badge: 'bg-amber-50 text-amber-700 border-amber-200' },
  REVISI: { label: 'Perlu Revisi', badge: 'bg-rose-50 text-rose-700 border-rose-200' },
  DISETUJUI: { label: 'Disetujui', badge: 'bg-emerald-50 text-emerald-700 border-emerald-200' }
}

const reports = ref<FinalReport[]>([])
const loading = ref(false)
const errorMessage = ref('')
const filterStatus = ref('')
const searchQuery = ref('')

const showForm = ref(false)
const saving = ref(false)
const formError = ref('')
const emptyForm = () => ({
  id: 0,
  intern_id: 0,
  intern_name: '',
  institution: '',
  title: '',
  document_url: '',
  submitted_at: '',
  status: 'BELUM_DIKUMPULKAN',
  score: '' as number | string,
  reviewer_notes: ''
})
const form = ref(emptyForm())

const interns = computed(() =>
  (authStore.usersList || []).filter((u: any) => {
    const role = String(u.role || '').toLowerCase()
    return role === 'intern' || role === 'peserta'
  })
)

const filteredReports = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return reports.value.filter(r => {
    if (filterStatus.value && r.status !== filterStatus.value) return false
    if (!q) return true
    return r.intern_name.toLowerCase().includes(q) ||
      r.institution.toLowerCase().includes(q) ||
      r.title.toLowerCase().includes(q)
  })
})

function formatDate(iso: string) {
  if (!iso) return '-'
  const [y, m, d] = iso.split('-')
  return `${d}/${m}/${y}`
}

function onInternChange() {
  const intern: any = interns.value.find((u: any) => u.id === form.value.intern_id)
  if (!intern) return
  form.value.intern_name = intern.name
  // jabatan holds the intern's school / campus in this system
  if (!form.value.institution) form.value.institution = intern.jabatan || ''
}

async function loadReports() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.get('/api/admin/final-reports')
    reports.value = res.data.reports || []
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal memuat laporan akhir.'
  } finally {
    loading.value = false
  }
}

function openForm(r?: FinalReport) {
  formError.value = ''
  form.value = r ? { ...r, score: r.score ?? '' } : emptyForm()
  showForm.value = true
}

async function saveReport() {
  saving.value = true
  formError.value = ''
  const scoreRaw = String(form.value.score ?? '').trim()
  const payload = {
    intern_id: form.value.intern_id,
    intern_name: form.value.intern_name,
    institution: form.value.institution,
    title: form.value.title,
    document_url: form.value.document_url,
    submitted_at: form.value.submitted_at,
    status: form.value.status,
    score: scoreRaw === '' ? null : Number(scoreRaw),
    reviewer_notes: form.value.reviewer_notes
  }
  try {
    if (form.value.id) {
      await axios.put(`/api/admin/final-reports/${form.value.id}`, payload)
    } else {
      await axios.post('/api/admin/final-reports', payload)
    }
    showForm.value = false
    await loadReports()
  } catch (err: any) {
    formError.value = err.response?.data?.error || 'Gagal menyimpan laporan akhir.'
  } finally {
    saving.value = false
  }
}

async function removeReport(r: FinalReport) {
  if (!confirm(`Hapus laporan akhir "${r.title}" milik ${r.intern_name}?`)) return
  try {
    await axios.delete(`/api/admin/final-reports/${r.id}`)
    await loadReports()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal menghapus laporan akhir.'
  }
}

onMounted(async () => {
  await Promise.all([loadReports(), authStore.fetchUsers()])
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
