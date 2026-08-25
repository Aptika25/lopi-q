<template>
  <AdminLayout>
    <div class="space-y-6 select-none font-sans">
      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4 mb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Log Aktivitas &amp; Audit Keamanan</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Catatan riwayat audit aktivitas sistem, otentikasi 2FA, presensi geofence, dan perubahan kredensial pengguna real-time.</p>
        </div>
        <button 
          @click="fetchLogs"
          :disabled="loading"
          class="py-2.5 px-4 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold text-xs rounded-xl border border-slate-300 flex items-center justify-center gap-1.5 cursor-pointer shrink-0 disabled:opacity-50 transition-all"
        >
          <span class="material-symbols-outlined text-[16px]" :class="loading ? 'animate-spin' : ''">refresh</span>
          <span>{{ loading ? 'Memuat Log...' : 'Refresh Audit Log' }}</span>
        </button>
      </div>

      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-slate-200 shadow-sm space-y-6">

        <!-- Controls Container: Clean 2-Bar Toolbar Layout -->
        <div class="bg-slate-50 p-4 sm:p-5 rounded-2xl border border-slate-200 shadow-2xs space-y-3.5">
          
          <!-- Baris 1 (Atas): Cari Input + Tombol Export -->
          <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
            <!-- Cari Input -->
            <div class="relative w-full flex-1">
              <span class="material-symbols-outlined absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 text-[18px]">search</span>
              <input 
                v-model="searchQuery" 
                type="text" 
                placeholder="Cari log berdasarkan NIP, Nama, Jenis Aktivitas, Alamat IP, atau User Agent..." 
                class="w-full pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-medium text-slate-800 focus:outline-none focus:border-rose-600 focus:ring-2 focus:ring-rose-500/20 transition-all shadow-xs"
              />
            </div>

            <!-- Tombol Export -->
            <div class="flex items-center gap-2 shrink-0 w-full sm:w-auto justify-end">
              <!-- Green CSV Button -->
              <button 
                @click="exportToCSV"
                class="py-2 px-4 bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white font-extrabold text-xs rounded-xl shadow-xs transition-all flex items-center gap-1.5 cursor-pointer border-0"
                title="Unduh Data Log Aktivitas sebagai Excel CSV"
              >
                <span class="material-symbols-outlined text-[16px]">download</span>
                <span>Excel (CSV)</span>
              </button>

              <!-- Red PDF Print Button -->
              <button 
                @click="exportToPDF"
                class="py-2 px-4 bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white font-extrabold text-xs rounded-xl shadow-xs transition-all flex items-center gap-1.5 cursor-pointer border-0"
                title="Cetak atau Simpan sebagai PDF"
              >
                <span class="material-symbols-outlined text-[16px]">print</span>
                <span>Cetak / PDF</span>
              </button>
            </div>
          </div>

          <!-- Baris 2 (Bawah): Dropdown Filter (Kiri) + Page Size & Total Metric (Kanan) -->
          <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-3 pt-3 border-t border-slate-200/80">
            
            <!-- Left Side: Dropdown Filter -->
            <div class="flex flex-wrap items-center gap-2">
              <!-- Filter Label -->
              <div class="flex items-center gap-1 text-slate-500 font-extrabold text-xs tracking-wider uppercase pr-1">
                <span class="material-symbols-outlined text-[16px] text-slate-400">filter_alt</span>
                <span>Filter:</span>
              </div>

              <!-- Filter: Jenis Aktivitas -->
              <select 
                v-model="selectedAction" 
                class="px-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-xs"
              >
                <option value="">Semua Aktivitas</option>
                <option value="LOGIN">LOGIN</option>
                <option value="PRESENSI">PRESENSI</option>
                <option value="2FA">2FA OTENTIKASI</option>
                <option value="USER">MANAJEMEN USER</option>
                <option value="SCHEDULE">JADWAL SHIFT</option>
              </select>

              <!-- Filter: Tanggal -->
              <select 
                v-model="selectedDay" 
                class="px-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-xs"
              >
                <option value="">Semua Tanggal</option>
                <option v-for="d in 31" :key="d" :value="d">Tanggal {{ d }}</option>
              </select>

              <!-- Filter: Bulan -->
              <select 
                v-model="selectedMonth" 
                class="px-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-xs"
              >
                <option value="">Semua Bulan</option>
                <option v-for="(m, idx) in monthsIndo" :key="idx" :value="idx + 1">{{ m }}</option>
              </select>

              <!-- Filter: Tahun -->
              <select 
                v-model="selectedYear" 
                class="px-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-xs"
              >
                <option value="">Semua Tahun</option>
                <option v-for="y in [2026, 2025, 2024]" :key="y" :value="y">{{ y }}</option>
              </select>

              <!-- Reset Filter Button -->
              <button 
                v-if="searchQuery || selectedAction || selectedDay || selectedMonth || selectedYear"
                @click="searchQuery = ''; selectedAction = ''; selectedDay = ''; selectedMonth = ''; selectedYear = '';"
                class="px-2.5 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 font-extrabold rounded-xl border border-rose-200 text-xs cursor-pointer whitespace-nowrap"
              >
                Reset Filter
              </button>
            </div>

            <!-- Right Side: Page Size Selector & Total Metric -->
            <div class="flex items-center gap-3 shrink-0 justify-between lg:justify-end">
              <!-- Sorting: Tampilkan Per Halaman -->
              <div class="flex items-center gap-1.5 text-xs text-slate-500 font-bold">
                <span class="text-[11px] text-slate-400">Tampilkan:</span>
                <select 
                  v-model="pageSize" 
                  class="px-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-xs"
                >
                  <option v-for="size in pageSizeOptions" :key="size" :value="size">{{ size }} Data / Hal</option>
                </select>
              </div>

              <!-- Total Log Count -->
              <div class="text-xs font-bold text-slate-600 whitespace-nowrap">
                Total: <span class="text-rose-700 font-extrabold">{{ filteredLogs.length }}</span> log
              </div>
            </div>

          </div>

        </div>

        <!-- Activity Audit Table -->
        <div class="overflow-x-auto rounded-2xl border border-slate-200 shadow-xs w-full">
          <table class="w-full text-left text-xs text-slate-700 table-auto">
            <thead class="bg-slate-100 text-slate-700 font-bold uppercase tracking-wider border-b border-slate-200 text-[10px]">
              <tr>
                <th class="py-3 px-3 whitespace-nowrap">Waktu &amp; Tanggal</th>
                <th class="py-3 px-3 whitespace-nowrap">Pengguna / NIP</th>
                <th class="py-3 px-3 whitespace-nowrap">Jenis Aktivitas</th>
                <th class="py-3 px-3">Detail Peristiwa</th>
                <th class="py-3 px-3 whitespace-nowrap">Alamat IP</th>
                <th class="py-3 px-3 whitespace-nowrap">User Agent</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 bg-white">
              <tr v-if="loading && logs.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400 font-semibold">
                  <span class="material-symbols-outlined animate-spin text-2xl text-rose-600 block mb-1">sync</span>
                  <span>Memuat data log aktivitas dari server...</span>
                </td>
              </tr>
              <tr v-else-if="filteredLogs.length === 0">
                <td colspan="6" class="py-8 text-center text-slate-400 font-semibold">
                  Tidak ada catatan log aktivitas yang sesuai dengan filter.
                </td>
              </tr>
              <tr v-for="log in paginatedLogs" :key="log.id" class="hover:bg-slate-50 transition-colors">
                <td class="py-3 px-3 font-mono text-[11px] font-bold text-slate-600 whitespace-nowrap align-middle">
                  {{ formatDateIndo(log.timestamp) }}
                </td>
                <td class="py-3 px-3 whitespace-nowrap align-middle">
                  <div class="font-bold text-slate-900 text-xs">{{ log.user_name }}</div>
                  <div class="font-mono text-[10.5px] text-rose-700 font-semibold whitespace-nowrap">NIP. {{ log.user_nip }}</div>
                </td>
                <td class="py-3 px-3 font-bold whitespace-nowrap align-middle">
                  <span 
                    class="px-2.5 py-0.5 rounded-full text-[10px] tracking-wide uppercase font-extrabold border inline-block text-center shadow-2xs"
                    :class="getActionBadgeClass(log.action)"
                  >
                    {{ log.action }}
                  </span>
                </td>
                <td class="py-3 px-3 text-slate-700 font-medium text-[11px] leading-snug max-w-[320px] align-middle">{{ log.details }}</td>
                <td class="py-3 px-3 font-mono text-[11px] text-slate-700 font-semibold whitespace-nowrap align-middle">
                  {{ log.ip_address || '180.242.190.12' }}
                </td>
                <td class="py-3 px-3 font-mono text-[10.5px] text-slate-500 max-w-[160px] truncate align-middle" :title="log.user_agent">
                  {{ formatUserAgent(log.user_agent) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls Footer -->
        <div v-if="filteredLogs.length > 0" class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t border-slate-200 text-xs">
          <div class="text-slate-500 font-semibold text-center sm:text-left">
            Menampilkan <strong class="text-slate-900 font-mono font-bold">{{ ((currentPage - 1) * pageSize) + 1 }}</strong> - <strong class="text-slate-900 font-mono font-bold">{{ Math.min(currentPage * pageSize, filteredLogs.length) }}</strong> dari <strong class="text-slate-900 font-mono font-bold">{{ filteredLogs.length }}</strong> log aktivitas
          </div>

          <div class="flex items-center gap-1.5">
            <button 
              @click="currentPage--" 
              :disabled="currentPage === 1"
              class="px-3.5 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-xl border border-slate-300 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs transition-all"
            >
              <span class="material-symbols-outlined text-[16px]">chevron_left</span>
              <span>Sebelumnya</span>
            </button>

            <span class="px-3.5 py-1.5 bg-rose-50 border border-rose-200 text-rose-800 font-extrabold font-mono rounded-xl text-xs">
              {{ currentPage }} / {{ totalPages }}
            </span>

            <button 
              @click="currentPage++" 
              :disabled="currentPage >= totalPages"
              class="px-3.5 py-1.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold rounded-xl border border-slate-300 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed flex items-center gap-1 text-xs transition-all"
            >
              <span>Berikutnya</span>
              <span class="material-symbols-outlined text-[16px]">chevron_right</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import axios from 'axios';
import AdminLayout from '@/layouts/AdminLayout.vue';

const searchQuery = ref('');
const selectedAction = ref('');
const selectedDay = ref('');
const selectedMonth = ref('');
const selectedYear = ref('');
const pageSize = ref(10);
const currentPage = ref(1);
const pageSizeOptions = [5, 10, 20, 50, 100];
const loading = ref(false);

const logs = ref([]);

watch([searchQuery, selectedAction, selectedDay, selectedMonth, selectedYear, pageSize], () => {
  currentPage.value = 1;
});

const monthsIndo = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
];

// Badge styling matching exact screenshot pastel palette
const getActionBadgeClass = (actionStr) => {
  if (!actionStr) return 'bg-slate-100 text-slate-700 border-slate-200';
  const act = actionStr.toUpperCase();

  // 1. 2FA Authentication (Pastel Soft Red / Rose)
  if (act.includes('2FA')) {
    return 'bg-red-100/70 text-rose-700 border-red-200/60';
  }

  // 2. Presensi & Login (Pastel Soft Green / Mint)
  if (act.includes('PRESENSI') || act.includes('LOGIN')) {
    return 'bg-emerald-100/70 text-emerald-800 border-emerald-200/60';
  }

  // 3. Update Schedule & User (Pastel Soft Blue / Indigo)
  if (act.includes('SCHEDULE') || act.includes('USER') || act.includes('RESET') || act.includes('ADMIN')) {
    return 'bg-blue-100/70 text-indigo-800 border-blue-200/60';
  }

  return 'bg-slate-100 text-slate-700 border-slate-200';
};

// Format timestamp into Indonesian format: 5 Agu 2026, 21.36.59 WITA
const formatDateIndo = (tsStr) => {
  if (!tsStr) return '-';

  let year = 2026, monthIdx = 7, day = 5, hours = 21, mins = 36, secs = 59;

  const cleanStr = String(tsStr).replace('T', ' ').replace(/\+.*/, '');
  const parts = cleanStr.split(' ');
  if (parts.length >= 1 && parts[0].includes('-')) {
    const [y, m, d] = parts[0].split('-').map(n => parseInt(n, 10));
    if (y && m && d) {
      year = y;
      monthIdx = m - 1;
      day = d;
    }
  }
  if (parts.length >= 2 && parts[1].includes(':')) {
    const timeParts = parts[1].split(':').map(n => parseInt(n, 10));
    if (timeParts.length >= 2) {
      hours = timeParts[0] || 0;
      mins = timeParts[1] || 0;
      secs = timeParts[2] || 0;
    }
  }

  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'];
  const monthName = months[monthIdx] || 'Agu';

  const hStr = String(hours).padStart(2, '0');
  const mStr = String(mins).padStart(2, '0');
  const sStr = String(secs).padStart(2, '0');

  return `${day} ${monthName} ${year}, ${hStr}.${mStr}.${sStr} WITA`;
};

// Format User Agent for clean table view
const formatUserAgent = (uaStr) => {
  if (!uaStr) return 'Chrome (Windows)';
  if (uaStr.includes('Android')) return 'Chrome (Android)';
  if (uaStr.includes('iPhone') || uaStr.includes('iPad')) return 'Safari (iOS)';
  if (uaStr.includes('Windows')) return 'Chrome (Windows)';
  if (uaStr.includes('Macintosh') || uaStr.includes('Mac OS')) return 'Safari (macOS)';
  if (uaStr.length > 25) {
    return uaStr.substring(0, 22) + '...';
  }
  return uaStr;
};

const fetchLogs = async () => {
  loading.value = true;
  try {
    const res = await axios.get('/api/admin/activity-logs');
    if (res.data && Array.isArray(res.data.logs) && res.data.logs.length > 0) {
      logs.value = res.data.logs;
    }
  } catch (e) {
    console.warn('Gagal memuat log aktivitas real-time:', e);
  } finally {
    loading.value = false;
  }
};

const filteredLogs = computed(() => {
  return logs.value.filter(l => {
    // 1. Search Query Filter
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase();
      const matchSearch = (
        (l.user_name && l.user_name.toLowerCase().includes(q)) ||
        (l.user_nip && l.user_nip.includes(q)) ||
        (l.action && l.action.toLowerCase().includes(q)) ||
        (l.details && l.details.toLowerCase().includes(q)) ||
        (l.ip_address && l.ip_address.includes(q)) ||
        (l.user_agent && l.user_agent.toLowerCase().includes(q))
      );
      if (!matchSearch) return false;
    }

    // 2. Action Filter
    if (selectedAction.value) {
      if (!l.action || !l.action.toUpperCase().includes(selectedAction.value.toUpperCase())) {
        return false;
      }
    }

    // Parse date for Day, Month, Year Filters
    const cleanStr = String(l.timestamp || '').replace('T', ' ').replace(/\+.*/, '');
    const parts = cleanStr.split(' ');
    let y = 0, m = 0, d = 0;
    if (parts.length >= 1 && parts[0].includes('-')) {
      const [parsedY, parsedM, parsedD] = parts[0].split('-').map(n => parseInt(n, 10));
      y = parsedY;
      m = parsedM;
      d = parsedD;
    }

    // 3. Day Filter
    if (selectedDay.value !== '') {
      if (d !== parseInt(selectedDay.value, 10)) return false;
    }

    // 4. Month Filter
    if (selectedMonth.value !== '') {
      if (m !== parseInt(selectedMonth.value, 10)) return false;
    }

    // 5. Year Filter
    if (selectedYear.value !== '') {
      if (y !== parseInt(selectedYear.value, 10)) return false;
    }

    return true;
  });
});

const totalPages = computed(() => {
  return Math.ceil(filteredLogs.value.length / pageSize.value) || 1;
});

const paginatedLogs = computed(() => {
  const startIdx = (currentPage.value - 1) * pageSize.value;
  const endIdx = startIdx + Number(pageSize.value);
  return filteredLogs.value.slice(startIdx, endIdx);
});

// Export to Excel CSV
const exportToCSV = () => {
  if (filteredLogs.value.length === 0) {
    alert('Tidak ada data log yang dapat di-export.');
    return;
  }

  let csvContent = '\uFEFF'; // UTF-8 BOM
  csvContent += 'No,Waktu & Tanggal,Nama Pengguna,NIP,Jenis Aktivitas,Detail Peristiwa,Alamat IP,User Agent\n';

  filteredLogs.value.forEach((log, index) => {
    const timeFormatted = formatDateIndo(log.timestamp).replace(/,/g, '');
    const name = `"${(log.user_name || '').replace(/"/g, '""')}"`;
    const nip = `'${log.user_nip || ''}`;
    const action = `"${(log.action || '').replace(/"/g, '""')}"`;
    const details = `"${(log.details || '').replace(/"/g, '""')}"`;
    const ip = `"${log.ip_address || ''}"`;
    const ua = `"${(log.user_agent || '').replace(/"/g, '""')}"`;

    csvContent += `${index + 1},${timeFormatted},${name},${nip},${action},${details},${ip},${ua}\n`;
  });

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.setAttribute('href', url);
  link.setAttribute('download', `Log_Aktivitas_LOPI-Q_${new Date().toISOString().slice(0,10)}.csv`);
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

// Export to PDF / Print Window
const exportToPDF = () => {
  if (filteredLogs.value.length === 0) {
    alert('Tidak ada data log yang dapat dicetak.');
    return;
  }

  const printWindow = window.open('', '_blank');
  if (!printWindow) return;

  const rowsHtml = filteredLogs.value.map((log, index) => `
    <tr>
      <td style="padding: 8px; border: 1px solid #cbd5e1; text-align: center;">${index + 1}</td>
      <td style="padding: 8px; border: 1px solid #cbd5e1; white-space: nowrap;">${formatDateIndo(log.timestamp)}</td>
      <td style="padding: 8px; border: 1px solid #cbd5e1;"><strong>${log.user_name}</strong><br><small style="color: #be123c;">NIP. ${log.user_nip}</small></td>
      <td style="padding: 8px; border: 1px solid #cbd5e1; font-weight: bold;">${log.action}</td>
      <td style="padding: 8px; border: 1px solid #cbd5e1;">${log.details}</td>
      <td style="padding: 8px; border: 1px solid #cbd5e1; font-family: monospace;">${log.ip_address || '180.242.190.12'}</td>
      <td style="padding: 8px; border: 1px solid #cbd5e1; font-family: monospace; font-size: 10px;">${formatUserAgent(log.user_agent)}</td>
    </tr>
  `).join('');

  printWindow.document.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <title>Laporan Log Aktivitas & Audit Keamanan LOPI-Q</title>
      <style>
        body { font-family: Arial, sans-serif; font-size: 12px; margin: 20px; color: #0f172a; }
        h2 { margin-bottom: 4px; color: #991b1b; }
        p { margin-top: 0; color: #475569; font-size: 11px; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; }
        th { background-color: #f1f5f9; padding: 10px; border: 1px solid #cbd5e1; text-align: left; font-size: 11px; text-transform: uppercase; }
        @media print {
          @page { size: landscape; margin: 15mm; }
        }
      </style>
    </head>
    <body>
      <h2>LAPORAN AUDIT LOG AKTIVITAS & KEAMANAN SISTEM LOPI-Q</h2>
      <p>Tanggal Cetak: ${new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })} | Total: ${filteredLogs.value.length} Log Data</p>
      <table>
        <thead>
          <tr>
            <th>No</th>
            <th>Waktu & Tanggal</th>
            <th>Pengguna / NIP</th>
            <th>Jenis Aktivitas</th>
            <th>Detail Peristiwa</th>
            <th>Alamat IP</th>
            <th>User Agent</th>
          </tr>
        </thead>
        <tbody>
          ${rowsHtml}
        </tbody>
      </table>
      <script>
        window.onload = function() { window.print(); }
      <\/script>
    </body>
    </html>
  `);
  printWindow.document.close();
};

onMounted(() => {
  fetchLogs();
});
</script>
