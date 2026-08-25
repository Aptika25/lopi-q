<template>
  <div class="w-full bg-background text-on-background font-body flex flex-col items-center select-none">
    
    <!-- Hero Section -->
    <section class="w-full flex flex-col md:flex-row items-center justify-between gap-12 py-8 md:py-16">
      <div class="flex-1 flex flex-col gap-6 items-start text-left">
        <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-surface-container-lowest rounded-lg border border-outline-variant">
          <span class="material-symbols-outlined text-primary text-[18px]">verified_user</span>
          <span class="text-xs font-bold text-on-surface-variant">Sistem Keamanan Terverifikasi</span>
        </div>

        <h1 class="text-3xl sm:text-5xl font-display font-black text-on-surface leading-tight tracking-tight">
          Sistem Presensi Geofence Call Taker <span class="text-primary">LOPI-Q</span>
        </h1>

        <p class="text-sm sm:text-base font-body text-on-surface-variant leading-relaxed max-w-xl">
          Sistem absensi berbasis lokasi untuk petugas Call Center 112. Pastikan Anda berada dalam radius 2.0 meter dari pusat komando sebelum melakukan pemindaian QR untuk mengaktifkan status siaga.
        </p>

        <div class="flex flex-col sm:flex-row gap-4 mt-2 w-full sm:w-auto">
          <router-link 
            to="/calltaker/scan" 
            class="bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs sm:text-sm px-6 py-3.5 rounded-xl shadow-md transition-all flex items-center justify-center gap-2 decoration-none w-full sm:w-auto border-0"
          >
            <span class="material-symbols-outlined text-[20px]">qr_code_scanner</span>
            <span>Buka Kamera Scan Presensi</span>
          </router-link>

          <router-link 
            v-if="!authStore.isAuthenticated"
            to="/login" 
            class="bg-surface-container-lowest text-primary border border-outline-variant hover:bg-surface-container-low font-bold text-xs sm:text-sm px-6 py-3.5 rounded-xl shadow-xs transition-colors flex items-center justify-center gap-2 decoration-none w-full sm:w-auto"
          >
            <span class="material-symbols-outlined text-[20px]">login</span>
            <span>Masuk ke Sistem (Login)</span>
          </router-link>

          <router-link 
            v-else-if="authStore.isAdmin"
            to="/admin" 
            class="bg-surface-container-lowest text-primary border border-outline-variant hover:bg-surface-container-low font-bold text-xs sm:text-sm px-6 py-3.5 rounded-xl shadow-xs transition-colors flex items-center justify-center gap-2 decoration-none w-full sm:w-auto"
          >
            <span class="material-symbols-outlined text-[20px]">dashboard</span>
            <span>Dashboard Super Admin</span>
          </router-link>
        </div>
      </div>

      <!-- Right Side: Hero Image Container -->
      <div class="flex-1 relative w-full flex justify-center md:justify-end">
        <div class="relative w-full max-w-sm aspect-[9/16] rounded-2xl overflow-hidden shadow-xl border-8 border-surface-container-lowest bg-surface-container-low">
          <img 
            alt="Smartphone scanning QR code" 
            class="w-full h-full object-cover" 
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuBMZ16TENRe5WDhN2GeDwz3PnLVS5RRNM0-KzK0kE50tRQdN_-ETgKY9Jjlfuk0PA-uTZz0It2_Cvaveov_0sq7sgED4CMXibe32w8wtqPk39loeMJkOs8UzAPtMwrBhxICYBB-8iN8G1TPSNMVj-4bBQGXBf2LagAYgU3zdroZYOEh-95joJnd17mlbG8l_u34wlW28Fowj-nJTnpmYspSKdkUy6TLxP6drWdbmMjHErS21B5AOw0t"
          />
        </div>
      </div>
    </section>

    <!-- Features Bento Grid -->
    <section class="w-full py-12 border-t border-outline-variant/30">
      <div class="text-center mb-12 space-y-2">
        <h2 class="text-2xl sm:text-3xl font-display font-black text-on-surface tracking-tight">Fitur Utama Presensi</h2>
        <p class="text-xs sm:text-sm text-on-surface-variant max-w-2xl mx-auto">Sistem yang dirancang untuk memastikan kehadiran petugas dengan tingkat akurasi dan keamanan tinggi.</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto justify-center">
        <!-- Feature 1 -->
        <div class="bg-surface-container-lowest p-8 rounded-3xl shadow-sm border border-outline-variant hover:shadow-md transition-all flex flex-col gap-4 items-center text-center">
          <div class="w-12 h-12 rounded-2xl bg-secondary-fixed text-on-secondary-fixed-variant flex items-center justify-center mb-2">
            <span class="material-symbols-outlined text-[24px]">my_location</span>
          </div>
          <h3 class="text-lg font-display font-bold text-on-surface">Geofence Radius 2.0m</h3>
          <p class="text-xs text-on-surface-variant leading-relaxed">Sistem akan menolak presensi jika lokasi Anda berada di luar radius 2 meter dari pusat komando LOPI-Q.</p>
        </div>

        <!-- Feature 2 -->
        <div class="bg-surface-container-lowest p-8 rounded-3xl shadow-sm border border-outline-variant hover:shadow-md transition-all flex flex-col gap-4 items-center text-center">
          <div class="w-12 h-12 rounded-2xl bg-primary-fixed text-primary flex items-center justify-center mb-2">
            <span class="material-symbols-outlined text-[24px]">qr_code_scanner</span>
          </div>
          <h3 class="text-lg font-display font-bold text-on-surface">Pemindaian QR Kamera</h3>
          <p class="text-xs text-on-surface-variant leading-relaxed">Gunakan kamera perangkat Anda untuk memindai kode QR dinamis yang tersedia di layar utama pos jaga.</p>
        </div>
      </div>
    </section>

    <!-- Bottom CTA Banner -->
    <section class="w-full mt-4 mb-8">
      <div class="bg-surface-container-lowest rounded-3xl p-8 sm:p-10 flex flex-col md:flex-row items-center justify-between gap-8 border border-outline-variant relative overflow-hidden shadow-xs">
        <div class="absolute top-0 right-0 w-64 h-64 bg-primary-container/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/4 pointer-events-none"></div>
        <div class="relative z-10 flex flex-col gap-2 text-center md:text-left">
          <h3 class="text-xl sm:text-2xl font-display font-black text-on-surface tracking-tight">Siap Melakukan Presensi Siaga 112?</h3>
          <p class="text-xs sm:text-sm text-on-surface-variant">Pastikan lokasi Anda aktif dan arahkan kamera ke QR Code.</p>
        </div>

        <router-link 
          to="/calltaker/scan"
          class="relative z-10 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs sm:text-sm px-7 py-3.5 rounded-xl shadow-md transition-colors flex items-center gap-2 whitespace-nowrap decoration-none border-0"
        >
          <span>Mulai Presensi Sekarang</span>
          <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
</script>
