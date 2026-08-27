<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">

      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">admin_panel_settings</span>
            Manajemen Akun Administrator
          </h1>
          <p class="text-sm text-[#574146] mt-1 font-medium">LOPI-Q (Logbook, Online Presence, and Internship Quality Management System)</p>
        </div>
        <button
          v-if="authStore.isSuperAdmin"
          @click="openAddDialog"
          class="w-full sm:w-auto px-5 py-2.5 bg-[#F06292] hover:bg-[#ab2c5d] text-white font-bold text-xs rounded-lg transition-all flex items-center justify-center gap-2 cursor-pointer border-0 shadow-xs"
        >
          <span class="material-symbols-outlined text-base">person_add</span>
          <span>Tambah Pembimbing / Admin</span>
        </button>
      </div>

      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" :class="['flex items-center gap-2.5 p-3.5 rounded-xl text-xs font-semibold border w-full shadow-xs',
          toast.success ? 'bg-[#E8F5E9] border-[#A5D6A7] text-[#1B5E20]' : 'bg-[#FCE4EC] border-[#F8BBD0] text-[#F06292]'
        ]">
          <span class="material-symbols-outlined text-lg shrink-0">
            {{ toast.success ? 'check_circle' : 'error' }}
          </span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Main Content Card -->
      <div class="bg-white/85 backdrop-blur-md border border-[#F8BBD0] rounded-xl p-6 shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col gap-6">

        <!-- Loading State -->
        <div v-if="loading" class="flex flex-col items-center justify-center py-14 gap-3 text-[#8a7176]">
          <span class="material-symbols-outlined text-[40px] animate-spin text-[#f06292]">sync</span>
          <span class="text-xs font-semibold">Mengambil data administrator...</span>
        </div>

        <div v-else class="w-full">
          <!-- Desktop Table View -->
          <div class="hidden md:block overflow-x-auto w-full border border-[#F8BBD0] rounded-xl">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0] text-[11px] font-bold text-[#574146] uppercase tracking-wider">
                  <th class="py-3.5 px-6 min-w-[260px] whitespace-nowrap">Nama &amp; NIP</th>
                  <th class="py-3.5 px-6">Email / Jabatan</th>
                  <th class="py-3.5 px-6 whitespace-nowrap">Role</th>
                  <th class="py-3.5 px-6">Izin Akses</th>
                  <th class="py-3.5 px-6">Status 2FA</th>
                  <th class="py-3.5 px-6">Status Akun</th>
                  <th class="py-3.5 px-6 text-center w-32">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#F8BBD0] text-xs text-[#1b1c1c] bg-white">
                <tr v-for="admin in admins" :key="admin.id" class="hover:bg-[#FCE4EC]/30 transition-colors">
                  <!-- Avatar & Identity -->
                  <td class="py-4 px-6 min-w-[260px] whitespace-nowrap">
                    <div class="flex items-center gap-3">
                      <div class="h-9 w-9 rounded-lg shrink-0 font-bold text-xs flex items-center justify-center border shadow-xs"
                        :class="admin.role === 'superadmin' ? 'bg-[#ffd9e4] text-[#ab2c5d] border-[#F8BBD0]' : 'bg-[#FCE4EC] text-[#f06292] border-[#F8BBD0]'"
                      >
                        {{ admin.name?.charAt(0).toUpperCase() }}
                      </div>
                      <div class="min-w-0">
                        <div class="font-bold text-[#1b1c1c] leading-tight whitespace-nowrap">{{ admin.name }}</div>
                        <div class="font-mono text-[10px] text-[#574146] mt-0.5 whitespace-nowrap">NIP. {{ admin.nip || '-' }}</div>
                      </div>
                    </div>
                  </td>

                  <!-- Email & Jabatan -->
                  <td class="py-4 px-6">
                    <div class="font-mono text-[11px] text-[#1b1c1c] font-bold">{{ admin.email }}</div>
                    <div class="text-[10px] text-[#574146] mt-0.5">{{ admin.jabatan || 'Pembimbing Magang' }}</div>
                  </td>

                  <!-- Role Badge -->
                  <td class="py-4 px-6 whitespace-nowrap">
                    <span :class="['inline-flex items-center px-3 py-1 rounded-full text-[10px] font-bold tracking-wider uppercase border whitespace-nowrap shadow-2xs',
                      admin.role === 'superadmin' ? 'bg-[#ffd9e4] border-[#F8BBD0] text-[#ab2c5d]' : 'bg-[#FCE4EC] border-[#F8BBD0] text-[#f06292]']"
                    >
                      {{ admin.role === 'superadmin' ? 'SUPER ADMIN' : 'PEMBIMBING' }}
                    </span>
                  </td>

                  <!-- Permissions -->
                  <td class="py-4 px-6">
                    <div v-if="admin.role === 'superadmin'" class="text-[10px] text-[#1B5E20] font-bold flex items-center gap-1">
                      <span class="material-symbols-outlined text-[14px]">verified_user</span>
                      Full Access
                    </div>
                    <div v-else class="flex flex-wrap gap-1 max-w-[220px]">
                      <span v-for="p in (admin.permissions || [])" :key="p"
                        class="inline-flex items-center px-2 py-0.5 rounded bg-[#FCE4EC] border border-[#F8BBD0] text-[#ab2c5d] text-[9px] font-bold"
                      >{{ formatPermLabel(p) }}</span>
                      <span v-if="!admin.permissions || admin.permissions.length === 0" class="text-[#8a7176] italic text-[10px]">
                        Akses standar pembimbing
                      </span>
                    </div>
                  </td>

                  <!-- Status 2FA Badge -->
                  <td class="py-4 px-6 whitespace-nowrap">
                    <span :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[10px] font-bold border shadow-2xs',
                      admin.totp_enabled 
                        ? 'bg-[#E8F5E9] border-[#A5D6A7] text-[#1B5E20]' 
                        : 'bg-[#FFF8E1] border-[#FFE082] text-[#F57F17]']">
                      <span class="material-symbols-outlined text-[14px]">
                        {{ admin.totp_enabled ? 'verified_user' : 'shield_lock' }}
                      </span>
                      <span>{{ admin.totp_enabled ? '2FA Terverifikasi' : 'Belum Setup 2FA' }}</span>
                    </span>
                  </td>

                  <!-- Active Toggle -->
                  <td class="py-4 px-6">
                    <div class="flex items-center gap-2">
                      <button
                        type="button"
                        @click="handleToggleActive(admin)"
                        :disabled="!authStore.isSuperAdmin || admin.id === authStore.user?.id || togglingId === admin.id"
                        :class="[
                          admin.is_active ? 'bg-[#f06292]' : 'bg-slate-200',
                          'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed'
                        ]"
                      >
                        <span :class="[admin.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                      </button>
                      <span :class="['text-[10px] font-bold', admin.is_active ? 'text-[#ab2c5d]' : 'text-slate-400']">
                        {{ admin.is_active ? 'Aktif' : 'Nonaktif' }}
                      </span>
                    </div>
                  </td>

                  <!-- Action Buttons -->
                  <td class="py-4 px-6 text-center">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        v-if="authStore.isSuperAdmin"
                        @click="openEditDialog(admin)"
                        class="p-1.5 text-[#f06292] hover:text-[#ab2c5d] hover:bg-[#FCE4EC] rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Edit Admin / Pembimbing"
                      >
                        <span class="material-symbols-outlined text-base">edit</span>
                      </button>
                      <button
                        v-if="authStore.isSuperAdmin && admin.totp_enabled"
                        @click="handleReset2FA(admin)"
                        class="p-1.5 text-rose-600 hover:text-rose-700 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer border-0 bg-transparent"
                        title="Reset 2FA"
                      >
                        <span class="material-symbols-outlined text-base">lock_reset</span>
                      </button>
                      <span v-if="!authStore.isSuperAdmin" class="text-slate-300 text-[10px] italic">—</span>
                    </div>
                  </td>
                </tr>
                <tr v-if="admins.length === 0">
                  <td colspan="7" class="py-12 text-center text-[#8a7176] font-semibold text-xs">Belum ada administrator atau pembimbing terdaftar.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Mobile Card List View -->
          <div class="md:hidden flex flex-col gap-3 w-full">
            <div v-for="admin in admins" :key="admin.id" class="p-4 bg-white border border-[#F8BBD0] rounded-xl flex flex-col gap-3 shadow-xs">
              <!-- Card Header -->
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-2.5">
                  <div class="h-9 w-9 rounded-lg font-bold text-xs flex items-center justify-center shrink-0 border"
                    :class="admin.role === 'superadmin' ? 'bg-[#ffd9e4] text-[#ab2c5d] border-[#F8BBD0]' : 'bg-[#FCE4EC] text-[#f06292] border-[#F8BBD0]'"
                  >
                    {{ admin.name?.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <div class="font-bold text-[#1b1c1c] text-xs leading-tight">{{ admin.name }}</div>
                    <div class="font-mono text-[10px] text-[#574146]">{{ admin.email }}</div>
                  </div>
                </div>
                <span :class="['inline-flex items-center px-2.5 py-0.5 rounded-full text-[9px] font-bold tracking-wider uppercase border whitespace-nowrap shrink-0',
                  admin.role === 'superadmin' ? 'bg-[#ffd9e4] border-[#F8BBD0] text-[#ab2c5d]' : 'bg-[#FCE4EC] border-[#F8BBD0] text-[#f06292]']"
                >{{ admin.role === 'superadmin' ? 'SUPER ADMIN' : 'PEMBIMBING' }}</span>
              </div>

              <!-- Card Body -->
              <div class="flex flex-col gap-1 text-[11px] border-t border-[#F8BBD0] pt-2">
                <div class="flex justify-between">
                  <span class="text-[10px] text-[#574146] font-semibold">NIP:</span>
                  <span class="font-mono text-[#1b1c1c] font-bold">{{ admin.nip || '-' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-[10px] text-[#574146] font-semibold">Jabatan:</span>
                  <span class="text-[#1b1c1c] text-right max-w-[160px] truncate font-medium">{{ admin.jabatan || 'Pembimbing Magang' }}</span>
                </div>
              </div>

              <!-- Card Footer -->
              <div class="flex items-center justify-between border-t border-[#F8BBD0] pt-3">
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    @click="handleToggleActive(admin)"
                    :disabled="!authStore.isSuperAdmin || admin.id === authStore.user?.id || togglingId === admin.id"
                    :class="[admin.is_active ? 'bg-[#f06292]' : 'bg-slate-200', 'relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 focus:outline-none disabled:opacity-40 disabled:cursor-not-allowed']"
                  >
                    <span :class="[admin.is_active ? 'translate-x-4' : 'translate-x-0', 'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200']"></span>
                  </button>
                  <span :class="['text-[10px] font-bold', admin.is_active ? 'text-[#ab2c5d]' : 'text-slate-400']">
                    {{ admin.is_active ? 'Aktif' : 'Nonaktif' }}
                  </span>
                </div>
                <div class="flex items-center gap-1.5" v-if="authStore.isSuperAdmin">
                  <button @click="openEditDialog(admin)"
                    class="p-2 text-[#f06292] hover:bg-[#FCE4EC] rounded-lg transition-all cursor-pointer border border-[#F8BBD0] bg-white"
                    title="Edit"
                  ><span class="material-symbols-outlined text-base">edit</span></button>
                  <button
                    v-if="authStore.isSuperAdmin && admin.totp_enabled"
                    @click="handleReset2FA(admin)"
                    class="p-2 text-rose-600 hover:bg-rose-50 rounded-lg transition-all cursor-pointer border border-rose-200 bg-white flex items-center justify-center"
                    title="Reset 2FA"
                  ><span class="material-symbols-outlined text-base">lock_reset</span></button>
                </div>
              </div>
            </div>
            <div v-if="admins.length === 0" class="py-12 text-center text-[#8a7176] font-semibold text-xs">Belum ada administrator atau pembimbing.</div>
          </div>
        </div>
      </div>
    </div>

    <!-- ===== ADD / EDIT MODAL ===== -->
    <transition name="fade">
      <div v-if="dialogOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click.self="closeDialog">
        <div class="w-full max-w-lg bg-white rounded-xl border border-[#F8BBD0] shadow-2xl p-6 flex flex-col gap-5 max-h-[90vh] overflow-y-auto">

          <!-- Modal Header -->
          <div class="flex justify-between items-center border-b border-[#F8BBD0] pb-3">
            <h3 class="font-bold text-[#1b1c1c] text-base flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">{{ isEdit ? 'edit_square' : 'person_add' }}</span>
              {{ isEdit ? 'Edit Akun Administrator / Pembimbing' : 'Tambah Administrator / Pembimbing Baru' }}
            </h3>
            <button @click="closeDialog" class="p-1 rounded-full hover:bg-[#FCE4EC] text-[#574146] border-0 bg-transparent cursor-pointer">
              <span class="material-symbols-outlined text-lg">close</span>
            </button>
          </div>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3.5 rounded-lg bg-[#FCE4EC] border border-[#F8BBD0] text-[#F06292] text-xs flex items-start gap-2.5">
            <span class="material-symbols-outlined text-base shrink-0 mt-0.5">error</span>
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Form -->
          <form @submit.prevent="submitForm" class="flex flex-col gap-4">

            <!-- Row: Nama -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Nama Lengkap <span class="text-[#f06292]">*</span></label>
              <input type="text" v-model="form.name" required placeholder="Contoh: Pembimbing Magang A"
                class="w-full px-4 py-2.5 bg-white border border-[#F8BBD0] focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 rounded-lg focus:outline-none transition-all text-xs font-semibold text-[#1b1c1c]"
              />
            </div>

            <!-- Row: NIP -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">NIP / ID Pembimbing</label>
              <input type="text" v-model="form.nip" placeholder="Contoh: 199801012022011001"
                class="w-full px-4 py-2.5 bg-white border border-[#F8BBD0] focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 rounded-lg focus:outline-none transition-all text-xs font-semibold text-[#1b1c1c]"
              />
            </div>

            <!-- Row: Jabatan -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Jabatan / Unit Kerja</label>
              <input type="text" v-model="form.jabatan" placeholder="Contoh: Pembimbing Lapangan / Instansi"
                class="w-full px-4 py-2.5 bg-white border border-[#F8BBD0] focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 rounded-lg focus:outline-none transition-all text-xs font-semibold text-[#1b1c1c]"
              />
            </div>

            <!-- Row: Email -->
            <div class="flex flex-col gap-1">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Email <span class="text-[#f06292]">*</span></label>
              <input type="email" v-model="form.email" required placeholder="Contoh: pembimbing@bulukumbakab.go.id"
                class="w-full px-4 py-2.5 bg-white border border-[#F8BBD0] focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 rounded-lg focus:outline-none transition-all text-xs font-semibold text-[#1b1c1c]"
              />
            </div>

            <!-- Row: Password -->
            <div class="flex flex-col gap-1.5">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">
                {{ isEdit ? 'Reset Password Baru' : 'Password Akses' }}
                <span v-if="!isEdit" class="text-[#f06292]">*</span>
              </label>
              <div class="flex gap-2">
                <div class="relative flex-grow">
                  <input
                    :type="showPassword ? 'text' : 'password'"
                    v-model="form.password"
                    :required="!isEdit"
                    :placeholder="isEdit ? 'Biarkan kosong jika tidak ingin diubah' : 'Minimal 6 karakter'"
                    class="w-full pl-4 pr-10 py-2.5 bg-white border border-[#F8BBD0] focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 rounded-lg focus:outline-none transition-all text-xs font-semibold text-[#1b1c1c]"
                  />
                  <button type="button" @click="showPassword = !showPassword"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-[#574146] hover:text-[#ab2c5d] border-0 bg-transparent cursor-pointer"
                  >
                    <span class="material-symbols-outlined text-base">{{ showPassword ? 'visibility' : 'visibility_off' }}</span>
                  </button>
                </div>
                <button type="button" @click="generatePassword"
                  class="py-2.5 px-3 bg-[#FCE4EC] hover:bg-[#ffd9e4] border border-[#F8BBD0] rounded-lg text-xs font-bold transition-all cursor-pointer flex items-center gap-1 shrink-0 text-[#ab2c5d]"
                >
                  <span class="material-symbols-outlined text-base">vpn_key</span>
                  <span>Acak</span>
                </button>
              </div>
            </div>

            <!-- Row: Role Selector (Super Admin vs Pembimbing) -->
            <div class="flex flex-col gap-2">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Role / Peran Akses <span class="text-[#f06292]">*</span></label>
              <div class="grid grid-cols-2 gap-3">
                <label :class="['flex items-center gap-2 p-3 rounded-lg border transition-all select-none cursor-pointer text-xs',
                  form.role === 'superadmin' ? 'border-[#f06292] bg-[#ffd9e4] text-[#ab2c5d]' : 'border-[#F8BBD0] hover:bg-[#FCE4EC] text-[#574146]']"
                >
                  <input type="radio" v-model="form.role" value="superadmin" @change="syncDefaultPermissions" class="hidden" />
                  <span class="material-symbols-outlined text-xl">shield</span>
                  <div class="flex flex-col text-left">
                    <span class="font-bold">Super Admin</span>
                    <span class="text-[9px] opacity-80 font-normal">Akses Penuh Sistem</span>
                  </div>
                </label>

                <label :class="['flex items-center gap-2 p-3 rounded-lg border transition-all select-none cursor-pointer text-xs',
                  form.role === 'admin' ? 'border-[#f06292] bg-[#FCE4EC] text-[#ab2c5d]' : 'border-[#F8BBD0] hover:bg-[#FCE4EC] text-[#574146]']"
                >
                  <input type="radio" v-model="form.role" value="admin" @change="syncDefaultPermissions" class="hidden" />
                  <span class="material-symbols-outlined text-xl">supervisor_account</span>
                  <div class="flex flex-col text-left">
                    <span class="font-bold">Pembimbing</span>
                    <span class="text-[9px] opacity-80 font-normal">Pendampingan Intern</span>
                  </div>
                </label>
              </div>
            </div>

            <!-- Row: Permissions (only when role = admin/pembimbing) -->
            <div v-if="form.role !== 'superadmin'" class="flex flex-col gap-2">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Izin Akses Khusus Pembimbing</label>
              <div class="flex flex-col gap-0 p-3 bg-white rounded-lg border border-[#F8BBD0]">
                <label v-for="perm in availablePermissions" :key="perm.value"
                  class="flex items-center gap-2.5 text-xs font-medium py-2 cursor-pointer border-b border-[#F8BBD0]/50 last:border-0"
                >
                  <input type="checkbox" v-model="form.permissions" :value="perm.value"
                    class="w-4 h-4 text-[#f06292] border-[#F8BBD0] rounded accent-[#f06292]"
                  />
                  <div class="flex flex-col">
                    <span class="font-bold text-[#1b1c1c]">{{ perm.label }}</span>
                    <span class="text-[10px] text-[#574146]">{{ perm.desc }}</span>
                  </div>
                </label>
              </div>
            </div>
            <div v-else class="flex flex-col gap-2">
              <label class="text-[10px] font-bold text-[#574146] uppercase tracking-wider">Izin Akses Khusus</label>
              <div class="p-3 bg-[#E8F5E9] border border-[#A5D6A7] text-[#1B5E20] rounded-lg text-xs font-bold flex items-center gap-2">
                <span class="material-symbols-outlined text-base">verified_user</span>
                Super Admin memiliki akses penuh ke seluruh modul sistem LOPI-Q.
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-3 justify-end mt-2 pt-3 border-t border-[#F8BBD0]">
              <button type="button" @click="closeDialog"
                class="py-2 px-4 border border-[#F8BBD0] hover:bg-[#FCE4EC] transition-colors text-xs rounded-lg cursor-pointer bg-white font-bold text-[#574146]"
              >Batal</button>
              <button type="submit" :disabled="submitLoading"
                class="py-2 px-5 bg-[#ab2c5d] hover:bg-[#8b0e45] text-white text-xs rounded-lg transition-all flex items-center gap-1.5 cursor-pointer border-0 font-bold shadow-xs disabled:opacity-60"
              >
                <span v-if="submitLoading" class="animate-spin material-symbols-outlined text-base">sync</span>
                <span>{{ isEdit ? 'Simpan Perubahan' : 'Tambah Pembimbing' }}</span>
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
      <div class="bg-white rounded-xl max-w-[360px] w-full p-6 text-center border border-[#F8BBD0] shadow-2xl relative overflow-hidden flex flex-col items-center">
        <div class="relative flex items-center justify-center w-14 h-14 rounded-xl bg-[#FCE4EC] border border-[#F8BBD0] mb-4">
          <span class="material-symbols-outlined text-2xl text-[#f06292]">lock_reset</span>
        </div>
        <h3 class="font-bold text-base text-[#1b1c1c] mb-2 tracking-tight">{{ confirmModal.title }}</h3>
        <p class="text-xs text-[#574146] leading-relaxed px-2 mb-6">{{ confirmModal.description }}</p>
        <div class="flex gap-3 w-full">
          <button type="button" @click="confirmModal.show = false"
            class="flex-grow py-2.5 px-4 rounded-lg border border-[#F8BBD0] text-[#574146] hover:bg-[#FCE4EC] text-xs font-bold transition-all cursor-pointer bg-white flex items-center justify-center gap-1"
          >
            <span class="material-symbols-outlined text-sm">arrow_back</span>Batal
          </button>
          <button type="button" @click="confirmModal.onConfirm"
            class="flex-grow py-2.5 px-4 rounded-lg bg-[#ab2c5d] hover:bg-[#8b0e45] text-white text-xs font-bold transition-all cursor-pointer border-0 shadow-xs flex items-center justify-center gap-1"
          >
            <span class="material-symbols-outlined text-sm">lock_reset</span>Reset 2FA
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

const toast = ref({ show: false, success: true, message: '' })

const confirmModal = ref({
  show: false,
  title: '',
  description: '',
  onConfirm: () => {}
})

const availablePermissions = [
  { value: 'view_dashboard', label: 'Lihat Dashboard Overview', desc: 'Akses statistik presensi & keaktifan intern' },
  { value: 'manage_interns', label: 'Manajemen Peserta Magang', desc: 'Akses tambah, edit & nonaktifkan data intern' },
  { value: 'view_logs', label: 'Lihat Activity Log', desc: 'Akses logbook harian & audit trail presensi' },
  { value: 'view_reports', label: 'Lihat & Download Rekap Kehadiran', desc: 'Akses ekspor laporan presensi bulanan' }
]

const form = ref({
  name: '',
  nip: '',
  jabatan: 'Pembimbing Lapangan',
  email: '',
  password: '',
  role: 'admin',
  permissions: ['view_dashboard', 'manage_interns', 'view_logs', 'view_reports']
})

const showToastMsg = (success: boolean, message: string) => {
  toast.value = { show: true, success, message }
  setTimeout(() => { toast.value.show = false }, 4000)
}

const formatPermLabel = (val: string) => {
  const found = availablePermissions.find(p => p.value === val)
  return found ? found.label : val
}

const syncDefaultPermissions = () => {
  if (form.value.role === 'superadmin') {
    form.value.permissions = availablePermissions.map(p => p.value)
  }
}

const openAddDialog = () => {
  isEdit.value = false
  editUserId.value = null
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    name: '',
    nip: '',
    jabatan: 'Pembimbing Lapangan',
    email: '',
    password: '',
    role: 'admin',
    permissions: ['view_dashboard', 'manage_interns', 'view_logs', 'view_reports']
  }
  dialogOpen.value = true
}

const openEditDialog = (admin: any) => {
  isEdit.value = true
  editUserId.value = admin.id
  errorMessage.value = ''
  showPassword.value = false
  form.value = {
    name: admin.name || '',
    nip: admin.nip || '',
    jabatan: admin.jabatan || 'Pembimbing Lapangan',
    email: admin.email || '',
    password: '',
    role: admin.role || 'admin',
    permissions: Array.isArray(admin.permissions) ? [...admin.permissions] : []
  }
  dialogOpen.value = true
}

const closeDialog = () => {
  dialogOpen.value = false
}

const generatePassword = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_'
  let pass = ''
  for (let i = 0; i < 10; i++) {
    pass += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.value.password = pass
  showPassword.value = true
}

const fetchAdmins = async () => {
  loading.value = true
  try {
    const res = await authStore.fetchAdminsList()
    if (res.success && Array.isArray(res.admins)) {
      admins.value = res.admins
    }
  } catch (err) {
    console.warn('Gagal memuat list admin.')
  } finally {
    loading.value = false
  }
}

const handleToggleActive = async (admin: any) => {
  togglingId.value = admin.id
  try {
    const res = await authStore.toggleAdminStatus(admin.id, !admin.is_active)
    if (res.success) {
      admin.is_active = !admin.is_active
      showToastMsg(true, `Status akun ${admin.name} berhasil diubah menjadi ${admin.is_active ? 'Aktif' : 'Nonaktif'}.`)
    } else {
      showToastMsg(false, res.error || 'Gagal mengubah status akun.')
    }
  } catch (err: any) {
    showToastMsg(false, 'Terjadi kesalahan sistem saat mengubah status.')
  } finally {
    togglingId.value = null
  }
}

const handleReset2FA = (admin: any) => {
  confirmModal.value = {
    show: true,
    title: 'Reset 2FA Administrator',
    description: `Apakah Anda yakin ingin mereset Google Authenticator (2FA) untuk ${admin.name}? Pengguna harus melakukan pemindaian QR 2FA kembali saat login.`,
    onConfirm: async () => {
      confirmModal.value.show = false
      try {
        const res = await authStore.resetAdmin2FA(admin.id)
        if (res.success) {
          admin.totp_enabled = false
          showToastMsg(true, `2FA untuk ${admin.name} berhasil direset.`)
        } else {
          showToastMsg(false, res.error || 'Gagal mereset 2FA.')
        }
      } catch (err: any) {
        showToastMsg(false, 'Terjadi kesalahan saat mereset 2FA.')
      }
    }
  }
}

const submitForm = async () => {
  errorMessage.value = ''
  submitLoading.value = true
  try {
    if (isEdit.value && editUserId.value) {
      const payload: any = {
        name: form.value.name,
        nip: form.value.nip,
        jabatan: form.value.jabatan,
        email: form.value.email,
        role: form.value.role,
        permissions: form.value.permissions
      }
      if (form.value.password) {
        payload.password = form.value.password
      }
      const res = await authStore.updateAdminAccount(editUserId.value, payload)
      if (res.success) {
        showToastMsg(true, 'Data akun administrator/pembimbing berhasil diperbarui!')
        dialogOpen.value = false
        await fetchAdmins()
      } else {
        errorMessage.value = res.error || 'Gagal memperbarui akun.'
      }
    } else {
      const res = await authStore.createAdminAccount(form.value)
      if (res.success) {
        showToastMsg(true, 'Akun administrator/pembimbing baru berhasil ditambahkan!')
        dialogOpen.value = false
        await fetchAdmins()
      } else {
        errorMessage.value = res.error || 'Gagal menambahkan akun.'
      }
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Terjadi kesalahan koneksi sistem.'
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchAdmins()
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
