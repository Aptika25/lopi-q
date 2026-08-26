<template>
  <AdminLayout>
    <div class="w-full space-y-6 select-none font-sans text-slate-800">
      
      <!-- Header -->
      <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 border-b border-[#ddbfc5]/60 pb-6 w-full">
        <div>
          <h1 class="text-2xl md:text-3xl font-extrabold text-[#1b1c1c] tracking-tight flex items-center gap-2">
            <span class="material-symbols-outlined text-[#ab2c5d] text-[32px] fill" style="font-variation-settings: 'FILL' 1;">history</span>
            Activity Log &amp; Audit Trail
          </h1>
          <p class="text-sm text-[#574146] mt-1">Monitor system actions, user authentication, and attendance events real-time.</p>
        </div>
        
        <div class="flex items-center gap-3 shrink-0 flex-wrap sm:flex-nowrap">
          <button 
            @click="fetchLogs"
            :disabled="loading"
            class="py-2.5 px-4 bg-slate-100 hover:bg-slate-200 text-[#574146] font-bold text-xs rounded-lg border border-[#ddbfc5] flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50 transition-all shadow-xs"
          >
            <span class="material-symbols-outlined text-base" :class="loading ? 'animate-spin' : ''">refresh</span>
            <span>{{ loading ? 'Memuat Log...' : 'Refresh Audit' }}</span>
          </button>

          <button 
            @click="exportToCSV"
            class="py-2.5 px-4 bg-[#ffd9e4] text-[#ab2c5d] hover:bg-[#fec1d6] font-semibold text-xs rounded-full transition-colors flex items-center gap-2 cursor-pointer shadow-xs border border-[#ddbfc5] whitespace-nowrap"
          >
            <span class="material-symbols-outlined text-base">download</span>
            Export CSV
          </button>

          <button 
            @click="exportToPDF"
            class="py-2.5 px-4 bg-[#ab2c5d] text-white hover:bg-[#5e002b] font-semibold text-xs rounded-full transition-colors flex items-center gap-2 cursor-pointer shadow-xs border-0 whitespace-nowrap"
          >
            <span class="material-symbols-outlined text-base">print</span>
            Cetak / PDF
          </button>
        </div>
      </div>

      <!-- Filters (Bento Layout) -->
      <section class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Search Filter (Col-span 2) -->
        <div class="bg-white/80 backdrop-blur-md rounded-xl p-4 md:col-span-2 border border-[#ddbfc5] shadow-xs flex flex-col justify-center">
          <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-2 flex items-center gap-1">
            <span class="material-symbols-outlined text-sm">search</span> Search Activities
          </label>
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-lg">search</span>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Cari user, NIP, aksi, IP, atau User Agent..." 
              class="w-full pl-10 pr-4 py-2 bg-white border border-[#ddbfc5] rounded-lg text-xs font-medium text-[#1b1c1c] focus:outline-none focus:border-[#f06292] focus:ring-1 focus:ring-[#f06292]/50 transition-all shadow-xs"
            />
          </div>
        </div>

        <!-- Action / Status Filter -->
        <div class="bg-white/80 backdrop-blur-md rounded-xl p-4 border border-[#ddbfc5] shadow-xs flex flex-col justify-center">
          <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-2 flex items-center gap-1">
            <span class="material-symbols-outlined text-sm">filter_list</span> Jenis Aktivitas
          </label>
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-lg">tune</span>
            <select 
              v-model="selectedAction" 
              class="w-full pl-10 pr-8 py-2 bg-white border border-[#ddbfc5] rounded-lg text-xs font-bold text-[#1b1c1c] appearance-none focus:outline-none focus:border-[#f06292] shadow-xs cursor-pointer"
            >
              <option value="">Semua Aktivitas</option>
              <option value="LOGIN">LOGIN</option>
              <option value="PRESENSI">PRESENSI</option>
              <option value="2FA">2FA OTENTIKASI</option>
              <option value="USER">MANAJEMEN USER</option>
              <option value="SCHEDULE">JADWAL SHIFT</option>
            </select>
          </div>
        </div>

        <!-- Month Filter -->
        <div class="bg-white/80 backdrop-blur-md rounded-xl p-4 border border-[#ddbfc5] shadow-xs flex flex-col justify-center">
          <label class="text-[11px] font-bold text-[#574146] uppercase tracking-wider mb-2 flex items-center gap-1">
            <span class="material-symbols-outlined text-sm">calendar_today</span> Filter Bulan
          </label>
          <div class="relative">
            <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-[#8a7176] text-lg">date_range</span>
            <select 
              v-model="selectedMonth" 
              class="w-full pl-10 pr-8 py-2 bg-white border border-[#ddbfc5] rounded-lg text-xs font-bold text-[#1b1c1c] appearance-none focus:outline-none focus:border-[#f06292] shadow-xs cursor-pointer"
            >
              <option value="">Semua Bulan</option>
              <option v-for="(m, idx) in monthsIndo" :key="idx" :value="idx + 1">{{ m }}</option>
            </select>
          </div>
        </div>
      </section>

      <!-- Active Filters Tag & Reset -->
      <div v-if="searchQuery || selectedAction || selectedMonth" class="flex items-center gap-2 pt-1">
        <span class="text-xs text-[#574146] font-bold">Filter Aktif:</span>
        <button 
          @click="searchQuery = ''; selectedAction = ''; selectedMonth = '';"
          class="px-3 py-1 bg-[#ffd9e4] text-[#ab2c5d] hover:bg-[#fec1d6] font-bold text-xs rounded-full border border-[#ddbfc5] cursor-pointer flex items-center gap-1 transition-colors"
        >
          <span class="material-symbols-outlined text-sm">close</span> Reset Filter
        </button>
      </div>

      <!-- Activity Table Container (Chronos Template Styling) -->
      <section class="bg-white/90 backdrop-blur-md rounded-2xl border border-[#ddbfc5] overflow-hidden shadow-sm flex-1">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead class="bg-[#ffd9e4]/40 border-b border-[#ddbfc5]/60">
              <tr>
                <th class="text-[11px] font-bold text-[#574146] uppercase px-6 py-4 tracking-wider">User / Entity</th>
                <th class="text-[11px] font-bold text-[#574146] uppercase px-6 py-4 tracking-wider">Action Performed</th>
                <th class="text-[11px] font-bold text-[#574146] uppercase px-6 py-4 tracking-wider">Timestamp</th>
                <th class="text-[11px] font-bold text-[#574146] uppercase px-6 py-4 tracking-wider">Status</th>
                <th class="text-[11px] font-bold text-[#574146] uppercase px-6 py-4 tracking-wider">Client Info</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#ddbfc5]/30">
              <tr v-if="loading && logs.length === 0">
                <td colspan="5" class="px-6 py-12 text-center text-[#574146] font-semibold">
                  <span class="material-symbols-outlined animate-spin text-3xl text-[#ab2c5d] block mb-2">sync</span>
                  <span>Memuat catatan audit log dari server...</span>
                </td>
              </tr>
              <tr v-else-if="filteredLogs.length === 0">
                <td colspan="5" class="px-6 py-12 text-center text-[#574146] font-semibold">
                  Tidak ada catatan log aktivitas yang sesuai dengan kriteria filter.
                </td>
              </tr>
              <tr 
                v-for="log in paginatedLogs" 
                :key="log.id" 
                class="hover:bg-[#ffd9e4]/10 transition-colors group"
              >
                <!-- User / Entity -->
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="h-9 w-9 rounded-full bg-[#f06292]/20 flex items-center justify-center text-[#ab2c5d] font-bold text-xs border border-[#fec1d6]">
                      {{ getUserInitials(log.user_name) }}
                    </div>
                    <div>
                      <p class="text-xs font-bold text-[#1b1c1c]">{{ log.user_name || 'System Auto' }}</p>
                      <p class="text-[10.5px] font-mono text-[#ab2c5d]">NIP. {{ log.user_nip || '-' }}</p>
                    </div>
                  </div>
                </td>

                <!-- Action Performed -->
                <td class="px-6 py-4">
                  <p class="text-xs font-bold text-[#1b1c1c] flex items-center gap-2">
                    <span class="material-symbols-outlined text-[#8a7176] text-base">{{ getActionIcon(log.action) }}</span>
                    <span>{{ log.action }}</span>
                  </p>
                  <p class="text-[11px] text-[#574146] mt-0.5 max-w-xs truncate" :title="log.details">{{ log.details }}</p>
                </td>

                <!-- Timestamp -->
                <td class="px-6 py-4">
                  <p class="text-xs font-mono text-[#574146] whitespace-nowrap">{{ formatDateIndo(log.timestamp) }}</p>
                </td>

                <!-- Status -->
                <td class="px-6 py-4">
                  <span 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-[10px] font-extrabold uppercase tracking-wide border shadow-2xs"
                    :class="getActionBadgeClass(log.action)"
                  >
                    {{ getStatusText(log.action) }}
                  </span>
                </td>

                <!-- Client Info (IP & User Agent) -->
                <td class="px-6 py-4">
                  <div class="text-[11px] font-mono text-[#1b1c1c] font-semibold">{{ log.ip_address || '180.242.190.12' }}</div>
                  <div class="text-[10px] text-[#8a7176] font-mono truncate max-w-[140px]" :title="log.user_agent">
                    {{ formatUserAgent(log.user_agent) }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="filteredLogs.length > 0" class="bg-white px-6 py-4 border-t border-[#ddbfc5]/40 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p class="text-xs text-[#574146]">
            Showing <strong class="font-mono text-[#1b1c1c]">{{ ((currentPage - 1) * pageSize) + 1 }}</strong> to <strong class="font-mono text-[#1b1c1c]">{{ Math.min(currentPage * pageSize, filteredLogs.length) }}</strong> of <strong class="font-mono text-[#1b1c1c]">{{ filteredLogs.length }}</strong> entries
          </p>

          <div class="flex items-center gap-2">
            <button 
              @click="currentPage--" 
              :disabled="currentPage === 1"
              class="w-8 h-8 rounded-lg border border-[#ddbfc5] flex items-center justify-center text-[#574146] hover:bg-slate-100 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
            >
              <span class="material-symbols-outlined text-base">chevron_left</span>
            </button>

            <span class="px-3 py-1 bg-[#f06292] text-white font-bold font-mono rounded-lg text-xs shadow-xs">
              {{ currentPage }} / {{ totalPages }}
            </span>

            <button 
              @click="currentPage++" 
              :disabled="currentPage >= totalPages"
              class="w-8 h-8 rounded-lg border border-[#ddbfc5] flex items-center justify-center text-[#574146] hover:bg-slate-100 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
            >
              <span class="material-symbols-outlined text-base">chevron_right</span>
            </button>
          </div>
        </div>
      </section>

    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import axios from 'axios';
import AdminLayout from '@/layouts/AdminLayout.vue';

const searchQuery = ref('');
const selectedAction = ref('');
const selectedMonth = ref('');
const pageSize = ref(10);
const currentPage = ref(1);
const loading = ref(false);

const logs = ref([]);

watch([searchQuery, selectedAction, selectedMonth, pageSize], () => {
  currentPage.value = 1;
});

const monthsIndo = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
];

const getUserInitials = (nameStr) => {
  if (!nameStr) return 'SYS';
  const parts = nameStr.trim().split(' ');
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  return nameStr.substring(0, 2).toUpperCase();
};

const getActionIcon = (actionStr) => {
  if (!actionStr) return 'info';
  const act = actionStr.toUpperCase();
  if (act.includes('LOGIN')) return 'login';
  if (act.includes('PRESENSI')) return 'fingerprint';
  if (act.includes('2FA')) return 'security';
  if (act.includes('SCHEDULE')) return 'edit_document';
  return 'history';
};

const getStatusText = (actionStr) => {
  if (!actionStr) return 'SUCCESS';
  const act = actionStr.toUpperCase();
  if (act.includes('FAIL') || act.includes('GAGAL')) return 'FAILED';
  if (act.includes('WARN')) return 'WARNING';
  return 'SUCCESS';
};

const getActionBadgeClass = (actionStr) => {
  if (!actionStr) return 'bg-[#E8F5E9] text-[#1B5E20] border-[#A5D6A7]';
  const act = actionStr.toUpperCase();
  if (act.includes('FAIL') || act.includes('GAGAL')) {
    return 'bg-[#FCE4EC] text-[#C2185B] border-[#F8BBD0]';
  }
  if (act.includes('WARN')) {
    return 'bg-[#FFF8E1] text-[#F57F17] border-[#FFE082]';
  }
  return 'bg-[#E8F5E9] text-[#1B5E20] border-[#A5D6A7]';
};

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

  return `${day} ${monthName} ${year}, ${hStr}:${mStr}`;
};

const formatUserAgent = (uaStr) => {
  if (!uaStr) return 'Chrome (Windows)';
  if (uaStr.includes('Android')) return 'Chrome (Android)';
  if (uaStr.includes('iPhone') || uaStr.includes('iPad')) return 'Safari (iOS)';
  if (uaStr.includes('Windows')) return 'Chrome (Windows)';
  if (uaStr.includes('Macintosh') || uaStr.includes('Mac OS')) return 'Safari (macOS)';
  if (uaStr.length > 20) {
    return uaStr.substring(0, 18) + '...';
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

    if (selectedAction.value) {
      if (!l.action || !l.action.toUpperCase().includes(selectedAction.value.toUpperCase())) {
        return false;
      }
    }

    const cleanStr = String(l.timestamp || '').replace('T', ' ').replace(/\+.*/, '');
    const parts = cleanStr.split(' ');
    let m = 0;
    if (parts.length >= 1 && parts[0].includes('-')) {
      const [, parsedM] = parts[0].split('-').map(n => parseInt(n, 10));
      m = parsedM;
    }

    if (selectedMonth.value !== '') {
      if (m !== parseInt(selectedMonth.value, 10)) return false;
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

const exportToCSV = () => {
  if (filteredLogs.value.length === 0) {
    alert('Tidak ada data log yang dapat di-export.');
    return;
  }

  let csvContent = '\uFEFF';
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
      <td style="padding: 8px; border: 1px solid #cbd5e1;"><strong>${log.user_name || 'System Auto'}</strong><br><small style="color: #be123c;">NIP. ${log.user_nip || '-'}</small></td>
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
        h2 { margin-bottom: 4px; color: #ab2c5d; }
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

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
