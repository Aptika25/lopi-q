<template>
  <div class="bg-white rounded-3xl space-y-6">
    <!-- Section Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-slate-200">
      <div>
        <h2 class="text-xl font-display font-black text-slate-900 flex items-center gap-2">
          <span class="material-symbols-outlined text-rose-700 text-[24px]">manage_accounts</span>
          Manajemen Peserta Magang & Akun User
        </h2>
        <p class="text-xs text-slate-500 mt-1">Kelola data Peserta Magangs NTPD 112 Kabupaten Bulukumba, atur password, dan reset 2FA.</p>
      </div>

      <button 
        @click="openAddModal"
        class="px-4 py-2.5 bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-bold text-xs rounded-xl shadow-md transition-all flex items-center justify-center gap-1.5 cursor-pointer border-0"
      >
        <span class="material-symbols-outlined text-[18px]">person_add</span>
        <span>Tambah Peserta Magang Baru</span>
      </button>
    </div>

    <!-- Search & Filters -->
    <div class="flex items-center space-x-3">
      <div class="relative w-full">
        <span class="material-symbols-outlined absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-[18px]">search</span>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Cari berdasarkan NIP, Email Dinas, Nama, Jabatan, atau Unit Kerja..." 
          class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white transition-all shadow-xs"
        />
      </div>
    </div>

    <!-- Users Table -->
    <div class="overflow-x-auto rounded-2xl border border-slate-200 shadow-xs">
      <table class="w-full text-left text-xs text-slate-700">
        <thead class="bg-slate-100 text-slate-700 font-bold uppercase tracking-wider border-b border-slate-200">
          <tr>
            <th class="py-3.5 px-4">Peserta Magang / NIP / Email Dinas</th>
            <th class="py-3.5 px-4">Jabatan</th>
            <th class="py-3.5 px-4">Unit Kerja</th>
            <th class="py-3.5 px-4">Role</th>
            <th class="py-3.5 px-4">2FA Google</th>
            <th class="py-3.5 px-4">Status</th>
            <th class="py-3.5 px-4 text-right">Aksi Super Admin</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-200 bg-white">
          <tr v-for="u in filteredUsers" :key="u.id" class="hover:bg-slate-50 transition-colors">
            <td class="py-3.5 px-4">
              <div class="font-bold text-slate-900 text-sm">{{ u.name }}</div>
              <div class="font-mono text-[11px] text-rose-700 font-semibold mt-0.5">NIP. {{ u.nip || '-' }}</div>
              <div class="font-mono text-[11px] text-slate-500">Email: {{ u.email || '-' }}</div>
            </td>
            <td class="py-3.5 px-4 text-slate-700 font-medium">{{ u.jabatan || '-' }}</td>
            <td class="py-3.5 px-4 text-slate-700 font-medium">{{ u.unit_kerja || '-' }}</td>
            <td class="py-3.5 px-4 whitespace-nowrap">
              <span 
                class="inline-flex items-center px-3 py-1 rounded-full text-[10px] font-extrabold uppercase tracking-wider whitespace-nowrap border shadow-2xs"
                :class="{
                  'bg-rose-50 text-rose-700 border-rose-200': u.role === 'superadmin',
                  'bg-indigo-50 text-indigo-700 border-indigo-200': u.role === 'admin',
                  'bg-slate-100 text-slate-700 border-slate-200': u.role === 'intern' || u.role === 'intern'
                }"
              >
                {{ u.role === 'superadmin' ? 'SUPER ADMIN' : (u.role === 'admin' ? 'ADMIN' : 'PESERTA MAGANG') }}
              </span>
            </td>
            <td class="py-3.5 px-4">
              <span 
                class="px-2.5 py-0.5 rounded-full text-[10px] font-bold"
                :class="u.totp_enabled ? 'bg-emerald-100 text-emerald-800 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
              >
                {{ u.totp_enabled ? '✓ Aktif' : 'Non-aktif' }}
              </span>
            </td>
            <td class="py-3.5 px-4">
              <button 
                @click="toggleActive(u)"
                class="px-2.5 py-0.5 rounded-full text-[10px] font-bold transition-all cursor-pointer border"
                :class="u.is_active ? 'bg-emerald-100 text-emerald-800 border-emerald-200 hover:bg-emerald-200' : 'bg-rose-100 text-rose-800 border-rose-200 hover:bg-rose-200'"
              >
                {{ u.is_active ? 'Aktif' : 'Dinonaktifkan' }}
              </button>
            </td>
            <td class="py-3.5 px-4 text-right space-x-1.5">
              <button 
                @click="openEditModal(u)"
                class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded-xl text-[11px] font-bold border border-slate-300 transition-all cursor-pointer inline-flex items-center gap-1"
                title="Atur Password & Edit"
              >
                <span class="material-symbols-outlined text-[14px]">key</span>
                <span>Atur Password</span>
              </button>

              <button 
                v-if="u.totp_enabled"
                @click="reset2fa(u)"
                class="px-3 py-1.5 bg-amber-100 hover:bg-amber-200 text-amber-900 rounded-xl text-[11px] font-bold border border-amber-300 transition-all cursor-pointer inline-flex items-center gap-1"
                title="Reset 2FA Google Authenticator"
              >
                <span class="material-symbols-outlined text-[14px]">lock_reset</span>
                <span>Reset 2FA</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Edit/Add User Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-xs">
      <div class="bg-white w-full max-w-md rounded-3xl p-6 relative border border-slate-200 shadow-2xl">
        <h3 class="text-lg font-display font-black text-slate-900 mb-4 flex items-center gap-2">
          <span class="material-symbols-outlined text-rose-700">edit_square</span>
          <span>{{ isEditing ? `Edit Akun: ${form.name}` : 'Tambah Peserta Magang / Admin Baru' }}</span>
        </h3>

        <form @submit.prevent="saveUser" class="space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-slate-700 font-bold mb-1">NIP Pegawai</label>
              <input 
                v-model="form.nip" 
                type="text" 
                placeholder="199405032025211138" 
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl font-mono text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
              />
            </div>
            <div>
              <label class="block text-slate-700 font-bold mb-1">Email Dinas</label>
              <input 
                v-model="form.email" 
                type="email" 
                placeholder="nama@bulukumbakab.go.id" 
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl font-mono text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
              />
            </div>
          </div>

          <div>
            <label class="block text-slate-700 font-bold mb-1">Nama Lengkap & Gelar</label>
            <input 
              v-model="form.name" 
              type="text" 
              required 
              placeholder="Contoh: A.Mappalua, S.Pd" 
              class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-slate-700 font-bold mb-1">Jabatan</label>
              <input 
                v-model="form.jabatan" 
                type="text" 
                placeholder="PENATA LAYANAN OPERASIONAL" 
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
              />
            </div>
            <div>
              <label class="block text-slate-700 font-bold mb-1">Unit Kerja</label>
              <input 
                v-model="form.unit_kerja" 
                type="text" 
                placeholder="Dinas Sosial" 
                class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
              />
            </div>
          </div>

          <div>
            <label class="block text-slate-700 font-bold mb-1">Peran (Role)</label>
            <select v-model="form.role" class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white">
              <option value="intern">Peserta Magang</option>
              <option value="admin">Admin</option>
              <option value="superadmin">Super Admin</option>
            </select>
          </div>

          <div>
            <label class="block text-slate-700 font-bold mb-1">
              {{ isEditing ? 'Password Baru (Kosongkan jika tidak diubah)' : 'Password Awal' }}
            </label>
            <input 
              v-model="form.password" 
              type="password" 
              :required="!isEditing"
              placeholder="Masukkan password baru..." 
              class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 focus:outline-none focus:border-rose-600 focus:bg-white"
            />
          </div>

          <div class="flex justify-end space-x-2 pt-4 border-t border-slate-200">
            <button type="button" @click="showModal = false" class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl font-bold border border-slate-300 cursor-pointer">
              Batal
            </button>
            <button type="submit" class="px-5 py-2 bg-gradient-to-r from-rose-700 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white font-bold rounded-xl shadow-md cursor-pointer border-0">
              Simpan Data
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
const searchQuery = ref('');
const showModal = ref(false);
const isEditing = ref(false);
const editingId = ref(null);

const form = ref({
  nip: '',
  email: '',
  name: '',
  jabatan: '',
  unit_kerja: '',
  role: 'intern',
  password: ''
});

const filteredUsers = computed(() => {
  if (!searchQuery.value) return authStore.usersList;
  const q = searchQuery.value.toLowerCase();
  return authStore.usersList.filter(u => 
    (u.nip && u.nip.toLowerCase().includes(q)) ||
    (u.email && u.email.toLowerCase().includes(q)) ||
    (u.name && u.name.toLowerCase().includes(q)) ||
    (u.jabatan && u.jabatan.toLowerCase().includes(q)) ||
    (u.unit_kerja && u.unit_kerja.toLowerCase().includes(q))
  );
});

function openAddModal() {
  isEditing.value = false;
  editingId.value = null;
  form.value = {
    nip: '',
    email: '',
    name: '',
    jabatan: 'OPERATOR LAYANAN OPERASIONAL',
    unit_kerja: 'Dinas Sosial',
    role: 'intern',
    password: ''
  };
  showModal.value = true;
}

function openEditModal(u) {
  isEditing.value = true;
  editingId.value = u.id;
  form.value = {
    nip: u.nip || '',
    email: u.email || '',
    name: u.name || '',
    jabatan: u.jabatan || '',
    unit_kerja: u.unit_kerja || '',
    role: u.role || 'intern',
    password: ''
  };
  showModal.value = true;
}

async function saveUser() {
  try {
    if (isEditing.value) {
      await authStore.updateUser(editingId.value, form.value);
    } else {
      await authStore.createUser(form.value);
    }
    showModal.value = false;
  } catch (err) {
    alert(authStore.error || "Gagal menyimpan pengguna.");
  }
}

async function toggleActive(u) {
  try {
    await authStore.toggleUserActive(u.id, !u.is_active);
  } catch (err) {
    alert("Gagal merubah status keaktifan.");
  }
}

async function reset2fa(u) {
  if (!confirm(`Reset 2FA Google Authenticator untuk ${u.name}?`)) return;
  try {
    await authStore.resetUser2fa(u.id);
    alert(`2FA untuk ${u.name} berhasil direset.`);
  } catch (err) {
    alert("Gagal mereset 2FA.");
  }
}

onMounted(() => {
  authStore.fetchUsers();
});
</script>
