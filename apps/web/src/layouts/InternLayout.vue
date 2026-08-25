<template>
  <div class="min-h-screen bg-slate-50 text-slate-800 antialiased font-body flex select-none relative overflow-x-hidden pb-24 lg:pb-0">

    <!-- ========== SIDEBAR (desktop only, hidden on mobile) ========== -->
    <aside class="hidden lg:flex fixed inset-y-0 left-0 w-64 h-screen bg-white border-r border-slate-200/80 z-50 flex-col justify-between py-5 px-3 lg:sticky lg:top-0 lg:h-screen lg:relative shrink-0">
      <!-- Branding with 112 Logo Icon -->
      <div class="px-3 mb-6">
        <router-link to="/intern/dashboard" class="flex items-center space-x-3 text-slate-900 decoration-none">
          <div class="relative flex items-center justify-center w-11 h-11 rounded-2xl bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white font-display font-black text-xl shadow-md shrink-0">
            <span class="tracking-tighter">112</span>
            <span class="absolute -top-1 -right-1 flex h-3.5 w-3.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3.5 w-3.5 bg-emerald-500 border-2 border-white"></span>
            </span>
          </div>
          <div>
            <h1 class="text-base font-display font-black text-slate-900 tracking-tight leading-none">GARDA<span class="text-rose-700">112</span></h1>
            <p class="text-[9px] font-bold text-rose-800 tracking-wider uppercase mt-1">PESERTA MAGANG PANEL</p>
          </div>
        </router-link>
      </div>

      <!-- Navigation Menu Links -->
      <nav class="flex-grow flex flex-col gap-1 overflow-y-auto">
        <div class="text-[9px] font-extrabold text-slate-400 uppercase tracking-widest px-4 mt-4 mb-1.5 select-none">
          Menu Operasional Peserta Magang
        </div>

        <div v-for="item in menuItems" :key="item.to" class="relative">
          <span
            v-if="$route.path === item.to"
            class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-rose-700 rounded-r-full"
          ></span>
          <router-link
            :to="item.to"
            :class="[
              $route.path === item.to
                ? 'bg-rose-50 text-rose-700 font-semibold'
                : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900 font-medium',
              'w-full py-3 px-4 rounded-xl flex items-center gap-2.5 text-xs transition-all duration-200 cursor-pointer text-left border-0 bg-transparent decoration-none'
            ]"
          >
            <span :class="['material-symbols-outlined text-[18px]', $route.path === item.to ? 'fill text-rose-700' : '']">{{ item.icon }}</span>
            <span>{{ item.label }}</span>
          </router-link>
        </div>
      </nav>

      <!-- Bottom User Section & Logout -->
      <div class="border-t border-slate-200/80 pt-3.5 flex flex-col gap-3">
        <div class="flex items-center gap-2.5 px-3">
          <div class="h-9 w-9 rounded-full bg-gradient-to-br from-rose-700 to-rose-600 text-white flex items-center justify-center font-bold text-sm shadow-sm shrink-0">
            {{ authStore.user?.name ? authStore.user.name.charAt(0).toUpperCase() : 'M' }}
          </div>
          <div class="flex flex-col justify-center min-w-0 text-left">
            <span class="text-xs font-bold text-slate-900 truncate leading-tight block">{{ authStore.user?.name || 'Peserta Magang' }}</span>
            <span class="text-[9px] text-slate-400 font-mono uppercase truncate leading-none mt-0.5 block">PESERTA MAGANG</span>
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

    <!-- ========== RIGHT SIDE: Header + Content ========== -->
    <div class="flex-grow flex flex-col min-w-0 h-screen overflow-y-auto relative bg-slate-50">

      <!-- ===== TOP NAVBAR ===== -->
      <header class="sticky top-0 z-30 bg-white/95 backdrop-blur-md border-b border-slate-200/60 shadow-xs min-h-14">

        <!-- MOBILE top navbar: logo only -->
        <div class="lg:hidden flex items-center justify-between px-4 h-14">
          <!-- Logo LOPI-Q -->
          <router-link to="/intern/dashboard" class="flex items-center gap-2.5 decoration-none">
            <div class="relative flex items-center justify-center w-9 h-9 rounded-xl bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white font-display font-black text-base shadow-md shrink-0">
              <span class="tracking-tighter">112</span>
              <span class="absolute -top-0.5 -right-0.5 flex h-2.5 w-2.5">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500 border-2 border-white"></span>
              </span>
            </div>
            <div>
              <div class="text-sm font-display font-black text-slate-900 tracking-tight leading-none">GARDA<span class="text-rose-700">112</span></div>
              <div class="text-[8px] font-bold text-rose-800 tracking-widest uppercase leading-none mt-0.5">PESERTA MAGANG PANEL</div>
            </div>
          </router-link>

          <!-- Online status pill -->
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[9px] font-bold bg-emerald-500/10 text-emerald-700 border border-emerald-500/20">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping shrink-0"></span>
            Online
          </span>
        </div>

        <!-- DESKTOP top navbar: title + breadcrumb + status -->
        <div class="hidden lg:flex items-center justify-between px-6 h-16">
          <div>
            <h2 class="text-base font-bold text-slate-900 capitalize leading-tight">{{ pageTitle }}</h2>
            <div class="flex items-center gap-1 text-[10px] text-slate-400 font-medium mt-0.5">
              <span>Peserta Magang</span>
              <span class="material-symbols-outlined text-[10px]">chevron_right</span>
              <span class="capitalize">{{ breadcrumbSub }}</span>
            </div>
          </div>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-[10px] font-bold bg-emerald-500/10 text-emerald-700 border border-emerald-500/20 shadow-xs">
            <span class="w-1.5 h-1.5 mr-1.5 rounded-full bg-emerald-500 animate-ping"></span>
            Sistem Online
          </span>
        </div>
      </header>

      <!-- Content Area -->
      <main class="p-4 lg:p-6 flex-grow flex flex-col gap-6 max-w-7xl w-full mx-auto">
        <slot />
      </main>

      <!-- Footer (desktop only) -->
      <footer class="hidden lg:block mt-20 border-t border-slate-200/60 bg-white/80 backdrop-blur-md py-6">
        <div class="max-w-7xl mx-auto px-4 md:px-10 xl:px-16 text-center text-xs">
          <p class="text-slate-400 font-medium">
            &copy; {{ new Date().getFullYear() }} <span class="font-semibold text-slate-800">LOPI-Q Kabupaten Bulukumba</span>. Developed by 
            <a href="https://diskominfo.bulukumbakab.go.id" target="_blank" rel="noopener noreferrer" class="font-semibold text-rose-700 hover:underline">Diskominfo dan Persandian Kab. Bulukumba</a>.
          </p>
        </div>
      </footer>
    </div>

    <!-- ========== MOBILE BOTTOM NAVBAR (clean & prominent HERO Scan QR) ========== -->
    <nav class="fixed bottom-0 left-0 right-0 z-50 lg:hidden">
      <!-- Clean backdrop glass -->
      <div class="absolute inset-0 bg-white/95 backdrop-blur-xl border-t border-slate-200/80 shadow-[0_-4px_24px_rgba(0,0,0,0.06)]"></div>

      <!-- Nav items container -->
      <div class="relative flex items-center justify-around px-4 pt-1.5 pb-safe-offset-2" style="padding-bottom: max(0.6rem, env(safe-area-inset-bottom))">

        <!-- 1. DASHBOARD BUTTON (Left) -->
        <router-link
          to="/intern/dashboard"
          class="flex flex-col items-center gap-0.5 flex-1 py-1 decoration-none group transition-all duration-200"
        >
          <div
            :class="[
              'w-10 h-8 rounded-xl flex items-center justify-center transition-all duration-200',
              $route.path === '/intern/dashboard' ? 'bg-rose-50 text-rose-700' : 'text-slate-400 group-hover:text-slate-600'
            ]"
          >
            <span :class="['material-symbols-outlined text-[22px]', $route.path === '/intern/dashboard' ? 'fill text-rose-700' : '']">grid_view</span>
          </div>
          <span
            :class="[
              'text-[10px] font-bold tracking-tight transition-colors',
              $route.path === '/intern/dashboard' ? 'text-rose-700 font-extrabold' : 'text-slate-400'
            ]"
          >Dashboard</span>
        </router-link>

        <!-- 2. SCAN QR BUTTON (Center - Prominent Hero Floating Button) -->
        <router-link
          to="/intern/scan"
          class="flex flex-col items-center gap-0.5 flex-1 relative decoration-none group -mt-6"
        >
          <div
            :class="[
              'w-14 h-14 rounded-full flex items-center justify-center transition-all duration-300 border-4 border-white shadow-xl',
              $route.path === '/intern/scan'
                ? 'bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white shadow-rose-600/50 scale-110 ring-4 ring-rose-500/25'
                : 'bg-gradient-to-tr from-rose-600 to-amber-500 text-white shadow-rose-500/40 hover:scale-105'
            ]"
          >
            <span class="material-symbols-outlined text-2xl fill text-white">qr_code_scanner</span>
          </div>
          <span
            :class="[
              'text-[10px] tracking-tight transition-colors mt-0.5',
              $route.path === '/intern/scan' ? 'text-rose-700 font-black' : 'text-slate-700 font-extrabold'
            ]"
          >Scan QR</span>
        </router-link>

        <!-- 3. RIWAYAT BUTTON (Right) -->
        <router-link
          to="/intern/history"
          class="flex flex-col items-center gap-0.5 flex-1 py-1 decoration-none group transition-all duration-200"
        >
          <div
            :class="[
              'w-10 h-8 rounded-xl flex items-center justify-center transition-all duration-200',
              $route.path === '/intern/history' ? 'bg-rose-50 text-rose-700' : 'text-slate-400 group-hover:text-slate-600'
            ]"
          >
            <span :class="['material-symbols-outlined text-[22px]', $route.path === '/intern/history' ? 'fill text-rose-700' : '']">history_edu</span>
          </div>
          <span
            :class="[
              'text-[10px] font-bold tracking-tight transition-colors',
              $route.path === '/intern/history' ? 'text-rose-700 font-extrabold' : 'text-slate-400'
            ]"
          >Riwayat</span>
        </router-link>

      </div>
    </nav>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const menuItems = [
  { to: '/intern/dashboard', label: 'Dashboard Status Siaga', icon: 'dashboard' },
  { to: '/intern/scan', label: 'Scan QR Presensi', icon: 'qr_code_scanner' },
  { to: '/intern/history', label: 'Laporan Kehadiran Saya', icon: 'description' }
];

const navItems = [
  { to: '/intern/dashboard', label: 'Dashboard', icon: 'grid_view' },
  { to: '/intern/scan', label: 'Scan QR', icon: 'qr_code_scanner' },
  { to: '/intern/history', label: 'Riwayat', icon: 'description' }
];

const pageTitle = computed(() => {
  const path = route.path;
  if (path === '/intern/history') return 'Laporan Kehadiran Saya';
  if (path === '/intern/scan') return 'Scan QR Presensi';
  return 'Dashboard Status Siaga';
});

const breadcrumbSub = computed(() => {
  const path = route.path;
  if (path === '/intern/history') return 'Riwayat Presensi';
  if (path === '/intern/scan') return 'Scan QR';
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
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}

.material-symbols-outlined.fill {
  font-variation-settings: 'FILL' 1, 'wght' 600, 'GRAD' 0, 'opsz' 24;
}
</style>
