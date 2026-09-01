<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">
      
      <!-- Header -->
      <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">security</span>
            Security Settings
          </h1>
          <p class="text-sm text-[#574146] mt-1">Manage your account security, passwords, and active sessions to keep your profile protected.</p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <span class="px-3.5 py-1.5 bg-[#ffd9e4] text-[#ab2c5d] border border-[#ddbfc5] rounded-full text-xs font-bold shrink-0 flex items-center gap-1.5 whitespace-nowrap">
            <span class="w-2 h-2 rounded-full bg-[#ab2c5d] animate-pulse"></span>
            <span>2FA TERPROTEKSI</span>
          </span>
        </div>
      </div>

      <!-- Main Layout Grid (Left: Password & Sessions, Right: 2FA & Permissions) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        
        <!-- Left Column: Password & Active Sessions (lg:col-span-8) -->
        <div class="lg:col-span-8 flex flex-col gap-6">
          
          <!-- Password Management Card -->
          <section class="bg-white/90 backdrop-blur-md rounded-xl border border-[#ddbfc5] p-6 shadow-xs space-y-4">
            <div class="flex items-center gap-3 border-b border-[#ddbfc5]/50 pb-3">
              <div class="w-10 h-10 rounded-full bg-[#ffd9e4] flex items-center justify-center text-[#ab2c5d]">
                <span class="material-symbols-outlined text-xl">key</span>
              </div>
              <div>
                <h3 class="font-bold text-base text-[#1b1c1c]">Change Password</h3>
                <p class="text-xs text-[#574146]">Perbarui kata sandi akun admin Anda secara berkala</p>
              </div>
            </div>

            <form @submit.prevent="handleUpdatePassword" class="space-y-4 pt-2">
              <div class="space-y-1.5">
                <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider">Current Password</label>
                <input 
                  v-model="passwordForm.currentPassword"
                  type="password"
                  placeholder="Enter current password"
                  required
                  class="w-full bg-white border border-[#ddbfc5] rounded-lg px-4 py-2.5 text-xs text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/50 transition-all placeholder-[#8a7176] shadow-xs"
                />
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="space-y-1.5">
                  <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider">New Password</label>
                  <input 
                    v-model="passwordForm.newPassword"
                    type="password"
                    placeholder="New password (min 8 char)"
                    required
                    class="w-full bg-white border border-[#ddbfc5] rounded-lg px-4 py-2.5 text-xs text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/50 transition-all placeholder-[#8a7176] shadow-xs"
                  />
                </div>

                <div class="space-y-1.5">
                  <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider">Confirm New Password</label>
                  <input 
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    placeholder="Confirm new password"
                    required
                    class="w-full bg-white border border-[#ddbfc5] rounded-lg px-4 py-2.5 text-xs text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/50 transition-all placeholder-[#8a7176] shadow-xs"
                  />
                </div>
              </div>

              <div class="pt-2 flex justify-end">
                <button 
                  type="submit"
                  :disabled="updatingPassword"
                  class="bg-[#f06292] text-white px-5 py-2.5 rounded-lg font-bold text-xs shadow-xs hover:bg-[#ab2c5d] transition-colors active:scale-95 duration-150 border-0 cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
                >
                  <span class="material-symbols-outlined text-base">update</span>
                  <span>{{ updatingPassword ? 'Saving...' : 'Update Password' }}</span>
                </button>
              </div>
            </form>
          </section>

          <!-- Active Sessions Card -->
          <section class="bg-white/90 backdrop-blur-md rounded-xl border border-[#ddbfc5] overflow-hidden shadow-xs">
            <div class="p-6 border-b border-[#ddbfc5]/60 bg-gradient-to-b from-white to-[#ffd9e4]/20 flex justify-between items-center">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-[#ffd9e4] flex items-center justify-center text-[#ab2c5d]">
                  <span class="material-symbols-outlined text-xl">devices</span>
                </div>
                <div>
                  <h3 class="font-bold text-base text-[#1b1c1c]">Active Sessions</h3>
                  <p class="text-xs text-[#574146]">Perangkat aktif yang saat ini terhubung dengan akun Anda</p>
                </div>
              </div>
              <button 
                @click="logoutAllSessions"
                class="text-[#ab2c5d] font-bold text-xs uppercase hover:text-[#5e002b] transition-colors cursor-pointer border-0 bg-transparent"
              >
                Log out all
              </button>
            </div>

            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr class="bg-[#ffd9e4]/30 border-b border-[#ddbfc5]/50">
                    <th class="py-3 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Device</th>
                    <th class="py-3 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Location / IP</th>
                    <th class="py-3 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Last Active</th>
                    <th class="py-3 px-6 text-right"></th>
                  </tr>
                </thead>
                <tbody class="text-xs divide-y divide-[#ddbfc5]/30">
                  <tr v-for="session in activeSessions" :key="session.id" class="hover:bg-[#ffd9e4]/10 transition-colors">
                    <td class="py-4 px-6">
                      <div class="flex items-center gap-3">
                        <span class="material-symbols-outlined text-[#8a7176] text-xl">{{ session.icon }}</span>
                        <div>
                          <p class="font-bold text-[#1b1c1c]">{{ session.deviceName }}</p>
                          <p class="text-[#574146] text-[11px]">{{ session.browser }}</p>
                        </div>
                      </div>
                    </td>
                    <td class="py-4 px-6 text-[#574146]">
                      <div class="font-medium text-[#1b1c1c]">{{ session.location }}</div>
                      <div class="text-[10px] font-mono text-[#8a7176]">{{ session.ipAddress }}</div>
                    </td>
                    <td class="py-4 px-6">
                      <span 
                        v-if="session.isCurrent" 
                        class="inline-flex items-center px-2.5 py-0.5 rounded-full bg-[#E8F5E9] text-[#1B5E20] font-bold text-[10px] uppercase tracking-wide border border-[#A5D6A7]"
                      >
                        Active Now
                      </span>
                      <span v-else class="text-[#574146] text-[11px]">{{ session.lastActive }}</span>
                    </td>
                    <td class="py-4 px-6 text-right">
                      <button 
                        v-if="!session.isCurrent"
                        @click="revokeSession(session.id)"
                        class="text-rose-600 hover:text-rose-800 transition-colors p-1.5 rounded-full hover:bg-rose-50 cursor-pointer border-0 bg-transparent"
                        title="Revoke session"
                      >
                        <span class="material-symbols-outlined text-lg">logout</span>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

        </div>

        <!-- Right Column: 2FA & Permissions (lg:col-span-4) -->
        <div class="lg:col-span-4 flex flex-col gap-6">
          
          <!-- 2FA Setup Card -->
          <section class="bg-gradient-to-br from-white to-[#ffd9e4]/20 rounded-xl border border-[#ddbfc5] p-6 shadow-xs space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-[#ffd9e4] flex items-center justify-center text-[#ab2c5d]">
                  <span class="material-symbols-outlined text-xl">shield_lock</span>
                </div>
                <div>
                  <h3 class="font-bold text-base text-[#1b1c1c]">Two-Factor Auth</h3>
                  <p class="text-[11px] text-[#574146]">Google Authenticator TOTP</p>
                </div>
              </div>
              <!-- Custom Toggle Switch -->
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="twoFactorEnabled" @change="toggle2FA" class="sr-only peer" />
                <div class="w-11 h-6 bg-slate-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#f06292]"></div>
              </label>
            </div>

            <p class="text-xs text-[#574146] leading-relaxed">
              Tambahkan lapisan keamanan ekstra dengan kode verifikasi 6 digit setiap kali Anda melakukan login.
            </p>

            <div class="space-y-2.5 pt-1">
              <button 
                @click="openSetupModal"
                class="w-full bg-[#ffd9e4] text-[#ab2c5d] px-4 py-2.5 rounded-lg font-bold text-xs hover:bg-[#fec1d6] transition-colors flex items-center justify-center gap-2 border border-[#ddbfc5] cursor-pointer shadow-xs"
              >
                <span class="material-symbols-outlined text-base">qr_code_scanner</span>
                <span>Setup Authenticator App</span>
              </button>
              
              <button 
                @click="showBackupCodes"
                class="w-full bg-slate-100 text-[#1b1c1c] px-4 py-2.5 rounded-lg font-bold text-xs hover:bg-slate-200 transition-colors flex items-center justify-center gap-2 border border-[#ddbfc5] cursor-pointer shadow-xs"
              >
                <span class="material-symbols-outlined text-base">key</span>
                <span>View Emergency Backup Codes</span>
              </button>
            </div>
          </section>

          <!-- Access Level & Permissions Card -->
          <section class="bg-white/90 backdrop-blur-md rounded-xl border border-[#ddbfc5] p-6 shadow-xs space-y-4">
            <div class="flex items-center gap-3 border-b border-[#ddbfc5]/50 pb-3">
              <div class="w-10 h-10 rounded-full bg-[#ffd9e4] flex items-center justify-center text-[#ab2c5d]">
                <span class="material-symbols-outlined text-xl">admin_panel_settings</span>
              </div>
              <div>
                <h3 class="font-bold text-base text-[#1b1c1c]">Access Level</h3>
                <p class="text-xs text-[#574146]">Otoritas dan perizinan role</p>
              </div>
            </div>

            <div class="bg-[#ffd9e4]/30 rounded-lg p-4 border border-[#ddbfc5] space-y-1">
              <div class="flex justify-between items-center">
                <span class="font-bold text-sm text-[#1b1c1c]">Super Administrator</span>
                <span class="inline-flex items-center px-2 py-0.5 rounded-full bg-[#E8F5E9] text-[#1B5E20] font-bold text-[10px] uppercase border border-[#A5D6A7]">Verified</span>
              </div>
              <p class="text-xs text-[#574146]">Full access to system settings, user management, and Posko 112 control.</p>
            </div>

            <div class="space-y-3 pt-1">
              <h4 class="text-[11px] font-bold text-[#574146] uppercase tracking-wider">Key Permissions</h4>
              <ul class="space-y-2 text-xs text-[#1b1c1c]">
                <li class="flex items-center gap-2">
                  <span class="material-symbols-outlined text-[#ab2c5d] text-base">check_circle</span>
                  <span>Manage User Roles &amp; Admin Users</span>
                </li>
                <li class="flex items-center gap-2">
                  <span class="material-symbols-outlined text-[#ab2c5d] text-base">check_circle</span>
                  <span>View System Audit Logs &amp; Activity</span>
                </li>
                <li class="flex items-center gap-2">
                  <span class="material-symbols-outlined text-[#ab2c5d] text-base">check_circle</span>
                  <span>Configure Security Policies &amp; Geofence</span>
                </li>
              </ul>
            </div>
          </section>

        </div>

      </div>

    </div>
  </AdminLayout>
</template>

<script setup>
import { ref } from 'vue';
import AdminLayout from '@/layouts/AdminLayout.vue';

const twoFactorEnabled = ref(true);
const updatingPassword = ref(false);

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const activeSessions = ref([
  {
    id: 1,
    icon: 'laptop_mac',
    deviceName: 'Windows Workstation / Chrome',
    browser: 'Chrome 126.0 - Windows 11',
    location: 'Bulukumba, Indonesia',
    ipAddress: '180.242.190.12',
    lastActive: 'Active Now',
    isCurrent: true
  },
  {
    id: 2,
    icon: 'smartphone',
    deviceName: 'Android Mobile / Chrome',
    browser: 'Chrome Mobile - Android 14',
    location: 'Bulukumba, Indonesia',
    ipAddress: '180.242.191.45',
    lastActive: '2 jam yang lalu',
    isCurrent: false
  }
]);

function handleUpdatePassword() {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('Kata sandi baru dan konfirmasi kata sandi tidak cocok.');
    return;
  }
  if (passwordForm.value.newPassword.length < 8) {
    alert('Kata sandi baru minimal 8 karakter.');
    return;
  }

  updatingPassword.value = true;
  setTimeout(() => {
    updatingPassword.value = false;
    alert('Kata sandi berhasil diperbarui.');
    passwordForm.value.currentPassword = '';
    passwordForm.value.newPassword = '';
    passwordForm.value.confirmPassword = '';
  }, 800);
}

function toggle2FA() {
  const statusStr = twoFactorEnabled.value ? 'diaktifkan' : 'dinonaktifkan';
  alert(`Otentikasi Dua Faktor (2FA) berhasil ${statusStr}.`);
}

function openSetupModal() {
  alert('Silakan scan QR Code di aplikasi Google Authenticator untuk mengonfigurasi 2FA.');
}

function showBackupCodes() {
  alert('Kode Darurat Backup 2FA:\n1. 8K2L-9M3P\n2. 4N7Q-1W5X\n3. 9Y3Z-6V8B\n4. 2C4D-7E9F');
}

function logoutAllSessions() {
  if (confirm('Apakah Anda yakin ingin mengakhiri semua sesi aktif lainnya?')) {
    activeSessions.value = activeSessions.value.filter(s => s.isCurrent);
    alert('Semua sesi lainnya telah diakhiri.');
  }
}

function revokeSession(id) {
  activeSessions.value = activeSessions.value.filter(s => s.id !== id);
  alert('Sesi berhasil diakhiri.');
}
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>

