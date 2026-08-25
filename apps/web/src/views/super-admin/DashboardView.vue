<template>
  <AdminLayout>
    <div class="space-y-6 select-none font-sans">
      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Dashboard Overview Utama</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Monitoring presensi Call Taker secara real-time, radar Geofence Posko Siaga 112, dan kontrol sistem.</p>
        </div>
      </div>
      <div class="bg-gradient-to-r from-slate-900 via-rose-950 to-slate-900 rounded-3xl p-6 sm:p-8 text-white shadow-xl relative overflow-hidden border border-rose-900/40">
        <!-- Decorative Glow Blur -->
        <div class="absolute -top-12 -right-12 w-64 h-64 bg-rose-600/20 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -bottom-12 -left-12 w-64 h-64 bg-amber-500/10 rounded-full blur-3xl pointer-events-none"></div>

        <div class="relative z-10 space-y-2">
          <div class="inline-flex items-center gap-2 px-3 py-1 bg-rose-500/20 text-rose-300 rounded-full text-xs font-bold border border-rose-500/30">
            <span class="w-2 h-2 rounded-full bg-rose-500 animate-ping"></span>
            <span>PUSAT KONTROL UTAMA LOPI-Q</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-display font-black tracking-tight text-white">
            Dashboard Super Admin Posko 112
          </h1>
          <p class="text-xs sm:text-sm text-slate-300 max-w-2xl leading-relaxed">
            Monitoring presensi Call Taker secara real-time, pantau radar Geofence Posko Siaga 112 Kabupaten Bulukumba, dan kelola otentikasi dua langkah (2FA).
          </p>
        </div>
      </div>

      <!-- KPI Summary Stat Cards Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <!-- Stat Card 1 -->
        <div class="bg-white rounded-3xl p-5 border border-slate-200 shadow-xs space-y-3 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-[11px] font-extrabold text-slate-400 uppercase tracking-wider">Total Call Taker</span>
            <div class="w-10 h-10 bg-rose-50 text-rose-700 rounded-2xl flex items-center justify-center border border-rose-100">
              <span class="material-symbols-outlined text-[22px]">headphones</span>
            </div>
          </div>
          <div>
            <div class="text-2xl font-display font-black text-slate-900">10 Petugas</div>
            <div class="text-xs text-slate-500 mt-0.5">Terdaftar dari 5 Instansi OPD</div>
          </div>
        </div>

        <!-- Stat Card 2 -->
        <div class="bg-white rounded-3xl p-5 border border-slate-200 shadow-xs space-y-3 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-[11px] font-extrabold text-slate-400 uppercase tracking-wider">Presensi Siaga Hari Ini</span>
            <div class="w-10 h-10 bg-emerald-50 text-emerald-700 rounded-2xl flex items-center justify-center border border-emerald-100">
              <span class="material-symbols-outlined text-[22px]">how_to_reg</span>
            </div>
          </div>
          <div>
            <div class="text-2xl font-display font-black text-emerald-700">100% Valid</div>
            <div class="text-xs text-slate-500 mt-0.5">Verifikasi Geofence &amp; Kamera</div>
          </div>
        </div>

        <!-- Stat Card 3 -->
        <div class="bg-white rounded-3xl p-5 border border-slate-200 shadow-xs space-y-3 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-[11px] font-extrabold text-slate-400 uppercase tracking-wider">Geofence Posko 112</span>
            <div class="w-10 h-10 bg-sky-50 text-sky-700 rounded-2xl flex items-center justify-center border border-sky-100">
              <span class="material-symbols-outlined text-[22px]">radar</span>
            </div>
          </div>
          <div>
            <div class="text-2xl font-display font-black text-slate-900">Radius {{ poskoInfo.radius }}m</div>
            <div class="text-xs text-slate-500 mt-0.5 font-mono truncate">{{ poskoInfo.lat.toFixed(4) }}, {{ poskoInfo.lng.toFixed(4) }}</div>
          </div>
        </div>

        <!-- Stat Card 4 -->
        <div class="bg-white rounded-3xl p-5 border border-slate-200 shadow-xs space-y-3 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-[11px] font-extrabold text-slate-400 uppercase tracking-wider">Keamanan 2FA Google</span>
            <div class="w-10 h-10 bg-indigo-50 text-indigo-700 rounded-2xl flex items-center justify-center border border-indigo-100">
              <span class="material-symbols-outlined text-[22px]">verified_user</span>
            </div>
          </div>
          <div>
            <div class="text-2xl font-display font-black text-indigo-900">Terproteksi</div>
            <div class="text-xs text-slate-500 mt-0.5">Otentikasi TOTP Aktif</div>
          </div>
        </div>
      </div>

      <!-- Quick Shortcuts Nav Cards -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-xs space-y-4">
        <h3 class="text-base font-display font-black text-slate-900 flex items-center gap-2">
          <span class="material-symbols-outlined text-rose-700">grid_view</span>
          <span>Akses Cepat Pengelolaan Sistem</span>
        </h3>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
          <router-link to="/admin/calltakers" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">badge</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Call Takers</div>
          </router-link>

          <router-link to="/admin/admins" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">shield_person</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Admin User</div>
          </router-link>

          <router-link to="/admin/attendance-recap" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">assessment</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Rekap Presensi</div>
          </router-link>

          <router-link to="/admin/location" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">pin_drop</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Lokasi Posko</div>
          </router-link>

          <router-link to="/admin/schedules" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">calendar_month</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Jadwal Shift</div>
          </router-link>

          <router-link to="/admin/security" class="p-4 bg-slate-50 hover:bg-rose-50/60 border border-slate-200 hover:border-rose-300 rounded-2xl text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-rose-700 rounded-xl flex items-center justify-center mx-auto shadow-2xs group-hover:scale-110 transition-transform">
              <span class="material-symbols-outlined text-[20px]">lock_reset</span>
            </div>
            <div class="text-xs font-bold text-slate-800 group-hover:text-rose-900">Keamanan 2FA</div>
          </router-link>
        </div>
      </div>

      <!-- Pending Leave & Shift Requests Section + Geofence Posko Monitor Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        
        <!-- Left: Leave & Shift Replacement Approval Widget (7 cols) -->
        <div class="lg:col-span-7 bg-white rounded-3xl p-6 border border-slate-200 shadow-xs space-y-4">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <h3 class="text-base font-display font-black text-slate-900 flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">assignment_ind</span>
              <span>Pengajuan Sakit &amp; Tukar Shift Call Taker</span>
            </h3>
            <span class="px-2.5 py-0.5 bg-amber-100 text-amber-900 border border-amber-200 text-[10px] font-bold rounded-full">
              {{ leaveRequests.length }} Pengajuan
            </span>
          </div>

          <div v-if="leaveRequests.length === 0" class="text-center py-8 text-xs text-slate-400 space-y-2">
            <span class="material-symbols-outlined text-3xl text-slate-300">task_alt</span>
            <p>Belum ada pengajuan sakit atau tukar shift baru dari Call Taker.</p>
          </div>

          <div v-else class="space-y-3">
            <div v-for="req in leaveRequests" :key="req.id" class="p-4 bg-slate-50 border border-slate-200 rounded-2xl space-y-2">
              <div class="flex items-center justify-between text-xs">
                <div class="font-bold text-slate-900">{{ req.user_name }} <span class="font-normal text-slate-500">({{ req.category }})</span></div>
                <span 
                  class="px-2 py-0.5 rounded-full text-[10px] font-extrabold uppercase"
                  :class="{
                    'bg-amber-100 text-amber-800 border border-amber-200': req.status === 'PENDING',
                    'bg-emerald-100 text-emerald-800 border border-emerald-200': req.status === 'APPROVED',
                    'bg-rose-100 text-rose-800 border border-rose-200': req.status === 'REJECTED'
                  }"
                >
                  {{ req.status }}
                </span>
              </div>
              <p class="text-xs text-slate-600 font-medium">
                Alasan: "{{ req.reason }}"
              </p>
              <div class="text-[11px] text-slate-500 flex flex-wrap items-center justify-between gap-2 pt-1 border-t border-slate-200">
                <span>Tanggal: <strong>{{ req.shift_date }}</strong> &bull; Pengganti: <strong>{{ req.replacement_name || '-' }}</strong></span>
                
                <div v-if="req.status === 'PENDING'" class="flex items-center gap-1.5">
                  <button @click="updateLeaveStatus(req.id, 'APPROVED')" class="px-2.5 py-1 bg-emerald-700 hover:bg-emerald-800 text-white rounded-lg text-[10px] font-bold border-0 cursor-pointer">Setujui</button>
                  <button @click="updateLeaveStatus(req.id, 'REJECTED')" class="px-2.5 py-1 bg-slate-200 hover:bg-slate-300 text-slate-800 rounded-lg text-[10px] font-bold border-0 cursor-pointer">Tolak</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Posko 112 Information Card (5 cols) -->
        <div class="lg:col-span-5 bg-white rounded-3xl p-6 border border-slate-200 shadow-xs space-y-4">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <h3 class="text-base font-display font-black text-slate-900 flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">emergency</span>
              <span>Informasi Posko Siaga 112</span>
            </h3>
            <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-extrabold rounded-full">TERHUBUNG</span>
          </div>

          <div class="space-y-3 text-xs">
            <div class="p-3.5 bg-slate-50 border border-slate-200 rounded-2xl space-y-1">
              <div class="text-[10px] font-extrabold text-slate-400 uppercase tracking-widest">Nama Posko Siaga</div>
              <div class="font-bold text-slate-900 text-sm">{{ poskoInfo.name }}</div>
            </div>

            <div class="p-3.5 bg-slate-50 border border-slate-200 rounded-2xl space-y-1">
              <div class="text-[10px] font-extrabold text-slate-400 uppercase tracking-widest">Alamat Posko</div>
              <div class="font-medium text-slate-700">{{ poskoInfo.address }}</div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div class="p-3 bg-slate-50 border border-slate-200 rounded-xl space-y-1">
                <div class="text-[10px] font-bold text-slate-400 uppercase">Latitude</div>
                <div class="font-mono font-bold text-rose-700 text-xs">{{ poskoInfo.lat.toFixed(6) }}</div>
              </div>
              <div class="p-3 bg-slate-50 border border-slate-200 rounded-xl space-y-1">
                <div class="text-[10px] font-bold text-slate-400 uppercase">Longitude</div>
                <div class="font-mono font-bold text-rose-700 text-xs">{{ poskoInfo.lng.toFixed(6) }}</div>
              </div>
            </div>
          </div>
        </div>

      </div>

    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import AdminLayout from '@/layouts/AdminLayout.vue';

const poskoInfo = ref({
  name: 'Posko Siaga NTPD 112 Kabupaten Bulukumba',
  address: 'Jl. Jend. Sudirman No. 1, Caile, Kec. Ujung Bulu, Kabupaten Bulukumba, Sulawesi Selatan',
  lat: -5.5645,
  lng: 120.1945,
  radius: 2.0
});

const leaveRequests = ref([]);

const fetchPoskoInfo = async () => {
  try {
    const res = await axios.get('/api/admin/location');
    if (res.data && res.data.success) {
      if (res.data.name) poskoInfo.value.name = res.data.name;
      if (res.data.address) poskoInfo.value.address = res.data.address;
      if (res.data.latitude) poskoInfo.value.lat = res.data.latitude;
      if (res.data.longitude) poskoInfo.value.lng = res.data.longitude;
      if (res.data.radius_meters) poskoInfo.value.radius = res.data.radius_meters;
    }
  } catch (e) {}
};

const fetchLeaveRequests = async () => {
  try {
    const res = await axios.get('/api/presensi/leave-requests');
    if (res.data && Array.isArray(res.data.requests)) {
      leaveRequests.value = res.data.requests;
    }
  } catch (e) {}
};

const updateLeaveStatus = async (id, newStatus) => {
  try {
    const item = leaveRequests.value.find(r => r.id === id);
    if (item) {
      item.status = newStatus;
      await axios.put('/api/presensi/leave-requests', leaveRequests.value);
    }
  } catch (e) {}
};

onMounted(() => {
  fetchPoskoInfo();
  fetchLeaveRequests();
});
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
