<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">
      
      <!-- Page Header -->
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">group</span>
            Manajemen Peserta Magang
          </h1>
          <p class="text-sm text-[#574146] mt-1">Kelola data peserta magang, akun Gmail, jurusan, dan asal sekolah/universitas.</p>
        </div>
      </div>

      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" :class="['flex items-center gap-2.5 p-3.5 rounded-xl text-xs font-semibold border w-full shadow-xs',
          toast.success ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-rose-50 border-rose-200 text-rose-800'
        ]">
          <span class="material-symbols-outlined text-[18px] shrink-0">
            {{ toast.success ? 'check_circle' : 'error' }}
          </span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Action Bar (Glass Card Layout) -->
      <div class="bg-white/85 backdrop-blur-md rounded-xl border border-[#F8BBD0] p-4 shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="flex items-center w-full sm:w-auto">
          <!-- Search Input -->
          <div class="relative w-full sm:w-80 group">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] group-focus-within:text-[#ab2c5d] transition-colors text-sm">search</span>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Cari nama, NISN, Gmail, jurusan, atau sekolah..." 
              class="w-full pl-10 pr-4 py-2.5 rounded-lg border border-[#F8BBD0] bg-white text-[#1b1c1c] focus:outline-none focus:ring-2 focus:ring-[#f06292]/30 focus:border-[#f06292] transition-all text-xs font-medium"
            />
          </div>
        </div>

        <!-- Add User Button -->
        <button 
          @click="openAddDialog"
          class="bg-[#F06292] hover:bg-[#ab2c5d] text-white font-bold text-xs px-6 py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 shadow-[0px_10px_30px_rgba(240,98,146,0.15)] w-full sm:w-auto shrink-0 border-0 cursor-pointer active:scale-95"
        >
          <span class="material-symbols-outlined text-base">add</span>
          <span>Tambah Peserta Magang</span>
        </button>
      </div>

      <!-- Data Table Card -->
      <div class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        
        <!-- Loading State -->
        <div v-if="loading" class="flex flex-col items-center justify-center py-14 gap-3 text-[#8a7176]">
          <span class="material-symbols-outlined text-[40px] animate-spin">sync</span>
          <span class="text-xs font-medium">Mengambil data peserta magang...</span>
        </div>

        <div v-else class="w-full">
          <div class="overflow-x-auto w-full">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0]">
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Pengguna (Nama &amp; NISN)</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Gmail / Email</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Jurusan</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Asal Sekolah / Universitas</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Status 2FA</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Status Akun</th>
                  <th class="py-4 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="text-xs text-[#1b1c1c] bg-white divide-y divide-[#F8BBD0]">
                <tr 
                  v-for="ct in filteredCallTakers" 
                  :key="ct.id" 
                  class="border-b border-[#F8BBD0] hover:bg-[#FCE4EC]/40 transition-colors"
                >
                  <!-- User Avatar & Identity (Nama & NISN/NIM) -->
                  <td class="py-4 px-6">
                    <div class="flex items-center gap-3">
                      <div class="h-10 w-10 rounded-full bg-[#fec1d6] overflow-hidden flex-shrink-0 flex items-center justify-center text-[#ab2c5d] font-bold text-sm shadow-xs border border-[#ddbfc5]">
                        {{ ct.name?.charAt(0).toUpperCase() }}
                      </div>
                      <div>
                        <div class="font-bold text-[#1b1c1c] text-sm">{{ ct.name }}</div>
                        <div class="font-mono text-[11px] text-[#8a7176]">NISN: {{ ct.nip || '-' }}</div>
                      </div>
                    </div>
                  </td>

                  <!-- Gmail / Email -->
                  <td class="py-4 px-6 text-[#574146] font-mono text-xs">
                    {{ ct.email }}
                  </td>

                  <!-- Jurusan & Asal Sekolah/University -->
                  <td class="py-4 px-6 text-[#574146]">
                    <div class="font-bold text-[#1b1c1c]">{{ ct.unit_kerja || 'Rekayasa Perangkat Lunak' }}</div>
                  </td>

                  <!-- Role / Jabatan Badge -->
                  <td class="py-4 px-6">
                    <div class="font-semibold text-[#1b1c1c] uppercase text-[11px]">{{ ct.jabatan || 'SMK Negeri 1 Bulukumba' }}</div>
                    <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[9px] font-bold uppercase tracking-wider bg-[#FCE4EC] text-[#F06292] border border-[#F06292]/20 mt-0.5">
                      INTERN
                    </span>
                  </td>

                  <!-- Status 2FA Badge -->
                  <td class="py-4 px-6 whitespace-nowrap">
                    <span 
                      v-if="ct.totp_enabled" 
                      class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-200"
                    >
                      <span class="material-symbols-outlined text-[14px] text-emerald-600">verified_user</span>
                      <span>2FA Terverifikasi</span>
                    </span>
                    <span 
                      v-else 
                      class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-bold bg-slate-100 text-slate-500 border border-slate-200"
                    >
                      <span class="material-symbols-outlined text-[14px] text-slate-400">gpp_maybe</span>
                      <span>Belum 2FA</span>
                    </span>
                  </td>

                  <!-- Status (Active / Inactive Toggle Switch) -->
                  <td class="py-4 px-6">
                    <div class="inline-flex items-center gap-2">
                      <button
                        type="button"
                        @click="handleToggleActive(ct)"
                        :disabled="togglingId === ct.id"
                        :title="ct.is_active ? 'Klik untuk menonaktifkan' : 'Klik untuk mengaktifkan'"
                        :class="[
                          ct.is_active ? 'bg-emerald-500' : 'bg-slate-300',
                          'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed'
                        ]"
                      >
                        <span :class="[ct.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                      </button>
                      <span :class="['text-xs font-bold', ct.is_active ? 'text-emerald-600' : 'text-slate-400']">
                        {{ ct.is_active ? 'Aktif' : 'Nonaktif' }}
                      </span>
                    </div>
                  </td>

                  <!-- Action Buttons (Edit & Reset 2FA) -->
                  <td class="py-4 px-6 text-right">
                    <div class="flex justify-end gap-2">
                      <button 
                        @click="openEditDialog(ct)"
                        class="p-2 text-[#574146] hover:text-[#ab2c5d] transition-colors rounded-lg hover:bg-[#ffd9e4] border-0 cursor-pointer bg-transparent"
                        title="Edit Data Peserta Magang"
                      >
                        <span class="material-symbols-outlined text-base">edit</span>
                      </button>
                      <!-- Tombol reset 2FA hanya muncul jika Status 2FA Terverifikasi (ct.totp_enabled === true) -->
                      <button 
                        v-if="ct.totp_enabled"
                        @click="handleReset2FA(ct)"
                        class="p-2 text-rose-600 hover:text-rose-800 transition-colors rounded-lg hover:bg-rose-50 border-0 cursor-pointer bg-transparent"
                        title="Reset 2FA Peserta Magang"
                      >
                        <span class="material-symbols-outlined text-base">lock_reset</span>
                      </button>
                    </div>
                  </td>
                </tr>

                <!-- Empty State -->
                <tr v-if="filteredCallTakers.length === 0">
                  <td colspan="7" class="py-6 text-center">
                    <div class="inline-flex items-center justify-center gap-1.5 text-xs text-slate-500 font-medium">
                      <span class="material-symbols-outlined text-slate-400" style="font-size: 14px !important;">info</span>
                      <span>Tidak ada data peserta magang yang ditemukan.</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination Container -->
          <div class="bg-white border-t border-[#F8BBD0] px-6 py-4 flex items-center justify-between">
            <span class="text-xs text-[#574146] font-medium">
              Menampilkan 1-{{ filteredCallTakers.length }} dari {{ callTakers.length }} pengguna
            </span>
            <div class="flex gap-2">
              <button disabled class="px-3 py-1.5 rounded border border-[#F8BBD0] text-[#574146] opacity-50 cursor-not-allowed bg-white">
                <span class="material-symbols-outlined text-sm">chevron_left</span>
              </button>
              <button class="px-3 py-1.5 rounded border border-[#F06292] bg-[#FCE4EC] text-[#ab2c5d] font-bold text-xs">1</button>
              <button disabled class="px-3 py-1.5 rounded border border-[#F8BBD0] text-[#574146] opacity-50 cursor-not-allowed bg-white">
                <span class="material-symbols-outlined text-sm">chevron_right</span>
              </button>
            </div>
          </div>
        </div>

      </div>

    </div>

    <!-- ===== ADD / EDIT DIALOG MODAL ===== -->
    <transition name="fade">
      <div v-if="dialogOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="closeDialog">
        <div class="bg-white rounded-2xl max-w-lg w-full p-6 border border-[#ddbfc5] shadow-2xl space-y-4">
          <div class="flex justify-between items-center border-b border-[#ddbfc5]/60 pb-3">
            <h3 class="font-bold text-base text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">{{ isEdit ? 'edit' : 'person_add' }}</span>
              <span>{{ isEdit ? 'Edit Data Peserta Magang' : 'Tambah Peserta Magang Baru' }}</span>
            </h3>
            <button @click="closeDialog" class="text-slate-400 hover:text-slate-600 border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-xl">close</span>
            </button>
          </div>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-800 text-xs font-semibold">
            {{ errorMessage }}
          </div>

          <form @submit.prevent="submitForm" class="space-y-4 text-xs">
            <!-- Nama Lengkap -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">Nama Lengkap <span class="text-rose-500">*</span></label>
              <input v-model="form.name" type="text" placeholder="Contoh: Sarah Jenkins" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>

            <!-- NISN / NIM -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">NISN / NIM <span class="text-rose-500">*</span></label>
              <input v-model="form.nip" type="text" placeholder="Contoh: 0051234567 atau 2024001" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>

            <!-- Gmail / Email -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">Gmail / Email Akses <span class="text-rose-500">*</span></label>
              <input v-model="form.email" type="email" placeholder="sarah.j@gmail.com" required :disabled="isEdit" class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292] disabled:bg-slate-100" />
            </div>

            <!-- Jurusan -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">Jurusan <span class="text-rose-500">*</span></label>
              <input v-model="form.unit_kerja" type="text" placeholder="Contoh: Rekayasa Perangkat Lunak / Teknik Informatika" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>

            <!-- Asal Sekolah / University -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">Asal Sekolah / University <span class="text-rose-500">*</span></label>
              <input v-model="form.jabatan" type="text" placeholder="Contoh: SMK Negeri 1 Bulukumba / Universitas Negeri Makassar" required class="w-full px-3.5 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" />
            </div>

            <!-- Password Akses -->
            <div class="space-y-1">
              <label class="font-bold text-[#574146] uppercase text-[10px]">
                {{ isEdit ? 'Reset Password (Opsional)' : 'Password Akses' }} <span v-if="!isEdit" class="text-rose-500">*</span>
              </label>
              <div class="flex gap-2">
                <div class="relative flex-grow">
                  <input 
                    :type="showPassword ? 'text' : 'password'" 
                    v-model="form.password" 
                    :required="!isEdit" 
                    placeholder="Minimal 6 karakter" 
                    class="w-full pl-3.5 pr-10 py-2 border border-[#ddbfc5] rounded-lg focus:outline-none focus:border-[#f06292]" 
                  />
                  <button 
                    type="button" 
                    @click="showPassword = !showPassword"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-[#ab2c5d] border-0 bg-transparent cursor-pointer flex items-center justify-center p-0.5 transition-colors"
                    :title="showPassword ? 'Sembunyikan Password' : 'Tampilkan Password'"
                  >
                    <span class="material-symbols-outlined text-[18px]">
                      {{ showPassword ? 'visibility_off' : 'visibility' }}
                    </span>
                  </button>
                </div>
                <button type="button" @click="generatePassword" class="px-3 py-2 bg-slate-100 border border-slate-200 hover:bg-slate-200 transition-colors rounded-lg font-bold text-xs text-[#ab2c5d] shrink-0 cursor-pointer">Acak</button>
              </div>
            </div>

            <div class="flex justify-end gap-3 pt-3 border-t border-[#ddbfc5]/60">
              <button type="button" @click="closeDialog" class="px-4 py-2 border border-[#ddbfc5] text-[#574146] rounded-lg font-bold bg-white cursor-pointer">Batal</button>
              <button type="submit" :disabled="submitLoading" class="px-5 py-2 bg-[#F06292] hover:bg-[#ab2c5d] text-white rounded-lg font-bold cursor-pointer border-0 shadow-xs">
                {{ isEdit ? 'Simpan Perubahan' : 'Tambah Peserta Magang' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ===== CONFIRM MODAL (Reset 2FA) ===== -->
    <transition name="fade">
      <div v-if="confirmModal.show" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs" @click.self="confirmModal.show = false">
        <div class="bg-white rounded-2xl max-w-sm w-full p-6 text-center border border-[#ddbfc5] shadow-2xl space-y-4">
          <div class="w-14 h-14 rounded-full bg-rose-50 text-rose-600 flex items-center justify-center mx-auto border border-rose-100">
            <span class="material-symbols-outlined text-2xl">lock_reset</span>
          </div>
          <h3 class="font-bold text-base text-[#1b1c1c]">{{ confirmModal.title }}</h3>
          <p class="text-xs text-[#574146] leading-relaxed">{{ confirmModal.description }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="confirmModal.show = false" class="flex-1 py-2.5 border border-[#ddbfc5] rounded-lg text-xs font-bold text-[#574146] bg-white cursor-pointer">Batal</button>
            <button type="button" @click="confirmModal.onConfirm" class="flex-1 py-2.5 bg-rose-600 hover:bg-rose-700 text-white rounded-lg text-xs font-bold border-0 cursor-pointer shadow-xs">Reset 2FA</button>
          </div>
        </div>
      </div>
    </transition>
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

// ===== FILTERED COMPUTED =====
const filteredCallTakers = computed(() => {
  if (!searchQuery.value) return callTakers.value
  const q = searchQuery.value.toLowerCase()
  return callTakers.value.filter((ct: any) => {
    return (
      ct.name?.toLowerCase().includes(q) ||
      ct.nip?.toLowerCase().includes(q) ||
      ct.email?.toLowerCase().includes(q) ||
      ct.unit_kerja?.toLowerCase().includes(q) ||
      ct.jabatan?.toLowerCase().includes(q)
    )
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

// ===== LOAD PESERTA MAGANGS =====
const loadCallTakers = async () => {
  loading.value = true
  try {
    await authStore.fetchUsers()
    const allUsers = authStore.usersList || []
    callTakers.value = allUsers.filter((u: any) => {
      const role = (u.role || u.Role || '').toLowerCase()
      return role === 'intern' || role === 'peserta'
    })
  } catch (err) {
    showToast(false, 'Gagal memuat data peserta magang.')
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
    name: '',
    nip: '',
    email: '',
    unit_kerja: 'Rekayasa Perangkat Lunak',
    jabatan: 'SMK Negeri 1 Bulukumba',
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
    unit_kerja: ct.unit_kerja || 'Rekayasa Perangkat Lunak',
    jabatan: ct.jabatan || 'SMK Negeri 1 Bulukumba',
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
    showToast(true, `Status peserta magang ${ct.name} berhasil ${newState ? 'diaktifkan' : 'dinonaktifkan'}.`)
  } catch (err) {
    showToast(false, 'Gagal mengubah status peserta magang.')
  } finally {
    togglingId.value = null
  }
}

// ===== RESET 2FA =====
const handleReset2FA = (ct: any) => {
  triggerConfirm(
    'Reset Keamanan 2FA?',
    `Apakah Anda yakin ingin menonaktifkan Google Authenticator (2FA) untuk peserta magang ${ct.name}? Status 2FA akan ter-reset kembali menjadi Belum 2FA.`,
    async () => {
      submitLoading.value = true
      try {
        await authStore.resetUser2fa(ct.id)
        ct.totp_enabled = false
        ct.totp_secret = ''

        // Update local array & authStore
        const target = (callTakers.value || []).find((u: any) => u.id === ct.id || u.email === ct.email)
        if (target) {
          target.totp_enabled = false
          target.totp_secret = ''
        }

        const storeTarget = (authStore.usersList || []).find((u: any) => u.id === ct.id || u.email === ct.email)
        if (storeTarget) {
          storeTarget.totp_enabled = false
          storeTarget.totp_secret = ''
        }

        showToast(true, `2FA untuk ${ct.name} berhasil di-reset menjadi Belum 2FA.`)
        await loadCallTakers()
      } catch (err) {
        showToast(false, 'Gagal mereset 2FA Peserta Magang.')
      } finally {
        submitLoading.value = false
      }
    }
  )
}

// ===== SUBMIT FORM =====
const submitForm = async () => {
  errorMessage.value = ''

  if (!form.value.name.trim()) { errorMessage.value = 'Nama lengkap wajib diisi.'; return }
  if (!form.value.nip.trim()) { errorMessage.value = 'NISN / NIM wajib diisi.'; return }
  if (!isEdit.value) {
    if (!form.value.email || !form.value.password) { errorMessage.value = 'Gmail dan Password wajib diisi.'; return }
    if (form.value.password.length < 6) { errorMessage.value = 'Password minimal 6 karakter.'; return }
  }

  submitLoading.value = true
  try {
    if (isEdit.value && editUserId.value !== null) {
      const res = await authStore.updateUser(editUserId.value, {
        name: form.value.name.trim(),
        nip: form.value.nip.trim(),
        email: form.value.email.trim(),
        unit_kerja: form.value.unit_kerja.trim() || 'Rekayasa Perangkat Lunak',
        jabatan: form.value.jabatan.trim() || 'SMK Negeri 1 Bulukumba',
        role: 'intern',
        password: form.value.password
      })
      if (res && res.success !== false) {
        showToast(true, `Data Peserta Magang ${form.value.name} berhasil diperbarui.`)
        closeDialog()
        await loadCallTakers()
      } else {
        errorMessage.value = res?.error || authStore.error || 'Gagal memperbarui data Peserta Magang.'
      }
    } else {
      const res = await authStore.createUser({
        name: form.value.name.trim(),
        nip: form.value.nip.trim(),
        email: form.value.email.trim(),
        password: form.value.password,
        unit_kerja: form.value.unit_kerja.trim() || 'Rekayasa Perangkat Lunak',
        jabatan: form.value.jabatan.trim() || 'SMK Negeri 1 Bulukumba',
        role: 'intern'
      })
      if (res && res.success !== false) {
        // Instant optimism: push to local array and authStore
        const createdObj = res.user || {
          id: Date.now(),
          name: form.value.name.trim(),
          nip: form.value.nip.trim(),
          email: form.value.email.trim(),
          unit_kerja: form.value.unit_kerja.trim() || 'Rekayasa Perangkat Lunak',
          jabatan: form.value.jabatan.trim() || 'SMK Negeri 1 Bulukumba',
          role: 'intern',
          is_active: true
        }
        callTakers.value = callTakers.value.filter((u: any) => u.email !== createdObj.email && (!createdObj.nip || u.nip !== createdObj.nip))
        callTakers.value.unshift(createdObj)

        // Also push to authStore.usersList immediately for Dashboard sync
        if (!authStore.usersList.some((u: any) => u.email === createdObj.email || (u.nip && u.nip === createdObj.nip))) {
          authStore.usersList.unshift(createdObj)
        }

        showToast(true, `Peserta Magang ${form.value.name} berhasil ditambahkan!`)
        closeDialog()
      } else {
        errorMessage.value = res?.error || authStore.error || 'Gagal menambahkan Peserta Magang.'
      }
    }
  } catch (err: any) {
    errorMessage.value = err.message || authStore.error || 'Terjadi kesalahan saat menambahkan peserta magang.'
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






