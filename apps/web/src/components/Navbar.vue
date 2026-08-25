<template>
  <header class="border-b border-slate-200 bg-white/90 backdrop-blur-md sticky top-0 z-40 shadow-xs select-none">
    <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 h-16 sm:h-20 flex items-center justify-between gap-2">
      <!-- Logo & Title -->
      <router-link to="/" class="flex items-center space-x-2.5 sm:space-x-3 text-slate-900 decoration-none min-w-0">
        <div class="relative flex items-center justify-center w-9 h-9 sm:w-11 sm:h-11 rounded-xl sm:rounded-2xl bg-gradient-to-tr from-rose-700 via-rose-600 to-amber-500 text-white font-display font-black text-lg sm:text-xl shadow-md shrink-0">
          <span class="tracking-tighter">112</span>
          <span class="absolute -top-1 -right-1 flex h-3 w-3 sm:h-3.5 sm:w-3.5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 sm:h-3.5 sm:w-3.5 bg-emerald-500 border-2 border-white"></span>
          </span>
        </div>
        <div class="flex flex-col min-w-0 justify-center">
          <h1 class="text-base sm:text-xl font-display font-black text-slate-900 tracking-tight leading-none">
            GARDA<span class="text-rose-700">112</span>
          </h1>
          <p class="text-[10px] sm:text-xs font-semibold text-rose-800 tracking-wider uppercase mt-1 leading-none whitespace-nowrap">
            NTPD 112 Kab. Bulukumba
          </p>
        </div>
      </router-link>

      <!-- User Action Button Based on Role -->
      <div v-if="authStore.user" class="flex items-center shrink-0">
        <!-- Super Admin Only Button -->
        <router-link 
          v-if="authStore.isAdmin"
          to="/admin"
          class="px-3 sm:px-5 py-2 sm:py-2.5 rounded-xl text-[11px] sm:text-xs font-bold bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white shadow-md transition-all decoration-none flex items-center gap-1.5 whitespace-nowrap"
        >
          <span class="material-symbols-outlined text-[16px] sm:text-[18px]">dashboard</span>
          <span class="hidden sm:inline">Dashboard Super Admin</span>
          <span class="sm:hidden">Dashboard</span>
        </router-link>

        <!-- Call Taker Only Button -->
        <router-link 
          v-else-if="authStore.isCallTaker"
          to="/calltaker/dashboard"
          class="px-3 sm:px-5 py-2 sm:py-2.5 rounded-xl text-[11px] sm:text-xs font-bold bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white shadow-md transition-all decoration-none flex items-center gap-1.5 whitespace-nowrap"
        >
          <span class="material-symbols-outlined text-[16px] sm:text-[18px]">support_agent</span>
          <span class="hidden sm:inline">Halaman Call Taker</span>
          <span class="sm:hidden">Call Taker</span>
        </router-link>
      </div>

      <!-- Guest Login Button: Navigate to /login -->
      <div v-else class="shrink-0">
        <router-link 
          to="/login"
          class="bg-gradient-to-r from-rose-700 via-rose-600 to-amber-600 hover:from-rose-800 hover:to-amber-700 text-white text-[11px] sm:text-xs font-bold px-3.5 sm:px-5 py-2 sm:py-2.5 rounded-xl shadow-md transition-all flex items-center gap-1.5 decoration-none whitespace-nowrap"
        >
          <span class="material-symbols-outlined text-[16px] sm:text-[18px]">login</span>
          <span>Masuk</span>
        </router-link>
      </div>
    </div>
  </header>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';

const authStore = useAuthStore();
</script>
