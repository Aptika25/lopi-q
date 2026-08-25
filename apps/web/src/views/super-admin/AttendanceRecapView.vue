<template>
  <AdminLayout>
    <div class="space-y-4 select-none font-sans">
      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4 mb-4">
        <div>
          <div class="flex items-center gap-2 flex-wrap">
            <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Rekapan Kehadiran &amp; Presensi Call Taker</h2>
            <span 
              class="px-2.5 py-0.5 text-[11px] font-black uppercase rounded-full tracking-wider border inline-flex items-center gap-1"
              :class="(schedulesData?.shiftMode || 2) === 3 ? 'bg-purple-100 text-purple-800 border-purple-300' : 'bg-emerald-100 text-emerald-800 border-emerald-300'"
            >
              ⚡ Skema Shift Aktif: {{ schedulesData?.shiftMode || 2 }} Shift / Hari
            </span>
          </div>
          <p class="font-sans text-slate-500 mt-1 text-xs">
            Rekapitulasi riwayat absensi, jam masuk/pulang, dan laporan log presensi siaga petugas Call Taker 112.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-wrap shrink-0">
          <button
            @click="exportExcel"
            class="h-9 px-4 text-xs font-semibold rounded-xl flex items-center gap-1.5 cursor-pointer whitespace-nowrap border-0
                   bg-gradient-to-b from-emerald-600 to-emerald-700 text-white shadow-xs
                   hover:from-emerald-500 hover:to-emerald-600 active:scale-[0.97] transition-all"
          >
            <span class="material-symbols-outlined text-[15px]">table_view</span> Excel
          </button>

          <button
            @click="exportPdf"
            class="h-9 px-4 text-xs font-semibold rounded-xl flex items-center gap-1.5 cursor-pointer whitespace-nowrap border-0
                   bg-gradient-to-b from-rose-600 to-rose-700 text-white shadow-xs
                   hover:from-rose-500 hover:to-rose-600 active:scale-[0.97] transition-all"
          >
            <span class="material-symbols-outlined text-[15px]">picture_as_pdf</span> PDF
          </button>

          <button
            @click="exportDoc"
            class="h-9 px-4 text-xs font-semibold rounded-xl flex items-center gap-1.5 cursor-pointer whitespace-nowrap border-0
                   bg-gradient-to-b from-blue-600 to-blue-700 text-white shadow-xs
                   hover:from-blue-500 hover:to-blue-600 active:scale-[0.97] transition-all"
          >
            <span class="material-symbols-outlined text-[15px]">description</span> Doc (Word)
          </button>
        </div>
      </div>

      <!-- Ringkasan Statistik Kehadiran -->
      <div class="grid grid-cols-2 sm:grid-cols-5 gap-3">
        <div class="bg-white p-3.5 rounded-xl border border-slate-200 shadow-2xs space-y-1">
          <div class="text-[10px] font-extrabold text-slate-400 uppercase tracking-wider">Total Penugasan</div>
          <div class="text-xl font-display font-black text-slate-900 font-mono">{{ statsSummary.total }} <span class="text-xs font-normal text-slate-400">Petugas</span></div>
        </div>
        <div class="bg-white p-3.5 rounded-xl border border-emerald-200 shadow-2xs space-y-1">
          <div class="text-[10px] font-extrabold text-emerald-700 uppercase tracking-wider flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span> Hadir Lengkap
          </div>
          <div class="text-xl font-display font-black text-emerald-950 font-mono">{{ statsSummary.hadir }}</div>
        </div>
        <div class="bg-white p-3.5 rounded-xl border border-amber-200 shadow-2xs space-y-1">
          <div class="text-[10px] font-extrabold text-amber-700 uppercase tracking-wider flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-amber-500"></span> Presensi Parsial
          </div>
          <div class="text-xl font-display font-black text-amber-950 font-mono">{{ statsSummary.parsial }}</div>
        </div>
        <div class="bg-white p-3.5 rounded-xl border border-blue-200 shadow-2xs space-y-1">
          <div class="text-[10px] font-extrabold text-blue-700 uppercase tracking-wider flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-blue-500"></span> Sakit / Izin
          </div>
          <div class="text-xl font-display font-black text-blue-950 font-mono">{{ statsSummary.sakit }}</div>
        </div>
        <div class="bg-white p-3.5 rounded-xl border border-rose-200 shadow-2xs space-y-1 col-span-2 sm:col-span-1">
          <div class="text-[10px] font-extrabold text-rose-700 uppercase tracking-wider flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-rose-500 animate-pulse"></span> Tidak Hadir (Alpha)
          </div>
          <div class="text-xl font-display font-black text-rose-950 font-mono">{{ statsSummary.alpha }}</div>
        </div>
      </div>

      <div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <!-- Notification -->
        <div v-if="actionMsg" class="mx-4 sm:mx-6 mt-3 px-4 py-2.5 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl text-xs font-semibold flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="material-symbols-outlined text-emerald-600 text-base">check_circle</span>
            {{ actionMsg }}
          </div>
          <button @click="actionMsg = ''" class="text-emerald-600 border-0 bg-transparent cursor-pointer ml-4 text-sm">✕</button>
        </div>

        <!-- Mode Tab Switcher: Hari Ini vs Bulanan -->
        <div class="flex items-center gap-2 border-b border-slate-200/80 px-4 sm:px-6 pt-3 bg-slate-50/70">
          <button
            @click="activeTab = 'today'"
            class="pb-2.5 px-4 text-xs font-bold transition-all relative border-0 bg-transparent cursor-pointer flex items-center gap-1.5"
            :class="activeTab === 'today' ? 'text-rose-700 font-black' : 'text-slate-500 hover:text-slate-800'"
          >
            <span class="material-symbols-outlined text-[17px]">today</span>
            <span>Presensi Hari Ini ({{ todayFormatted }})</span>
            <div v-if="activeTab === 'today'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-rose-600 rounded-t-full"></div>
          </button>

          <button
            @click="activeTab = 'monthly'"
            class="pb-2.5 px-4 text-xs font-bold transition-all relative border-0 bg-transparent cursor-pointer flex items-center gap-1.5"
            :class="activeTab === 'monthly' ? 'text-rose-700 font-black' : 'text-slate-500 hover:text-slate-800'"
          >
            <span class="material-symbols-outlined text-[17px]">calendar_month</span>
            <span>Rekapan Bulanan (Tgl 1 s/d Today)</span>
            <div v-if="activeTab === 'monthly'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-rose-600 rounded-t-full"></div>
          </button>
        </div>

        <!-- Filter Bar Container (2-Baris Rapi & Responsive) -->
        <div class="p-4 sm:p-5 border-b border-slate-200/80 bg-slate-50/50 space-y-3">
          
          <!-- Baris 1: Search + Filter Dropdowns -->
          <div class="flex flex-wrap items-center gap-2">
            <!-- Search Bar -->
            <div class="relative flex-1 min-w-[200px]">
              <span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-[16px]">search</span>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari nama, NIP, atau OPD..."
                class="w-full pl-9 pr-3 py-1.5 bg-white border border-slate-200 rounded-xl text-xs font-medium text-slate-800 focus:outline-none focus:border-rose-500 focus:ring-2 focus:ring-rose-500/20 transition-all shadow-xs"
              />
            </div>

            <!-- Filter Bulan & Tahun -->
            <div class="flex items-center gap-1 bg-rose-50/80 border border-rose-200/90 rounded-xl px-2.5 h-8 shadow-xs">
              <span class="material-symbols-outlined text-rose-700 text-[16px]">calendar_month</span>
              <select v-model="selectedFilterMonth" class="h-full bg-transparent text-xs text-rose-900 font-extrabold focus:outline-none cursor-pointer border-0">
                <option v-for="(m, idx) in monthNames" :key="idx" :value="idx">{{ m }}</option>
              </select>
              <select v-model="selectedFilterYear" class="h-full bg-transparent text-xs text-rose-900 font-extrabold focus:outline-none cursor-pointer border-0">
                <option v-for="y in [2025, 2026, 2027]" :key="y" :value="y">{{ y }}</option>
              </select>
            </div>

            <!-- Filter Shift -->
            <select v-model="selectedShift" class="h-8 px-3 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 focus:outline-none focus:ring-2 focus:ring-rose-500/20 cursor-pointer font-semibold shadow-xs">
              <option value="">Semua Shift</option>
              <option value="Shift 1">Shift 1</option>
              <option value="Shift 2">Shift 2</option>
              <option value="Shift 3">Shift 3</option>
            </select>

            <!-- Filter Tim -->
            <select v-model="selectedTeam" class="h-8 px-3 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 focus:outline-none focus:ring-2 focus:ring-rose-500/20 cursor-pointer font-semibold shadow-xs">
              <option value="">Semua Tim</option>
              <option value="A">Tim A</option>
              <option value="B">Tim B</option>
              <option value="C">Tim C</option>
              <option value="D">Tim D</option>
              <option value="E">Tim E</option>
            </select>

            <!-- Filter OPD -->
            <select v-model="selectedOpd" class="h-8 px-3 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 focus:outline-none focus:ring-2 focus:ring-rose-500/20 cursor-pointer shadow-xs">
              <option value="">Semua OPD</option>
              <option value="Dinas Sosial">Dinas Sosial</option>
              <option value="Badan Penanggulangan Bencana Daerah">BPBD</option>
              <option value="Dinas Kesehatan">Dinas Kesehatan</option>
              <option value="Dinas Perhubungan">Dinas Perhubungan</option>
              <option value="Satpol">Satpol PP</option>
            </select>

            <!-- Filter Status -->
            <select v-model="selectedStatus" class="h-8 px-3 bg-white border border-slate-200 rounded-xl text-xs text-slate-700 focus:outline-none focus:ring-2 focus:ring-rose-500/20 cursor-pointer shadow-xs">
              <option value="">Semua Status</option>
              <option value="Hadir Lengkap">🟢 Hadir Lengkap</option>
              <option value="Masuk Saja">🟡 Masuk Saja</option>
              <option value="Pulang Saja">⚪ Pulang Saja</option>
              <option value="Sakit / Izin (Resmi)">🔵 Sakit / Izin (Resmi)</option>
              <option value="Tidak Hadir (Alpha)">🔴 Tidak Hadir (Alpha)</option>
            </select>

            <!-- Reset Filter Button -->
            <button 
              v-if="searchQuery || selectedShift || selectedTeam || selectedOpd || selectedStatus"
              @click="searchQuery = ''; selectedShift = ''; selectedTeam = ''; selectedOpd = ''; selectedStatus = '';"
              class="h-8 px-3 bg-rose-50 hover:bg-rose-100 text-rose-700 font-extrabold rounded-xl border border-rose-200 text-xs cursor-pointer whitespace-nowrap transition-all"
            >
              Reset Filter
            </button>
          </div>

          <!-- Baris 2: Sub-header Status / Informasi & Counter Badge (Rapi, Terbaca Jelas & Tidak Terpotong) -->
          <div class="flex flex-wrap items-center justify-between gap-3 pt-2.5 border-t border-slate-200/60 text-xs">
            <div class="flex items-center gap-2 text-slate-500 font-medium">
              <span class="w-2 h-2 rounded-full bg-rose-600"></span>
              <span>
                Mode Laporan: <strong class="text-slate-800 font-bold">{{ activeTab === 'today' ? 'Presensi Hari Ini' : 'Rekapan Bulanan' }}</strong>
              </span>
            </div>

            <!-- Counter Badge Rapi -->
            <div class="px-3 py-1 bg-white border border-slate-200 rounded-full text-slate-600 font-semibold text-xs shadow-2xs flex items-center gap-1.5 shrink-0">
              <span>Menampilkan</span>
              <strong class="font-mono font-extrabold text-rose-700 text-sm">{{ groupedRecap.length }}</strong>
              <span>petugas {{ activeTab === 'today' ? '(Hari Ini)' : '(Bulanan)' }}</span>
            </div>
          </div>

        </div>

        <!-- Table -->
        <div class="overflow-x-auto w-full">
          <table class="w-full text-left text-xs table-auto">
            <thead>
              <tr class="border-b border-slate-100">
                <th class="py-3 px-3 text-[10px] font-extrabold text-slate-500 uppercase tracking-wider w-14 min-w-[3.5rem] whitespace-nowrap text-center">No.</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Tanggal</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Shift</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Tim</th>
                <th class="py-3 px-3 text-[10px] font-semibold text-slate-400 uppercase tracking-wider">Nama Petugas</th>
                <th class="py-3 px-3 text-[10px] font-semibold text-slate-400 uppercase tracking-wider">OPD Unit Kerja</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Jam Masuk</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Jam Pulang</th>
                <th class="py-3 px-2.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider whitespace-nowrap">Status</th>
                <th class="py-3 px-2.5 text-center text-[10px] font-semibold text-slate-400 uppercase tracking-wider w-20 whitespace-nowrap">Aksi Admin</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, idx) in groupedRecap"
                :key="row.key"
                class="border-b border-slate-50 hover:bg-slate-50/70 transition-colors"
                :class="{ 'bg-rose-50/30': row.status === 'Tidak Hadir (Alpha)' }"
              >
                <td class="py-3 px-2 text-center font-mono text-xs text-slate-500 font-extrabold align-middle">{{ idx + 1 }}</td>
                <td class="py-3 px-2.5 font-mono text-[11px] text-slate-500 whitespace-nowrap align-middle">{{ row.date }}</td>
                <td class="py-3 px-2.5 whitespace-nowrap align-middle">
                  <span v-if="row.shift" class="px-2 py-0.5 rounded-lg text-[10px] font-extrabold bg-slate-100 text-slate-700 border border-slate-200">
                    {{ row.shift }}
                  </span>
                  <span v-else class="text-slate-300 text-xs">—</span>
                </td>
                <td class="py-3 px-2.5 whitespace-nowrap align-middle">
                  <span v-if="row.team" class="px-2 py-0.5 rounded-lg text-[10px] font-black border" :class="getTeamBadgeClass(row.team)">
                    Tim {{ row.team }}
                  </span>
                  <span v-else class="text-slate-300 text-xs">—</span>
                </td>
                <td class="py-3 px-3 align-middle">
                  <div class="font-semibold text-slate-800 text-xs flex items-center gap-1.5 flex-wrap">
                    <span class="whitespace-nowrap">{{ row.name }}</span>
                    <span v-if="row.isReplacer" class="px-1.5 py-0.2 bg-amber-100 text-amber-800 text-[9px] font-extrabold rounded border border-amber-200">Pengganti</span>
                  </div>
                  <div class="font-mono text-[10px] text-slate-400 mt-0.5 whitespace-nowrap">NIP. {{ row.nip }}</div>
                </td>
                <td class="py-3 px-3 text-[11px] text-slate-600 align-middle leading-snug max-w-[200px]">{{ row.unit }}</td>
                <td class="py-3 px-2.5 whitespace-nowrap align-middle">
                  <span v-if="row.masuk" class="inline-flex items-center gap-1 font-mono text-[11px] font-medium text-emerald-700">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 flex-shrink-0"></span>
                    <span>{{ row.masuk }} <span class="text-[10px] font-normal text-emerald-600">WITA</span></span>
                  </span>
                  <span v-else-if="row.status === 'Tidak Hadir (Alpha)'" class="text-rose-500 font-mono text-[11px] font-semibold">Absen</span>
                  <span v-else class="text-slate-300 text-xs">—</span>
                </td>
                <td class="py-3 px-2.5 whitespace-nowrap align-middle">
                  <span v-if="row.pulang" class="inline-flex items-center gap-1 font-mono text-[11px] font-medium text-orange-600">
                    <span class="w-1.5 h-1.5 rounded-full bg-orange-400 flex-shrink-0"></span>
                    <span>{{ row.pulang }} <span class="text-[10px] font-normal text-orange-500">WITA</span></span>
                  </span>
                  <span v-else-if="row.status === 'Tidak Hadir (Alpha)'" class="text-rose-500 font-mono text-[11px] font-semibold">Absen</span>
                  <span v-else class="text-slate-300 text-xs">—</span>
                </td>
                <td class="py-3 px-2.5 align-middle whitespace-nowrap">
                  <span
                    class="inline-flex items-center px-2 py-0.5 rounded-lg text-[10px] font-semibold"
                    :class="{
                      'bg-emerald-50 text-emerald-700 border border-emerald-100': row.status === 'Hadir Lengkap',
                      'bg-amber-50 text-amber-700 border border-amber-100': row.status === 'Masuk Saja',
                      'bg-slate-100 text-slate-500 border border-slate-200': row.status === 'Pulang Saja',
                      'bg-blue-50 text-blue-700 border border-blue-200': row.status === 'Sakit / Izin (Resmi)',
                      'bg-rose-50 text-rose-700 border border-rose-200 font-bold': row.status === 'Tidak Hadir (Alpha)'
                    }"
                  >
                    <span class="w-1.5 h-1.5 rounded-full mr-1 flex-shrink-0"
                      :class="{
                        'bg-emerald-500': row.status === 'Hadir Lengkap',
                        'bg-amber-500': row.status === 'Masuk Saja',
                        'bg-slate-400': row.status === 'Pulang Saja',
                        'bg-blue-500': row.status === 'Sakit / Izin (Resmi)',
                        'bg-rose-500': row.status === 'Tidak Hadir (Alpha)'
                      }"
                    ></span>
                    {{ row.status }}
                  </span>
                </td>
                <td class="py-3 px-2.5 align-middle text-center whitespace-nowrap">
                  <button
                    @click="openEditModal(row)"
                    class="px-2 py-1 text-[11px] font-bold rounded-lg inline-flex items-center justify-center gap-1 cursor-pointer transition-all border
                           bg-white border-slate-300 text-slate-700 hover:bg-rose-50 hover:text-rose-700 hover:border-rose-300 shadow-2xs"
                    title="Edit Status / Jam Presensi Manual"
                  >
                    <span class="material-symbols-outlined text-[13px]">edit</span>
                    <span>Edit</span>
                  </button>
                </td>
              </tr>
              <tr v-if="groupedRecap.length === 0">
                <td colspan="10" class="py-12 text-center">
                  <div class="flex flex-col items-center gap-2 text-slate-400">
                    <span class="material-symbols-outlined text-[36px] text-slate-200">how_to_reg</span>
                    <span class="text-xs">{{ loading ? 'Memuat data dari database...' : 'Tidak ada data yang sesuai.' }}</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

      </div>
    </div>

    <!-- ========== MODAL EDIT PRESENSI MANUAL (SUPER ADMIN) ========== -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
      <div class="bg-white rounded-2xl border border-slate-200 shadow-2xl max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <!-- Header -->
        <div class="bg-slate-900 text-white px-6 py-4 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="material-symbols-outlined text-rose-500 text-xl">edit_square</span>
            <h3 class="font-bold text-sm">Edit Presensi Manual Call Taker</h3>
          </div>
          <button @click="showEditModal = false" class="text-slate-400 hover:text-white border-0 bg-transparent cursor-pointer text-base">✕</button>
        </div>

        <!-- Body -->
        <div class="p-6 space-y-4 text-xs">
          <!-- Officer summary box -->
          <div class="bg-slate-50 p-3.5 rounded-xl border border-slate-200 space-y-1">
            <div class="flex items-center justify-between">
              <span class="font-extrabold text-slate-900 text-sm">{{ editForm.name }}</span>
              <div class="flex items-center gap-1.5">
                <span v-if="editForm.shift" class="px-2 py-0.5 rounded text-[10px] font-extrabold bg-rose-50 text-rose-700 border border-rose-200">{{ editForm.shift }}</span>
                <span v-if="editForm.team" class="px-2 py-0.5 rounded text-[10px] font-black border" :class="getTeamBadgeClass(editForm.team)">Tim {{ editForm.team }}</span>
              </div>
            </div>
            <div class="text-slate-500 font-mono">NIP. {{ editForm.nip }} · {{ editForm.unit }}</div>
            <div class="text-rose-700 font-bold mt-1">Tanggal Presensi: {{ editForm.date }}</div>
          </div>

          <!-- Info Helper for Night Shift (Cross-midnight) -->
          <div v-if="isNightShiftEdit" class="p-3 bg-amber-50 border border-amber-200 rounded-xl text-[11px] text-amber-900 flex items-start gap-2">
            <span class="material-symbols-outlined text-[18px] text-amber-600 shrink-0 mt-0.5">info</span>
            <div>
              <strong>Akomodasi Shift Malam (Lintas Hari):</strong>
              <p class="mt-0.5 text-[10.5px] leading-relaxed">
                <span v-if="(schedulesData?.shiftMode || 2) === 2">
                  Sistem aktif: <strong>2 Shift/Hari</strong>. Shift 2 adalah Shift Malam (Jam Masuk: <code class="bg-amber-100/80 px-1 py-0.5 rounded font-mono">20:00:00</code> &amp; Jam Pulang: <code class="bg-amber-100/80 px-1 py-0.5 rounded font-mono">08:00:00</code> keesokan harinya).
                </span>
                <span v-else>
                  Sistem aktif: <strong>3 Shift/Hari</strong>. Shift 3 adalah Shift Dini Hari (Jam Masuk: <code class="bg-amber-100/80 px-1 py-0.5 rounded font-mono">23:30:00</code> &amp; Jam Pulang: <code class="bg-amber-100/80 px-1 py-0.5 rounded font-mono">07:30:00</code> keesokan harinya).
                </span>
                System backend otomatis menyimpan jam pulang pada tanggal hari berikutnya sehingga rekapan tetap akurat.
              </p>
            </div>
          </div>

          <!-- Input Jam Masuk -->
          <div class="space-y-1">
            <label class="font-bold text-slate-700 block flex items-center justify-between">
              <span>Jam Masuk (WITA)</span>
              <span class="text-[10px] text-slate-400 font-normal">Format: HH:MM:SS</span>
            </label>
            <input
              v-model="editForm.timeMasuk"
              type="text"
              placeholder="HH:MM:SS (contoh: 08:00:00)"
              class="w-full px-3 py-2 border border-slate-200 rounded-xl font-mono text-xs focus:outline-none focus:border-rose-500 bg-white"
            />
          </div>

          <!-- Input Jam Pulang -->
          <div class="space-y-1">
            <label class="font-bold text-slate-700 block flex items-center justify-between">
              <span>Jam Pulang (WITA)</span>
              <span class="text-[10px] text-slate-400 font-normal">Format: HH:MM:SS</span>
            </label>
            <input
              v-model="editForm.timePulang"
              type="text"
              placeholder="HH:MM:SS (contoh: 20:00:00)"
              class="w-full px-3 py-2 border border-slate-200 rounded-xl font-mono text-xs focus:outline-none focus:border-rose-500 bg-white"
            />
          </div>

          <!-- Dropdown Status Kehadiran -->
          <div class="space-y-1">
            <label class="font-bold text-slate-700 block">Status Kehadiran</label>
            <select
              v-model="editForm.status"
              class="w-full px-3 py-2 border border-slate-200 rounded-xl text-xs text-slate-800 font-bold focus:outline-none focus:border-rose-500 cursor-pointer"
            >
              <option value="Hadir Lengkap">🟢 Hadir Lengkap (Masuk & Pulang)</option>
              <option value="Masuk Saja">🟡 Masuk Saja</option>
              <option value="Pulang Saja">⚪ Pulang Saja</option>
              <option value="Sakit / Izin (Resmi)">🔵 Sakit / Izin (Resmi)</option>
              <option value="Tidak Hadir (Alpha)">🔴 Tidak Hadir (Alpha)</option>
            </select>
          </div>

          <!-- Input Catatan Admin -->
          <div class="space-y-1">
            <label class="font-bold text-slate-700 block">Catatan / Alasan Edit Admin</label>
            <textarea
              v-model="editForm.note"
              rows="2"
              placeholder="Contoh: HP petugas error camera saat scan, di-input manual oleh Admin."
              class="w-full px-3 py-2 border border-slate-200 rounded-xl text-xs text-slate-800 focus:outline-none focus:border-rose-500 resize-none"
            ></textarea>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="bg-slate-50 px-6 py-3.5 border-t border-slate-200 flex items-center justify-end gap-2">
          <button
            @click="showEditModal = false"
            class="px-4 py-2 bg-white border border-slate-200 hover:bg-slate-100 text-slate-700 text-xs font-bold rounded-xl cursor-pointer"
          >
            Batal
          </button>
          <button
            @click="saveManualPresensi"
            :disabled="savingEdit"
            class="px-5 py-2 bg-gradient-to-r from-rose-700 to-amber-600 hover:from-rose-600 hover:to-amber-500 text-white text-xs font-bold rounded-xl shadow-xs cursor-pointer disabled:opacity-50 flex items-center gap-1.5"
          >
            <span class="material-symbols-outlined text-[16px]">save</span>
            <span>{{ savingEdit ? 'Menyimpan...' : 'Simpan Perubahan Presensi' }}</span>
          </button>
        </div>
      </div>
    </div>

  </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import AdminLayout from '@/layouts/AdminLayout.vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const searchQuery = ref('');
const selectedOpd = ref('');
const attendanceRecords = ref([]);
const loading = ref(true);
const cleaning = ref(false);
const actionMsg = ref('');

const opdMap = {
  '199405032025211138': 'Dinas Sosial',
  '198703042025211061': 'Badan Penanggulangan Bencana Daerah',
  '200002062025211166': 'Dinas Kesehatan',
  '199002152025211114': 'Dinas Perhubungan',
  '199110052025211087': 'Satpol, Pemadam Kebakaran dan Penyelamatan',
  '198611302025211101': 'Dinas Sosial',
  '198603042025211147': 'Badan Penanggulangan Bencana Daerah',
  '197608022006041017': 'Dinas Kesehatan',
  '197605022006041017': 'Dinas Kesehatan',
  '198607122025211089': 'Dinas Perhubungan',
  '199603282025211050': 'Satpol, Pemadam Kebakaran dan Penyelamatan',
  'mappalua': 'Dinas Sosial',
  'suherman': 'Badan Penanggulangan Bencana Daerah',
  'riswandi': 'Dinas Kesehatan',
  'abil': 'Dinas Perhubungan',
  'imam': 'Satpol, Pemadam Kebakaran dan Penyelamatan',
  'rahim': 'Dinas Sosial',
  'munawir': 'Badan Penanggulangan Bencana Daerah',
  'abdullah': 'Dinas Kesehatan',
  'ismail': 'Dinas Perhubungan',
  'aldi': 'Satpol, Pemadam Kebakaran dan Penyelamatan'
};

function resolveUnitKerja(rawNip, rawName, dbUnit) {
  if (dbUnit && !['Diskominfo / OPD Bulukumba', 'Diskominfo Kab. Bulukumba', 'Diskominfo'].includes(dbUnit)) return dbUnit;
  const cleanNip = String(rawNip || '').replace(/\s+/g, '');
  if (cleanNip && opdMap[cleanNip]) return opdMap[cleanNip];
  const lowerName = String(rawName || '').toLowerCase();
  for (const key in opdMap) { if (lowerName.includes(key)) return opdMap[key]; }
  return 'Dinas Sosial';
}

const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
];

const now = new Date();
const activeTab = ref('today'); // 'today' | 'monthly'
const selectedTeam = ref('');
const selectedShift = ref('');
const selectedStatus = ref('');
const selectedFilterMonth = ref(now.getMonth());
const selectedFilterYear = ref(now.getFullYear());

const schedulesData = ref(null);
const leaveRequests = ref([]);

const todayFormatted = computed(() => {
  const now = new Date();
  return `${String(now.getDate()).padStart(2, '0')}-${String(now.getMonth() + 1).padStart(2, '0')}-${now.getFullYear()}`;
});

function getTeamBadgeClass(teamCode) {
  switch (String(teamCode || '').toUpperCase()) {
    case 'A': return 'bg-emerald-50 text-emerald-800 border-emerald-200';
    case 'B': return 'bg-amber-50 text-amber-800 border-amber-200';
    case 'C': return 'bg-indigo-50 text-indigo-800 border-indigo-200';
    case 'D': return 'bg-purple-50 text-purple-800 border-purple-200';
    case 'E': return 'bg-cyan-50 text-cyan-800 border-cyan-200';
    default: return 'bg-slate-100 text-slate-700 border-slate-200';
  }
}

const resolveLatestName = (nip, fallbackName) => {
  if (!nip && !fallbackName) return 'Call Taker';
  const cleanNip = String(nip || '').replace(/\s+/g, '');
  const fallbackClean = String(fallbackName || '').toLowerCase().trim();
  const firstName = fallbackClean.split(/[\s,.]+/)[0]; // e.g. "abdullah", "munawir"

  const userList = authStore?.usersList || [];
  if (!userList || userList.length === 0) return fallbackName || 'Call Taker';

  // 1. Try exact NIP match
  let user = userList.find((u) => {
    const uNip = String(u.nip || u.Nip || u.NIP || '').replace(/\s+/g, '');
    return uNip && uNip === cleanNip;
  });

  // 2. Try prefix/suffix NIP match (e.g. 1976...2006041017 handling 05 vs 08 month digit)
  if (!user && cleanNip.length >= 10) {
    const prefix = cleanNip.slice(0, 4);
    const suffix = cleanNip.slice(-8);
    user = userList.find((u) => {
      const uNip = String(u.nip || u.Nip || u.NIP || '').replace(/\s+/g, '');
      return uNip && uNip.startsWith(prefix) && uNip.endsWith(suffix);
    });
  }

  // 3. Try first name match (e.g. "abdullah")
  if (!user && firstName && firstName.length > 2) {
    user = userList.find((u) => {
      const uName = String(u.name || u.Name || '').toLowerCase();
      return uName.includes(firstName);
    });
  }

  return user?.name || user?.Name || fallbackName || 'Call Taker';
};

const fetchAttendanceHistory = async () => {
  loading.value = true;
  try {
    if (authStore && typeof authStore.fetchUsers === 'function') {
      try {
        await authStore.fetchUsers();
      } catch (err) {
        console.warn('Could not fetch users list:', err);
      }
    }
    const [resHist, resSched, resLeave] = await Promise.allSettled([
      axios.get('/api/presensi/history'),
      axios.get('/api/admin/schedules'),
      axios.get('/api/admin/leave-requests')
    ]);

    if (resHist.status === 'fulfilled' && resHist.value.data?.history) {
      const rawList = resHist.value.data.history || [];
      attendanceRecords.value = rawList.map((item) => {
        const rawTs = String(item.timestamp || '').replace('T', ' ').replace('Z', '').replace(/\+.*/, '');
        const parts = rawTs.split(' ');
        let datePart = parts[0] || '';
        let timePart = parts[1] || parts[0] || '--:--:--';
        if (datePart.includes('-')) {
          const p = datePart.split('-');
          if (p[0].length === 4) {
            datePart = `${p[2].padStart(2, '0')}-${p[1].padStart(2, '0')}-${p[0]}`;
          }
        }
        const cleanTime = (timePart.match(/\d{2}:\d{2}:\d{2}/) || ['--:--:--'])[0];
        const nip = item.user_nip || '-';
        return {
          id: item.id,
          date: datePart,
          time: cleanTime,
          nip,
          name: resolveLatestName(nip, item.user_name),
          unit: resolveUnitKerja(nip, item.user_name, item.unit_kerja),
          type: item.type || 'MASUK'
        };
      });
    }

    if (resSched.status === 'fulfilled' && resSched.value.data?.schedules) {
      schedulesData.value = resSched.value.data.schedules;
    }

    if (resLeave.status === 'fulfilled' && resLeave.value.data?.requests) {
      leaveRequests.value = resLeave.value.data.requests || [];
    }
  } catch (e) {
    console.error('Failed to fetch attendance data:', e);
  } finally {
    loading.value = false;
  }
};

const resolveFallbackShift = (dateStr, teamCode, clockInTime) => {
  if (schedulesData.value && Array.isArray(schedulesData.value.daysInMonth)) {
    const day = schedulesData.value.daysInMonth.find((d) => {
      const dStr = d.dateStr || d.date;
      if (!dStr) return false;
      const parts = dStr.split('-');
      const formatted = `${parts[2]}-${parts[1]}-${parts[0]}`;
      return formatted === dateStr || dStr === dateStr;
    });
    if (day) {
      if (day.shift1 === teamCode) return 'Shift 1';
      if (day.shift2 === teamCode) return 'Shift 2';
      if (day.shift3 === teamCode) return 'Shift 3';
    }
  }
  if (clockInTime) {
    const h = parseInt(clockInTime.split(':')[0] || '8', 10);
    if (h >= 6 && h < 16) return 'Shift 1';
    if (h >= 16 && h < 24) return 'Shift 2';
    return 'Shift 2';
  }
  return 'Shift 1';
};

const cleanNipStr = (nip) => String(nip || '').replace(/\D/g, '');

const normalizeDateStr = (rawDate) => {
  if (!rawDate) return '';
  const clean = String(rawDate).trim().split('T')[0].split(' ')[0];
  const parts = clean.split(/[-/]/);
  if (parts.length !== 3) return clean;
  let d = '', m = '', y = '';
  if (parts[0].length === 4) {
    y = parts[0];
    m = parts[1].padStart(2, '0');
    d = parts[2].padStart(2, '0');
  } else if (parts[2].length === 4) {
    d = parts[0].padStart(2, '0');
    m = parts[1].padStart(2, '0');
    y = parts[2];
  } else {
    return clean;
  }
  return `${d}-${m}-${y}`;
};

// 1. Full monthly recap starting from Date 1 (chronological 01-MM-YYYY -> 02-MM-YYYY -> ...)
const allMonthlyRecap = computed(() => {
  const map = new Map();
  const now = new Date();
  const todayStr = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;

  // Build scheduled duty roster from schedulesData if available
  if (schedulesData.value && Array.isArray(schedulesData.value.daysInMonth)) {
    const teams = schedulesData.value.teams || [];
    const masterMembers = new Map();
    const nipToTeamCode = new Map();

    teams.forEach((t) => {
      if (t.members) {
        t.members.forEach((m) => {
          if (m.nip) {
            const cleanNip = cleanNipStr(m.nip);
            masterMembers.set(cleanNip, m);
            nipToTeamCode.set(cleanNip, t.code);
          }
        });
      }
    });

    schedulesData.value.daysInMonth.forEach((day) => {
      const dStr = day.dateStr || day.date;
      if (!dStr) return;
      const dateFormatted = normalizeDateStr(dStr);
      if (!dateFormatted) return;

      // Convert dateFormatted (DD-MM-YYYY) back to YYYY-MM-DD for comparison with todayStr
      const [d, m, y] = dateFormatted.split('-');
      const dayIsoStr = `${y}-${m}-${d}`;
      if (dayIsoStr > todayStr) return; // Process from date 1 up to today

      // Dynamic Architecture: Use historical day.teamsSnapshot if available, fallback to current teams
      const rawDayTeams = (day.teamsSnapshot && Array.isArray(day.teamsSnapshot) && day.teamsSnapshot.length > 0)
        ? day.teamsSnapshot
        : teams;

      // Preserve historical Riswandi Risman in Tim B for past dates before Aug 13, 2026 cutoff
      const dayTeams = JSON.parse(JSON.stringify(rawDayTeams));
      if (dayIsoStr < '2026-08-13') {
        dayTeams.forEach((t) => {
          if (t.code === 'B' && Array.isArray(t.members) && t.members.length > 0) {
            const m0 = t.members[0];
            if (m0 && (cleanNipStr(m0.nip) === '199107052025211081' || m0.name.includes('Aswar'))) {
              t.members[0] = {
                name: 'Riswandi Risman',
                nip: '20000206 202521 1 166',
                jabatan: 'OPERATOR LAYANAN OPERASIONAL',
                unit: 'Dinas Kesehatan'
              };
            }
          }
        });
      }

      const dayMasterMembers = new Map();
      dayTeams.forEach((t) => {
        if (t.members) {
          t.members.forEach((m) => {
            if (m.nip) {
              const cleanNip = cleanNipStr(m.nip);
              dayMasterMembers.set(cleanNip, m);
            }
          });
        }
      });

      const shiftMode = schedulesData.value?.shiftMode || 2;
      const activeSlotTeams = [
        { code: day.shift1, slot: 'Shift 1' },
        { code: day.shift2, slot: 'Shift 2' }
      ];
      if (shiftMode === 3 && day.shift3) activeSlotTeams.push({ code: day.shift3, slot: 'Shift 3' });

      // Collect scheduled officers for active teams using dayTeams
      const scheduledOfficers = new Map();
      activeSlotTeams.forEach(({ code, slot }) => {
        const teamObj = dayTeams.find(t => t.code === code);
        if (teamObj && Array.isArray(teamObj.members)) {
          teamObj.members.forEach((m) => {
            if (m.nip) {
              const cleanNip = cleanNipStr(m.nip);
              scheduledOfficers.set(cleanNip, {
                nip: m.nip,
                name: m.name,
                unit: resolveUnitKerja(m.nip, m.name, m.unit),
                team: code,
                shift: slot,
                isReplacer: false
              });
            }
          });
        }
      });

      // Apply individual replacements for this day
      ['replacementsShift1', 'replacementsShift2', 'replacementsShift3'].forEach((rKey, idx) => {
        const slotTeamCode = activeSlotTeams[idx]?.code || 'A';
        const slotName = activeSlotTeams[idx]?.slot || `Shift ${idx + 1}`;
        const repls = day[rKey] || [];
        repls.forEach((r) => {
          if (r.replacedNip) scheduledOfficers.delete(cleanNipStr(r.replacedNip));
          if (r.replacerNip) {
            const replacerObj = dayMasterMembers.get(cleanNipStr(r.replacerNip)) || masterMembers.get(cleanNipStr(r.replacerNip));
            scheduledOfficers.set(cleanNipStr(r.replacerNip), {
              nip: r.replacerNip,
              name: r.replacerName || (replacerObj ? replacerObj.name : 'Pengganti'),
              unit: resolveUnitKerja(r.replacerNip, r.replacerName, replacerObj?.unit),
              team: slotTeamCode,
              shift: slotName,
              isReplacer: true
            });
          }
        });
      });

      // Insert scheduled officers into matrix
      scheduledOfficers.forEach((officer, cleanNip) => {
        const key = `${cleanNip}__${dateFormatted}`;
        map.set(key, {
          key,
          date: dateFormatted,
          nip: officer.nip,
          name: officer.name,
          unit: officer.unit,
          team: officer.team,
          shift: officer.shift,
          isReplacer: officer.isReplacer,
          masuk: null,
          pulang: null,
          ids: [],
          isScheduled: true
        });
      });
    });
  }

  // Helper to parse DD-MM-YYYY date string into a Date object
  const parseDateStr = (str) => {
    if (!str) return new Date();
    const normalized = normalizeDateStr(str);
    const [d, m, y] = normalized.split('-');
    return new Date(parseInt(y, 10), parseInt(m, 10) - 1, parseInt(d, 10));
  };

  // Helper to format Date object into DD-MM-YYYY string
  const formatDateObj = (dt) => {
    return `${String(dt.getDate()).padStart(2, '0')}-${String(dt.getMonth() + 1).padStart(2, '0')}-${dt.getFullYear()}`;
  };

  // Step 1: Shift-aware assignment of MASUK and PULANG records based on officer's scheduled shift
  attendanceRecords.value.forEach((rec) => {
    const cleanNip = String(rec.nip || '').replace(/\D/g, '');
    const recDateNorm = normalizeDateStr(rec.date);
    const todayKey = `${cleanNip}__${recDateNorm}`;
    const recTime = rec.time || '00:00:00';
    const hour = parseInt(recTime.split(':')[0] || '0', 10);

    const recDt = parseDateStr(recDateNorm);
    const prevDt = new Date(recDt.getTime() - 24 * 60 * 60 * 1000);
    const yesterdayDateNorm = formatDateObj(prevDt);
    const yesterdayKey = `${cleanNip}__${yesterdayDateNorm}`;

    // Primary Map Lookup
    let todayEntry = map.get(todayKey);
    // Secondary Fallback Lookup across Map values by NIP and Normalized Date
    if (!todayEntry) {
      for (const entry of map.values()) {
        if (cleanNipStr(entry.nip) === cleanNip && normalizeDateStr(entry.date) === recDateNorm) {
          todayEntry = entry;
          break;
        }
      }
    }

    let yesterdayEntry = map.get(yesterdayKey);
    if (!yesterdayEntry) {
      for (const entry of map.values()) {
        if (cleanNipStr(entry.nip) === cleanNip && normalizeDateStr(entry.date) === yesterdayDateNorm) {
          yesterdayEntry = entry;
          break;
        }
      }
    }

    const currentShiftMode = schedulesData.value?.shiftMode || 2;
    const yesterdayMasukHour = yesterdayEntry?.masuk ? parseInt(yesterdayEntry.masuk.split(':')[0] || '0', 10) : -1;
    const isYesterdayNightShift = yesterdayEntry && (
      (currentShiftMode === 2 && yesterdayEntry.shift === 'Shift 2') ||
      (currentShiftMode === 3 && yesterdayEntry.shift === 'Shift 3') ||
      (yesterdayMasukHour >= 18)
    );

    const isTodayNightShift = todayEntry && (
      (currentShiftMode === 2 && todayEntry.shift === 'Shift 2') ||
      (currentShiftMode === 3 && todayEntry.shift === 'Shift 3')
    );

    let targetEntry = null;
    let isMasukTarget = false;

    if (currentShiftMode === 3) {
      // ===== 3-SHIFT SYSTEM (Shift 1 Pagi 08-16, Shift 2 Sore 16-24, Shift 3 Malam 00-08) =====
      const isYesterdayShift3 = yesterdayEntry && yesterdayEntry.shift === 'Shift 3';
      const todayShift = todayEntry?.shift;

      if (hour < 12) {
        if (isYesterdayShift3 && (yesterdayEntry.masuk || !yesterdayEntry.pulang)) {
          targetEntry = yesterdayEntry;
          isMasukTarget = false;
        } else if (todayShift === 'Shift 1' || todayShift === 'Shift 3') {
          targetEntry = todayEntry;
          isMasukTarget = true;
        } else {
          targetEntry = todayEntry;
          isMasukTarget = hour < 12;
        }
      } else if (hour >= 12 && hour < 18) {
        if (todayShift === 'Shift 2') {
          targetEntry = todayEntry;
          isMasukTarget = true;
        } else if (todayShift === 'Shift 1') {
          targetEntry = todayEntry;
          isMasukTarget = false;
        } else {
          targetEntry = todayEntry;
          isMasukTarget = hour < 15;
        }
      } else {
        if (todayShift === 'Shift 3') {
          targetEntry = todayEntry;
          isMasukTarget = true;
        } else if (todayShift === 'Shift 2') {
          targetEntry = todayEntry;
          isMasukTarget = false;
        } else {
          targetEntry = todayEntry;
          isMasukTarget = true;
        }
      }
    } else {
      // ===== 2-SHIFT SYSTEM (Shift 1 Pagi 08-20, Shift 2 Malam 20-08) =====
      if (hour < 12) {
        if (isYesterdayNightShift && (yesterdayEntry.masuk || !yesterdayEntry.pulang)) {
          targetEntry = yesterdayEntry;
          isMasukTarget = false;
        } else {
          targetEntry = todayEntry;
          isMasukTarget = true;
        }
      } else {
        if (isTodayNightShift) {
          targetEntry = todayEntry;
          isMasukTarget = true;
        } else {
          targetEntry = todayEntry;
          isMasukTarget = false;
        }
      }
    }

    let fallbackTeam = '';
    if (schedulesData.value && Array.isArray(schedulesData.value.daysInMonth)) {
      const dayObj = schedulesData.value.daysInMonth.find(d => normalizeDateStr(d.dateStr || d.date) === recDateNorm);
      const dayTeams = (dayObj?.teamsSnapshot && Array.isArray(dayObj.teamsSnapshot) && dayObj.teamsSnapshot.length > 0)
        ? dayObj.teamsSnapshot
        : (schedulesData.value.teams || []);
      for (const t of dayTeams) {
        if (t.members && t.members.some(m => m.nip && cleanNipStr(m.nip) === cleanNip)) {
          fallbackTeam = t.code; break;
        }
      }
    } else if (schedulesData.value && Array.isArray(schedulesData.value.teams)) {
      for (const t of schedulesData.value.teams) {
        if (t.members && t.members.some(m => m.nip && cleanNipStr(m.nip) === cleanNip)) {
          fallbackTeam = t.code; break;
        }
      }
    }

    if (!targetEntry) {
      if (!map.has(todayKey)) {
        map.set(todayKey, {
          key: todayKey,
          date: recDateNorm,
          nip: rec.nip,
          name: rec.name,
          unit: rec.unit,
          team: fallbackTeam,
          shift: null,
          isReplacer: false,
          masuk: null,
          pulang: null,
          ids: [],
          isScheduled: false
        });
      }
      targetEntry = map.get(todayKey);
      isMasukTarget = rec.type === 'MASUK' || hour < 14;
    }

    if (targetEntry) {
      if (!targetEntry.team && fallbackTeam) targetEntry.team = fallbackTeam;
      if (isMasukTarget) {
        if (!targetEntry.masuk) targetEntry.masuk = recTime;
        else if (!targetEntry.pulang && recTime > targetEntry.masuk) targetEntry.pulang = recTime;
      } else {
        if (!targetEntry.pulang) targetEntry.pulang = recTime;
        else if (!targetEntry.masuk && recTime < targetEntry.pulang) targetEntry.masuk = recTime;
      }
      if (rec.id && !targetEntry.ids.includes(rec.id)) targetEntry.ids.push(rec.id);
    }
  });

function checkIsMasukLate(shift, masukTime, shiftMode) {
  if (!masukTime) return false;
  const timeClean = String(masukTime).trim().split(' ')[0] || '';
  const parts = timeClean.split(':');
  if (parts.length < 2) return false;
  const hour = parseInt(parts[0], 10);
  const min = parseInt(parts[1], 10);
  const totalMins = hour * 60 + min;

  const mode = shiftMode || 2;
  const s = String(shift || '');

  if (s.includes('1')) {
    const scheduledStartMins = (mode === 3) ? (7 * 60 + 30) : (8 * 60);
    const lateThresholdMins = scheduledStartMins + 30;
    if (totalMins >= 18 * 60) {
      return totalMins > (20 * 60 + 30);
    }
    return totalMins > lateThresholdMins;
  }
  if (s.includes('2')) {
    if (mode === 3) {
      return totalMins > (16 * 60);
    } else {
      if (totalMins >= 18 * 60) {
        return totalMins > (20 * 60 + 30);
      }
      return totalMins > (16 * 60 + 30);
    }
  }
  if (s.includes('3')) {
    return totalMins > 0 && totalMins < 7 * 60;
  }
  return false;
}

  // Final Deduplication & Merge Pass: Ensure exactly 1 row per officer per date
  const deduplicatedMap = new Map();

  for (const row of map.values()) {
    const cleanNip = cleanNipStr(row.nip);
    const dateNorm = normalizeDateStr(row.date);
    const dedupKey = `${cleanNip}__${dateNorm}`;

    if (!deduplicatedMap.has(dedupKey)) {
      deduplicatedMap.set(dedupKey, { ...row, date: dateNorm });
    } else {
      const existing = deduplicatedMap.get(dedupKey);
      // Merge properties into existing row (prioritizing scheduled team and non-null check-in/out times)
      if ((!existing.team || existing.team === '--') && row.team && row.team !== '--') {
        existing.team = row.team;
      }
      if (!existing.shift && row.shift) existing.shift = row.shift;
      if (!existing.masuk && row.masuk) existing.masuk = row.masuk;
      if (!existing.pulang && row.pulang) existing.pulang = row.pulang;
      if (row.isScheduled) existing.isScheduled = true;
      if (row.isReplacer) existing.isReplacer = true;
      if (Array.isArray(row.ids)) {
        row.ids.forEach((id) => {
          if (!existing.ids.includes(id)) existing.ids.push(id);
        });
      }
    }
  }

  // Resolve status & shift fallback & sort chronologically starting from Date 1
  return Array.from(deduplicatedMap.values())
    .map((row) => {
      const shift = row.shift || resolveFallbackShift(row.date, row.team, row.masuk);
      const isLate = checkIsMasukLate(shift, row.masuk, schedulesData.value?.shiftMode || 2);

      let status = 'Tidak Hadir (Alpha)';
      if (row.masuk && row.pulang) {
        status = isLate ? 'Hadir Terlambat' : 'Hadir Lengkap';
      } else if (row.masuk) {
        status = isLate ? 'Hadir Terlambat' : 'Masuk Saja';
      } else if (row.pulang) {
        status = 'Pulang Saja';
      } else {
        const hasLeave = leaveRequests.value.some((lr) => {
          const lrNip = cleanNipStr(lr.user_nip);
          const cleanRowNip = cleanNipStr(row.nip);
          const nipMatch = lrNip && cleanRowNip && lrNip === cleanRowNip;
          const statusMatch = lr.status === 'APPROVED';
          return nipMatch && statusMatch;
        });
        if (hasLeave) {
          status = 'Sakit / Izin (Resmi)';
        } else {
          status = 'Tidak Hadir (Alpha)';
        }
      }
      const name = resolveLatestName(row.nip, row.name);
      return { ...row, name, shift, status, isLate };
    })
    .sort((a, b) => {
      const aParts = a.date.split('-');
      const bParts = b.date.split('-');
      const aISO = `${aParts[2] || '0000'}-${aParts[1] || '00'}-${aParts[0] || '00'}`;
      const bISO = `${bParts[2] || '0000'}-${bParts[1] || '00'}-${bParts[0] || '00'}`;
      return aISO.localeCompare(bISO) || (a.shift || '').localeCompare(b.shift || '') || (a.team || '').localeCompare(b.team || '') || a.name.localeCompare(b.name);
    });
});

// 2. Filtered monthly recap (search + Month/Year + Shift + Team + OPD + status)
const allMonthlyRecapFiltered = computed(() => {
  const q = searchQuery.value.toLowerCase();
  const teamFilter = selectedTeam.value.toUpperCase();
  const shiftFilter = selectedShift.value;
  const opdFilter = selectedOpd.value.toLowerCase();
  const statusFilter = selectedStatus.value;
  const targetMonth = selectedFilterMonth.value;
  const targetYear = selectedFilterYear.value;

  return allMonthlyRecap.value.filter((row) => {
    const matchSearch = !q || row.name.toLowerCase().includes(q) || row.nip.includes(q) || row.unit.toLowerCase().includes(q);
    const matchTeam = !teamFilter || (row.team && row.team.toUpperCase() === teamFilter);
    const matchShift = !shiftFilter || (row.shift && row.shift === shiftFilter);
    const matchOpd = !opdFilter || row.unit.toLowerCase().includes(opdFilter);
    const matchStatus = !statusFilter || row.status === statusFilter;

    let matchMonthYear = true;
    if (row.date && row.date.includes('-')) {
      const parts = row.date.split('-');
      if (parts.length >= 3) {
        const m = parseInt(parts[1], 10) - 1;
        const y = parseInt(parts[2], 10);
        matchMonthYear = (m === targetMonth && y === targetYear);
      }
    }

    return matchSearch && matchTeam && matchShift && matchOpd && matchStatus && matchMonthYear;
  });
});


// 3. View table recap: switch between Today vs Monthly
const groupedRecap = computed(() => {
  if (activeTab.value === 'today') {
    return allMonthlyRecapFiltered.value.filter((r) => r.date === todayFormatted.value);
  }
  return allMonthlyRecapFiltered.value;
});

// 4. Summary metrics calculation
const statsSummary = computed(() => {
  const list = groupedRecap.value;
  let hadir = 0, parsial = 0, sakit = 0, alpha = 0;
  list.forEach((r) => {
    if (r.status === 'Hadir Lengkap') hadir++;
    else if (r.status === 'Masuk Saja' || r.status === 'Pulang Saja') parsial++;
    else if (r.status === 'Sakit / Izin (Resmi)') sakit++;
    else if (r.status === 'Tidak Hadir (Alpha)') alpha++;
  });
  return { total: list.length, hadir, parsial, sakit, alpha };
});

// 5. Admin Manual Presensi Edit Modal
const showEditModal = ref(false);
const savingEdit = ref(false);
const editForm = ref({
  name: '',
  nip: '',
  unit: '',
  team: '',
  shift: 'Shift 1',
  date: '',
  timeMasuk: '',
  timePulang: '',
  status: 'Hadir Lengkap',
  note: ''
});

const isNightShiftEdit = computed(() => {
  const shiftMode = schedulesData.value?.shiftMode || 2;
  if (shiftMode === 2) {
    return editForm.value.shift && editForm.value.shift.includes('2');
  }
  return editForm.value.shift && editForm.value.shift.includes('3');
});

function openEditModal(row) {
  const shiftMode = schedulesData.value?.shiftMode || 2;
  const isShift2 = row.shift && row.shift.includes('2');
  const isShift3 = row.shift && row.shift.includes('3');

  let defaultMasuk = '08:00:00';
  let defaultPulang = '20:00:00';

  if (shiftMode === 2) {
    if (isShift2) {
      defaultMasuk = '20:00:00';
      defaultPulang = '08:00:00';
    } else {
      defaultMasuk = '08:00:00';
      defaultPulang = '20:00:00';
    }
  } else {
    if (isShift3) {
      defaultMasuk = '23:30:00';
      defaultPulang = '07:30:00';
    } else if (isShift2) {
      defaultMasuk = '16:00:00';
      defaultPulang = '24:00:00';
    } else {
      defaultMasuk = '08:00:00';
      defaultPulang = '16:00:00';
    }
  }

  editForm.value = {
    name: row.name,
    nip: row.nip,
    unit: row.unit,
    team: row.team || '',
    shift: row.shift || 'Shift 1',
    date: row.date,
    timeMasuk: row.masuk || defaultMasuk,
    timePulang: row.pulang || defaultPulang,
    status: row.status,
    note: ''
  };
  showEditModal.value = true;
}

async function saveManualPresensi() {
  savingEdit.value = true;
  try {
    const res = await axios.post('/api/admin/presensi/manual-entry', {
      nip: editForm.value.nip,
      name: editForm.value.name,
      date: editForm.value.date,
      time_masuk: editForm.value.timeMasuk,
      time_pulang: editForm.value.timePulang,
      status: editForm.value.status,
      note: editForm.value.note
    });

    if (res.data?.success) {
      actionMsg.value = res.data.message || `Presensi ${editForm.value.name} berhasil diperbarui oleh Admin.`;
      showEditModal.value = false;
      await fetchAttendanceHistory();
    }
  } catch (e) {
    alert(e.response?.data?.error || 'Gagal menyimpan perubahan presensi.');
  } finally {
    savingEdit.value = false;
  }
}

// Export Excel (Filtered or Full List)
async function exportExcel() {
  try {
    const list = groupedRecap.value;
    if (!list || list.length === 0) {
      alert('Tidak ada data presensi yang dapat diexport.');
      return;
    }
    const XLSX = await import('xlsx');
    const headers = ['No.', 'Tanggal', 'Shift', 'Tim', 'Nama Call Taker', 'Keterangan', 'NIP', 'OPD Unit Kerja', 'Jam Masuk (WITA)', 'Jam Pulang (WITA)', 'Status Kehadiran'];
    const rows = list.map((r, i) => [
      i + 1,
      r.date,
      r.shift || '-',
      r.team ? `Tim ${r.team}` : '-',
      r.name,
      r.isReplacer ? 'Pengganti' : 'Terjadwal',
      r.nip,
      r.unit,
      r.masuk || '-',
      r.pulang || '-',
      r.status
    ]);
    const ws = XLSX.utils.aoa_to_sheet([headers, ...rows]);
    ws['!cols'] = [{ wch: 5 }, { wch: 14 }, { wch: 10 }, { wch: 10 }, { wch: 28 }, { wch: 12 }, { wch: 22 }, { wch: 36 }, { wch: 18 }, { wch: 18 }, { wch: 16 }];
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Rekap Presensi');
    XLSX.writeFile(wb, `Rekapan_Presensi_LOPI-Q_${new Date().toISOString().slice(0, 10)}.xlsx`);
  } catch (e) {
    alert('Gagal mengunduh Excel: ' + (e.message || e));
  }
}

// Export PDF (Filtered or Full List)
async function exportPdf() {
  try {
    const list = groupedRecap.value;
    if (!list || list.length === 0) {
      alert('Tidak ada data presensi yang dapat diexport ke PDF.');
      return;
    }
    const { default: jsPDF } = await import('jspdf');
    const { default: autoTable } = await import('jspdf-autotable');
    const doc = new jsPDF({ orientation: 'landscape', unit: 'mm', format: 'a4' });
    const W = doc.internal.pageSize.getWidth();

    const monthLabel = `${monthNames[selectedFilterMonth.value]} ${selectedFilterYear.value}`;

    doc.setFont('helvetica', 'bold'); doc.setFontSize(13);
    doc.text('PEMERINTAH KABUPATEN BULUKUMBA', W / 2, 12, { align: 'center' });
    doc.setFontSize(11);
    doc.text('POSKO TERPADU NTPD 112', W / 2, 17.5, { align: 'center' });
    doc.setFont('helvetica', 'bold'); doc.setFontSize(9.5);
    doc.text(`Laporan Rekapan Kehadiran & Presensi Call Taker Bulan ${monthLabel}`, W / 2, 23.5, { align: 'center' });
    doc.setFont('helvetica', 'normal'); doc.setFontSize(8);
    doc.text(`Dicetak Pada: ${new Date().toLocaleString('id-ID', { timeZone: 'Asia/Makassar' })} WITA`, W / 2, 28.5, { align: 'center' });
    doc.setLineWidth(0.4);
    doc.line(10, 31, W - 10, 31);

    autoTable(doc, {
      startY: 34,
      margin: { left: 10, right: 10 },
      head: [['No.', 'Tanggal', 'Shift', 'Tim', 'Nama & NIP', 'Ket.', 'OPD Unit Kerja', 'Jam Masuk', 'Jam Pulang', 'Status']],
      body: list.map((r, i) => [
        i + 1,
        r.date,
        r.shift || '—',
        r.team ? `Tim ${r.team}` : '—',
        `${r.name}\nNIP. ${r.nip}`,
        r.isReplacer ? 'Pengganti' : 'Terjadwal',
        r.unit,
        r.masuk ? `${r.masuk} WITA` : '—',
        r.pulang ? `${r.pulang} WITA` : '—',
        r.status
      ]),
      theme: 'grid',
      styles: { font: 'helvetica', fontSize: 7.5, cellPadding: 3, verticalAlign: 'middle' },
      headStyles: { fillColor: [30, 41, 59], textColor: 255, fontStyle: 'bold', fontSize: 8, halign: 'center' },
      alternateRowStyles: { fillColor: [248, 250, 252] },
      columnStyles: {
        0: { halign: 'center', cellWidth: 14 },
        1: { halign: 'center', cellWidth: 22 },
        2: { halign: 'center', cellWidth: 16 },
        3: { halign: 'center', cellWidth: 14 },
        4: { halign: 'left', cellWidth: 50 },
        5: { halign: 'center', cellWidth: 20 },
        6: { halign: 'left' },
        7: { halign: 'center', cellWidth: 25 },
        8: { halign: 'center', cellWidth: 25 },
        9: { halign: 'center', cellWidth: 28 }
      }
    });

    const totalPages = doc.internal.getNumberOfPages();
    const H = doc.internal.pageSize.getHeight();
    for (let i = 1; i <= totalPages; i++) {
      doc.setPage(i);
      doc.setFont('helvetica', 'normal');
      doc.setFontSize(7.5);
      doc.setTextColor(120);
      doc.text(`Halaman ${i} dari ${totalPages}`, W / 2, H - 5, { align: 'center' });
    }

    doc.save(`Rekapan_Presensi_LOPI-Q_${selectedFilterYear.value}_${String(selectedFilterMonth.value + 1).padStart(2, '0')}.pdf`);
  } catch (e) {
    alert('Gagal mengunduh PDF: ' + (e.message || e));
  }
}

// Export Word / DOC (Filtered or Full List)
function exportDoc() {
  try {
    const list = groupedRecap.value;
    if (!list || list.length === 0) {
      alert('Tidak ada data presensi yang dapat diexport ke Word.');
      return;
    }
    const monthLabel = `${monthNames[selectedFilterMonth.value]} ${selectedFilterYear.value}`;
    const content = `
      <html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word' xmlns='http://www.w3.org/TR/REC-html40'>
      <head>
        <meta charset='utf-8'>
        <title>Laporan Rekapan Kehadiran LOPI-Q</title>
        <style>
          @page WordSection1 { size: 841.9pt 595.3pt; mso-page-orientation: landscape; margin: 36pt; }
          div.WordSection1 { page: WordSection1; }
          body { font-family: Arial, sans-serif; font-size: 11pt; color: #0f172a; }
          .header { text-align: center; border-bottom: 2pt solid #1e293b; padding-bottom: 8pt; margin-bottom: 12pt; }
          .header h1 { font-size: 14pt; font-weight: bold; margin: 0; text-transform: uppercase; color: #0f172a; }
          .header h2 { font-size: 12pt; font-weight: bold; margin: 2pt 0; color: #1e293b; text-transform: uppercase; }
          .header p { font-size: 10.5pt; margin: 2pt 0; color: #334155; }
          table { width: 100%; border-collapse: collapse; margin-top: 10pt; }
          th { background-color: #1e293b; color: #ffffff; font-weight: bold; padding: 6pt; border: 1pt solid #94a3b8; font-size: 10pt; text-align: center; }
          td { padding: 5pt 6pt; border: 1pt solid #cbd5e1; font-size: 10pt; vertical-align: middle; }
          tr:nth-child(even) { background-color: #f8fafc; }
          .text-center { text-align: center; }
          .status-lengkap { color: #15803d; font-weight: bold; }
          .status-masuk { color: #b45309; font-weight: bold; }
          .status-pulang { color: #475569; font-weight: bold; }
          .status-sakit { color: #1d4ed8; font-weight: bold; }
          .status-alpha { color: #be123c; font-weight: bold; }
          .footer { margin-top: 20pt; font-size: 9.5pt; color: #64748b; text-align: center; }
        </style>
      </head>
      <body>
        <div class="WordSection1">
          <div class="header">
            <h1>PEMERINTAH KABUPATEN BULUKUMBA</h1>
            <h2>POSKO TERPADU NTPD 112</h2>
            <p><b>Laporan Rekapan Kehadiran &amp; Presensi Call Taker Bulan ${monthLabel}</b></p>
            <p style="font-size: 9.5pt; color: #64748b;">Dicetak Pada: ${new Date().toLocaleString('id-ID', { timeZone: 'Asia/Makassar' })} WITA</p>
          </div>

          <table>
            <thead>
              <tr>
                <th style="width: 5%; text-align: center;">No.</th>
                <th style="width: 10%; text-align: center;">Tanggal</th>
                <th style="width: 8%; text-align: center;">Shift</th>
                <th style="width: 7%; text-align: center;">Tim</th>
                <th style="width: 23%;">Nama &amp; NIP</th>
                <th style="width: 20%;">OPD Unit Kerja</th>
                <th style="width: 10%; text-align: center;">Jam Masuk</th>
                <th style="width: 10%; text-align: center;">Jam Pulang</th>
                <th style="width: 7%; text-align: center;">Status</th>
              </tr>
            </thead>
            <tbody>
              ${list.map((r, i) => `
                <tr>
                  <td class="text-center">${i + 1}</td>
                  <td class="text-center">${r.date}</td>
                  <td class="text-center">${r.shift || '—'}</td>
                  <td class="text-center"><strong>${r.team ? 'Tim ' + r.team : '—'}</strong></td>
                  <td>
                    <strong>${r.name}</strong>
                    ${r.isReplacer ? '<span style="display:inline-block; margin-left:6pt; background:#fef3c7; color:#92400e; border:1pt solid #f59e0b; border-radius:3pt; padding:1pt 4pt; font-size:8.5pt; font-weight:bold;">Pengganti</span>' : ''}<br/>
                    <span style="color:#64748b; font-size:9.5pt;">NIP. ${r.nip}</span>
                  </td>
                  <td>${r.unit}</td>
                  <td class="text-center">${r.masuk ? r.masuk + ' WITA' : '—'}</td>
                  <td class="text-center">${r.pulang ? r.pulang + ' WITA' : '—'}</td>
                  <td class="text-center ${r.status === 'Hadir Lengkap' ? 'status-lengkap' : (r.status === 'Masuk Saja' ? 'status-masuk' : (r.status === 'Sakit / Izin (Resmi)' ? 'status-sakit' : (r.status === 'Tidak Hadir (Alpha)' ? 'status-alpha' : 'status-pulang')))}"> ${r.status}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>

          <div class="footer">
            <p>Halaman 1</p>
          </div>
        </div>
      </body>
      </html>
    `;

    const blob = new Blob(['\ufeff', content], { type: 'application/msword' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `Rekapan_Presensi_LOPI-Q_${selectedFilterYear.value}_${String(selectedFilterMonth.value + 1).padStart(2, '0')}.doc`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  } catch (e) {
    alert('Gagal mengunduh file Word: ' + (e.message || e));
  }
}

onMounted(() => { fetchAttendanceHistory(); });
</script>

<style scoped>
.material-symbols-outlined {
  font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
}
</style>

