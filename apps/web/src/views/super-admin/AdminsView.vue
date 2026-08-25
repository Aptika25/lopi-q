<template>
  <AdminLayout>
    <div class="flex flex-col gap-6 w-full select-none">

      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Manajemen Akun Administrator</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Kelola akun Super Admin &amp; Admin Posko NTPD 112 Bulukumba beserta hak akses masing-masing.</p>
        </div>
        <button
          v-if="authStore.isSuperAdmin"
          @click="openAddDialog"
          class="w-full sm:w-auto py-2.5 px-4 bg-rose-700 hover:bg-rose-800 text-white font-bold text-xs rounded-xl shadow-sm transition-all flex items-center justify-center gap-1.5 cursor-pointer border-0"
        >
          <span class="material-symbols-outlined text-[16px]">person_add</span>
          <span>Tambah Administrator</span>
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

      <!-- Main Content -->
      <div class="bg-white border border-slate-200 rounded-3xl p-6 shadow-sm flex flex-col gap-6">

        <!-- Loading State -->
        <div v-if="loading" class="flex flex-col items-center justify-center py-14 gap-3 text-slate-400">
          <span class="material-symbols-outlined text-[40px] animate-spin">sync</span>
          <span class="text-xs">Mengambil data administrator...</span>
        </div>

        <div v-else class="w-full">
          <!-- Desktop Table View -->
          <div class="hidden md:block overflow-x-auto w-full border border-slate-200 rounded-2xl">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">
                  <th class="py-3.5 px-4 min-w-[260px] whitespace-nowrap">Nama &amp; NIP</th>
                  <th class="py-3.5 px-4">Email / Jabatan</th>
                  <th class="py-3.5 px-4 whitespace-nowrap">Role</th>
                  <th class="py-3.5 px-4">Izin Akses</th>
                  <th class="py-3.5 px-4">Status 2FA</th>
                  <th class="py-3.5 px-4">Status Akun</th>
                  <th class="py-3.5 px-4 text-center w-32">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 text-xs text-slate-700">
                <tr v-for="admin in admins" :key="admin.id" class="hover:bg-slate-50/60 transition-colors">
                  <!-- Avatar & Identity -->
                  <td class="py-3.5 px-4 min-w-[260px] whitespace-nowrap">
                    <div class="flex items-center gap-3">
                      <div class="h-9 w-9 rounded-xl shrink-0 font-black text-xs flex items-center justify-center shadow-sm"
                        :class="admin.role === 'superadmin' ? 'bg-rose-100 text-rose-800 border border-rose-200' : 'bg-indigo-100 text-indigo-800 border border-indigo-200'"
                      >
                        {{ admin.name?.charAt(0).toUpperCase() }}
                      </div>
                      <div class="min-w-0">
                        <div class="font-bold text-slate-900 leading-tight whitespace-nowrap">{{ admin.name }}</div>
                        <div class="font-mono text-[10px] text-slate-400 mt-0.5 whitespace-nowrap">NIP. {{ admin.nip || '-' }}</div>
                      </div>
                    </div>
                  </td>

                  <!-- Email & Jabatan -->
                  <td class="py-3.5 px-4">
                    <div class="font-mono text-[11px] text-slate-600">{{ admin.email }}</div>
                    <div class="text-[10px] text-slate-400 mt-0.5">{{ admin.jabatan || 'Diskominfo Kab. Bulukumba' }}</div>
                  </td>

                  <!-- Role Badge -->
                  <td class="py-3.5 px-4 whitespace-nowrap">
                    <span :class="['inline-flex items-center px-3 py-1 rounded-full text-[10px] font-extrabold tracking-wider uppercase border whitespace-nowrap shadow-2xs',
                      admin.role === 'superadmin' ? 'bg-rose-50 border-rose-200 text-rose-700' : 'bg-indigo-50 border-indigo-200 text-indigo-700']"
                    >
                      {{ admin.role === 'superadmin' ? 'SUPER ADMIN' : 'ADMIN' }}
                    </span>
                  </td>

                  <!-- Permissions -->
                  <td class="py-3.5 px-4">
                    <div v-if="admin.role === 'superadmin'" class="text-[10px] text-emerald-700 font-bold flex items-center gap-1">
                      <span class="material-symbols-outlined text-[14px]">verified_user</span>
                      Full Access
                    </div>
                    <div v-else class="flex flex-wrap gap-1 max-w-[220px]">
                      <span v-for="p in (admin.permissions || [])" :key="p"
                        class="inline-flex items-center px-1.5 py-0.5 rounded bg-slate-100 text-slate-700 text-[9px] font-medium"
                      >{{ formatPermLabel(p) }}</span>
                      <span v-if="!admin.permissions || admin.permissions.length === 0" class="text-slate-400 italic text-[10px]">
                        Tidak ada izin
                      </span>
                    </div>
                  </td>

                  <!-- Status 2FA Badge -->
                  <td class="py-3.5 px-4 whitespace-nowrap">
                    <span :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-extrabold border shadow-2xs',
                      admin.totp_enabled 
                        ? 'bg-emerald-50 border-emerald-200/80 text-emerald-700' 
                        : 'bg-amber-50 border-amber-200/80 text-amber-700']">
                      <span class="material-symbols-outlined text-[14px]" :class="admin.totp_enabled ? 'text-emerald-600' : 'text-amber-600'">
                        {{ admin.totp_enabled ? 'verified_user' : 'shield_lock' }}
                      </span>
                      <span>{{ admin.totp_enabled ? '2FA Terverifikasi' : 'Belum Setup 2FA' }}</span>
                    </span>
                  </td>

                  <!-- Active Toggle -->
                  <td class="py-3.5 px-4">
                    <div class="flex items-center gap-2">
                      <button
                        type="button"
                        @click="handleToggleActive(admin)"
                        :disabled="!authStore.isSuperAdmin || admin.id === authStore.user?.id || togglingId === admin.id"
                        :class="[
                          admin.is_active ? 'bg-emerald-500' : 'bg-slate-200',
                          'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed'
                        ]"
                      >
                        <span :class="[admin.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                      </button>
                      <span :class="['text-[10px] font-bold', admin.is_active ? 'text-emerald-600' : 'text-slate-400']">
                        {{ admin.is_active ? 'Aktif' : 'Nonaktif' }}
                      </span>
                    </div>
                  </td>

                  <!-- Action Buttons -->
                  <td class="py-3.5 px-4 text-center">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        v-if="authStore.isSuperAdmin"
                        @click="openEditDialog(admin)"
                        class="p-1.5 text-amber-600 hover:text-amber-700 hover:bg-amber-50 rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Edit Administrator"
                      >
                        <span class="material-symbols-outlined text-[16px]">edit</span>
                      </button>
                      <button
                        v-if="authStore.isSuperAdmin && admin.totp_enabled"
                        @click="handleReset2FA(admin)"
                        class="p-1.5 text-rose-600 hover:text-rose-700 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Reset 2FA Administrator"
                      >
                        <span class="material-symbols-outlined text-[16px]">lock_reset</span>
                      </button>
                      <span v-if="!authStore.isSuperAdmin" class="text-slate-300 text-[10px] italic">—</span>
                    </div>
                  </td>
                </tr>
                <tr v-if="admins.length === 0">
                  <td colspan="6" class="py-12 text-center text-slate-400 font-medium">Belum ada administrator terdaftar.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Mobile Card List View -->
          <div class="md:hidden flex flex-col gap-3 w-full">
            <div v-for="admin in admins" :key="admin.id" class="p-4 bg-slate-50 border border-slate-200/60 rounded-2xl flex flex-col gap-3">
              <!-- Card Header -->
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-2.5">
                  <div class="h-9 w-9 rounded-xl font-black text-xs flex items-center justify-center shrink-0"
                    :class="admin.role === 'superadmin' ? 'bg-rose-100 text-rose-800 border border-rose-200' : 'bg-indigo-100 text-indigo-800 border border-indigo-200'"
                  >
                    {{ admin.name?.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <div class="font-bold text-slate-900 text-xs leading-tight">{{ admin.name }}</div>
                    <div class="font-mono text-[10px] text-slate-400">{{ admin.email }}</div>
                  </div>
                </div>
                <span :class="['inline-flex items-center px-2.5 py-0.5 rounded-full text-[9px] font-extrabold tracking-wider uppercase border whitespace-nowrap shrink-0',
                  admin.role === 'superadmin' ? 'bg-rose-50 border-rose-200 text-rose-700' : 'bg-indigo-50 border-indigo-200 text-indigo-700']"
                >{{ admin.role === 'superadmin' ? 'SUPER ADMIN' : 'ADMIN' }}</span>
              </div>

              <!-- Card Body -->
              <div class="flex flex-col gap-1 text-[11px] border-t border-slate-200/60 pt-2">
                <div class="flex justify-between">
                  <span class="text-[10px] text-slate-400 font-semibold">NIP.</span>
                  <span class="font-mono text-slate-600">{{ admin.nip || '-' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-[10px] text-slate-400 font-semibold">Jabatan:</span>
                  <span class="text-slate-600 text-right max-w-[160px] truncate">{{ admin.jabatan || 'Diskominfo' }}</span>
                </div>
                <div class="flex flex-col gap-1 mt-1">
                  <span class="text-[10px] text-slate-400 font-semibold">Izin Akses:</span>
                  <div v-if="admin.role === 'superadmin'" class="text-[10px] text-emerald-700 font-bold flex items-center gap-1">
                    <span class="material-symbols-outlined text-[13px]">verified_user</span> Full Access
                  </div>
                  <div v-else class="flex flex-wrap gap-1">
                    <span v-for="p in (admin.permissions || [])" :key="p"
                      class="px-1.5 py-0.5 rounded bg-white border border-slate-200 text-slate-700 text-[9px] font-medium"
                    >{{ formatPermLabel(p) }}</span>
                    <span v-if="!admin.permissions || admin.permissions.length === 0" class="text-slate-400 italic text-[10px]">Tidak ada</span>
                  </div>
                </div>
              </div>

              <!-- Card Footer -->
              <div class="flex items-center justify-between border-t border-slate-200/60 pt-3">
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    @click="handleToggleActive(admin)"
                    :disabled="!authStore.isSuperAdmin || admin.id === authStore.user?.id || togglingId === admin.id"
                    :class="[admin.is_active ? 'bg-emerald-500' : 'bg-slate-200', 'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed']"
                  >
                    <span :class="[admin.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                  </button>
                  <span :class="['text-[10px] font-bold', admin.is_active ? 'text-emerald-600' : 'text-slate-400']">
                    {{ admin.is_active ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </div>
                <div class="flex items-center gap-1.5" v-if="authStore.isSuperAdmin">
                  <button @click="openEditDialog(admin)"
                    class="p-2 text-amber-600 hover:bg-amber-50 rounded-xl transition-all cursor-pointer border border-amber-200/40 bg-white"
                    title="Edit"
                  ><span class="material-symbols-outlined text-[16px]">edit</span></button>
                  <button
                    v-if="authStore.isSuperAdmin && admin.totp_enabled"
                    @click="handleReset2FA(admin)"
                    class="p-2 text-rose-600 hover:bg-rose-50 rounded-xl transition-all cursor-pointer border border-rose-200/40 bg-white flex items-center justify-center"
                    title="Reset 2FA"
                  ><span class="material-symbols-outlined text-[16px]">lock_reset</span></button>
                </div>
              </div>
            </div>
            <div v-if="admins.length === 0" class="py-12 text-center text-slate-400 font-medium text-xs">Belum ada administrator.</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===== ADD / EDIT MODAL ===== -->
    <transition name="fade">
      <div v-if="dialogOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click.self="closeDialog">
        <div class="w-full max-w-lg bg-white rounded-3xl border border-slate-200 shadow-2xl p-6 flex flex-col gap-5 max-h-[90vh] overflow-y-auto">

          <!-- Modal Header -->
          <div class="flex justify-between items-center border-b border-slate-100 pb-3">
            <h3 class="font-display font-black text-slate-900 text-base flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">{{ isEdit ? 'edit_square' : 'person_add' }}</span>
              {{ isEdit ? 'Edit Akun Administrator' : 'Tambah Administrator Baru' }}
            </h3>
            <button @click="closeDialog" class="p-1 rounded-full hover:bg-slate-100 text-slate-400 border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-[18px]">close</span>
            </button>
          </div>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-red-50 border border-red-200 text-red-800 text-xs flex items-start gap-2.5">
            <span class="material-symbols-outlined text-[16px] shrink-0 mt-0.5">error</span>
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Form -->
          <form @submit.prevent="submitForm" class="flex flex-col gap-4">

            <!-- Row: Nama -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Nama Lengkap <span class="text-rose-500">*</span></label>
              <input type="text" v-model="form.name" required placeholder="Contoh: Muhammad Aswan"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Row: NIP -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">NIP</label>
              <input type="text" v-model="form.nip" placeholder="Contoh: 199801012022011001"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Row: Jabatan -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Jabatan</label>
              <input type="text" v-model="form.jabatan" placeholder="Contoh: Pranata Komputer Ahli Pertama"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Row: Email -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Email Dinas <span class="text-rose-500">*</span></label>
              <input type="email" v-model="form.email" required placeholder="Contoh: aswan@bulukumbakab.go.id"
                class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 focus:border-rose-500 focus:ring-2 focus:ring-rose-500/10 rounded-xl focus:outline-none transition-all text-xs"
              />
            </div>

            <!-- Row: Password -->
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

            <!-- Row: Role Selector -->
            <div class="flex flex-col gap-2">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Role / Peran Sistem <span class="text-rose-500">*</span></label>
              <div class="grid grid-cols-2 gap-2.5">
                <label :class="['flex items-center gap-2 p-3 rounded-xl border transition-all select-none cursor-pointer text-xs',
                  form.role === 'superadmin' ? 'border-rose-500 bg-rose-50 text-rose-800' : 'border-slate-200 hover:bg-slate-50 text-slate-600']"
                >
                  <input type="radio" v-model="form.role" value="superadmin" @change="syncDefaultPermissions" class="hidden" />
                  <span class="material-symbols-outlined text-[20px]">shield</span>
                  <div class="flex flex-col text-left">
                    <span class="font-bold">Super Admin</span>
                    <span class="text-[9px] opacity-70 font-normal">Akses penuh</span>
                  </div>
                </label>
                <label :class="['flex items-center gap-2 p-3 rounded-xl border transition-all select-none cursor-pointer text-xs',
                  form.role === 'admin' ? 'border-indigo-500 bg-indigo-50 text-indigo-800' : 'border-slate-200 hover:bg-slate-50 text-slate-600']"
                >
                  <input type="radio" v-model="form.role" value="admin" @change="syncDefaultPermissions" class="hidden" />
                  <span class="material-symbols-outlined text-[20px]">manage_accounts</span>
                  <div class="flex flex-col text-left">
                    <span class="font-bold">Admin</span>
                    <span class="text-[9px] opacity-70 font-normal">Akses terbatas</span>
                  </div>
                </label>
              </div>
            </div>

            <!-- Row: Permissions (only when role = admin) -->
            <div v-if="form.role !== 'superadmin'" class="flex flex-col gap-2">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Izin Akses Khusus</label>
              <div class="flex flex-col gap-0 p-3 bg-slate-50 rounded-2xl border border-slate-200">
                <label v-for="perm in availablePermissions" :key="perm.value"
                  class="flex items-center gap-2.5 text-xs font-medium py-2 cursor-pointer border-b border-slate-100 last:border-0"
                >
                  <input type="checkbox" v-model="form.permissions" :value="perm.value"
                    class="w-4 h-4 text-rose-600 border-slate-300 rounded accent-rose-600"
                  />
                  <div class="flex flex-col">
                    <span class="font-semibold text-slate-700">{{ perm.label }}</span>
                    <span class="text-[10px] text-slate-400">{{ perm.desc }}</span>
                  </div>
                </label>
              </div>
            </div>
            <div v-else class="flex flex-col gap-2">
              <label class="text-[10px] font-extrabold text-slate-500 uppercase tracking-widest">Izin Akses Khusus</label>
              <div class="p-3 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl text-xs font-semibold flex items-center gap-2">
                <span class="material-symbols-outlined text-[18px]">verified_user</span>
                Super Admin memiliki akses penuh ke semua fitur sistem.
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
                <span>{{ isEdit ? 'Simpan Perubahan' : 'Tambah Administrator' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- ===== CUSTOM CONFIRM DIALOG (Reset 2FA) ===== -->
    <div v-if="confirmModal.show"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm"
      @click.self="confirmModal.show = false"
    >
      <div class="bg-white rounded-[28px] max-w-[360px] w-full p-6 text-center border border-slate-100 shadow-2xl relative overflow-hidden flex flex-col items-center">
        <div class="absolute -top-12 -right-12 w-28 h-28 bg-rose-100/40 rounded-full blur-2xl pointer-events-none"></div>
        <div class="absolute -bottom-12 -left-12 w-28 h-28 bg-slate-100 rounded-full blur-2xl pointer-events-none"></div>
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
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AdminLayout from '@/layouts/AdminLayout.vue'

const authStore = useAuthStore()

// ===== STATE =====
const admins = ref<any[]>([])
const loading = ref(false)
const togglingId = ref<number | null>(null)
const submitLoading = ref(false)
const dialogOpen = ref(false)
const isEdit = ref(false)
const editUserId = ref<number | null>(null)
const errorMessage = ref('')
const showPassword = ref(false)

const form = ref({
  email: '',
  password: '',
  name: '',
  nip: '',
  jabatan: '',
  unit_kerja: '',
  role: 'admin',
  permissions: [] as string[]
})

const toast = ref({ show: false, success: true, message: '' })

const confirmModal = ref({
  show: false,
  title: '',
  description: '',
  onConfirm: () => {}
})

// ===== PERMISSIONS DEFINITION =====
const availablePermissions = [
  { value: 'manage_users',      label: 'Manajemen User & Admin',   desc: 'Tambah, edit, nonaktifkan user dan admin' },
  { value: 'manage_attendance', label: 'Kelola Data Presensi',      desc: 'Lihat dan kelola rekapan kehadiran call taker' },
  { value: 'manage_locations',  label: 'Kelola Lokasi & QR Code',   desc: 'Tambah, edit, dan hapus posko & QR Code' },
  { value: 'view_reports',      label: 'Akses Laporan & Rekapan',   desc: 'Unduh dan cetak laporan kehadiran' },
]

const formatPermLabel = (p: string) => {
  const found = availablePermissions.find(x => x.value === p)
  return found ? found.label : p
}

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

// ===== LOAD ADMINS =====
const loadAdmins = async () => {
  loading.value = true
  try {
    await authStore.fetchUsers()
    // Filter only superadmin & admin, sort primary superadmin first, then role hierarchy & name
    const roleOrder: Record<string, number> = { superadmin: 1, admin: 2 }
    admins.value = (authStore.usersList || [])
      .filter((u: any) => u.role === 'superadmin' || u.role === 'admin')
      .sort((a: any, b: any) => {
        const isPrimaryA = a.id === 1 || a.email === 'aswan@bulukumbakab.go.id'
        const isPrimaryB = b.id === 1 || b.email === 'aswan@bulukumbakab.go.id'
        if (isPrimaryA && !isPrimaryB) return -1
        if (!isPrimaryA && isPrimaryB) return 1

        const orderA = roleOrder[a.role] ?? 99
        const orderB = roleOrder[b.role] ?? 99
        if (orderA !== orderB) return orderA - orderB
        return a.name.localeCompare(b.name)
      })
  } catch (err) {
    showToast(false, 'Gagal memuat data administrator.')
  } finally {
    loading.value = false
  }
}

onMounted(() => loadAdmins())

// ===== DIALOG OPEN/CLOSE =====
const openAddDialog = () => {
  isEdit.value = false
  editUserId.value = null
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    email: '', password: '', name: '', nip: '', jabatan: '',
    unit_kerja: '', role: 'admin',
    permissions: ['manage_attendance', 'manage_locations', 'view_reports']
  }
  dialogOpen.value = true
}

const openEditDialog = (admin: any) => {
  isEdit.value = true
  editUserId.value = admin.id
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    email: admin.email,
    password: '',
    name: admin.name,
    nip: admin.nip || '',
    jabatan: admin.jabatan || '',
    unit_kerja: admin.unit_kerja || '',
    role: admin.role,
    permissions: [...(admin.permissions || [])]
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
  for (let i = 0; i < 12; i++) pass += chars.charAt(Math.floor(Math.random() * chars.length))
  form.value.password = pass
  showPassword.value = true
}

// ===== SYNC DEFAULT PERMISSIONS =====
const syncDefaultPermissions = () => {
  if (form.value.role === 'superadmin') {
    form.value.permissions = ['manage_users', 'manage_attendance', 'manage_locations', 'view_reports']
  } else {
    form.value.permissions = ['manage_attendance', 'manage_locations', 'view_reports']
  }
}

// ===== TOGGLE ACTIVE =====
const handleToggleActive = async (admin: any) => {
  if (admin.id === authStore.user?.id) {
    showToast(false, 'Anda tidak dapat menonaktifkan akun Anda sendiri.')
    return
  }
  const newState = !admin.is_active
  togglingId.value = admin.id
  try {
    await authStore.toggleUserActive(admin.id, newState)
    admin.is_active = newState
    showToast(true, `Akun ${admin.name} berhasil ${newState ? 'diaktifkan' : 'dinonaktifkan'}.`)
  } catch (err) {
    showToast(false, 'Gagal mengubah status keaktifan akun.')
  } finally {
    togglingId.value = null
  }
}

// ===== RESET 2FA =====
const handleReset2FA = (admin: any) => {
  triggerConfirm(
    'Reset Keamanan 2FA?',
    `Apakah Anda yakin ingin menonaktifkan Google Authenticator (2FA) untuk administrator ${admin.name}? Administrator harus melakukan setup ulang 2FA saat login berikutnya.`,
    async () => {
      submitLoading.value = true
      try {
        await authStore.resetUser2fa(admin.id)
        showToast(true, `2FA untuk ${admin.name} berhasil dinonaktifkan.`)
        await loadAdmins()
      } catch (err) {
        showToast(false, 'Gagal mereset 2FA administrator.')
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
  if (!isEdit.value) {
    if (!form.value.email || !form.value.password) { errorMessage.value = 'Email dan Password wajib diisi.'; return }
    if (form.value.password.length < 6) { errorMessage.value = 'Password minimal 6 karakter.'; return }
  }

  const permissions = form.value.role === 'superadmin'
    ? ['manage_users', 'manage_attendance', 'manage_locations', 'view_reports']
    : form.value.permissions

  submitLoading.value = true
  try {
    if (isEdit.value && editUserId.value !== null) {
      const res = await authStore.updateUser(editUserId.value, {
        name: form.value.name,
        nip: form.value.nip,
        jabatan: form.value.jabatan,
        unit_kerja: form.value.unit_kerja,
        role: form.value.role,
        permissions,
        password: form.value.password
      })
      if (res.success) {
        showToast(true, 'Data administrator berhasil diperbarui.')
        closeDialog()
        await loadAdmins()
      } else {
        errorMessage.value = res.error || 'Gagal memperbarui administrator.'
      }
    } else {
      const res = await authStore.createUser({
        email: form.value.email,
        password: form.value.password,
        name: form.value.name,
        nip: form.value.nip,
        jabatan: form.value.jabatan,
        unit_kerja: form.value.unit_kerja,
        role: form.value.role,
        permissions
      })
      if (res.success) {
        showToast(true, 'Administrator baru berhasil ditambahkan!')
        closeDialog()
        await loadAdmins()
      } else {
        errorMessage.value = res.error || 'Gagal menambahkan administrator.'
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
