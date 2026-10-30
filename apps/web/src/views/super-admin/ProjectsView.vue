<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">

      <!-- Header Section -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill">rocket_launch</span>
            Proyek Magang
          </h1>
          <p class="text-sm text-[#574146] mt-1">Kelola proyek yang dikerjakan peserta magang beserta progres dan tenggat waktunya.</p>
        </div>

        <button
          @click="openForm()"
          class="shrink-0 px-4 py-2 bg-[#ab2c5d] hover:bg-[#881b47] text-white rounded-lg text-xs font-bold shadow-sm flex items-center gap-1.5 border-0 cursor-pointer transition-colors"
        >
          <span class="material-symbols-outlined text-base">add</span>
          <span>Tambah Proyek</span>
        </button>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <div v-for="card in summaryCards" :key="card.key" class="bg-white/90 rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-[#8a7176]">{{ card.label }}</div>
          <div class="text-2xl font-extrabold text-[#1b1c1c] mt-1">{{ card.value }}</div>
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
            placeholder="Cari proyek, pembimbing, atau peserta..."
            class="w-full pl-10 pr-4 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-medium text-[#1b1c1c] focus:outline-none focus:border-[#f06292]"
          />
        </div>
      </div>

      <div v-if="errorMessage" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs font-bold">{{ errorMessage }}</div>

      <!-- Projects Table -->
      <div class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                <th class="py-3.5 px-6 min-w-[240px]">Proyek</th>
                <th class="py-3.5 px-6 min-w-[180px]">Peserta</th>
                <th class="py-3.5 px-6 w-40">Progres</th>
                <th class="py-3.5 px-6 w-32">Tenggat</th>
                <th class="py-3.5 px-6 w-32">Status</th>
                <th class="py-3.5 px-6 w-24 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="text-xs text-[#1b1c1c] bg-white divide-y divide-[#F8BBD0]">
              <tr v-if="loading">
                <td colspan="6" class="py-6 text-center text-slate-500 font-medium">Memuat data proyek...</td>
              </tr>
              <tr v-for="p in filteredProjects" v-else :key="p.id" class="hover:bg-[#FCE4EC]/30 transition-colors">
                <td class="py-4 px-6 align-top">
                  <div class="font-bold text-xs leading-tight">{{ p.title }}</div>
                  <div v-if="p.description" class="text-[11px] text-[#574146] leading-relaxed mt-0.5 line-clamp-2">{{ p.description }}</div>
                  <div v-if="p.mentor" class="text-[10px] text-[#8a7176] font-semibold mt-1">Pembimbing: {{ p.mentor }}</div>
                </td>
                <td class="py-4 px-6 align-top">
                  <div v-if="p.members.length" class="flex flex-wrap gap-1">
                    <span v-for="m in p.members" :key="m.id" class="px-2 py-0.5 rounded-full bg-[#FFF5F8] border border-[#F8BBD0] text-[10px] font-semibold text-[#ab2c5d]">{{ m.name }}</span>
                  </div>
                  <span v-else class="text-[11px] text-slate-400">Belum ada peserta</span>
                </td>
                <td class="py-4 px-6 align-top">
                  <div class="flex items-center gap-2">
                    <div class="flex-1 h-2 rounded-full bg-[#FCE4EC] overflow-hidden">
                      <div class="h-full bg-[#f06292] rounded-full" :style="{ width: p.progress + '%' }"></div>
                    </div>
                    <span class="font-mono font-bold text-[11px] text-[#574146] w-9 text-right">{{ p.progress }}%</span>
                  </div>
                </td>
                <td class="py-4 px-6 align-top font-mono font-bold whitespace-nowrap" :class="isOverdue(p) ? 'text-rose-600' : 'text-[#574146]'">
                  {{ formatDate(p.due_date) }}
                </td>
                <td class="py-4 px-6 align-top">
                  <span class="px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase border whitespace-nowrap" :class="statusMeta[p.status]?.badge">
                    {{ statusMeta[p.status]?.label || p.status }}
                  </span>
                </td>
                <td class="py-4 px-6 align-top text-right whitespace-nowrap">
                  <button @click="openForm(p)" class="p-1.5 rounded-lg hover:bg-[#FCE4EC] text-[#ab2c5d] border-0 bg-transparent cursor-pointer" title="Ubah">
                    <span class="material-symbols-outlined text-base">edit</span>
                  </button>
                  <button @click="removeProject(p)" class="p-1.5 rounded-lg hover:bg-rose-50 text-rose-600 border-0 bg-transparent cursor-pointer" title="Hapus">
                    <span class="material-symbols-outlined text-base">delete</span>
                  </button>
                </td>
              </tr>
              <tr v-if="!loading && filteredProjects.length === 0">
                <td colspan="6" class="py-6 text-center">
                  <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
                    <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
                    <span>Belum ada proyek magang.</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Form Modal -->
      <div v-if="showForm" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="showForm = false">
        <form @submit.prevent="saveProject" class="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto p-6 space-y-4 border border-[#F8BBD0] shadow-2xl text-xs">
          <div class="flex justify-between items-center border-b border-[#F8BBD0]/60 pb-3">
            <h2 class="text-base font-bold text-[#1b1c1c]">{{ form.id ? 'Ubah Proyek' : 'Tambah Proyek' }}</h2>
            <button type="button" @click="showForm = false" class="p-1 rounded-full hover:bg-slate-100 border-0 bg-transparent cursor-pointer text-slate-400">
              <span class="material-symbols-outlined">close</span>
            </button>
          </div>

          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Judul Proyek <span class="text-rose-500">*</span></label>
            <input v-model="form.title" required type="text" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Deskripsi</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] resize-none"></textarea>
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Pembimbing</label>
            <input v-model="form.mentor" type="text" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
          </div>
          <div class="space-y-1">
            <label class="font-bold text-[#574146]">Peserta Magang</label>
            <div class="max-h-36 overflow-y-auto border border-[#ddbfc5] rounded-lg p-2 space-y-1">
              <label v-for="intern in interns" :key="intern.id" class="flex items-center gap-2 px-1.5 py-1 rounded hover:bg-[#FFF5F8] cursor-pointer">
                <input type="checkbox" :value="intern.id" v-model="form.memberIds" class="accent-[#ab2c5d]" />
                <span class="font-semibold">{{ intern.name }}</span>
                <span class="text-[10px] text-[#8a7176]">{{ intern.jabatan }}</span>
              </label>
              <div v-if="interns.length === 0" class="text-[11px] text-slate-400 px-1.5 py-1">Belum ada peserta magang terdaftar.</div>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Status</label>
              <select v-model="form.status" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] bg-white">
                <option v-for="(meta, key) in statusMeta" :key="key" :value="key">{{ meta.label }}</option>
              </select>
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Progres ({{ form.progress }}%)</label>
              <input v-model.number="form.progress" type="range" min="0" max="100" step="5" class="w-full accent-[#ab2c5d] mt-2" />
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tanggal Mulai</label>
              <input v-model="form.start_date" type="date" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>
            <div class="space-y-1">
              <label class="font-bold text-[#574146]">Tenggat</label>
              <input v-model="form.due_date" type="date" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>
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

interface Member { id: number; name: string }
interface Project {
  id: number
  title: string
  description: string
  mentor: string
  members: Member[]
  status: string
  progress: number
  start_date: string
  due_date: string
}

const authStore = useAuthStore()

const statusMeta: Record<string, { label: string; badge: string }> = {
  PERENCANAAN: { label: 'Perencanaan', badge: 'bg-sky-50 text-sky-700 border-sky-200' },
  BERJALAN: { label: 'Berjalan', badge: 'bg-amber-50 text-amber-700 border-amber-200' },
  SELESAI: { label: 'Selesai', badge: 'bg-emerald-50 text-emerald-700 border-emerald-200' },
  DITUNDA: { label: 'Ditunda', badge: 'bg-slate-100 text-slate-600 border-slate-200' }
}

const projects = ref<Project[]>([])
const loading = ref(false)
const errorMessage = ref('')
const filterStatus = ref('')
const searchQuery = ref('')

const showForm = ref(false)
const saving = ref(false)
const formError = ref('')
const emptyForm = () => ({
  id: 0,
  title: '',
  description: '',
  mentor: '',
  memberIds: [] as number[],
  status: 'PERENCANAAN',
  progress: 0,
  start_date: '',
  due_date: ''
})
const form = ref(emptyForm())

const interns = computed(() =>
  (authStore.usersList || []).filter((u: any) => {
    const role = String(u.role || '').toLowerCase()
    return role === 'intern' || role === 'peserta'
  })
)

const summaryCards = computed(() => [
  { key: 'total', label: 'Total Proyek', value: projects.value.length },
  { key: 'run', label: 'Berjalan', value: projects.value.filter(p => p.status === 'BERJALAN').length },
  { key: 'done', label: 'Selesai', value: projects.value.filter(p => p.status === 'SELESAI').length },
  { key: 'late', label: 'Lewat Tenggat', value: projects.value.filter(isOverdue).length }
])

const filteredProjects = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return projects.value.filter(p => {
    if (filterStatus.value && p.status !== filterStatus.value) return false
    if (!q) return true
    return p.title.toLowerCase().includes(q) ||
      p.mentor.toLowerCase().includes(q) ||
      p.members.some(m => m.name.toLowerCase().includes(q))
  })
})

function isOverdue(p: Project) {
  if (!p.due_date || p.status === 'SELESAI') return false
  return p.due_date < new Date().toISOString().substring(0, 10)
}

function formatDate(iso: string) {
  if (!iso) return '-'
  const [y, m, d] = iso.split('-')
  return `${d}/${m}/${y}`
}

async function loadProjects() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axios.get('/api/admin/projects')
    projects.value = res.data.projects || []
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal memuat data proyek.'
  } finally {
    loading.value = false
  }
}

function openForm(p?: Project) {
  formError.value = ''
  form.value = p
    ? { ...p, memberIds: p.members.map(m => m.id) }
    : emptyForm()
  showForm.value = true
}

async function saveProject() {
  saving.value = true
  formError.value = ''
  const internById = new Map(interns.value.map((u: any) => [u.id, u.name]))
  const existingNames = new Map((projects.value.find(p => p.id === form.value.id)?.members || []).map(m => [m.id, m.name]))
  const payload = {
    title: form.value.title,
    description: form.value.description,
    mentor: form.value.mentor,
    // Keep the stored name for members whose account no longer appears in the intern list
    members: form.value.memberIds.map(id => ({ id, name: internById.get(id) || existingNames.get(id) || `Peserta #${id}` })),
    status: form.value.status,
    progress: form.value.progress,
    start_date: form.value.start_date,
    due_date: form.value.due_date
  }
  try {
    if (form.value.id) {
      await axios.put(`/api/admin/projects/${form.value.id}`, payload)
    } else {
      await axios.post('/api/admin/projects', payload)
    }
    showForm.value = false
    await loadProjects()
  } catch (err: any) {
    formError.value = err.response?.data?.error || 'Gagal menyimpan proyek.'
  } finally {
    saving.value = false
  }
}

async function removeProject(p: Project) {
  if (!confirm(`Hapus proyek "${p.title}"?`)) return
  try {
    await axios.delete(`/api/admin/projects/${p.id}`)
    await loadProjects()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Gagal menghapus proyek.'
  }
}

onMounted(async () => {
  await Promise.all([loadProjects(), authStore.fetchUsers()])
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
