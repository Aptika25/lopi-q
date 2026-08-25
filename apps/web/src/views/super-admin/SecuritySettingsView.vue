<template>
  <AdminLayout>
    <div class="flex flex-col gap-6 w-full animate-fade-in text-left">
      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Pengaturan Keamanan &amp; Otentikasi 2FA</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Konfigurasi kebijakan keamanan global, otentikasi dua faktor (2FA), dan pembatasan presensi geofence Posko 112.</p>
        </div>
        <span class="px-3.5 py-1.5 bg-rose-50 text-rose-700 border border-rose-200/80 rounded-full text-xs font-bold shrink-0">
          TERPROTEKSI 2FA
        </span>
      </div>

      <!-- Main Content Container Card -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200/80 shadow-sm space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs">
          
          <!-- Card 1: 2FA Policy -->
          <div class="bg-slate-50 border border-slate-200/80 rounded-2xl p-5 space-y-4">
            <h3 class="font-display font-bold text-slate-900 flex items-center gap-2 text-sm">
              <span class="material-symbols-outlined text-rose-700 text-[20px]">lock_reset</span>
              <span>Kebijakan 2FA Google Authenticator</span>
            </h3>

            <div class="space-y-3">
              <div class="flex items-center justify-between p-3.5 bg-white rounded-xl border border-slate-200/80">
                <div>
                  <div class="font-bold text-slate-900">Wajibkan 2FA Semua Pengguna</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Seluruh Call Taker &amp; Admin wajib mengaktifkan TOTP saat login pertama kali.</div>
                </div>
                <input type="checkbox" checked class="w-4 h-4 text-rose-600 rounded cursor-pointer" />
              </div>

              <div class="flex items-center justify-between p-3.5 bg-white rounded-xl border border-slate-200/80">
                <div>
                  <div class="font-bold text-slate-900">Jumlah Kode Backup Pemulihan</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Setiap pengguna diberikan kode pemulihan darurat 8-karakter.</div>
                </div>
                <span class="font-mono font-bold text-rose-700 bg-rose-50 px-2.5 py-1 rounded-lg border border-rose-200">4 Kode</span>
              </div>
            </div>
          </div>

          <!-- Card 2: Geofence Radius -->
          <div class="bg-slate-50 border border-slate-200/80 rounded-2xl p-5 space-y-4">
            <h3 class="font-display font-bold text-slate-900 flex items-center gap-2 text-sm">
              <span class="material-symbols-outlined text-rose-700 text-[20px]">radar</span>
              <span>Pembatasan Radius Presensi Geofence</span>
            </h3>

            <div class="space-y-3">
              <div class="flex items-center justify-between p-3.5 bg-white rounded-xl border border-slate-200/80">
                <div>
                  <div class="font-bold text-slate-900">Batas Maksimal Radius</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Toleransi jarak presensi Call Taker dari Posko Siaga 112 Bulukumba.</div>
                </div>
                <span class="font-bold text-emerald-700 bg-emerald-50 px-2.5 py-1 rounded-lg border border-emerald-200">2.0 Meter</span>
              </div>

              <div class="flex items-center justify-between p-3.5 bg-white rounded-xl border border-slate-200/80">
                <div>
                  <div class="font-bold text-slate-900">Validasi Kamera Real-Time</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Wajib mengambil foto kamera Call Taker saat melakukan presensi siaga.</div>
                </div>
                <input type="checkbox" checked class="w-4 h-4 text-rose-600 rounded cursor-pointer" />
              </div>
            </div>
          </div>

          <!-- Card 3: Session & Password Policy -->
          <div class="bg-slate-50 border border-slate-200/80 rounded-2xl p-5 space-y-4 md:col-span-2">
            <h3 class="font-display font-bold text-slate-900 flex items-center gap-2 text-sm">
              <span class="material-symbols-outlined text-rose-700 text-[20px]">vpn_key</span>
              <span>Kebijakan Password &amp; Sesi JWT</span>
            </h3>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="p-3.5 bg-white rounded-xl border border-slate-200/80 space-y-1">
                <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Masa Berlaku Token Sesi JWT</div>
                <div class="font-bold text-slate-900">24 Jam (Otomatis Expired)</div>
              </div>
              <div class="p-3.5 bg-white rounded-xl border border-slate-200/80 space-y-1">
                <div class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Ketentuan Sandi Pengguna</div>
                <div class="font-bold text-slate-900">Minimal 8 Karakter (Kombinasi Huruf &amp; Angka)</div>
              </div>
            </div>

            <div class="pt-2 flex justify-end">
              <button 
                @click="saveSettings"
                class="px-5 py-2.5 bg-rose-700 hover:bg-rose-800 text-white font-bold text-xs rounded-xl shadow-sm transition-all cursor-pointer border-0 flex items-center gap-1.5"
              >
                <span class="material-symbols-outlined text-[16px]">save</span>
                <span>Simpan Pengaturan Keamanan</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup>
import AdminLayout from '@/layouts/AdminLayout.vue';

function saveSettings() {
  alert('Pengaturan keamanan global berhasil diperbarui.');
}
</script>
