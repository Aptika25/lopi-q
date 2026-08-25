<template>
  <AdminLayout>
    <div class="flex flex-col gap-6 w-full select-none">

      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Manajemen Call Taker</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Kelola data seluruh petugas Call Taker 112 lintas instansi/dinas Kabupaten Bulukumba.</p>
        </div>
        <button
          v-if="authStore.isAdmin"
          @click="openAddDialog"
          class="w-full sm:w-auto py-2.5 px-4 bg-rose-700 hover:bg-rose-800 text-white font-bold text-xs rounded-xl shadow-sm transition-all flex items-center justify-center gap-1.5 cursor-pointer border-0 shrink-0"
        >
          <span class="material-symbols-outlined text-[16px]">person_add</span>
          <span>Tambah Call Taker Baru</span>
        </button>
      </div>

      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" :class="['flex items-center gap-2.5 p-3.5 rounded-2xl text-xs font-semibold border w-full shadow-sm',
          toast.success ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-red-50 border-red-200 text-red-800'
        ]">
          <span class="material-symbols-outlined text-[18px] shrink-0">
            {{ toast.success ? 'check_circle' : 'error' }}
          </span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Search & Filter Controls -->
      <div class="bg-white border border-slate-200 rounded-3xl p-5 shadow-sm flex flex-col md:flex-row items-center justify-between gap-4">
        <!-- Search Input -->
        <div class="relative w-full md:w-80">
          <span class="absolute left-3.5 top-1/2 -translate-y-1/2 material-symbols-outlined text-slate-400 text-[18px]">search</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari Nama, NIP, atau Email..."
            class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
          />
        </div>

        <!-- Filter Instansi / Unit Kerja -->
        <div class="flex items-center gap-2 w-full md:w-auto">
          <span class="text-xs font-bold text-slate-500 shrink-0">Unit Kerja:</span>
          <select
            v-model="selectedUnit"
            class="w-full md:w-64 px-3 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 rounded-xl text-xs font-semibold focus:outline-none"
          >
            <option value="">Semua Unit Kerja / Instansi</option>
            <option v-for="unit in availableUnits" :key="unit" :value="unit">{{ unit }}</option>
          </select>
        </div>
      </div>

      <!-- Main Content Container -->
      <div class="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col gap-6">

        <!-- Loading State -->
        <div v-if="loading" class="flex flex-col items-center justify-center py-14 gap-3 text-slate-400">
          <span class="material-symbols-outlined text-[40px] animate-spin">sync</span>
          <span class="text-xs">Mengambil data petugas Call Taker...</span>
        </div>

        <div v-else class="w-full">
          <!-- Desktop Table View -->
          <div class="hidden md:block overflow-x-auto w-full border border-slate-200 rounded-2xl">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">
                  <th class="py-3.5 px-4 min-w-[260px] whitespace-nowrap">Nama Call Taker &amp; NIP.</th>
                  <th class="py-3.5 px-4">Email</th>
                  <th class="py-3.5 px-4">Unit Kerja / Instansi</th>
                  <th class="py-3.5 px-4">Jabatan</th>
                  <th class="py-3.5 px-4">Status 2FA</th>
                  <th class="py-3.5 px-4">Status Akun</th>
                  <th class="py-3.5 px-4 text-center w-28">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-xs text-slate-700">
                <tr v-for="ct in filteredCallTakers" :key="ct.id" class="hover:bg-slate-50/60 transition-colors">
                  <!-- Avatar & Identity (Name & NIP) -->
                  <td class="py-3.5 px-4 min-w-[260px] whitespace-nowrap">
                    <div class="flex items-center gap-3">
                      <div class="h-9 w-9 rounded-xl shrink-0 font-black text-xs flex items-center justify-center shadow-sm bg-rose-50 text-rose-800 border border-rose-200">
                        {{ ct.name?.charAt(0).toUpperCase() }}
                      </div>
                      <div class="min-w-0">
                        <div class="font-bold text-slate-900 leading-tight whitespace-nowrap">{{ ct.name }}</div>
                        <div class="font-mono text-[10px] text-slate-400 mt-0.5 whitespace-nowrap">NIP. {{ ct.nip || '-' }}</div>
                      </div>
                    </div>
                  </td>

                  <!-- Email -->
                  <td class="py-3.5 px-4">
                    <div class="font-mono text-[11px] text-slate-600">{{ ct.email }}</div>
                  </td>

                  <!-- Unit Kerja -->
                  <td class="py-3.5 px-4">
                    <span class="inline-flex items-center px-2.5 py-1 rounded-lg bg-slate-100 border border-slate-200 text-slate-700 text-[11px] font-semibold">
                      {{ ct.unit_kerja || 'Diskominfo Bulukumba' }}
                    </span>
                  </td>

                  <!-- Jabatan -->
                  <td class="py-3.5 px-4 text-slate-600 font-medium text-[11px]">
                    {{ ct.jabatan || 'OPERATOR LAYANAN OPERASIONAL' }}
                  </td>

                  <!-- Status 2FA Badge -->
                  <td class="py-3.5 px-4 whitespace-nowrap">
                    <span :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-extrabold border shadow-2xs',
                      ct.totp_enabled 
                        ? 'bg-emerald-50 border-emerald-200/80 text-emerald-700' 
                        : 'bg-amber-50 border-amber-200/80 text-amber-700']">
                      <span class="material-symbols-outlined text-[14px]" :class="ct.totp_enabled ? 'text-emerald-600' : 'text-amber-600'">
                        {{ ct.totp_enabled ? 'verified_user' : 'shield_lock' }}
                      </span>
                      <span>{{ ct.totp_enabled ? '2FA Terverifikasi' : 'Belum Setup 2FA' }}</span>
                    </span>
                  </td>

                  <!-- Active Toggle -->
                  <td class="py-3.5 px-4">
                    <div class="flex items-center gap-2">
                      <button
                        type="button"
                        @click="handleToggleActive(ct)"
                        :disabled="!authStore.isAdmin || togglingId === ct.id"
                        :class="[
                          ct.is_active ? 'bg-emerald-500' : 'bg-slate-200',
                          'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed'
                        ]"
                      >
                        <span :class="[ct.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                      </button>
                      <span :class="['text-[10px] font-bold', ct.is_active ? 'text-emerald-600' : 'text-slate-400']">
                        {{ ct.is_active ? 'Aktif' : 'Nonaktif' }}
                      </span>
                    </div>
                  </td>

                  <!-- Action Buttons -->
                  <td class="py-3.5 px-4 text-center">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        @click="openEditDialog(ct)"
                        class="p-1.5 text-amber-600 hover:text-amber-700 hover:bg-amber-50 rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Edit Data Call Taker"
                      >
                        <span class="material-symbols-outlined text-[16px]">edit</span>
                      </button>
                      <button
                        v-if="authStore.isAdmin && ct.totp_enabled"
                        @click="handleReset2FA(ct)"
                        class="p-1.5 text-rose-600 hover:text-rose-700 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Reset 2FA Call Taker"
                      >
                        <span class="material-symbols-outlined text-[16px]">lock_reset</span>
                      </button>
                    </div>
                  </td>
                </tr>
                <tr v-if="filteredCallTakers.length === 0">
                  <td colspan="7" class="py-12 text-center text-slate-400 font-medium">Tidak ada data petugas Call Taker yang sesuai.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Mobile Card List View -->
          <div class="md:hidden flex flex-col gap-3 w-full">
            <div v-for="ct in filteredCallTakers" :key="ct.id" class="p-4 bg-slate-50 border border-slate-200/60 rounded-2xl flex flex-col gap-3">
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-2.5">
                  <div class="h-9 w-9 rounded-xl font-black text-xs flex items-center justify-center shrink-0 bg-rose-50 text-rose-800 border border-rose-200">
                    {{ ct.name?.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <div class="font-bold text-slate-900 text-xs leading-tight">{{ ct.name }}</div>
                    <div class="font-mono text-[10px] text-rose-700 font-semibold">NIP. {{ ct.nip || '-' }}</div>
                  </div>
                </div>
                <span :class="['inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[9px] font-bold border',
                  ct.totp_enabled ? 'bg-emerald-50 border-emerald-200 text-emerald-700' : 'bg-amber-50 border-amber-200 text-amber-700']">
                  {{ ct.totp_enabled ? '2FA Aktif' : 'Belum 2FA' }}
                </span>
              </div>

              <div class="flex flex-col gap-1 text-[11px] border-t border-slate-200/60 pt-2">
                <div class="flex justify-between">
                  <span class="text-[10px] text-slate-400 font-semibold">Email:</span>
                  <span class="font-mono text-slate-600">{{ ct.email }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-[10px] text-slate-400 font-semibold">Unit Kerja:</span>
                  <span class="text-slate-700 font-bold text-right max-w-[170px] truncate">{{ ct.unit_kerja || 'Diskominfo' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-[10px] text-slate-400 font-semibold">Jabatan:</span>
                  <span class="text-slate-600 text-right max-w-[170px] truncate">{{ ct.jabatan || '-' }}</span>
                </div>
              </div>

              <div class="flex items-center justify-between border-t border-slate-200/60 pt-3">
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    @click="handleToggleActive(ct)"
                    :disabled="!authStore.isAdmin || togglingId === ct.id"
                    :class="[ct.is_active ? 'bg-emerald-500' : 'bg-slate-200', 'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed']"
                  >
                    <span :class="[ct.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                  </button>
                  <span :class="['text-[10px] font-bold', ct.is_active ? 'text-emerald-600' : 'text-slate-400']">
                    {{ ct.is_active ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </div>
                <div class="flex items-center gap-1.5">
                  <button @click="openEditDialog(ct)"
                    class="p-2 text-amber-600 hover:bg-amber-50 rounded-xl transition-all cursor-pointer border border-amber-200/40 bg-white"
                    title="Edit"
                  ><span class="material-symbols-outlined text-[16px]">edit</span></button>
                  <button v-if="authStore.isAdmin && ct.totp_enabled" @click="handleReset2FA(ct)"
                    class="p-2 text-rose-600 hover:bg-rose-50 rounded-xl transition-all cursor-pointer border border-rose-200/40 bg-white flex items-center justify-center"
                    title="Reset 2FA"
                  ><span class="material-symbols-outlined text-[16px]">lock_reset</span></button>
                </div>
              </div>
            </div>
            <div v-if="filteredCallTakers.length === 0" class="py-12 text-center text-slate-400 font-medium text-xs">Tidak ada data Call Taker.</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===== ADD / EDIT MODAL ===== -->
    <transition name="fade">
      <div v-if="dialogOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click.self="closeDialog">
        <div class="w-full max-w-lg bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 flex flex-col gap-5 max-h-[90vh] overflow-y-auto">

          <div class="flex justify-between items-center border-b border-slate-100 pb-3">
            <h3 class="font-display font-black text-slate-900 text-base flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">{{ isEdit ? 'edit_square' : 'person_add' }}</span>
              {{ isEdit ? 'Edit Data Call Taker' : 'Tambah Call Taker Baru' }}
            </h3>
            <button @click="closeDialog" class="p-1 rounded-full hover:bg-slate-100 text-slate-400 border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-[18px]">close</span>
            </button>
          </div>

          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-red-50 border border-red-200 text-red-800 text-xs flex items-start gap-2.5">
            <span class="material-symbols-outlined text-[16px] shrink-0 mt-0.5">error</span>
            <span>{{ errorMessage }}</span>
          </div>

          <form @submit.prevent="submitForm" class="flex flex-col gap-4">

            <!-- Nama -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Nama Lengkap <span class="text-rose-500">*</span></label>
              <input type="text" v-model="form.name" required placeholder="Contoh: A.Mappalua, S.Pd"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- NIP -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">NIP <span class="text-rose-500">*</span></label>
              <input type="text" v-model="form.nip" required placeholder="Contoh: 199405032025211138"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Email -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Email Dinas <span class="text-rose-500">*</span></label>
              <input type="email" v-model="form.email" required placeholder="Contoh: mappalua@bulukumbakab.go.id"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Unit Kerja / Instansi -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Unit Kerja / Instansi <span class="text-rose-500">*</span></label>
              <select v-model="form.unit_kerja" required class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 rounded-xl text-xs font-semibold focus:outline-none">
                <option value="" disabled>Pilih Unit Kerja / Instansi</option>
                <option v-for="unit in availableUnits" :key="unit" :value="unit">{{ unit }}</option>
              </select>
            </div>

            <!-- Jabatan -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Jabatan</label>
              <input type="text" v-model="form.jabatan" placeholder="Contoh: PENATA LAYANAN OPERASIONAL"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Password Field -->
            <div class="flex flex-col gap-1.5">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">
                {{ isEdit ? 'Reset Password Baru' : 'Password Akses' }}
                <span v-if="!isEdit" class="text-rose-500">*</span>
              </label>
              <div class="flex gap-2">
                <div class="relative flex-grow">
                  <input
                    :type="showPassword ? 'text' : 'password'"
                    v-model="form.password"
                    :required="!isEdit"
                    :placeholder="isEdit ? 'Biarkan kosong jika tidak ingin diubah' : 'Minimal 6 karakter'"
                    class="w-full pl-4 pr-10 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
                  />
                  <button type="button" @click="showPassword = !showPassword"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-rose-600 border-0 bg-transparent cursor-pointer"
                  >
                    <span class="material-symbols-outlined text-[18px]">{{ showPassword ? 'visibility' : 'visibility_off' }}</span>
                  </button>
                </div>
                <button type="button" @click="generatePassword"
                  class="py-2.5 px-3 bg-slate-100 hover:bg-slate-200 border border-slate-200 rounded-xl text-xs font-bold transition-all cursor-pointer flex items-center gap-1 shrink-0 text-rose-700"
                >
                  <span class="material-symbols-outlined text-[16px]">vpn_key</span>
                  <span>Acak</span>
                </button>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-3 justify-end mt-2 pt-3 border-t border-slate-100">
              <button type="button" @click="closeDialog"
                class="py-2.5 px-5 border border-slate-200 hover:bg-slate-50 transition-colors text-xs rounded-xl cursor-pointer bg-white font-semibold text-slate-600"
              >Batal</button>
              <button type="submit" :disabled="submitLoading"
                class="py-2.5 px-5 bg-gradient-to-r from-rose-700 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white text-xs rounded-xl transition-all flex items-center gap-1.5 cursor-pointer border-0 font-bold shadow-md disabled:opacity-60"
              >
                <span v-if="submitLoading" class="animate-spin material-symbols-outlined text-[16px]">sync</span>
                <span>{{ isEdit ? 'Simpan Perubahan' : 'Tambah Call Taker' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ===== CONFIRM MODAL (Reset 2FA) ===== -->
    <div v-if="confirmModal.show"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm"
      @click.self="confirmModal.show = false"
    >
      <div class="bg-white rounded-[28px] max-w-[360px] w-full p-6 text-center border border-slate-100 shadow-2xl relative overflow-hidden flex flex-col items-center">
        <div class="relative flex items-center justify-center w-16 h-16 rounded-2xl bg-rose-50 border border-rose-100 mb-4 shadow-sm">
          <span class="material-symbols-outlined text-[28px] text-rose-500">lock_reset</span>
        </div>
        <h3 class="font-display text-base font-black text-slate-900 mb-2 tracking-tight">{{ confirmModal.title }}</h3>
        <p class="text-xs text-slate-500 leading-relaxed px-2 mb-6">{{ confirmModal.description }}</p>
        <div class="flex gap-3 w-full">
          <button type="button" @click="confirmModal.show = false"
            class="flex-grow py-3 px-4 rounded-xl border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs font-bold transition-all cursor-pointer bg-white flex items-center justify-center gap-1.5"
          >
            <span class="material-symbols-outlined text-[15px]">arrow_back</span>Batal
          </button>
          <button type="button" @click="confirmModal.onConfirm"
            class="flex-grow py-3 px-4 rounded-xl bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition-all cursor-pointer border-0 shadow-sm flex items-center justify-center gap-1.5"
          >
            <span class="material-symbols-outlined text-[15px]">lock_reset</span>Reset 2FA
          </button>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AdminLayout from '@/layouts/AdminLayout.vue'

const authStore = useAuthStore()

// ===== STATE =====
const callTakers = ref<any[]>([])
const loading = ref(false)
const togglingId = ref<number | null>(null)
const submitLoading = ref(false)
const dialogOpen = ref(false)
const isEdit = ref(false)
const editUserId = ref<number | null>(null)
const errorMessage = ref('')
const showPassword = ref(false)
const searchQuery = ref('')
const selectedUnit = ref('')

const form = ref({
  name: '',
  nip: '',
  email: '',
  unit_kerja: '',
  jabatan: '',
  password: ''
})

const toast = ref({ show: false, success: true, message: '' })

const confirmModal = ref({
  show: false,
  title: '',
  description: '',
  onConfirm: () => {}
})

// ===== UNITS LIST =====
const availableUnits = [
  'Dinas Sosial',
  'Badan Penanggulangan Bencana Daerah',
  'Dinas Kesehatan',
  'Dinas Perhubungan',
  'Satpol, Pemadam Kebakaran dan Penyelamatan',
  'Diskominfo Kab. Bulukumba'
]

// ===== FILTERED COMPUTED =====
const filteredCallTakers = computed(() => {
  return callTakers.value.filter((ct: any) => {
    const matchesSearch = !searchQuery.value || 
      ct.name?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      ct.nip?.includes(searchQuery.value) ||
      ct.email?.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesUnit = !selectedUnit.value || ct.unit_kerja === selectedUnit.value

    return matchesSearch && matchesUnit
  })
})

// ===== TOAST =====
const showToast = (success: boolean, message: string) => {
  toast.value = { show: true, success, message }
  setTimeout(() => { toast.value.show = false }, 4000)
}

// ===== CONFIRM MODAL =====
const triggerConfirm = (title: string, description: string, onConfirm: () => void) => {
  confirmModal.value = {
    show: true,
    title,
    description,
    onConfirm: () => {
      onConfirm()
      confirmModal.value.show = false
    }
  }
}

// ===== LOAD CALL TAKERS =====
const loadCallTakers = async () => {
  loading.value = true
  try {
    await authStore.fetchUsers()
    callTakers.value = (authStore.usersList || [])
      .filter((u: any) => u.role === 'call_taker')
      .sort((a: any, b: any) => a.name.localeCompare(b.name))
  } catch (err) {
    showToast(false, 'Gagal memuat data petugas Call Taker.')
  } finally {
    loading.value = false
  }
}

onMounted(() => loadCallTakers())

// ===== DIALOG OPEN/CLOSE =====
const openAddDialog = () => {
  isEdit.value = false
  editUserId.value = null
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    name: '', nip: '', email: '',
    unit_kerja: 'Dinas Sosial', jabatan: 'PENATA LAYANAN OPERASIONAL',
    password: ''
  }
  dialogOpen.value = true
}

const openEditDialog = (ct: any) => {
  isEdit.value = true
  editUserId.value = ct.id
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    name: ct.name,
    nip: ct.nip || '',
    email: ct.email,
    unit_kerja: ct.unit_kerja || 'Diskominfo Kab. Bulukumba',
    jabatan: ct.jabatan || '',
    password: ''
  }
  dialogOpen.value = true
}

const closeDialog = () => {
  dialogOpen.value = false
  errorMessage.value = ''
  showPassword.value = false
}

// ===== PASSWORD GENERATOR =====
const generatePassword = () => {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%'
  let pass = ''
  for (let i = 0; i < 10; i++) pass += chars.charAt(Math.floor(Math.random() * chars.length))
  form.value.password = pass
  showPassword.value = true
}

// ===== TOGGLE ACTIVE =====
const handleToggleActive = async (ct: any) => {
  const newState = !ct.is_active
  togglingId.value = ct.id
  try {
    await authStore.toggleUserActive(ct.id, newState)
    ct.is_active = newState
    showToast(true, `Akun ${ct.name} berhasil ${newState ? 'diaktifkan' : 'dinonaktifkan'}.`)
  } catch (err) {
    showToast(false, 'Gagal mengubah status keaktifan akun.')
  } finally {
    togglingId.value = null
  }
}

// ===== RESET 2FA =====
const handleReset2FA = (ct: any) => {
  triggerConfirm(
    'Reset Keamanan 2FA?',
    `Apakah Anda yakin ingin menonaktifkan Google Authenticator (2FA) untuk petugas Call Taker ${ct.name}? Petugas harus melakukan scan ulang QR 2FA saat login berikutnya.`,
    async () => {
      submitLoading.value = true
      try {
        await authStore.resetUser2fa(ct.id)
        showToast(true, `2FA untuk ${ct.name} berhasil dinonaktifkan.`)
        await loadCallTakers()
      } catch (err) {
        showToast(false, 'Gagal mereset 2FA Call Taker.')
      } finally {
        submitLoading.value = false
      }
    }
  )
}

// ===== SUBMIT FORM =====
const submitForm = async () => {
  submitLoading.value = false
  errorMessage.value = ''

  if (!form.value.name.trim()) { errorMessage.value = 'Nama lengkap wajib diisi.'; return }
  if (!form.value.nip.trim()) { errorMessage.value = 'NIP wajib diisi.'; return }
  if (!isEdit.value) {
    if (!form.value.email || !form.value.password) { errorMessage.value = 'Email dan Password wajib diisi.'; return }
    if (form.value.password.length < 6) { errorMessage.value = 'Password minimal 6 karakter.'; return }
  }

  submitLoading.value = true
  try {
    if (isEdit.value && editUserId.value !== null) {
      const res = await authStore.updateUser(editUserId.value, {
        name: form.value.name,
        nip: form.value.nip,
        email: form.value.email,
        unit_kerja: form.value.unit_kerja,
        jabatan: form.value.jabatan,
        role: 'call_taker',
        password: form.value.password
      })
      if (res.success) {
        showToast(true, 'Data Call Taker berhasil diperbarui.')
        closeDialog()
        await loadCallTakers()
      } else {
        errorMessage.value = res.error || 'Gagal memperbarui data Call Taker.'
      }
    } else {
      const res = await authStore.createUser({
        name: form.value.name,
        nip: form.value.nip,
        email: form.value.email,
        password: form.value.password,
        unit_kerja: form.value.unit_kerja,
        jabatan: form.value.jabatan,
        role: 'call_taker'
      })
      if (res.success) {
        showToast(true, 'Call Taker baru berhasil ditambahkan!')
        closeDialog()
        await loadCallTakers()
      } else {
        errorMessage.value = res.error || 'Gagal menambahkan Call Taker.'
      }
    }
  } catch (err: any) {
    errorMessage.value = authStore.error || 'Terjadi kesalahan saat memproses permintaan.'
  } finally {
    submitLoading.value = false
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1; }
</style>
