<template>
  <div class="min-h-screen bg-background text-on-surface antialiased font-body flex select-none relative overflow-x-hidden">
    
    <!-- Mobile Sidebar Backdrop Overlay -->
    <transition name="fade">
      <div 
        v-if="mobileMenuOpen" 
        class="fixed inset-0 bg-slate-950/60 backdrop-blur-sm z-40 lg:hidden" 
        @click="mobileMenuOpen = false"
      ></div>
    </transition>

    <!-- Sidebar -->
    <aside 
      :class="[
        mobileMenuOpen ? 'translate-x-0' : '-translate-x-full',
        'fixed inset-y-0 left-0 w-64 h-screen bg-surface border-r border-outline-variant/30 z-50 flex flex-col justify-between py-5 px-3 transition-transform duration-300 ease-in-out lg:translate-x-0 lg:sticky lg:top-0 lg:h-screen lg:relative lg:inset-auto shrink-0'
      ]"
    >
      <!-- Branding with LOPI-Q Logo Icon -->
      <div class="px-3 mb-6">
        <router-link to="/admin" class="flex items-center space-x-3 text-on-surface decoration-none">
          <div class="relative flex items-center justify-center w-11 h-11 rounded-2xl bg-gradient-to-tr from-primary via-primary-container to-secondary text-on-primary font-display font-black text-xl shadow-md shrink-0">
            <span class="tracking-tighter">LQ</span>
            <span class="absolute -top-1 -right-1 flex h-3.5 w-3.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3.5 w-3.5 bg-emerald-500 border-2 border-white"></span>
            </span>
          </div>
          <div>
            <h1 class="text-base font-display font-black text-on-surface tracking-tight leading-none">LOPI<span class="text-primary">-Q</span></h1>
            <p class="text-[9px] font-bold text-primary tracking-wider uppercase mt-1">PRESENSI MAGANG BULUKUMBA</p>
          </div>
        </router-link>
      </div>

      <!-- Navigation Menu Links (Collapsible) -->
      <nav class="flex-grow flex flex-col gap-1.5 overflow-y-auto pr-0.5">
        <!-- Single Item: Dashboard -->
        <div class="relative">
          <span
            v-if="$route.path === '/admin'"
            class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-5 bg-primary rounded-r-full"
          ></span>

          <router-link
            to="/admin"
            @click="mobileMenuOpen = false"
            :class="[
              $route.path === '/admin'
                ? 'bg-primary-fixed/40 text-primary font-bold'
                : 'text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface font-medium',
              'w-full py-2.5 px-3 rounded-xl flex items-center gap-2.5 text-xs transition-all duration-200 cursor-pointer text-left border-0 bg-transparent decoration-none'
            ]"
          >
            <span :class="['material-symbols-outlined text-[18px]', $route.path === '/admin' ? 'fill text-primary' : '']">dashboard</span>
            <span>Dashboard Utama</span>
          </router-link>
        </div>

        <!-- Collapsible Menu Groups -->
        <div v-for="group in menuGroups" :key="group.key" class="relative flex flex-col">
          <!-- Collapsible Header Button -->
          <button
            @click="toggleGroup(group.key)"
            :class="[
              isGroupActive(group)
                ? 'text-primary font-bold bg-primary-fixed/20'
                : 'text-on-surface hover:bg-surface-container-low font-semibold',
              'w-full py-2.5 px-3 rounded-xl flex items-center gap-2.5 text-xs transition-all duration-200 cursor-pointer text-left border-0 bg-transparent'
            ]"
          >
            <span :class="['material-symbols-outlined text-[18px]', isGroupActive(group) ? 'text-primary' : 'text-on-surface-variant']">{{ group.icon }}</span>
            <span class="flex-1 truncate">{{ group.title }}</span>
            <span 
              class="material-symbols-outlined text-[16px] transition-transform duration-200 text-outline shrink-0"
              :class="{ 'rotate-180': expandedGroups[group.key] }"
            >
              keyboard_arrow_down
            </span>
          </button>

          <!-- Dropdown Items -->
          <div 
            v-if="expandedGroups[group.key]" 
            class="pl-3 flex flex-col gap-1 mt-1 transition-all duration-200 animate-fade-in"
          >
            <div v-for="item in group.items" :key="item.to" class="relative">
              <!-- Active child marker -->
              <span
                v-if="$route.path === item.to"
                class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-4 bg-primary rounded-r-full"
              ></span>

              <router-link
                :to="item.to"
                @click="mobileMenuOpen = false"
                :class="[
                  $route.path === item.to
                    ? 'bg-primary-fixed/40 text-primary font-bold'
                    : 'text-on-surface-variant hover:bg-surface-container-low hover:text-on-surface font-medium',
                  'w-full py-2 px-3 rounded-xl flex items-center gap-2 text-[11px] transition-all duration-200 cursor-pointer text-left border-0 bg-transparent decoration-none'
                ]"
              >
                <span :class="['material-symbols-outlined text-[16px]', $route.path === item.to ? 'fill text-primary' : '']">{{ item.icon }}</span>
                <span class="truncate">{{ item.label }}</span>
              </router-link>
            </div>
          </div>
        </div>
      </nav>

      <!-- Bottom User Section & Logout -->
      <div class="border-t border-outline-variant/40 pt-3.5 flex flex-col gap-3">
        <div class="flex items-center gap-2.5 px-3">
          <div class="h-9 w-9 rounded-full bg-gradient-to-br from-primary to-primary-container text-on-primary flex items-center justify-center font-bold text-sm shadow-sm shrink-0">
            {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'M' }}
          </div>
          <div class="flex flex-col justify-center min-w-0 text-left">
            <span class="text-xs font-bold text-on-surface truncate leading-tight block">{{ authStore.user?.name || 'Muhammad Aswan, S.T.' }}</span>
            <span class="text-[9px] text-on-surface-variant font-mono uppercase truncate leading-none mt-0.5 block">SUPERADMIN</span>
          </div>
        </div>

        <button 
          @click="handleLogout"
          class="w-full py-2 px-3 bg-red-500/10 hover:bg-red-600 text-red-600 hover:text-white text-xs rounded-xl border border-red-500/20 hover:border-red-600 transition-all active:scale-[0.96] flex items-center justify-center gap-2 cursor-pointer duration-200 font-semibold"
        >
          <span>Keluar Sesi</span>
          <span class="material-symbols-outlined text-[16px]">logout</span>
        </button>
      </div>
    </aside>

    <!-- Right Side: Header + Content Wrapper -->
    <div class="flex-grow flex flex-col min-w-0 h-screen overflow-y-auto relative bg-background">
      
      <!-- Top Navbar -->
      <header class="sticky top-0 z-30 bg-surface/95 backdrop-blur-md border-b border-outline-variant/30 py-3 px-6 flex items-center justify-between shadow-xs min-h-16">
        
        <!-- Left side: Hamburger Toggle & Title -->
        <div class="flex items-center gap-3 text-left">
          <button 
            @click="mobileMenuOpen = !mobileMenuOpen" 
            class="lg:hidden p-2 rounded-xl hover:bg-surface-container-low transition-all active:scale-95 text-on-surface-variant cursor-pointer flex items-center justify-center border-0 bg-transparent"
          >
            <span class="material-symbols-outlined text-[22px]">menu</span>
          </button>
          
          <div>
            <h2 class="text-sm md:text-base font-bold text-on-surface capitalize leading-tight">
              {{ pageTitle }}
            </h2>
            <div class="hidden md:flex items-center gap-1 text-[10px] text-on-surface-variant font-medium">
              <span>Portal Admin</span>
              <span class="material-symbols-outlined text-[10px]">chevron_right</span>
              <span class="capitalize">{{ breadcrumbSub }}</span>
            </div>
          </div>
        </div>

        <!-- Right side: Online Status -->
        <div class="flex items-center gap-4">
          <span class="inline-flex items-center px-3 py-1 rounded-full text-[10px] font-bold bg-emerald-500/10 text-emerald-700 border border-emerald-500/20 shadow-xs">
            <span class="w-1.5 h-1.5 mr-1.5 rounded-full bg-emerald-500 animate-ping"></span>
            Sistem Online
          </span>
        </div>
      </header>

      <!-- Content Area -->
      <main class="p-4 sm:p-6 flex-grow flex flex-col gap-6 max-w-[1750px] w-full mx-auto">
        <slot />
      </main>

      <!-- Footer -->
      <footer class="mt-20 border-t border-outline-variant/30 bg-surface/80 backdrop-blur-md py-6">
        <div class="max-w-7xl mx-auto px-4 md:px-10 xl:px-16 text-center text-xs">
          <p class="text-on-surface-variant font-medium">
            &copy; {{ new Date().getFullYear() }} <span class="font-semibold text-on-surface">LOPI-Q Kabupaten Bulukumba</span>. Developed by 
            <a href="https://diskominfo.bulukumbakab.go.id" target="_blank" rel="noopener noreferrer" class="font-semibold text-primary hover:underline">Diskominfo dan Persandian Kab. Bulukumba</a>.
          </p>
        </div>
      </footer>
    </div>

    <!-- ========== MODAL KONFIRMASI KELUAR SESI ========== -->
    <Teleport to="body">
      <div v-if="showLogoutModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs animate-in fade-in duration-200 select-none">
        <div class="bg-surface rounded-3xl border border-outline-variant/40 shadow-2xl max-w-sm w-full p-6 text-center space-y-4 animate-in zoom-in-95 duration-150">
          <div class="w-14 h-14 rounded-full bg-rose-100 text-rose-600 flex items-center justify-center mx-auto shadow-inner">
            <span class="material-symbols-outlined text-3xl">logout</span>
          </div>
          
          <div class="space-y-1">
            <h3 class="text-base font-extrabold text-on-surface">Konfirmasi Keluar Sesi</h3>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Apakah Anda yakin ingin keluar dari akun <strong class="text-on-surface">{{ authStore.user?.name || 'Admin LOPI-Q' }}</strong>? Anda harus memasukkan kredensial login kembali untuk masuk.
            </p>
          </div>

          <div class="grid grid-cols-2 gap-2.5 pt-2">
            <button
              @click="showLogoutModal = false"
              class="w-full py-2.5 px-4 bg-surface-container hover:bg-surface-container-high text-on-surface font-bold text-xs rounded-xl transition-all cursor-pointer border-0"
            >
              Batal
            </button>
            <button
              @click="executeLogout"
              class="w-full py-2.5 px-4 bg-gradient-to-r from-primary to-primary-container hover:opacity-90 text-on-primary font-extrabold text-xs rounded-xl shadow-md transition-all cursor-pointer border-0 active:scale-95"
            >
              Ya, Keluar
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();
const mobileMenuOpen = ref(false);
const showLogoutModal = ref(false);

const menuGroups = [
  {
    key: 'users',
    title: 'Kelola Akun & Pengguna',
    icon: 'manage_accounts',
    items: [
      { to: '/admin/admins', label: 'Administrator & Pembimbing', icon: 'admin_panel_settings' },
      { to: '/admin/interns', label: 'Peserta Magang / Intern', icon: 'groups' }
    ]
  },
  {
    key: 'attendance',
    title: 'Presensi & Aktivitas Magang',
    icon: 'fact_check',
    items: [
      { to: '/admin/attendance-recap', label: 'Rekapitulasi Kehadiran', icon: 'assignment' },
      { to: '/admin/location', label: 'Lokasi & QR Code Presensi', icon: 'distance' },
      { to: '/admin/activities', label: 'Jurnal & Laporan Aktivitas', icon: 'article' }
    ]
  },
  {
    key: 'security',
    title: 'Keamanan & Audit Sistem',
    icon: 'security',
    items: [
      { to: '/admin/activity-logs', label: 'Log Audit Aktivitas', icon: 'history_toggle_off' },
      { to: '/admin/security', label: 'Konfigurasi Keamanan & 2FA', icon: 'shield' }
    ]
  }
];

// Expanded groups state
const expandedGroups = ref({
  users: false,
  attendance: false,
  security: false
});

// Auto-expand group if any menu item inside it matches current route path
watch(() => route.path, (currentPath) => {
  for (const group of menuGroups) {
    if (group.items.some(item => item.to === currentPath)) {
      expandedGroups.value[group.key] = true;
    }
  }
}, { immediate: true });

function toggleGroup(key) {
  expandedGroups.value[key] = !expandedGroups.value[key];
}

function isGroupActive(group) {
  return group.items.some(item => item.to === route.path);
}

const pageTitle = computed(() => {
  const path = route.path;
  if (path === '/admin/admins') return 'Administrator & Pembimbing';
  if (path === '/admin/interns') return 'Peserta Magang / Intern';
  if (path === '/admin/attendance-recap') return 'Rekapitulasi Kehadiran';
  if (path === '/admin/location') return 'Lokasi & QR Code Presensi';
  if (path === '/admin/activities') return 'Jurnal & Laporan Aktivitas';
  if (path === '/admin/activity-logs') return 'Log Audit Aktivitas';
  if (path === '/admin/security') return 'Konfigurasi Keamanan & 2FA';
  return 'Dashboard Utama';
});

const breadcrumbSub = computed(() => {
  const path = route.path;
  if (path === '/admin/admins') return 'Admin';
  if (path === '/admin/interns') return 'Intern';
  if (path === '/admin/attendance-recap') return 'Presensi';
  if (path === '/admin/location') return 'Geofence';
  if (path === '/admin/activities') return 'Jurnal';
  if (path === '/admin/activity-logs') return 'Audit';
  if (path === '/admin/security') return 'Keamanan';
  return 'Dashboard';
});

function handleLogout() {
  showLogoutModal.value = true;
}

function executeLogout() {
  showLogoutModal.value = false;
  authStore.logout();
  window.location.href = '/login';
}
</script>

<style scoped>
.material-symbols-outlined {
  font-variation-settings: 'FILL' 0, 'wght' 500, 'GRAD' 0, 'opsz' 24;
}
.material-symbols-outlined.fill {
  font-variation-settings: 'FILL' 1, 'wght' 600, 'GRAD' 0, 'opsz' 24;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
