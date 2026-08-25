<template>
  <div class="w-full min-h-screen bg-surface text-on-surface font-body select-none relative overflow-hidden">
    <!-- Decorative Glowing Blobs -->
    <div class="absolute top-0 right-0 -mr-32 -mt-32 w-[600px] h-[600px] bg-primary-container/20 rounded-full blur-[100px] opacity-70 pointer-events-none"></div>
    <div class="absolute bottom-0 left-0 -ml-32 -mb-32 w-[500px] h-[500px] bg-primary-fixed/40 rounded-full blur-[80px] opacity-70 pointer-events-none"></div>

    <div class="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-8 pb-20 relative z-10">
      <!-- HERO SECTION -->
      <section class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-16 items-center py-8 lg:py-16">
        <!-- Left Hero Content -->
        <div class="flex flex-col gap-6 max-w-2xl text-left">
          <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-white border border-outline-variant/40 w-fit shadow-xs">
            <span class="w-2.5 h-2.5 rounded-full bg-primary-container animate-pulse"></span>
            <span class="font-label-sm text-xs font-bold text-primary">Sistem Otomasi Magang Terverifikasi</span>
          </div>

          <h1 class="font-display text-2xl sm:text-4xl lg:text-[42px] lg:leading-[50px] text-on-surface font-extrabold tracking-tight">
            Presensi & Logbook <br class="hidden sm:inline" />
            Aktivitas Magang <span class="whitespace-nowrap gradient-text">LOPI-Q</span>
          </h1>

          <p class="text-sm sm:text-base font-body text-on-surface-variant leading-relaxed">
            <strong>LOPI-Q</strong> (Layanan Otomasi Presensi, Informasi Aktivitas, dan Quality Monitoring Magang) menghadirkan pelacakan waktu real-time, geofencing presisi radius 2.0 meter, dan analitik jurnal harian untuk Siswa & Mahasiswa Magang Diskominfo Kabupaten Bulukumba.
          </p>

          <!-- Action Buttons -->
          <div class="flex flex-wrap items-center gap-4 mt-2">
            <router-link 
              to="/intern/scan" 
              class="px-6 py-3.5 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs sm:text-sm rounded-xl shadow-[0_8px_20px_rgba(240,98,146,0.25)] hover:translate-y-[-2px] transition-all duration-300 flex items-center gap-2 decoration-none border-0 cursor-pointer"
            >
              <span class="material-symbols-outlined text-[20px]">qr_code_scanner</span>
              <span>Buka Kamera Scan Presensi</span>
            </router-link>

            <router-link 
              v-if="!authStore.isAuthenticated"
              to="/login" 
              class="px-6 py-3.5 bg-white text-primary border border-outline-variant hover:bg-surface transition-colors font-bold text-xs sm:text-sm rounded-xl shadow-xs flex items-center gap-2 decoration-none cursor-pointer"
            >
              <span class="material-symbols-outlined text-[20px]">login</span>
              <span>Masuk ke Sistem (Login)</span>
            </router-link>

            <router-link 
              v-else-if="authStore.isAdmin"
              to="/admin" 
              class="px-6 py-3.5 bg-white text-primary border border-outline-variant hover:bg-surface transition-colors font-bold text-xs sm:text-sm rounded-xl shadow-xs flex items-center gap-2 decoration-none cursor-pointer"
            >
              <span class="material-symbols-outlined text-[20px]">dashboard</span>
              <span>Dashboard Super Admin</span>
            </router-link>
          </div>

          <!-- Trusted Avatars Stack -->
          <div class="flex items-center gap-4 mt-4 text-on-surface-variant/80">
            <div class="flex -space-x-3">
              <img class="w-10 h-10 rounded-full border-2 border-white object-cover shadow-xs" alt="User Avatar" src="https://lh3.googleusercontent.com/aida-public/AB6AXuC6cLDkwdfsRMFinQjZkMDNnGM1aL9qRxOJR7jWQkM1BQHt72SVLE0kTOOnXodLuy_5EvrCNSSGKrjZkgV9MWj0rSlY7CCc5aYP1kFqyAQ6aCbD9uJeyj0bUXwP4I_EF2Y7wY1UVSPjp15STJ5Gif1aS2GrQ6HT0uwh09WZ0q4ipf4RILdrXrnEgsoy7mPqPYPVuJtr_xnn-vBrvWeoml48wKbu503PsMGNv-pU6fKv3IoXe5srF6Upbg">
              <img class="w-10 h-10 rounded-full border-2 border-white object-cover shadow-xs" alt="User Avatar" src="https://lh3.googleusercontent.com/aida-public/AB6AXuDu6w-fD_NyQjNBXRG_iHwYJBWvVbrgkzAsTwV3cEWnjoOtti7iYwAGfPTKQG8iHs09b4nEyBDa6iflpp8mJjS0YaPxE3WZp2Iyl0yI26bwQCQctKauMrtWx8EsGXMpxGkOJHWL3srQClbJ4pRXELIFgyz0Bel_cAv0s0i8MzJVy6ZHhFg0xnEbZZaVWHfR47Stfe33XH_JWY2C0Ji9gu4d2XTiGENDMHkae_Bpgaz-R3YcDFsQ73AuMw">
              <img class="w-10 h-10 rounded-full border-2 border-white object-cover shadow-xs" alt="User Avatar" src="https://lh3.googleusercontent.com/aida-public/AB6AXuDITpk-Vcdx2W1vPkb9zCTCWglNyEdkrt0eOBHLp2-dv-VCUUFrYvXkiP4ylgziGnf3z79c_geFmha6Jyxra0Nu2WIrEWxZQI1RFprn0eRVVEnY-U_rNvh0brRroDiMge7mWiyluek0QYV9CKtweIQZQHfksgM-sHRTQjt9W46bq59QEfcDPbRsAcLcjWHBEbXTsi4lDAlqaxLnwQSUpFWfdHy7QiV807RfFllz4HQ4CfRIAOcmwI6olg">
            </div>
            <span class="text-xs font-semibold text-on-surface-variant">Dipercaya oleh Siswa & Mahasiswa Magang Diskominfo Bulukumba</span>
          </div>
        </div>

        <!-- Right Side: Glassmorphism Bento Grid Mockup (Matching Attached Image) -->
        <div class="relative w-full h-[480px] sm:h-[540px] lg:h-[580px] perspective-1000">
          <div class="absolute inset-0 bg-gradient-to-tr from-primary-container/20 to-primary/10 rounded-3xl transform rotate-2 scale-105 opacity-60"></div>
          
          <div class="glass-panel soft-shadow w-full h-full rounded-3xl relative overflow-hidden flex flex-col p-5 z-10 border border-white/80 bg-white/75 backdrop-blur-xl">
            <!-- Mockup Top Bar -->
            <div class="flex items-center justify-between pb-4 border-b border-outline-variant/30 mb-4">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-primary-container text-white flex items-center justify-center shadow-xs">
                  <span class="material-symbols-outlined text-[20px]">dashboard</span>
                </div>
                <span class="font-display text-base font-black text-primary tracking-tight">Overview Presensi</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="px-2.5 py-1 rounded-full bg-green-100 text-green-700 font-bold text-[11px] flex items-center gap-1.5 border border-green-200">
                  <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span> Live Active
                </span>
              </div>
            </div>

            <!-- Inner Bento Grid -->
            <div class="grid grid-cols-2 gap-4 flex-grow">
              <!-- Stats Card 1 -->
              <div class="bg-white rounded-2xl p-4 border border-outline-variant/30 shadow-xs flex flex-col justify-between">
                <div class="flex items-center gap-2 text-on-surface-variant mb-1">
                  <span class="material-symbols-outlined text-[18px] text-primary">group</span>
                  <span class="text-xs font-bold">Hadir Hari Ini</span>
                </div>
                <div class="flex items-end justify-between mt-2">
                  <span class="font-display text-3xl font-black text-on-surface">142</span>
                  <span class="text-green-600 font-bold text-xs flex items-center bg-green-50 px-1.5 py-0.5 rounded-md">
                    <span class="material-symbols-outlined text-[14px]">trending_up</span> 98%
                  </span>
                </div>
                <div class="w-full h-10 mt-3 bg-primary-fixed/20 rounded-xl overflow-hidden relative border border-primary-fixed/30">
                  <div class="absolute top-0 left-0 h-full w-[98%] bg-primary-container/20 rounded-xl"></div>
                  <svg class="absolute inset-0 w-full h-full text-primary-container" preserveAspectRatio="none" viewBox="0 0 100 100">
                    <path d="M0,100 L0,50 Q25,20 50,60 T100,30 L100,100 Z" fill="currentColor" opacity="0.3"></path>
                    <path d="M0,100 L0,60 Q25,30 50,70 T100,40 L100,100 Z" fill="currentColor" opacity="0.6"></path>
                  </svg>
                </div>
              </div>

              <!-- Map / Geofence Card -->
              <div class="bg-white rounded-2xl p-4 border border-outline-variant/30 shadow-xs col-span-1 row-span-2 relative overflow-hidden flex flex-col group cursor-pointer hover:border-primary/40 transition-all">
                <div class="flex items-center justify-between text-on-surface-variant mb-2 relative z-10">
                  <span class="text-[11px] font-bold bg-white/90 px-2.5 py-1 rounded-lg backdrop-blur border border-outline-variant/30 text-primary shadow-2xs">Live Geofence 2.0m</span>
                  <span class="material-symbols-outlined text-[18px] text-primary">location_on</span>
                </div>
                <div class="absolute inset-0 w-full h-full mt-10 rounded-xl overflow-hidden border border-slate-100">
                  <img 
                    class="w-full h-full object-cover opacity-85 mix-blend-multiply transition-transform duration-700 group-hover:scale-105" 
                    alt="Geofence Radius Map"
                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuA0rxjJfenIM3-ZdIYk5-gVK6107Jed3ke0Bnm-6FnVXyuB60l98NMdIJ0kcPBHFJuVWg7m0lVVPoZ2Z-PQ7dd-xRGEjvVOqJHxrOx7DNKrTm51OwgLczylisfyxzWN_TeKO5gCZosPthJVsQAmlVEpHSuM_beeus5LO9Hydlcu2m6UBE9K195gY-M9AJsccCt3Qp6kcoApctIf_B_3SDKYM7QB_tbu5pS1oI62EjaUnTYQQmE99kN4aQ"
                  />
                  <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                    <div class="w-24 h-24 rounded-full border-2 border-primary-container bg-primary-container/15 animate-ping absolute -top-12 -left-12"></div>
                    <div class="w-3.5 h-3.5 bg-primary-container rounded-full shadow-[0_0_12px_rgba(240,98,146,0.9)] relative z-10 border-2 border-white"></div>
                  </div>
                </div>
              </div>

              <!-- Activity List Card -->
              <div class="bg-white rounded-2xl p-4 border border-outline-variant/30 shadow-xs flex flex-col">
                <span class="text-xs font-bold text-on-surface-variant mb-2.5">Aktivitas Terkini</span>
                <div class="space-y-2.5">
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-lg bg-primary-fixed/40 text-primary flex items-center justify-center font-bold text-[11px] shrink-0">AR</div>
                    <div class="flex flex-col min-w-0">
                      <span class="font-bold text-xs text-on-surface leading-tight truncate">Andi R.</span>
                      <span class="text-[10px] text-on-surface-variant font-medium">Clock In - 08:02 AM</span>
                    </div>
                  </div>
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-lg bg-primary-fixed/40 text-primary flex items-center justify-center font-bold text-[11px] shrink-0">BW</div>
                    <div class="flex flex-col min-w-0">
                      <span class="font-bold text-xs text-on-surface leading-tight truncate">Budi W.</span>
                      <span class="text-[10px] text-on-surface-variant font-medium">Clock In - 08:15 AM</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- FEATURES SECTION -->
      <section class="w-full py-16 border-t border-outline-variant/30 relative z-10" id="fitur">
        <div class="text-center max-w-2xl mx-auto mb-16 space-y-3">
          <span class="text-xs font-bold text-primary tracking-widest uppercase block">Kemampuan Platform LOPI-Q</span>
          <h2 class="font-display text-2xl sm:text-4xl font-extrabold text-on-surface tracking-tight">Semua yang Anda Butuhkan untuk Mengelola Magang</h2>
          <p class="text-xs sm:text-sm text-on-surface-variant leading-relaxed">Desain presisi yang mengombinasikan kecepatan presensi QR dengan geofencing 2.0 meter.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="glass-panel p-8 rounded-3xl hover:-translate-y-2 transition-all duration-300 group shadow-xs">
            <div class="w-12 h-12 rounded-2xl bg-primary-fixed/40 text-primary flex items-center justify-center mb-6 group-hover:bg-primary-container group-hover:text-white transition-all">
              <span class="material-symbols-outlined text-[24px]">bolt</span>
            </div>
            <h3 class="text-lg font-display font-bold text-on-surface mb-2">Pelacakan Real-time</h3>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Pantau status kehadiran peserta magang secara instan. Data tersinkronisasi ke cloud tanpa delay, memastikan akurasi operasional.
            </p>
          </div>

          <!-- Feature 2 -->
          <div class="glass-panel p-8 rounded-3xl hover:-translate-y-2 transition-all duration-300 group shadow-xs">
            <div class="w-12 h-12 rounded-2xl bg-primary-fixed/40 text-primary flex items-center justify-center mb-6 group-hover:bg-primary-container group-hover:text-white transition-all">
              <span class="material-symbols-outlined text-[24px]">my_location</span>
            </div>
            <h3 class="text-lg font-display font-bold text-on-surface mb-2">Geofencing Presisi 2.0m</h3>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Tentukan radius valid untuk absensi. Peserta magang hanya dapat melakukan clock-in jika berada di lokasi kerja Diskominfo Bulukumba.
            </p>
          </div>

          <!-- Feature 3 -->
          <div class="glass-panel p-8 rounded-3xl hover:-translate-y-2 transition-all duration-300 group shadow-xs">
            <div class="w-12 h-12 rounded-2xl bg-primary-fixed/40 text-primary flex items-center justify-center mb-6 group-hover:bg-primary-container group-hover:text-white transition-all">
              <span class="material-symbols-outlined text-[24px]">monitoring</span>
            </div>
            <h3 class="text-lg font-display font-bold text-on-surface mb-2">Analitik & Logbook Harian</h3>
            <p class="text-xs text-on-surface-variant leading-relaxed">
              Pencatatan aktivitas jurnal harian magang secara otomatis. Dapatkan laporan terperinci mengenai jam kerja dan kegiatan peserta dalam satu klik.
            </p>
          </div>
        </div>
      </section>

      <!-- BOTTOM CTA BANNER -->
      <section class="w-full mt-4 mb-8 relative z-10" id="kontak">
        <div class="glass-panel rounded-3xl p-8 sm:p-12 flex flex-col md:flex-row items-center justify-between gap-8 relative overflow-hidden shadow-sm bg-white/80 border border-outline-variant/40">
          <div class="absolute top-0 right-0 w-64 h-64 bg-primary-container/15 rounded-full blur-3xl -translate-y-1/2 translate-x-1/4 pointer-events-none"></div>
          <div class="relative z-10 flex flex-col gap-2 text-center md:text-left">
            <h3 class="text-xl sm:text-2xl font-display font-black text-on-surface tracking-tight">Siap Melakukan Presensi & Logbook Magang?</h3>
            <p class="text-xs sm:text-sm text-on-surface-variant">Pastikan lokasi GPS Anda aktif dan arahkan kamera ke Kode QR di lokasi magang Diskominfo Bulukumba.</p>
          </div>

          <router-link 
            to="/intern/scan"
            class="relative z-10 bg-primary-container hover:bg-on-primary-fixed-variant text-on-primary font-bold text-xs sm:text-sm px-7 py-3.5 rounded-xl shadow-md hover:shadow-lg transition-all flex items-center gap-2 whitespace-nowrap decoration-none border-0 cursor-pointer"
          >
            <span>Mulai Presensi Sekarang</span>
            <span class="material-symbols-outlined text-[18px]">arrow_forward</span>
          </router-link>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
</script>
