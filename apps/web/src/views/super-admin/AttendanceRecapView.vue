<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">
      
      <!-- Toast Notification -->
      <transition name="fade">
        <div v-if="toast.show" class="flex items-center gap-2.5 p-3.5 rounded-xl text-xs font-semibold border w-full shadow-xs bg-[#E8F5E9] border-[#A5D6A7] text-[#1B5E20]">
          <span class="material-symbols-outlined text-lg shrink-0">check_circle</span>
          <span>{{ toast.message }}</span>
        </div>
      </transition>

      <!-- Header Section -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">history_edu</span>
            Rekap Kehadiran
          </h1>
          <p class="text-sm text-[#574146] mt-1">Tinjauan komprehensif kehadiran intern untuk periode yang dipilih.</p>
        </div>

        <!-- Excel & PDF Export/Import Action Buttons -->
        <div class="flex items-center gap-2.5 flex-wrap shrink-0">
          <!-- File Import Input -->
          <label class="bg-white hover:bg-[#FCE4EC] text-[#ab2c5d] border border-[#F8BBD0] px-4 py-2.5 rounded-lg font-bold text-xs transition-all cursor-pointer shadow-xs flex items-center gap-1.5">
            <span class="material-symbols-outlined text-base">upload_file</span>
            <span>Impor Excel / PDF</span>
            <input type="file" @change="handleImportFile" accept=".xlsx,.xls,.pdf" class="hidden" />
          </label>

          <!-- Export Excel Button -->
          <button 
            @click="exportExcel"
            class="bg-[#2E7D32] hover:bg-[#1B5E20] text-white px-4 py-2.5 rounded-lg font-bold text-xs transition-all border-0 cursor-pointer shadow-xs flex items-center gap-1.5"
          >
            <span class="material-symbols-outlined text-base">table_view</span>
            <span>Ekspor EXCEL</span>
          </button>

          <!-- Export PDF Button -->
          <button 
            @click="exportPdf"
            class="bg-[#F06292] hover:bg-[#ab2c5d] text-white px-4 py-2.5 rounded-lg font-bold text-xs transition-all border-0 cursor-pointer shadow-xs flex items-center gap-1.5"
          >
            <span class="material-symbols-outlined text-base">picture_as_pdf</span>
            <span>Ekspor PDF</span>
          </button>
        </div>
      </div>

      <!-- Stats Overview (Bento Grid Style) -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!-- Stat Card 1: Total Hadir -->
        <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <span class="material-symbols-outlined text-[#ab2c5d] bg-[#ffd9e4] p-2 rounded-lg text-xl">how_to_reg</span>
          </div>
          <div>
            <p class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-1">Total Hadir</p>
            <p class="text-2xl font-extrabold text-[#1b1c1c] font-mono">{{ stats.totalHadir }}</p>
          </div>
        </div>

        <!-- Stat Card 2: Persentase Tepat Waktu -->
        <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <span class="material-symbols-outlined text-[#ab2c5d] bg-[#ffd9e4] p-2 rounded-lg text-xl">schedule</span>
          </div>
          <div>
            <p class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-1">Persentase Tepat Waktu</p>
            <p class="text-2xl font-extrabold text-[#1b1c1c] font-mono">{{ stats.pctTepatWaktu }}%</p>
          </div>
        </div>

        <!-- Stat Card 3: Total Terlambat -->
        <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <span class="material-symbols-outlined text-[#F57F17] bg-[#FFF8E1] p-2 rounded-lg text-xl">history_toggle_off</span>
          </div>
          <div>
            <p class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-1">Total Terlambat</p>
            <p class="text-2xl font-extrabold text-[#1b1c1c] font-mono">{{ stats.totalTerlambat }}</p>
          </div>
        </div>

        <!-- Stat Card 4: Total Absen -->
        <div class="bg-white/85 backdrop-blur-md rounded-xl p-6 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col justify-between">
          <div class="flex justify-between items-start mb-4">
            <span class="material-symbols-outlined text-rose-600 bg-rose-50 p-2 rounded-lg text-xl">person_off</span>
          </div>
          <div>
            <p class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-1">Total Absen</p>
            <p class="text-2xl font-extrabold text-[#1b1c1c] font-mono">{{ stats.totalAbsen }}</p>
          </div>
        </div>
      </div>

      <!-- Filters Bar (Glass Card Layout) -->
      <div class="bg-white/85 backdrop-blur-md rounded-xl p-4 border border-[#F8BBD0] shadow-[0px_10px_30px_rgba(240,98,146,0.05)] flex flex-col md:flex-row gap-4 items-center justify-between">
        <div class="flex flex-col md:flex-row gap-3 w-full md:w-auto">
          <!-- Month Selector -->
          <div class="relative w-full md:w-48">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">calendar_month</span>
            <input 
              v-model="filterMonth" 
              type="month"
              class="w-full pl-10 pr-3 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all cursor-pointer"
            />
          </div>

          <!-- Department Selector -->
          <div class="relative w-full md:w-60">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">domain</span>
            <select 
              v-model="filterDept" 
              class="w-full pl-10 pr-8 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-bold text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 appearance-none transition-all cursor-pointer"
            >
              <option value="">Semua Departemen</option>
              <option value="Engineering">Engineering</option>
              <option value="Design">Design</option>
              <option value="Marketing">Marketing</option>
              <option value="Diskominfo">Diskominfo Bulukumba</option>
            </select>
          </div>
        </div>

        <!-- Search Input -->
        <div class="relative w-full md:w-64">
          <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-sm">search</span>
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Cari nama intern..." 
            class="w-full pl-10 pr-4 py-2 bg-white border border-[#F8BBD0] rounded-lg text-xs font-medium text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/30 transition-all"
          />
        </div>
      </div>

      <!-- Data Table Card -->
      <div class="bg-white/90 backdrop-blur-md rounded-xl border border-[#F8BBD0] overflow-hidden shadow-[0px_10px_30px_rgba(240,98,146,0.05)]">
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-[#FCE4EC] border-b border-[#F8BBD0]">
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Nama Intern</th>
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Tanggal</th>
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Jam Masuk</th>
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Jam Keluar</th>
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Total Jam</th>
                <th class="py-3.5 px-6 text-[11px] font-bold text-[#574146] uppercase tracking-wider">Status</th>
              </tr>
            </thead>
            <tbody class="text-xs text-[#1b1c1c] bg-white divide-y divide-[#F8BBD0]">
              <tr 
                v-for="row in filteredRecords" 
                :key="row.id" 
                class="border-b border-[#F8BBD0] hover:bg-[#FCE4EC]/30 transition-colors"
              >
                <td class="py-4 px-6 font-bold text-[#1b1c1c]">{{ row.name }}</td>
                <td class="py-4 px-6 text-[#574146]">{{ row.date }}</td>
                <td class="py-4 px-6 font-mono" :class="row.clockIn === '--:--' ? 'text-slate-400' : 'text-[#1b1c1c]'">{{ row.clockIn }}</td>
                <td class="py-4 px-6 font-mono" :class="row.clockOut === '--:--' ? 'text-slate-400' : 'text-[#1b1c1c]'">{{ row.clockOut }}</td>
                <td class="py-4 px-6 font-mono text-[#574146]">{{ row.totalHours }}</td>
                <td class="py-4 px-6">
                  <span 
                    v-if="row.status === 'HADIR'" 
                    class="px-3 py-1 rounded-full bg-[#E8F5E9] text-[#1B5E20] text-[10px] uppercase font-bold border border-[#A5D6A7]"
                  >
                    Hadir
                  </span>
                  <span 
                    v-else-if="row.status === 'TERLAMBAT'" 
                    class="px-3 py-1 rounded-full bg-[#FFF8E1] text-[#F57F17] text-[10px] uppercase font-bold border border-[#FFE082]"
                  >
                    Terlambat
                  </span>
                  <span 
                    v-else-if="row.status === 'ABSEN'" 
                    class="px-3 py-1 rounded-full bg-[#FCE4EC] text-[#F06292] text-[10px] uppercase font-bold border border-[#F8BBD0]"
                  >
                    Absen
                  </span>
                  <span 
                    v-else 
                    class="px-3 py-1 rounded-full bg-slate-100 text-slate-600 text-[10px] uppercase font-bold border border-slate-200"
                  >
                    {{ row.status }}
                  </span>
                </td>
              </tr>

              <!-- Empty State -->
              <tr v-if="filteredRecords.length === 0">
                <td colspan="6" class="py-12 text-center text-[#8a7176]">
                  <span class="material-symbols-outlined text-4xl block mb-2 opacity-50">event_busy</span>
                  <span class="text-xs font-semibold">Belum ada data rekapitulasi presensi.</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Container -->
        <div class="flex items-center justify-between px-6 py-4 bg-white border-t border-[#F8BBD0]">
          <div class="text-xs text-[#574146] font-medium">
            Menampilkan {{ filteredRecords.length > 0 ? '1' : '0' }}-{{ filteredRecords.length }} dari {{ records.length }} data
          </div>
          <div class="flex gap-2">
            <button disabled class="p-1 rounded border border-[#F8BBD0] text-[#574146] opacity-50 cursor-not-allowed bg-white">
              <span class="material-symbols-outlined text-sm">chevron_left</span>
            </button>
            <button class="w-8 h-8 rounded bg-[#F06292] text-white font-bold text-xs flex items-center justify-center border-0">1</button>
            <button disabled class="p-1 rounded border border-[#F8BBD0] text-[#574146] opacity-50 cursor-not-allowed bg-white">
              <span class="material-symbols-outlined text-sm">chevron_right</span>
            </button>
          </div>
        </div>
      </div>

    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, computed } from 'vue';
import AdminLayout from '@/layouts/AdminLayout.vue';

const filterMonth = ref('2026-08');
const filterDept = ref('');
const searchQuery = ref('');
const toast = ref({ show: false, message: '' });

const records = ref([]);

const showToast = (msg) => {
  toast.value = { show: true, message: msg };
  setTimeout(() => { toast.value.show = false; }, 3500);
};

const stats = computed(() => {
  if (records.value.length === 0) {
    return {
      totalHadir: 0,
      pctTepatWaktu: 0,
      totalTerlambat: 0,
      totalAbsen: 0
    };
  }
  const hadir = records.value.filter(r => r.status === 'HADIR').length;
  const terlambat = records.value.filter(r => r.status === 'TERLAMBAT').length;
  const absen = records.value.filter(r => r.status === 'ABSEN').length;
  const pct = Math.round((hadir / records.value.length) * 100);
  return {
    totalHadir: hadir,
    pctTepatWaktu: pct,
    totalTerlambat: terlambat,
    totalAbsen: absen
  };
});

const filteredRecords = computed(() => {
  return records.value.filter(r => {
    const matchesSearch = !searchQuery.value || r.name.toLowerCase().includes(searchQuery.value.toLowerCase());
    const matchesDept = !filterDept.value || r.dept === filterDept.value;
    return matchesSearch && matchesDept;
  });
});

function exportExcel() {
  showToast('Mengunduh Laporan Rekap Kehadiran format EXCEL (.xlsx)...');
}

function exportPdf() {
  showToast('Mengunduh Laporan Rekap Kehadiran format PDF (.pdf)...');
}

function handleImportFile(event) {
  const file = event.target.files[0];
  if (file) {
    showToast(`File "${file.name}" berhasil diimpor ke data Rekap Kehadiran!`);
  }
}
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
.material-symbols-outlined.fill { font-variation-settings: 'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
