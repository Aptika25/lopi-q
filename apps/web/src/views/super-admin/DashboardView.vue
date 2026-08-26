<template>
  <AdminLayout>
    <div class="space-y-6 select-none font-sans text-slate-800">
      
      <!-- Dashboard Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/70 pb-4">
        <div>
          <h1 class="text-2xl sm:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill">shield_person</span>
            Welcome back, Admin.
          </h1>
          <p class="text-sm text-[#574146] mt-1">Here's what's happening today in the LOPI-Q system (Posko Siaga 112).</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#fec1d6]/30 text-[#ab2c5d] rounded-full text-xs font-bold border border-[#fec1d6]">
            <span class="w-2 h-2 rounded-full bg-[#ab2c5d] animate-ping"></span>
            <span>SYSTEM OPTIMAL</span>
          </span>
        </div>
      </div>

      <!-- Stats Bento Grid (4 Cards) -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
        <!-- Card 1: Interns / Petugas -->
        <div class="bg-gradient-to-br from-white to-[#fec1d6]/10 p-6 rounded-xl border border-[#ddbfc5] hover:shadow-lg hover:shadow-[#f06292]/10 transition-all flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <div class="p-3 bg-[#f06292]/10 rounded-lg text-[#ab2c5d]">
              <span class="material-symbols-outlined fill text-2xl">school</span>
            </div>
            <span class="text-[11px] font-bold uppercase text-[#574146] bg-slate-100 px-2.5 py-1 rounded-full">Active</span>
          </div>
          <div>
            <h3 class="text-3xl font-extrabold text-[#1b1c1c] mb-1">10 Petugas</h3>
            <p class="text-xs text-[#574146]">Total Peserta Magang Terdaftar</p>
          </div>
        </div>

        <!-- Card 2: Visitors / Presensi -->
        <div class="bg-gradient-to-br from-white to-[#fec1d6]/10 p-6 rounded-xl border border-[#ddbfc5] hover:shadow-lg hover:shadow-[#f06292]/10 transition-all flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <div class="p-3 bg-[#fec1d6]/30 rounded-lg text-[#805062]">
              <span class="material-symbols-outlined fill text-2xl">group</span>
            </div>
            <span class="text-[11px] font-bold uppercase text-[#ab2c5d] bg-[#f06292]/20 px-2.5 py-1 rounded-full">100% Valid</span>
          </div>
          <div>
            <h3 class="text-3xl font-extrabold text-[#1b1c1c] mb-1">Presensi Hari Ini</h3>
            <p class="text-xs text-[#574146]">Verifikasi Geofence &amp; Kamera</p>
          </div>
        </div>

        <!-- Card 3: Geofence Posko -->
        <div class="bg-gradient-to-br from-white to-[#fec1d6]/10 p-6 rounded-xl border border-[#ddbfc5] hover:shadow-lg hover:shadow-[#f06292]/10 transition-all flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <div class="p-3 bg-[#f06292]/10 rounded-lg text-[#ab2c5d]">
              <span class="material-symbols-outlined fill text-2xl">radar</span>
            </div>
            <span class="text-[11px] font-bold uppercase text-emerald-700 bg-emerald-100 px-2.5 py-1 rounded-full">Secure</span>
          </div>
          <div>
            <h3 class="text-xl font-bold text-[#1b1c1c] mb-1">Radius {{ poskoInfo.radius }}m</h3>
            <p class="text-xs font-mono text-[#574146] truncate">{{ poskoInfo.lat.toFixed(4) }}, {{ poskoInfo.lng.toFixed(4) }}</p>
          </div>
        </div>

        <!-- Card 4: Logbook / Leave Requests -->
        <div class="bg-gradient-to-br from-white to-[#fec1d6]/10 p-6 rounded-xl border border-[#ddbfc5] hover:shadow-lg hover:shadow-[#f06292]/10 transition-all flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <div class="p-3 bg-[#fec1d6]/30 rounded-lg text-[#805062]">
              <span class="material-symbols-outlined fill text-2xl">library_books</span>
            </div>
            <span class="text-[11px] font-bold uppercase text-[#574146] bg-slate-100 px-2.5 py-1 rounded-full">Pending</span>
          </div>
          <div>
            <h3 class="text-3xl font-extrabold text-[#1b1c1c] mb-1">{{ leaveRequests.length }}</h3>
            <p class="text-xs text-[#574146]">Pengajuan Sakit &amp; Tukar Shift</p>
          </div>
        </div>
      </div>

      <!-- Quick Shortcuts Grid -->
      <div class="bg-white rounded-xl border border-[#ddbfc5] p-6 shadow-xs space-y-4">
        <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
          <span class="material-symbols-outlined text-[#ab2c5d]">grid_view</span>
          <span>Akses Cepat Pengelolaan Sistem</span>
        </h3>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
          <router-link to="/admin/interns" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">badge</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Peserta Magang</div>
          </router-link>

          <router-link to="/admin/admins" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">manage_accounts</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Admin Management</div>
          </router-link>

          <router-link to="/admin/attendance-recap" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">assessment</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Rekap Presensi</div>
          </router-link>

          <router-link to="/admin/location" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">pin_drop</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Lokasi Posko</div>
          </router-link>

          <router-link to="/admin/schedules" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">calendar_month</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Jadwal Shift</div>
          </router-link>

          <router-link to="/admin/security" class="p-4 bg-slate-50 hover:bg-[#fec1d6]/20 border border-[#ddbfc5] hover:border-[#ab2c5d] rounded-lg text-center space-y-2 transition-all group decoration-none">
            <div class="w-10 h-10 bg-white text-[#ab2c5d] rounded-lg flex items-center justify-center mx-auto shadow-xs group-hover:scale-105 transition-transform">
              <span class="material-symbols-outlined text-[22px]">security</span>
            </div>
            <div class="text-xs font-bold text-[#1b1c1c] group-hover:text-[#ab2c5d]">Keamanan 2FA</div>
          </router-link>
        </div>
      </div>

      <!-- Content Split Area -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
        
        <!-- Recent Activity Table / Leave Requests (lg:col-span-2) -->
        <div class="lg:col-span-2 bg-white rounded-xl border border-[#ddbfc5] overflow-hidden flex flex-col shadow-xs">
          <div class="p-6 border-b border-[#ddbfc5]/50 flex justify-between items-center bg-gradient-to-r from-white to-[#fec1d6]/10">
            <div>
              <h3 class="font-bold text-lg text-[#1b1c1c]">Pengajuan Sakit &amp; Tukar Shift</h3>
              <p class="text-xs text-[#574146] mt-0.5">Daftar pengajuan izin dan pergantian piket petugas magang</p>
            </div>
            <span class="px-3 py-1 bg-[#fec1d6]/30 text-[#ab2c5d] text-xs font-bold rounded-full border border-[#fec1d6]">
              {{ leaveRequests.length }} Total
            </span>
          </div>

          <div v-if="leaveRequests.length === 0" class="text-center py-12 text-sm text-[#574146] space-y-2">
            <span class="material-symbols-outlined text-4xl text-[#8a7176]">task_alt</span>
            <p>Belum ada pengajuan sakit atau tukar shift baru dari Peserta Magang.</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-[#f4dce4]/40 border-b border-[#ddbfc5]/50 text-xs uppercase text-[#574146] font-semibold">
                  <th class="py-3.5 px-5">Peserta / Modul</th>
                  <th class="py-3.5 px-5">Alasan</th>
                  <th class="py-3.5 px-5">Tanggal Shift</th>
                  <th class="py-3.5 px-5">Status / Aksi</th>
                </tr>
              </thead>
              <tbody class="text-xs divide-y divide-[#ddbfc5]/30">
                <tr v-for="req in leaveRequests" :key="req.id" class="hover:bg-[#fec1d6]/10 transition-colors">
                  <td class="py-4 px-5">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 rounded-full bg-[#e4e2e2] flex items-center justify-center text-[#574146]">
                        <span class="material-symbols-outlined text-sm">person</span>
                      </div>
                      <div>
                        <p class="font-bold text-[#1b1c1c]">{{ req.user_name }}</p>
                        <p class="text-[11px] text-[#574146]">{{ req.category }}</p>
                      </div>
                    </div>
                  </td>
                  <td class="py-4 px-5 text-[#1b1c1c] max-w-xs truncate">
                    "{{ req.reason }}"
                  </td>
                  <td class="py-4 px-5 text-[#574146]">
                    <div>{{ req.shift_date }}</div>
                    <div class="text-[10px] text-[#8a7176]">Pengganti: {{ req.replacement_name || '-' }}</div>
                  </td>
                  <td class="py-4 px-5">
                    <div class="flex items-center gap-2">
                      <span 
                        class="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase"
                        :class="{
                          'bg-[#FFF8E1] text-[#F57F17]': req.status === 'PENDING',
                          'bg-[#E8F5E9] text-[#1B5E20]': req.status === 'APPROVED',
                          'bg-[#FFEBEE] text-[#C62828]': req.status === 'REJECTED'
                        }"
                      >
                        {{ req.status }}
                      </span>
                      <div v-if="req.status === 'PENDING'" class="flex items-center gap-1">
                        <button @click="updateLeaveStatus(req.id, 'APPROVED')" class="px-2 py-0.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded text-[10px] font-bold border-0 cursor-pointer">Setujui</button>
                        <button @click="updateLeaveStatus(req.id, 'REJECTED')" class="px-2 py-0.5 bg-rose-600 hover:bg-rose-700 text-white rounded text-[10px] font-bold border-0 cursor-pointer">Tolak</button>
                      </div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Quick Actions & Posko Info (lg:col-span-1) -->
        <div class="flex flex-col gap-6">
          <!-- Quick Actions Card -->
          <div class="bg-white rounded-xl border border-[#ddbfc5] p-6 flex flex-col shadow-xs">
            <h3 class="font-bold text-base text-[#1b1c1c] mb-4">Aksi Cepat Admin</h3>
            <div class="space-y-3">
              <router-link to="/admin/interns" class="w-full bg-[#f06292] text-white py-2.5 px-4 rounded-lg font-semibold text-xs flex items-center justify-center gap-2 hover:bg-[#ab2c5d] transition-colors shadow-xs decoration-none">
                <span class="material-symbols-outlined text-sm">person_add</span>
                Kelola Peserta Magang
              </router-link>
              <router-link to="/admin/schedules" class="w-full bg-[#fec1d6]/30 text-[#ab2c5d] py-2.5 px-4 rounded-lg font-semibold text-xs flex items-center justify-center gap-2 hover:bg-[#fec1d6]/60 transition-colors decoration-none">
                <span class="material-symbols-outlined text-sm">calendar_month</span>
                Jadwal Shift &amp; Presensi
              </router-link>
              <router-link to="/admin/attendance-recap" class="w-full border border-[#ddbfc5] text-[#574146] py-2.5 px-4 rounded-lg font-semibold text-xs flex items-center justify-center gap-2 hover:bg-slate-50 transition-colors decoration-none">
                <span class="material-symbols-outlined text-sm">download</span>
                Export Rekap Laporan
              </router-link>
            </div>
          </div>

          <!-- Posko Info Banner Card -->
          <div class="bg-white rounded-xl border border-[#ddbfc5] p-5 space-y-3 shadow-xs">
            <div class="flex items-center justify-between border-b border-[#ddbfc5]/50 pb-2.5">
              <h4 class="font-bold text-xs uppercase text-[#574146] flex items-center gap-1.5">
                <span class="material-symbols-outlined text-[#ab2c5d] text-base">emergency</span>
                Informasi Posko 112
              </h4>
              <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full">TERHUBUNG</span>
            </div>
            <div class="text-xs space-y-2">
              <div>
                <span class="text-[10px] uppercase font-bold text-[#8a7176]">Nama Posko</span>
                <p class="font-bold text-[#1b1c1c] mt-0.5">{{ poskoInfo.name }}</p>
              </div>
              <div>
                <span class="text-[10px] uppercase font-bold text-[#8a7176]">Alamat</span>
                <p class="text-[#574146] leading-relaxed mt-0.5">{{ poskoInfo.address }}</p>
              </div>
            </div>
          </div>

        </div>

      </div>

      <!-- Microservices & System Infrastructure Status Section -->
      <div class="bg-white rounded-xl border border-[#ddbfc5] p-6 shadow-xs space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-[#ddbfc5]/50 pb-3">
          <div>
            <h3 class="text-base font-bold text-[#1b1c1c] flex items-center gap-2">
              <span class="material-symbols-outlined text-[#ab2c5d]">dns</span>
              <span>Status Layanan Microservices LOPI-Q</span>
            </h3>
            <p class="text-xs text-[#574146] mt-0.5">Monitoring kesehatan container backend, API gateway, &amp; database PostgreSQL</p>
          </div>
          <span class="px-3 py-1 bg-emerald-100 text-emerald-800 text-xs font-bold rounded-full flex items-center gap-1.5 w-fit">
            <span class="w-2 h-2 rounded-full bg-emerald-600 animate-pulse"></span>
            <span>Semua Layanan Berjalan</span>
          </span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Service 1: API Gateway -->
          <div class="p-4 bg-slate-50 border border-[#ddbfc5] rounded-xl flex items-center justify-between hover:border-[#ab2c5d] transition-all">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-[#f06292]/10 rounded-lg text-[#ab2c5d]">
                <span class="material-symbols-outlined text-xl">api</span>
              </div>
              <div>
                <h4 class="font-bold text-xs text-[#1b1c1c]">API Gateway</h4>
                <p class="text-[11px] text-[#574146] font-mono">Port: 5000 &bull; REST</p>
              </div>
            </div>
            <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full">ACTIVE</span>
          </div>

          <!-- Service 2: Auth Service -->
          <div class="p-4 bg-slate-50 border border-[#ddbfc5] rounded-xl flex items-center justify-between hover:border-[#ab2c5d] transition-all">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-[#f06292]/10 rounded-lg text-[#ab2c5d]">
                <span class="material-symbols-outlined text-xl">vpn_key</span>
              </div>
              <div>
                <h4 class="font-bold text-xs text-[#1b1c1c]">Auth Service</h4>
                <p class="text-[11px] text-[#574146] font-mono">Port: 50051 &bull; gRPC</p>
              </div>
            </div>
            <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full">ACTIVE</span>
          </div>

          <!-- Service 3: User Service -->
          <div class="p-4 bg-slate-50 border border-[#ddbfc5] rounded-xl flex items-center justify-between hover:border-[#ab2c5d] transition-all">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-[#f06292]/10 rounded-lg text-[#ab2c5d]">
                <span class="material-symbols-outlined text-xl">group_work</span>
              </div>
              <div>
                <h4 class="font-bold text-xs text-[#1b1c1c]">User Service</h4>
                <p class="text-[11px] text-[#574146] font-mono">Port: 50052 &bull; gRPC</p>
              </div>
            </div>
            <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full">ACTIVE</span>
          </div>

          <!-- Service 4: PostgreSQL Database -->
          <div class="p-4 bg-slate-50 border border-[#ddbfc5] rounded-xl flex items-center justify-between hover:border-[#ab2c5d] transition-all">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-[#fec1d6]/30 rounded-lg text-[#805062]">
                <span class="material-symbols-outlined text-xl">database</span>
              </div>
              <div>
                <h4 class="font-bold text-xs text-[#1b1c1c]">PostgreSQL DB</h4>
                <p class="text-[11px] text-[#574146] font-mono">Port: 5432 &bull; Apps</p>
              </div>
            </div>
            <span class="px-2 py-0.5 bg-emerald-100 text-emerald-800 text-[10px] font-bold rounded-full">ONLINE</span>
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
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
