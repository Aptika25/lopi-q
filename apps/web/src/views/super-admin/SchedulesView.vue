<template>
  <AdminLayout>
    <div class="space-y-6 font-sans select-none">
      
      <!-- Page Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200/60 pb-4">
        <div>
          <h2 class="font-display font-bold text-slate-900 text-base md:text-lg">Manajemen Jadwal Shift &amp; Tim Peserta Magang</h2>
          <p class="font-sans text-slate-500 mt-1 text-xs">Pengaturan dinamis pembagian tim, roster shift bulanan, override petugas pengganti (standby backup), dan ekspor dokumen Word (.doc).</p>
        </div>

        <div class="flex flex-wrap items-center gap-2 shrink-0">
          <button 
            @click="generateAutoRotation"
            class="py-2.5 px-4 bg-amber-400 hover:bg-amber-300 text-amber-950 font-bold text-xs rounded-xl shadow-xs transition-all flex items-center gap-1.5 cursor-pointer border-0"
          >
            <span class="material-symbols-outlined text-[16px]">auto_mode</span>
            <span>Generate Rotasi Shift</span>
          </button>

          <button 
            @click="exportToWordDoc"
            class="py-2.5 px-4 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-bold text-xs rounded-xl shadow-xs transition-all flex items-center gap-1.5 cursor-pointer border-0"
          >
            <span class="material-symbols-outlined text-[16px]">description</span>
            <span>Unduh Surat Tugas (.doc)</span>
          </button>

          <button 
            @click="saveSchedules"
            :disabled="saving"
            class="py-2.5 px-4 bg-gradient-to-r from-emerald-600 to-emerald-700 hover:from-emerald-700 hover:to-emerald-800 text-white font-bold text-xs rounded-xl shadow-xs disabled:opacity-50 transition-all flex items-center gap-1.5 cursor-pointer border-0"
          >
            <span class="material-symbols-outlined text-[16px]">save</span>
            <span>{{ saving ? 'Menyimpan...' : 'Simpan Jadwal' }}</span>
          </button>
        </div>
      </div>

      <!-- Main Navigation Tabs -->
      <div class="flex items-center gap-3 border-b border-slate-200 pb-1">
        <button 
          @click="activeTab = 'roster'"
          class="px-5 py-3 font-extrabold text-xs rounded-2xl transition-all border-0 cursor-pointer flex items-center gap-2"
          :class="activeTab === 'roster' ? 'bg-rose-700 text-white shadow-md' : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'"
        >
          <span class="material-symbols-outlined text-[18px]">calendar_month</span>
          <span>1. Pembagian Tim & Roster Shift</span>
        </button>

        <button 
          @click="activeTab = 'requests'"
          class="px-5 py-3 font-extrabold text-xs rounded-2xl transition-all border-0 cursor-pointer flex items-center gap-2 relative"
          :class="activeTab === 'requests' ? 'bg-rose-700 text-white shadow-md' : 'bg-white text-slate-600 hover:bg-slate-100 border border-slate-200'"
        >
          <span class="material-symbols-outlined text-[18px]">assignment_ind</span>
          <span>2. Pengajuan Sakit & Tukar Shift</span>
          <span v-if="pendingRequestsCount > 0" class="ml-1 px-2 py-0.5 bg-amber-400 text-amber-950 font-black text-[10px] rounded-full animate-bounce">
            {{ pendingRequestsCount }} Pending
          </span>
        </button>
      </div>

      <!-- Alert Toast -->
      <div v-if="alertMessage" class="p-4 bg-emerald-50 border border-emerald-200 rounded-2xl text-emerald-900 text-xs font-bold flex items-center justify-between shadow-xs">
        <div class="flex items-center gap-2">
          <span class="material-symbols-outlined text-emerald-600">check_circle</span>
          <span>{{ alertMessage }}</span>
        </div>
        <button @click="alertMessage = ''" class="text-emerald-700 hover:text-emerald-950 border-0 bg-transparent cursor-pointer font-bold">✕</button>
      </div>

      <!-- TAB 1: ROSTER SHIFT & PEMBAGIAN TIM DYNAMIS -->
      <div v-if="activeTab === 'roster'" class="space-y-6">
        
        <!-- Controls & Filter Bar -->
        <div class="bg-white rounded-3xl p-5 border border-slate-200 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4">
          <!-- Month & Year Selector -->
          <div class="flex items-center gap-3">
            <span class="material-symbols-outlined text-rose-700 text-[22px]">calendar_month</span>
            <div class="flex items-center gap-2">
              <select 
                v-model="selectedMonth" 
                @change="onMonthYearChange"
                class="px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer"
              >
                <option v-for="(m, idx) in monthNames" :key="idx" :value="idx">{{ m }}</option>
              </select>
              <select 
                v-model="selectedYear" 
                @change="onMonthYearChange"
                class="px-3 py-2 bg-slate-50 border border-slate-300 rounded-xl text-xs font-bold text-slate-800 focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer"
              >
                <option v-for="y in [now.getFullYear(), now.getFullYear() + 1]" :key="y" :value="y">{{ y }}</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Section 0: Penetapan Mode Operasional Shift (2 Shift vs 3 Shift) -->
        <div class="bg-gradient-to-r from-slate-900 via-rose-950 to-slate-900 rounded-3xl p-6 border border-rose-900/50 shadow-md text-white space-y-4">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-white/10 pb-4">
            <div class="flex items-start gap-3">
              <div class="p-3 bg-rose-500/20 border border-rose-400/30 rounded-2xl text-rose-400">
                <span class="material-symbols-outlined text-2xl">tune</span>
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <h2 class="text-lg font-display font-black text-white">Penetapan Mode Sistem Shift Siaga (2 Shift vs 3 Shift)</h2>
                  <span 
                    class="px-2.5 py-0.5 text-[11px] font-black uppercase rounded-full tracking-wider border"
                    :class="shiftMode === 3 ? 'bg-purple-500/20 text-purple-300 border-purple-400/40' : 'bg-emerald-500/20 text-emerald-300 border-emerald-400/40'"
                  >
                    ⚡ Mode Aktif: {{ shiftMode }} Shift / Hari
                  </span>
                </div>
                <p class="text-xs text-slate-300 mt-1">Penetapan skema jam kerja ini otomatis menyesuaikan pendeteksian scan presensi pada aplikasi Peserta Magang dan tabel rekapan admin.</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-1">
            <!-- Option 2 Shift -->
            <div 
              @click="setShiftMode(2)"
              class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex flex-col justify-between"
              :class="shiftMode === 2 
                ? 'bg-rose-900/40 border-rose-500 ring-2 ring-rose-500/30 shadow-lg shadow-rose-950/50 text-white' 
                : 'bg-slate-800/40 border-slate-700/60 hover:border-slate-500 text-slate-400 hover:text-slate-200'"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2 font-black text-sm">
                  <span class="material-symbols-outlined text-rose-400">schedule</span>
                  <span>Sistem 2 Shift / Hari (12 Jam Kerja)</span>
                </div>
                <span v-if="shiftMode === 2" class="material-symbols-outlined text-rose-400 font-bold">check_circle</span>
              </div>
              <ul class="text-xs space-y-1 text-slate-300">
                <li>• <strong>Shift 1 (Pagi/Siang)</strong>: 08.00 – 20.00 WITA</li>
                <li>• <strong>Shift 2 (Malam)</strong>: 20.00 – 08.00 WITA (Esok Hari)</li>
              </ul>
            </div>

            <!-- Option 3 Shift -->
            <div 
              @click="setShiftMode(3)"
              class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex flex-col justify-between"
              :class="shiftMode === 3 
                ? 'bg-purple-900/40 border-purple-500 ring-2 ring-purple-500/30 shadow-lg shadow-purple-950/50 text-white' 
                : 'bg-slate-800/40 border-slate-700/60 hover:border-slate-500 text-slate-400 hover:text-slate-200'"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2 font-black text-sm">
                  <span class="material-symbols-outlined text-purple-400">more_time</span>
                  <span>Sistem 3 Shift / Hari (8 Jam Kerja)</span>
                </div>
                <span v-if="shiftMode === 3" class="material-symbols-outlined text-purple-400 font-bold">check_circle</span>
              </div>
              <ul class="text-xs space-y-1 text-slate-300">
                <li>• <strong>Shift 1 (Pagi)</strong>: 08.00 – 16.00 WITA</li>
                <li>• <strong>Shift 2 (Sore)</strong>: 16.00 – 24.00 WITA</li>
                <li>• <strong>Shift 3 (Malam)</strong>: 00.00 – 08.00 WITA</li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Section 1: Dynamic Team Assignment (Manual & Auto-Shuffle) -->
        <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-4">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-100 pb-3">
            <div>
              <h2 class="text-base font-display font-black text-slate-800 flex items-center gap-2">
                <span class="material-symbols-outlined text-rose-700">groups</span>
                <span>Pembagian Tim Peserta Magang Dinamis (2 Petugas per Tim)</span>
              </h2>
              <p class="text-xs text-slate-500 mt-0.5">Atur anggota Tim A–E secara manual melalui dropdown, atau acak tim secara otomatis.</p>
            </div>

            <div class="flex items-center gap-2">
              <button 
                @click="shuffleTeams"
                class="px-3.5 py-2 bg-rose-50 hover:bg-rose-100 text-rose-800 font-extrabold text-xs rounded-xl border border-rose-200 flex items-center gap-1.5 cursor-pointer"
                title="Acak pembagian tim secara otomatis"
              >
                <span class="material-symbols-outlined text-[16px]">casino</span>
                <span>Acak Pembagian Tim</span>
              </button>
            </div>
          </div>

          <!-- 5 Dynamic Team Cards -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
            <div 
              v-for="team in teams" 
              :key="team.code" 
              class="p-4 rounded-3xl border transition-all space-y-3.5"
              :class="team.colorClass"
            >
              <!-- Team Header -->
              <div class="flex items-center justify-between border-b border-slate-200/80 pb-2.5">
                <span class="font-display font-black text-lg tracking-tight text-slate-900">TIM {{ team.code }}</span>
                <span class="text-[10px] font-extrabold px-2.5 py-0.5 rounded-full bg-white/90 border border-slate-200 shadow-2xs text-slate-700">2 PETUGAS</span>
              </div>

              <!-- Members List -->
              <div class="space-y-3.5">
                <div v-for="(member, idx) in team.members" :key="idx" class="space-y-1.5">
                  <div class="flex items-center justify-between">
                    <span class="text-[10px] font-extrabold text-slate-400 tracking-wider uppercase">Anggota {{ idx + 1 }}</span>
                    <span v-if="member.nip && isMemberInactive(member.nip)" class="text-[9px] font-black px-1.5 py-0.5 bg-amber-100 text-amber-900 rounded border border-amber-300 animate-pulse">⚠️ Nonaktif</span>
                  </div>
                  
                  <!-- Clean Dropdown with Clear/Empty option -->
                  <div class="relative">
                    <select 
                      v-model="member.nip"
                      @change="onMemberSelect(team.code, idx, member.nip)"
                      class="w-full py-2 px-2 border rounded-xl text-xs font-bold focus:outline-none focus:ring-2 focus:ring-rose-500 cursor-pointer shadow-2xs hover:border-slate-400 transition-all truncate"
                      :class="[
                        member.nip && isMemberInactive(member.nip) ? 'bg-amber-50 border-amber-400 text-amber-950 font-black' : (member.nip ? 'bg-white border-slate-300/80 text-slate-900' : 'bg-slate-50 border-dashed border-slate-300 text-slate-400')
                      ]"
                    >
                      <!-- Empty / clear slot option -->
                      <option value="">— Kosongkan Slot —</option>
                      <!-- All officers not yet used in other slots (includes current selection naturally) -->
                      <option 
                        v-for="off in availableOfficersFor(team.code, idx)" 
                        :key="off.nip" 
                        :value="off.nip"
                      >
                        {{ off.name }}{{ isMemberInactive(off.nip) ? ' (Nonaktif)' : '' }}
                      </option>
                    </select>
                  </div>

                  <!-- Officer Meta Info Card (only shown when slot is filled) -->
                  <div v-if="member.nip" class="p-2.5 rounded-xl border space-y-1.5 transition-all shadow-2xs" :class="isMemberInactive(member.nip) ? 'bg-amber-50/90 border-amber-300' : 'bg-white/80 border-slate-200/80'">
                    <div class="text-[10px] text-rose-800 font-mono font-extrabold">NIP. {{ member.nip }}</div>
                    <div class="text-[10px] text-slate-800 font-extrabold leading-tight uppercase tracking-tight">{{ member.jabatan }}</div>
                    <div v-if="isMemberInactive(member.nip)" class="mt-1 p-1.5 bg-amber-100 text-amber-950 rounded-lg border border-amber-300 flex items-center gap-1 text-[9.5px] font-black">
                      <span class="material-symbols-outlined text-[13px] text-amber-700">warning</span>
                      <span>Akun Nonaktif - Silakan Ganti</span>
                    </div>
                    <div v-else class="inline-block text-[9px] font-extrabold px-2 py-0.5 bg-slate-100 text-slate-700 rounded-md border border-slate-200/90 mt-0.5">
                      {{ member.unit }}
                    </div>
                  </div>

                  <!-- Empty Slot Indicator -->
                  <div v-else class="p-2.5 bg-slate-50/80 rounded-xl border border-dashed border-slate-300 text-center">
                    <div class="text-[10px] text-slate-400 font-bold">Slot Kosong</div>
                    <div class="text-[9px] text-slate-300">Pilih petugas dari dropdown</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Section 2: Tabel Jadwal Shift Bulanan & Override -->
        <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-4">
          <div class="flex items-center justify-between border-b border-slate-100 pb-3">
            <div>
              <h2 class="text-base font-display font-black text-slate-800 flex items-center gap-2">
                <span class="material-symbols-outlined text-rose-700">edit_calendar</span>
                <span>Jadwal Shift Petugas Bulan {{ monthNames[selectedMonth] }} {{ selectedYear }}</span>
              </h2>
              <p class="text-xs text-slate-500 mt-0.5">Atur tim utama & tentukan **Petugas Pengganti (Standby Backup)** jika ada petugas yang berhalangan/sakit.</p>
            </div>

            <div class="text-xs font-mono font-bold text-slate-600 bg-slate-100 px-3 py-1.5 rounded-xl border border-slate-200">
              Total {{ daysInMonth.length }} Hari
            </div>
          </div>

          <!-- Roster Schedule Table -->
          <div class="overflow-x-auto rounded-2xl border border-slate-200 shadow-xs">
            <table class="w-full text-left text-xs text-slate-800">
              <thead class="bg-slate-100 text-slate-700 font-bold uppercase tracking-wider border-b border-slate-200">
                <tr>
                  <th class="py-3.5 px-4 w-44">Hari / Tanggal</th>
                  <th v-if="shiftMode === 2" class="py-3.5 px-4 min-w-[200px]">
                    <div>Shift Pagi (08.00 – 20.00)</div>
                    <div class="text-[9px] font-medium text-emerald-700 normal-case tracking-normal mt-0.5">Tim bertugas + pengganti jika ada sakit</div>
                  </th>
                  <th v-if="shiftMode === 2" class="py-3.5 px-4 min-w-[200px]">
                    <div>Shift Malam (20.00 – 08.00)</div>
                    <div class="text-[9px] font-medium text-amber-700 normal-case tracking-normal mt-0.5">Tim bertugas + pengganti jika ada sakit</div>
                  </th>
                  <th v-if="shiftMode === 3" class="py-3.5 px-4 min-w-[200px]">Shift Pagi (08.00 – 16.00)</th>
                  <th v-if="shiftMode === 3" class="py-3.5 px-4 min-w-[200px]">Shift Sore (16.00 – 24.00)</th>
                  <th v-if="shiftMode === 3" class="py-3.5 px-4 min-w-[200px]">Shift Malam (00.00 – 08.00)</th>
                  <th class="py-3.5 px-4 min-w-[220px]">
                    <div>Petugas Cadangan (Libur / Standby)</div>
                    <div class="text-[9px] font-medium text-blue-700 normal-case tracking-normal mt-0.5">Siap dipanggil jika ada petugas shift yang berhalangan</div>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-200 bg-white">
                <tr v-for="day in daysInMonth" :key="day.dateStr" class="transition-colors align-top" :class="isDayPast(day.dateStr) ? 'bg-slate-50/60' : 'hover:bg-slate-50'">
                  <!-- Day & Date Column -->
                  <td class="py-3 px-4 font-mono font-bold whitespace-nowrap" :class="day.isSunday ? 'text-rose-700' : 'text-slate-900'">
                    <div class="font-sans text-xs font-bold">{{ day.dayName }}</div>
                    <div class="text-[11px] text-slate-500 font-mono">{{ day.formattedDate }}</div>
                    <div v-if="isDayPast(day.dateStr)" class="mt-1 inline-flex items-center gap-1 px-1.5 py-0.5 bg-slate-200/80 text-slate-600 rounded text-[9px] font-extrabold border border-slate-300">
                      <span class="material-symbols-outlined text-[10px]">lock</span>
                      <span>Selesai (Terkunci)</span>
                    </div>
                    <div v-else-if="day.isSunday" class="mt-1 inline-flex items-center gap-1 px-1.5 py-0.5 bg-rose-100 text-rose-700 rounded text-[9px] font-extrabold">
                      <span class="material-symbols-outlined text-[10px]">weekend</span>
                      <span>Libur</span>
                    </div>
                  </td>

                  <!-- Shift 1 / Pagi Column -->
                  <td class="py-2 px-3 space-y-2">
                    <!-- Tim yang bertugas -->
                    <select 
                      v-model="day.shift1" 
                      @change="calculateOffTeams(day)"
                      :disabled="isDayPast(day.dateStr)"
                      class="w-full py-1.5 px-2 border text-emerald-950 font-black rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      :class="isDayPast(day.dateStr) ? 'opacity-50 cursor-not-allowed bg-slate-100 border-slate-200' : 'bg-emerald-50 border-emerald-300 cursor-pointer'"
                    >
                      <option v-for="t in teamCodes" :key="t" :value="t">Tim {{ t }} — {{ getTeamNames(t, day) }}</option>
                    </select>

                    <!-- Anggota tim bertugas shift 1 dengan pengganti per-orang -->
                    <div class="space-y-2 pt-1">
                      <div v-for="m in getTeamMembers(day.shift1, day)" :key="m.nip" class="rounded-xl border border-emerald-200 overflow-hidden">
                        <!-- Orang yang bertugas -->
                        <div class="flex items-center gap-2 px-2 py-1.5 bg-emerald-50">
                          <span class="material-symbols-outlined text-emerald-600 text-[13px] shrink-0">person</span>
                          <div class="min-w-0 flex-1">
                            <div class="text-[10px] font-black text-emerald-900 leading-tight truncate">{{ m.name }}</div>
                            <div class="text-[9px] text-slate-400 font-mono">{{ m.unit.replace('Badan Penanggulangan Bencana Daerah','BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan','Satpol PP').replace('Dinas Perhubungan','Dishub').replace('Dinas Kesehatan','Dinkes').replace('Dinas Sosial','Dinsos') }}</div>
                          </div>
                        </div>
                        <!-- Pengganti untuk orang ini jika sakit -->
                        <div v-if="day.offTeams.length > 0" class="px-2 py-1.5 bg-white space-y-1">
                          <div class="text-[8px] text-rose-600 font-extrabold uppercase tracking-wider flex items-center gap-0.5">
                            <span class="material-symbols-outlined text-[10px]">swap_horiz</span>
                            <span>Jika sakit, diganti oleh:</span>
                          </div>
                          <select
                            :disabled="isDayPast(day.dateStr)"
                            :value="getReplacerFor(day.replacementsShift1, m.nip)"
                            @change="setReplacerFor(day, 'replacementsShift1', m, ($event.target as HTMLSelectElement).value)"
                            class="w-full py-1 px-2 border text-slate-800 font-bold rounded-lg text-[10px]"
                            :class="isDayPast(day.dateStr) ? 'opacity-50 cursor-not-allowed bg-slate-100 border-slate-200' : 'bg-slate-50 border-rose-200 cursor-pointer focus:ring-1 focus:ring-rose-400'"
                          >
                            <option value="">— Tidak ada pengganti —</option>
                            <optgroup v-for="tc in day.offTeams" :key="tc" :label="'Tim ' + tc + ' (Libur)'">
                              <option v-for="r in getTeamMembers(tc, day)" :key="r.nip" :value="r.nip">
                                {{ r.name }} · {{ r.unit.replace('Badan Penanggulangan Bencana Daerah','BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan','Satpol PP').replace('Dinas Perhubungan','Dishub').replace('Dinas Kesehatan','Dinkes').replace('Dinas Sosial','Dinsos') }}
                              </option>
                            </optgroup>
                          </select>
                          <!-- Konfirmasi pengganti aktif -->
                          <div v-if="getReplacerFor(day.replacementsShift1, m.nip)" class="flex items-center gap-1 text-[9px] text-emerald-700 font-bold bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-200">
                            <span class="material-symbols-outlined text-[11px]">check_circle</span>
                            <span>{{ getReplacerNameFor(day.replacementsShift1, m.nip) }} → Ganti {{ m.name.split(',')[0] }}</span>
                          </div>
                        </div>
                        <div v-else class="px-2 py-1 bg-slate-50 text-[9px] text-slate-400 italic">Tidak ada tim libur hari ini</div>
                      </div>
                    </div>
                  </td>

                  <!-- Shift 2 / Malam or Sore Column -->
                  <td class="py-2 px-3 space-y-2">
                    <!-- Tim yang bertugas -->
                    <select 
                      v-model="day.shift2" 
                      @change="calculateOffTeams(day)"
                      :disabled="isDayPast(day.dateStr)"
                      class="w-full py-1.5 px-2 border text-amber-950 font-black rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-amber-500"
                      :class="isDayPast(day.dateStr) ? 'opacity-50 cursor-not-allowed bg-slate-100 border-slate-200' : 'bg-amber-50 border-amber-300 cursor-pointer'"
                    >
                      <option v-for="t in teamCodes" :key="t" :value="t">Tim {{ t }} — {{ getTeamNames(t, day) }}</option>
                    </select>

                    <!-- Anggota tim bertugas shift 2 dengan pengganti per-orang -->
                    <div class="space-y-2 pt-1">
                      <div v-for="m in getTeamMembers(day.shift2, day)" :key="m.nip" class="rounded-xl border border-amber-200 overflow-hidden">
                        <!-- Orang yang bertugas -->
                        <div class="flex items-center gap-2 px-2 py-1.5 bg-amber-50">
                          <span class="material-symbols-outlined text-amber-600 text-[13px] shrink-0">person</span>
                          <div class="min-w-0 flex-1">
                            <div class="text-[10px] font-black text-amber-900 leading-tight truncate">{{ m.name }}</div>
                            <div class="text-[9px] text-slate-400 font-mono">{{ m.unit.replace('Badan Penanggulangan Bencana Daerah','BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan','Satpol PP').replace('Dinas Perhubungan','Dishub').replace('Dinas Kesehatan','Dinkes').replace('Dinas Sosial','Dinsos') }}</div>
                          </div>
                        </div>
                        <!-- Pengganti untuk orang ini jika sakit -->
                        <div v-if="day.offTeams.length > 0" class="px-2 py-1.5 bg-white space-y-1">
                          <div class="text-[8px] text-rose-600 font-extrabold uppercase tracking-wider flex items-center gap-0.5">
                            <span class="material-symbols-outlined text-[10px]">swap_horiz</span>
                            <span>Jika sakit, diganti oleh:</span>
                          </div>
                          <select
                            :disabled="isDayPast(day.dateStr)"
                            :value="getReplacerFor(day.replacementsShift2, m.nip)"
                            @change="setReplacerFor(day, 'replacementsShift2', m, ($event.target as HTMLSelectElement).value)"
                            class="w-full py-1 px-2 border text-slate-800 font-bold rounded-lg text-[10px]"
                            :class="isDayPast(day.dateStr) ? 'opacity-50 cursor-not-allowed bg-slate-100 border-slate-200' : 'bg-slate-50 border-rose-200 cursor-pointer focus:ring-1 focus:ring-rose-400'"
                          >
                            <option value="">— Tidak ada pengganti —</option>
                            <optgroup v-for="tc in day.offTeams" :key="tc" :label="'Tim ' + tc + ' (Libur)'">
                              <option v-for="r in getTeamMembers(tc, day)" :key="r.nip" :value="r.nip">
                                {{ r.name }} · {{ r.unit.replace('Badan Penanggulangan Bencana Daerah','BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan','Satpol PP').replace('Dinas Perhubungan','Dishub').replace('Dinas Kesehatan','Dinkes').replace('Dinas Sosial','Dinsos') }}
                              </option>
                            </optgroup>
                          </select>
                          <!-- Konfirmasi pengganti aktif -->
                          <div v-if="getReplacerFor(day.replacementsShift2, m.nip)" class="flex items-center gap-1 text-[9px] text-emerald-700 font-bold bg-emerald-50 px-1.5 py-0.5 rounded border border-emerald-200">
                            <span class="material-symbols-outlined text-[11px]">check_circle</span>
                            <span>{{ getReplacerNameFor(day.replacementsShift2, m.nip) }} → Ganti {{ m.name.split(',')[0] }}</span>
                          </div>
                        </div>
                        <div v-else class="px-2 py-1 bg-slate-50 text-[9px] text-slate-400 italic">Tidak ada tim libur hari ini</div>
                      </div>
                    </div>
                  </td>

                  <!-- Shift 3 / Malam Column (if 3-shift mode) -->
                  <td v-if="shiftMode === 3" class="py-2 px-3 space-y-2">
                    <select 
                      v-model="day.shift3" 
                      @change="calculateOffTeams(day)"
                      :disabled="isDayPast(day.dateStr)"
                      class="w-full py-1.5 px-2 border text-indigo-950 font-black rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500"
                      :class="isDayPast(day.dateStr) ? 'opacity-50 cursor-not-allowed bg-slate-100 border-slate-200' : 'bg-indigo-50 border-indigo-300 cursor-pointer'"
                    >
                      <option v-for="t in teamCodes" :key="t" :value="t">Tim {{ t }} — {{ getTeamNames(t, day) }}</option>
                    </select>

                    <!-- Anggota tim bertugas shift 3 -->
                    <div class="space-y-0.5 pl-1">
                      <div v-for="m in getTeamMembers(day.shift3, day)" :key="m.nip" class="flex items-start gap-1 text-[10px] text-indigo-900">
                        <span class="material-symbols-outlined text-indigo-600 text-[11px] mt-0.5 shrink-0">person</span>
                        <div>
                          <div class="font-bold leading-tight">{{ m.name }}</div>
                          <div class="text-slate-400 text-[9px] font-mono">{{ m.unit.replace('Badan Penanggulangan Bencana Daerah', 'BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan', 'Satpol PP').replace('Dinas Perhubungan', 'Dishub').replace('Dinas Kesehatan', 'Dinkes').replace('Dinas Sosial', 'Dinsos') }}</div>
                        </div>
                      </div>
                    </div>
                  </td>

                  <!-- Off / Standby Backup Column — clearly shows who is available to replace -->
                  <td class="py-3 px-3">
                    <div v-if="day.offTeams.length > 0" class="space-y-2">
                      <div class="text-[9px] font-extrabold text-blue-700 uppercase tracking-wider flex items-center gap-1 mb-1">
                        <span class="material-symbols-outlined text-[11px]">emergency</span>
                        <span>Cadangan siap dipanggil:</span>
                      </div>
                      <div v-for="tc in day.offTeams" :key="tc" class="bg-blue-50 border border-blue-200 rounded-xl p-2 space-y-1.5">
                        <!-- Team header -->
                        <div class="flex items-center justify-between">
                          <div class="flex items-center gap-1.5">
                            <div class="w-5 h-5 rounded-lg bg-blue-600 text-white flex items-center justify-center text-[9px] font-black shrink-0">{{ tc }}</div>
                            <span class="text-[10px] font-extrabold text-blue-900">Tim {{ tc }}</span>
                          </div>
                          <span class="px-1.5 py-0.5 bg-blue-100 text-blue-700 rounded text-[8px] font-extrabold border border-blue-200">LIBUR</span>
                        </div>
                        <!-- Member list with full info -->
                        <div v-for="m in getTeamMembers(tc, day)" :key="m.nip" class="flex items-start gap-1.5 pl-1">
                          <span class="material-symbols-outlined text-blue-400 text-[12px] mt-0.5 shrink-0">person_check</span>
                          <div class="min-w-0">
                            <div class="text-[10px] font-bold text-slate-900 leading-tight">{{ m.name }}</div>
                            <div class="text-[9px] text-blue-700 font-mono">{{ m.nip }}</div>
                            <div class="text-[9px] text-slate-500">{{ m.unit.replace('Badan Penanggulangan Bencana Daerah', 'BPBD').replace('Satpol, Pemadam Kebakaran dan Penyelamatan', 'Satpol PP').replace('Dinas Perhubungan', 'Dishub').replace('Dinas Kesehatan', 'Dinkes').replace('Dinas Sosial', 'Dinsos') }}</div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div v-else class="text-center py-3">
                      <span class="material-symbols-outlined text-slate-300 text-xl">group_work</span>
                      <div class="text-[9px] text-slate-400 font-bold mt-0.5">Semua tim bertugas</div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Official Notes Footer -->
          <div class="bg-rose-50/70 border border-rose-200 p-4 rounded-2xl text-xs text-rose-950 space-y-1">
            <div class="font-bold text-rose-900 mb-1 flex items-center gap-1.5">
              <span class="material-symbols-outlined text-[16px]">assignment_late</span>
              <span>Catatan Resmi Surat Tugas Shift:</span>
            </div>
            <p><strong>a.</strong> Untuk pergantian shift diwajibkan untuk melapor satu hari sebelum pergantian shift.</p>
            <p><strong>b.</strong> Untuk absen alasan sakit wajib menyertakan surat keterangan sakit.</p>
            <p><strong>c.</strong> Petugas yang libur dapat dijadikan <em>standby backup</em> jika ada yang berhalangan hadir atau sakit.</p>
          </div>
        </div>
      </div>

      <!-- TAB 2: PERSETUJUAN PENGAJUAN SAKIT & TUKAR SHIFT -->
      <div v-if="activeTab === 'requests'" class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm space-y-6">
        <div class="flex items-center justify-between border-b border-slate-100 pb-4">
          <div>
            <h2 class="text-base font-display font-black text-slate-900 flex items-center gap-2">
              <span class="material-symbols-outlined text-rose-700">assignment_ind</span>
              <span>Daftar Pengajuan Sakit / Izin / Tukar Shift Peserta Magang</span>
            </h2>
            <p class="text-xs text-slate-500 mt-0.5">Verifikasi dokumen bukti dokter dan setujui penugasan petugas pengganti (standby backup).</p>
          </div>

          <button 
            @click="fetchLeaveRequests"
            class="px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs rounded-xl border border-slate-300 cursor-pointer flex items-center gap-1.5"
          >
            <span class="material-symbols-outlined text-[16px]">refresh</span>
            <span>Refresh Data</span>
          </button>
        </div>

        <!-- Request List Table -->
        <div v-if="leaveRequests.length === 0" class="text-center py-12 bg-slate-50 rounded-2xl border border-dashed border-slate-300">
          <span class="material-symbols-outlined text-slate-300 text-5xl">task_alt</span>
          <p class="text-xs font-bold text-slate-500 mt-2">Belum ada pengajuan sakit atau tukar shift dari Peserta Magang.</p>
        </div>

        <div v-else class="overflow-x-auto rounded-2xl border border-slate-200">
          <table class="w-full text-left text-xs text-slate-800">
            <thead class="bg-slate-100 text-slate-700 font-bold uppercase border-b border-slate-200">
              <tr>
                <th class="py-3.5 px-4">Tanggal Pengajuan</th>
                <th class="py-3.5 px-4">Nama Petugas</th>
                <th class="py-3.5 px-4">Kategori</th>
                <th class="py-3.5 px-4">Tanggal Shift</th>
                <th class="py-3.5 px-4">Petugas Pengganti (Standby)</th>
                <th class="py-3.5 px-4">Bukti / Alasan</th>
                <th class="py-3.5 px-4">Status</th>
                <th class="py-3.5 px-4 text-center">Aksi Persetujuan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 bg-white">
              <tr v-for="req in leaveRequests" :key="req.id" class="hover:bg-slate-50 transition-colors">
                <td class="py-3 px-4 font-mono text-slate-500 whitespace-nowrap">{{ req.created_at }}</td>
                <td class="py-3 px-4">
                  <div class="font-bold text-slate-900">{{ req.user_name }}</div>
                  <div class="text-[10px] font-mono text-rose-700">NIP. {{ req.user_nip }}</div>
                </td>
                <td class="py-3 px-4">
                  <span 
                    class="px-2.5 py-1 text-[11px] font-bold rounded-full border"
                    :class="{
                      'bg-rose-50 text-rose-800 border-rose-200': req.category === 'Sakit',
                      'bg-amber-50 text-amber-800 border-amber-200': req.category === 'Izin',
                      'bg-blue-50 text-blue-800 border-blue-200': req.category === 'Tukar Shift'
                    }"
                  >
                    {{ req.category }}
                  </span>
                </td>
                <td class="py-3 px-4 font-mono font-bold text-slate-900 whitespace-nowrap">{{ req.shift_date }}</td>
                <td class="py-3 px-4">
                  <div v-if="req.replacement_name" class="font-bold text-emerald-800 bg-emerald-50 px-2.5 py-1 rounded-lg border border-emerald-200 text-[11px] inline-block">
                    {{ req.replacement_name }}
                  </div>
                  <span v-else class="text-slate-400 italic">Ditentukan Admin</span>
                </td>
                <td class="py-3 px-4">
                  <p class="text-[11px] text-slate-700 font-medium mb-1">{{ req.reason }}</p>
                  <button 
                    v-if="req.attachment_url" 
                    @click="previewImage(req.attachment_url)"
                    class="text-[10px] font-extrabold text-blue-700 hover:underline border-0 bg-transparent cursor-pointer flex items-center gap-1"
                  >
                    <span class="material-symbols-outlined text-[14px]">attachment</span>
                    <span>Lihat Surat Dokter / Lampiran</span>
                  </button>
                </td>
                <td class="py-3 px-4 whitespace-nowrap">
                  <span 
                    class="px-2.5 py-1 rounded-full text-[10px] font-extrabold border"
                    :class="{
                      'bg-amber-100 text-amber-900 border-amber-300': req.status === 'PENDING',
                      'bg-emerald-100 text-emerald-900 border-emerald-300': req.status === 'APPROVED',
                      'bg-rose-100 text-rose-900 border-rose-300': req.status === 'REJECTED'
                    }"
                  >
                    {{ req.status }}
                  </span>
                </td>
                <td class="py-3 px-4 text-center whitespace-nowrap">
                  <div v-if="req.status === 'PENDING'" class="flex items-center justify-center gap-2">
                    <button 
                      @click="approveRequest(req)"
                      class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs rounded-xl border-0 shadow-2xs cursor-pointer flex items-center gap-1"
                      title="Setujui dan ganti petugas di roster"
                    >
                      <span class="material-symbols-outlined text-[16px]">check</span>
                      <span>Setujui</span>
                    </button>
                    <button 
                      @click="rejectRequest(req)"
                      class="px-3 py-1.5 bg-rose-100 hover:bg-rose-200 text-rose-800 font-bold text-xs rounded-xl border border-rose-200 cursor-pointer flex items-center gap-1"
                      title="Tolak pengajuan"
                    >
                      <span class="material-symbols-outlined text-[16px]">close</span>
                      <span>Tolak</span>
                    </button>
                  </div>
                  <span v-else class="text-[10px] text-slate-400 font-bold">Selesai</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Preview Image Modal -->
      <div v-if="showPreviewModal" class="fixed inset-0 bg-slate-900/70 backdrop-blur-xs z-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-3xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
          <div class="flex items-center justify-between border-b border-slate-200 pb-3">
            <h3 class="text-sm font-bold text-slate-900">Bukti Lampiran Surat Dokter</h3>
            <button @click="showPreviewModal = false" class="text-slate-400 hover:text-slate-700 text-lg font-bold border-0 bg-transparent cursor-pointer">✕</button>
          </div>
          <div class="flex justify-center p-2 bg-slate-100 rounded-2xl border border-slate-200">
            <img :src="previewImageUrl" alt="Surat Dokter" class="max-h-96 rounded-xl object-contain" />
          </div>
          <div class="flex justify-end">
            <button @click="showPreviewModal = false" class="px-4 py-2 bg-slate-800 text-white text-xs font-bold rounded-xl cursor-pointer border-0">Tutup</button>
          </div>
        </div>
      </div>

      <!-- Printable Surat Tugas Modal / Preview -->
      <div v-if="showPrintModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs z-50 flex items-center justify-center p-4 overflow-y-auto">
        <div class="bg-white rounded-3xl max-w-4xl w-full p-6 sm:p-8 space-y-6 shadow-2xl max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between border-b border-slate-200 pb-4">
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined text-blue-600 text-2xl">description</span>
              <h3 class="text-lg font-display font-black text-slate-900">Pratinjau Dokumen Surat Tugas (.doc)</h3>
            </div>
            <button @click="showPrintModal = false" class="text-slate-400 hover:text-slate-700 text-xl font-bold border-0 bg-transparent cursor-pointer">✕</button>
          </div>

          <!-- Printable Document Container Preview -->
          <div class="bg-white p-8 border border-slate-300 shadow-sm rounded-xl space-y-6 text-slate-900 font-sans text-xs">
            <!-- Header Surat -->
            <div class="text-left space-y-1 pb-4 border-b-2 border-slate-900">
              <div class="font-bold uppercase tracking-wider text-[11px]">LAMPIRAN SURAT TUGAS</div>
              <div class="font-bold">Nomor : ............................................................</div>
              <div class="font-bold">Tanggal : 30 Juli 2026</div>
            </div>

            <!-- Judul Surat -->
            <div class="text-center py-2 space-y-1">
              <h2 class="font-black text-sm uppercase tracking-wide">DAFTAR NAMA PETUGAS PESERTA MAGANG LAYANAN NTPD 112</h2>
              <h3 class="font-bold text-xs uppercase">KAB. BULUKUMBA</h3>
            </div>

            <!-- Tabel Tim -->
            <table class="w-full border-collapse border border-slate-900 text-xs text-left">
              <thead>
                <tr class="bg-slate-100 font-bold uppercase border-b border-slate-900 text-center">
                  <th class="border border-slate-900 py-2 px-3 w-12">TIM</th>
                  <th class="border border-slate-900 py-2 px-3">NAMA / NIP</th>
                  <th class="border border-slate-900 py-2 px-3">JABATAN</th>
                  <th class="border border-slate-900 py-2 px-3">UNIT KERJA</th>
                </tr>
              </thead>
              <tbody>
                <template v-for="team in teams" :key="team.code">
                  <tr v-for="(m, idx) in team.members" :key="idx" class="border-b border-slate-900">
                    <td v-if="idx === 0" :rowspan="2" class="border border-slate-900 font-black text-center text-lg align-middle">{{ team.code }}</td>
                    <td class="border border-slate-900 py-2 px-3 font-semibold">{{ m.name }} / NIP. {{ m.nip }}</td>
                    <td class="border border-slate-900 py-2 px-3 uppercase font-medium">{{ m.jabatan }}</td>
                    <td class="border border-slate-900 py-2 px-3">{{ m.unit }}</td>
                  </tr>
                </template>
              </tbody>
            </table>

            <!-- Title Jadwal -->
            <div class="text-center pt-4 space-y-1">
              <h2 class="font-black text-sm uppercase tracking-wide">JADWAL SHIFT PETUGAS PESERTA MAGANG LAYANAN NTPD 112</h2>
              <h3 class="font-bold text-xs uppercase">KAB. BULUKUMBA BULAN {{ monthNames[selectedMonth] }} {{ selectedYear }}</h3>
            </div>

            <!-- Tabel Schedule Roster -->
            <table class="w-full border-collapse border border-slate-900 text-xs text-center">
              <thead>
                <tr class="bg-slate-100 font-bold uppercase border-b border-slate-900">
                  <th class="border border-slate-900 py-2 px-3">Hari / Tanggal</th>
                  <th class="border border-slate-900 py-2 px-3">Shift Pagi (08.00 – 20.00)</th>
                  <th class="border border-slate-900 py-2 px-3">Shift Malam (20.00 – 08.00)</th>
                  <th class="border border-slate-900 py-2 px-3">Keterangan (Libur/Cadangan)</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="day in daysInMonth" :key="day.dateStr" class="border-b border-slate-900">
                  <td class="border border-slate-900 py-1.5 px-3 text-left font-mono">{{ day.dayName }} / {{ day.formattedDate }}</td>
                  <td class="border border-slate-900 py-1.5 px-3 font-bold">
                    Tim {{ day.shift1 }}
                    <span v-if="day.replacementsShift1 && day.replacementsShift1.length > 0" class="block text-[10px] text-amber-800 font-normal">
                      (Pengganti: {{ day.replacementsShift1.map(r => r.replacerName.split(',')[0]).join(', ') }})
                    </span>
                  </td>
                  <td class="border border-slate-900 py-1.5 px-3 font-bold">
                    Tim {{ day.shift2 }}
                    <span v-if="day.replacementsShift2 && day.replacementsShift2.length > 0" class="block text-[10px] text-amber-800 font-normal">
                      (Pengganti: {{ day.replacementsShift2.map(r => r.replacerName.split(',')[0]).join(', ') }})
                    </span>
                  </td>
                  <td class="border border-slate-900 py-1.5 px-3">Tim {{ day.offTeams.join(', ') }}</td>
                </tr>
              </tbody>
            </table>

            <!-- Catatan Footer -->
            <div class="pt-2 text-xs space-y-1 font-medium">
              <div class="font-bold">Catatan :</div>
              <div>a. Untuk pergantian shift di wajibkan untuk melapor satu hari sebelum pergantian shift.</div>
              <div>b. Untuk absen alasan sakit wajib menyertakan surat keterangan sakit.</div>
              <div>c. Petugas yang libur dapat dijadikan <em>standby backup</em> jika ada yang berhalangan hadir atau sakit.</div>
            </div>
          </div>

          <div class="flex justify-end gap-3 border-t border-slate-200 pt-4">
            <button @click="showPrintModal = false" class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-800 font-bold text-xs rounded-xl border border-slate-300 cursor-pointer">Tutup</button>
            <button @click="exportToWordDoc" class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-extrabold text-xs rounded-xl shadow-md cursor-pointer border-0 flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">download</span>
              <span>Unduh Dokumen Word (.doc)</span>
            </button>
            <button @click="printDocument" class="px-6 py-2.5 bg-rose-700 hover:bg-rose-800 text-white font-extrabold text-xs rounded-xl shadow-md cursor-pointer border-0 flex items-center gap-2">
              <span class="material-symbols-outlined text-[18px]">print</span>
              <span>Cetak Langsung</span>
            </button>
          </div>
        </div>
      </div>
      <!-- Floating Save Icon Button (FAB) at Bottom Right -->
      <div class="fixed bottom-6 right-6 z-50">
        <button 
          @click="saveSchedules"
          :disabled="saving"
          class="w-14 h-14 bg-gradient-to-r from-emerald-600 to-teal-700 hover:from-emerald-700 hover:to-teal-800 text-white rounded-full shadow-2xl hover:shadow-[0_10px_25px_rgba(16,185,129,0.5)] transform hover:scale-105 active:scale-95 transition-all flex items-center justify-center cursor-pointer border border-emerald-400/40 disabled:opacity-50 disabled:cursor-not-allowed group backdrop-blur-md relative"
          title="Simpan Perubahan Jadwal Shift"
        >
          <span class="material-symbols-outlined text-[24px] group-hover:scale-110 transition-transform" :class="saving ? 'animate-spin' : ''">
            {{ saving ? 'sync' : 'save' }}
          </span>
          <span v-if="!saving" class="absolute top-1 right-1 w-3 h-3 rounded-full bg-emerald-300 animate-ping"></span>
        </button>
      </div>

    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import AdminLayout from '@/layouts/AdminLayout.vue'
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const monthNames = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']
const teamCodes = ['A', 'B', 'C', 'D', 'E']

const activeTab = ref<'roster' | 'requests'>('roster')
const now = new Date()
const selectedMonth = ref(now.getMonth()) // 0-indexed, dinamis sesuai bulan sekarang
const selectedYear = ref(now.getFullYear()) // dinamis sesuai tahun sekarang
const shiftMode = ref<number>(2) // 2 shift vs 3 shift

const alertMessage = ref('')
const saving = ref(false)
const showPrintModal = ref(false)
const showPreviewModal = ref(false)
const previewImageUrl = ref('')

// Helper for digits-only NIP comparisons
const cleanNipDigits = (nip?: string) => String(nip || '').replace(/\D/g, '')

// Dynamic List of Peserta Magangs
const masterOfficers = ref<any[]>([
  { name: 'A.Mappalua, S.Pd', nip: '19940503 202521 1 138', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Dinas Sosial' },
  { name: 'Suherman, S.Pd', nip: '19870304 202521 1 061', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Badan Penanggulangan Bencana Daerah' },
  { name: 'Riswandi Risman', nip: '20000206 202521 1 166', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Kesehatan' },
  { name: 'Abil Kizri', nip: '19900215 202521 1 114', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Perhubungan' },
  { name: 'Imam Ardiyansah', nip: '19911005 202521 1 087', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Satpol, Pemadam Kebakaran dan Penyelamatan' },
  { name: 'Abd.Rahim', nip: '19861130 202521 1 101', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Sosial' },
  { name: 'Munawir Syadzali,S.Pd', nip: '19860304 202521 1 147', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Badan Penanggulangan Bencana Daerah' },
  { name: 'Abdullah,S.Kep.,Ns.,M.Kep', nip: '19760802 200604 1 017', jabatan: 'PERENCANA', unit: 'Dinas Kesehatan' },
  { name: 'Ismail, S.Sos', nip: '19860712 202521 1 089', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Dinas Perhubungan' },
  { name: 'Aldi Afdali Saputra', nip: '19960328 202521 1 050', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Satpol, Pemadam Kebakaran dan Penyelamatan' }
])

// Load live active call takers from database
const loadLiveCallTakers = async () => {
  try {
    if (authStore && typeof authStore.fetchUsers === 'function') {
      await authStore.fetchUsers()
    }
    const allUsers = authStore?.usersList || []
    const activeCallTakers = allUsers.filter((u: any) => {
      const isCT = u.role === 'intern' || u.role === 'intern'
      const isActive = u.is_active !== false && u.status !== 'INACTIVE'
      return isCT && isActive
    })

    if (activeCallTakers.length > 0) {
      masterOfficers.value = activeCallTakers.map((u: any) => ({
        name: u.name,
        nip: u.nip || '',
        jabatan: u.jabatan || 'OPERATOR LAYANAN OPERASIONAL',
        unit: u.unit_kerja || 'Diskominfo Kab. Bulukumba'
      }))
    }
  } catch (err) {
    console.warn('[SchedulesView] Gagal memuat data live call takers:', err)
  }
}

// Check if a member NIP is currently non-active
const isMemberInactive = (nip?: string) => {
  if (!nip) return false
  const clean = cleanNipDigits(nip)
  if (!clean) return false
  return !masterOfficers.value.some(o => cleanNipDigits(o.nip) === clean)
}

// 5 Dynamic Teams matching ST
const teams = ref([
  {
    code: 'A',
    colorClass: 'bg-rose-50/80 border-rose-200 text-rose-950',
    members: [
      { name: 'A.Mappalua, S.Pd', nip: '19940503 202521 1 138', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Dinas Sosial' },
      { name: 'Suherman, S.Pd', nip: '19870304 202521 1 061', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Badan Penanggulangan Bencana Daerah' }
    ]
  },
  {
    code: 'B',
    colorClass: 'bg-amber-50/80 border-amber-200 text-amber-950',
    members: [
      { name: 'Riswandi Risman', nip: '20000206 202521 1 166', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Kesehatan' },
      { name: 'Abil Kizri', nip: '19900215 202521 1 114', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Perhubungan' }
    ]
  },
  {
    code: 'C',
    colorClass: 'bg-emerald-50/80 border-emerald-200 text-emerald-950',
    members: [
      { name: 'Imam Ardiyansah', nip: '19911005 202521 1 087', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Satpol, Pemadam Kebakaran dan Penyelamatan' },
      { name: 'Abd.Rahim', nip: '19861130 202521 1 101', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Dinas Sosial' }
    ]
  },
  {
    code: 'D',
    colorClass: 'bg-teal-50/80 border-teal-200 text-teal-950',
    members: [
      { name: 'Munawir Syadzali,S.Pd', nip: '19860304 202521 1 147', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Badan Penanggulangan Bencana Daerah' },
      { name: 'Abdullah,S.Kep.,Ns.,M.Kep', nip: '19760802 200604 1 017', jabatan: 'PERENCANA', unit: 'Dinas Kesehatan' }
    ]
  },
  {
    code: 'E',
    colorClass: 'bg-indigo-50/80 border-indigo-200 text-indigo-950',
    members: [
      { name: 'Ismail, S.Sos', nip: '19860712 202521 1 089', jabatan: 'PENATA LAYANAN OPERASIONAL', unit: 'Dinas Perhubungan' },
      { name: 'Aldi Afdali Saputra', nip: '19960328 202521 1 050', jabatan: 'OPERATOR LAYANAN OPERASIONAL', unit: 'Satpol, Pemadam Kebakaran dan Penyelamatan' }
    ]
  }
])

// Shuffle / Randomize Teams
const shuffleTeams = () => {
  const list = masterOfficers.value.length > 0 ? [...masterOfficers.value] : []
  const shuffled = list.sort(() => Math.random() - 0.5)
  teams.value.forEach((team, tIdx) => {
    if (shuffled[tIdx * 2]) team.members[0] = { ...shuffled[tIdx * 2] }
    if (shuffled[tIdx * 2 + 1]) team.members[1] = { ...shuffled[tIdx * 2 + 1] }
  })
  alertMessage.value = 'Pembagian tim berhasil diacak secara otomatis (2 orang per tim).'
}

// On Manual Dropdown Officer Select
const onMemberSelect = (teamCode: string, memberIdx: number, selectedNip: string) => {
  const team = teams.value.find(t => t.code === teamCode)
  if (!team) return

  // If empty nip selected, clear the slot
  if (!selectedNip) {
    team.members[memberIdx] = { name: '', nip: '', jabatan: '', unit: '' }
    return
  }

  const targetClean = cleanNipDigits(selectedNip)
  let targetOfficer = masterOfficers.value.find(o => cleanNipDigits(o.nip) === targetClean)
  if (!targetOfficer) {
    for (const t of teams.value) {
      const found = t.members.find(m => cleanNipDigits(m.nip) === targetClean)
      if (found) { targetOfficer = found; break }
    }
  }

  if (targetOfficer) {
    team.members[memberIdx] = { ...targetOfficer }
  }
}

interface PersonReplacement {
  replacedNip: string
  replacedName: string
  replacerNip: string
  replacerName: string
}

interface DaySchedule {
  dateStr: string
  dayName: string
  formattedDate: string
  isSunday: boolean
  shift1: string
  shift2: string
  shift3?: string
  replacementsShift1: PersonReplacement[]
  replacementsShift2: PersonReplacement[]
  replacementsShift3?: PersonReplacement[]
  offTeams: string[]
  teamsSnapshot?: any[]
}

const todayDateStr = computed(() => {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
})

const isDayPast = (dateStr: string) => {
  if (!dateStr) return false
  return dateStr < todayDateStr.value
}

const daysInMonth = ref<DaySchedule[]>([])

interface LeaveRequestItem {
  id: number
  created_at: string
  user_name: string
  user_nip: string
  category: 'Sakit' | 'Izin' | 'Tukar Shift'
  shift_date: string
  replacement_name?: string
  reason: string
  attachment_url?: string
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
}

const leaveRequests = ref<LeaveRequestItem[]>([
  {
    id: 101,
    created_at: '2026-07-31 08:30',
    user_name: 'A.Mappalua, S.Pd',
    user_nip: '19940503 202521 1 138',
    category: 'Sakit',
    shift_date: '01 Agustus 2026',
    replacement_name: 'Ismail, S.Sos',
    reason: 'Demam tinggi dan flu berat, wajib istirahat sesuai rujukan dokter RSUD Bulukumba.',
    attachment_url: 'https://images.unsplash.com/photo-1584515979956-d9f6e5d09982?w=500&auto=format&fit=crop&q=60',
    status: 'PENDING'
  }
])

const pendingRequestsCount = computed(() => {
  return leaveRequests.value.filter(r => r.status === 'PENDING').length
})

const getTeamNames = (code?: string, day?: any) => {
  if (!code) return ''
  let teamList = (day && day.teamsSnapshot && Array.isArray(day.teamsSnapshot) && day.teamsSnapshot.length > 0)
    ? day.teamsSnapshot
    : teams.value
  if (day && day.dateStr && day.dateStr < '2026-08-13') {
    teamList = JSON.parse(JSON.stringify(teamList))
    teamList.forEach((t: any) => {
      if (t.code === 'B' && Array.isArray(t.members) && t.members.length > 0) {
        const m0 = t.members[0]
        if (m0 && (cleanNipDigits(m0.nip) === '199107052025211081' || m0.name.includes('Aswar'))) {
          t.members[0] = {
            name: 'Riswandi Risman',
            nip: '20000206 202521 1 166',
            jabatan: 'OPERATOR LAYANAN OPERASIONAL',
            unit: 'Dinas Kesehatan'
          }
        }
      }
    })
  }
  const t = teamList.find((item: any) => item.code === code)
  if (!t) return ''
  return t.members.map((m: any) => m.name.split(',')[0]).join(' & ')
}

const getTeamMembers = (code?: string, day?: any) => {
  if (!code) return []
  let teamList = (day && day.teamsSnapshot && Array.isArray(day.teamsSnapshot) && day.teamsSnapshot.length > 0)
    ? day.teamsSnapshot
    : teams.value
  if (day && day.dateStr && day.dateStr < '2026-08-13') {
    teamList = JSON.parse(JSON.stringify(teamList))
    teamList.forEach((t: any) => {
      if (t.code === 'B' && Array.isArray(t.members) && t.members.length > 0) {
        const m0 = t.members[0]
        if (m0 && (cleanNipDigits(m0.nip) === '199107052025211081' || m0.name.includes('Aswar'))) {
          t.members[0] = {
            name: 'Riswandi Risman',
            nip: '20000206 202521 1 166',
            jabatan: 'OPERATOR LAYANAN OPERASIONAL',
            unit: 'Dinas Kesehatan'
          }
        }
      }
    })
  }
  const t = teamList.find((item: any) => item.code === code)
  return t ? t.members : []
}

// ─── Per-individual replacement helpers ───

// Get the replacer NIP for a given replacedNip in a replacement array
const getReplacerFor = (replacements?: PersonReplacement[], replacedNip?: string): string => {
  if (!replacements || !replacedNip) return ''
  const replacedClean = cleanNipDigits(replacedNip)
  const entry = replacements.find(r => cleanNipDigits(r.replacedNip) === replacedClean)
  return entry ? entry.replacerNip : ''
}

// Get the replacer name for display
const getReplacerNameFor = (replacements?: PersonReplacement[], replacedNip?: string): string => {
  if (!replacements || !replacedNip) return ''
  const replacedClean = cleanNipDigits(replacedNip)
  const entry = replacements.find(r => cleanNipDigits(r.replacedNip) === replacedClean)
  return entry ? entry.replacerName : ''
}

// Set or clear a replacer for a specific member in a shift
const setReplacerFor = (
  day: DaySchedule,
  shiftKey: 'replacementsShift1' | 'replacementsShift2' | 'replacementsShift3',
  replacedMember: { nip: string; name: string },
  replacerNip: string
) => {
  if (!day[shiftKey]) day[shiftKey] = []
  const replacedClean = cleanNipDigits(replacedMember.nip)
  // Remove existing entry for this replacedNip
  day[shiftKey] = day[shiftKey]!.filter(r => cleanNipDigits(r.replacedNip) !== replacedClean)
  if (!replacerNip) return // cleared
  // Find replacer info from masterOfficers or active teams
  const replacerClean = cleanNipDigits(replacerNip)
  let replacer = masterOfficers.value.find(o => cleanNipDigits(o.nip) === replacerClean)
  if (!replacer) {
    for (const t of teams.value) {
      const found = t.members.find(m => cleanNipDigits(m.nip) === replacerClean)
      if (found) { replacer = found; break }
    }
  }

  if (replacer) {
    day[shiftKey]!.push({
      replacedNip: replacedMember.nip,
      replacedName: replacedMember.name,
      replacerNip: replacer.nip,
      replacerName: replacer.name
    })
  }
}

// Returns officers NOT already assigned to another slot (excluding current slot)
const availableOfficersFor = (teamCode: string, memberIdx: number) => {
  const currentTeam = teams.value.find(t => t.code === teamCode)
  const currentMember = currentTeam?.members[memberIdx]

  // Collect all NIPs currently in use across all teams, except the current slot
  const usedNips = new Set<string>()
  teams.value.forEach(t => {
    t.members.forEach((m, idx) => {
      const isSelf = t.code === teamCode && idx === memberIdx
      if (!isSelf && m.nip) {
        const c = cleanNipDigits(m.nip)
        if (c) usedNips.add(c)
      }
    })
  })

  // Return only ACTIVE officers whose NIP is not yet used in another slot
  const list = masterOfficers.value.filter(o => {
    const c = cleanNipDigits(o.nip)
    return c && !usedNips.has(c)
  })

  // Preserve current assigned member if not in list
  if (currentMember && currentMember.nip) {
    const currClean = cleanNipDigits(currentMember.nip)
    if (currClean && !list.some(o => cleanNipDigits(o.nip) === currClean)) {
      list.unshift({ ...currentMember })
    }
  }

  // Deduplicate list by NIP
  const seen = new Set<string>()
  const uniqueList: any[] = []
  list.forEach(o => {
    const c = cleanNipDigits(o.nip)
    if (c && !seen.has(c)) {
      seen.add(c)
      uniqueList.push(o)
    } else if (!c && !uniqueList.includes(o)) {
      uniqueList.push(o)
    }
  })

  return uniqueList
}

const setShiftMode = (mode: number) => {
  shiftMode.value = mode
  generateAutoRotation()
}

const calculateOffTeams = (day: DaySchedule) => {
  const active = [day.shift1, day.shift2]
  if (shiftMode.value === 3 && day.shift3) {
    active.push(day.shift3)
  }
  day.offTeams = teamCodes.filter(c => !active.includes(c))
}

const generateAutoRotation = () => {
  const year = selectedYear.value
  const month = selectedMonth.value
  const numDays = new Date(year, month + 1, 0).getDate()

  const dayNamesIndo = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']

  const rotationPattern = [
    { shift1: 'A', shift2: 'B' },
    { shift1: 'C', shift2: 'D' },
    { shift1: 'E', shift2: 'A' },
    { shift1: 'B', shift2: 'C' },
    { shift1: 'D', shift2: 'E' }
  ]

  const rotationPattern3Shift = [
    { shift1: 'A', shift2: 'B', shift3: 'C' },
    { shift1: 'D', shift2: 'E', shift3: 'A' },
    { shift1: 'B', shift2: 'C', shift3: 'D' },
    { shift1: 'E', shift2: 'A', shift3: 'B' },
    { shift1: 'C', shift2: 'D', shift3: 'E' }
  ]

  const list: DaySchedule[] = []
  for (let d = 1; d <= numDays; d++) {
    const dateObj = new Date(year, month, d)
    const dayOfWeek = dateObj.getDay()
    const dayName = dayNamesIndo[dayOfWeek]
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const formattedDate = `${String(d).padStart(2, '0')} ${monthNames[month]} ${year}`

    const patIdx = (d - 1) % 5
    let s1 = 'A'
    let s2 = 'B'
    let s3 = 'C'

    if (shiftMode.value === 2) {
      s1 = rotationPattern[patIdx].shift1
      s2 = rotationPattern[patIdx].shift2
    } else {
      s1 = rotationPattern3Shift[patIdx].shift1
      s2 = rotationPattern3Shift[patIdx].shift2
      s3 = rotationPattern3Shift[patIdx].shift3
    }

    const daySchedule: DaySchedule = {
      dateStr,
      dayName,
      formattedDate,
      isSunday: dayOfWeek === 0,
      shift1: s1,
      shift2: s2,
      shift3: s3,
      replacementsShift1: [],
      replacementsShift2: [],
      replacementsShift3: [],
      offTeams: [],
      teamsSnapshot: JSON.parse(JSON.stringify(teams.value))
    }
    calculateOffTeams(daySchedule)
    list.push(daySchedule)
  }

  daysInMonth.value = list
  alertMessage.value = `Rotasi shift ${shiftMode.value} shift/hari untuk bulan ${monthNames[month]} ${year} berhasil di-generate.`
}

const onMonthYearChange = () => {
  generateAutoRotation()
}

const fetchSchedules = async () => {
  try {
    await loadLiveCallTakers()
    const userMap = new Map<string, string>()
    if (authStore?.usersList && Array.isArray(authStore.usersList)) {
      authStore.usersList.forEach((u: any) => {
        if (u.nip && u.name) {
          userMap.set(u.nip.replace(/\s+/g, ''), u.name)
        }
      })
    }

    const res = await axios.get('/api/admin/schedules')
    if (res.data && res.data.schedules && res.data.schedules.daysInMonth) {
      daysInMonth.value = res.data.schedules.daysInMonth
      shiftMode.value = res.data.schedules.shiftMode || 2
      selectedMonth.value = res.data.schedules.selectedMonth ?? 7
      selectedYear.value = res.data.schedules.selectedYear ?? 2026
      if (res.data.schedules.teams) {
        teams.value = res.data.schedules.teams
        teams.value.forEach((t: any) => {
          if (t.members && Array.isArray(t.members)) {
            t.members.forEach((m: any) => {
              if (m.nip) {
                const cleanM = cleanNipDigits(m.nip)
                const masterMatch = masterOfficers.value.find(o => cleanNipDigits(o.nip) === cleanM)
                if (masterMatch) {
                  m.nip = masterMatch.nip
                  m.name = masterMatch.name
                  m.jabatan = masterMatch.jabatan
                  m.unit = masterMatch.unit
                } else if (userMap.has(cleanM)) {
                  m.name = userMap.get(cleanM)
                }
              }
            })
          }
        })
      }
    } else {
      generateAutoRotation()
    }
  } catch (e) {
    generateAutoRotation()
  }
}

const fetchLeaveRequests = async () => {
  try {
    const res = await axios.get('/api/admin/leave-requests')
    if (res.data && res.data.requests && Array.isArray(res.data.requests) && res.data.requests.length > 0) {
      leaveRequests.value = res.data.requests
    }
  } catch (e) {
    // keep default seed
  }
}

const saveSchedules = async () => {
  saving.value = true
  try {
    const today = new Date()
    const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`

    // DYNAMIC ARCHITECTURE FOR ALL TEAMS & OFFICERS:
    // Freeze teamsSnapshot for past dates (< todayStr), only update for today & future dates (>= todayStr)
    daysInMonth.value.forEach((day: any) => {
      if (day.dateStr >= todayStr) {
        day.teamsSnapshot = JSON.parse(JSON.stringify(teams.value))
      } else if (!day.teamsSnapshot || !Array.isArray(day.teamsSnapshot) || day.teamsSnapshot.length === 0) {
        day.teamsSnapshot = JSON.parse(JSON.stringify(teams.value))
      }
    })

    const payload = {
      selectedMonth: selectedMonth.value,
      selectedYear: selectedYear.value,
      shiftMode: shiftMode.value,
      teams: teams.value,
      daysInMonth: daysInMonth.value
    }
    const res = await axios.put('/api/admin/schedules', payload)
    if (res.data && res.data.success) {
      alertMessage.value = res.data.message || 'Jadwal shift berhasil disimpan ke server!'
    }
  } catch (err: any) {
    alertMessage.value = err.response?.data?.error || 'Gagal menyimpan jadwal.'
  } finally {
    saving.value = false
  }
}

const saveLeaveRequestsToBackend = async () => {
  try {
    await axios.put('/api/admin/leave-requests', leaveRequests.value)
  } catch (e) {}
}

const approveRequest = (req: LeaveRequestItem) => {
  req.status = 'APPROVED'
  const targetDay = daysInMonth.value.find(d => d.formattedDate.toLowerCase().includes(req.shift_date.toLowerCase()) || d.dateStr.includes(req.shift_date))
  if (targetDay && req.replacement_name && req.user_nip) {
    // Tambahkan ke replacementsShift1 dengan format per-individu
    if (!targetDay.replacementsShift1) targetDay.replacementsShift1 = []
    // Cari data replacer dari masterOfficers
    const replacer = masterOfficers.value.find(o => o.name === req.replacement_name || o.nip === req.replacement_name)
    const replaced = masterOfficers.value.find(o => o.nip === req.user_nip) || { nip: req.user_nip, name: req.user_name }
    if (replacer) {
      targetDay.replacementsShift1 = targetDay.replacementsShift1.filter(r => r.replacedNip !== replaced.nip)
      targetDay.replacementsShift1.push({
        replacedNip: replaced.nip,
        replacedName: replaced.name,
        replacerNip: replacer.nip,
        replacerName: replacer.name
      })
    }
  }

  saveLeaveRequestsToBackend()
  saveSchedules()
  alertMessage.value = `Pengajuan ${req.category} atas nama ${req.user_name} TELAH DISETUJUI. Petugas Pengganti ${req.replacement_name || ''} otomatis ditempatkan di Roster.`
}

const rejectRequest = (req: LeaveRequestItem) => {
  req.status = 'REJECTED'
  saveLeaveRequestsToBackend()
  alertMessage.value = `Pengajuan ${req.category} atas nama ${req.user_name} ditolak.`
}

const previewImage = (url: string) => {
  previewImageUrl.value = url
  showPreviewModal.value = true
}

// Build Complete HTML String for Surat Tugas (without hardcoded number)
const buildSuratTugasHTML = () => {
  const monthStr = monthNames[selectedMonth.value]
  const yearStr = selectedYear.value
  
  let teamsRows = ''
  teams.value.forEach(team => {
    team.members.forEach((m, idx) => {
      teamsRows += `
        <tr style="border-bottom: 1px solid #000;">
          ${idx === 0 ? `<td rowspan="2" style="border: 1px solid #000; text-align: center; font-weight: bold; font-size: 13px; vertical-align: middle;">${team.code}</td>` : ''}
          <td style="border: 1px solid #000; padding: 6px; font-weight: bold;">${m.name} / NIP. ${m.nip}</td>
          <td style="border: 1px solid #000; padding: 6px; text-transform: uppercase;">${m.jabatan}</td>
          <td style="border: 1px solid #000; padding: 6px;">${m.unit}</td>
        </tr>
      `
    })
  })

  let rosterRows = ''
  daysInMonth.value.forEach(day => {
    const repl1Text = day.replacementsShift1 && day.replacementsShift1.length > 0
      ? ` <br><small style="color:#b45309;">(Pengganti: ${day.replacementsShift1.map(r => r.replacerName.split(',')[0]).join(', ')})</small>`
      : ''
    const repl2Text = day.replacementsShift2 && day.replacementsShift2.length > 0
      ? ` <br><small style="color:#b45309;">(Pengganti: ${day.replacementsShift2.map(r => r.replacerName.split(',')[0]).join(', ')})</small>`
      : ''
    const shift1Text = `Tim ${day.shift1}` + repl1Text
    const shift2Text = `Tim ${day.shift2}` + repl2Text
    
    rosterRows += `
      <tr style="border-bottom: 1px solid #000;">
        <td style="border: 1px solid #000; padding: 5px 8px; font-family: monospace; text-align: left;">${day.dayName} / ${day.formattedDate}</td>
        <td style="border: 1px solid #000; padding: 5px 8px; text-align: center; font-weight: bold;">${shift1Text}</td>
        <td style="border: 1px solid #000; padding: 5px 8px; text-align: center; font-weight: bold;">${shift2Text}</td>
        <td style="border: 1px solid #000; padding: 5px 8px; text-align: center;">Tim ${day.offTeams.join(', ')}</td>
      </tr>
    `
  })

  return `
    <div style="font-family: Arial, 'Times New Roman', sans-serif; font-size: 11px; color: #000; line-height: 1.4; padding: 20px;">
      <!-- Header Surat tanpa nomor hardcoded -->
      <div style="padding-bottom: 8px; border-bottom: 2px solid #000; margin-bottom: 15px;">
        <div style="font-weight: bold; text-transform: uppercase; font-size: 11px;">LAMPIRAN SURAT TUGAS</div>
        <div style="font-weight: bold;">Nomor : ............................................................</div>
        <div style="font-weight: bold;">Tanggal : 30 Juli 2026</div>
      </div>

      <!-- Judul -->
      <div style="text-align: center; margin-bottom: 15px;">
        <h2 style="margin: 0; font-size: 13px; font-weight: bold; text-transform: uppercase;">DAFTAR NAMA PETUGAS PESERTA MAGANG LAYANAN NTPD 112</h2>
        <h3 style="margin: 3px 0 0 0; font-size: 12px; font-weight: bold; text-transform: uppercase;">KAB. BULUKUMBA</h3>
      </div>

      <!-- Tabel Tim -->
      <table style="width: 100%; border-collapse: collapse; margin-bottom: 20px; font-size: 11px;">
        <thead>
          <tr style="background-color: #f1f5f9; text-align: center; font-weight: bold;">
            <th style="border: 1px solid #000; padding: 6px; width: 40px;">TIM</th>
            <th style="border: 1px solid #000; padding: 6px;">NAMA / NIP</th>
            <th style="border: 1px solid #000; padding: 6px;">JABATAN</th>
            <th style="border: 1px solid #000; padding: 6px;">UNIT KERJA</th>
          </tr>
        </thead>
        <tbody>
          ${teamsRows}
        </tbody>
      </table>

      <!-- Judul Jadwal -->
      <div style="text-align: center; margin-top: 20px; margin-bottom: 15px;">
        <h2 style="margin: 0; font-size: 13px; font-weight: bold; text-transform: uppercase;">JADWAL SHIFT PETUGAS PESERTA MAGANG LAYANAN NTPD 112</h2>
        <h3 style="margin: 3px 0 0 0; font-size: 12px; font-weight: bold; text-transform: uppercase;">KAB. BULUKUMBA BULAN ${monthStr.toUpperCase()} ${yearStr}</h3>
      </div>

      <!-- Tabel Jadwal -->
      <table style="width: 100%; border-collapse: collapse; margin-bottom: 15px; font-size: 11px;">
        <thead>
          <tr style="background-color: #f1f5f9; text-align: center; font-weight: bold;">
            <th style="border: 1px solid #000; padding: 6px;">Hari / Tanggal</th>
            <th style="border: 1px solid #000; padding: 6px;">Shift Pagi (08.00 – 20.00)</th>
            <th style="border: 1px solid #000; padding: 6px;">Shift Malam (20.00 – 08.00)</th>
            <th style="border: 1px solid #000; padding: 6px;">Keterangan (Libur/Cadangan)</th>
          </tr>
        </thead>
        <tbody>
          ${rosterRows}
        </tbody>
      </table>

      <!-- Catatan Footer -->
      <div style="margin-top: 15px; font-size: 11px;">
        <div style="font-weight: bold;">Catatan :</div>
        <div>a. Untuk pergantian shift di wajibkan untuk melapor satu hari sebelum pergantian shift.</div>
        <div>b. Untuk absen alasan sakit wajib menyertakan surat keterangan sakit.</div>
        <div>c. Petugas yang libur dapat dijadikan <em>standby backup</em> jika ada yang berhalangan hadir atau sakit.</div>
      </div>
    </div>
  `
}

const exportToWordDoc = () => {
  const contentHTML = buildSuratTugasHTML()
  const monthStr = monthNames[selectedMonth.value]
  const yearStr = selectedYear.value
  
  const header = `<html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word' xmlns='http://www.w3.org/TR/REC-html40'>
  <head>
    <meta charset='utf-8'>
    <title>Lampiran Surat Tugas NTPD 112</title>
  </head>
  <body>`
  
  const footer = `</body></html>`
  const fullDoc = header + contentHTML + footer
  
  const blob = new Blob(['\ufeff' + fullDoc], {
    type: 'application/msword'
  })
  
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `Surat_Tugas_CallTaker_112_${monthStr}_${yearStr}.doc`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  alertMessage.value = `Dokumen Word (Surat_Tugas_CallTaker_112_${monthStr}_${yearStr}.doc) berhasil diunduh.`
  showPrintModal.value = false
}

const printDocument = () => {
  const contentHTML = buildSuratTugasHTML()
  const win = window.open('', '_blank')
  if (win) {
    win.document.write(`
      <html>
        <head>
          <title>Lampiran Surat Tugas Shift Peserta Magang 112</title>
        </head>
        <body>
          ${contentHTML}
        </body>
      </html>
    `)
    win.document.close()
    win.focus()
    setTimeout(() => {
      win.print()
    }, 500)
  }
}

onMounted(() => {
  fetchSchedules()
  fetchLeaveRequests()
})
</script>

<style scoped>
.material-symbols-outlined { font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24; }
</style>
